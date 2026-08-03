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
	"fmt"
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
	Current string // the running release, as given
	Latest  string // the newest published release, bare semver
	Known   bool   // false when the release list could not be read
	// SourceBuild is a build with no release to compare: the stamp is "dev"
	// or absent. Comparable is false for that and for a stamp this command
	// cannot read as a release. Both state a fact about the running binary
	// and are set whether or not the release list answered, so a caller
	// cannot read an outage as an unreadable stamp — keeping "up to date" and
	// "could not tell" apart is the whole reason this command exists.
	SourceBuild bool
	Comparable  bool
	Behind      bool // the running release is older than Latest
	// Ahead is a running stamp newer than anything published: an unreleased
	// build. It is kept apart from "current" because calling it the newest
	// release would state a fact about a version the release list never named.
	Ahead  bool
	Recipe []string // the installation route for this platform, one command per line
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

	// What the running binary is does not depend on whether the release list
	// answered, so it is settled before anything reaches the network.
	current, prerelease, comparable := parseCurrent(opts.Current)
	report := Report{
		Current:     opts.Current,
		SourceBuild: opts.Current == "" || opts.Current == "dev",
		Comparable:  comparable,
	}

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

	// The reported version is rebuilt from the numbers that were validated,
	// never echoed from the answer. What this command prints is a command the
	// user is told to run, so no byte of it may come from the network
	// unexamined — a tag is accepted only if it is exactly three numbers, and
	// then only those numbers are used.
	latest, ok := parseVersion(tag)
	if !ok {
		return report
	}
	report.Latest = fmt.Sprintf("%d.%d.%d", latest[0], latest[1], latest[2])
	report.Known = true
	// A pre-release orders before the release it names: 0.13.0-rc1 is behind
	// 0.13.0, and the moment the final ships is exactly when its user needs to
	// hear that.
	report.Behind = comparable && (older(current, latest) || (current == latest && prerelease))
	// A pre-release of an unpublished version — 0.13.0-rc1 while 0.12.0 is
	// newest — is ahead, not current: its numbers name nothing published.
	report.Ahead = comparable && !report.Behind && older(latest, current)
	if report.Behind {
		report.Recipe = recipe(opts.Platform, report.Latest)
	}
	return report
}

// AmbientBudget is the request budget for a check nobody asked for. It is
// shorter than the requested check's, because a notice that rides along with
// an ordinary command must never be something the user can feel; the day-long
// cache means it is paid at most once a day (PDR-023).
const AmbientBudget = time.Second

// QuietLine is the one-line form of a behind report: what `clue latest
// --quiet` prints, and what a workflow command's ambient notice prints. Both
// callers share it so there is one phrasing of this fact rather than two that
// drift.
func QuietLine(r Report) string {
	return fmt.Sprintf("clue %s is behind %s — run \"clue latest\" for the upgrade recipe", r.Current, r.Latest)
}

// Notice is the ambient advisory a workflow command prints when a newer
// release exists, and "" when there is nothing to say (PDR-023).
//
// It is the same check the command exposes, with a shorter budget and one
// caller-facing difference: only "behind" produces a line. Current, ahead, a
// source build, an unreadable stamp, and every way of failing to reach the
// release list are all silence — the caller did not ask, so anything other
// than the one fact worth interrupting for is noise.
//
// Whether a notice may be printed at all is the caller's gate, not this
// function's: reading the environment and which command ran belongs where
// those live.
func Notice(opts Options) string {
	if opts.Timeout == 0 {
		opts.Timeout = AmbientBudget
	}
	report := Check(opts)
	if !report.Known || !report.Behind {
		return ""
	}
	return QuietLine(report)
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

// parseCurrent reads the running stamp, and reports whether it named a
// pre-release. Unlike a tag from the network it is local, it never reaches the
// printed recipe, and it may legitimately carry a pre-release or build suffix
// — a release candidate is stamped 0.13.0-rc1 — so the suffix is trimmed
// before comparing. A "-" suffix is a pre-release and orders before the
// release whose numbers it carries; "+" is build metadata and orders with it.
// A source build reports "dev" and has no release to compare at all, the same
// exemption the drift rule makes.
func parseCurrent(v string) (numbers [3]int, prerelease, ok bool) {
	if v == "" || v == "dev" {
		return [3]int{}, false, false
	}
	if i := strings.IndexAny(v, "-+"); i >= 0 {
		prerelease = v[i] == '-'
		v = v[:i]
	}
	numbers, ok = parseVersion(v)
	return numbers, prerelease && ok, ok
}

// older compares two releases by their numbers.
func older(c, l [3]int) bool {
	for i := range c {
		if c[i] != l[i] {
			return c[i] < l[i]
		}
	}
	return false
}

// parseVersion reads exactly major.minor.patch, with an optional leading "v"
// and nothing else. A pre-release or build suffix is refused rather than
// trimmed: the published newest release never carries one, and accepting a
// suffix would mean the string that was checked and the string that is printed
// are not the same string. Anything it cannot read is not an error — it makes
// the answer unknown, and unknown reports nothing.
func parseVersion(v string) ([3]int, bool) {
	var out [3]int
	v = strings.TrimPrefix(v, "v")
	parts := strings.Split(v, ".")
	if len(parts) != 3 {
		return out, false
	}
	for i, p := range parts {
		// Atoi accepts a sign and Go's underscores; only decimal digits are a
		// version component, and an empty part is not one either.
		if p == "" || strings.ContainsFunc(p, func(r rune) bool { return r < '0' || r > '9' }) {
			return out, false
		}
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
