package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func readRepoFile(t *testing.T, parts ...string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(append([]string{"..", ".."}, parts...)...))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func TestAC245_UnitPositive_TrackedFormAsksForTheChangeWithoutABrief(t *testing.T) {
	form := readRepoFile(t, ".github", "PULL_REQUEST_TEMPLATE", "tracked.md")
	for _, required := range []string{"## Summary", "## Verification", "## Cliewen proposal", "## Traceability and Decisions"} {
		if !strings.Contains(form, required) {
			t.Errorf("tracked form is missing %q", required)
		}
	}
	if strings.Contains(form, "Acceptance brief") || strings.Contains(form, "acceptance brief") {
		t.Error("tracked form still carries an acceptance brief")
	}
}

func TestAC245_UnitNegative_NoWorkflowReadsThePullRequestBody(t *testing.T) {
	for _, rel := range []string{"ci.yml", "clue-validation.yml"} {
		workflow := readRepoFile(t, ".github", "workflows", rel)
		for _, banned := range []string{"PR_BODY", "pull_request.body", "Require a completed acceptance brief"} {
			if strings.Contains(workflow, banned) {
				t.Errorf("%s still gates on the pull-request body (%q)", rel, banned)
			}
		}
	}
}

func TestAC246_UnitPositive_TrackedChecklistNamesTheProposalsHumanEvidence(t *testing.T) {
	form := readRepoFile(t, ".github", "PULL_REQUEST_TEMPLATE", "tracked.md")
	for _, required := range []string{"Human proof recorded in the proposal's `## Human evidence` section", "what was not observed"} {
		if !strings.Contains(form, required) {
			t.Errorf("tracked checklist is missing %q", required)
		}
	}
}

func TestAC246_UnitNegative_HumanProofHasNoSeparateConfirmationSection(t *testing.T) {
	form := readRepoFile(t, ".github", "PULL_REQUEST_TEMPLATE", "tracked.md")
	if n := strings.Count(form, "## Human evidence"); n != 1 {
		t.Fatalf("expected the proposal's human evidence section to be named once, got %d", n)
	}
	if strings.Contains(form, "\n## Human evidence") || strings.Contains(form, "## Human-class confirmation") {
		t.Fatal("Human-class proof belongs in the proposal, not in a pull-request section")
	}
}
