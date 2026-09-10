package corpus

import (
	"strings"
	"testing"
)

func TestAC186_UnitPositive_UnfinishedMilestonesPrioritizeDoingAndDraftIsSeparate(t *testing.T) {
	c := &Corpus{Artifacts: []*Artifact{
		{Type: "plan", Status: "active", Path: "docs/plans/P-002.md", Body: "| ID | Milestone | Status | Evidence |\n|---|---|---|---|\n| M-002 | Later | todo | |\n"},
		{Type: "plan", Status: "active", Path: "docs/plans/P-001.md", Body: "| ID | Milestone | Exit criterion | Status | Evidence |\n|---|---|---|---|---|\n| M-001 | First | Do first | doing | |\n"},
		{Type: "plan", Status: "draft", Path: "docs/plans/P-003.md", Body: "| ID | Milestone | Status | Evidence |\n|---|---|---|---|\n| M-003 | Proposed | todo | |\n"},
	}}
	actionable := UnfinishedMilestones(c)
	if len(actionable) != 2 || actionable[0].ID != "M-001" || actionable[1].ID != "M-002" {
		t.Fatalf("actionable order = %#v", actionable)
	}
	if actionable[0].ExitCriterion != "Do first" {
		t.Fatalf("exit criterion = %q", actionable[0].ExitCriterion)
	}
	draft := DraftUnfinishedMilestones(c)
	if len(draft) != 1 || draft[0].ID != "M-003" {
		t.Fatalf("draft milestones = %#v", draft)
	}
}

func TestAC186_UnitNegative_PlanMilestonesIgnoreExamplesAndFinishedRows(t *testing.T) {
	plan := &Artifact{Type: "plan", Body: strings.TrimSpace(`
| ID | Milestone | Status |
|---|---|---|
| M-001 | Done | done |
| M-002 | Dropped | dropped |
| M-003 | Work | todo |

    | ID | Milestone | Status |
    |---|---|---|
    | M-004 | Example | todo |

<details>
| ID | Milestone | Status |
|---|---|---|
| M-005 | HTML example | todo |
</details>
`)}
	got := PlanMilestones(plan)
	if len(got) != 3 || got[0].ID != "M-001" || got[2].ID != "M-003" {
		t.Fatalf("parsed milestones = %#v", got)
	}
	if got[1].Status != "dropped" {
		t.Fatalf("dropped status = %q", got[1].Status)
	}
}

func TestAC187_UnitPositive_LedgerMilestoneIdentitiesSpanEveryPlanStatusInIDOrder(t *testing.T) {
	c := &Corpus{Artifacts: []*Artifact{
		{Type: "plan", Status: "active", Path: "docs/plans/P-002.md", Body: "| ID | Milestone | Status | Evidence |\n|---|---|---|---|\n| M-090 | Live | todo | |\n"},
		{Type: "plan", Status: "completed", Path: "docs/plans/P-001.md", Body: "| ID | Milestone | Status | Evidence |\n|---|---|---|---|\n| M-001 | Done | done | |\n| M-043 | Withdrawn | dropped | |\n"},
	}}
	got := LedgerMilestoneIdentities(c)
	if len(got) != 3 || got[0].ID != "M-001" || got[1].ID != "M-043" || got[2].ID != "M-090" {
		t.Fatalf("ledger milestone identities = %#v", got)
	}
}

func TestAC187_UnitNegative_LedgerMilestoneIdentitiesEmptyWithoutPlans(t *testing.T) {
	c := &Corpus{Artifacts: []*Artifact{{Type: "analysis", Path: "docs/analysis/AN-001.md"}}}
	if got := LedgerMilestoneIdentities(c); len(got) != 0 {
		t.Fatalf("ledger milestone identities = %#v, want none", got)
	}
}
