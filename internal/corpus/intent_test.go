package corpus

import (
	"strings"
	"testing"
)

// visionBody is a well-formed vision: the judge checks nothing about what it
// says, so the content here is deliberately trivial.
const visionBody = "---\nid: VIS-001\ntype: vision\nstatus: active\nlinks: []\ntitle: A product\n---\n\n# VIS-001\n\nWhat it is for.\n"

const useCaseBody = "---\nid: UC-001\ntype: use-case\nstatus: active\nlinks: [G-001, CAP-001]\ntitle: A journey\n---\n\n# UC-001\n\n## Actors\n\nSomeone.\n\n## Trigger\n\nThey ask.\n\n## Main flow\n\n1. It happens.\n\n## Outcome\n\nIt happened.\n"

const capabilityBody = "---\nid: CAP-001\ntype: capability\nstatus: active\nlinks: [G-001]\ntitle: A capability\ngoal: G-001\n---\n\n# CAP-001\n"

// intentFiles is validFiles plus a capability, a vision, and one use case, with
// every index row the layout rules require.
func intentFiles(extra map[string]string) map[string]string {
	base := with(validFiles, map[string]string{
		"docs/README.md":                            "# Corpus\n\n<!-- clue:index:start -->\n- [vision.md](vision.md)\n- [goals/](goals/README.md)\n- [plans/](plans/README.md)\n- [use-cases/](use-cases/README.md)\n- [capabilities/](capabilities/README.md)\n<!-- clue:index:end -->\n",
		"docs/vision.md":                            visionBody,
		"docs/capabilities/README.md":               "# Capabilities\n\n<!-- clue:index:start -->\n- [CAP-001](CAP-001-thing/README.md)\n<!-- clue:index:end -->\n",
		"docs/capabilities/CAP-001-thing/README.md": capabilityBody,
		"docs/use-cases/README.md":                  "# Use cases\n\n<!-- clue:index:start -->\n- [UC-001](UC-001-a-journey.md)\n<!-- clue:index:end -->\n",
		"docs/use-cases/UC-001-a-journey.md":        useCaseBody,
	})
	return with(base, extra)
}

// AC-162: one vision, at one address, checked for form and for nothing it says.
func TestAC162_UnitPositive_OneVisionAtItsAddressPasses(t *testing.T) {
	if issues := run(t, intentFiles(nil), true); len(issues) != 0 {
		t.Fatalf("expected a green corpus with one vision, got %v", issues)
	}
	// A corpus that states no direction at all is equally valid: nothing
	// requires the artifact to exist (ADR-067).
	if issues := run(t, validFiles, true); len(issues) != 0 {
		t.Fatalf("expected a green corpus with no vision, got %v", issues)
	}
}

func TestAC162_UnitNegative_SecondVisionWrongAddressAndBootstrapFail(t *testing.T) {
	second := strings.Replace(visionBody, "id: VIS-001", "id: VIS-002", 1)
	issues := run(t, intentFiles(map[string]string{
		"docs/goals/VIS-002-rival.md": second,
		"docs/goals/README.md":        "# Goals\n\n<!-- clue:index:start -->\n- [G-001](G-001-first.md)\n- [VIS-002](VIS-002-rival.md)\n<!-- clue:index:end -->\n",
	}), false)
	assertIssue(t, issues, "a corpus has one vision")
	assertIssue(t, issues, "a vision lives at docs/vision.md")

	// An artifact squatting the vision's address is its own defect, because
	// its repair is to move that artifact rather than to move a vision.
	issues = run(t, intentFiles(map[string]string{
		"docs/vision.md": "---\nid: G-002\ntype: goal\nstatus: proposed\nlinks: []\ntitle: Not a vision\n---\n\n# G-002\n",
	}), false)
	assertIssue(t, issues, "this address belongs to the corpus vision")

	// An unreplaced bootstrap is the one content rule, and it is about a
	// marker rather than about meaning.
	issues = run(t, intentFiles(map[string]string{
		"docs/vision.md": visionBody + "\n" + VisionBootstrapMarker + "\n",
	}), false)
	assertIssue(t, issues, "still the scaffold bootstrap")
}

// AC-163: a use case is checked for form when it exists, and never required.
func TestAC163_UnitPositive_UseCasesAreCheckedForFormAndNeverRequired(t *testing.T) {
	if issues := run(t, intentFiles(nil), true); len(issues) != 0 {
		t.Fatalf("expected a well-formed use case to pass, got %v", issues)
	}
	// The goal and the capability in this corpus have no use case naming
	// them once it is removed, and nothing reports that.
	noUseCases := intentFiles(nil)
	delete(noUseCases, "docs/use-cases/UC-001-a-journey.md")
	delete(noUseCases, "docs/use-cases/README.md")
	noUseCases["docs/README.md"] = strings.Replace(noUseCases["docs/README.md"], "- [use-cases/](use-cases/README.md)\n", "", 1)
	if issues := run(t, noUseCases, true); len(issues) != 0 {
		t.Fatalf("a corpus with no use cases must be valid, got %v", issues)
	}
}

func TestAC163_UnitNegative_MisplacedMisnamedUnlinkedOrIncompleteUseCaseFails(t *testing.T) {
	issues := run(t, intentFiles(map[string]string{
		"docs/goals/UC-001-a-journey.md":     useCaseBody,
		"docs/use-cases/UC-001-a-journey.md": "",
		"docs/goals/README.md":               "# Goals\n\n<!-- clue:index:start -->\n- [G-001](G-001-first.md)\n- [UC-001](UC-001-a-journey.md)\n<!-- clue:index:end -->\n",
	}), false)
	assertIssue(t, issues, "use cases live in docs/use-cases")

	issues = run(t, intentFiles(map[string]string{
		"docs/use-cases/UC-001-a-journey.md": "",
		"docs/use-cases/journey.md":          useCaseBody,
		"docs/use-cases/README.md":           "# Use cases\n\n<!-- clue:index:start -->\n- [UC-001](journey.md)\n<!-- clue:index:end -->\n",
	}), false)
	assertIssue(t, issues, "filename matching their id")

	issues = run(t, intentFiles(map[string]string{
		"docs/use-cases/UC-001-a-journey.md": strings.Replace(useCaseBody, "links: [G-001, CAP-001]", "links: [CAP-001]", 1),
	}), false)
	assertIssue(t, issues, "names no goal it serves")

	issues = run(t, intentFiles(map[string]string{
		"docs/use-cases/UC-001-a-journey.md": strings.Replace(useCaseBody, "links: [G-001, CAP-001]", "links: [G-001]", 1),
	}), false)
	assertIssue(t, issues, "names no capability it crosses")

	issues = run(t, intentFiles(map[string]string{
		"docs/use-cases/UC-001-a-journey.md": strings.Replace(useCaseBody, "## Trigger\n\nThey ask.\n\n", "", 1),
	}), false)
	assertIssue(t, issues, "## Trigger")
}

// AC-164: the derived state, without a coverage figure. The command's own
// rendering is covered beside the command.
func TestAC164_UnitPositive_IntentStateNamesTheVisionAndTheJourneys(t *testing.T) {
	c, issues := Scan(writeCorpus(t, intentFiles(map[string]string{
		"docs/vision.md": strings.Replace(visionBody, "links: []", "links: []\nprovenance: inferred\nreversal-cost: low", 1),
	})))
	if len(issues) != 0 {
		t.Fatal(issues)
	}
	state := Intent(c)
	if !state.Vision.Present || state.Vision.ID != "VIS-001" || !state.Vision.Inferred {
		t.Fatalf("unexpected vision state: %+v", state.Vision)
	}
	if len(state.UseCases) != 1 || state.UseCases[0].ID != "UC-001" || len(state.UseCases[0].Capabilities) != 1 {
		t.Fatalf("unexpected use-case state: %+v", state.UseCases)
	}
}

func TestAC164_UnitNegative_AnAbsentVisionIsAStateAndNotAnIssue(t *testing.T) {
	c, issues := Scan(writeCorpus(t, validFiles))
	if len(issues) != 0 {
		t.Fatal(issues)
	}
	if state := Intent(c); state.Vision.Present || len(state.UseCases) != 0 {
		t.Fatalf("expected an empty intent state, got %+v", state)
	}
	if issues := Validate(c, Options{}); len(issues) != 0 {
		t.Fatalf("an absent vision must not be an issue, got %v", issues)
	}
}

// AC-165: the naming that lets a capability reach the journey governing it.
func TestAC165_UnitPositive_UseCasesNamingAnArtifactAreFoundInPathOrder(t *testing.T) {
	c, issues := Scan(writeCorpus(t, intentFiles(map[string]string{
		"docs/use-cases/UC-002-another.md": strings.Replace(strings.Replace(useCaseBody, "id: UC-001", "id: UC-002", 1), "A journey", "Another journey", 1),
		"docs/use-cases/README.md":         "# Use cases\n\n<!-- clue:index:start -->\n- [UC-001](UC-001-a-journey.md)\n- [UC-002](UC-002-another.md)\n<!-- clue:index:end -->\n",
	})))
	if len(issues) != 0 {
		t.Fatal(issues)
	}
	naming := UseCasesNaming(c, "CAP-001")
	if len(naming) != 2 || naming[0].ID != "UC-001" || naming[1].ID != "UC-002" {
		t.Fatalf("expected both use cases in path order, got %+v", naming)
	}
}

func TestAC165_UnitNegative_AnArtifactNoUseCaseNamesIsNamedByNothing(t *testing.T) {
	c, issues := Scan(writeCorpus(t, intentFiles(nil)))
	if len(issues) != 0 {
		t.Fatal(issues)
	}
	if naming := UseCasesNaming(c, "P-001"); len(naming) != 0 {
		t.Fatalf("expected no use case to name the plan, got %+v", naming)
	}
	// A use case never names itself, so a slice rooted at one gains nothing.
	if naming := UseCasesNaming(c, "UC-001"); len(naming) != 0 {
		t.Fatalf("expected a use case not to name itself, got %+v", naming)
	}
}
