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
		"docs/capabilities/CAP-101/README.md":   "---\nid: CAP-101\ntype: capability\nstatus: active\nlinks: [M-101]\ntitle: Capability\ngoal: G-101\n---\n\n# Capability\n",
		"docs/capabilities/CAP-101/criteria.md": "---\nid: CAP-101-criteria\ntype: criteria\nstatus: active\nlinks: [CAP-101]\ntitle: Criteria\n---\n\n```gherkin\n@AC-101\nScenario: Works\n  Test-type: Unit\n```\n",
		"docs/analysis/AN-101-z.md":             "---\nid: AN-101\ntype: analysis\nstatus: active\nlinks: [G-101]\ntitle: Z\n---\n",
		"docs/analysis/AN-102-a.md":             "---\nid: AN-102\ntype: analysis\nstatus: active\nlinks: [G-101]\ntitle: A\n---\n",
		"docs/analysis/AN-999-unrelated.md":     "---\nid: AN-999\ntype: analysis\nstatus: active\nlinks: []\ntitle: Unrelated\n---\n",
	}
}

func TestAC133_UnitPositive_ContextFollowsTransitiveLinksAndEmbeddedIDs(t *testing.T) {
	c, issues := Scan(writeCorpus(t, contextFiles()))
	if len(issues) != 0 {
		t.Fatalf("scan: %v", issues)
	}

	got, _, unfollowed, err := Context(c, "AC-101", ContextOptions{Depth: DepthAll})
	if err != nil {
		t.Fatal(err)
	}
	if len(unfollowed) != 0 {
		t.Fatalf("unfollowed edges in a sound corpus: %v", unfollowed)
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

	fromMilestone, _, _, err := Context(c, "M-101", ContextOptions{Depth: DepthAll})
	if err != nil {
		t.Fatal(err)
	}
	if fromMilestone[0].ID != "P-101" {
		t.Fatalf("milestone owner = %s, want P-101", fromMilestone[0].ID)
	}

	c.ByID["CAP-101"][0].Links = []string{"AN-102", "AN-101"}
	sameDepth, _, _, err := Context(c, "CAP-101", ContextOptions{Depth: DepthAll})
	if err != nil {
		t.Fatal(err)
	}
	if sameDepth[1].Path != "docs/analysis/AN-101-z.md" || sameDepth[2].Path != "docs/analysis/AN-102-a.md" {
		t.Fatalf("same-depth artifacts not path ordered: %v, %v", sameDepth[1].Path, sameDepth[2].Path)
	}
}

// A milestone is declared by its row in a plan's milestone table. Another plan
// that only mentions the ID in prose — corpus-global numbering continues from
// the previous campaign — is referencing it, not declaring it.
func TestAC133_UnitPositive_ContextTreatsProseMilestoneMentionsAsReferences(t *testing.T) {
	files := contextFiles()
	files["docs/plans/P-102-successor.md"] = "---\nid: P-102\ntype: plan\nstatus: active\nlinks: [G-101]\ntitle: Successor\n---\n\nMilestone numbering continues from P-101 (M-101).\n\n| M-102 | Next | todo |\n"
	c, issues := Scan(writeCorpus(t, files))
	if len(issues) != 0 {
		t.Fatalf("scan: %v", issues)
	}

	got, _, _, err := Context(c, "M-101", ContextOptions{Depth: DepthAll})
	if err != nil {
		t.Fatalf("prose mention made a declared milestone unresolvable: %v", err)
	}
	if got[0].ID != "P-101" {
		t.Fatalf("milestone owner = %s, want P-101", got[0].ID)
	}
}

// An edge the slice cannot follow is reported, not fatal: focused reading has
// to keep working on exactly the corpus a reader is repairing.
func TestAC133_UnitPositive_ContextReportsUnfollowedEdgesAndKeepsSlicing(t *testing.T) {
	files := contextFiles()
	files["docs/analysis/AN-101-z.md"] = "---\nid: AN-101\ntype: analysis\nstatus: active\nlinks: [G-101, CAP-404]\ntitle: Z\n---\n"
	c, issues := Scan(writeCorpus(t, files))
	if len(issues) != 0 {
		t.Fatalf("scan: %v", issues)
	}

	got, _, unfollowed, err := Context(c, "AN-101", ContextOptions{Depth: DepthAll})
	if err != nil {
		t.Fatalf("a dangling link ended the slice: %v", err)
	}
	var paths []string
	for _, artifact := range got {
		paths = append(paths, artifact.Path)
	}
	want := []string{"docs/analysis/AN-101-z.md", "docs/goals/G-101-goal.md"}
	if !reflect.DeepEqual(paths, want) {
		t.Fatalf("context paths = %v, want %v", paths, want)
	}
	if len(unfollowed) != 1 || !strings.Contains(unfollowed[0].Msg, "CAP-404") {
		t.Fatalf("unfollowed edge not reported: %v", unfollowed)
	}
}

func TestAC133_UnitNegative_ContextRejectsUnknownAndAmbiguousIDs(t *testing.T) {
	files := contextFiles()
	c, issues := Scan(writeCorpus(t, files))
	if len(issues) != 0 {
		t.Fatalf("scan: %v", issues)
	}
	if _, _, _, err := Context(c, "CAP-404", ContextOptions{Depth: DepthAll}); err == nil || !strings.Contains(err.Error(), "CAP-404") {
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
	if _, _, _, err := Context(c, "AC-101", ContextOptions{Depth: DepthAll}); err == nil || !strings.Contains(err.Error(), "ambiguous") {
		t.Fatalf("ambiguous ID error = %v", err)
	}
}

// The default slice is bounded. A corpus densifies as campaigns close, and an
// unbounded closure from a hub artifact stops being a slice and becomes the
// repository; the bound is what keeps `clue context` answering the question it
// was built for.
func TestAC133_UnitPositive_ContextBoundsTheSliceAndNamesWhatItHeldBack(t *testing.T) {
	c, issues := Scan(writeCorpus(t, contextFiles()))
	if len(issues) != 0 {
		t.Fatalf("scan: %v", issues)
	}

	// CAP-101 -> M-101 (P-101) -> G-101 is three hops of reachable graph.
	got, frontier, _, err := Context(c, "CAP-101", ContextOptions{Depth: 1})
	if err != nil {
		t.Fatal(err)
	}
	var paths []string
	for _, artifact := range got {
		paths = append(paths, artifact.Path)
	}
	want := []string{"docs/capabilities/CAP-101/README.md", "docs/plans/P-101-plan.md"}
	if !reflect.DeepEqual(paths, want) {
		t.Fatalf("depth-1 slice = %v, want %v", paths, want)
	}
	if len(frontier) != 1 || frontier[0].Artifact.ID != "G-101" || frontier[0].Hops != 2 {
		t.Fatalf("the bound must name what it held back, got %v", frontier)
	}

	// Depth 0 is the root alone, and it is not the zero-value default by
	// accident: an unset bound has to mean something safe.
	rootOnly, frontier, _, err := Context(c, "CAP-101", ContextOptions{Depth: 0})
	if err != nil {
		t.Fatal(err)
	}
	if len(rootOnly) != 1 || rootOnly[0].ID != "CAP-101" || len(frontier) != 2 {
		t.Fatalf("depth-0 slice = %v, frontier = %v", rootOnly, frontier)
	}

	// Widening restores the whole closure and empties the frontier.
	widened, frontier, _, err := Context(c, "CAP-101", ContextOptions{Depth: DepthAll})
	if err != nil {
		t.Fatal(err)
	}
	if len(widened) != 3 || len(frontier) != 0 {
		t.Fatalf("--depth=all must reach everything and hold nothing back: %d artifacts, %d frontier", len(widened), len(frontier))
	}
}

// An edge leaving an artifact the reader never sees describes a part of the
// graph this slice does not cover. Reporting it turns a bounded read into
// someone else's backlog, so the bound silences it while `clue validate`
// stays the judge of graph health.
func TestAC133_UnitNegative_ContextReportsNoUnfollowedEdgeFromAnExcludedArtifact(t *testing.T) {
	files := contextFiles()
	files["docs/goals/G-101-goal.md"] = "---\nid: G-101\ntype: goal\nstatus: accepted\nlinks: [CAP-404]\ntitle: Goal\n---\n\n# Goal\n"
	c, issues := Scan(writeCorpus(t, files))
	if len(issues) != 0 {
		t.Fatalf("scan: %v", issues)
	}

	// G-101 sits two hops out, beyond a depth-1 bound.
	if _, _, unfollowed, err := Context(c, "CAP-101", ContextOptions{Depth: 1}); err != nil || len(unfollowed) != 0 {
		t.Fatalf("a dangling edge outside the slice was reported: %v (err %v)", unfollowed, err)
	}
	if _, _, unfollowed, err := Context(c, "CAP-101", ContextOptions{Depth: DepthAll}); err != nil || len(unfollowed) != 1 {
		t.Fatalf("a dangling edge inside the slice must be reported: %v (err %v)", unfollowed, err)
	}
}
