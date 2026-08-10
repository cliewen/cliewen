package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAC133_UnitPositive_ContextCommandPrintsCompleteArtifacts(t *testing.T) {
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

func TestAC133_UnitNegative_ContextCommandNamesUnknownIDWithoutArtifactOutput(t *testing.T) {
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

// The bound is only honest when the output says where it stopped. A reader who
// cannot see that neighbours were withheld cannot decide to widen, which is the
// silent cap PDR-034 refused.
func TestAC133_UnitPositive_ContextCommandNamesTheFrontierAndWidens(t *testing.T) {
	root := t.TempDir()
	writeContextFile(t, root, "docs/goals/G-101.md", "---\nid: G-101\ntype: goal\nstatus: accepted\nlinks: []\ntitle: Distant goal\n---\n\n# Distant goal body\n")
	writeContextFile(t, root, "docs/plans/P-101.md", "---\nid: P-101\ntype: plan\nstatus: active\nlinks: [G-101]\ntitle: Plan\n---\n\n| M-101 | Do it | todo |\n")
	writeContextFile(t, root, "docs/capabilities/CAP-101/README.md", "---\nid: CAP-101\ntype: capability\nstatus: active\nlinks: [P-101]\ntitle: Capability\ngoal: G-101\n---\n\n# Capability\n")

	var out, errOut bytes.Buffer
	if code := runContext([]string{"--stats", "CAP-101", root}, &out, &errOut); code != 0 {
		t.Fatalf("context exit = %d, stderr = %s", code, errOut.String())
	}
	got := out.String()
	if strings.Contains(got, "# Distant goal body") {
		t.Fatalf("an artifact two hops out was printed by the default slice:\n%s", got)
	}
	if !strings.Contains(got, "G-101 | Distant goal") {
		t.Fatalf("the frontier does not name the artifact the bound held back:\n%s", got)
	}
	if !strings.Contains(got, "--depth=all") || !strings.Contains(got, "byte(s) printed") {
		t.Fatalf("the slice does not state how to widen or what it cost:\n%s", got)
	}

	out.Reset()
	errOut.Reset()
	if code := runContext([]string{"--depth=all", "CAP-101", root}, &out, &errOut); code != 0 {
		t.Fatalf("widened context exit = %d, stderr = %s", code, errOut.String())
	}
	got = out.String()
	if !strings.Contains(got, "# Distant goal body") {
		t.Fatalf("--depth=all did not reach the whole closure:\n%s", got)
	}
	if strings.Contains(got, "frontier:") {
		t.Fatalf("--depth=all held something back:\n%s", got)
	}
}

// A depth that is neither a number nor "all" is a usage error, not a silent
// fallback to the default: a mistyped bound must never look like a full read.
func TestAC133_UnitNegative_ContextCommandRejectsAnUnreadableDepth(t *testing.T) {
	var out, errOut bytes.Buffer
	if code := runContext([]string{"--depth=deep", "CAP-404", t.TempDir()}, &out, &errOut); code != 2 {
		t.Fatalf("context exit = %d, want 2", code)
	}
	if out.Len() != 0 {
		t.Fatalf("an unreadable depth emitted output: %q", out.String())
	}
	if !strings.Contains(errOut.String(), "--depth") {
		t.Fatalf("stderr does not name the offending flag: %s", errOut.String())
	}
}
