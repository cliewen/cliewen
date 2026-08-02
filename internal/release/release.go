// Package release answers one question the repository cannot answer for
// itself: is there a newer clue than the one running?
//
// It reaches the network, so it lives outside the deterministic judge and
// never contributes to a validation verdict (ADR-042). It reports and writes
// nothing: no file in the repository with or without flags, and never the
// binary it is running as. Every way of failing to get an answer — offline, a
// timeout, a rate limit, a body it does not recognize — means "could not
// tell", which is not a defect in the repository and never an error here.
package release

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// latestURL is the published release list. Only the newest release's tag is
// read from it; nothing else about a release reaches this command.
const latestURL = "https://api.github.com/repos/cliewen/cliewen/releases/latest"

const (
	defaultTimeout = 3 * time.Second
	// defaultTTL bounds how long a cached answer stands. A release is cut at
	// most a few times a week and this check runs at every session start, so a
	// day is short enough to be useful and long enough that the ordinary cost
	// of running it is zero requests.
	defaultTTL = 24 * time.Hour
)

// Platform is the machine the command is running on, injected so all three
// installation routes are provable without three machines.
type Platform struct {
	OS   string // runtime.GOOS
	Arch string // runtime.GOARCH
}

// Options configures a check. Every field that reaches outside the process —
// the network, the clock, the filesystem, the machine's identity — is
// injectable, so no test depends on a live service or on which host ran it.
type Options struct {
	Current  string        // the running release, "dev" or "" when this is a source build
	Platform Platform      // nil-valued fields mean "no prebuilt asset", which is the go install route
	Client   *http.Client  // injected for tests; nil uses a default client
	Timeout  time.Duration // request budget; zero uses the default
	CacheDir string        // injected for tests; empty uses the user's cache directory
	TTL      time.Duration // how long a cached answer stands; zero uses the default
	Now      func() time.Time
}

// Report is what a check found. It never carries an error: a check that could
// not reach the release list reports Known == false and says so calmly.
type Report struct {
	Current string   // the running release, as given
	Latest  string   // the newest published release, bare semver
	Known   bool     // false when the release list could not be read
	Behind  bool     // the running release is older than Latest
	Recipe  []string // the installation route for this platform, one command per line
}

// cached is the stored answer. Only the tag and when it was fetched are kept:
// anything else would be a second copy of a fact that already has a home.
type cached struct {
	Tag     string    `json:"tag"`
	Fetched time.Time `json:"fetched"`
}

// Check reports whether a newer release exists. It never returns an error,
// because none of its failure modes says anything about the repository.
func Check(opts Options) Report {
	now := opts.Now
	if now == nil {
		now = time.Now
	}
	ttl := opts.TTL
	if ttl == 0 {
		ttl = defaultTTL
	}

	report := Report{Current: opts.Current}

	tag, ok := readCache(opts.CacheDir, now(), ttl)
	if !ok {
		tag, ok = fetchLatest(opts)
		if ok {
			// A cache that cannot be written is absence, not an error: the
			// answer is already in hand and the only cost is asking again.
			writeCache(opts.CacheDir, cached{Tag: tag, Fetched: now()})
		}
	}
	if !ok {
		return report
	}

	report.Latest = strings.TrimPrefix(tag, "v")
	report.Known = true
	report.Behind = olderThan(opts.Current, report.Latest)
	if report.Behind {
		report.Recipe = recipe(opts.Platform, report.Latest)
	}
	return report
}

// recipe is the installation route for one machine. Printing all three would
// make every reader skip past two wrong lines to find theirs (ADR-042).
func recipe(p Platform, version string) []string {
	switch {
	case p.OS == "windows" && prebuilt(p):
		return []string{`irm https://cliewen.dev/install.ps1 | iex`}
	case prebuilt(p):
		return []string{`curl -fsSL https://cliewen.dev/install.sh | sh`}
	default:
		// No published asset matches this machine, so the source route is the
		// only one that works. It is pinned rather than @latest: the point of
		// the line is to reach the release just named.
		return []string{"go install github.com/cliewen/cliewen/cmd/clue@v" + version}
	}
}

// prebuilt reports whether a release publishes a binary for this machine. The
// matrix is .goreleaser.yaml's, and it is a public contract (ADR-030): a
// platform outside it is not a failure, it is the go install route.
func prebuilt(p Platform) bool {
	switch p.OS {
	case "windows", "darwin", "linux":
	default:
		return false
	}
	return p.Arch == "amd64" || p.Arch == "arm64"
}

// olderThan compares two bare semver releases. A source build reports "dev"
// and has no release to compare, so it is never behind — the same exemption
// the drift rule makes, for the same reason: there is nothing to compare with.
func olderThan(current, latest string) bool {
	if current == "" || current == "dev" || latest == "" {
		return false
	}
	c, okC := parseVersion(current)
	l, okL := parseVersion(latest)
	if !okC || !okL {
		return false
	}
	for i := range c {
		if c[i] != l[i] {
			return c[i] < l[i]
		}
	}
	return false
}

// parseVersion reads major.minor.patch, discarding any pre-release or build
// suffix. A version it cannot read is not an error: it makes the comparison
// unknown, and unknown means nothing is reported.
func parseVersion(v string) ([3]int, bool) {
	var out [3]int
	v = strings.TrimPrefix(v, "v")
	if i := strings.IndexAny(v, "-+"); i >= 0 {
		v = v[:i]
	}
	parts := strings.Split(v, ".")
	if len(parts) != 3 {
		return out, false
	}
	for i, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil || n < 0 {
			return out, false
		}
		out[i] = n
	}
	return out, true
}

// fetchLatest asks the release list for the newest tag. Every failure —
// transport, status, and a body whose shape is not the one expected — returns
// the same answer: not known.
func fetchLatest(opts Options) (string, bool) {
	timeout := opts.Timeout
	if timeout == 0 {
		timeout = defaultTimeout
	}
	client := &http.Client{Timeout: timeout}
	if opts.Client != nil {
		c := *opts.Client
		client = &c
		if client.Timeout == 0 {
			client.Timeout = timeout
		}
	}
	req, err := http.NewRequest(http.MethodGet, latestURL, nil)
	if err != nil {
		return "", false
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	resp, err := client.Do(req)
	if err != nil {
		return "", false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		// A rate limit arrives here as 403 or 429, and is read exactly like an
		// outage: the command could not tell, and says so.
		return "", false
	}
	// The body is bounded: this command must cost a session start almost
	// nothing, and an unexpected response is already handled as unknown.
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return "", false
	}
	var payload struct {
		Tag string `json:"tag_name"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", false
	}
	if _, ok := parseVersion(payload.Tag); !ok {
		return "", false
	}
	return payload.Tag, true
}

// cachePath is where the answer is kept: machine state, in the user's cache
// directory, never in the repository (ADR-042). An empty path means the
// directory could not be resolved, and the check runs without a cache.
func cachePath(dir string) string {
	if dir == "" {
		d, err := os.UserCacheDir()
		if err != nil {
			return ""
		}
		dir = filepath.Join(d, "cliewen")
	}
	return filepath.Join(dir, "latest-release.json")
}

// readCache returns a stored answer that is still inside its lifetime.
// Missing, unreadable, unparseable, and stale are one answer — absence.
func readCache(dir string, now time.Time, ttl time.Duration) (string, bool) {
	path := cachePath(dir)
	if path == "" {
		return "", false
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", false
	}
	var c cached
	if err := json.Unmarshal(data, &c); err != nil {
		return "", false
	}
	if _, ok := parseVersion(c.Tag); !ok {
		return "", false
	}
	// A cache stamped in the future is as untrustworthy as a stale one — a
	// clock moved backwards would otherwise freeze the answer indefinitely.
	age := now.Sub(c.Fetched)
	if age < 0 || age >= ttl {
		return "", false
	}
	return c.Tag, true
}

// writeCache stores an answer, and ignores every reason it could not. The
// caller already has the answer; a cache that can fail a command is worse
// than no cache.
func writeCache(dir string, c cached) {
	path := cachePath(dir)
	if path == "" {
		return
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return
	}
	data, err := json.Marshal(c)
	if err != nil {
		return
	}
	_ = os.WriteFile(path, data, 0o644)
}
