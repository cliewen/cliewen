package corpus

import "testing"

func coverageFor(t *testing.T, files map[string]string, cap string) string {
	t.Helper()
	c, issues := Scan(writeCorpus(t, files))
	for _, i := range issues {
		t.Fatalf("unexpected scan issue: %v", i)
	}
	for _, cc := range Coverage(c) {
		if cc.Capability == cap {
			return cc.State
		}
	}
	t.Fatalf("capability %s not found in coverage report", cap)
	return ""
}

func TestAC047_UnitPositive_FullySatisfiedCapabilityIsCovered(t *testing.T) {
	files := capFiles("active")
	files["pkg/x_test.go"] = "package x\n\nfunc TestAC101_Works(t *testing.T) {}\n"
	if state := coverageFor(t, files, "CAP-101"); state != "covered" {
		t.Fatalf("expected covered, got %s", state)
	}
}

func TestAC047_UnitNegative_PartialAndGapCapabilitiesReportSeparately(t *testing.T) {
	files := capFiles("active")
	files["docs/capabilities/CAP-101-x/criteria.md"] = "---\nid: CAP-101-criteria\ntype: criteria\nstatus: active\nlinks: [CAP-101]\ntitle: X criteria\n---\n\n```gherkin\nFeature: X\n\n  @AC-101\n  Scenario: proven\n    Given a thing\n    Then it works\n\n  @AC-102 @draft\n  Scenario: not proven yet\n    Given a thing\n    Then it will work\n```\n"
	files["pkg/x_test.go"] = "package x\n\nfunc TestAC101_Works(t *testing.T) {}\n"
	if state := coverageFor(t, files, "CAP-101"); state != "partial" {
		t.Fatalf("expected partial, got %s", state)
	}

	files["docs/capabilities/CAP-101-x/criteria.md"] = "---\nid: CAP-101-criteria\ntype: criteria\nstatus: active\nlinks: [CAP-101]\ntitle: X criteria\n---\n\n```gherkin\nFeature: X\n\n  @AC-101 @draft\n  Scenario: not proven yet\n    Given a thing\n    Then it will work\n```\n"
	delete(files, "pkg/x_test.go")
	if state := coverageFor(t, files, "CAP-101"); state != "gap" {
		t.Fatalf("expected gap, got %s", state)
	}
}
