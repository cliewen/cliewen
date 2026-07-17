package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// AC-002: the command-level path — init then validate exits 0 with no
// manual edits in between.
func TestAC002_InitThenValidateExitsZero(t *testing.T) {
	root := t.TempDir()
	var out bytes.Buffer
	if code := runInit([]string{root}, &out, &out); code != 0 {
		t.Fatalf("init: expected exit 0, got %d\n%s", code, out.String())
	}
	if code := runValidate([]string{root}); code != 0 {
		t.Fatalf("validate after init: expected exit 0, got %d", code)
	}
}

// AC-003: validate exits non-zero on a scaffolded corpus with a broken
// link (the message content is asserted in the scaffold package tests).
func TestAC003_ValidateExitsNonZeroOnBrokenScaffold(t *testing.T) {
	root := t.TempDir()
	var out bytes.Buffer
	if code := runInit([]string{root}, &out, &out); code != 0 {
		t.Fatalf("init: expected exit 0, got %d", code)
	}
	bad := "---\nid: G-001\ntype: goal\nstatus: proposed\nlinks: [G-999]\ntitle: First goal\n---\n"
	if err := os.WriteFile(filepath.Join(root, "docs", "goals", "G-001-first.md"), []byte(bad), 0o644); err != nil {
		t.Fatal(err)
	}
	if code := runInit([]string{root}, &out, &out); code != 0 { // re-index the new artifact
		t.Fatalf("re-init: expected exit 0, got %d", code)
	}
	if code := runValidate([]string{root}); code != 1 {
		t.Fatalf("expected exit 1 on dangling link, got %d", code)
	}
}

// AC-025: the report tells the user what was skipped and what to do next.
func TestAC025_InitReportsSkipsAndNextStep(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "AGENTS.md"), []byte("mine"), 0o644); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if code := runInit([]string{root}, &out, &out); code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	report := out.String()
	if !strings.Contains(report, "exists   AGENTS.md") {
		t.Fatalf("report does not name the skipped file:\n%s", report)
	}
	if !strings.Contains(report, "clue validate") {
		t.Fatalf("report does not point at the next step:\n%s", report)
	}
}

// Unit: a scaffold error exits 1 and lands on the error writer, not the
// report writer — scripted callers filter stderr.
func TestUnit_InitErrorGoesToErrWriterAndExitsOne(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "docs", "goals"), 0o755); err != nil {
		t.Fatal(err)
	}
	malformed := "# Goals\n\n<!-- clue:index:end -->\n"
	if err := os.WriteFile(filepath.Join(root, "docs", "goals", "README.md"), []byte(malformed), 0o644); err != nil {
		t.Fatal(err)
	}
	var out, errOut bytes.Buffer
	if code := runInit([]string{root}, &out, &errOut); code != 1 {
		t.Fatalf("expected exit 1 on malformed markers, got %d\nout: %s\nerr: %s", code, out.String(), errOut.String())
	}
	if !strings.Contains(errOut.String(), "docs/goals/README.md") {
		t.Fatalf("error writer does not name the offending file:\nout: %s\nerr: %s", out.String(), errOut.String())
	}
}

// Unit: a pre-existing docs folder without a README is named in the
// report, so the first validate is not red without warning.
func TestUnit_InitReportsMissingFolderReadme(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "docs", "notes"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "notes", "note.md"), []byte("# A note\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if code := runInit([]string{root}, &out, &out); code != 0 {
		t.Fatalf("expected exit 0, got %d\n%s", code, out.String())
	}
	if !strings.Contains(out.String(), "missing  docs/notes/README.md") {
		t.Fatalf("report does not name the missing README:\n%s", out.String())
	}
}
