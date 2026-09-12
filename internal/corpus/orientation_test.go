package corpus

import "testing"

func TestAC190_UnitPositive_OrientationFindsOpenChangesAndProposedGoalsInStableOrder(t *testing.T) {
	c := &Corpus{Artifacts: []*Artifact{
		{ID: "G-002", Type: "goal", Status: "proposed", Path: "docs/goals/G-002.md"},
		{ID: "CH-002", Type: "change", Status: "open", Path: "changes/CH-002/proposal.md"},
		{ID: "G-001", Type: "goal", Status: "proposed", Path: "docs/goals/G-001.md"},
		{ID: "CH-001", Type: "change", Status: "open", Path: "changes/CH-001/proposal.md"},
	}}
	changes := OpenChanges(c)
	goals := ProposedGoals(c)
	if len(changes) != 2 || changes[0].ID != "CH-001" || changes[1].ID != "CH-002" {
		t.Fatalf("open changes = %#v", changes)
	}
	if len(goals) != 2 || goals[0].ID != "G-001" || goals[1].ID != "G-002" {
		t.Fatalf("proposed goals = %#v", goals)
	}
}

func TestAC190_UnitNegative_OrientationExcludesFinishedAndAcceptedArtifacts(t *testing.T) {
	c := &Corpus{Artifacts: []*Artifact{
		{ID: "G-001", Type: "goal", Status: "accepted", Path: "docs/goals/G-001.md"},
		{ID: "P-001", Type: "plan", Status: "completed", Path: "docs/plans/P-001.md"},
		{ID: "CH-001", Type: "imported-change", Status: "complete", Path: "docs/imported-changes/CH-001.md"},
	}}
	if got := OpenChanges(c); len(got) != 0 {
		t.Fatalf("open changes = %#v, want none", got)
	}
	if got := ProposedGoals(c); len(got) != 0 {
		t.Fatalf("proposed goals = %#v, want none", got)
	}
}
