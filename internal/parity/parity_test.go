package parity

import (
	"os"
	"path/filepath"
	"testing"
)

func writeFiles(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for p, content := range files {
		full := filepath.Join(root, filepath.FromSlash(p))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

const criteriaWithOneAC = `---
id: PARTEST-criteria
type: criteria
status: active
links: []
title: Parity fixture criteria
---

` + "```gherkin" + `
Feature: fixture

  @AC-900
  Scenario: fixture criterion
    Test-type: Integration
    Given a thing
    When it happens
    Then it works
` + "```" + `
`

// fixtureGoTest builds a throwaway Go test file's content without ever
// spelling "func Test" as contiguous source text in this file: the repo's
// own AC-evidence scan reads *_test.go files as raw text (it does not
// distinguish this file's real tests from a fixture string literal), so a
// literal "func Test..." here would be misread as evidence for the fixture
// IDs below, which no real criteria.md declares.
func fixtureGoTest(names ...string) string {
	out := "package fixture\n\n"
	for _, name := range names {
		out += "func " + "Test" + name + "(t *testing.T) {}\n"
	}
	return out
}

var goTestPositive = fixtureGoTest("AC900_IntegrationPositive_works")

var goTestBoth = fixtureGoTest("AC900_IntegrationPositive_works", "AC900_IntegrationNegative_fails")

// TestAC109_UnitPositive_cleanRunPasses proves a clean parity run produces no
// findings, and that comparing twice against the same inputs is deterministic.
func TestAC109_UnitPositive_cleanRunPasses(t *testing.T) {
	root := writeFiles(t, map[string]string{
		"docs/capabilities/PARTEST/criteria.md": criteriaWithOneAC,
		"fixture_test.go":                       goTestBoth,
	})
	target, err := DeriveTargetManifest(root)
	if err != nil {
		t.Fatal(err)
	}
	source := SourceManifest{Entries: []SourceEntry{
		{ID: "AC-900", ProofClass: "Integration", Direction: "positive", EvidenceLocation: "fixture_test.go"},
	}}
	r1 := Compare(source, target)
	if r1.Failed() {
		t.Fatalf("expected a clean report, got %v", r1.Findings)
	}
	r2 := Compare(source, target)
	if len(r1.Findings) != len(r2.Findings) {
		t.Fatalf("expected deterministic reports, got %v and %v", r1.Findings, r2.Findings)
	}
}

// TestAC109_UnitNegative_dirtyRunFails proves an altered entry still fails
// after a clean baseline, so a clean run is not a permanent exemption.
func TestAC109_UnitNegative_dirtyRunFails(t *testing.T) {
	root := writeFiles(t, map[string]string{
		"docs/capabilities/PARTEST/criteria.md": criteriaWithOneAC,
		"fixture_test.go":                       goTestPositive,
	})
	target, err := DeriveTargetManifest(root)
	if err != nil {
		t.Fatal(err)
	}
	source := SourceManifest{Entries: []SourceEntry{
		{ID: "AC-900", ProofClass: "Integration", Direction: "negative", EvidenceLocation: "fixture_test.go"},
	}}
	if r := Compare(source, target); !r.Failed() {
		t.Fatal("expected a failing report for a direction the target does not carry")
	}
}

// TestAC110_UnitPositive_excludedEntryIsNotMissing proves a declared
// exclusion with a reason is not reported as a missing criterion.
func TestAC110_UnitPositive_excludedEntryIsNotMissing(t *testing.T) {
	target := TargetManifest{Entries: map[string]TargetEntry{}}
	source := SourceManifest{Entries: []SourceEntry{
		{ID: "AC-901", Excluded: true, Reason: "not carried forward"},
	}}
	if r := Compare(source, target); r.Failed() {
		t.Fatalf("expected an excluded entry to produce no finding, got %v", r.Findings)
	}
}

// TestAC110_UnitNegative_missingCriterionFails proves a non-excluded source
// entry with no target counterpart is reported as missing.
func TestAC110_UnitNegative_missingCriterionFails(t *testing.T) {
	target := TargetManifest{Entries: map[string]TargetEntry{}}
	source := SourceManifest{Entries: []SourceEntry{
		{ID: "AC-901", ProofClass: "Unit", Direction: "positive", EvidenceLocation: "x_test.go"},
	}}
	r := Compare(source, target)
	if !r.Failed() || r.Findings[0].Class != ClassMissingCriterion {
		t.Fatalf("expected a missing-criterion finding, got %v", r.Findings)
	}
}

// TestAC111_UnitPositive_unmentionedTargetIsIgnored proves target evidence
// for an ID the source manifest never mentions and that carries no
// classified evidence is not reported as orphaned.
func TestAC111_UnitPositive_unmentionedTargetIsIgnored(t *testing.T) {
	target := TargetManifest{Entries: map[string]TargetEntry{
		"AC-902": {ID: "AC-902"},
	}}
	source := SourceManifest{}
	if r := Compare(source, target); r.Failed() {
		t.Fatalf("expected no finding for an evidence-less, unmentioned target entry, got %v", r.Findings)
	}
}

// TestAC111_UnitNegative_orphanedTagFails proves classified target evidence
// for an ID absent from the source manifest entirely fails as orphaned.
func TestAC111_UnitNegative_orphanedTagFails(t *testing.T) {
	target := TargetManifest{Entries: map[string]TargetEntry{
		"AC-902": {ID: "AC-902", EvidenceLocations: []string{"x_test.go"}},
	}}
	source := SourceManifest{}
	r := Compare(source, target)
	if !r.Failed() || r.Findings[0].Class != ClassOrphanedTag {
		t.Fatalf("expected an orphaned-tag finding, got %v", r.Findings)
	}
}

// TestAC112_UnitPositive_unchangedEvidenceIsNotReported proves an entry whose
// direction and location match the target produces no finding.
func TestAC112_UnitPositive_unchangedEvidenceIsNotReported(t *testing.T) {
	target := TargetManifest{Entries: map[string]TargetEntry{
		"AC-903": {ID: "AC-903", ProofClass: "Unit", Directions: []string{"positive"}, EvidenceLocations: []string{"x_test.go"}},
	}}
	source := SourceManifest{Entries: []SourceEntry{
		{ID: "AC-903", ProofClass: "Unit", Direction: "positive", EvidenceLocation: "x_test.go"},
	}}
	if r := Compare(source, target); r.Failed() {
		t.Fatalf("expected no finding for unchanged evidence, got %v", r.Findings)
	}
}

// TestAC112_UnitNegative_changedLocationFails proves a relocated evidence
// reference fails as changed evidence.
func TestAC112_UnitNegative_changedLocationFails(t *testing.T) {
	target := TargetManifest{Entries: map[string]TargetEntry{
		"AC-903": {ID: "AC-903", ProofClass: "Unit", Directions: []string{"positive"}, EvidenceLocations: []string{"moved_test.go"}},
	}}
	source := SourceManifest{Entries: []SourceEntry{
		{ID: "AC-903", ProofClass: "Unit", Direction: "positive", EvidenceLocation: "x_test.go"},
	}}
	r := Compare(source, target)
	if !r.Failed() || r.Findings[0].Class != ClassChangedEvidence {
		t.Fatalf("expected a changed-evidence finding, got %v", r.Findings)
	}
}

// TestAC113_UnitPositive_matchingFingerprintPasses proves a manifest whose
// source-revision matches the ledger-recorded revision produces no finding.
func TestAC113_UnitPositive_matchingFingerprintPasses(t *testing.T) {
	target := TargetManifest{Entries: map[string]TargetEntry{
		"AC-904": {ID: "AC-904", ProofClass: "Unit", Directions: []string{"positive"}, EvidenceLocations: []string{"x_test.go"}, SourceRevision: "rev-1"},
	}}
	source := SourceManifest{SourceRevision: "rev-1", Entries: []SourceEntry{
		{ID: "AC-904", ProofClass: "Unit", Direction: "positive", EvidenceLocation: "x_test.go"},
	}}
	if r := Compare(source, target); r.Failed() {
		t.Fatalf("expected no finding for a matching fingerprint, got %v", r.Findings)
	}
}

// TestAC113_UnitNegative_staleFingerprintFails proves a manifest revision
// disagreeing with the ledger-recorded revision fails as stale.
func TestAC113_UnitNegative_staleFingerprintFails(t *testing.T) {
	target := TargetManifest{Entries: map[string]TargetEntry{
		"AC-904": {ID: "AC-904", ProofClass: "Unit", Directions: []string{"positive"}, EvidenceLocations: []string{"x_test.go"}, SourceRevision: "rev-2"},
	}}
	source := SourceManifest{SourceRevision: "rev-1", Entries: []SourceEntry{
		{ID: "AC-904", ProofClass: "Unit", Direction: "positive", EvidenceLocation: "x_test.go"},
	}}
	r := Compare(source, target)
	if !r.Failed() || r.Findings[0].Class != ClassStaleFingerprint {
		t.Fatalf("expected a stale-fingerprint finding, got %v", r.Findings)
	}
}

// TestAC114_UnitPositive_justifiedDispositionPasses proves a source entry
// carrying a disposition and justification produces no finding for a
// target criterion classified draft, Human, or retired.
func TestAC114_UnitPositive_justifiedDispositionPasses(t *testing.T) {
	target := TargetManifest{Entries: map[string]TargetEntry{
		"AC-905": {ID: "AC-905", Draft: true},
	}}
	source := SourceManifest{Entries: []SourceEntry{
		{ID: "AC-905", Disposition: DispositionDraft, Justification: "plan door P-012 M-060; source openspec/specs/foo/spec.md#L20"},
	}}
	if r := Compare(source, target); r.Failed() {
		t.Fatalf("expected no finding for a justified disposition, got %v", r.Findings)
	}
}

// TestAC114_UnitNegative_unjustifiedDispositionFails proves a target
// criterion classified draft, Human, or retired with no matching source
// disposition and justification fails.
func TestAC114_UnitNegative_unjustifiedDispositionFails(t *testing.T) {
	target := TargetManifest{Entries: map[string]TargetEntry{
		"AC-905": {ID: "AC-905", Human: true},
	}}
	source := SourceManifest{Entries: []SourceEntry{
		{ID: "AC-905", ProofClass: "Human"},
	}}
	r := Compare(source, target)
	if !r.Failed() || r.Findings[0].Class != ClassUnjustifiedDisposition {
		t.Fatalf("expected an unjustified-disposition finding, got %v", r.Findings)
	}
}
