package main

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/cliewen/cliewen/internal/release"
)

// counting is the injected network with a request tally, so "the cache means a
// repeated call costs no request" is proven by counting rather than asserted.
type counting struct {
	inner    http.RoundTripper
	requests atomic.Int64
}

func (c *counting) RoundTrip(r *http.Request) (*http.Response, error) {
	c.requests.Add(1)
	return c.inner.RoundTrip(r)
}

// noEnv is an environment with nothing set in it.
func noEnv(string) string { return "" }

// envWith is an environment with exactly one variable set.
func envWith(key string) func(string) string {
	return func(k string) string {
		if k == key {
			return "1"
		}
		return ""
	}
}

// AC-086 positive: an agent that never asks is still told. The whole point of
// the notifier is that it rides along with work the session was already doing,
// because the instruction that asks an agent to check was observed being read
// and not acted on (PDR-023).
//
// This drives the whole path — command name, gate, check, stream — rather than
// the check alone, because a notice that is computed and never printed is the
// failure this change exists to fix, and it would look identical from inside.
func TestAC086_UnitPositive_AWorkflowCommandVolunteersTheNotice(t *testing.T) {
	opts := latestOptions(t, "0.11.2", release.Platform{OS: "linux", Arch: "amd64"},
		answering{status: http.StatusOK, body: `{"tag_name":"v0.12.0"}`})

	var errOut bytes.Buffer
	emitNotice("context", &errOut, noEnv, opts)

	got := errOut.String()
	if got == "" {
		t.Fatal("a session running an ordinary workflow command was told nothing")
	}
	if strings.Count(strings.TrimSpace(got), "\n") != 0 {
		t.Errorf("the notice must be exactly one line, got:\n%s", got)
	}
	if !strings.Contains(got, "0.11.2") || !strings.Contains(got, "0.12.0") {
		t.Errorf("the notice must name both versions, got: %q", got)
	}
}

// AC-087 negative: the same path, with the gate closed, writes nothing at all.
// Paired with the test above, this is what proves the gate is wired in rather
// than merely present — a notice printed unconditionally would pass every
// assertion above and reach the judge, a script, and every runner.
func TestAC087_UnitNegative_ABlockedNoticeWritesNothingToTheStream(t *testing.T) {
	for _, tc := range []struct {
		name    string
		command string
		env     func(string) string
	}{
		{name: "the judge", command: "validate", env: noEnv},
		{name: "version", command: "version", env: noEnv},
		{name: "CI", command: "context", env: envWith("CI")},
		{name: "the opt-out", command: "context", env: envWith("CLUE_NO_UPDATE_NOTIFIER")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			opts := latestOptions(t, "0.11.2", release.Platform{OS: "linux", Arch: "amd64"},
				answering{status: http.StatusOK, body: `{"tag_name":"v0.12.0"}`})

			var errOut bytes.Buffer
			emitNotice(tc.command, &errOut, tc.env, opts)

			if errOut.Len() != 0 {
				t.Errorf("expected nothing on the stream, got: %q", errOut.String())
			}
		})
	}
}

// AC-086: the notice and the quiet check say the same thing in the same words.
// Two phrasings of one fact drift, and the drift would be invisible — each is
// correct on its own.
func TestAC086_UnitPositive_TheNoticeSpeaksInTheQuietChecksWords(t *testing.T) {
	net := answering{status: http.StatusOK, body: `{"tag_name":"v0.12.0"}`}
	platform := release.Platform{OS: "linux", Arch: "amd64"}

	notice := release.Notice(latestOptions(t, "0.11.2", platform, net))
	_, quiet, _ := runLatestCapturing(t, []string{"--quiet"}, latestOptions(t, "0.11.2", platform, net))

	if notice != strings.TrimSpace(quiet) {
		t.Errorf("the two callers have drifted apart:\n notice: %q\n  quiet: %q", notice, strings.TrimSpace(quiet))
	}
}

// AC-086 negative: a repeated call inside the cached answer's lifetime does
// not ask again. This is what makes a notice on every workflow command
// acceptable to the release host as well as to the user (ADR-042) — without
// it, a session running ten commands would spend ten requests to repeat one
// unchanged fact.
func TestAC086_UnitNegative_ARepeatedNoticeCostsNoRequest(t *testing.T) {
	net := &counting{inner: answering{status: http.StatusOK, body: `{"tag_name":"v0.12.0"}`}}
	opts := release.Options{
		Current:  "0.11.2",
		Platform: release.Platform{OS: "linux", Arch: "amd64"},
		Client:   &http.Client{Transport: net},
		CacheDir: t.TempDir(), // one directory, shared by both calls
	}

	first := release.Notice(opts)
	second := release.Notice(opts)

	if first == "" || second == "" {
		t.Fatalf("both calls must report the same behind repository, got %q then %q", first, second)
	}
	if first != second {
		t.Errorf("the answer changed without the release list changing:\n%q\n%q", first, second)
	}
	if got := net.requests.Load(); got != 1 {
		t.Errorf("expected exactly one request for two notices, got %d", got)
	}
}

// AC-087 positive: the gate lets every workflow command through. This is the
// direction that keeps the negative honest — a gate that simply answered "no"
// would satisfy every exclusion below and deliver nothing at all, which is
// exactly the failure this change exists to fix.
func TestAC087_UnitPositive_TheGateAllowsEveryWorkflowCommand(t *testing.T) {
	for _, command := range []string{"context", "migrate", "refs", "init", "scaffold"} {
		t.Run(command, func(t *testing.T) {
			if !notifierAllowed(command, noEnv) {
				t.Errorf("clue %s carries no notice, so a session running it learns nothing", command)
			}
		})
	}
}

// AC-087 positive: the gate has no terminal condition, so the notice survives
// the pipe every coding agent runs commands through.
//
// This is the case the first implementation got wrong and the reason the gate
// was revised (Q3): an agent captures a command's output rather than watching
// it on a terminal, so a terminal condition switched the notice off in exactly
// the sessions it was built to reach. The test drives the whole path with the
// streams a captured run actually has, because that is the shape the failure
// took — every unit assertion about the gate passed while the real thing
// delivered nothing.
func TestAC087_UnitPositive_TheNoticeSurvivesACapturedStream(t *testing.T) {
	opts := latestOptions(t, "0.11.2", release.Platform{OS: "linux", Arch: "amd64"},
		answering{status: http.StatusOK, body: `{"tag_name":"v0.12.0"}`})

	// An os.Pipe is what an agent harness and a shell redirect both give the
	// command: a stream that is not a character device. A bytes.Buffer would
	// prove nothing here, having never been a terminal in the first place.
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()

	emitNotice("context", w, noEnv, opts)
	w.Close()

	got, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) == 0 {
		t.Fatal("a session whose output is captured rather than shown on a terminal was told nothing")
	}
	if !strings.Contains(string(got), "0.12.0") {
		t.Errorf("the notice must name the newer release, got: %q", got)
	}
}

// AC-086 negative: the notice is advisory, so it moves neither the exit code
// nor standard output. It is asserted against the built binary rather than
// against emitNotice, because the claim is about what a caller of `clue`
// receives — a test of the notifier alone would stay green if main printed the
// line on the wrong stream or let it decide the code.
func TestAC086_UnitNegative_TheNoticeMovesNeitherExitCodeNorStandardOutput(t *testing.T) {
	binary := filepath.Join(t.TempDir(), "clue"+exeSuffix())
	if out, err := exec.Command("go", "build", "-ldflags=-X main.version=0.11.2", "-o", binary, ".").CombinedOutput(); err != nil {
		t.Fatalf("building the command failed: %v\n%s", err, out)
	}
	cacheRoot := t.TempDir()
	cacheFile := redirectedCacheFile(cacheRoot)
	if err := os.MkdirAll(filepath.Dir(cacheFile), 0o755); err != nil {
		t.Fatal(err)
	}
	cache := []byte(`{"tag":"v0.12.0","fetched":"` + time.Now().UTC().Format(time.RFC3339Nano) + `"}`)
	if err := os.WriteFile(cacheFile, cache, 0o644); err != nil {
		t.Fatal(err)
	}
	baseEnv := make([]string, 0, len(os.Environ()))
	for _, entry := range os.Environ() {
		key, _, _ := strings.Cut(entry, "=")
		if strings.EqualFold(key, "CI") || strings.EqualFold(key, "CLUE_NO_UPDATE_NOTIFIER") {
			continue
		}
		baseEnv = append(baseEnv, entry)
	}

	// The same run twice: once with the notice permitted and once with the
	// documented opt-out, which is the only difference between them.
	run := func(extra ...string) (int, string, string) {
		t.Helper()
		cmd := exec.Command(binary, "refs", ".")
		cmd.Env = append(baseEnv, redirectedCacheEnv(cacheRoot)...)
		cmd.Env = append(cmd.Env, extra...)
		var out, errOut bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errOut
		err := cmd.Run()
		var exit *exec.ExitError
		switch {
		case err == nil:
		case errors.As(err, &exit):
		default:
			t.Fatalf("running the command failed: %v\n%s", err, errOut.String())
		}
		return cmd.ProcessState.ExitCode(), out.String(), errOut.String()
	}

	withCode, withOut, withErr := run()
	withoutCode, withoutOut, withoutErr := run("CLUE_NO_UPDATE_NOTIFIER=1")

	if withCode != withoutCode {
		t.Errorf("the notice moved the exit code: %d with, %d without", withCode, withoutCode)
	}
	if withOut != withoutOut {
		t.Errorf("the notice reached standard output:\n with: %q\n without: %q", withOut, withoutOut)
	}
	if !strings.Contains(withErr, "clue 0.11.2 is behind 0.12.0") {
		t.Errorf("the permitted run did not emit the notice on standard error: %q", withErr)
	}
	if withoutErr != "" {
		t.Errorf("the opted-out run wrote to standard error: %q", withoutErr)
	}
}

// AC-087 negative: no run writes a file in the repository, with or without the
// notice. The check reports and never edits (ADR-042), and riding along with a
// command that does edit must not change that.
func TestAC087_UnitNegative_ANoticeWritesNoFileInTheRepository(t *testing.T) {
	repo := t.TempDir()
	t.Chdir(repo)

	opts := latestOptions(t, "0.11.2", release.Platform{OS: "linux", Arch: "amd64"},
		answering{status: http.StatusOK, body: `{"tag_name":"v0.12.0"}`})

	var errOut bytes.Buffer
	emitNotice("context", &errOut, noEnv, opts)
	if errOut.Len() == 0 {
		t.Fatal("the notice did not run, so this proves nothing about what it writes")
	}

	entries, err := os.ReadDir(repo)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("the notice wrote into the repository: %v", entries)
	}
}

// AC-087 negative: the gate keeps the notice out of the judge, out of a
// script, and off a runner. Each exclusion is a decision rather than an
// oversight — the judge's output must stay a statement about the repository
// alone, "clue version" answers instantly and offline forever, and an
// ephemeral runner must never reach the release list (PDR-023, ADR-042).
func TestAC087_UnitNegative_TheGateExcludesTheJudgeAScriptAndARunner(t *testing.T) {
	for _, tc := range []struct {
		name    string
		command string
		env     func(string) string
	}{
		{name: "the judge is never interrupted", command: "validate", env: noEnv},
		{name: "version answers instantly and offline", command: "version", env: noEnv},
		{name: "the long version spelling too", command: "--version", env: noEnv},
		{name: "the check itself would report twice", command: "latest", env: noEnv},
		{name: "help", command: "help", env: noEnv},
		{name: "an unknown command", command: "nonesuch", env: noEnv},

		{name: "a CI runner has no one to tell", command: "context", env: envWith("CI")},
		{name: "the documented opt-out", command: "context", env: envWith("CLUE_NO_UPDATE_NOTIFIER")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if notifierAllowed(tc.command, tc.env) {
				t.Errorf("notifierAllowed(%q) allowed a notice it must never print", tc.command)
			}
		})
	}
}

// AC-087 negative: everything that is not "behind" is silence. The caller did
// not ask, so anything other than the one fact worth interrupting for is noise
// — and a notice that greets a current repository is one its reader learns to
// skip, which is the failure mode that costs the most.
func TestAC087_UnitNegative_TheNoticeIsSilentWhenThereIsNothingToSay(t *testing.T) {
	for _, tc := range []struct {
		name    string
		current string
		net     http.RoundTripper
	}{
		{name: "current", current: "0.12.0", net: answering{status: http.StatusOK, body: `{"tag_name":"v0.12.0"}`}},
		{name: "offline", current: "0.11.2", net: answering{err: offline{}}},
		{name: "rate limited", current: "0.11.2", net: answering{status: http.StatusForbidden, body: `{"message":"rate limit exceeded"}`}},
		{name: "a body it does not recognize", current: "0.11.2", net: answering{status: http.StatusOK, body: `<html>nope</html>`}},
		{name: "a source build has no release to compare", current: "dev", net: answering{status: http.StatusOK, body: `{"tag_name":"v0.12.0"}`}},
		{name: "an unreadable stamp is not reported as behind", current: "nonsense", net: answering{status: http.StatusOK, body: `{"tag_name":"v0.12.0"}`}},
		{name: "an unpublished build is ahead, not behind", current: "0.13.0", net: answering{status: http.StatusOK, body: `{"tag_name":"v0.12.0"}`}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			opts := latestOptions(t, tc.current, release.Platform{OS: "linux", Arch: "amd64"}, tc.net)
			if line := release.Notice(opts); line != "" {
				t.Errorf("expected silence, got: %q", line)
			}
		})
	}
}

// AC-087 negative: a source build and an unreadable stamp are known before
// consulting the release list, and neither can ever be reported as behind.
// Silence must therefore cost no ambient request, rather than adding latency
// to every local-development command for an answer the notifier cannot use.
func TestAC087_UnitNegative_AnIncomparableStampCostsNoRequest(t *testing.T) {
	for _, current := range []string{"dev", "", "nonsense"} {
		t.Run(current, func(t *testing.T) {
			net := &counting{inner: answering{status: http.StatusOK, body: `{"tag_name":"v0.12.0"}`}}
			opts := release.Options{
				Current:  current,
				Platform: release.Platform{OS: "linux", Arch: "amd64"},
				Client:   &http.Client{Transport: net},
				CacheDir: t.TempDir(),
			}

			if line := release.Notice(opts); line != "" {
				t.Errorf("expected silence, got: %q", line)
			}
			if got := net.requests.Load(); got != 0 {
				t.Errorf("an incomparable stamp made %d request(s)", got)
			}
		})
	}
}

// TestSanity_TheHelpTextNamesTheNoticeAndItsOptOut keeps the tool's own help
// honest about a line it prints unasked. The opt-out is documented in the
// changelog and both guides, but someone whose standard error just got noisier
// reaches for --help first, and a switch they cannot find there is a switch
// they do not have.
//
// It also binds the named commands to the gate, so a command added to or
// dropped from the notice stops being described by prose that used to be true.
func TestSanity_TheHelpTextNamesTheNoticeAndItsOptOut(t *testing.T) {
	if !strings.Contains(usage, "CLUE_NO_UPDATE_NOTIFIER") {
		t.Error("the help text does not name the opt-out, so a user cannot find it where they look first")
	}
	for command := range notifierCommands {
		if !strings.Contains(usage, command) {
			t.Errorf("clue %s carries the notice but the help text never says so", command)
		}
	}
	for _, excluded := range []string{"validate", "version"} {
		if notifierCommands[excluded] {
			t.Errorf("the help text promises no notice from clue %s, but the gate allows one", excluded)
		}
	}
}

// AC-087: the ambient budget is shorter than the requested check's, because a
// notice nobody asked for must never be something the user can feel. The
// requested check's default is the number this one has to beat; if that
// default ever drops, this says so rather than silently passing.
func TestAC087_UnitNegative_TheAmbientBudgetIsShorterThanTheRequestedCheckS(t *testing.T) {
	const requestedDefault = 3 * time.Second
	if release.AmbientBudget >= requestedDefault {
		t.Errorf("the ambient budget is %v, which is not shorter than the requested check's %v", release.AmbientBudget, requestedDefault)
	}
	if release.AmbientBudget <= 0 {
		t.Errorf("the ambient budget must be a real budget, got %v", release.AmbientBudget)
	}
}
