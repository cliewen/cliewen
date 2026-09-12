package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const wallChecklist = ".github/cliewen-wall.md"

// AC-199 positive: init writes the checklist beside the caller, covering the
// whole PDR-021 boundary, and the caller points an adopter at it.
func TestAC199_UnitPositive_InitShipsTheMergeBoundaryChecklist(t *testing.T) {
	root, _ := runInto(t)
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(wallChecklist)))
	if err != nil {
		t.Fatalf("init did not write the checklist: %v", err)
	}
	for _, want := range []string{
		"Branch deletion is blocked",
		"Force pushes are blocked",
		"Changes arrive only through a pull request",
		"Merge commits are the only allowed merge method",
		"Review conversations must be resolved",
		"The `validate` check must pass",
		"Nobody is on the bypass list",
		"https://cliewen.dev/ci-wall",
	} {
		if !strings.Contains(string(data), want) {
			t.Errorf("the checklist does not say %q", want)
		}
	}
	caller, err := os.ReadFile(filepath.Join(root, ".github", "workflows", "clue.yml"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(caller), wallChecklist) {
		t.Error("the generated caller does not point to the checklist")
	}
}

// AC-199 negative: the checklist speaks in host-neutral terms rather than one
// host's menus, and migration never demands it from an existing adopter.
func TestAC199_UnitNegative_ChecklistIsHostNeutralAndNotAManagedCarrier(t *testing.T) {
	data, err := templates.ReadFile("templates/github/cliewen-wall.md")
	if err != nil {
		t.Fatal(err)
	}
	for _, menu := range []string{"Settings →", "Rulesets →", "New branch ruleset"} {
		if strings.Contains(string(data), menu) {
			t.Errorf("the checklist walks one host's menus (%q); that detail belongs to the guide", menu)
		}
	}
	carriers, err := ManagedCarrierFiles()
	if err != nil {
		t.Fatal(err)
	}
	if _, managed := carriers[wallChecklist]; managed {
		t.Error("the checklist is a managed carrier, so migration would demand it from every existing adopter")
	}
}
