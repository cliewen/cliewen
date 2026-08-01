// Package refs resolves the external addresses a corpus points at.
//
// ADR-040 splits the check in two. The judge enforces that an external
// reference names its target, offline and deterministically. Following that
// target needs the network, so it lives here — outside `clue validate`, unable
// to reach a validation verdict, and never a required status check. A verdict
// that depended on another system's present state would stop being
// reproducible against a pinned revision and would let someone else's outage
// condemn a corpus that has not changed.
package refs

import (
	"errors"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

// Outcome is how a target answered. ADR-040 fixes these five and their
// meanings; only Gone is an error.
type Outcome string

const (
	// Reachable — the target answered and is there.
	Reachable Outcome = "reachable"
	// Restricted — the target is not readable from here. Two answers prove
	// that. A 401 or 403 says so outright. A 404 whose owner root still
	// answers says it obliquely: GitHub returns 404 rather than 403 for a
	// private repository, deliberately, so that a refusal cannot confirm the
	// repository exists. Taking that 404 at face value would report a correct
	// reference as rot, so the owner probe decides between the two.
	//
	// The trade-off runs in the safe direction. A genuinely deleted repository
	// whose owner still exists is reported restricted rather than gone, so
	// this can miss real rot — but it can never fail a corpus that is right,
	// which is the only failure ADR-040 treats as an error.
	Restricted Outcome = "restricted"
	// Redirected — the target moved. The new location is offered as a rewrite.
	Redirected Outcome = "redirected"
	// Gone — 404 or 410. The reference points at nothing. The one error.
	Gone Outcome = "gone"
	// Unreachable — a timeout, a DNS failure, a rate limit. Reported as
	// unknown, never as invalid.
	Unreachable Outcome = "unreachable"
)

// Reference is one external address found in the corpus.
type Reference struct {
	Path   string // repo-relative file, forward slashes
	Line   int    // 1-based line in that file
	URL    string
	Result Outcome
	Target string // the new location, when Redirected
	Detail string // why, when Unreachable
	Frozen bool   // in pinned history: reported, never rewritten
}

// Report is a resolved corpus.
type Report struct {
	References []Reference
	Applied    bool
}

// Counts returns how many references landed in each outcome.
func (r Report) Counts() map[Outcome]int {
	counts := map[Outcome]int{}
	for _, ref := range r.References {
		counts[ref.Result]++
	}
	return counts
}

// HasErrors reports whether any reference points at nothing. Restricted and
// Unreachable are deliberately excluded: neither says the corpus is wrong.
func (r Report) HasErrors() bool {
	for _, ref := range r.References {
		if ref.Result == Gone {
			return true
		}
	}
	return false
}

// Options configures a resolution run.
type Options struct {
	Apply   bool          // rewrite redirected references in place
	Timeout time.Duration // per-request budget; zero uses the default
	Workers int           // concurrent requests; zero uses the default
	Client  *http.Client  // injected for tests; nil uses a default client
}

const (
	defaultTimeout = 10 * time.Second
	defaultWorkers = 8
)

// Resolve harvests every external address under root and follows it.
func Resolve(root string, opts Options) (Report, error) {
	refs, err := harvest(root)
	if err != nil {
		return Report{}, err
	}
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
	// Not following redirects is the behaviour, not a default: a followed hop
	// reports the destination as reachable and lets the stale reference
	// survive unnoticed. An injected client gets the same policy, so no caller
	// — including a test — can quietly turn the finding off.
	client.CheckRedirect = func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	}
	workers := opts.Workers
	if workers == 0 {
		workers = defaultWorkers
	}

	// Resolve each distinct address once: a corpus cites the same target from
	// several files, and asking repeatedly is rude to the host and slower for
	// no added truth.
	distinct := map[string]int{}
	for _, r := range refs {
		distinct[r.URL] = 0
	}
	type answer struct {
		result Outcome
		target string
		detail string
	}
	answers := map[string]answer{}
	var mu sync.Mutex
	var wg sync.WaitGroup
	queue := make(chan string)
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for u := range queue {
				res, target, detail := probe(client, u)
				mu.Lock()
				answers[u] = answer{res, target, detail}
				mu.Unlock()
			}
		}()
	}
	for u := range distinct {
		queue <- u
	}
	close(queue)
	wg.Wait()

	for i := range refs {
		a := answers[refs[i].URL]
		refs[i].Result = a.result
		refs[i].Target = a.target
		refs[i].Detail = a.detail
	}
	sort.Slice(refs, func(i, j int) bool {
		if refs[i].Path != refs[j].Path {
			return refs[i].Path < refs[j].Path
		}
		if refs[i].Line != refs[j].Line {
			return refs[i].Line < refs[j].Line
		}
		return refs[i].URL < refs[j].URL
	})

	report := Report{References: refs}
	if opts.Apply {
		if err := rewrite(root, refs); err != nil {
			return report, err
		}
		report.Applied = true
	}
	return report, nil
}

// probe asks for one address and classifies the answer.
func probe(client *http.Client, raw string) (Outcome, string, string) {
	// HEAD first: the corpus needs the status, not the body. Hosts that refuse
	// HEAD get a GET, because a 405 says nothing about whether the target
	// exists.
	res, target, detail, ok := request(client, http.MethodHead, raw)
	if ok && res == Gone {
		// A host that answers 404 to HEAD but serves the page on GET is common
		// enough that calling it gone without asking properly would be wrong.
		if r2, t2, d2, ok2 := request(client, http.MethodGet, raw); ok2 {
			res, target, detail = r2, t2, d2
		}
	}
	if !ok {
		r2, t2, d2, ok2 := request(client, http.MethodGet, raw)
		if !ok2 {
			return Unreachable, "", detail
		}
		res, target, detail = r2, t2, d2
	}
	if res == Gone {
		if owner, live := ownerAnswers(client, raw); live {
			return Restricted, "", "owner " + owner + " answers but the target does not; likely private rather than deleted"
		}
	}
	return res, target, detail
}

// ownerAnswers reports whether the address's owner root still exists. It is
// the evidence that separates a private target from a deleted one on a host
// that refuses to distinguish them itself.
func ownerAnswers(client *http.Client, raw string) (string, bool) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", false
	}
	segments := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(segments) < 2 || segments[0] == "" {
		// With no owner segment there is nothing narrower than the host to
		// ask about, and "the host is up" says nothing about the target.
		return "", false
	}
	root := &url.URL{Scheme: u.Scheme, Host: u.Host, Path: "/" + segments[0]}
	res, _, _, ok := request(client, http.MethodHead, root.String())
	if !ok {
		return "", false
	}
	return root.String(), res == Reachable || res == Redirected || res == Restricted
}

// hostToken returns the credential to send to host, and only to that host. A
// token is scoped to the service that issued it, so attaching it to an
// arbitrary wiki or documentation host in the corpus would hand that host a
// GitHub credential it was never meant to see.
func hostToken(host string) string {
	if host != "github.com" && host != "api.github.com" && !strings.HasSuffix(host, ".github.com") {
		return ""
	}
	for _, name := range []string{"GITHUB_TOKEN", "GH_TOKEN"} {
		if v := strings.TrimSpace(os.Getenv(name)); v != "" {
			return v
		}
	}
	return ""
}

func request(client *http.Client, method, raw string) (Outcome, string, string, bool) {
	req, err := http.NewRequest(method, raw, nil)
	if err != nil {
		return Unreachable, "", err.Error(), false
	}
	// A credential turns a private target from an oblique 404 into a plain
	// answer — but only where the credential is honoured. A bearer token is
	// inert on github.com, whose web pages authenticate by session cookie, so
	// it is worth sending only when the address has an API equivalent to ask
	// instead. Where it has none, the plain web request is made with no
	// credential and no API media type: attaching either would hand a token to
	// a request that cannot use it, and the JSON media type makes the web host
	// answer 406 to a page that is perfectly fine.
	//
	// The corpus still stores the plain address. This is how the resolver
	// asks, not what the corpus says.
	if tok := hostToken(req.URL.Host); tok != "" {
		if api := apiEquivalent(req.URL); api != "" {
			if r, err := http.NewRequest(method, api, nil); err == nil {
				req = r
				req.Header.Set("Authorization", "Bearer "+tok)
				req.Header.Set("Accept", "application/vnd.github+json")
			}
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		var uerr *url.Error
		detail := err.Error()
		if errors.As(err, &uerr) && uerr.Timeout() {
			detail = "timeout"
		}
		return Unreachable, "", detail, false
	}
	defer resp.Body.Close()
	switch {
	case resp.StatusCode == http.StatusUnauthorized, resp.StatusCode == http.StatusForbidden:
		// A throttled forge answers 403 too. Calling that "restricted" would
		// report every affected address as probably-private and destroy the
		// signal separating the two outcomes exists to preserve.
		if resp.Header.Get("Retry-After") != "" || resp.Header.Get("X-RateLimit-Remaining") == "0" {
			return Unreachable, "", "rate limited", true
		}
		return Restricted, "", "", true
	case resp.StatusCode == http.StatusNotFound, resp.StatusCode == http.StatusGone:
		return Gone, "", "", true
	case resp.StatusCode >= 300 && resp.StatusCode < 400:
		loc := resp.Header.Get("Location")
		if loc == "" {
			return Unreachable, "", "redirect without a location", true
		}
		if abs, err := resolveRelative(raw, loc); err == nil {
			loc = abs
		}
		return Redirected, loc, "", true
	case resp.StatusCode == http.StatusTooManyRequests:
		return Unreachable, "", "rate limited", true
	case resp.StatusCode >= 200 && resp.StatusCode < 300:
		return Reachable, "", "", true
	default:
		return Unreachable, "", resp.Status, true
	}
}

func resolveRelative(base, loc string) (string, error) {
	b, err := url.Parse(base)
	if err != nil {
		return "", err
	}
	l, err := url.Parse(loc)
	if err != nil {
		return "", err
	}
	return b.ResolveReference(l).String(), nil
}

var (
	fenceRe    = regexp.MustCompile("^\\s*(```|~~~)")
	mdTargetRe = regexp.MustCompile(`\]\((https?://[^)\s]+)\)`)
	autolinkRe = regexp.MustCompile(`<(https?://[^>\s]+)>`)
	bareURLRe  = regexp.MustCompile(`(^|[\s(])(https?://[^\s)\]<>]+)`)
)

// harvest collects every http(s) address in corpus prose, under docs/ and
// changes/ alike, so the resolver answers about the same corpus the judge
// validated.
//
// A clue: identity is deliberately not collected. It names a corpus artifact
// in another repository, and following it would mean knowing where that
// repository currently lives — the address coupling the scheme exists to
// avoid.
func harvest(root string) ([]Reference, error) {
	var refs []Reference
	var walkErr error
	for _, dir := range []string{"docs", "changes"} {
		if err := harvestTree(root, filepath.Join(root, dir), &refs); err != nil {
			walkErr = err
		}
	}
	return refs, walkErr
}

// harvestTree collects addresses under one corpus tree. Both docs/ and
// changes/ are read, because the judge validates the form in both and a
// resolver that saw less would answer about a different corpus.
func harvestTree(root, tree string, refs *[]Reference) error {
	err := filepath.Walk(tree, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		frozen := isFrozen(string(data))
		inFence := false
		for i, line := range strings.Split(string(data), "\n") {
			if fenceRe.MatchString(line) {
				inFence = !inFence
				continue
			}
			if inFence {
				continue
			}
			for _, u := range addressesIn(line) {
				*refs = append(*refs, Reference{Path: rel, Line: i + 1, URL: u, Frozen: frozen})
			}
		}
		return nil
	})
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func addressesIn(line string) []string {
	seen := map[string]bool{}
	var out []string
	add := func(u string) {
		u = strings.TrimRight(u, ".,;:")
		if u == "" || seen[u] {
			return
		}
		seen[u] = true
		out = append(out, u)
	}
	for _, m := range mdTargetRe.FindAllStringSubmatch(line, -1) {
		add(m[1])
	}
	for _, m := range autolinkRe.FindAllStringSubmatch(line, -1) {
		add(m[1])
	}
	// Bare URLs last: a URL already captured as a link target must not be
	// counted twice.
	for _, m := range bareURLRe.FindAllStringSubmatch(line, -1) {
		add(m[2])
	}
	return out
}

// rewrite replaces every redirected address with the location its host gave.
// Only redirects are written: a gone reference has no correct replacement to
// offer, and inventing one would be guessing.
//
// Frozen history is never rewritten. A completed plan records what was
// observed at the time — "the guide at this address returned HTTP 200" is a
// statement about that address on that day, and pointing it at wherever the
// content lives now would falsify the observation rather than repair it. Those
// references are reported and left for a human, who can decide whether the
// sentence is about the address or merely uses it.
func rewrite(root string, refs []Reference) error {
	byFile := map[string][]Reference{}
	for _, r := range refs {
		if r.Result == Redirected && r.Target != "" && !r.Frozen {
			byFile[r.Path] = append(byFile[r.Path], r)
		}
	}
	for rel, list := range byFile {
		full := filepath.Join(root, filepath.FromSlash(rel))
		data, err := os.ReadFile(full)
		if err != nil {
			return err
		}
		// Rewriting by plain substring replacement corrupts two things. A
		// renamed repository's address is a prefix of every address below it,
		// so replacing it would rewrite a longer address that was never
		// classified. And a fenced block is deliberately not harvested, so
		// editing inside one would change an example or a quoted log the run
		// never looked at. Replacing on the classified lines only, matched to
		// the end of the address, avoids both.
		lines := strings.Split(string(data), "\n")
		byLine := map[int][]Reference{}
		for _, r := range list {
			byLine[r.Line] = append(byLine[r.Line], r)
		}
		for n, refsOnLine := range byLine {
			if n < 1 || n > len(lines) {
				continue
			}
			// Longest first, so a prefix never consumes a longer sibling.
			sort.Slice(refsOnLine, func(i, j int) bool {
				return len(refsOnLine[i].URL) > len(refsOnLine[j].URL)
			})
			for _, r := range refsOnLine {
				lines[n-1] = replaceWholeAddress(lines[n-1], r.URL, r.Target)
			}
		}
		if err := os.WriteFile(full, []byte(strings.Join(lines, "\n")), 0o644); err != nil {
			return err
		}
	}
	return nil
}

// replaceWholeAddress swaps an address only where it ends there, so a shorter
// address is never substituted into a longer one that shares its prefix.
func replaceWholeAddress(line, from, to string) string {
	var b strings.Builder
	for {
		i := strings.Index(line, from)
		if i < 0 {
			b.WriteString(line)
			return b.String()
		}
		end := i + len(from)
		if end < len(line) && !isAddressBoundary(line[end]) {
			b.WriteString(line[:end])
			line = line[end:]
			continue
		}
		b.WriteString(line[:i])
		b.WriteString(to)
		line = line[end:]
	}
}

func isAddressBoundary(b byte) bool {
	// Sentence punctuation belongs here because harvesting strips it: an
	// address ending a sentence is stored without the stop, so a boundary set
	// that omitted it would skip the rewrite and still report the run as
	// applied — telling the user a repair happened that did not. A carriage
	// return is the same hazard on a checkout without LF normalisation.
	switch b {
	case ' ', '\t', '\r', '\n', ')', ']', '>', ',', ';', ':', '.', '!', '?', '"', '\'', '`', '*', '|':
		return true
	}
	return false
}

// apiEquivalent maps a github.com web address to the API address that answers
// the same question for a credential holder. It returns "" for a shape it does
// not recognise, which leaves the web request — and the owner probe behind it
// — exactly as they were.
//
// This is forge-specific knowledge, and it is deliberately confined here. What
// the corpus records stays a plain address with no forge assumed; only the act
// of asking knows that GitHub answers bearer tokens on one host and not the
// other.
func apiEquivalent(u *url.URL) string {
	if u.Host != "github.com" {
		return ""
	}
	seg := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(seg) < 2 {
		return ""
	}
	owner, repo := seg[0], seg[1]
	base := "https://api.github.com/repos/" + owner + "/" + repo
	switch {
	case len(seg) == 2:
		return base
	case len(seg) == 4 && seg[2] == "pull":
		return base + "/pulls/" + seg[3]
	case len(seg) == 4 && seg[2] == "issues":
		return base + "/issues/" + seg[3]
	case len(seg) == 4 && seg[2] == "commit":
		return base + "/commits/" + seg[3]
	case len(seg) >= 4 && seg[2] == "releases" && seg[3] == "tag":
		return base + "/releases/tags/" + strings.Join(seg[4:], "/")
	default:
		return ""
	}
}

var (
	typeRe   = regexp.MustCompile(`(?m)^type:\s*(\S+)\s*$`)
	statusRe = regexp.MustCompile(`(?m)^status:\s*(\S+)\s*$`)
)

// isFrozen reports whether an artifact is pinned history. A completed plan is
// the corpus's own frozen record — never deleted, never restated — so an
// address inside it is part of what was observed, not a pointer to be kept
// current.
func isFrozen(text string) bool {
	t := typeRe.FindStringSubmatch(text)
	s := statusRe.FindStringSubmatch(text)
	return t != nil && s != nil && t[1] == "plan" && s[1] == "completed"
}
