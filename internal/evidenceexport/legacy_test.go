package evidenceexport

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSanity_CompatibilityExportKeepsExistingRepositoryProof(t *testing.T) {
	root, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	producer, err := Collect(root, "source-regression")
	if err != nil || len(producer.Diagnostics) != 0 || len(producer.References) < 100 {
		t.Fatalf("source export: refs=%d diagnostics=%v error=%v", len(producer.References), producer.Diagnostics, err)
	}
}

func TestSanity_CompatibilityJVMAttribution(t *testing.T) {
	prefixes := map[string]bool{"AC": true, "PDO": true, "SNAP-SQS": true}
	for _, tc := range []struct {
		name, source, file string
		refs, issues       int
	}{
		{"literal", "class ExampleTest {\n@Test\n@Tag(\"AC-001\")\n@Tag(\"unit\")\n@Tag(\"positive\")\nvoid works() {}\n}", "ExampleTest.java", 1, 0},
		{"named", "class ExampleTest {\nvoid testPDO_115_UnitPositive_works() {}\n}", "ExampleTest.java", 1, 0},
		{"kotlin", "class ExampleTest {\n@Test\n@Tag(\"SNAP_SQS_001\")\n@Tag(\"integration\")\n@Tag(\"negative\")\nfun works() {}\n}", "ExampleTest.kt", 1, 0},
		{"legacy", "@Test\n@Tag(\"AC-001\")\nvoid works() {}", "ExampleTest.java", 1, 0},
		{"parameterized", "@ParameterizedTest\n@Tag(\"AC-001\")\n@Tag(\"e2e\")\n@Tag(\"negative\")\nvoid works(String value) {}", "ExampleTest.java", 1, 0},
		{"matching", "@Test\n@Tag(\"AC-001\")\n@Tag(\"Unit\")\n@Tag(\"positive\")\nvoid testAC001_UnitPositive_works() {}", "ExampleTest.java", 1, 0},
		{"container", "@Tag(\"AC-001\")\nclass ExampleTest {}", "ExampleTest.java", 0, 1},
		{"dynamic", "@Test\n@Tag(dynamic)\nvoid works() {}", "ExampleTest.java", 0, 1},
		{"dangling", "@Tag(\"AC-001\")", "ExampleTest.java", 0, 1},
		{"detached", "@Tag(\"AC-001\")\nint value = 1;\nvoid works() {}", "ExampleTest.java", 0, 1},
		{"unannotated", "@Tag(\"AC-001\")\nvoid works() {}", "ExampleTest.java", 0, 1},
		{"disagree", "@Test\n@Tag(\"AC-002\")\nvoid testAC001_UnitPositive_works() {}", "ExampleTest.java", 2, 1},
		{"multiple", "@Test\n@Tag(\"AC-001\")\n@Tag(\"PDO-115\")\nvoid works() {}", "ExampleTest.java", 2, 1},
		{"malformed", "@Test\n@Tag(\"AC-bad\")\nvoid works() {}", "ExampleTest.java", 0, 1},
		{"ordinary", "@Test\n@Tag(\"java-17\")\nvoid ordinary() {}", "ExampleTest.java", 0, 0},
		{"strings", "// @Tag(\"AC-001\")\nString text = \"@Tag(\\\"AC-001\\\")\";\nString block = \"\"\"\n@Tag(\"AC-001\")\nvoid testAC001_UnitPositive_fake() {}\n\"\"\";", "ExampleTest.java", 0, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			refs, issues := harvestJVMEvidence(tc.source, tc.file, prefixes)
			if len(refs) != tc.refs || len(issues) != tc.issues {
				t.Fatalf("refs=%v issues=%v", refs, issues)
			}
		})
	}
}

func TestSanity_CompatibilityCucumberAndDiagnostics(t *testing.T) {
	for _, tc := range []struct {
		source       string
		refs, issues int
	}{
		{"@AC-001 @unit @positive\nScenario: proof\n Given an input", 1, 0},
		{"@AC_001 @integration @negative\nScenario Outline: proof\n Given an input", 1, 0},
		{"@AC-001\nScenario: legacy", 1, 0},
		{"@AC-001 @unit @e2e @positive\nScenario: ambiguous", 0, 1},
		{"@AC-bad\nScenario: malformed", 0, 1},
		{"@ordinary\nScenario: helper", 0, 0},
		{"@AC-001\nFeature: container", 0, 0},
	} {
		root := t.TempDir()
		if err := os.WriteFile(filepath.Join(root, "example.feature"), []byte(tc.source), 0644); err != nil {
			t.Fatal(err)
		}
		producer, err := Collect(root, "cucumber")
		if err != nil || len(producer.References) != tc.refs || len(producer.Diagnostics) != tc.issues {
			t.Fatalf("%s: %+v %v", tc.source, producer, err)
		}
		if err := Fixture(root); err != nil {
			t.Fatal(err)
		}
	}
}
