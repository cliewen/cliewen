package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestAC067_UnitPositive_CoverageListsForeignPointersApart reads what the
// command prints. The criterion describes user-visible output, so asserting a
// string the test builds itself would prove nothing — deleting the print would
// leave such a test green.
func TestAC067_UnitPositive_CoverageListsForeignPointersApart(t *testing.T) {
	root := validCorpus(t)
	writeFile(t, root, "docs/README.md", "# Corpus\n\n<!-- clue:index:start -->\n- [goals/](goals/README.md)\n- [analysis/](analysis/README.md)\n<!-- clue:index:end -->\n")
	writeFile(t, root, "docs/analysis/README.md", "# Analysis\n\n<!-- clue:index:start -->\n- [AN-101](AN-101-x.md)\n<!-- clue:index:end -->\n")
	writeFile(t, root, "docs/analysis/AN-101-x.md", "---\nid: AN-101\ntype: analysis\nstatus: active\nlinks: []\ntitle: X\nprovenance: inferred\nreversal-cost: low\n---\n\nProven by clue:robocode-dev/tank-royale@384d27d5/BR-001 upstream.\n")

	code, out := runValidateCapturingStdout(t, []string{"--coverage", root})
	if code != 0 {
		t.Fatalf("expected a valid corpus, got exit %d: %s", code, out)
	}
	want := "clue:robocode-dev/tank-royale@384d27d5/BR-001: named but locally unproven"
	if !strings.Contains(out, want) {
		t.Fatalf("expected the pointer listed, got:\n%s", out)
	}
	// It must sit apart from the capability states, never among them: a reader
	// scanning for covered/partial/gap must not meet a pointer in that column.
	lines := strings.Split(strings.TrimSpace(out), "\n")
	lastState, pointerAt := -1, -1
	for i, ln := range lines {
		switch {
		case strings.HasSuffix(ln, ": covered"), strings.HasSuffix(ln, ": partial"), strings.HasSuffix(ln, ": gap"):
			lastState = i
			if strings.HasPrefix(ln, "clue:") {
				t.Fatalf("a pointer was printed as a coverage state: %q", ln)
			}
		case ln == want:
			pointerAt = i
		}
	}
	if pointerAt < 0 || (lastState >= 0 && pointerAt < lastState) {
		t.Fatalf("expected the pointer after the capability states, got:\n%s", out)
	}
}

// TestAC068_UnitNegative_RefsExitsNonZeroOnlyForGone drives the command itself,
// so the exit code the criterion names is proven rather than inferred from the
// report type.
func TestAC068_UnitNegative_RefsExitsNonZeroOnlyForGone(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/owner/live/forbidden", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	})
	mux.HandleFunc("/owner", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	write := func(t *testing.T, body string) string {
		t.Helper()
		root := t.TempDir()
		dir := filepath.Join(root, "docs", "analysis")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		text := "---\nid: AN-001\ntype: analysis\nstatus: active\nlinks: []\ntitle: t\n---\n\n" + body + "\n"
		if err := os.WriteFile(filepath.Join(dir, "AN-001-x.md"), []byte(text), 0o644); err != nil {
			t.Fatal(err)
		}
		return root
	}

	var out, errOut bytes.Buffer
	restricted := write(t, "restricted "+srv.URL+"/owner/live/forbidden")
	if code := runRefs([]string{restricted}, &out, &errOut); code != 0 {
		t.Fatalf("a restricted address must not fail the command, got exit %d: %s", code, out.String())
	}

	out.Reset()
	gone := write(t, "deleted "+srv.URL+"/nobody/vanished")
	if code := runRefs([]string{gone}, &out, &errOut); code != 1 {
		t.Fatalf("a gone address must exit 1, got %d: %s", code, out.String())
	}

	out.Reset()
	if code := runRefs([]string{"a", "b"}, &out, &errOut); code != 2 {
		t.Fatalf("two paths is a usage error, got exit %d", code)
	}
}

// TestAC069_UnitPositive_PinnedHistoryIsMarkedInTheReport reads the printed
// note. The criterion describes user-visible output, and deleting the marker
// left every other test green.
func TestAC069_UnitPositive_PinnedHistoryIsMarkedInTheReport(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/owner/live/ok", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("/owner/live/moved", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", "/owner/live/ok")
		w.WriteHeader(http.StatusMovedPermanently)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	root := t.TempDir()
	for rel, body := range map[string]string{
		"docs/analysis/AN-001-x.md": "---\nid: AN-001\ntype: analysis\nstatus: active\nlinks: []\ntitle: t\n---\n\nordinary " + srv.URL + "/owner/live/moved\n",
		"docs/plans/P-003-x.md":     "---\nid: P-003\ntype: plan\nstatus: completed\nlinks: []\ntitle: t\n---\n\nthe guide at " + srv.URL + "/owner/live/moved returned HTTP 200\n",
	} {
		full := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	var out, errOut bytes.Buffer
	if code := runRefs([]string{root}, &out, &errOut); code != 0 {
		t.Fatalf("a redirect is not an error, got exit %d: %s", code, out.String())
	}
	got := out.String()
	var ordinary, frozen string
	for _, ln := range strings.Split(got, "\n") {
		switch {
		case strings.Contains(ln, "AN-001"):
			ordinary = ln
		case strings.Contains(ln, "P-003"):
			frozen = ln
		}
	}
	if ordinary == "" || frozen == "" {
		t.Fatalf("expected both redirects named, got:\n%s", got)
	}
	if !strings.Contains(frozen, "pinned history") {
		t.Fatalf("the completed plan must be marked as pinned history: %q", frozen)
	}
	if strings.Contains(ordinary, "pinned history") {
		t.Fatalf("an ordinary artifact must not be marked pinned: %q", ordinary)
	}
}
