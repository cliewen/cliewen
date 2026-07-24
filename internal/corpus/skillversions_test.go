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

func writeSkillNamed(t *testing.T, root, dir, manifest, content string) {
	t.Helper()
	p := filepath.Join(root, ".agents", "skills", dir, manifest)
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

func markedSkill(version string) string {
	return "---\ncliewen-skill: true\nversion: " + version + "\n---\n"
}

// AC-029: marked skills join the Cliewen set while an unmarked
// third-party skill is ignored.
func TestAC029_OwnershipMarkerScopesValidation(t *testing.T) {
	root := t.TempDir()
	writeSkill(t, root, "clue-delta", markedSkill("0.1.0"))
	writeSkill(t, root, "third-party", "---\nversion: 9.9.9\n---\n")
	if issues := checkSkillVersions(&Corpus{Root: root}, "0.1.0"); len(issues) != 0 {
		t.Fatalf("unmarked third-party skill affected validation: %v", issues)
	}
}

// AC-029 (negative): a present ownership marker must be boolean true.
func TestAC029_InvalidOwnershipMarkerFails(t *testing.T) {
	root := t.TempDir()
	writeSkill(t, root, "third-party", "---\ncliewen-skill: \"true\"\nversion: 0.1.0\n---\n")
	issues := checkSkillVersions(&Corpus{Root: root}, "dev")
	if !anyMsg(issues, "must be boolean true") {
		t.Fatalf("expected malformed ownership marker issue, got %v", issues)
	}
}

// AC-029 (diagnostics): a declared marker inside an unterminated YAML
// block is not silently treated as an unmarked third-party skill.
func TestAC029_UnterminatedOwnershipFrontmatterFails(t *testing.T) {
	root := t.TempDir()
	writeSkill(t, root, "third-party", "---\ncliewen-skill: true\nversion: 0.1.0\n")
	issues := checkSkillVersions(&Corpus{Root: root}, "dev")
	if !anyMsg(issues, "unterminated YAML block") {
		t.Fatalf("expected unterminated ownership frontmatter issue, got %v", issues)
	}
}

// AC-029 (negative): a nested key in malformed third-party
// frontmatter does not become a top-level ownership declaration.
func TestAC029_MalformedNestedMarkerDoesNotClaimThirdPartySkill(t *testing.T) {
	root := t.TempDir()
	writeSkill(t, root, "third-party", "---\nmetadata:\n  cliewen-skill: true\nbroken: [\n---\n")
	if issues := checkSkillVersions(&Corpus{Root: root}, "dev"); len(issues) != 0 {
		t.Fatalf("nested marker claimed an unmarked third-party skill: %v", issues)
	}
}

// AC-029 (diagnostics): a delimiter lookalike before a top-level marker
// does not hide a malformed Cliewen ownership declaration.
func TestAC029_DelimiterLookalikeDoesNotHideMalformedMarkedSkill(t *testing.T) {
	root := t.TempDir()
	writeSkill(t, root, "third-party", "---\nversion: [broken\n---not-a-delimiter\ncliewen-skill: true\n---\n")
	issues := checkSkillVersions(&Corpus{Root: root}, "dev")
	if !anyMsg(issues, "does not parse") {
		t.Fatalf("expected malformed marked frontmatter issue, got %v", issues)
	}
}

// AC-030: an unmarked skill in a canonical legacy slot fails toward
// reinstalling the Cliewen skill set.
func TestAC030_UnmarkedLegacyCliewenSkillFails(t *testing.T) {
	root := t.TempDir()
	writeSkill(t, root, "clue-delta", "---\nversion: 0.3.0\n---\n")
	issues := checkSkillVersions(&Corpus{Root: root}, "0.4.0")
	if !anyMsg(issues, "legacy Cliewen skill") || !anyMsg(issues, "reinstall") {
		t.Fatalf("expected legacy reinstall guidance, got %v", issues)
	}
}

// AC-030 (negative): an unmarked skill outside the reserved legacy names
// remains outside Cliewen's validation scope.
func TestAC030_UnmarkedNonCliewenSkillPasses(t *testing.T) {
	root := t.TempDir()
	writeSkill(t, root, "my-skill", "# no Cliewen frontmatter\n")
	if issues := checkSkillVersions(&Corpus{Root: root}, "0.4.0"); len(issues) != 0 {
		t.Fatalf("unmarked non-Cliewen skill should be ignored, got %v", issues)
	}
}

// AC-031: a marked Cliewen skill without a version stamp fails.
func TestAC031_MarkedSkillWithoutVersionFails(t *testing.T) {
	root := t.TempDir()
	writeSkill(t, root, "clue-delta", "---\ncliewen-skill: true\n---\n")
	if !anyMsg(checkSkillVersions(&Corpus{Root: root}, "dev"), "no version stamp") {
		t.Fatal("expected a missing-version-stamp issue")
	}
}

// AC-031 (negative): a marked, stamped skill raises no stamp issue.
func TestAC031_MarkedSkillWithVersionPasses(t *testing.T) {
	root := t.TempDir()
	writeSkill(t, root, "clue-delta", markedSkill("0.1.0"))
	if issues := checkSkillVersions(&Corpus{Root: root}, "dev"); len(issues) != 0 {
		t.Fatalf("marked, stamped skill should pass, got %v", issues)
	}
}

// AC-031 (diagnostics): a non-string version is named as such — YAML
// reads `version: 1.0` as a number.
func TestAC031_NonStringVersionIsNamedAsSuch(t *testing.T) {
	root := t.TempDir()
	writeSkill(t, root, "clue-delta", "---\ncliewen-skill: true\nversion: 1.0\n---\n")
	issues := checkSkillVersions(&Corpus{Root: root}, "dev")
	if !anyMsg(issues, "not a string") {
		t.Fatalf("expected a not-a-string issue, got %v", issues)
	}
	if anyMsg(issues, "no version stamp") {
		t.Fatal("a present-but-mistyped stamp must not be reported as missing")
	}
}

// AC-031 (diagnostics): malformed frontmatter that declares Cliewen
// ownership is a parse failure, not a missing stamp.
func TestAC031_MalformedMarkedFrontmatterIsNamedAsSuch(t *testing.T) {
	root := t.TempDir()
	writeSkill(t, root, "custom-cliewen", "---\ncliewen-skill: true\nversion: [unclosed\n---\n")
	issues := checkSkillVersions(&Corpus{Root: root}, "dev")
	if !anyMsg(issues, "does not parse") {
		t.Fatalf("expected a frontmatter parse issue, got %v", issues)
	}
	if anyMsg(issues, "no version stamp") {
		t.Fatal("unparseable frontmatter must not be reported as a missing stamp")
	}
}

// AC-032: marked skills that disagree on a version fail.
func TestAC032_DivergentMarkedSkillVersionsFail(t *testing.T) {
	root := t.TempDir()
	writeSkill(t, root, "clue-delta", markedSkill("0.1.0"))
	writeSkill(t, root, "clue-plan", markedSkill("0.2.0"))
	if !anyMsg(checkSkillVersions(&Corpus{Root: root}, "dev"), "disagrees") {
		t.Fatal("expected a disagreement issue across marked skills")
	}
}

// AC-032 (negative): marked skills that agree pass, regardless of an
// unmarked skill's version.
func TestAC032_ConsistentMarkedSkillVersionsPass(t *testing.T) {
	root := t.TempDir()
	writeSkill(t, root, "clue-delta", markedSkill("0.1.0"))
	writeSkill(t, root, "clue-plan", markedSkill("0.1.0"))
	writeSkill(t, root, "third-party", "---\nversion: 8.0.0\n---\n")
	if issues := checkSkillVersions(&Corpus{Root: root}, "dev"); len(issues) != 0 {
		t.Fatalf("consistent marked skills should pass, got %v", issues)
	}
}

// AC-032 (diagnostics): the marked outlier is named even when it sorts
// first — the reference is the version the majority carries.
func TestAC032_OutlierIsNamedEvenWhenSortedFirst(t *testing.T) {
	root := t.TempDir()
	writeSkill(t, root, "a-outlier", markedSkill("0.9.0"))
	writeSkill(t, root, "clue-delta", markedSkill("0.1.0"))
	writeSkill(t, root, "clue-plan", markedSkill("0.1.0"))
	issues := checkSkillVersions(&Corpus{Root: root}, "dev")
	if len(issues) != 1 {
		t.Fatalf("expected exactly the outlier to be reported, got %v", issues)
	}
	if !strings.Contains(issues[0].Path, "a-outlier") || !strings.Contains(issues[0].Msg, "0.9.0") {
		t.Fatalf("expected a-outlier (0.9.0) to be the skill named, got %v", issues[0])
	}
}

// AC-033: a released binary whose marked skills differ reports drift.
func TestAC033_ReleasedBinaryDriftFails(t *testing.T) {
	root := t.TempDir()
	writeSkill(t, root, "clue-delta", markedSkill("0.1.0"))
	if !anyMsg(checkSkillVersions(&Corpus{Root: root}, "0.2.0"), "drift") {
		t.Fatal("expected drift against a released binary")
	}
}

// AC-033 (negative): a dev build skips comparison, a matching release
// passes, and unmarked skills do not drift.
func TestAC033_DevMatchingAndUnmarkedSkillsDoNotDrift(t *testing.T) {
	root := t.TempDir()
	writeSkill(t, root, "clue-delta", markedSkill("0.1.0"))
	writeSkill(t, root, "third-party", "---\nversion: 9.9.9\n---\n")
	if anyMsg(checkSkillVersions(&Corpus{Root: root}, "dev"), "drift") {
		t.Fatal("a dev build must not report drift")
	}
	if issues := checkSkillVersions(&Corpus{Root: root}, "0.1.0"); len(issues) != 0 {
		t.Fatalf("matching release with unmarked neighbor should pass, got %v", issues)
	}
}

// AC-037: a manifest named SKILL.md joins the managed set exactly as a
// lowercase skill.md would, so the verdict does not depend on the host
// filesystem. The drift finding names the actual manifest file, which the
// fixed-path lookup could not do.
func TestAC037_UppercaseManifestJoinsManagedSet(t *testing.T) {
	root := t.TempDir()
	writeSkillNamed(t, root, "clue-delta", "SKILL.md", markedSkill("0.1.0"))
	issues := checkSkillVersions(&Corpus{Root: root}, "0.2.0")
	if !anyMsg(issues, "drift") {
		t.Fatalf("an uppercase-manifest skill must join the set and drift against 0.2.0, got %v", issues)
	}
	named := false
	for _, i := range issues {
		if strings.Contains(i.Path, "SKILL.md") {
			named = true
		}
	}
	if !named {
		t.Fatalf("the drift finding must name the actual manifest SKILL.md, got %v", issues)
	}
}

// AC-037 (negative): a single directory holding two case-variants of the
// manifest name is reported as an ambiguity rather than silently resolving to
// one. Only a case-sensitive filesystem can hold both files.
func TestAC037_TwoCaseVariantsAreReportedAmbiguous(t *testing.T) {
	root := t.TempDir()
	writeSkillNamed(t, root, "clue-delta", "skill.md", markedSkill("0.1.0"))
	writeSkillNamed(t, root, "clue-delta", "SKILL.md", markedSkill("0.1.0"))
	entries, err := os.ReadDir(filepath.Join(root, ".agents", "skills", "clue-delta"))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) < 2 {
		t.Skip("case-insensitive filesystem cannot hold two manifest case-variants")
	}
	if !anyMsg(checkSkillVersions(&Corpus{Root: root}, "0.1.0"), "case-variants") {
		t.Fatal("expected an ambiguity issue naming the multiple case-variants")
	}
}

// AC-037 (negative): the ambiguity is reported for an unmarked third-party
// skill too. Ownership lives inside the manifest that cannot be chosen, so
// this is the one place an unmarked directory does not simply drop out
// (ADR-028).
func TestAC037_AmbiguityIsReportedForAnUnmarkedSkill(t *testing.T) {
	root := t.TempDir()
	writeSkillNamed(t, root, "third-party", "skill.md", "---\nversion: 9.9.9\n---\n")
	writeSkillNamed(t, root, "third-party", "SKILL.md", "---\nversion: 9.9.9\n---\n")
	entries, err := os.ReadDir(filepath.Join(root, ".agents", "skills", "third-party"))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) < 2 {
		t.Skip("case-insensitive filesystem cannot hold two manifest case-variants")
	}
	issues := checkSkillVersions(&Corpus{Root: root}, "0.1.0")
	if !anyMsg(issues, "case-variants") {
		t.Fatalf("an unmarked skill's manifest ambiguity must still be reported, got %v", issues)
	}
	if len(issues) != 1 {
		t.Fatalf("the ambiguity ends the directory; no further finding is expected, got %v", issues)
	}
}

// AC-037: only an entry resolving to a regular file is a manifest. A
// directory named skill.md sitting in a legacy Cliewen slot is not one, so it
// neither enrolls nor trips the reinstall finding (ADR-028).
func TestAC037_ADirectoryNamedSkillMdIsNotAManifest(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".agents", "skills", "clue-delta", "skill.md"), 0o755); err != nil {
		t.Fatal(err)
	}
	if issues := checkSkillVersions(&Corpus{Root: root}, "0.1.0"); len(issues) != 0 {
		t.Fatalf("a directory named skill.md is not a manifest, got %v", issues)
	}
}

// Unit: a repo with no skills folder has nothing to check.
func TestUnit_NoSkillsFolderIsClean(t *testing.T) {
	if issues := checkSkillVersions(&Corpus{Root: t.TempDir()}, "0.1.0"); len(issues) != 0 {
		t.Fatalf("expected no issues without a skills folder, got %v", issues)
	}
}
