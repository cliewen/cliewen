package corpus

import (
	"reflect"
	"strings"
	"testing"
)

func contextFiles() map[string]string {
	return map[string]string{
		"docs/goals/G-101-goal.md":              "---\nid: G-101\ntype: goal\nstatus: accepted\nlinks: []\ntitle: Goal\n---\n\n# Goal\n",
		"docs/plans/P-101-plan.md":              "---\nid: P-101\ntype: plan\nstatus: active\nlinks: [G-101]\ntitle: Plan\n---\n\n| M-101 | Do it | todo |\n",
		"docs/capabilities/CAP-101/README.md":   "---\nid: CAP-101\ntype: capability\nstatus: active\nlinks: [M-101]\ntitle: Capability\n---\n\n# Capability\n",
		"docs/capabilities/CAP-101/criteria.md": "---\nid: CAP-101-criteria\ntype: criteria\nstatus: active\nlinks: [CAP-101]\ntitle: Criteria\n---\n\n```gherkin\n@AC-101\nScenario: Works\n  Test-type: Unit\n```\n",
		"docs/analysis/AN-101-z.md":             "---\nid: AN-101\ntype: analysis\nstatus: active\nlinks: [G-101]\ntitle: Z\n---\n",
		"docs/analysis/AN-102-a.md":             "---\nid: AN-102\ntype: analysis\nstatus: active\nlinks: [G-101]\ntitle: A\n---\n",
		"docs/analysis/AN-999-unrelated.md":     "---\nid: AN-999\ntype: analysis\nstatus: active\nlinks: []\ntitle: Unrelated\n---\n",
	}
}

func TestAC053_UnitPositive_ContextFollowsTransitiveLinksAndEmbeddedIDs(t *testing.T) {
	c, issues := Scan(writeCorpus(t, contextFiles()))
	if len(issues) != 0 {
		t.Fatalf("scan: %v", issues)
	}

	got, err := Context(c, "AC-101")
	if err != nil {
		t.Fatal(err)
	}
	var paths []string
	for _, artifact := range got {
		paths = append(paths, artifact.Path)
	}
	want := []string{
		"docs/capabilities/CAP-101/criteria.md",
		"docs/capabilities/CAP-101/README.md",
		"docs/plans/P-101-plan.md",
		"docs/goals/G-101-goal.md",
	}
	if !reflect.DeepEqual(paths, want) {
		t.Fatalf("context paths = %v, want %v", paths, want)
	}

	fromMilestone, err := Context(c, "M-101")
	if err != nil {
		t.Fatal(err)
	}
	if fromMilestone[0].ID != "P-101" {
		t.Fatalf("milestone owner = %s, want P-101", fromMilestone[0].ID)
	}

	c.ByID["CAP-101"][0].Links = []string{"AN-102", "AN-101"}
	sameDepth, err := Context(c, "CAP-101")
	if err != nil {
		t.Fatal(err)
	}
	if sameDepth[1].Path != "docs/analysis/AN-101-z.md" || sameDepth[2].Path != "docs/analysis/AN-102-a.md" {
		t.Fatalf("same-depth artifacts not path ordered: %v, %v", sameDepth[1].Path, sameDepth[2].Path)
	}
}

func TestAC053_UnitNegative_ContextRejectsUnknownAndAmbiguousIDs(t *testing.T) {
	files := contextFiles()
	c, issues := Scan(writeCorpus(t, files))
	if len(issues) != 0 {
		t.Fatalf("scan: %v", issues)
	}
	if _, err := Context(c, "CAP-404"); err == nil || !strings.Contains(err.Error(), "CAP-404") {
		t.Fatalf("unknown ID error = %v", err)
	}

	files["docs/capabilities/CAP-101/other-criteria.md"] = strings.Replace(
		files["docs/capabilities/CAP-101/criteria.md"],
		"id: CAP-101-criteria",
		"id: CAP-101-other-criteria",
		1,
	)
	c, issues = Scan(writeCorpus(t, files))
	if len(issues) != 0 {
		t.Fatalf("scan ambiguous corpus: %v", issues)
	}
	if _, err := Context(c, "AC-101"); err == nil || !strings.Contains(err.Error(), "ambiguous") {
		t.Fatalf("ambiguous ID error = %v", err)
	}
}
