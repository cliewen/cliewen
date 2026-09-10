package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAC186_UnitPositive_NextCommandReportsRecommendationAlternativesAndDrafts(t *testing.T) {
	root := t.TempDir()
	writeNextFile(t, root, "docs/plans/P-001.md", "---\nid: P-001\ntype: plan\nstatus: active\nlinks: []\ntitle: First plan\n---\n\n| ID | Milestone | Exit criterion | Status | Evidence |\n|---|---|---|---|---|\n| M-001 | First | Do first | todo | |\n")
	writeNextFile(t, root, "docs/plans/P-002.md", "---\nid: P-002\ntype: plan\nstatus: active\nlinks: []\ntitle: Second plan\n---\n\n| ID | Milestone | Status | Evidence |\n|---|---|---|---|\n| M-002 | Resume | doing | |\n")
	writeNextFile(t, root, "docs/plans/P-003.md", "---\nid: P-003\ntype: plan\nstatus: draft\nlinks: []\ntitle: Proposed plan\n---\n\n| ID | Milestone | Status | Evidence |\n|---|---|---|---|\n| M-003 | Proposed | todo | |\n")

	var out, errOut bytes.Buffer
	if code := runNext([]string{root}, &out, &errOut); code != 0 {
		t.Fatalf("next exit = %d, stderr = %s", code, errOut.String())
	}
	got := out.String()
	if !strings.Contains(got, "P-002/M-002 | doing | Resume") || !strings.Contains(got, "Other active alternatives: 1") {
		t.Fatalf("recommendation or alternatives missing: %s", got)
	}
	if !strings.Contains(got, "P-003/M-003 | todo | Proposed") || !strings.Contains(got, "not actionable") {
		t.Fatalf("draft milestone was not explained: %s", got)
	}

	out.Reset()
	errOut.Reset()
	if code := runNext([]string{"--all", root}, &out, &errOut); code != 0 {
		t.Fatalf("next --all exit = %d, stderr = %s", code, errOut.String())
	}
	got = out.String()
	if strings.Index(got, "P-002/M-002") > strings.Index(got, "P-001/M-001") {
		t.Fatalf("doing milestone did not precede todo: %s", got)
	}
	activeSection := strings.SplitN(got, "Proposed unfinished milestones", 2)[0]
	if strings.Contains(activeSection, "P-003/M-003 |") {
		t.Fatalf("draft milestone leaked into active --all list: %s", got)
	}
}

func TestAC186_UnitNegative_NextCommandDoesNotMutateAndHandlesNoActiveWork(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "docs", "plans", "P-001.md")
	writeNextFile(t, root, "docs/plans/P-001.md", "---\nid: P-001\ntype: plan\nstatus: draft\nlinks: []\ntitle: Proposed\n---\n\n| ID | Milestone | Status |\n|---|---|---|\n| M-001 | Proposed | todo |\n")
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var out, errOut bytes.Buffer
	if code := runNext([]string{root}, &out, &errOut); code != 0 {
		t.Fatalf("next exit = %d, stderr = %s", code, errOut.String())
	}
	if !strings.Contains(out.String(), "No actionable milestone in active plans.") || !strings.Contains(out.String(), "P-001/M-001") {
		t.Fatalf("no-active report missing: %s", out.String())
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(before, after) {
		t.Fatalf("next mutated the plan")
	}
}

func writeNextFile(t *testing.T, root, rel, content string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
