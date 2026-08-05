package migrate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// promotedNotice returns the MIG-007 notice in a plan, or nil.
func promotedNotice(plan MigrationPlan) *Notice {
	for i := range plan.Notices {
		if plan.Notices[i].Migration == MigrationPromotedConstraints {
			return &plan.Notices[i]
		}
	}
	return nil
}

const promotedFrontmatter = "---\nid: C-001\ntype: constraint\nstatus: active\nlinks: []\ntitle: Markdown prose is never hard-wrapped\nsource: %SOURCE%\nenforcement: %CLASS%\n---\n\n# C-001\n\nOne line per paragraph.\n"

func writeConstraint(t *testing.T, root, name, source, class string) {
	t.Helper()
	dir := filepath.Join(root, "docs", "constraints")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	body := strings.NewReplacer("%SOURCE%", source, "%CLASS%", class).Replace(promotedFrontmatter)
	if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

// TestAC064_UnitPositive_PromotedConstraintIsReported proves MIG-007 names the
// register entry whose promotion trigger this release satisfied.
func TestAC064_UnitPositive_PromotedConstraintIsReported(t *testing.T) {
	root := migrationFixture(t, "")
	writeConstraint(t, root, "C-001-no-hard-wrapped-markdown.md", scaffoldedConstraintSource, "agent")

	plan, err := Plan(root, Options{})
	if err != nil {
		t.Fatal(err)
	}
	notice := promotedNotice(plan)
	if notice == nil {
		t.Fatalf("no MIG-007 notice for a scaffolded constraint still awaiting a machine check: %+v", plan.Notices)
	}
	if notice.Path != "docs/constraints/C-001-no-hard-wrapped-markdown.md" {
		t.Fatalf("notice names %q rather than the register entry", notice.Path)
	}
	for _, want := range []string{"enforcement: machine", "Checked by", "does not edit"} {
		if !strings.Contains(notice.Message, want) {
			t.Fatalf("notice does not state %q: %s", want, notice.Message)
		}
	}
	// A notice never blocks: an adopter's register being one class stale is not
	// a reason to refuse to refresh their skills.
	for _, f := range plan.Findings {
		if f.Migration == MigrationPromotedConstraints {
			t.Fatalf("MIG-007 raised a blocking finding: %+v", f)
		}
	}
}

// TestAC064_UnitNegative_OwnAndPromotedConstraintsAreLeftAlone proves the
// report is narrow. A constraint the adopter wrote carries their own promotion
// trigger, which this project is in no position to call satisfied, and a
// scaffolded one already promoted is simply current.
func TestAC064_UnitNegative_OwnAndPromotedConstraintsAreLeftAlone(t *testing.T) {
	for _, tc := range []struct {
		name   string
		source string
		class  string
	}{
		{name: "the adopter's own agent-enforced rule", source: "our engineering handbook", class: "agent"},
		{name: "an already promoted scaffolded rule", source: scaffoldedConstraintSource, class: "machine"},
		{name: "a scaffolded rule declared human", source: scaffoldedConstraintSource, class: "human"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := migrationFixture(t, "")
			writeConstraint(t, root, "C-001-a-rule.md", tc.source, tc.class)

			plan, err := Plan(root, Options{})
			if err != nil {
				t.Fatal(err)
			}
			if notice := promotedNotice(plan); notice != nil {
				t.Fatalf("MIG-007 reported %s: %+v", tc.name, notice)
			}
		})
	}
}
