package main

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

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
		"rate limit":        {status: http.StatusForbidden, body: `{"message":"API rate limit exceeded"}`},
		"unrecognized body": {status: http.StatusOK, body: `<html>maintenance</html>`},
	}
	for name, net := range cases {
		opts := latestOptions(t, "0.11.2", release.Platform{OS: "linux", Arch: "amd64"}, net)
		code, out, errOut := runLatestCapturing(t, nil, opts)
		if code != 0 {
			t.Fatalf("%s: an unanswerable check is not a failure, got exit %d", name, code)
		}
		if errOut != "" {
			t.Fatalf("%s: nothing belongs on stderr, got %q", name, errOut)
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
