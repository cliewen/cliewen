package corpus

import "testing"

func TestAC135_UnitPositive_ReadCostBacklogsNameMultiDocumentArtifactsAndWideSlices(t *testing.T) {
	root := &Artifact{ID: "G-001", Path: "docs/goals/G-001-root.md", Links: []string{"G-002", "G-003", "G-004", "G-005", "G-006", "G-007", "G-008", "G-009"}}
	artifacts := []*Artifact{
		root,
		{ID: "AN-001", Path: "docs/analysis/AN-001.md", Type: "analysis", Status: "active", Body: "# First document\n\n# Second document\n\n```markdown\n# Example only\n```\n"},
		{ID: "P-001", Path: "docs/plans/P-001.md", Type: "plan", Status: "completed", Body: "# Frozen first\n\n# Frozen second\n"},
	}
	for _, id := range root.Links {
		artifacts = append(artifacts, &Artifact{ID: id, Path: "docs/goals/" + id + ".md"})
	}
	c := &Corpus{Artifacts: artifacts}

	multi := MultiDocumentBacklog(c)
	if len(multi) != 1 || multi[0].Artifact.ID != "AN-001" || multi[0].Documents != 2 {
		t.Fatalf("multi-document backlog = %#v, want active AN-001 with two documents", multi)
	}
	wide := ContextSliceBudgetBacklog(c)
	if len(wide) != 1 || wide[0].Identity != "G-001" || wide[0].Artifacts != DefaultContextSliceBudget+1 {
		t.Fatalf("context-budget backlog = %#v, want G-001 with %d artifacts", wide, DefaultContextSliceBudget+1)
	}
}

func TestAC135_UnitNegative_ReadCostBacklogsIgnoreExamplesFrozenPlansAndSlicesAtBudget(t *testing.T) {
	root := &Artifact{ID: "G-001", Path: "docs/goals/G-001-root.md", Links: []string{"G-002", "G-003", "G-004", "G-005", "G-006", "G-007", "G-008"}}
	artifacts := []*Artifact{
		root,
		{ID: "AN-001", Path: "docs/analysis/AN-001.md", Type: "analysis", Status: "active", Body: "# One document\n\n```markdown\n# Example only\n```\n"},
		{ID: "P-001", Path: "docs/plans/P-001.md", Type: "plan", Status: "completed", Body: "# Frozen first\n\n# Frozen second\n"},
		{ID: "P-002", Path: "docs/plans/P-002.md", Type: "plan", Status: "completed", Links: []string{"G-001", "G-002", "G-003", "G-004", "G-005", "G-006", "G-007", "G-008", "P-001"}},
	}
	for _, id := range root.Links {
		artifacts = append(artifacts, &Artifact{ID: id, Path: "docs/goals/" + id + ".md"})
	}
	c := &Corpus{Artifacts: artifacts}

	if multi := MultiDocumentBacklog(c); len(multi) != 0 {
		t.Fatalf("multi-document backlog = %#v, want none", multi)
	}
	if wide := ContextSliceBudgetBacklog(c); len(wide) != 0 {
		t.Fatalf("context-budget backlog = %#v, want none", wide)
	}
}
