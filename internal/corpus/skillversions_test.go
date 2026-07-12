package corpus

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeSkill(t *testing.T, root, name, content string) {
	t.Helper()
	p := filepath.Join(root, ".agents", "skills", name, "skill.md")
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func anyMsg(issues []Issue, sub string) bool {
	for _, i := range issues {
		if strings.Contains(i.Msg, sub) {
			return true
		}
	}
	return false
}

// AC-020: a skill without a version stamp fails.
func TestAC020_SkillWithoutVersionFails(t *testing.T) {
	root := t.TempDir()
	writeSkill(t, root, "clue-delta", "# clue-delta\n\nno frontmatter, no stamp\n")
	if !anyMsg(checkSkillVersions(&Corpus{Root: root}, "dev"), "no version stamp") {
		t.Fatal("expected a missing-version-stamp issue")
	}
}

// AC-020 (negative): a stamped skill raises no missing-stamp issue.
func TestAC020_SkillWithVersionPasses(t *testing.T) {
	root := t.TempDir()
	writeSkill(t, root, "clue-delta", "---\nversion: 0.1.0\n---\n\n# clue-delta\n")
	if anyMsg(checkSkillVersions(&Corpus{Root: root}, "dev"), "no version stamp") {
		t.Fatal("a stamped skill should not be flagged")
	}
}

// AC-021: skills that disagree on a version fail.
func TestAC021_DivergentSkillVersionsFail(t *testing.T) {
	root := t.TempDir()
	writeSkill(t, root, "clue-delta", "---\nversion: 0.1.0\n---\n")
	writeSkill(t, root, "clue-plan", "---\nversion: 0.2.0\n---\n")
	if !anyMsg(checkSkillVersions(&Corpus{Root: root}, "dev"), "disagrees") {
		t.Fatal("expected a disagreement issue across skills")
	}
}

// AC-021 (negative): skills that agree pass the set-consistency check.
func TestAC021_ConsistentSkillVersionsPass(t *testing.T) {
	root := t.TempDir()
	writeSkill(t, root, "clue-delta", "---\nversion: 0.1.0\n---\n")
	writeSkill(t, root, "clue-plan", "---\nversion: 0.1.0\n---\n")
	if anyMsg(checkSkillVersions(&Corpus{Root: root}, "dev"), "disagrees") {
		t.Fatal("agreeing skills should not be flagged")
	}
}

// AC-022: a released binary whose skills differ from it reports drift.
func TestAC022_ReleasedBinaryDriftFails(t *testing.T) {
	root := t.TempDir()
	writeSkill(t, root, "clue-delta", "---\nversion: 0.1.0\n---\n")
	if !anyMsg(checkSkillVersions(&Corpus{Root: root}, "0.2.0"), "drift") {
		t.Fatal("expected drift against a released binary")
	}
}

// AC-022 (negative): a dev build skips the comparison, and a matching
// release does not drift.
func TestAC022_DevSkipsDriftAndMatchingReleasePasses(t *testing.T) {
	root := t.TempDir()
	writeSkill(t, root, "clue-delta", "---\nversion: 0.1.0\n---\n")
	if anyMsg(checkSkillVersions(&Corpus{Root: root}, "dev"), "drift") {
		t.Fatal("a dev build must not report drift")
	}
	if anyMsg(checkSkillVersions(&Corpus{Root: root}, "0.1.0"), "drift") {
		t.Fatal("a release matching the skills must not drift")
	}
}

// Unit: a repo with no skills folder has nothing to check.
func TestUnit_NoSkillsFolderIsClean(t *testing.T) {
	if issues := checkSkillVersions(&Corpus{Root: t.TempDir()}, "0.1.0"); len(issues) != 0 {
		t.Fatalf("expected no issues without a skills folder, got %v", issues)
	}
}
