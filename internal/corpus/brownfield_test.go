package corpus

import "testing"

func extendedCriteria(prefix, ids string) string {
	return "---\nid: CAP-101-criteria\ntype: criteria\nstatus: active\nlinks: [CAP-101]\ntitle: X criteria\nac-prefix: " + prefix + "\n---\n\n```gherkin\nFeature: X\n\n" + ids + "```\n"
}

func TestAC059_UnitPositive_ExtendedIDsPassAndCoverage(t *testing.T) {
	files := capFiles("active")
	files["docs/capabilities/CAP-101-x/README.md"] = "---\nid: CAP-101\ntype: capability\nstatus: active\nlinks: [SNAP-SQS-002b]\ntitle: X\n---\n"
	files["docs/capabilities/CAP-101-x/criteria.md"] = extendedCriteria("SNAP-SQS", "  @SNAP-SQS-001\n  Scenario: first\n    Test-type: Unit\n    Given a thing\n    Then it works\n\n  @SNAP-SQS-002b\n  Scenario: second\n    Test-type: Unit\n    Given a thing\n    Then it works\n")
	files["pkg/x_test.go"] = "package x\n\nfunc TestSNAPSQS001_UnitPositive_Accepts(t *testing.T) {}\nfunc TestSNAPSQS001_UnitNegative_Rejects(t *testing.T) {}\nfunc TestSNAPSQS002b_UnitPositive_Accepts(t *testing.T) {}\nfunc TestSNAPSQS002b_UnitNegative_Rejects(t *testing.T) {}\n"
	if issues := run(t, files, false); len(issues) != 0 {
		t.Fatalf("extended IDs should validate, got %v", issues)
	}
	if state := coverageFor(t, files, "CAP-101"); state != "covered" {
		t.Fatalf("extended IDs should produce covered capability, got %s", state)
	}
	c, issues := Scan(writeCorpus(t, files))
	if len(issues) != 0 {
		t.Fatalf("extended corpus scan failed, got %v", issues)
	}
	if artifacts, unfollowed, err := Context(c, "SNAP-SQS-002b"); err != nil || len(unfollowed) != 0 || len(artifacts) != 2 || artifacts[0].ID != "CAP-101-criteria" || artifacts[1].ID != "CAP-101" {
		t.Fatalf("extended criterion context = artifacts %v, unfollowed %v, err %v", artifacts, unfollowed, err)
	}
}

func TestAC059_UnitNegative_ExtendedIDRulesRejectMalformedAndCollidingForms(t *testing.T) {
	tests := map[string]struct {
		files map[string]string
		want  string
	}{
		"wrong namespace": {
			files: func() map[string]string {
				files := capFiles("active")
				files["docs/capabilities/CAP-101-x/criteria.md"] = extendedCriteria("SNAP-SQS", "  @ADP-045b\n  Scenario: foreign\n    Given a thing\n    Then it fails\n")
				return files
			}(),
			want: "outside this file's namespace SNAP-SQS",
		},
		"malformed suffix": {
			files: func() map[string]string {
				files := capFiles("active")
				files["docs/capabilities/CAP-101-x/criteria.md"] = extendedCriteria("SNAP-SQS", "  @SNAP-SQS-001B\n  Scenario: malformed\n    Given a thing\n    Then it fails\n")
				return files
			}(),
			want: "not a canonical acceptance-criterion ID",
		},
		"lowercase identity": {
			files: func() map[string]string {
				files := capFiles("active")
				files["docs/capabilities/CAP-101-x/criteria.md"] = extendedCriteria("SNAP-SQS", "  @snap-sqs-001\n  Scenario: lowercase\n    Given a thing\n    Then it fails\n")
				return files
			}(),
			want: "not a canonical acceptance-criterion ID",
		},
		"underscore declaration": {
			files: func() map[string]string {
				files := capFiles("active")
				files["docs/capabilities/CAP-101-x/criteria.md"] = extendedCriteria("SNAP-SQS", "  @SNAP_SQS_001\n  Scenario: carrier alias is not a declaration\n    Given a thing\n    Then it fails\n")
				return files
			}(),
			want: "not a canonical acceptance-criterion ID",
		},
		"normalized prefix collision": {
			files: func() map[string]string {
				files := capFiles("active")
				files["docs/capabilities/CAP-101-x/criteria.md"] = extendedCriteria("SNAP-SQS", "  @SNAP-SQS-001\n  Scenario: first\n    Given a thing\n    Then it works\n")
				files["docs/capabilities/CAP-101-x/other-criteria.md"] = "---\nid: CAP-101-other-criteria\ntype: criteria\nstatus: active\nlinks: [CAP-101]\ntitle: Other criteria\nac-prefix: SNAPSQS\n---\n\n```gherkin\nFeature: Other\n\n  @SNAPSQS-002\n  Scenario: collision\n    Given a thing\n    Then it fails\n```\n"
				return files
			}(),
			want: "collides with",
		},
		"duplicate identity": {
			files: func() map[string]string {
				files := capFiles("active")
				files["docs/capabilities/CAP-101-x/criteria.md"] = extendedCriteria("SNAP-SQS", "  @SNAP-SQS-001\n  Scenario: first\n    Given a thing\n    Then it works\n\n  @SNAP-SQS-001\n  Scenario: duplicate\n    Given a thing\n    Then it fails\n")
				return files
			}(),
			want: "duplicate declaration of SNAP-SQS-001",
		},
		"retired suffix reference": {
			files: func() map[string]string {
				files := capFiles("active")
				files["docs/capabilities/CAP-101-x/criteria.md"] = extendedCriteria("ADP", "  @ADP-045b @retired\n  Scenario: retired\n    Given a thing\n    Then it used to work\n")
				files["jvm/RetiredTest.java"] = `import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

class RetiredTest {
  @Tag("ADP_045b")
  @Test
  void retired() {}
}
`
				return files
			}(),
			want: "retired ADP-045b",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assertIssue(t, run(t, test.files, false), test.want)
		})
	}
}

func TestAC060_UnitPositive_ExtendedIDsWorkAcrossEvidenceCarriers(t *testing.T) {
	files := capFiles("active")
	files["docs/capabilities/CAP-101-x/criteria.md"] = extendedCriteria("SNAP-SQS", "  @SNAP-SQS-001\n  Scenario: JVM and Go\n    Test-type: Unit\n    Given a thing\n    Then it works\n")
	files["docs/capabilities/CAP-101-x/other-criteria.md"] = "---\nid: CAP-101-other-criteria\ntype: criteria\nstatus: active\nlinks: [CAP-101]\ntitle: Other criteria\nac-prefix: ADP\n---\n\n```gherkin\nFeature: Other\n\n  @ADP-045b\n  Scenario: Cucumber\n    Test-type: Unit\n    Given a thing\n    Then it works\n```\n"
	files["pkg/x_test.go"] = "package x\n\nfunc TestSNAPSQS001_UnitPositive_Go(t *testing.T) {}\nfunc TestSNAPSQS001_UnitNegative_Go(t *testing.T) {}\nfunc TestADP045b_UnitPositive_Go(t *testing.T) {}\nfunc TestADP045b_UnitNegative_Go(t *testing.T) {}\n"
	files["jvm/ExtendedTest.java"] = `package x;

import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

class ExtendedTest {
  @Tag("SNAP_SQS_001")
  @Tag("Unit")
  @Tag("Positive")
  @Test
  void snapPositive() {}

  @Tag("SNAP_SQS_001")
  @Tag("Unit")
  @Tag("Negative")
  @Test
  void snapNegative() {}

  @Tag("ADP_045b")
  @Tag("Unit")
  @Tag("Positive")
  @Test
  void adpPositive() {}

  @Tag("ADP_045b")
  @Tag("Unit")
  @Tag("Negative")
  @Test
  void adpNegative() {}

  void testSNAPSQS001_UnitPositive_named() {}
  void testSNAPSQS001_UnitNegative_named() {}
  void testADP045b_UnitPositive_named() {}
  void testADP045b_UnitNegative_named() {}
}
`
	files["features/extended.feature"] = `Feature: extended IDs

  @SNAP-SQS-001 @unit @positive
  Scenario: Cucumber positive
    Given a thing
    Then it works

  @SNAP_SQS_001 @unit @positive
  Scenario: Cucumber normalized positive
    Given a thing
    Then it works

  @ADP-045b @unit @negative
  Scenario: Cucumber negative
    Given a thing
    Then it fails
`
	if issues := run(t, files, false); len(issues) != 0 {
		t.Fatalf("extended evidence carriers should validate, got %v", issues)
	}
}

func TestAC060_UnitNegative_ExtendedEvidenceRejectsMalformedCarrier(t *testing.T) {
	files := capFiles("active")
	files["docs/capabilities/CAP-101-x/criteria.md"] = extendedCriteria("SNAP-SQS", "  @SNAP-SQS-001\n  Scenario: malformed carrier\n    Test-type: Unit\n    Given a thing\n    Then it works\n")
	files["pkg/x_test.go"] = "package x\n\nfunc TestSNAPSQS001B_UnitPositive_Malformed(t *testing.T) {}\n"
	assertIssue(t, run(t, files, false), "declares no purpose")

	files["pkg/x_test.go"] = "package x\n\nfunc TestSNAPSQS001_UnitPositive_Valid(t *testing.T) {}\n"
	files["jvm/ExtendedTest.java"] = `class ExtendedTest {
  void testSNAPSQS001B_UnitNegative_Malformed() {}
}
`
	assertIssue(t, run(t, files, false), "malformed or ambiguous normalized criterion prefix")
}
