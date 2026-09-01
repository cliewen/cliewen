package corpus

import "testing"

// retiredTargetFiles holds a corpus where AN-001 has been retired: the file
// is gone and ADR-001 carries the pointer forward in supersedes.
var retiredTargetFiles = map[string]string{
	"docs/README.md":                      "# Corpus\n\n<!-- clue:index:start -->\n- [goals/](goals/README.md)\n- [plans/](plans/README.md)\n- [decisions/](decisions/README.md)\n- [architecture/](architecture/README.md)\n- [design/](design/README.md)\n<!-- clue:index:end -->\n",
	"docs/architecture/README.md":         "# Architecture\n\nStructure.\n\n<!-- clue:index:start -->\n<!-- clue:index:end -->\n",
	"docs/design/README.md":               "# Design\n\nCross-cutting behaviour.\n\n<!-- clue:index:start -->\n<!-- clue:index:end -->\n",
	"docs/decisions/README.md":            "# Decisions\n\n<!-- clue:index:start -->\n- [ADR-001](ADR-001-successor.md)\n<!-- clue:index:end -->\n",
	"docs/decisions/ADR-001-successor.md": "---\nid: ADR-001\ntype: decision\nstatus: inferred\nlinks: []\nsupersedes: [AN-001]\ntitle: Successor\nauthor: agent\naccepted-by: []\n---\n\n# ADR-001\n",
	"docs/plans/README.md":                "# Plans\n\n<!-- clue:index:start -->\n- [P-001](P-001-baseline.md)\n<!-- clue:index:end -->\n",
}

func planLinking(status, target string) string {
	return "---\nid: P-001\ntype: plan\nstatus: " + status + "\nlinks: [G-001, " + target + "]\ntitle: Baseline\n---\n\n| M-001 | do it | done |\n"
}

// A finished campaign records what it referenced while it ran. C-008 freezes
// the file, so demanding a repoint would ask for an edit the guard forbids.
func TestAC157_UnitPositive_ACompletedPlanKeepsALinkToARetiredArtifact(t *testing.T) {
	files := with(validFiles, with(retiredTargetFiles, map[string]string{
		"docs/plans/P-001-baseline.md": planLinking("completed", "AN-001"),
	}))
	if issues := run(t, files, false); len(issues) != 0 {
		t.Fatalf("a completed plan's link to a retired artifact should pass, got %v", issues)
	}
}

// The allowance is for frozen history only. An active plan can still be
// edited, so it repoints like anything else.
func TestAC157_UnitNegative_AnActivePlanStillRepoints(t *testing.T) {
	files := with(validFiles, with(retiredTargetFiles, map[string]string{
		"docs/plans/P-001-baseline.md": planLinking("active", "AN-001"),
	}))
	assertIssue(t, run(t, files, false), "link AN-001 was retired — repoint to its successor ADR-001")
}

// A link to an ID nothing ever declared is a typo, not history, and the
// allowance must not swallow it.
func TestAC157_UnitNegative_ACompletedPlanStillFailsOnANeverDeclaredID(t *testing.T) {
	files := with(validFiles, with(retiredTargetFiles, map[string]string{
		"docs/plans/P-001-baseline.md": planLinking("completed", "AN-404"),
	}))
	assertIssue(t, run(t, files, false), "link AN-404 resolves to no artifact")
}
