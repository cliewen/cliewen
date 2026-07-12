package corpus

import "testing"

// capFiles extends validFiles with a capability whose criteria declare AC-101.
func capFiles(criteriaStatus string) map[string]string {
	return with(validFiles, map[string]string{
		"docs/README.md":                          "# Corpus\n\n<!-- clue:index:start -->\n- [goals/](goals/README.md)\n- [plans/](plans/README.md)\n- [capabilities/](capabilities/README.md)\n<!-- clue:index:end -->\n",
		"docs/capabilities/README.md":             "# Capabilities\n\n<!-- clue:index:start -->\n- [CAP-101](CAP-101-x/README.md)\n<!-- clue:index:end -->\n",
		"docs/capabilities/CAP-101-x/README.md":   "---\nid: CAP-101\ntype: capability\nstatus: active\nlinks: []\ntitle: X\n---\n",
		"docs/capabilities/CAP-101-x/criteria.md": "---\nid: CAP-101-criteria\ntype: criteria\nstatus: " + criteriaStatus + "\nlinks: [CAP-101]\ntitle: X criteria\n---\n\n```gherkin\nFeature: X\n\n  @AC-101\n  Scenario: it works\n    Given a thing\n    Then it works\n```\n",
	})
}

// AC-009: an AC in active criteria without a test fails; with a test it passes.
func TestAC009_ActiveACWithoutTestFails(t *testing.T) {
	files := capFiles("active")
	assertIssue(t, run(t, files, false), "AC-101 has no test")

	files["pkg/x_test.go"] = "package x\n\nfunc TestAC101_Works(t *testing.T) {}\n"
	if issues := run(t, files, false); len(issues) != 0 {
		t.Fatalf("AC-101 has a test; expected no issues, got %v", issues)
	}
}

// AC-009 (negative side): draft criteria are exempt from the contract.
func TestAC009_DraftCriteriaExempt(t *testing.T) {
	if issues := run(t, capFiles("draft"), false); len(issues) != 0 {
		t.Fatalf("draft criteria are exempt; expected no issues, got %v", issues)
	}
}

// AC-010: a test referencing an AC no criteria.md declares fails.
func TestAC010_UnknownACReferenceFails(t *testing.T) {
	files := with(validFiles, map[string]string{
		"pkg/x_test.go": "package x\n\nfunc TestAC999_Ghost(t *testing.T) {}\n",
	})
	issues := run(t, files, false)
	assertIssue(t, issues, "AC-999")
	assertIssue(t, issues, "pkg/x_test.go")
}

// AC-012: a retired AC needs no test; a surviving test referencing it fails.
func TestAC012_RetiredACTombstone(t *testing.T) {
	files := capFiles("active")
	files["docs/capabilities/CAP-101-x/criteria.md"] = "---\nid: CAP-101-criteria\ntype: criteria\nstatus: active\nlinks: [CAP-101]\ntitle: X criteria\n---\n\n```gherkin\nFeature: X\n\n  @AC-101 @retired\n  Scenario: it used to work this way\n    Given a thing\n    Then it worked\n```\n"
	if issues := run(t, files, false); len(issues) != 0 {
		t.Fatalf("retired AC needs no test; expected no issues, got %v", issues)
	}
	files["pkg/x_test.go"] = "package x\n\nfunc TestAC101_Works(t *testing.T) {}\n"
	issues := run(t, files, false)
	assertIssue(t, issues, "retired AC-101")
}

// AC-013: the same AC declared twice fails, naming both files.
func TestAC013_DuplicateACDeclarationFails(t *testing.T) {
	files := capFiles("active")
	files["pkg/x_test.go"] = "package x\n\nfunc TestAC101_Works(t *testing.T) {}\n"
	files["docs/capabilities/CAP-101-x/criteria.md"] = "---\nid: CAP-101-criteria\ntype: criteria\nstatus: active\nlinks: [CAP-101]\ntitle: X criteria\n---\n\n```gherkin\nFeature: X\n\n  @AC-101\n  Scenario: one\n    Then fine\n\n  @AC-101\n  Scenario: two\n    Then clash\n```\n"
	issues := run(t, files, false)
	assertIssue(t, issues, "duplicate declaration of AC-101")
}

// AC-011: every test declares its purpose; Unit/Sanity/Arch need no AC.
func TestAC011_UnclassifiedTestFails(t *testing.T) {
	files := with(validFiles, map[string]string{
		"pkg/x_test.go": "package x\n\nfunc TestSomethingUseful(t *testing.T) {}\n",
	})
	issues := run(t, files, false)
	assertIssue(t, issues, "TestSomethingUseful")
	assertIssue(t, issues, "declares no purpose")

	files["pkg/x_test.go"] = "package x\n\nfunc TestUnit_Something(t *testing.T) {}\n\nfunc TestSanity_Env(t *testing.T) {}\n\nfunc TestArch_Layering(t *testing.T) {}\n\nfunc TestMain(m *testing.M) {}\n"
	if issues := run(t, files, false); len(issues) != 0 {
		t.Fatalf("purpose-classified tests need no AC; expected no issues, got %v", issues)
	}
}
