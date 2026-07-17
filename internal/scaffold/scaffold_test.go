package scaffold

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cliewen/cliewen/internal/corpus"
)

func runInto(t *testing.T) (string, *Report) {
	t.Helper()
	root := t.TempDir()
	rep, err := Run(root)
	if err != nil {
		t.Fatal(err)
	}
	return root, rep
}

func validateAt(t *testing.T, root string) []corpus.Issue {
	t.Helper()
	c, issues := corpus.Scan(root)
	return append(issues, corpus.Validate(c, corpus.Options{})...)
}

// AC-002: init produces a corpus that validate accepts unchanged.
func TestAC002_InitOutputPassesValidateUnchanged(t *testing.T) {
	root, _ := runInto(t)
	if issues := validateAt(t, root); len(issues) > 0 {
		t.Fatalf("expected a green corpus, got issues: %v", issues)
	}
	// The pair is coherent for a released binary too: the emitted skills'
	// stamp matches what a release stamped with the same version expects.
	version, err := PairVersion()
	if err != nil {
		t.Fatal(err)
	}
	c, _ := corpus.Scan(root)
	if issues := corpus.Validate(c, corpus.Options{Version: version}); len(issues) > 0 {
		t.Fatalf("emitted skills drift from pair version %s: %v", version, issues)
	}
}

// AC-002 negative: the green result is not vacuous — validate really
// judges the generated corpus and catches damage to it.
func TestAC002_DamagedScaffoldIsCaught(t *testing.T) {
	root, _ := runInto(t)
	if err := os.Remove(filepath.Join(root, "docs", "goals", "README.md")); err != nil {
		t.Fatal(err)
	}
	if issues := validateAt(t, root); len(issues) == 0 {
		t.Fatal("expected issues after damaging the scaffold, got none")
	}
}

// AC-003: a broken link in a scaffolded corpus fails loudly, naming the
// offending file and the missing ID.
func TestAC003_BrokenLinkNamesFileAndMissingID(t *testing.T) {
	root, _ := runInto(t)
	bad := "docs/goals/G-001-first.md"
	content := "---\nid: G-001\ntype: goal\nstatus: proposed\nlinks: [G-999]\ntitle: First goal\n---\n"
	if err := os.WriteFile(filepath.Join(root, filepath.FromSlash(bad)), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Run(root); err != nil { // re-index so only the dangling link is at fault
		t.Fatal(err)
	}
	issues := validateAt(t, root)
	if len(issues) == 0 {
		t.Fatal("expected the dangling link to be reported")
	}
	found := false
	for _, is := range issues {
		if is.Path == bad && strings.Contains(is.Msg, "G-999") {
			found = true
		}
	}
	if !found {
		t.Fatalf("no issue names both %s and G-999: %v", bad, issues)
	}
}

// AC-003 negative: resolving the link restores green — the failure was
// the link, not the scaffold.
func TestAC003_ResolvedLinkRestoresGreen(t *testing.T) {
	root, _ := runInto(t)
	content := "---\nid: G-001\ntype: goal\nstatus: proposed\nlinks: []\ntitle: First goal\n---\n"
	if err := os.WriteFile(filepath.Join(root, "docs", "goals", "G-001-first.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Run(root); err != nil {
		t.Fatal(err)
	}
	if issues := validateAt(t, root); len(issues) > 0 {
		t.Fatalf("expected green after fixing, got: %v", issues)
	}
}

// AC-024: a re-run regenerates README index blocks from folder contents
// and leaves prose outside the markers alone.
func TestAC024_RerunIndexesNewArtifactAndKeepsProse(t *testing.T) {
	root, _ := runInto(t)
	readme := filepath.Join(root, "docs", "goals", "README.md")
	prose := "# Goals\n\nHand-written local prose that must survive.\n\n<!-- clue:index:start -->\n<!-- clue:index:end -->\n"
	if err := os.WriteFile(readme, []byte(prose), 0o644); err != nil {
		t.Fatal(err)
	}
	artifact := "---\nid: G-001\ntype: goal\nstatus: proposed\nlinks: []\ntitle: First goal\n---\n"
	if err := os.WriteFile(filepath.Join(root, "docs", "goals", "G-001-first.md"), []byte(artifact), 0o644); err != nil {
		t.Fatal(err)
	}
	rep, err := Run(root)
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
	if !strings.Contains(text, "Hand-written local prose that must survive.") {
		t.Fatalf("prose outside the markers was touched:\n%s", text)
	}
	indexed := false
	for _, p := range rep.Indexed {
		if p == "docs/goals/README.md" {
			indexed = true
		}
	}
	if !indexed {
		t.Fatalf("expected docs/goals/README.md in the indexed report, got %v", rep.Indexed)
	}
	if issues := validateAt(t, root); len(issues) > 0 {
		t.Fatalf("expected green after re-index, got: %v", issues)
	}
}

// AC-024: a taxonomy README that predates init and has no markers gains
// an appended index block — prose intact, validate green afterwards.
func TestAC024_MarkerlessReadmeGainsAppendedBlock(t *testing.T) {
	root := t.TempDir()
	prose := "# My docs\n\nPre-existing prose without any markers.\n"
	if err := os.MkdirAll(filepath.Join(root, "docs", "goals"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "README.md"), []byte(prose), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "goals", "README.md"), []byte("# Goals, my way\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Run(root); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filepath.Join(root, "docs", "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(got)
	if !strings.Contains(text, "Pre-existing prose without any markers.") {
		t.Fatalf("pre-existing prose was lost:\n%s", text)
	}
	if !strings.Contains(text, indexStart) || !strings.Contains(text, indexEnd) {
		t.Fatalf("no index block was appended:\n%s", text)
	}
	if issues := validateAt(t, root); len(issues) > 0 {
		t.Fatalf("expected green after marker append, got: %v", issues)
	}
}

// AC-024 negative: with nothing new to index, a re-run changes no file.
func TestAC024_RerunOnUnchangedTreeIsANoOp(t *testing.T) {
	root, _ := runInto(t)
	before := snapshot(t, root)
	rep, err := Run(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(rep.Created) != 0 || len(rep.Indexed) != 0 {
		t.Fatalf("expected nothing created or re-indexed, got created=%v indexed=%v", rep.Created, rep.Indexed)
	}
	after := snapshot(t, root)
	if len(before) != len(after) {
		t.Fatalf("file count changed: %d -> %d", len(before), len(after))
	}
	for p, b := range before {
		if after[p] != b {
			t.Fatalf("file %s changed on a no-op re-run", p)
		}
	}
}

// AC-025: an existing file is never overwritten — it is skipped and the
// report says so.
func TestAC025_ExistingFileIsSkippedAndReported(t *testing.T) {
	root := t.TempDir()
	own := "# My own routing hub\n"
	if err := os.WriteFile(filepath.Join(root, "AGENTS.md"), []byte(own), 0o644); err != nil {
		t.Fatal(err)
	}
	rep, err := Run(root)
	if err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filepath.Join(root, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != own {
		t.Fatalf("existing AGENTS.md was overwritten:\n%s", got)
	}
	found := false
	for _, p := range rep.Skipped {
		if p == "AGENTS.md" {
			found = true
		}
	}
	if !found {
		t.Fatalf("AGENTS.md not reported as skipped: %v", rep.Skipped)
	}
}

// AC-025 negative: skipping is per file, not per run — everything the
// existing file did not shadow is still created.
func TestAC025_SkipIsPerFileNotPerRun(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "AGENTS.md"), []byte("mine"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Run(root); err != nil {
		t.Fatal(err)
	}
	for _, rel := range []string{"docs/README.md", "docs/constraints/C-001-no-hard-wrapped-markdown.md", ".agents/skills/clue-delta/skill.md", ".claude/skills/clue-delta/SKILL.md", ".github/workflows/clue.yml"} {
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(rel))); err != nil {
			t.Fatalf("%s was not created: %v", rel, err)
		}
	}
}

// The embedded skill copies must stay byte-identical to the canonical
// skills in .agents/skills — the duplication go:embed forces (dotted
// directories are invisible to the Go toolchain) is held by this test.
func TestSanity_EmbeddedSkillsMatchCanonicalSkills(t *testing.T) {
	canonical := filepath.Join("..", "..", ".agents", "skills")
	seen := map[string]bool{}
	err := fs.WalkDir(templates, "templates/skills", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel := strings.TrimPrefix(p, "templates/skills/")
		seen[rel] = true
		embedded, _ := templates.ReadFile(p)
		disk, derr := os.ReadFile(filepath.Join(canonical, filepath.FromSlash(rel)))
		if derr != nil {
			t.Errorf("embedded %s has no canonical twin: %v", rel, derr)
			return nil
		}
		if strings.ReplaceAll(string(disk), "\r\n", "\n") != strings.ReplaceAll(string(embedded), "\r\n", "\n") {
			t.Errorf("embedded %s differs from .agents/skills — re-copy the skills into internal/scaffold/templates/skills", rel)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	_ = filepath.WalkDir(canonical, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(canonical, p)
		if !seen[filepath.ToSlash(rel)] {
			t.Errorf("canonical skill file %s is missing from the embedded templates", filepath.ToSlash(rel))
		}
		return nil
	})
}

// The CI template's version pin comes from the embedded skills' stamp;
// no placeholder may survive into the emitted workflow.
func TestUnit_WorkflowVersionSubstituted(t *testing.T) {
	root, _ := runInto(t)
	data, err := os.ReadFile(filepath.Join(root, ".github", "workflows", "clue.yml"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), versionPlaceholder) {
		t.Fatal("workflow still contains the version placeholder")
	}
	version, err := PairVersion()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "CLUE_VERSION: "+version) {
		t.Fatalf("workflow does not pin CLUE_VERSION to %s", version)
	}
}

// Unit: a lone or reversed index marker is ambiguous — Run refuses with
// an error naming the file instead of guessing at the block's bounds,
// and the file is left byte-for-byte untouched (the prose promise).
func TestUnit_MalformedMarkersErrorAndLeaveFileUntouched(t *testing.T) {
	cases := map[string]string{
		"lone end":      "# Goals\n\nProse that must survive.\n\n" + indexEnd + "\n",
		"lone start":    "# Goals\n\n" + indexStart + "\n\nProse that must survive.\n",
		"reversed pair": "# Goals\n\n" + indexEnd + "\n\nProse that must survive.\n\n" + indexStart + "\n",
	}
	for name, content := range cases {
		t.Run(name, func(t *testing.T) {
			root := t.TempDir()
			readme := filepath.Join(root, "docs", "goals", "README.md")
			if err := os.MkdirAll(filepath.Dir(readme), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(readme, []byte(content), 0o644); err != nil {
				t.Fatal(err)
			}
			_, err := Run(root)
			if err == nil {
				t.Fatal("expected an error on malformed markers, got none")
			}
			if !strings.Contains(err.Error(), "docs/goals/README.md") {
				t.Fatalf("error does not name the offending file: %v", err)
			}
			got, rerr := os.ReadFile(readme)
			if rerr != nil {
				t.Fatal(rerr)
			}
			if string(got) != content {
				t.Fatalf("malformed README was modified:\nbefore: %q\nafter:  %q", content, got)
			}
		})
	}
}

// Unit: a pre-existing docs folder without the README validate requires
// is named in the report — init does not invent the file.
func TestUnit_MissingFolderReadmeIsReportedNotInvented(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "docs", "notes"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "notes", "note.md"), []byte("# A note\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	rep, err := Run(root)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, p := range rep.MissingReadmes {
		if p == "docs/notes/README.md" {
			found = true
		}
	}
	if !found {
		t.Fatalf("docs/notes/README.md not reported missing: %v", rep.MissingReadmes)
	}
	if _, err := os.Stat(filepath.Join(root, "docs", "notes", "README.md")); err == nil {
		t.Fatal("init invented a README it must not create")
	}
}

// Unit: a CRLF checkout keeps its line endings — prose outside the
// markers is byte-for-byte untouched and generated index lines adopt
// the file's own style (the AC-024 prose promise on Windows).
func TestUnit_CrlfReadmeKeepsItsLineEndings(t *testing.T) {
	root, _ := runInto(t)
	readme := filepath.Join(root, "docs", "goals", "README.md")
	proseCRLF := "# Goals\r\n\r\nProse that must stay CRLF.\r\n\r\n"
	crlf := proseCRLF + indexStart + "\r\n" + indexEnd + "\r\n"
	if err := os.WriteFile(readme, []byte(crlf), 0o644); err != nil {
		t.Fatal(err)
	}
	artifact := "---\nid: G-001\ntype: goal\nstatus: proposed\nlinks: []\ntitle: First goal\n---\n"
	if err := os.WriteFile(filepath.Join(root, "docs", "goals", "G-001-first.md"), []byte(artifact), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Run(root); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(readme)
	if err != nil {
		t.Fatal(err)
	}
	text := string(got)
	if !strings.HasPrefix(text, proseCRLF) {
		t.Fatalf("CRLF prose outside the markers was rewritten:\n%q", text)
	}
	if !strings.Contains(text, "](G-001-first.md)\r\n") {
		t.Fatalf("generated index line does not use the file's CRLF endings:\n%q", text)
	}
}

// Unit: missing READMEs are found at every depth, the same recursive
// reading validate applies — not just at taxonomy level.
func TestUnit_MissingReadmeIsReportedRecursively(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "docs", "notes", "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "notes", "sub", "note.md"), []byte("# A note\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	rep, err := Run(root)
	if err != nil {
		t.Fatal(err)
	}
	missing := map[string]bool{}
	for _, p := range rep.MissingReadmes {
		missing[p] = true
	}
	for _, want := range []string{"docs/notes/README.md", "docs/notes/sub/README.md"} {
		if !missing[want] {
			t.Fatalf("%s not reported missing: %v", want, rep.MissingReadmes)
		}
	}
}

func snapshot(t *testing.T, root string) map[string]string {
	t.Helper()
	files := map[string]string{}
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		data, rerr := os.ReadFile(p)
		if rerr != nil {
			return rerr
		}
		rel, _ := filepath.Rel(root, p)
		files[filepath.ToSlash(rel)] = string(data)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return files
}
