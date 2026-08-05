package corpus

import (
	"strings"
	"testing"
)

// icFiles extends capFiles with an imported-change record whose proof-links
// table names criterion, and a supporting Go test so AC-101 itself is
// provable (checkACTests satisfied) independent of what checkImportedChanges
// asserts about it.
func icFiles(criteriaStatus, icStatus, criterion string) map[string]string {
	files := with(capFiles(criteriaStatus), map[string]string{
		"docs/README.md":                    "# Corpus\n\n<!-- clue:index:start -->\n- [goals/](goals/README.md)\n- [plans/](plans/README.md)\n- [capabilities/](capabilities/README.md)\n- [imported-changes/](imported-changes/README.md)\n<!-- clue:index:end -->\n",
		"docs/imported-changes/README.md":   "# Imported changes\n\n<!-- clue:index:start -->\n- [IC-101](IC-101-x.md)\n<!-- clue:index:end -->\n",
		"docs/imported-changes/IC-101-x.md": "---\nid: IC-101\ntype: imported-change\nstatus: " + icStatus + "\nlinks: []\ntitle: X import\nsource-revision: abc123\nsource-location: example/repo\n---\n\n## Proof links\n\n| Task | Criterion |\n|---|---|\n| do the thing | " + criterion + " |\n",
	})
	files["pkg/x_test.go"] = "package x\n\nfunc TestAC101_Works(t *testing.T) {}\n"
	return files
}

// AC-115: an imported-change record's origin, proof links, and dependency
// stay inspectable and clue validate accepts a well-formed record.
func TestAC115_UnitPositive_WellFormedImportedChangeRecordPasses(t *testing.T) {
	files := icFiles("active", "complete", "AC-101")
	issues := run(t, files, false)
	for _, is := range issues {
		if strings.Contains(is.Path, "IC-101") {
			t.Fatalf("expected the imported-change record to pass, got: %v", issues)
		}
	}
}

// AC-115 negative: a record with no "## Proof links" table declares no
// inspectable proof and clue validate rejects it.
func TestAC115_UnitNegative_MissingProofLinksTableFails(t *testing.T) {
	files := with(capFiles("active"), map[string]string{
		"docs/README.md":                    "# Corpus\n\n<!-- clue:index:start -->\n- [goals/](goals/README.md)\n- [plans/](plans/README.md)\n- [capabilities/](capabilities/README.md)\n- [imported-changes/](imported-changes/README.md)\n<!-- clue:index:end -->\n",
		"docs/imported-changes/README.md":   "# Imported changes\n\n<!-- clue:index:start -->\n- [IC-101](IC-101-x.md)\n<!-- clue:index:end -->\n",
		"docs/imported-changes/IC-101-x.md": "---\nid: IC-101\ntype: imported-change\nstatus: in-progress\nlinks: []\ntitle: X import\nsource-revision: abc123\nsource-location: example/repo\n---\n\nNo proof-links section at all.\n",
	})
	assertIssue(t, run(t, files, false), "imported-change has no proof-links table")
}

// AC-116: a complete record naming an unproven criterion (here, one no
// criteria.md declares at all) is rejected by clue validate.
func TestAC116_UnitNegative_CompleteRecordWithUnprovenCriterionFails(t *testing.T) {
	files := icFiles("active", "complete", "AC-999")
	assertIssue(t, run(t, files, false), "imported-change is complete but proof-linked criterion AC-999 does not exist, is @draft, or is retired (ADR-050)")
}

// AC-116 positive: a complete record whose proof-linked criterion exists,
// is undrafted, and is not retired passes.
func TestAC116_UnitPositive_CompleteRecordWithProvenCriterionPasses(t *testing.T) {
	files := icFiles("active", "complete", "AC-101")
	issues := run(t, files, false)
	for _, is := range issues {
		if strings.Contains(is.Msg, "does not exist, is @draft, or is retired") {
			t.Fatalf("expected no unproven-completion issue, got: %v", issues)
		}
	}
}

// AC-117: an in-progress record naming a criterion no criteria.md declares
// is not rejected for that reason — clue validate exempts unproven links
// while status is in-progress.
func TestAC117_UnitPositive_InProgressRecordWithUnprovenCriterionPasses(t *testing.T) {
	files := icFiles("active", "in-progress", "AC-999")
	issues := run(t, files, false)
	for _, is := range issues {
		if strings.Contains(is.Msg, "does not exist, is @draft, or is retired") {
			t.Fatalf("in-progress record should not be held to the complete standard, got: %v", issues)
		}
	}
}

// AC-117 negative: the identical unproven link fails once the same record's
// status claims complete instead — proving the exemption is specific to
// in-progress, not a blanket pass on this rule.
func TestAC117_UnitNegative_TheSameRecordFailsOnceItClaimsComplete(t *testing.T) {
	files := icFiles("active", "complete", "AC-999")
	assertIssue(t, run(t, files, false), "imported-change is complete but proof-linked criterion AC-999 does not exist, is @draft, or is retired (ADR-050)")
}
