package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAC053_UnitPositive_ContextCommandPrintsCompleteArtifacts(t *testing.T) {
	root := t.TempDir()
	writeContextFile(t, root, "docs/goals/G-101.md", "---\nid: G-101\ntype: goal\nstatus: accepted\nlinks: []\ntitle: Goal\n---\n\n# Complete goal body\n")
	writeContextFile(t, root, "docs/plans/P-101.md", "---\nid: P-101\ntype: plan\nstatus: active\nlinks: [G-101]\ntitle: Plan\n---\n\n| M-101 | Do it | todo |\n")

	var out, errOut bytes.Buffer
	if code := runContext([]string{"M-101", root}, &out, &errOut); code != 0 {
		t.Fatalf("context exit = %d, stderr = %s", code, errOut.String())
	}
	got := out.String()
	if !strings.HasPrefix(got, "===== P-101 | docs/plans/P-101.md =====\n---\n") {
		t.Fatalf("root artifact was not emitted first with complete frontmatter:\n%s", got)
	}
	if !strings.Contains(got, "# Complete goal body") {
		t.Fatalf("linked artifact body missing:\n%s", got)
	}
}

func TestAC053_UnitNegative_ContextCommandNamesUnknownIDWithoutArtifactOutput(t *testing.T) {
	var out, errOut bytes.Buffer
	if code := runContext([]string{"CAP-404", t.TempDir()}, &out, &errOut); code != 1 {
		t.Fatalf("context exit = %d, want 1", code)
	}
	if out.Len() != 0 {
		t.Fatalf("unknown ID emitted artifact output: %q", out.String())
	}
	if !strings.Contains(errOut.String(), "CAP-404") {
		t.Fatalf("stderr does not name unknown ID: %s", errOut.String())
	}
}

func writeContextFile(t *testing.T, root, rel, content string) {
	t.Helper()
	full := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
