package migrate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cliewen/cliewen/internal/role"
)

func noticesFor(plan MigrationPlan, migration string) []Notice {
	var out []Notice
	for _, n := range plan.Notices {
		if n.Migration == migration {
			out = append(out, n)
		}
	}
	return out
}

func findingsFor(plan MigrationPlan, migration string) []Finding {
	var out []Finding
	for _, f := range plan.Findings {
		if f.Migration == migration {
			out = append(out, f)
		}
	}
	return out
}

func changesFor(plan MigrationPlan, migration string) []Change {
	var out []Change
	for _, c := range plan.Changes {
		if c.Migration == migration {
			out = append(out, c)
		}
	}
	return out
}

// An undeclared repository is told, never blocked. Every repository
// onboarded before the marker existed lacks one.
func TestAC153_UnitPositive_AnUndeclaredRoleIsANoticeAndNeverAFinding(t *testing.T) {
	root := migrationFixture(t, "")
	plan, err := Plan(root, Options{ReversalCost: "low"})
	if err != nil {
		t.Fatal(err)
	}
	notices := noticesFor(plan, MigrationRoleMarker)
	if len(notices) != 1 {
		t.Fatalf("expected one role notice, got %+v", notices)
	}
	if !strings.Contains(notices[0].Message, "has not declared its Cliewen role") {
		t.Fatalf("role notice does not say what is missing: %q", notices[0].Message)
	}
	if f := findingsFor(plan, MigrationRoleMarker); len(f) != 0 {
		t.Fatalf("role marker produced blocking findings: %+v", f)
	}
	if c := changesFor(plan, MigrationRoleMarker); len(c) != 0 {
		t.Fatalf("role marker planned a write: %+v", c)
	}
}

// Migration cannot tell the source repository from an adopter, which is the
// ambiguity the marker exists to end, so it never writes one.
func TestAC153_UnitPositive_ADeclaredRoleIsNotReported(t *testing.T) {
	root := migrationFixture(t, "")
	if err := role.Write(root, role.Adopter); err != nil {
		t.Fatal(err)
	}
	plan, err := Plan(root, Options{ReversalCost: "low"})
	if err != nil {
		t.Fatal(err)
	}
	if n := noticesFor(plan, MigrationRoleMarker); len(n) != 0 {
		t.Fatalf("a declared repository was still reported: %+v", n)
	}
}

// spentCorpus writes a corpus holding one spent analysis and one that is
// not spent, on top of the migration fixture.
func spentCorpus(t *testing.T) string {
	t.Helper()
	root := migrationFixture(t, "")
	files := map[string]string{
		"docs/README.md":                      "# Corpus\n\n<!-- clue:index:start -->\n- [plans/](plans/README.md)\n- [analysis/](analysis/README.md)\n- [architecture/](architecture/README.md)\n- [design/](design/README.md)\n<!-- clue:index:end -->\n",
		"docs/plans/README.md":                "# Plans\n\n<!-- clue:index:start -->\n- [P-001](P-001-done.md)\n- [P-002](P-002-open.md)\n<!-- clue:index:end -->\n",
		"docs/plans/P-001-done.md":            "---\nid: P-001\ntype: plan\nstatus: completed\nlinks: []\ntitle: Done\n---\n\n| M-001 | done | done |\n",
		"docs/plans/P-002-open.md":            "---\nid: P-002\ntype: plan\nstatus: active\nlinks: []\ntitle: Open\n---\n\n| M-002 | open | todo |\n",
		"docs/architecture/README.md":         "# Architecture\n\nStructure.\n\n<!-- clue:index:start -->\n- [ARCH-001](ARCH-001-shape.md)\n<!-- clue:index:end -->\n",
		"docs/architecture/ARCH-001-shape.md": "---\nid: ARCH-001\ntype: architecture\nstatus: active\nlinks: []\ntitle: Shape\n---\n\n# ARCH-001\n",
		"docs/design/README.md":               "# Design\n\nCross-cutting behaviour.\n\n<!-- clue:index:start -->\n<!-- clue:index:end -->\n",
		"docs/analysis/README.md":             "# Analysis\n\n<!-- clue:index:start -->\n- [AN-001](AN-001-spent.md)\n- [AN-002](AN-002-live.md)\n<!-- clue:index:end -->\n",
		"docs/analysis/AN-001-spent.md":       "---\nid: AN-001\ntype: analysis\nstatus: active\nlinks: [P-001]\ncarried-by: [ARCH-001]\ntitle: Spent\n---\n\n# AN-001\n",
		"docs/analysis/AN-002-live.md":        "---\nid: AN-002\ntype: analysis\nstatus: active\nlinks: [P-002]\ncarried-by: [ARCH-001]\ntitle: Live\n---\n\n# AN-002\n",
	}
	for rel, content := range files {
		full := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func TestAC156_UnitPositive_ASpentAnalysisIsReportedWithItsEvidence(t *testing.T) {
	root := spentCorpus(t)
	plan, err := Plan(root, Options{ReversalCost: "low"})
	if err != nil {
		t.Fatal(err)
	}
	notices := noticesFor(plan, MigrationSpentAnalysis)
	if len(notices) != 1 {
		t.Fatalf("expected exactly the spent analysis to be reported, got %+v", notices)
	}
	if !strings.Contains(notices[0].Path, "AN-001") {
		t.Fatalf("the wrong analysis was reported: %q", notices[0].Path)
	}
	for _, want := range []string{"ARCH-001", "P-001", "reviewed change"} {
		if !strings.Contains(notices[0].Message, want) {
			t.Fatalf("notice does not name %q: %q", want, notices[0].Message)
		}
	}
}

// The report is inert: retirement is a deletion, and ADR-034 keeps a
// deletion inside a reviewed change.
func TestAC156_UnitNegative_ReportingASpentAnalysisNeverWritesOrBlocks(t *testing.T) {
	root := spentCorpus(t)
	plan, err := Plan(root, Options{ReversalCost: "low"})
	if err != nil {
		t.Fatal(err)
	}
	if f := findingsFor(plan, MigrationSpentAnalysis); len(f) != 0 {
		t.Fatalf("a spent analysis blocked the migration: %+v", f)
	}
	if c := changesFor(plan, MigrationSpentAnalysis); len(c) != 0 {
		t.Fatalf("a spent analysis planned a write: %+v", c)
	}
	spent := filepath.Join(root, "docs", "analysis", "AN-001-spent.md")
	before, err := os.ReadFile(spent)
	if err != nil {
		t.Fatal(err)
	}
	if err := Apply(root, plan); err != nil {
		t.Fatalf("apply failed: %v", err)
	}
	after, err := os.ReadFile(spent)
	if err != nil {
		t.Fatalf("applying the migration removed the reported analysis: %v", err)
	}
	if string(before) != string(after) {
		t.Fatalf("applying the migration rewrote the reported analysis:\n%s", after)
	}
}
