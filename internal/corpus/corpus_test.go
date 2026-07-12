package corpus

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// validFiles is a minimal corpus that must pass every rule.
var validFiles = map[string]string{
	"docs/README.md":               "# Corpus\n\n<!-- clue:index:start -->\n- [goals/](goals/README.md)\n- [plans/](plans/README.md)\n<!-- clue:index:end -->\n",
	"docs/goals/README.md":         "# Goals\n\n<!-- clue:index:start -->\n- [G-001](G-001-first.md)\n<!-- clue:index:end -->\n",
	"docs/goals/G-001-first.md":    "---\nid: G-001\ntype: goal\nstatus: accepted\nlinks: []\ntitle: First goal\n---\n\n# G-001\n",
	"docs/plans/README.md":         "# Plans\n\n<!-- clue:index:start -->\n- [P-001](P-001-baseline.md)\n<!-- clue:index:end -->\n",
	"docs/plans/P-001-baseline.md": "---\nid: P-001\ntype: plan\nstatus: active\nlinks: [G-001]\ntitle: Baseline\n---\n\n| M-001 | do it | todo |\n",
}

func writeCorpus(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for p, content := range files {
		full := filepath.Join(root, filepath.FromSlash(p))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func with(base map[string]string, extra map[string]string) map[string]string {
	out := map[string]string{}
	for k, v := range base {
		out[k] = v
	}
	for k, v := range extra {
		out[k] = v
	}
	return out
}

func run(t *testing.T, files map[string]string, forbid bool) []Issue {
	t.Helper()
	c, issues := Scan(writeCorpus(t, files))
	return append(issues, Validate(c, Options{ForbidChanges: forbid})...)
}

func assertIssue(t *testing.T, issues []Issue, substr string) {
	t.Helper()
	for _, i := range issues {
		if strings.Contains(i.String(), substr) {
			return
		}
	}
	t.Fatalf("expected an issue containing %q, got %v", substr, issues)
}

// AC-004: a valid corpus exits clean.
func TestAC004_ValidCorpusHasNoIssues(t *testing.T) {
	if issues := run(t, validFiles, true); len(issues) != 0 {
		t.Fatalf("expected no issues, got %v", issues)
	}
}

// AC-005: a missing core field is reported with file and field.
func TestAC005_MissingCoreFieldReported(t *testing.T) {
	files := with(validFiles, map[string]string{
		"docs/goals/G-001-first.md": "---\nid: G-001\ntype: goal\nlinks: []\ntitle: First goal\n---\n",
	})
	issues := run(t, files, false)
	assertIssue(t, issues, "G-001-first.md")
	assertIssue(t, issues, "status")
}

func TestAC005_MissingFrontmatterReported(t *testing.T) {
	files := with(validFiles, map[string]string{
		"docs/goals/G-002-bare.md": "# No frontmatter here\n",
		"docs/goals/README.md":     "# Goals\n\n<!-- clue:index:start -->\n- [G-001](G-001-first.md)\n- [G-002](G-002-bare.md)\n<!-- clue:index:end -->\n",
	})
	assertIssue(t, run(t, files, false), "missing frontmatter")
}

// AC-006: a link to an unknown ID is reported with file and ID.
func TestAC006_UnresolvedLinkReported(t *testing.T) {
	files := with(validFiles, map[string]string{
		"docs/goals/G-001-first.md": "---\nid: G-001\ntype: goal\nstatus: accepted\nlinks: [G-999]\ntitle: First goal\n---\n",
	})
	issues := run(t, files, false)
	assertIssue(t, issues, "G-001-first.md")
	assertIssue(t, issues, "G-999")
}

func TestAC006_MilestoneLinksResolveViaPlanBody(t *testing.T) {
	files := with(validFiles, map[string]string{
		"docs/goals/G-001-first.md": "---\nid: G-001\ntype: goal\nstatus: accepted\nlinks: [M-001]\ntitle: First goal\n---\n",
	})
	if issues := run(t, files, false); len(issues) != 0 {
		t.Fatalf("M-001 exists in P-001 body; expected no issues, got %v", issues)
	}
	files["docs/goals/G-001-first.md"] = "---\nid: G-001\ntype: goal\nstatus: accepted\nlinks: [M-999]\ntitle: First goal\n---\n"
	assertIssue(t, run(t, files, false), "M-999")
}

func TestDuplicateIDReported(t *testing.T) {
	files := with(validFiles, map[string]string{
		"docs/goals/G-001-copy.md": validFiles["docs/goals/G-001-first.md"],
		"docs/goals/README.md":     "# Goals\n\n<!-- clue:index:start -->\n- [G-001](G-001-first.md)\n- [copy](G-001-copy.md)\n<!-- clue:index:end -->\n",
	})
	assertIssue(t, run(t, files, false), "duplicate id G-001")
}

func TestStatusVocabEnforced(t *testing.T) {
	files := with(validFiles, map[string]string{
		"docs/goals/G-001-first.md": "---\nid: G-001\ntype: goal\nstatus: wip\nlinks: []\ntitle: First goal\n---\n",
	})
	assertIssue(t, run(t, files, false), "status wip not allowed for type goal")

	files["docs/goals/G-001-first.md"] = "---\nid: G-001\ntype: wizard\nstatus: accepted\nlinks: []\ntitle: First goal\n---\n"
	assertIssue(t, run(t, files, false), "unknown type wizard")
}

func TestFolderWithoutReadmeReported(t *testing.T) {
	files := with(validFiles, map[string]string{
		"docs/quality/QS-001-fast.md": "---\nid: QS-001\ntype: quality\nstatus: active\nlinks: []\ntitle: Fast\n---\n",
		"docs/README.md":              "# Corpus\n\n<!-- clue:index:start -->\n- [goals/](goals/README.md)\n- [plans/](plans/README.md)\n- [quality/](quality/QS-001-fast.md)\n<!-- clue:index:end -->\n",
	})
	assertIssue(t, run(t, files, false), "docs/quality: folder has no README.md")
}

// AC-007: index drift — a block link to a missing file, or a sibling
// artifact the block does not reference.
func TestAC007_IndexDriftReported(t *testing.T) {
	files := with(validFiles, map[string]string{
		"docs/goals/README.md": "# Goals\n\n<!-- clue:index:start -->\n- [G-001](G-001-first.md)\n- [gone](G-777-gone.md)\n<!-- clue:index:end -->\n",
	})
	assertIssue(t, run(t, files, false), "index references missing file G-777-gone.md")

	files = with(validFiles, map[string]string{
		"docs/goals/README.md": "# Goals\n\n<!-- clue:index:start -->\n<!-- clue:index:end -->\n",
	})
	assertIssue(t, run(t, files, false), "index does not reference sibling G-001-first.md")
}

func TestAC007_MissingIndexMarkersReported(t *testing.T) {
	files := with(validFiles, map[string]string{
		"docs/goals/README.md": "# Goals — hand-written, no markers\n",
	})
	assertIssue(t, run(t, files, false), "index markers missing")
}

// AC-008: the digest-before-merge gate.
func TestAC008_ForbidChangesGate(t *testing.T) {
	files := with(validFiles, map[string]string{
		"changes/CH-009-x/proposal.md": "---\nid: CH-009\ntype: change\nstatus: open\nlinks: [P-001]\ntitle: X\n---\n",
	})
	if issues := run(t, files, false); len(issues) != 0 {
		t.Fatalf("changes/ allowed without the gate; got %v", issues)
	}
	assertIssue(t, run(t, files, true), "digest before merge")
}

func TestCRLFFrontmatterParses(t *testing.T) {
	files := with(validFiles, map[string]string{
		"docs/goals/G-001-first.md": strings.ReplaceAll(validFiles["docs/goals/G-001-first.md"], "\n", "\r\n"),
	})
	if issues := run(t, files, false); len(issues) != 0 {
		t.Fatalf("CRLF corpus should validate; got %v", issues)
	}
}

// The dogfood test: this repository's own corpus must always be valid.
func TestRepoCorpusIsValid(t *testing.T) {
	root, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(root, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(root)
		if parent == root {
			t.Fatal("repo root with go.mod not found")
		}
		root = parent
	}
	c, issues := Scan(root)
	issues = append(issues, Validate(c, Options{ForbidChanges: false})...)
	if len(issues) != 0 {
		t.Fatalf("the repo's own corpus has issues:\n%v", issues)
	}
	if len(c.Artifacts) == 0 {
		t.Fatal("expected artifacts in the repo corpus")
	}
}
