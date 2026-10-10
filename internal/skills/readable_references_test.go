package skills

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAC232_UnitPositive_EveryStandaloneSkillRoutesToTheCanonicalReaderInstruction(t *testing.T) {
	root := t.TempDir()
	if e := Write(root); e != nil {
		t.Fatal(e)
	}
	var canonical string
	for _, output := range outputRoots {
		for _, name := range skillNames {
			entry, e := os.ReadFile(filepath.Join(root, filepath.FromSlash(output), name, "skill.md"))
			if e != nil {
				t.Fatal(e)
			}
			if !strings.Contains(string(entry), "Before writing reader-facing prose, a human report or a handoff, read [Readable references](references/readable-references.md)") {
				t.Fatalf("%s has no authoring route", name)
			}
			reference, e := os.ReadFile(filepath.Join(root, filepath.FromSlash(output), name, "references/readable-references.md"))
			if e != nil {
				t.Fatal(e)
			}
			if canonical == "" {
				canonical = string(reference)
			} else if canonical != string(reference) {
				t.Fatalf("%s has a second writing rule", name)
			}
			for _, want := range []string{"sentence itself", "descriptive link text", "Changelogs state", "frontmatter", "structured evidence", "same ID", "do not prove human comprehension"} {
				if !strings.Contains(string(reference), want) {
					t.Fatalf("%s misses %q", name, want)
				}
			}
		}
	}
	if drifts, e := Check(root); e != nil || len(drifts) != 0 {
		t.Fatalf("fresh distribution drift: %v %v", drifts, e)
	}
}

func TestAC232_UnitNegative_MissingOrEditedReaderInstructionsFailDistributionChecks(t *testing.T) {
	for _, output := range outputRoots {
		for _, edit := range []bool{false, true} {
			root := t.TempDir()
			if e := Write(root); e != nil {
				t.Fatal(e)
			}
			target := filepath.Join(root, filepath.FromSlash(output), "clue-analysis/references/readable-references.md")
			if edit {
				if e := os.WriteFile(target, []byte("an independently changed rule"), 0644); e != nil {
					t.Fatal(e)
				}
			} else if e := os.Remove(target); e != nil {
				t.Fatal(e)
			}
			drifts, e := Check(root)
			if e != nil || len(drifts) == 0 {
				t.Fatalf("drift ignored: %v %v", drifts, e)
			}
		}
	}
}
