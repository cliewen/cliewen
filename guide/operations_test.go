package guide

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var nextAction = regexp.MustCompile(`(?s)## Next\n\n\[[^\]]+\]\([^\)]+\)\n\z`)

func TestAC036_UnitPositive_OperationsGuideStatesSupportedBoundary(t *testing.T) {
	requiredByPage := map[string][]string{
		"adoption.md":        {"classifies and counts the pair", "Cucumber", "Test-type: Human", "@draft", "cannot check that the acceptance brief supplies Human proof", "does not run the tests"},
		"getting-started.md": {"classified positive and negative evidence", "Cucumber", "Test-type: Human", "@draft", "cannot check that the brief supplies the proof", "does not run tests"},
		"operations.md":      {"Go test names", "Java and Kotlin executables whose own JUnit", "same executable", "Cucumber scenario tags", "(single-direction)", "Test-type: Human", "@draft", "does not run your tests", "does not update installed files in the background", "Keep the binary, generated skills, and CI caller on the same release", "Recover without bypassing the evidence", "`clue init` reports a skipped file", "`clue validate` fails", "CI rejects a transient workspace", "Do not delete a rule", "foreign-soil trials, not adoptions"},
		"what-is-cliewen.md": {"classified positive and negative evidence", "Cucumber", "Test-type: Human", "@draft", "does not execute tests"},
		"design.md":          {"acceptance criterion → acceptance evidence", "Cucumber", "Human-class", "@draft", "does not execute tests or inspect the pull request acceptance brief"},
		"methodology.md":     {"Acceptance evidence", "Cucumber", "Test-type: Human", "@draft", "does not execute tests"},
		"change-loop.md":     {"classified by that type and positive/negative direction", "Cucumber", "Test-type: Human", "@draft", "one supported reference"},
	}
	for page, required := range requiredByPage {
		content, err := os.ReadFile(page)
		if err != nil {
			t.Fatal(err)
		}
		for _, want := range required {
			if !strings.Contains(string(content), want) {
				t.Errorf("%s is missing current evidence-model carrier %q", page, want)
			}
		}
	}
}

func TestAC058_UnitPositive_GuideStatesPerExecutableJVMCarrier(t *testing.T) {
	requiredByPage := map[string][]string{
		"adoption.md":        {"per-executable JVM JUnit method tags", "class tags and unrelated methods cannot supply missing parts"},
		"getting-started.md": {"same Java or Kotlin executable"},
		"operations.md":      {"stable JVM named-executable form", "metadata split across methods do not count"},
		"what-is-cliewen.md": {"per-executable Java/Kotlin JUnit method tags", "metadata split across methods"},
		"design.md":          {"stable JVM named-executable fallback", "diagnoses ambiguous, class-level, or unsupported evidence syntax"},
		"methodology.md":     {"all three parts on one supported Java or Kotlin executable"},
		"change-loop.md":     {"same supported Java or Kotlin executable"},
	}
	for page, required := range requiredByPage {
		content, err := os.ReadFile(page)
		if err != nil {
			t.Fatal(err)
		}
		for _, fragment := range required {
			if !strings.Contains(string(content), fragment) {
				t.Errorf("%s is missing per-executable JVM carrier %q", page, fragment)
			}
		}
	}
}

func TestAC058_UnitNegative_GuideRejectsFileLevelJVMEvidence(t *testing.T) {
	pages, err := filepath.Glob("*.md")
	if err != nil {
		t.Fatal(err)
	}
	for _, page := range pages {
		content, err := os.ReadFile(page)
		if err != nil {
			t.Fatal(err)
		}
		for _, stale := range []string{
			"JVM per-method purpose enforcement via an ArchUnit rule",
			"JVM tags are file-level",
			"file-level JVM harvesting",
		} {
			if strings.Contains(string(content), stale) {
				t.Errorf("%s still carries stale JVM evidence claim %q", page, stale)
			}
		}
	}
}

func TestAC036_UnitNegative_MultiplePrimaryActionsAreRejected(t *testing.T) {
	pages, err := filepath.Glob("*.md")
	if err != nil {
		t.Fatal(err)
	}
	for _, page := range pages {
		content, err := os.ReadFile(page)
		if err != nil {
			t.Fatal(err)
		}
		if !nextAction.Match(content) {
			t.Errorf("%s must end with exactly one Next action", page)
		}
	}
	if nextAction.MatchString("## Next\n\n[first](./one) and [second](./two)\n") {
		t.Fatal("multiple actions must not satisfy the next-action rule")
	}
}

// The evidence-model claims repaired by PDR-019 are stable enough to anchor,
// so no guide page may drift back to the pre-v0.9 contract.
func TestSanity_NoGuidePageCarriesAStaleEvidenceModelClaim(t *testing.T) {
	pages, err := filepath.Glob("*.md")
	if err != nil {
		t.Fatal(err)
	}
	stale := []string{
		"does not count or classify the pair",
		"does not distinguish or count the positive and negative pair",
		"Every active criterion gets focused positive and negative tests",
		"every active acceptance criterion must have focused positive and negative tests",
		"An acceptance criterion whose proof is inherently human has no equivalent home yet",
	}
	for _, page := range pages {
		content, err := os.ReadFile(page)
		if err != nil {
			t.Fatal(err)
		}
		for _, claim := range stale {
			if strings.Contains(string(content), claim) {
				t.Errorf("%s still carries stale evidence-model claim %q", page, claim)
			}
		}
	}
}

func TestAC040_CIWallMakesKnownFindingsDurable(t *testing.T) {
	content, err := os.ReadFile("ci-wall.md")
	if err != nil {
		t.Fatal(err)
	}
	for _, required := range []string{
		"Require conversation resolution before merging",
		"known agent-review findings remain blocking",
		"agents can warn about unresolved findings, but neither can block integration",
	} {
		if !strings.Contains(string(content), required) {
			t.Errorf("CI wall guide is missing durable-finding rule %q", required)
		}
	}
}

func TestAC132_UnitPositive_ChangeLoopCoordinatesOnlyTheSharedPR(t *testing.T) {
	content, err := os.ReadFile("change-loop.md")
	if err != nil {
		t.Fatal(err)
	}
	for _, required := range []string{
		"Separate authors keep separate branches from accepted `main`",
		"Any agent asked to fix one becomes the updater for that turn",
		"commits and pushes the repair with the turn that made it",
		"which returns a ready pull request to draft",
		"If another updater moved the head",
		"merging newer accepted `main` into the branch with a normal push instead of rewriting hosted history",
		"the agent says the PR is not merge-ready",
	} {
		if !strings.Contains(string(content), required) {
			t.Errorf("change-loop guide is missing multi-agent handoff rule %q", required)
		}
	}
}

func TestSanity_FullChangeMergeModeIsExplicitAcrossGuidanceAndProbe(t *testing.T) {
	requiredByPage := map[string][]string{
		"operations.md":      {"Full-change merge", "Human-controlled merge commits only", "disable **squash and merge** and **rebase and merge**", "outside Cliewen's supported full-change adoption path"},
		"change-loop.md":     {"human-controlled merge commit", "disable squash and rebase-and-merge", "hosted history is never rebased or rewritten"},
		"ci-wall.md":         {"Create a merge commit", "Squash and merge", "Rebase and merge", ".allow_merge_commit", ".allow_squash_merge", ".allow_rebase_merge", "not ready for a full Cliewen change"},
		"getting-started.md": {"protected default branch for human-controlled merge commits only"},
	}
	for page, required := range requiredByPage {
		content, err := os.ReadFile(page)
		if err != nil {
			t.Fatal(err)
		}
		for _, want := range required {
			if !strings.Contains(string(content), want) {
				t.Errorf("%s is missing merge-history contract %q", page, want)
			}
		}
	}

	// The forge's pull-request record is not the system of record. This exact
	// sentence made that claim before CH-090 and must not come back.
	for _, page := range []string{"operations.md", "change-loop.md", "ci-wall.md", "getting-started.md", "what-is-cliewen.md"} {
		content, err := os.ReadFile(page)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(content), "the merged pull request itself is the historical record") {
			t.Errorf("%s again treats the pull request as the provenance archive", page)
		}
	}

	// Guarding a list of equivalence phrasings nobody has written yet catches
	// nothing. The invariant that does hold: a page cannot raise squash or
	// rebase-and-merge without also placing them outside the support boundary,
	// so no page can describe them as an equivalent way to accept a full change.
	pages, err := filepath.Glob("*.md")
	if err != nil {
		t.Fatal(err)
	}
	boundaryMarkers := []string{"unsupported", "disable", "outside", "not ready"}
	for _, page := range pages {
		content, err := os.ReadFile(page)
		if err != nil {
			t.Fatal(err)
		}
		lowered := strings.ToLower(string(content))
		if !strings.Contains(lowered, "squash") && !strings.Contains(lowered, "rebase and merge") && !strings.Contains(lowered, "rebase-and-merge") {
			continue
		}
		carriesBoundary := false
		for _, marker := range boundaryMarkers {
			if strings.Contains(lowered, marker) {
				carriesBoundary = true
				break
			}
		}
		if !carriesBoundary {
			t.Errorf("%s names squash or rebase-and-merge without stating that they fall outside the supported full-change merge mode", page)
		}
	}
}
