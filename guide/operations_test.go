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
		"operations.md":      {"Go test names", "Java and Kotlin JUnit", "Cucumber scenario tags", "(single-direction)", "Test-type: Human", "@draft", "does not run your tests", "does not update installed files in the background", "Recover without bypassing the evidence", "`clue init` reports a skipped file", "`clue validate` fails", "CI rejects a transient workspace", "Do not delete a rule", "foreign-soil trials, not adoptions"},
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
	for _, stale := range []string{
		"does not count or classify the pair",
		"does not distinguish or count the positive and negative pair",
		"Every active criterion gets focused positive and negative tests",
		"every active acceptance criterion must have focused positive and negative tests",
		"An acceptance criterion whose proof is inherently human has no equivalent home yet",
	} {
		for _, page := range pages {
			content, err := os.ReadFile(page)
			if err != nil {
				t.Fatal(err)
			}
			if strings.Contains(string(content), stale) {
				t.Errorf("%s still carries stale evidence-model claim %q", page, stale)
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

func TestAC041_ChangeLoopCoordinatesOnlyTheSharedPR(t *testing.T) {
	content, err := os.ReadFile("change-loop.md")
	if err != nil {
		t.Fatal(err)
	}
	for _, required := range []string{
		"Separate authors keep separate branches from accepted `main`",
		"Any agent asked to fix one becomes the updater for that turn",
		"pushes without force",
		"If another updater moved the head",
		"merge newer accepted `main` into its branch with a normal push instead of rewriting hosted history",
		"the agent says the PR is not merge-ready",
	} {
		if !strings.Contains(string(content), required) {
			t.Errorf("change-loop guide is missing multi-agent handoff rule %q", required)
		}
	}
}
