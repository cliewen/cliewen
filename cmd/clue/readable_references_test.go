package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAC233_UnitPositive_ReportsNameReferencesAndKeepRawSource(t *testing.T) {
	root := intentCorpus(t, "---\nid: VIS-001\ntype: vision\nstatus: active\nlinks: []\ntitle: A product\n---\n\nA direction.\n")
	path := filepath.Join(root, "docs/goals/G-001-first.md")
	before, e := os.ReadFile(path)
	if e != nil {
		t.Fatal(e)
	}
	var out bytes.Buffer
	if code := runValidate([]string{"--coverage", "--intent", root}, &out); code != 0 {
		t.Fatalf("exit %d: %s", code, out.String())
	}
	for _, want := range []string{"A capability (CAP-001): gap", "Another capability (CAP-002): gap", "vision: A product (VIS-001)", "use case: A journey (UC-001)", "crosses A capability (CAP-001)", "goal: First goal (G-001)", "capabilities: A capability (CAP-001) (active), Another capability (CAP-002) (active)"} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("missing %q: %s", want, out.String())
		}
	}
	after, e := os.ReadFile(path)
	if e != nil || !bytes.Equal(before, after) {
		t.Fatal("human report changed stored source")
	}
	var context, errOut bytes.Buffer
	if code := runContext([]string{"G-001", root}, &context, &errOut); code != 0 {
		t.Fatalf("context %d: %s", code, errOut.String())
	}
	if !strings.Contains(context.String(), "===== First goal (G-001) | docs/goals/G-001-first.md =====") || !strings.Contains(context.String(), string(before)) {
		t.Fatal(context.String())
	}
}

func TestAC233_UnitNegative_NamePresentationPreservesDraftSelectionAndUsageErrors(t *testing.T) {
	root := t.TempDir()
	writeNextFile(t, root, "docs/plans/P-901.md", "---\nid: P-901\ntype: plan\nstatus: draft\nlinks: []\ntitle: Proposed clarity work\n---\n\n| ID | Milestone | Status |\n|---|---|---|\n| M-901 | Clarify the failed outcome | todo |\n")
	var out, errOut bytes.Buffer
	if code := runNext([]string{root}, &out, &errOut); code != 0 {
		t.Fatal(code)
	}
	if !strings.Contains(out.String(), "not actionable") || !strings.Contains(out.String(), "Clarify the failed outcome (M-901) | todo | plan: Proposed clarity work (P-901)") {
		t.Fatal(out.String())
	}
	out.Reset()
	errOut.Reset()
	if code := runContext([]string{"G-999", root}, &out, &errOut); code != 1 || out.Len() != 0 || !strings.Contains(errOut.String(), "G-999") {
		t.Fatalf("unknown context changed: %d %s %s", code, out.String(), errOut.String())
	}
	out.Reset()
	errOut.Reset()
	if code := runContext([]string{}, &out, &errOut); code != 2 {
		t.Fatalf("usage exit changed: %d", code)
	}
}

func TestAC234_UnitPositive_ValidationDiagnosticIncludesTheCriterionMeaning(t *testing.T) {
	root := validCorpus(t)
	writeFile(t, root, "docs/README.md", "---\ntype: index\ntitle: Corpus\n---\n\n<!-- clue:index:start -->\n- [goals/](goals/README.md)\n- [capabilities/](capabilities/README.md)\n<!-- clue:index:end -->\n")
	writeFile(t, root, "docs/capabilities/README.md", "---\ntype: index\ntitle: Capabilities\n---\n\n<!-- clue:index:start -->\n- [CAP-101](CAP-101-x/README.md)\n<!-- clue:index:end -->\n")
	writeFile(t, root, "docs/capabilities/CAP-101-x/README.md", "---\nid: CAP-101\ntype: capability\nstatus: active\nlinks: [G-001]\ngoal: G-001\ntitle: Backup status\n---\n")
	writeFile(t, root, "docs/capabilities/CAP-101-x/criteria.md", "---\nid: CAP-101-criteria\ntype: criteria\nstatus: active\nlinks: [CAP-101]\ntitle: Backup proof\n---\n\n```gherkin\n@AC-101\nScenario: Show the failed backup clearly\n  Test-type: Unit\n  Given a failed backup\n  Then its state is shown\n```\n")
	var out bytes.Buffer
	if code := runValidate([]string{root}, &out); code != 1 {
		t.Fatalf("missing machine proof was accepted: %d %s", code, out.String())
	}
	if !strings.Contains(out.String(), "Show the failed backup clearly (AC-101) has no test") {
		t.Fatal(out.String())
	}
}
