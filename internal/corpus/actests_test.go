package corpus

import (
	"fmt"
	"strings"
	"testing"
)

// capFiles extends validFiles with a capability whose criteria declare AC-101.
func capFiles(criteriaStatus string) map[string]string {
	return with(validFiles, map[string]string{
		"docs/README.md":                          "# Corpus\n\n<!-- clue:index:start -->\n- [goals/](goals/README.md)\n- [plans/](plans/README.md)\n- [capabilities/](capabilities/README.md)\n<!-- clue:index:end -->\n",
		"docs/capabilities/README.md":             "# Capabilities\n\n<!-- clue:index:start -->\n- [CAP-101](CAP-101-x/README.md)\n<!-- clue:index:end -->\n",
		"docs/capabilities/CAP-101-x/README.md":   "---\nid: CAP-101\ntype: capability\nstatus: active\nlinks: []\ntitle: X\n---\n",
		"docs/capabilities/CAP-101-x/criteria.md": "---\nid: CAP-101-criteria\ntype: criteria\nstatus: " + criteriaStatus + "\nlinks: [CAP-101]\ntitle: X criteria\n---\n\n```gherkin\nFeature: X\n\n  @AC-101\n  Scenario: it works\n    Given a thing\n    Then it works\n```\n",
	})
}

// AC-009 negative: an AC in active criteria without a test fails.
func TestAC009_UnitNegative_ActiveACWithoutTestFails(t *testing.T) {
	files := capFiles("active")
	assertIssue(t, run(t, files, false), "AC-101 has no test")
}

// AC-009 positive: a supported reference satisfies an active AC, while a
// whole draft criteria file remains exempt from the active-file contract.
func TestAC009_UnitPositive_ReferencedACAndDraftFilePass(t *testing.T) {
	files := capFiles("active")
	files["pkg/x_test.go"] = "package x\n\nfunc TestAC101_Works(t *testing.T) {}\n"
	if issues := run(t, files, false); len(issues) != 0 {
		t.Fatalf("AC-101 has a test; expected no issues, got %v", issues)
	}
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

func TestAC043_SecondTestTypeInScenarioFails(t *testing.T) {
	files := capFiles("active")
	files["docs/capabilities/CAP-101-x/criteria.md"] = "---\nid: CAP-101-criteria\ntype: criteria\nstatus: active\nlinks: [CAP-101]\ntitle: X criteria\n---\n\n```gherkin\nFeature: X\n\n  @AC-101\n  Scenario: it works\n    Test-type: Unit\n    Given a thing\n    Test-type: E2E\n    Then it works\n```\n"
	files["pkg/x_test.go"] = "package x\n\nfunc TestAC101_UnitPositive_Works(t *testing.T) {}\nfunc TestAC101_UnitNegative_Rejects(t *testing.T) {}\n"
	assertIssue(t, run(t, files, false), "AC-101 has a Test-type that is not the first non-blank scenario-body line")
}

func TestAC043_UnitPositive_JvmTagsCarryTypeAndDirection(t *testing.T) {
	files := capFiles("active")
	files["docs/capabilities/CAP-101-x/criteria.md"] = "---\nid: CAP-101-criteria\ntype: criteria\nstatus: active\nlinks: [CAP-101]\ntitle: X criteria\n---\n\n```gherkin\nFeature: X\n\n  @AC-101\n  Scenario: it works\n    Test-type: Integration\n    Given a thing\n    Then it works\n```\n"
	files["core/src/test/kotlin/XTest.kt"] = "class XTest {\n    @Test\n    @Tag(\"AC_101\")\n    @Tag(\"integration\")\n    @Tag(\"positive\")\n    fun works() {}\n\n    @Test\n    @Tag(\"AC_101\")\n    @Tag(\"integration\")\n    @Tag(\"negative\")\n    fun rejects() {}\n}\n"
	if issues := run(t, files, false); len(issues) != 0 {
		t.Fatalf("JVM tags should satisfy classified coverage, got %v", issues)
	}
}

func TestAC057_UnitPositive_MixedJVMMethodsRemainIsolated(t *testing.T) {
	files := classifiedJvmFiles(101, 110)
	var java strings.Builder
	java.WriteString("class MixedTest {\n")
	for id := 101; id <= 110; id++ {
		fmt.Fprintf(&java, "  @Test @Tag(\"AC_%d\") @Tag(\"unit\") @Tag(\"positive\") void accepts%d() {}\n", id, id)
		fmt.Fprintf(&java, "  @Test @Tag(\"AC_%d\") @Tag(\"unit\") @Tag(\"negative\") void rejects%d() {}\n", id, id)
	}
	java.WriteString("}\n")
	files["core/src/test/java/MixedTest.java"] = java.String()

	if issues := run(t, files, false); len(issues) != 0 {
		t.Fatalf("ten mixed AC pairs should remain attributable to their own methods, got %v", issues)
	}
}

func TestAC057_UnitNegative_JVMTagsDoNotCrossProductOrCreditAmbiguity(t *testing.T) {
	files := classifiedJvmFiles(101, 102)
	files["core/src/test/kotlin/MixedTest.kt"] = "class MixedTest {\n  @Test @Tag(\"AC_101\") @Tag(\"unit\") @Tag(\"positive\") fun accepts101() {}\n  @Test @Tag(\"AC_102\") @Tag(\"unit\") @Tag(\"negative\") fun rejects102() {}\n}\n"
	issues := run(t, files, false)
	assertIssue(t, issues, "AC-101 has no Unit negative evidence")
	assertIssue(t, issues, "AC-102 has no Unit positive evidence")

	files["core/src/test/kotlin/MixedTest.kt"] = "class MixedTest {\n  @Test @Tag(\"AC_101\") @Tag(\"AC_102\") @Tag(\"unit\") @Tag(\"positive\") fun ambiguous() {}\n}\n"
	issues = run(t, files, false)
	assertIssue(t, issues, "ambiguous metadata receives no classified credit")
	assertIssue(t, issues, "AC-101 has no Unit positive evidence")
}

func TestAC058_UnitPositive_SupportedJVMExecutableFormsReceiveCredit(t *testing.T) {
	for name, positiveAnnotation := range map[string]string{
		"ordinary":      "@Test",
		"parameterized": "@ParameterizedTest",
		"repeated":      "@RepeatedTest",
		"factory":       "@TestFactory",
		"template":      "@TestTemplate",
	} {
		t.Run(name, func(t *testing.T) {
			files := classifiedJvmFiles(101, 101)
			files["core/src/test/java/XTest.java"] = "class XTest {\n  " + positiveAnnotation + " @Tag(\"AC_101\") @Tag(\"unit\") @Tag(\"positive\") void accepts() {}\n  @Test @Tag(\"AC_101\") @Tag(\"unit\") @Tag(\"negative\") void rejects() {}\n}\n"
			if issues := run(t, files, false); len(issues) != 0 {
				t.Fatalf("%s executable should receive one evidence credit, got %v", name, issues)
			}
		})
	}

	t.Run("nested Kotlin", func(t *testing.T) {
		files := classifiedJvmFiles(101, 101)
		files["core/src/test/kotlin/XTest.kt"] = "class XTest {\n  @Nested\n  class Inner {\n    @ParameterizedTest @Tag(\"AC_101\") @Tag(\"unit\") @Tag(\"positive\") fun accepts() {}\n    @RepeatedTest @Tag(\"AC_101\") @Tag(\"unit\") @Tag(\"negative\") fun rejects() {}\n  }\n}\n"
		if issues := run(t, files, false); len(issues) != 0 {
			t.Fatalf("nested and parameterized Kotlin methods should receive evidence, got %v", issues)
		}
	})

	t.Run("named fallback", func(t *testing.T) {
		files := classifiedJvmFiles(101, 101)
		files["core/src/test/java/XTest.java"] = "class XTest {\n  @Test @Tag(\"slow\") void testAC101_UnitPositive_accepts() {}\n  @Test @Tag(\"slow\") void testAC101_UnitNegative_rejects() {}\n}\n"
		if issues := run(t, files, false); len(issues) != 0 {
			t.Fatalf("stable named JVM executables should receive evidence while unrelated runner tags remain metadata, got %v", issues)
		}
	})
}

func TestAC058_UnitNegative_UnsupportedJVMEvidenceIsDiagnosedAndIgnored(t *testing.T) {
	t.Run("unknown and retired references", func(t *testing.T) {
		files := classifiedJvmFiles(101, 101)
		files["core/src/test/kotlin/XTest.kt"] = "class XTest {\n  @Test @Tag(\"AC_999\") @Tag(\"unit\") @Tag(\"positive\") fun ghost() {}\n  @Test @Tag(\"AC_101\") @Tag(\"unit\") @Tag(\"positive\") fun accepts() {}\n  @Test @Tag(\"AC_101\") @Tag(\"unit\") @Tag(\"negative\") fun rejects() {}\n}\n"
		assertIssue(t, run(t, files, false), "references AC-999 which no criteria.md declares")

		files["docs/capabilities/CAP-101-x/criteria.md"] = strings.Replace(files["docs/capabilities/CAP-101-x/criteria.md"], "@AC-101", "@AC-101 @retired", 1)
		assertIssue(t, run(t, files, false), "references retired AC-101")
	})

	t.Run("class-level AC", func(t *testing.T) {
		files := classifiedJvmFiles(101, 101)
		files["core/src/test/java/XTest.java"] = "@Tag(\"AC_101\") @Tag(\"unit\")\nclass XTest {\n  @Test @Tag(\"positive\") void accepts() {}\n  @Test @Tag(\"negative\") void rejects() {}\n}\n"
		issues := run(t, files, false)
		assertIssue(t, issues, "class-level JVM tag \"AC_101\" cannot be attributed")
		assertIssue(t, issues, "AC-101 has no Unit positive evidence")
	})

	t.Run("unsupported multiline tag", func(t *testing.T) {
		files := classifiedJvmFiles(101, 101)
		files["core/src/test/java/XTest.java"] = "class XTest {\n  @Test\n  @Tag(\n    \"AC_101\")\n  void accepts() {}\n}\n"
		assertIssue(t, run(t, files, false), "unsupported JVM @Tag syntax")
	})

	t.Run("executable annotation cannot drift", func(t *testing.T) {
		files := classifiedJvmFiles(101, 101)
		files["core/src/test/java/XTest.java"] = "class XTest {\n  @Test\n  void multiline(\n  ) {}\n\n  @Tag(\"AC_101\") @Tag(\"unit\") @Tag(\"positive\")\n  void unannotated() {}\n}\n"
		issues := run(t, files, false)
		assertIssue(t, issues, "carries JVM AC tags without a supported executable annotation")
		assertIssue(t, issues, "AC-101 has no Unit positive evidence")
	})

	t.Run("disagreeing name and tags", func(t *testing.T) {
		files := classifiedJvmFiles(101, 101)
		files["core/src/test/java/XTest.java"] = "class XTest {\n  @Test @Tag(\"AC_101\") @Tag(\"unit\") @Tag(\"positive\") void testAC101_UnitNegative_conflict() {}\n}\n"
		assertIssue(t, run(t, files, false), "disagreeing or incomplete named and JUnit evidence carriers")
	})

	t.Run("proximity comments and strings", func(t *testing.T) {
		files := classifiedJvmFiles(101, 101)
		files["core/src/test/kotlin/XTest.kt"] = "class XTest {\n  // AC: AC-101 Unit positive\n  // @Test @Tag(\"AC_101\") @Tag(\"unit\") @Tag(\"positive\") fun commented() {}\n  val example = \"@Tag(\\\"AC_101\\\")\"\n  val block = \"\"\"\n    @Test @Tag(\"AC_101\") @Tag(\"unit\") @Tag(\"negative\") fun text() {}\n  \"\"\"\n}\n"
		issues := run(t, files, false)
		assertIssue(t, issues, "AC-101 has no test")
		if len(issues) != 3 {
			t.Fatalf("comments and strings should leave only the missing test and two missing directions, got %v", issues)
		}
		for _, issue := range issues {
			if strings.Contains(issue.Msg, "commented") || strings.Contains(issue.Msg, "text") {
				t.Fatalf("comments and strings must not become JVM evidence: %v", issues)
			}
		}
	})
}

func classifiedJvmFiles(first, last int) map[string]string {
	files := capFiles("active")
	var criteria strings.Builder
	criteria.WriteString("---\nid: CAP-101-criteria\ntype: criteria\nstatus: active\nlinks: [CAP-101]\ntitle: X criteria\n---\n\n```gherkin\nFeature: X\n")
	for id := first; id <= last; id++ {
		fmt.Fprintf(&criteria, "\n  @AC-%d\n  Scenario: criterion %d\n    Test-type: Unit\n    Given input %d\n    Then result %d\n", id, id, id, id)
	}
	criteria.WriteString("```\n")
	files["docs/capabilities/CAP-101-x/criteria.md"] = criteria.String()
	return files
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

func TestAC044_CucumberAmbiguousDirectionTagsFail(t *testing.T) {
	files := capFiles("active")
	files["docs/capabilities/CAP-101-x/criteria.md"] = "---\nid: CAP-101-criteria\ntype: criteria\nstatus: active\nlinks: [CAP-101]\ntitle: X criteria\n---\n\n```gherkin\nFeature: X\n\n  @AC-101\n  Scenario: it works\n    Test-type: E2E\n    Given a thing\n    Then it works\n```\n"
	files["features/x.feature"] = "@AC-101 @e2e @positive @negative\nScenario: it works\n"
	assertIssue(t, run(t, files, false), "Cucumber tag block for AC-101 must declare at most one test type and direction")
	assertIssue(t, run(t, files, false), "AC-101 has no E2E positive evidence")
}

func TestAC045_UnitPositive_HumanClassNeedsNoCodeTest(t *testing.T) {
	files := capFiles("active")
	files["docs/capabilities/CAP-101-x/criteria.md"] = "---\nid: CAP-101-criteria\ntype: criteria\nstatus: active\nlinks: [CAP-101]\ntitle: X criteria\n---\n\n```gherkin\nFeature: X\n\n  @AC-101\n  Scenario: a human confirms it\n    Test-type: Human\n    Given a thing only a person can judge\n    Then a human confirms it in the acceptance brief\n```\n"
	if issues := run(t, files, false); len(issues) != 0 {
		t.Fatalf("Human class needs no code test, got %v", issues)
	}
}

func TestAC045_UnitNegative_HumanWithSingleDirectionIsMalformed(t *testing.T) {
	files := capFiles("active")
	files["docs/capabilities/CAP-101-x/criteria.md"] = "---\nid: CAP-101-criteria\ntype: criteria\nstatus: active\nlinks: [CAP-101]\ntitle: X criteria\n---\n\n```gherkin\nFeature: X\n\n  @AC-101\n  Scenario: a human confirms it\n    Test-type: Human (single-direction)\n    Given a thing only a person can judge\n    Then a human confirms it in the acceptance brief\n```\n"
	issues := run(t, files, false)
	assertIssue(t, issues, "AC-101 declares Test-type: Human (single-direction), which the Human class does not use")
	if len(issues) != 1 {
		t.Fatalf("expected exactly one issue naming the malformed declaration, not a second misleading \"has no test\" issue; got %v", issues)
	}
}

func TestAC046_UnitPositive_DraftTagExemptsOneCriterion(t *testing.T) {
	files := capFiles("active")
	files["docs/capabilities/CAP-101-x/criteria.md"] = "---\nid: CAP-101-criteria\ntype: criteria\nstatus: active\nlinks: [CAP-101]\ntitle: X criteria\n---\n\n```gherkin\nFeature: X\n\n  @AC-101 @draft\n  Scenario: not proven yet\n    Given a thing\n    Then it will work\n```\n"
	if issues := run(t, files, false); len(issues) != 0 {
		t.Fatalf("@draft criterion needs no test, got %v", issues)
	}
}

func TestAC046_UnitNegative_RemovingDraftRestoresTheTestRequirement(t *testing.T) {
	files := capFiles("active")
	files["docs/capabilities/CAP-101-x/criteria.md"] = "---\nid: CAP-101-criteria\ntype: criteria\nstatus: active\nlinks: [CAP-101]\ntitle: X criteria\n---\n\n```gherkin\nFeature: X\n\n  @AC-101\n  Scenario: not proven yet\n    Given a thing\n    Then it will work\n```\n"
	assertIssue(t, run(t, files, false), "AC-101 has no test")
}
