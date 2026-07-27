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

// nsFiles extends capFiles with an ac-prefix: MG namespace declaring MG-101.
func nsFiles() map[string]string {
	files := capFiles("active")
	files["docs/capabilities/CAP-101-x/criteria.md"] = "---\nid: CAP-101-criteria\ntype: criteria\nstatus: active\nlinks: [CAP-101]\ntitle: X criteria\nac-prefix: MG\n---\n\n```gherkin\nFeature: X\n\n  @MG-101\n  Scenario: it works\n    Given a thing\n    Then it works\n```\n"
	return files
}

// AC-014: a criteria file declares ACs in its own namespace; the default
// namespace AC keeps working untouched (the rest of the suite is its proof).
func TestAC014_NamespacedACEnforced(t *testing.T) {
	files := nsFiles()
	assertIssue(t, run(t, files, false), "MG-101 has no test")

	files["pkg/x_test.go"] = "package x\n\nfunc TestMG101_Works(t *testing.T) {}\n"
	if issues := run(t, files, false); len(issues) != 0 {
		t.Fatalf("MG-101 has a Go test; expected no issues, got %v", issues)
	}
}

func TestAC014_InvalidPrefixReported(t *testing.T) {
	files := nsFiles()
	files["docs/capabilities/CAP-101-x/criteria.md"] = "---\nid: CAP-101-criteria\ntype: criteria\nstatus: active\nlinks: [CAP-101]\ntitle: X criteria\nac-prefix: mg1\n---\n"
	assertIssue(t, run(t, files, false), "ac-prefix must be uppercase")
}

// AC-015: a tag outside the file's declared namespace fails.
func TestAC015_ForeignNamespaceTagFails(t *testing.T) {
	files := nsFiles()
	files["docs/capabilities/CAP-101-x/criteria.md"] = "---\nid: CAP-101-criteria\ntype: criteria\nstatus: active\nlinks: [CAP-101]\ntitle: X criteria\nac-prefix: MG\n---\n\n```gherkin\nFeature: X\n\n  @PG-001\n  Scenario: wrong namespace\n    Then it fails\n```\n"
	issues := run(t, files, false)
	assertIssue(t, issues, "@PG-001")
	assertIssue(t, issues, "namespace MG")
}

// AC-016: a JVM @Tag in a declared namespace satisfies coverage; tags
// outside every namespace are runner metadata.
func TestAC016_JvmTagSatisfiesCoverage(t *testing.T) {
	files := nsFiles()
	files["core/src/test/kotlin/XTest.kt"] = "package x\n\nclass XTest {\n    @Test\n    @Tag(\"MG_101\")\n    @Tag(\"UNIT\")\n    fun `it works`() {}\n}\n"
	if issues := run(t, files, false); len(issues) != 0 {
		t.Fatalf("MG-101 has a tagged JVM test; expected no issues, got %v", issues)
	}

	files["core/src/test/java/YTests.java"] = "class YTests {\n    @Tag(\"INTEGRATION\")\n    void nothingToDeclare() {}\n}\n"
	if issues := run(t, files, false); len(issues) != 0 {
		t.Fatalf("runner-only tags are ignored; expected no issues, got %v", issues)
	}
}

// AC-017: a JVM tag referencing an unknown or retired AC fails.
func TestAC017_JvmUnknownOrRetiredTagFails(t *testing.T) {
	files := nsFiles()
	files["core/src/test/kotlin/XTest.kt"] = "class XTest {\n    @Tag(\"MG_101\")\n    @Tag(\"MG_999\")\n    fun f() {}\n}\n"
	issues := run(t, files, false)
	assertIssue(t, issues, "core/src/test/kotlin/XTest.kt")
	assertIssue(t, issues, "tag \"MG_999\" references MG-999 which no criteria.md declares")

	files["docs/capabilities/CAP-101-x/criteria.md"] = "---\nid: CAP-101-criteria\ntype: criteria\nstatus: active\nlinks: [CAP-101]\ntitle: X criteria\nac-prefix: MG\n---\n\n```gherkin\nFeature: X\n\n  @MG-101 @retired\n  Scenario: it used to work\n    Then it worked\n```\n"
	files["core/src/test/kotlin/XTest.kt"] = "class XTest {\n    @Tag(\"MG_101\")\n    fun f() {}\n}\n"
	assertIssue(t, run(t, files, false), "references retired MG-101")
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

func TestAC043_UnitPositive_ClassifiedEvidenceNeedsBothDirections(t *testing.T) {
	files := capFiles("active")
	files["docs/capabilities/CAP-101-x/criteria.md"] = "---\nid: CAP-101-criteria\ntype: criteria\nstatus: active\nlinks: [CAP-101]\ntitle: X criteria\n---\n\n```gherkin\nFeature: X\n\n  @AC-101\n  Scenario: it works\n    Test-type: Unit\n    Given a thing\n    Then it works\n```\n"
	files["pkg/x_test.go"] = "package x\n\nfunc TestAC101_UnitPositive_Works(t *testing.T) {}\n"
	assertIssue(t, run(t, files, false), "AC-101 has no Unit negative evidence")
}

func TestAC043_UnitNegative_ClassifiedEvidenceCoversThePair(t *testing.T) {
	files := capFiles("active")
	files["docs/capabilities/CAP-101-x/criteria.md"] = "---\nid: CAP-101-criteria\ntype: criteria\nstatus: active\nlinks: [CAP-101]\ntitle: X criteria\n---\n\n```gherkin\nFeature: X\n\n  @AC-101\n  Scenario: it works\n    Test-type: Unit\n    Given a thing\n    Then it works\n```\n"
	files["pkg/x_test.go"] = "package x\n\nfunc TestAC101_UnitPositive_Works(t *testing.T) {}\nfunc TestAC101_UnitNegative_Rejects(t *testing.T) {}\n"
	if issues := run(t, files, false); len(issues) != 0 {
		t.Fatalf("classified pair should cover the criterion, got %v", issues)
	}
}

func TestAC043_TestTypeAfterScenarioBodyLineFails(t *testing.T) {
	files := capFiles("active")
	files["docs/capabilities/CAP-101-x/criteria.md"] = "---\nid: CAP-101-criteria\ntype: criteria\nstatus: active\nlinks: [CAP-101]\ntitle: X criteria\n---\n\n```gherkin\nFeature: X\n\n  @AC-101\n  Scenario: it works\n    Given a thing\n    Test-type: Unit\n    Then it works\n```\n"
	files["pkg/x_test.go"] = "package x\n\nfunc TestAC101_UnitPositive_Works(t *testing.T) {}\nfunc TestAC101_UnitNegative_Rejects(t *testing.T) {}\n"
	assertIssue(t, run(t, files, false), "AC-101 has a Test-type that is not the first non-blank scenario-body line")
}

func TestAC043_TestTypeAfterContiguousScenarioTagsIsClassified(t *testing.T) {
	files := capFiles("active")
	files["docs/capabilities/CAP-101-x/criteria.md"] = "---\nid: CAP-101-criteria\ntype: criteria\nstatus: active\nlinks: [CAP-101]\ntitle: X criteria\n---\n\n```gherkin\nFeature: X\n\n  @AC-101\n  @slow\n  Scenario: it works\n    Test-type: Unit\n    Given a thing\n    Then it works\n```\n"
	files["pkg/x_test.go"] = "package x\n\nfunc TestAC101_UnitPositive_Works(t *testing.T) {}\n"
	assertIssue(t, run(t, files, false), "AC-101 has no Unit negative evidence")
}

func TestAC043_SeparatedTagsDoNotAttachToScenario(t *testing.T) {
	files := capFiles("active")
	files["docs/capabilities/CAP-101-x/criteria.md"] = "---\nid: CAP-101-criteria\ntype: criteria\nstatus: active\nlinks: [CAP-101]\ntitle: X criteria\n---\n\n```gherkin\nFeature: X\n\n  @AC-101\n\n  @slow\n  Scenario: it works\n    Test-type: Unit\n    Given a thing\n    Then it works\n```\n"
	files["pkg/x_test.go"] = "package x\n\nfunc TestAC101_UnitPositive_Works(t *testing.T) {}\n"
	if issues := run(t, files, false); len(issues) != 0 {
		t.Fatalf("a separated tag must not attach its Test-type to the AC, got %v", issues)
	}
}

func TestAC043_TestTypeAfterScenarioBodyTagFails(t *testing.T) {
	files := capFiles("active")
	files["docs/capabilities/CAP-101-x/criteria.md"] = "---\nid: CAP-101-criteria\ntype: criteria\nstatus: active\nlinks: [CAP-101]\ntitle: X criteria\n---\n\n```gherkin\nFeature: X\n\n  @AC-101\n  Scenario: it works\n    @slow\n    Test-type: Unit\n    Given a thing\n    Then it works\n```\n"
	files["pkg/x_test.go"] = "package x\n\nfunc TestAC101_UnitPositive_Works(t *testing.T) {}\nfunc TestAC101_UnitNegative_Rejects(t *testing.T) {}\n"
	assertIssue(t, run(t, files, false), "AC-101 has a Test-type that is not the first non-blank scenario-body line")
}

func TestAC043_UnitPositive_JvmTagsCarryTypeAndDirection(t *testing.T) {
	files := capFiles("active")
	files["docs/capabilities/CAP-101-x/criteria.md"] = "---\nid: CAP-101-criteria\ntype: criteria\nstatus: active\nlinks: [CAP-101]\ntitle: X criteria\n---\n\n```gherkin\nFeature: X\n\n  @AC-101\n  Scenario: it works\n    Test-type: Integration\n    Given a thing\n    Then it works\n```\n"
	files["core/src/test/kotlin/XTest.kt"] = "class XTest {\n    @Tag(\"AC_101\")\n    @Tag(\"integration\")\n    @Tag(\"positive\")\n    fun works() {}\n\n    @Tag(\"AC_101\")\n    @Tag(\"integration\")\n    @Tag(\"negative\")\n    fun rejects() {}\n}\n"
	if issues := run(t, files, false); len(issues) != 0 {
		t.Fatalf("JVM tags should satisfy classified coverage, got %v", issues)
	}
}

func TestAC043_UnitNegative_PerformanceNamesCarryTypeAndDirection(t *testing.T) {
	files := capFiles("active")
	files["docs/capabilities/CAP-101-x/criteria.md"] = "---\nid: CAP-101-criteria\ntype: criteria\nstatus: active\nlinks: [CAP-101]\ntitle: X criteria\n---\n\n```gherkin\nFeature: X\n\n  @AC-101\n  Scenario: it performs\n    Test-type: Performance\n    Given a thing\n    Then it responds within its bound\n```\n"
	files["pkg/x_test.go"] = "package x\n\nfunc TestAC101_PerformancePositive_MeetsBound(t *testing.T) {}\nfunc TestAC101_PerformanceNegative_ExceedsBound(t *testing.T) {}\n"
	if issues := run(t, files, false); len(issues) != 0 {
		t.Fatalf("Performance names should satisfy classified coverage, got %v", issues)
	}
}

func TestAC044_UnitPositive_CucumberPositiveNeedsNegative(t *testing.T) {
	files := capFiles("active")
	files["docs/capabilities/CAP-101-x/criteria.md"] = "---\nid: CAP-101-criteria\ntype: criteria\nstatus: active\nlinks: [CAP-101]\ntitle: X criteria\n---\n\n```gherkin\nFeature: X\n\n  @AC-101\n  Scenario: it works\n    Test-type: E2E\n    Given a thing\n    Then it works\n```\n"
	files["features/x.feature"] = "@AC-101 @e2e @positive\nScenario: it works\n"
	assertIssue(t, run(t, files, false), "AC-101 has no E2E negative evidence")
}

func TestAC044_UnitNegative_CucumberPairSatisfiesCoverage(t *testing.T) {
	files := capFiles("active")
	files["docs/capabilities/CAP-101-x/criteria.md"] = "---\nid: CAP-101-criteria\ntype: criteria\nstatus: active\nlinks: [CAP-101]\ntitle: X criteria\n---\n\n```gherkin\nFeature: X\n\n  @AC-101\n  Scenario: it works\n    Test-type: E2E\n    Given a thing\n    Then it works\n```\n"
	files["features/x.feature"] = "@AC-101 @e2e @positive\nScenario: it works\n\n@AC-101 @e2e @negative\nScenario: it fails\n"
	if issues := run(t, files, false); len(issues) != 0 {
		t.Fatalf("Cucumber pair should cover the criterion, got %v", issues)
	}
}

func TestAC044_CucumberMultiLineTagBlocksCarryClassifiedEvidence(t *testing.T) {
	files := capFiles("active")
	files["docs/capabilities/CAP-101-x/criteria.md"] = "---\nid: CAP-101-criteria\ntype: criteria\nstatus: active\nlinks: [CAP-101]\ntitle: X criteria\n---\n\n```gherkin\nFeature: X\n\n  @AC-101\n  Scenario: it works\n    Test-type: E2E\n    Given a thing\n    Then it works\n```\n"
	files["features/x.feature"] = "@AC-101\n@e2e @positive\nScenario: it works\n\n@AC-101\n@e2e @negative\nScenario: it fails\n"
	if issues := run(t, files, false); len(issues) != 0 {
		t.Fatalf("multi-line Cucumber tag blocks should satisfy classified coverage, got %v", issues)
	}
}

func TestAC044_CucumberNonScenarioTagBlocksDoNotCountAsEvidence(t *testing.T) {
	for name, feature := range map[string]string{
		"feature":  "@AC-101 @e2e @positive\nFeature: X\nScenario: it works\n",
		"rule":     "Feature: X\n@AC-101 @e2e @positive\nRule: a rule\nScenario: it works\n",
		"examples": "Feature: X\nScenario Outline: it works\n@AC-101 @e2e @positive\nExamples: data\n  | value |\n  | x |\n",
	} {
		t.Run(name, func(t *testing.T) {
			files := capFiles("active")
			files["docs/capabilities/CAP-101-x/criteria.md"] = "---\nid: CAP-101-criteria\ntype: criteria\nstatus: active\nlinks: [CAP-101]\ntitle: X criteria\n---\n\n```gherkin\nFeature: X\n\n  @AC-101\n  Scenario: it works\n    Test-type: E2E\n    Given a thing\n    Then it works\n```\n"
			files["features/x.feature"] = feature
			assertIssue(t, run(t, files, false), "AC-101 has no E2E positive evidence")
		})
	}
}
