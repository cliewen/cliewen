package guide

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var nextAction = regexp.MustCompile(`(?s)## Next\n\n\[[^\]]+\]\([^\)]+\)\n\z`)

func TestAC036_OperationsGuideStatesSupportedBoundary(t *testing.T) {
	content, err := os.ReadFile("operations.md")
	if err != nil {
		t.Fatal(err)
	}
	for _, required := range []string{
		"Go test names", "Java and Kotlin JUnit", "does not run your tests", "does not update installed files in the background", "Keep the binary, generated skills, and vendored CI binary on the same release", "Recover without bypassing the evidence", "foreign-soil trials, not adoptions",
	} {
		if !strings.Contains(string(content), required) {
			t.Errorf("operations guide is missing %q", required)
		}
	}
}

func TestAC036_EachGuidePageEndsWithOnePrimaryAction(t *testing.T) {
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
