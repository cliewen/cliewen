package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// AC-026: Regen references a new artifact in the taxonomy README index,
// keeps prose outside the markers, and the result is green.
func TestAC026_RegenIndexesNewArtifactAndKeepsProse(t *testing.T) {
	root, _ := runInto(t)
	readme := filepath.Join(root, "docs", "goals", "README.md")
	prose := "# Goals\n\nHand-written prose that must survive.\n\n<!-- clue:index:start -->\n<!-- clue:index:end -->\n"
	if err := os.WriteFile(readme, []byte(prose), 0o644); err != nil {
		t.Fatal(err)
	}
	artifact := "---\nid: G-001\ntype: goal\nstatus: proposed\nlinks: []\ntitle: First goal\n---\n"
	if err := os.WriteFile(filepath.Join(root, "docs", "goals", "G-001-first.md"), []byte(artifact), 0o644); err != nil {
		t.Fatal(err)
	}
	rep, err := Regen(root)
	if err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(readme)
	if err != nil {
		t.Fatal(err)
	}
	text := string(got)
	if !strings.Contains(text, "](G-001-first.md)") {
		t.Fatalf("index block does not reference the new artifact:\n%s", text)
	}
	if !strings.Contains(text, "Hand-written prose that must survive.") {
		t.Fatalf("prose outside the markers was touched:\n%s", text)
	}
	found := false
	for _, p := range rep.Indexed {
		if p == "docs/goals/README.md" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected docs/goals/README.md in the indexed report, got %v", rep.Indexed)
	}
	if issues := validateAt(t, root); len(issues) > 0 {
		t.Fatalf("expected green after regen, got: %v", issues)
	}
}

// AC-026 negative: with nothing new to index, Regen changes no file.
func TestAC026_RegenOnUnchangedTreeIsANoOp(t *testing.T) {
	root, _ := runInto(t)
	before := snapshot(t, root)
	rep, err := Regen(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(rep.Indexed) != 0 {
		t.Fatalf("expected nothing re-indexed, got %v", rep.Indexed)
	}
	after := snapshot(t, root)
	if len(before) != len(after) {
		t.Fatalf("file count changed: %d -> %d", len(before), len(after))
	}
	for p, b := range before {
		if after[p] != b {
			t.Fatalf("file %s changed on a no-op regen", p)
		}
	}
}

// AC-027: Regen touches only taxonomy README index blocks — no file is
// created and no non-README file is modified.
func TestAC027_RegenTouchesOnlyTaxonomyReadmes(t *testing.T) {
	root, _ := runInto(t)
	artifact := "---\nid: G-001\ntype: goal\nstatus: proposed\nlinks: []\ntitle: First goal\n---\n"
	if err := os.WriteFile(filepath.Join(root, "docs", "goals", "G-001-first.md"), []byte(artifact), 0o644); err != nil {
		t.Fatal(err)
	}
	before := snapshot(t, root)
	if _, err := Regen(root); err != nil {
		t.Fatal(err)
	}
	after := snapshot(t, root)
	if len(before) != len(after) {
		t.Fatalf("regen created or deleted files: %d -> %d", len(before), len(after))
	}
	for p, b := range before {
		if after[p] != b && !strings.HasSuffix(p, "README.md") {
			t.Fatalf("non-README file %s was modified", p)
		}
	}
}

// AC-027 negative: a root without a docs tree is a loud error and
// nothing is created — scaffold never materializes.
func TestAC027_NoDocsTreeIsAnErrorAndCreatesNothing(t *testing.T) {
	root := t.TempDir()
	if _, err := Regen(root); err == nil {
		t.Fatal("expected an error on a root without docs, got none")
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("scaffold created files in an empty root: %v", entries)
	}
}
