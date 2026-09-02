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
	if issues := activatedValidateAt(t, root); len(issues) > 0 {
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

// AC-026: a curated entry linking a live descendant of a subfolder
// covers that subfolder for the validator — the regenerator keeps it
// instead of replacing it with the generated README link.
func TestAC026_CuratedDescendantEntrySurvives(t *testing.T) {
	root, _ := runInto(t)
	artifact := "---\nid: G-001\ntype: goal\nstatus: proposed\nlinks: []\ntitle: First goal\n---\n"
	if err := os.WriteFile(filepath.Join(root, "docs", "goals", "G-001-first.md"), []byte(artifact), 0o644); err != nil {
		t.Fatal(err)
	}
	rootReadme := filepath.Join(root, "docs", "README.md")
	raw, err := os.ReadFile(rootReadme)
	if err != nil {
		t.Fatal(err)
	}
	curated := "- [Goals, curated](goals/G-001-first.md)"
	text := strings.Replace(string(raw), "- [goals/](goals/README.md) — G-xxx: who wants it, why", curated, 1)
	if !strings.Contains(text, curated) {
		t.Fatalf("fixture setup failed — generated goals line not found in:\n%s", raw)
	}
	if err := os.WriteFile(rootReadme, []byte(text), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Regen(root); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(rootReadme)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(got), curated) {
		t.Fatalf("curated descendant entry was replaced:\n%s", got)
	}
	if strings.Contains(string(got), "](goals/README.md)") {
		t.Fatalf("generated entry was appended although the subfolder is covered:\n%s", got)
	}
	if issues := activatedValidateAt(t, root); len(issues) > 0 {
		t.Fatalf("expected green with the curated entry, got: %v", issues)
	}
}

// AC-026: a curated line carrying several links covers every target it
// names — no generated duplicate for the later links, and a line whose
// first link is informational survives on the strength of its second.
func TestAC026_MultiLinkCuratedLineCoversAllTargets(t *testing.T) {
	root, _ := runInto(t)
	for _, a := range []struct{ file, id string }{
		{"G-001-first.md", "G-001"},
		{"G-002-second.md", "G-002"},
	} {
		content := "---\nid: " + a.id + "\ntype: goal\nstatus: proposed\nlinks: []\ntitle: A goal\n---\n"
		if err := os.WriteFile(filepath.Join(root, "docs", "goals", a.file), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	readme := filepath.Join(root, "docs", "goals", "README.md")
	curated := "- [overview](../README.md): [first](G-001-first.md) and [second](G-002-second.md)"
	prose := "# Goals\n\n<!-- clue:index:start -->\n" + curated + "\n<!-- clue:index:end -->\n"
	if err := os.WriteFile(readme, []byte(prose), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Regen(root); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(readme)
	if err != nil {
		t.Fatal(err)
	}
	text := string(got)
	if !strings.Contains(text, curated) {
		t.Fatalf("multi-link curated line was replaced:\n%s", text)
	}
	for _, target := range []string{"](G-001-first.md)", "](G-002-second.md)"} {
		if strings.Count(text, target) != 1 {
			t.Fatalf("generated duplicate for already-covered %s:\n%s", target, text)
		}
	}
	if issues := activatedValidateAt(t, root); len(issues) > 0 {
		t.Fatalf("expected green with the curated line, got: %v", issues)
	}
}

// AC-027: Regen touches only the taxonomy READMEs checkIndexes judges —
// docs/README.md and docs/<folder>/README.md; nothing is created, and a
// nested README stays byte-identical.
func TestAC027_RegenTouchesOnlyTaxonomyReadmes(t *testing.T) {
	root, _ := runInto(t)
	artifact := "---\nid: G-001\ntype: goal\nstatus: proposed\nlinks: []\ntitle: First goal\n---\n"
	if err := os.WriteFile(filepath.Join(root, "docs", "goals", "G-001-first.md"), []byte(artifact), 0o644); err != nil {
		t.Fatal(err)
	}
	nested := filepath.Join(root, "docs", "goals", "nested")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	nestedReadme := "# Nested\n\nNot a taxonomy README — must stay byte-identical.\n"
	if err := os.WriteFile(filepath.Join(nested, "README.md"), []byte(nestedReadme), 0o644); err != nil {
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
	isTaxonomyReadme := func(p string) bool {
		parts := strings.Split(p, "/")
		return p == "docs/README.md" || (len(parts) == 3 && parts[0] == "docs" && parts[2] == "README.md")
	}
	for p, b := range before {
		if after[p] != b && !isTaxonomyReadme(p) {
			t.Fatalf("file outside the taxonomy READMEs was modified: %s", p)
		}
	}
	if after["docs/goals/nested/README.md"] != nestedReadme {
		t.Fatal("nested README was modified")
	}
}

// AC-027 negative: a regular file named docs is not a docs tree — Regen
// errors instead of succeeding with nothing regenerated.
func TestAC027_DocsAsRegularFileIsAnError(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "docs"), []byte("not a tree"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Regen(root); err == nil {
		t.Fatal("expected an error when docs is a regular file, got none")
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

// AC-138: an appended index row states the record it links — id, title, and
// status for an ordinary artifact, enforcement for a constraint — rather than
// restating the filename the link already carries. The title here is
// YAML-quoted because its value contains a colon; the row must spell the
// parsed value, never the quoting.
func TestAC138_UnitPositive_AppendedRowsUseStatusOrConstraintEnforcement(t *testing.T) {
	root, _ := runInto(t)
	readme := filepath.Join(root, "docs", "goals", "README.md")
	prose := "# Goals\n\n<!-- clue:index:start -->\n<!-- clue:index:end -->\n"
	if err := os.WriteFile(readme, []byte(prose), 0o644); err != nil {
		t.Fatal(err)
	}
	artifact := "---\nid: G-001\ntype: goal\nstatus: proposed\nlinks: []\ntitle: \"First goal: it carries a colon\"\n---\n"
	if err := os.WriteFile(filepath.Join(root, "docs", "goals", "G-001-first.md"), []byte(artifact), 0o644); err != nil {
		t.Fatal(err)
	}
	constraint := "---\nid: C-999\ntype: constraint\nstatus: active\nlinks: [C-001]\ntitle: First constraint\nsource: C-001\nenforcement: machine\n---\n"
	if err := os.WriteFile(filepath.Join(root, "docs", "constraints", "C-999-first.md"), []byte(constraint), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Regen(root); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(readme)
	if err != nil {
		t.Fatal(err)
	}
	want := "- [G-001 — First goal: it carries a colon](G-001-first.md) · `proposed`"
	if !strings.Contains(string(got), want) {
		t.Fatalf("appended row does not state the record:\nwant %q\ngot:\n%s", want, got)
	}
	constraints, err := os.ReadFile(filepath.Join(root, "docs", "constraints", "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	constraintRow := "- [C-999 — First constraint](C-999-first.md) · `machine`"
	if !strings.Contains(string(constraints), constraintRow) {
		t.Fatalf("appended constraint row does not state enforcement:\nwant %q\ngot:\n%s", constraintRow, constraints)
	}
	if issues := activatedValidateAt(t, root); len(issues) > 0 {
		t.Fatalf("expected green after regen, got: %v", issues)
	}
}

// AC-096: an appended row says what the artifact is about. The lede beneath the
// H1 wins over prose under a later heading, and every inline link in the
// sentence is reduced to its label — including one inside a code span, whose
// placeholder target resolves to nothing and would leave the generator emitting
// a corpus the judge rejects.
func TestAC096_UnitPositive_AppendedRowSaysWhatTheArtifactIsAbout(t *testing.T) {
	root, _ := runInto(t)
	readme := filepath.Join(root, "docs", "goals", "README.md")
	prose := "# Goals\n\n<!-- clue:index:start -->\n<!-- clue:index:end -->\n"
	if err := os.WriteFile(readme, []byte(prose), 0o644); err != nil {
		t.Fatal(err)
	}
	// A lede directly beneath the H1, then a heading whose paragraph must lose.
	lede := "---\nid: G-001\ntype: goal\nstatus: proposed\nlinks: []\ntitle: First goal\n---\n\n" +
		"# G-001 — First goal\n\nA reader reaches the [thread](../other.md) in one hop, as `- [a](b)` shows. Second sentence.\n\n" +
		"## Context\n\nProse under a heading that must not be chosen.\n"
	if err := os.WriteFile(filepath.Join(root, "docs", "goals", "G-001-first.md"), []byte(lede), 0o644); err != nil {
		t.Fatal(err)
	}
	// No lede: the first paragraph under the first heading is read, and the
	// table, list, quote and fence above it are structure rather than prose.
	fallback := "---\nid: G-002\ntype: goal\nstatus: proposed\nlinks: []\ntitle: Second goal\n---\n\n" +
		"# G-002 — Second goal\n\n## Context\n\n| a | b |\n|---|---|\n\n- a list item\n\n> a quote\n\n```\nfenced. text\n```\n\nThe paragraph that should win.\n"
	if err := os.WriteFile(filepath.Join(root, "docs", "goals", "G-002-second.md"), []byte(fallback), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Regen(root); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(readme)
	if err != nil {
		t.Fatal(err)
	}
	got := string(raw)
	wantLede := "- [G-001 — First goal](G-001-first.md) · `proposed` — A reader reaches the thread in one hop, as `- a` shows."
	if !strings.Contains(got, wantLede) {
		t.Fatalf("lede must be seeded with every link reduced to its label:\nwant %q\ngot:\n%s", wantLede, got)
	}
	if strings.Contains(got, "](b)") {
		t.Fatalf("a link inside a code span must not survive into the row, got:\n%s", got)
	}
	wantFallback := "- [G-002 — Second goal](G-002-second.md) · `proposed` — The paragraph that should win."
	if !strings.Contains(got, wantFallback) {
		t.Fatalf("a body with no lede reads the first paragraph under its first heading:\nwant %q\ngot:\n%s", wantFallback, got)
	}
	if strings.Contains(got, "Second sentence") || strings.Contains(got, "must not be chosen") {
		t.Fatalf("only the first sentence of the chosen paragraph is seeded, got:\n%s", got)
	}
	if issues := activatedValidateAt(t, root); len(issues) > 0 {
		t.Fatalf("expected green after regen, got: %v", issues)
	}
}

// AC-096 negative: a sentence past the bound is cut at a word boundary, so a
// long opening paragraph cannot push the identity a reader scans for off the
// line, and no word is severed mid-way.
func TestAC096_UnitNegative_LongSentenceIsCutAtAWordBoundary(t *testing.T) {
	root, _ := runInto(t)
	readme := filepath.Join(root, "docs", "goals", "README.md")
	if err := os.WriteFile(readme, []byte("# Goals\n\n<!-- clue:index:start -->\n<!-- clue:index:end -->\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	long := strings.TrimSpace(strings.Repeat("elaboration ", 40)) + " end"
	body := "---\nid: G-001\ntype: goal\nstatus: proposed\nlinks: []\ntitle: First goal\n---\n\n# G-001\n\n" + long + ".\n"
	if err := os.WriteFile(filepath.Join(root, "docs", "goals", "G-001-first.md"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Regen(root); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(readme)
	if err != nil {
		t.Fatal(err)
	}
	row := ""
	for _, line := range strings.Split(string(raw), "\n") {
		if strings.Contains(line, "G-001-first.md") {
			row = line
		}
	}
	if !strings.HasSuffix(row, "…") {
		t.Fatalf("an over-long sentence must be cut and marked, got %q", row)
	}
	if strings.Contains(row, "end") {
		t.Fatalf("the cut must drop the tail of the sentence, got %q", row)
	}
	if strings.Contains(row, "elaborat…") || strings.Contains(row, "elaboratio…") {
		t.Fatalf("the cut must fall on a word boundary, got %q", row)
	}
	if issues := activatedValidateAt(t, root); len(issues) > 0 {
		t.Fatalf("expected green after regen, got: %v", issues)
	}
}

// AC-097 negative: an artifact whose frontmatter cannot be read degrades all
// the way to the plain link rather than acquiring a status badge or a
// description tail. A row is one shape or the other, never a third.
func TestAC097_UnitNegative_UnreadableFrontmatterGetsNeitherBadgeNorTail(t *testing.T) {
	root, _ := runInto(t)
	readme := filepath.Join(root, "docs", "goals", "README.md")
	if err := os.WriteFile(readme, []byte("# Goals\n\n<!-- clue:index:start -->\n<!-- clue:index:end -->\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// No frontmatter at all, but a body a description could have been read from.
	body := "# G-009 — Broken goal\n\nA sentence a seed would happily have used.\n"
	if err := os.WriteFile(filepath.Join(root, "docs", "goals", "G-009-broken.md"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Regen(root); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(readme)
	if err != nil {
		t.Fatal(err)
	}
	if want := "- [G-009-broken](G-009-broken.md)\n"; !strings.Contains(string(raw), want) {
		t.Fatalf("unreadable frontmatter must fall back to the plain link, want %q, got:\n%s", want, raw)
	}
	if strings.Contains(string(raw), "A sentence a seed would happily have used") {
		t.Fatalf("a plain-link row carries no description, got:\n%s", raw)
	}
	if strings.Contains(string(raw), "G-009-broken.md) ·") {
		t.Fatalf("a plain-link row carries no status badge, got:\n%s", raw)
	}
}

// AC-097: an artifact with no readable prose sentence keeps the row that states
// its record, with no trailing separator and no empty tail — a row is one shape
// or the other, never a third.
func TestAC097_UnitPositive_NoReadableSentenceLeavesNoEmptyTail(t *testing.T) {
	root, _ := runInto(t)
	readme := filepath.Join(root, "docs", "goals", "README.md")
	if err := os.WriteFile(readme, []byte("# Goals\n\n<!-- clue:index:start -->\n<!-- clue:index:end -->\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// Frontmatter and a heading, then nothing a sentence can be read from.
	body := "---\nid: G-001\ntype: goal\nstatus: proposed\nlinks: []\ntitle: First goal\n---\n\n# G-001 — First goal\n\n## Context\n\n- only a list item\n"
	if err := os.WriteFile(filepath.Join(root, "docs", "goals", "G-001-first.md"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Regen(root); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(readme)
	if err != nil {
		t.Fatal(err)
	}
	want := "- [G-001 — First goal](G-001-first.md) · `proposed`\n"
	if !strings.Contains(string(raw), want) {
		t.Fatalf("a row with no readable sentence must be exactly the stated record:\nwant %q\ngot:\n%s", want, raw)
	}
	if strings.Contains(string(raw), "`proposed` —") {
		t.Fatalf("a row must never carry an empty description tail, got:\n%s", raw)
	}
	if issues := activatedValidateAt(t, root); len(issues) > 0 {
		t.Fatalf("expected green after regen, got: %v", issues)
	}
}

// AC-160: the seed is a first draft and never an assertion — a description an
// author corrected by hand outlives regeneration, and nothing backfills a row
// that already exists.
func TestAC160_UnitPositive_CuratedDescriptionOutlivesRegeneration(t *testing.T) {
	root, _ := runInto(t)
	readme := filepath.Join(root, "docs", "goals", "README.md")
	curated := "- [G-001 — First goal](G-001-first.md) · `proposed` — What an author decided this actually means."
	if err := os.WriteFile(readme, []byte("# Goals\n\n<!-- clue:index:start -->\n"+curated+"\n<!-- clue:index:end -->\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	body := "---\nid: G-001\ntype: goal\nstatus: proposed\nlinks: []\ntitle: First goal\n---\n\n# G-001\n\nThe body sentence a seed would have used.\n"
	if err := os.WriteFile(filepath.Join(root, "docs", "goals", "G-001-first.md"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	rep, err := Regen(root)
	if err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(readme)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), curated) {
		t.Fatalf("a curated row must survive regeneration untouched, got:\n%s", raw)
	}
	if strings.Contains(string(raw), "The body sentence a seed would have used") {
		t.Fatalf("regeneration must not backfill over a curated description, got:\n%s", raw)
	}
	for _, rel := range rep.Indexed {
		if rel == "docs/goals/README.md" {
			t.Fatalf("a corpus whose rows are already curated must report nothing regenerated, got %v", rep.Indexed)
		}
	}
}

// AC-160 negative: preservation is not unconditional. A curated description
// buys the row nothing once its target is gone — the row is dropped, exactly as
// an uncurated one would be, so a deleted artifact cannot leave a dangling
// entry behind on the strength of having been described.
func TestAC160_UnitNegative_CuratedRowForAMissingTargetIsStillDropped(t *testing.T) {
	root, _ := runInto(t)
	readme := filepath.Join(root, "docs", "goals", "README.md")
	curated := "- [G-404 — Deleted goal](G-404-deleted.md) · `accepted` — A description an author wrote before the file was removed."
	if err := os.WriteFile(readme, []byte("# Goals\n\n<!-- clue:index:start -->\n"+curated+"\n<!-- clue:index:end -->\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// G-404-deleted.md is deliberately never written.
	if _, err := Regen(root); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(readme)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(raw), "G-404-deleted.md") {
		t.Fatalf("a row whose target is gone must be dropped however well curated, got:\n%s", raw)
	}
	if strings.Contains(string(raw), "before the file was removed") {
		t.Fatalf("the curated description must go with its row, got:\n%s", raw)
	}
	if issues := activatedValidateAt(t, root); len(issues) > 0 {
		t.Fatalf("expected green after regen, got: %v", issues)
	}
}

// AC-138 negative: a target with no readable identity or no required
// constraint enforcement falls back to the plain link instead of emitting a
// half-formed row, and a subfolder row states a section rather than a record.
func TestAC138_UnitNegative_UnreadableIdentityAndSubfolderRowsStayPlain(t *testing.T) {
	root, _ := runInto(t)
	readme := filepath.Join(root, "docs", "goals", "README.md")
	prose := "# Goals\n\n<!-- clue:index:start -->\n<!-- clue:index:end -->\n"
	if err := os.WriteFile(readme, []byte(prose), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "goals", "G-002-broken.md"), []byte("# no frontmatter at all\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// Frontmatter carrying an id and title but no status must degrade the same
	// way: a row is the stated shape or the plain link, never a third form
	// carrying an empty status badge.
	noStatus := "---\nid: G-004\ntype: goal\nlinks: []\ntitle: Statusless goal\n---\n"
	if err := os.WriteFile(filepath.Join(root, "docs", "goals", "G-004-statusless.md"), []byte(noStatus), 0o644); err != nil {
		t.Fatal(err)
	}
	noEnforcement := "---\nid: C-998\ntype: constraint\nstatus: active\nlinks: []\ntitle: Enforcementless constraint\nsource: C-001\n---\n"
	if err := os.WriteFile(filepath.Join(root, "docs", "constraints", "C-998-enforcementless.md"), []byte(noEnforcement), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "docs", "extra"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "extra", "README.md"), []byte("# Extra\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Regen(root); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(readme)
	if err != nil {
		t.Fatal(err)
	}
	if want := "- [G-002-broken](G-002-broken.md)\n"; !strings.Contains(string(got), want) {
		t.Fatalf("unreadable identity must fall back to the plain link, want %q, got:\n%s", want, got)
	}
	if strings.Contains(string(got), "G-002-broken.md) ·") {
		t.Fatalf("a row with no readable identity must carry no status:\n%s", got)
	}
	if want := "- [G-004-statusless](G-004-statusless.md)\n"; !strings.Contains(string(got), want) {
		t.Fatalf("an artifact without a status must fall back to the plain link, want %q, got:\n%s", want, got)
	}
	if strings.Contains(string(got), "· ``") {
		t.Fatalf("no row may carry an empty status badge:\n%s", got)
	}
	constraintsReadme, err := os.ReadFile(filepath.Join(root, "docs", "constraints", "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	if want := "- [C-998-enforcementless](C-998-enforcementless.md)\n"; !strings.Contains(string(constraintsReadme), want) {
		t.Fatalf("a constraint without enforcement must fall back to the plain link, want %q, got:\n%s", want, constraintsReadme)
	}
	if strings.Contains(string(constraintsReadme), "C-998-enforcementless.md) ·") {
		t.Fatalf("a constraint without enforcement must carry no badge:\n%s", constraintsReadme)
	}
	rootReadme, err := os.ReadFile(filepath.Join(root, "docs", "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	if want := "- [extra/](extra/README.md)"; !strings.Contains(string(rootReadme), want) {
		t.Fatalf("subfolder row must stay plain, want %q, got:\n%s", want, rootReadme)
	}
}
