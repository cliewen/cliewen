package release

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// errOffline is a transport failure with no server on the other end; errTimeout
// is one that timed out. They are distinct because the milestone names them
// separately, and identical in outcome because both mean "could not tell".
type errOffline struct{}

func (errOffline) Error() string { return "dial tcp: no such host" }

type errTimeout struct{}

func (errTimeout) Error() string   { return "context deadline exceeded" }
func (errTimeout) Timeout() bool   { return true }
func (errTimeout) Temporary() bool { return true }

// answering serves one canned reply to every request and counts how many it
// was asked for. The network is injected, so no test here reaches a live
// service and no verdict depends on which machine ran it (ADR-042).
type answering struct {
	status int
	body   string
	err    error
	calls  int
}

func (a *answering) RoundTrip(r *http.Request) (*http.Response, error) {
	a.calls++
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

func serving(status int, body string) *answering {
	return &answering{status: status, body: body}
}

func failing(err error) *answering { return &answering{err: err} }

func (a *answering) client() *http.Client { return &http.Client{Transport: a} }

// TestAC075_UnitPositive_EachPlatformGetsItsOwnRoute proves the branch the
// milestone names: the machine decides which single recipe is printed.
func TestAC075_UnitPositive_EachPlatformGetsItsOwnRoute(t *testing.T) {
	cases := []struct {
		platform Platform
		want     string
	}{
		{Platform{"windows", "amd64"}, "install.ps1"},
		{Platform{"windows", "arm64"}, "install.ps1"},
		{Platform{"darwin", "arm64"}, "install.sh"},
		{Platform{"linux", "amd64"}, "install.sh"},
		// No release publishes an asset for these, so the source route is the
		// only one that works — pinned to the release just reported.
		{Platform{"freebsd", "amd64"}, "go install github.com/cliewen/cliewen/cmd/clue@v0.12.0"},
		{Platform{"linux", "386"}, "go install github.com/cliewen/cliewen/cmd/clue@v0.12.0"},
	}
	for _, c := range cases {
		net := serving(http.StatusOK, `{"tag_name":"v0.12.0"}`)
		got := Check(Options{
			Current:  "0.11.2",
			Platform: c.platform,
			Client:   net.client(),
			CacheDir: t.TempDir(),
		})
		if !got.Known || !got.Behind || got.Latest != "0.12.0" {
			t.Fatalf("%v: expected a known newer release, got %+v", c.platform, got)
		}
		if len(got.Recipe) != 1 {
			t.Fatalf("%v: expected exactly one route, got %v", c.platform, got.Recipe)
		}
		if !strings.Contains(got.Recipe[0], c.want) {
			t.Fatalf("%v: expected %q, got %q", c.platform, c.want, got.Recipe[0])
		}
		// The two routes this machine is not on must not appear at all: a
		// reader who has to skip past wrong lines is what the branch avoids.
		for _, wrong := range []string{"install.ps1", "install.sh", "go install"} {
			if wrong != c.want && !strings.Contains(c.want, wrong) && strings.Contains(got.Recipe[0], wrong) {
				t.Fatalf("%v: route names another platform's recipe: %q", c.platform, got.Recipe[0])
			}
		}
	}
}

// TestAC075_UnitNegative_ACurrentOrSourceBuildIsNotBehind holds the other
// direction: nothing to catch up with means no recipe to print.
func TestAC075_UnitNegative_ACurrentOrSourceBuildIsNotBehind(t *testing.T) {
	for _, current := range []string{"0.12.0", "0.13.0", "dev", ""} {
		net := serving(http.StatusOK, `{"tag_name":"v0.12.0"}`)
		got := Check(Options{
			Current:  current,
			Platform: Platform{"linux", "amd64"},
			Client:   net.client(),
			CacheDir: t.TempDir(),
		})
		if !got.Known {
			t.Fatalf("%q: the release should still be known", current)
		}
		if got.Behind {
			t.Fatalf("%q: must not be reported behind %s", current, got.Latest)
		}
		if len(got.Recipe) != 0 {
			t.Fatalf("%q: expected no recipe, got %v", current, got.Recipe)
		}
	}
}

// TestAC077_UnitPositive_EveryDegradationIsUnknown covers all four the
// milestone names. Each means "could not tell", and that is one answer.
func TestAC077_UnitPositive_EveryDegradationIsUnknown(t *testing.T) {
	cases := map[string]*answering{
		"offline":             failing(errOffline{}),
		"timeout":             failing(errTimeout{}),
		"rate limit":          serving(http.StatusForbidden, `{"message":"API rate limit exceeded"}`),
		"too many requests":   serving(http.StatusTooManyRequests, ``),
		"unrecognized body":   serving(http.StatusOK, `<html>maintenance</html>`),
		"unrecognized fields": serving(http.StatusOK, `{"name":"0.12.0"}`),
		"unparseable tag":     serving(http.StatusOK, `{"tag_name":"nightly"}`),
	}
	for name, net := range cases {
		got := Check(Options{
			Current:  "0.11.2",
			Platform: Platform{"linux", "amd64"},
			Client:   net.client(),
			CacheDir: t.TempDir(),
		})
		if got.Known {
			t.Fatalf("%s: expected an unknown answer, got %+v", name, got)
		}
		if got.Behind || len(got.Recipe) != 0 {
			t.Fatalf("%s: an unknown answer must claim nothing, got %+v", name, got)
		}
	}
}

// TestAC077_UnitNegative_AnUnknownAnswerIsNeverCached keeps a failed request
// from silencing the next one: absence must not be stored as an answer.
func TestAC077_UnitNegative_AnUnknownAnswerIsNeverCached(t *testing.T) {
	dir := t.TempDir()
	down := failing(errOffline{})
	if got := Check(Options{Current: "0.11.2", Client: down.client(), CacheDir: dir}); got.Known {
		t.Fatalf("expected unknown, got %+v", got)
	}
	if entries, _ := os.ReadDir(dir); len(entries) != 0 {
		t.Fatalf("an unknown answer was written to the cache: %v", entries)
	}
	up := serving(http.StatusOK, `{"tag_name":"v0.12.0"}`)
	if got := Check(Options{Current: "0.11.2", Client: up.client(), CacheDir: dir}); !got.Known {
		t.Fatalf("the next call must ask again, got %+v", got)
	}
}

// TestAC078_UnitPositive_AFreshCacheCostsNoRequest is what makes running this
// at every session start acceptable to the host as well as to the user.
func TestAC078_UnitPositive_AFreshCacheCostsNoRequest(t *testing.T) {
	dir := t.TempDir()
	now := time.Date(2026, 8, 2, 12, 0, 0, 0, time.UTC)
	net := serving(http.StatusOK, `{"tag_name":"v0.12.0"}`)

	first := Check(Options{Current: "0.11.2", Client: net.client(), CacheDir: dir, Now: func() time.Time { return now }})
	if !first.Known || net.calls != 1 {
		t.Fatalf("expected one request and a known answer, got %+v after %d call(s)", first, net.calls)
	}
	// The cache is a file in the given directory, and nothing else was written.
	path := filepath.Join(dir, "latest-release.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("expected the answer cached at %s: %v", path, err)
	}
	var stored cached
	if err := json.Unmarshal(data, &stored); err != nil || stored.Tag != "v0.12.0" {
		t.Fatalf("unexpected cache contents %q (%v)", data, err)
	}

	later := now.Add(23 * time.Hour)
	second := Check(Options{Current: "0.11.2", Client: net.client(), CacheDir: dir, Now: func() time.Time { return later }})
	if net.calls != 1 {
		t.Fatalf("a fresh cache must cost no request, got %d call(s)", net.calls)
	}
	if !second.Known || second.Latest != "0.12.0" {
		t.Fatalf("expected the cached answer, got %+v", second)
	}
}

// TestAC078_UnitNegative_StaleAndBrokenCachesAreAbsence proves the three ways
// a cache stops counting. None of them is an error: a cache that can fail a
// command is worse than no cache.
func TestAC078_UnitNegative_StaleAndBrokenCachesAreAbsence(t *testing.T) {
	now := time.Date(2026, 8, 2, 12, 0, 0, 0, time.UTC)
	cases := map[string]string{
		"stale":           `{"tag":"v0.9.0","fetched":"2026-07-01T12:00:00Z"}`,
		"unreadable":      "\x00 not json at all",
		"nonsense":        `{"tag":"nightly","fetched":"2026-08-02T11:00:00Z"}`,
		"from the future": `{"tag":"v0.9.0","fetched":"2027-01-01T00:00:00Z"}`,
	}
	for name, contents := range cases {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "latest-release.json"), []byte(contents), 0o644); err != nil {
			t.Fatal(err)
		}
		net := serving(http.StatusOK, `{"tag_name":"v0.12.0"}`)
		got := Check(Options{Current: "0.11.2", Client: net.client(), CacheDir: dir, Now: func() time.Time { return now }})
		if net.calls != 1 {
			t.Fatalf("%s: expected the host to be asked, got %d call(s)", name, net.calls)
		}
		if !got.Known || got.Latest != "0.12.0" {
			t.Fatalf("%s: expected the fresh answer, got %+v", name, got)
		}
	}

	// A cache directory that cannot be written is absence too: the answer is
	// already in hand, so the only cost is asking again next time.
	blocked := filepath.Join(t.TempDir(), "file-where-a-directory-should-be")
	if err := os.WriteFile(blocked, []byte("not a directory"), 0o644); err != nil {
		t.Fatal(err)
	}
	net := serving(http.StatusOK, `{"tag_name":"v0.12.0"}`)
	got := Check(Options{Current: "0.11.2", Client: net.client(), CacheDir: filepath.Join(blocked, "cliewen")})
	if !got.Known || !got.Behind {
		t.Fatalf("an unwritable cache must not affect the answer, got %+v", got)
	}
}
