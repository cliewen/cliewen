package main

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/cliewen/cliewen/internal/release"
)

// answering is the injected network: the command's own network and platform
// are supplied by its caller, so what a user sees is proven without a live
// service and without three machines (ADR-042).
type answering struct {
	status int
	body   string
	err    error
}

func (a answering) RoundTrip(r *http.Request) (*http.Response, error) {
	if a.err != nil {
		return nil, a.err
	}
	return &http.Response{
		StatusCode: a.status,
		Body:       io.NopCloser(strings.NewReader(a.body)),
		Header:     make(http.Header),
		Request:    r,
	}, nil
}

type offline struct{}

func (offline) Error() string { return "dial tcp: no such host" }

// latestOptions builds a check that reads the injected network and an empty
// cache directory, so every run in these tests asks and answers deterministically.
func latestOptions(t *testing.T, current string, p release.Platform, rt http.RoundTripper) release.Options {
	t.Helper()
	return release.Options{
		Current:  current,
		Platform: p,
		Client:   &http.Client{Transport: rt},
		CacheDir: t.TempDir(),
	}
}

// runLatestCapturing drives the command and returns its exit code, stdout, and
// stderr. The criterion describes what a user sees, so a test that asserted a
// string it built itself would stay green after the print was deleted.
func runLatestCapturing(t *testing.T, args []string, opts release.Options) (int, string, string) {
	t.Helper()
	var out, errOut bytes.Buffer
	code := runLatest(args, &out, &errOut, opts)
	return code, out.String(), errOut.String()
}

// TestAC075_UnitPositive_LatestPrintsOneRouteAndTheMigrateSequence reads what
// the command actually printed on a Windows machine that is behind.
func TestAC075_UnitPositive_LatestPrintsOneRouteAndTheMigrateSequence(t *testing.T) {
	opts := latestOptions(t, "0.11.2", release.Platform{OS: "windows", Arch: "amd64"},
		answering{status: http.StatusOK, body: `{"tag_name":"v0.12.0"}`})
	code, out, errOut := runLatestCapturing(t, nil, opts)
	if code != 0 || errOut != "" {
		t.Fatalf("expected a clean exit, got %d with stderr %q", code, errOut)
	}
	for _, want := range []string{
		"clue 0.11.2", "0.12.0",
		"irm https://cliewen.dev/install.ps1 | iex",
		"clue migrate", "clue migrate --apply",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected %q in the report, got:\n%s", want, out)
		}
	}
	// The two routes this machine is not on must not appear at all.
	if strings.Contains(out, "install.sh") || strings.Contains(out, "go install") {
		t.Fatalf("a Windows run named another platform's route:\n%s", out)
	}
}

// TestAC075_UnitNegative_LatestReadsNoRepositoryAndWritesNothing holds the
// boundary between a reporting command and a background updater.
func TestAC075_UnitNegative_LatestReadsNoRepositoryAndWritesNothing(t *testing.T) {
	repo := t.TempDir()
	t.Chdir(repo)

	for _, args := range [][]string{nil, {"--quiet"}} {
		opts := latestOptions(t, "0.11.2", release.Platform{OS: "linux", Arch: "amd64"},
			answering{status: http.StatusOK, body: `{"tag_name":"v0.12.0"}`})
		if code, _, _ := runLatestCapturing(t, args, opts); code != 0 {
			t.Fatalf("%v: expected exit 0, got %d", args, code)
		}
		entries, err := os.ReadDir(repo)
		if err != nil {
			t.Fatal(err)
		}
		if len(entries) != 0 {
			t.Fatalf("%v: the command wrote into the repository: %v", args, entries)
		}
	}

	// A path would be a silent no-op rather than the scoping the other
	// commands give it, so it is a usage error rather than an ignored argument.
	opts := latestOptions(t, "0.11.2", release.Platform{OS: "linux", Arch: "amd64"},
		answering{status: http.StatusOK, body: `{"tag_name":"v0.12.0"}`})
	if code, _, errOut := runLatestCapturing(t, []string{repo}, opts); code != 2 || errOut == "" {
		t.Fatalf("expected a named usage error, got exit %d with stderr %q", code, errOut)
	}
}

// TestAC075_UnitPositive_TheTimeoutFlagReachesTheRequest holds the budget the
// usage text promises: the deadline the client applies is observable on the
// request, so removing the wiring cannot stay green.
func TestAC075_UnitPositive_TheTimeoutFlagReachesTheRequest(t *testing.T) {
	var budget time.Duration
	var seen bool
	rt := deadlineReading{observed: func(d time.Duration, ok bool) { budget, seen = d, ok }}
	opts := latestOptions(t, "0.11.2", release.Platform{OS: "linux", Arch: "amd64"}, rt)
	if code, _, _ := runLatestCapturing(t, []string{"--timeout=50ms"}, opts); code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if !seen {
		t.Fatal("the request carried no deadline — the flag never reached the client")
	}
	// Only the upper bound carries the proof: without the wiring the budget is
	// the 3s default. A lower bound would only measure the runner's load.
	if budget > 50*time.Millisecond {
		t.Fatalf("expected a budget of at most the 50ms asked for, got %v", budget)
	}
}

// deadlineReading reports the request's remaining budget, which is how the
// client's timeout becomes visible from inside a round trip.
type deadlineReading struct {
	observed func(time.Duration, bool)
}

func (d deadlineReading) RoundTrip(r *http.Request) (*http.Response, error) {
	if deadline, ok := r.Context().Deadline(); ok {
		d.observed(time.Until(deadline), true)
	} else {
		d.observed(0, false)
	}
	return answering{status: http.StatusOK, body: `{"tag_name":"v0.12.0"}`}.RoundTrip(r)
}

// TestAC075_UnitNegative_TheRealRunReportsOnTheRealMachine holds the one part
// of the wiring every other test injects past. A transposed field or a missing
// dispatch entry would leave the injected half fully proven and the shipped
// command reporting on a machine nobody is running.
func TestAC075_UnitNegative_TheRealRunReportsOnTheRealMachine(t *testing.T) {
	opts := runtimeCheck()
	if opts.Platform.OS != runtime.GOOS || opts.Platform.Arch != runtime.GOARCH {
		t.Fatalf("the check reports on %s/%s, not on the machine it runs on (%s/%s)",
			opts.Platform.OS, opts.Platform.Arch, runtime.GOOS, runtime.GOARCH)
	}
	if opts.Current != version {
		t.Fatalf("the check compares %q, not the running stamp %q", opts.Current, version)
	}
	// The command must also be reachable: the usage text promises it, and a
	// missing dispatch entry would make that promise a 404 in the terminal.
	if !strings.Contains(usage, "clue latest") {
		t.Fatal("the usage text does not name clue latest")
	}
	// Running the built command is the only way to prove the dispatch entry
	// exists. A 1ms budget guarantees the answer is "could not tell" whether
	// or not this machine has a network, which is the outcome being asserted:
	// silence and exit 0.
	//
	// It is built with the ordinary environment and only then run with a
	// redirected cache: the toolchain derives its own build cache from the
	// same directory, so redirecting before the build would recompile the
	// world into a temporary directory and delete it again.
	binary := filepath.Join(t.TempDir(), "clue"+exeSuffix())
	if out, err := exec.Command("go", "build", "-o", binary, ".").CombinedOutput(); err != nil {
		t.Fatalf("building the command failed: %v\n%s", err, out)
	}
	cmd := exec.Command(binary, "latest", "--quiet", "--timeout=1ms")
	cmd.Env = append(os.Environ(), redirectedCacheEnv(t.TempDir())...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("clue latest must exit 0 even when nothing answers: %v\n%s", err, out)
	}
	if len(out) != 0 {
		t.Fatalf("a quiet run that could not tell must print nothing, got:\n%s", out)
	}
}

// redirectedCacheEnv points the child's os.UserCacheDir at dir, per platform:
// darwin consults neither XDG_CACHE_HOME nor LocalAppData, so a run there
// would otherwise read and write the developer's own cache.
func redirectedCacheEnv(dir string) []string {
	switch runtime.GOOS {
	case "windows":
		return []string{"LocalAppData=" + dir}
	case "darwin":
		return []string{"HOME=" + dir}
	default:
		return []string{"XDG_CACHE_HOME=" + dir}
	}
}

func exeSuffix() string {
	if runtime.GOOS == "windows" {
		return ".exe"
	}
	return ""
}

// TestSanity_TheRoutesAreTheOnesActuallyPublished binds the recipe to its two
// sources of truth. The build matrix decides which machines have a prebuilt
// binary at all, and the install scripts are the commands the guide publishes;
// a route that drifts from either sends the user — by definition the one on the
// least-supported machine — to a command that fails in their terminal.
func TestSanity_TheRoutesAreTheOnesActuallyPublished(t *testing.T) {
	cfg := readGoreleaserConfig(t)
	if len(cfg.Builds) == 0 {
		t.Fatal("the release config declares no build")
	}
	published := cfg.Builds[0]

	// Every machine the release builds for must get a script route, and a
	// machine outside the matrix must get the source route.
	for _, goos := range published.Goos {
		for _, goarch := range published.Goarch {
			route := routeFor(t, release.Platform{OS: goos, Arch: goarch})
			if strings.Contains(route, "go install") {
				t.Errorf("%s/%s is a published release asset but the check offers the source route: %q", goos, goarch, route)
			}
		}
	}
	for _, p := range []release.Platform{{OS: "linux", Arch: "386"}, {OS: "freebsd", Arch: "amd64"}, {OS: "openbsd", Arch: "arm64"}} {
		if route := routeFor(t, p); !strings.Contains(route, "go install") {
			t.Errorf("%s/%s has no published asset but the check offers %q", p.OS, p.Arch, route)
		}
	}

	// The two script routes are the commands the guide publishes, character
	// for character: the check must not invent a third spelling.
	quickstart, err := os.ReadFile(filepath.Join("..", "..", "guide", "getting-started.md"))
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range []release.Platform{{OS: "windows", Arch: "amd64"}, {OS: "linux", Arch: "amd64"}} {
		route := routeFor(t, p)
		if !strings.Contains(string(quickstart), route) {
			t.Errorf("%s/%s is offered %q, which the quickstart does not publish", p.OS, p.Arch, route)
		}
	}
}

// routeFor is the single route the check would print on one machine.
func routeFor(t *testing.T, p release.Platform) string {
	t.Helper()
	got := release.Check(release.Options{
		Current:  "0.0.1",
		Platform: p,
		Client:   &http.Client{Transport: answering{status: http.StatusOK, body: `{"tag_name":"v0.12.0"}`}},
		CacheDir: t.TempDir(),
	})
	if len(got.Recipe) != 1 {
		t.Fatalf("%s/%s: expected exactly one route, got %v", p.OS, p.Arch, got.Recipe)
	}
	return got.Recipe[0]
}

// TestAC076_UnitPositive_QuietSaysOneLineWhenBehind is the shape a session
// start can afford.
func TestAC076_UnitPositive_QuietSaysOneLineWhenBehind(t *testing.T) {
	opts := latestOptions(t, "0.11.2", release.Platform{OS: "linux", Arch: "amd64"},
		answering{status: http.StatusOK, body: `{"tag_name":"v0.12.0"}`})
	code, out, errOut := runLatestCapturing(t, []string{"--quiet"}, opts)
	if code != 0 || errOut != "" {
		t.Fatalf("expected a clean exit, got %d with stderr %q", code, errOut)
	}
	if lines := strings.Count(strings.TrimSpace(out), "\n"); lines != 0 {
		t.Fatalf("expected exactly one line, got:\n%s", out)
	}
	if !strings.Contains(out, "0.11.2") || !strings.Contains(out, "0.12.0") {
		t.Fatalf("the one line must name both versions, got: %q", out)
	}
}

// TestAC076_UnitNegative_QuietIsSilentWhenThereIsNothingToSay covers the other
// direction: a check that greets a current repository every morning is a check
// its reader learns to ignore.
func TestAC076_UnitNegative_QuietIsSilentWhenThereIsNothingToSay(t *testing.T) {
	cases := map[string]answering{
		"current": {status: http.StatusOK, body: `{"tag_name":"v0.12.0"}`},
		"offline": {err: offline{}},
	}
	for name, net := range cases {
		current := "0.12.0"
		if name == "offline" {
			current = "0.11.2"
		}
		opts := latestOptions(t, current, release.Platform{OS: "linux", Arch: "amd64"}, net)
		code, out, errOut := runLatestCapturing(t, []string{"--quiet"}, opts)
		if code != 0 {
			t.Fatalf("%s: expected exit 0, got %d", name, code)
		}
		if out != "" || errOut != "" {
			t.Fatalf("%s: expected complete silence, got stdout %q stderr %q", name, out, errOut)
		}
	}
}

// TestAC077_UnitPositive_NotBeingAbleToTellExitsZeroAndSaysSo drives the
// command, so the exit code and the stream each word lands on are proven.
func TestAC077_UnitPositive_NotBeingAbleToTellExitsZeroAndSaysSo(t *testing.T) {
	cases := map[string]answering{
		"offline":           {err: offline{}},
		"timeout":           {err: timedOut{}},
		"rate limit":        {status: http.StatusForbidden, body: `{"message":"API rate limit exceeded"}`},
		"unrecognized body": {status: http.StatusOK, body: `<html>maintenance</html>`},
	}
	for name, net := range cases {
		for _, args := range [][]string{nil, {"--quiet"}} {
			opts := latestOptions(t, "0.11.2", release.Platform{OS: "linux", Arch: "amd64"}, net)
			code, out, errOut := runLatestCapturing(t, args, opts)
			if code != 0 {
				t.Fatalf("%s %v: an unanswerable check is not a failure, got exit %d", name, args, code)
			}
			if errOut != "" {
				t.Fatalf("%s %v: nothing belongs on stderr, got %q", name, args, errOut)
			}
			if len(args) > 0 {
				if out != "" {
					t.Fatalf("%s: the quiet run must print nothing at all, got %q", name, out)
				}
				continue
			}
			if !strings.Contains(out, "could not reach the release list") {
				t.Fatalf("%s: expected a calm report, got:\n%s", name, out)
			}
			// The one thing it must never do is claim the repository is current.
			if strings.Contains(out, "newest release") {
				t.Fatalf("%s: an unknown answer claimed the repository is current:\n%s", name, out)
			}
		}
	}
}

// timedOut is a request that exceeded its budget, which the command must read
// exactly as it reads an outage.
type timedOut struct{}

func (timedOut) Error() string { return "context deadline exceeded" }
func (timedOut) Timeout() bool { return true }

// TestAC077_UnitNegative_AnUnreadableStampIsNotReportedAsCurrent covers the
// case that looks like an answer and is not one: the host replied, and the
// running stamp is not a release this command can compare. Claiming currency
// there would collapse the very distinction the command exists to make.
func TestAC077_UnitNegative_AnUnreadableStampIsNotReportedAsCurrent(t *testing.T) {
	for _, current := range []string{"0.11", "0.11.2.3", "garbage"} {
		opts := latestOptions(t, current, release.Platform{OS: "linux", Arch: "amd64"},
			answering{status: http.StatusOK, body: `{"tag_name":"v0.12.0"}`})
		code, out, errOut := runLatestCapturing(t, nil, opts)
		if code != 0 || errOut != "" {
			t.Fatalf("%q: expected a clean exit, got %d with stderr %q", current, code, errOut)
		}
		if strings.Contains(out, "this is the newest release") {
			t.Fatalf("%q: a stamp that cannot be compared was reported as current:\n%s", current, out)
		}
		if !strings.Contains(out, "cannot be compared") || !strings.Contains(out, "0.12.0") {
			t.Fatalf("%q: expected the newest release named and the comparison declined, got:\n%s", current, out)
		}
	}
	// A pre-release stamp is comparable — it is trimmed, not refused — so the
	// user who is most likely to be behind is still told so.
	opts := latestOptions(t, "0.11.2-rc1", release.Platform{OS: "linux", Arch: "amd64"},
		answering{status: http.StatusOK, body: `{"tag_name":"v0.12.0"}`})
	_, out, _ := runLatestCapturing(t, nil, opts)
	if !strings.Contains(out, "0.12.0 is available") {
		t.Fatalf("a release candidate must still be told it is behind, got:\n%s", out)
	}
}

// TestAC077_UnitNegative_AKnownAnswerStillReports keeps the silence honest: it
// is the degradation that is quiet, not the command.
func TestAC077_UnitNegative_AKnownAnswerStillReports(t *testing.T) {
	opts := latestOptions(t, "0.12.0", release.Platform{OS: "darwin", Arch: "arm64"},
		answering{status: http.StatusOK, body: `{"tag_name":"v0.12.0"}`})
	code, out, errOut := runLatestCapturing(t, nil, opts)
	if code != 0 || errOut != "" {
		t.Fatalf("expected a clean exit, got %d with stderr %q", code, errOut)
	}
	if !strings.Contains(out, "this is the newest release") {
		t.Fatalf("expected the current report, got:\n%s", out)
	}
	if strings.Contains(out, "could not reach") {
		t.Fatalf("a known answer reported as unknown:\n%s", out)
	}
}
