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
