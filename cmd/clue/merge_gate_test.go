package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAC042_AcceptanceBriefCarriesHumanMergeQuestions(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", ".github", "pull_request_template.md"))
	if err != nil {
		t.Fatal(err)
	}
	for _, required := range []string{
		"## Acceptance brief",
		"Plan item and whether it remains wanted",
		"verifies-something-adjacent",
		"What becomes binding on merge",
		"one screen",
	} {
		if !strings.Contains(string(data), required) {
			t.Errorf("acceptance brief is missing %q", required)
		}
	}
}

func TestAC042_CIRejectsUnfilledAcceptanceBrief(t *testing.T) {
	for _, rel := range []string{
		filepath.Join("..", "..", ".github", "workflows", "ci.yml"),
		filepath.Join("..", "..", "internal", "scaffold", "templates", "github", "workflows", "clue.yml"),
	} {
		data, err := os.ReadFile(rel)
		if err != nil {
			t.Fatal(err)
		}
		for _, required := range []string{"Require a completed acceptance brief", "## Acceptance brief", "<!-- REQUIRED"} {
			if !strings.Contains(string(data), required) {
				t.Errorf("%s does not enforce %q", rel, required)
			}
		}
	}
}

func TestAC048_UnitPositive_BriefNamesHumanClassCriteriaAsTheirProof(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", ".github", "pull_request_template.md"))
	if err != nil {
		t.Fatal(err)
	}
	for _, required := range []string{
		"Test-type: Human",
		"this line is that criterion's proof, not a separate confirmation",
	} {
		if !strings.Contains(string(data), required) {
			t.Errorf("acceptance brief is missing %q", required)
		}
	}
}

func TestAC048_UnitNegative_HumanProofRidesTheExistingCriteriaLine(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", ".github", "pull_request_template.md"))
	if err != nil {
		t.Fatal(err)
	}
	if n := strings.Count(string(data), "Added or changed criteria"); n != 1 {
		t.Fatalf("expected exactly one criteria bullet in the acceptance brief, got %d", n)
	}
	if strings.Contains(string(data), "## Human-class confirmation") {
		t.Fatalf("Human-class proof must ride the existing criteria bullet, not a separate section")
	}
}
