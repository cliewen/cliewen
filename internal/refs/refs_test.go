package refs

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// corpusWith writes a docs tree containing one artifact per entry.
func corpusWith(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for name, body := range files {
		full := filepath.Join(root, "docs", filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

// forge answers each path with the status the test asked for. It stands in for
// every host a corpus points at, so the classification is exercised without a
// network.
func forge(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/owner/live/ok", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("/owner/live/forbidden", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	})
	mux.HandleFunc("/owner/live/moved", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", "/owner/live/ok")
		w.WriteHeader(http.StatusMovedPermanently)
	})
	// A private target and its live owner: the shape GitHub produces, where a
	// 404 is a refusal that declines to confirm the target exists.
	mux.HandleFunc("/owner/live/private", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	mux.HandleFunc("/owner", func(w http.ResponseWriter, r *http.Request) {})
	// Siblings sharing the redirected address's prefix. Each answers 200, so
	// each is classified reachable and must survive a neighbouring rewrite.
	mux.HandleFunc("/owner/live/moved.git", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("/owner/live/movedx", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("/owner/live/moved/deeper/page", func(w http.ResponseWriter, r *http.Request) {})
	// A deleted target whose owner is gone too.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

func resultFor(t *testing.T, report Report, contains string) Reference {
	t.Helper()
	for _, r := range report.References {
		if strings.Contains(r.URL, contains) {
			return r
		}
	}
	t.Fatalf("no reference containing %q in %v", contains, report.References)
	return Reference{}
}

func TestAC068_UnitPositive_EveryOutcomeIsClassified(t *testing.T) {
	srv := forge(t)
	body := "" +
		"reachable " + srv.URL + "/owner/live/ok\n" +
		"restricted " + srv.URL + "/owner/live/forbidden\n" +
		"redirected " + srv.URL + "/owner/live/moved\n" +
		"private " + srv.URL + "/owner/live/private\n" +
		"deleted " + srv.URL + "/nobody/vanished\n" +
		"identity clue:robocode-dev/tank-royale@384d27d5/BR-001\n"
	root := corpusWith(t, map[string]string{"analysis/AN-001.md": "---\ntype: analysis\nstatus: active\n---\n\n" + body})

	report, err := Resolve(root, Options{Client: srv.Client()})
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		url  string
		want Outcome
	}{
		{"/owner/live/ok", Reachable},
		{"/owner/live/forbidden", Restricted},
		{"/owner/live/moved", Redirected},
		{"/owner/live/private", Restricted}, // 404 whose owner answers
		{"/nobody/vanished", Gone},
	} {
		if got := resultFor(t, report, c.url).Result; got != c.want {
			t.Errorf("%s: got %s, want %s", c.url, got, c.want)
		}
	}
	if got := resultFor(t, report, "/owner/live/moved").Target; !strings.HasSuffix(got, "/owner/live/ok") {
		t.Errorf("expected the new location, got %q", got)
	}
	// The clue: identity names a corpus artifact elsewhere; following it would
	// need to know where that repository currently lives.
	for _, r := range report.References {
		if strings.HasPrefix(r.URL, "clue:") {
			t.Errorf("a clue: identity was resolved: %v", r)
		}
	}
}

func TestAC068_UnitNegative_OnlyGoneIsAnError(t *testing.T) {
	srv := forge(t)
	base := "---\ntype: analysis\nstatus: active\n---\n\n"
	root := corpusWith(t, map[string]string{
		"analysis/AN-001.md": base +
			"restricted " + srv.URL + "/owner/live/forbidden\n" +
			"private " + srv.URL + "/owner/live/private\n",
	})
	report, err := Resolve(root, Options{Client: srv.Client()})
	if err != nil {
		t.Fatal(err)
	}
	if report.HasErrors() {
		t.Fatal("restricted must never fail the run — it says nothing about whether the corpus is wrong")
	}

	gone := corpusWith(t, map[string]string{"analysis/AN-002.md": base + "deleted " + srv.URL + "/nobody/vanished\n"})
	report, err = Resolve(gone, Options{Client: srv.Client()})
	if err != nil {
		t.Fatal(err)
	}
	if !report.HasErrors() {
		t.Fatal("a gone reference is the one outcome that is an error")
	}
}

func TestAC069_UnitPositive_ApplyRewritesOnlyOrdinaryArtifacts(t *testing.T) {
	srv := forge(t)
	ordinary := "---\ntype: analysis\nstatus: active\n---\n\nmoved " + srv.URL + "/owner/live/moved\n"
	root := corpusWith(t, map[string]string{"analysis/AN-001.md": ordinary})

	report, err := Resolve(root, Options{Apply: true, Client: srv.Client()})
	if err != nil {
		t.Fatal(err)
	}
	if !report.Applied {
		t.Fatal("expected the run to report that it wrote")
	}
	got, err := os.ReadFile(filepath.Join(root, "docs", "analysis", "AN-001.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(got), "/owner/live/ok") {
		t.Fatalf("expected the address rewritten to its new location, got %q", got)
	}
}

func TestAC069_UnitNegative_PinnedHistoryAndGoneAreNeverRewritten(t *testing.T) {
	srv := forge(t)
	// A completed plan records what was observed at the time: repointing its
	// address would falsify the observation rather than repair it.
	plan := "---\ntype: plan\nstatus: completed\n---\n\nthe guide at " + srv.URL + "/owner/live/moved returned HTTP 200\n"
	gone := "---\ntype: analysis\nstatus: active\n---\n\ndeleted " + srv.URL + "/nobody/vanished\n"
	root := corpusWith(t, map[string]string{
		"plans/P-003.md":     plan,
		"analysis/AN-001.md": gone,
	})

	report, err := Resolve(root, Options{Apply: true, Client: srv.Client()})
	if err != nil {
		t.Fatal(err)
	}
	if got := resultFor(t, report, "/owner/live/moved"); !got.Frozen {
		t.Fatal("a completed plan's address must be reported as pinned history")
	}
	after, err := os.ReadFile(filepath.Join(root, "docs", "plans", "P-003.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != plan {
		t.Fatalf("pinned history was rewritten:\n%q", after)
	}
	unchanged, err := os.ReadFile(filepath.Join(root, "docs", "analysis", "AN-001.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(unchanged) != gone {
		t.Fatal("a gone address has no correct replacement to offer and must be left alone")
	}
}

func TestAC069_UnitNegative_CredentialGoesOnlyToItsIssuer(t *testing.T) {
	// A forge token attached to an arbitrary wiki would hand that host a
	// credential it was never meant to see.
	t.Setenv("GITHUB_TOKEN", "secret")
	if got := hostToken("wiki.example.com"); got != "" {
		t.Fatalf("a third-party host must receive no credential, got %q", got)
	}
	if got := hostToken("github.com"); got != "secret" {
		t.Fatalf("the issuing host should receive the credential, got %q", got)
	}
}

func TestAC068_UnitNegative_ATimeoutIsUnknownAndDoesNotFail(t *testing.T) {
	// A host that never answers must be reported as unknown: the corpus did
	// not change, and someone else's outage cannot condemn it.
	slow := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(300 * time.Millisecond)
	}))
	t.Cleanup(slow.Close)
	root := corpusWith(t, map[string]string{
		"analysis/AN-001.md": "---\ntype: analysis\nstatus: active\n---\n\nslow " + slow.URL + "/owner/thing\n",
	})
	report, err := Resolve(root, Options{Client: slow.Client(), Timeout: 20 * time.Millisecond})
	if err != nil {
		t.Fatal(err)
	}
	got := resultFor(t, report, "/owner/thing")
	if got.Result != Unreachable {
		t.Fatalf("expected unreachable, got %s", got.Result)
	}
	if report.HasErrors() {
		t.Fatal("an unreachable address must not fail the run")
	}
}

func TestAC069_UnitNegative_PreviewWritesNothing(t *testing.T) {
	// The preview/write split is the safety property the whole --apply design
	// rests on, so it is asserted rather than assumed.
	srv := forge(t)
	before := "---\ntype: analysis\nstatus: active\n---\n\nmoved " + srv.URL + "/owner/live/moved\n"
	root := corpusWith(t, map[string]string{"analysis/AN-001.md": before})

	report, err := Resolve(root, Options{Client: srv.Client()})
	if err != nil {
		t.Fatal(err)
	}
	if report.Applied {
		t.Fatal("a preview must not report that it wrote")
	}
	after, err := os.ReadFile(filepath.Join(root, "docs", "analysis", "AN-001.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != before {
		t.Fatalf("a preview wrote to the corpus:\n%q", after)
	}
}

func TestAC069_UnitNegative_ARenamedPrefixDoesNotCorruptLongerAddresses(t *testing.T) {
	// A renamed repository's address is a prefix of every address below it.
	// Replacing by plain substring would rewrite a longer address the run
	// never classified.
	srv := forge(t)
	line := "see " + srv.URL + "/owner/live/moved and " + srv.URL + "/owner/live/moved/deeper/page\n"
	root := corpusWith(t, map[string]string{
		"analysis/AN-001.md": "---\ntype: analysis\nstatus: active\n---\n\n" + line,
	})
	if _, err := Resolve(root, Options{Apply: true, Client: srv.Client()}); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filepath.Join(root, "docs", "analysis", "AN-001.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(got), "/owner/live/moved/deeper/page") {
		t.Fatalf("the longer address was corrupted by a prefix rewrite:\n%q", got)
	}
}

func TestAC068_UnitNegative_AThrottleIsNotARefusal(t *testing.T) {
	// A rate-limited 403 classified as restricted would report every affected
	// address as probably-private, destroying the signal that separating the
	// two outcomes exists to preserve.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "60")
		w.WriteHeader(http.StatusForbidden)
	}))
	t.Cleanup(srv.Close)
	root := corpusWith(t, map[string]string{
		"analysis/AN-001.md": "---\ntype: analysis\nstatus: active\n---\n\nthrottled " + srv.URL + "/owner/thing\n",
	})
	report, err := Resolve(root, Options{Client: srv.Client()})
	if err != nil {
		t.Fatal(err)
	}
	if got := resultFor(t, report, "/owner/thing").Result; got != Unreachable {
		t.Fatalf("a throttled 403 must be unknown, not a refusal: got %s", got)
	}
	if report.HasErrors() {
		t.Fatal("a throttle must not fail the run")
	}
}

func TestAC069_UnitNegative_AnAddressEndingASentenceIsStillRewritten(t *testing.T) {
	// Harvesting strips trailing sentence punctuation, so a boundary set that
	// omitted it would skip the rewrite and still report the run as applied —
	// telling the user a repair happened that did not.
	srv := forge(t)
	for _, suffix := range []string{".", ",", ":", ")", "!", "|", "...", ""} {
		t.Run("ends with "+suffix, func(t *testing.T) {
			root := corpusWith(t, map[string]string{
				"analysis/AN-001.md": "---\ntype: analysis\nstatus: active\n---\n\nSee " + srv.URL + "/owner/live/moved" + suffix + "\n",
			})
			if _, err := Resolve(root, Options{Apply: true, Client: srv.Client()}); err != nil {
				t.Fatal(err)
			}
			got, err := os.ReadFile(filepath.Join(root, "docs", "analysis", "AN-001.md"))
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(string(got), "/owner/live/ok") {
				t.Fatalf("the rewrite was skipped:\n%q", got)
			}
		})
	}
}

func TestAC069_UnitNegative_ASiblingSharingAPrefixIsNeverCorrupted(t *testing.T) {
	// Every character admitted as a boundary is a way for a classified address
	// to be substituted into a longer one the run never looked at. These are
	// the shapes that hazard actually takes.
	srv := forge(t)
	for _, sibling := range []string{
		"/owner/live/moved.git",
		"/owner/live/movedx",
		"/owner/live/moved/deeper/page",
		// A fragment and a query string are deliberately absent: neither
		// reaches the server as a distinct resource, so an address carrying
		// one is the same target and is meant to move with it.
	} {
		t.Run(sibling, func(t *testing.T) {
			line := "see " + srv.URL + "/owner/live/moved and " + srv.URL + sibling + "\n"
			root := corpusWith(t, map[string]string{
				"analysis/AN-001.md": "---\ntype: analysis\nstatus: active\n---\n\n" + line,
			})
			if _, err := Resolve(root, Options{Apply: true, Client: srv.Client()}); err != nil {
				t.Fatal(err)
			}
			got, err := os.ReadFile(filepath.Join(root, "docs", "analysis", "AN-001.md"))
			if err != nil {
				t.Fatal(err)
			}
			// Both halves of the trade-off, or the test guards nothing: a
			// rewrite that never happens also leaves the sibling intact.
			if !strings.Contains(string(got), "/owner/live/ok") {
				t.Fatalf("the classified address was not rewritten:\n%q", got)
			}
			if !strings.Contains(string(got), srv.URL+sibling) {
				t.Fatalf("the sibling address was corrupted:\n%q", got)
			}
		})
	}
}

func TestAC069_UnitNegative_ACarriageReturnDoesNotSkipTheRewrite(t *testing.T) {
	// A checkout without LF normalisation is the one environment where the
	// skipped-but-reported-applied regression would recur silently.
	srv := forge(t)
	root := corpusWith(t, map[string]string{
		"analysis/AN-001.md": "---\r\ntype: analysis\r\nstatus: active\r\n---\r\n\r\nSee " + srv.URL + "/owner/live/moved\r\n",
	})
	if _, err := Resolve(root, Options{Apply: true, Client: srv.Client()}); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filepath.Join(root, "docs", "analysis", "AN-001.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(got), "/owner/live/ok") {
		t.Fatalf("a CRLF line skipped the rewrite:\n%q", got)
	}
}

func TestAC068_UnitNegative_AQueryStringIsPartOfTheAddress(t *testing.T) {
	// A question mark opens a query string, not a sentence. Treating it as
	// prose truncates a legitimate address into one no host has, which then
	// answers 404 and condemns a corpus that is right.
	mux := http.NewServeMux()
	mux.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("name") == "" {
			w.WriteHeader(http.StatusNotFound)
		}
	})
	mux.HandleFunc("/owner", func(w http.ResponseWriter, r *http.Request) {})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	full := srv.URL + "/items?name=remote-containers"
	root := corpusWith(t, map[string]string{
		"analysis/AN-001.md": "---\ntype: analysis\nstatus: active\n---\n\n" +
			"bare " + full + "\nlinked [label](" + full + ")\n",
	})
	report, err := Resolve(root, Options{Client: srv.Client()})
	if err != nil {
		t.Fatal(err)
	}
	if report.HasErrors() {
		t.Fatalf("a query-string address must not be condemned: %v", report.References)
	}
	for _, r := range report.References {
		if r.URL != full {
			t.Fatalf("expected the whole address, got %q", r.URL)
		}
	}
}

func TestAC069_UnitNegative_ACodeSpanIsContentNotACitation(t *testing.T) {
	// A literal command example is content. Classifying it would let --apply
	// rewrite the example itself when the host redirects.
	srv := forge(t)
	body := "Run `git clone " + srv.URL + "/owner/live/moved` to start.\n"
	root := corpusWith(t, map[string]string{
		"analysis/AN-001.md": "---\ntype: analysis\nstatus: active\n---\n\n" + body,
	})
	report, err := Resolve(root, Options{Apply: true, Client: srv.Client()})
	if err != nil {
		t.Fatal(err)
	}
	if len(report.References) != 0 {
		t.Fatalf("a code span must not be harvested, got %v", report.References)
	}
	got, err := os.ReadFile(filepath.Join(root, "docs", "analysis", "AN-001.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(got), "/owner/live/moved") {
		t.Fatalf("the example was rewritten:\n%q", got)
	}
}

func TestAC068_UnitPositive_AnEmphasisedOrQuotedAddressIsStillSeen(t *testing.T) {
	// A silent miss cannot fail a corpus, but a rotted address inside bold
	// text or a tight table row would be invisible, which is worse than noisy.
	srv := forge(t)
	for _, shape := range []string{
		"**" + srv.URL + "/nobody/vanished**",
		"\"" + srv.URL + "/nobody/vanished\"",
		"|" + srv.URL + "/nobody/vanished|",
	} {
		t.Run(shape, func(t *testing.T) {
			root := corpusWith(t, map[string]string{
				"analysis/AN-001.md": "---\ntype: analysis\nstatus: active\n---\n\nsee " + shape + "\n",
			})
			report, err := Resolve(root, Options{Client: srv.Client()})
			if err != nil {
				t.Fatal(err)
			}
			got := resultFor(t, report, "/nobody/vanished")
			if got.URL != srv.URL+"/nobody/vanished" {
				t.Fatalf("expected the address without its wrapper, got %q", got.URL)
			}
		})
	}
}

func TestAC069_UnitNegative_ARewriteNeverReachesIntoACodeSpan(t *testing.T) {
	// The same address bare and in an example on one line: the citation moves,
	// the literal does not.
	srv := forge(t)
	line := "See " + srv.URL + "/owner/live/moved — run `git clone " + srv.URL + "/owner/live/moved` to fetch it.\n"
	root := corpusWith(t, map[string]string{
		"analysis/AN-001.md": "---\ntype: analysis\nstatus: active\n---\n\n" + line,
	})
	if _, err := Resolve(root, Options{Apply: true, Client: srv.Client()}); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filepath.Join(root, "docs", "analysis", "AN-001.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(got), "See "+srv.URL+"/owner/live/ok") {
		t.Fatalf("the citation was not rewritten:\n%q", got)
	}
	if !strings.Contains(string(got), "git clone "+srv.URL+"/owner/live/moved`") {
		t.Fatalf("the example inside the code span was rewritten:\n%q", got)
	}
}
