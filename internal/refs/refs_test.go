package refs

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
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
