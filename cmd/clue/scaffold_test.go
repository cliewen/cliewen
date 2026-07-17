package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// AC-026: the command-level path — scaffold re-indexes a grown corpus
// and validate accepts the result.
func TestAC026_ScaffoldThenValidateExitsZero(t *testing.T) {
	root := t.TempDir()
	var out bytes.Buffer
	if code := runInit([]string{root}, &out, &out); code != 0 {
		t.Fatalf("init: expected exit 0, got %d\n%s", code, out.String())
	}
	artifact := "---\nid: G-001\ntype: goal\nstatus: proposed\nlinks: []\ntitle: First goal\n---\n"
	if err := os.WriteFile(filepath.Join(root, "docs", "goals", "G-001-first.md"), []byte(artifact), 0o644); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	if code := runScaffold([]string{root}, &out, &out); code != 0 {
		t.Fatalf("scaffold: expected exit 0, got %d\n%s", code, out.String())
	}
	if !strings.Contains(out.String(), "indexed  docs/goals/README.md") {
		t.Fatalf("report does not name the regenerated index:\n%s", out.String())
	}
	if code := runValidate([]string{root}); code != 0 {
		t.Fatalf("validate after scaffold: expected exit 0, got %d", code)
	}
}

// AC-027: without a docs tree the command exits 1 with the error on the
// error writer — and materializes nothing.
func TestAC027_ScaffoldWithoutDocsExitsOneToStderr(t *testing.T) {
	root := t.TempDir()
	var out, errOut bytes.Buffer
	if code := runScaffold([]string{root}, &out, &errOut); code != 1 {
		t.Fatalf("expected exit 1, got %d\nout: %s\nerr: %s", code, out.String(), errOut.String())
	}
	if !strings.Contains(errOut.String(), "no docs tree") {
		t.Fatalf("error writer does not explain the failure:\nout: %s\nerr: %s", out.String(), errOut.String())
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("scaffold created files: %v", entries)
	}
}
