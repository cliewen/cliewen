package scaffold

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
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

func TestAC145_UnitPositive_InitAndScaffoldUseSubjectTypedDecisions(t *testing.T) {
	root, _ := runInto(t)
	readme := filepath.Join(root, "docs", "decisions", "README.md")
	content, err := os.ReadFile(readme)
	if err != nil {
		t.Fatal(err)
	}
	text := string(content)
	for _, want := range []string{"ADR-xxx", "PDR-xxx", "IDR-xxx", "Subject alone selects the type"} {
		if !strings.Contains(text, want) {
			t.Errorf("scaffolded decisions guidance is missing %q", want)
		}
	}
	if _, err := os.Stat(filepath.Join(root, "docs", "decisions", "log.md")); !os.IsNotExist(err) {
		t.Fatalf("init materialized a legacy decision log: %v", err)
	}
	before := text
	if _, err := Run(root); err != nil {
		t.Fatal(err)
	}
	after, err := os.ReadFile(readme)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != before {
		t.Fatalf("scaffold changed an already-current empty decision index:\n%s", after)
	}
	if issues := validateAt(t, root); len(issues) != 0 {
		t.Fatalf("subject-typed scaffold does not validate: %v", issues)
	}
}

func TestAC145_UnitNegative_ScaffoldDoesNotOverwriteExistingDecision(t *testing.T) {
	root, _ := runInto(t)
	rel := filepath.Join("docs", "decisions", "IDR-001-local-choice.md")
	decision := "---\nid: IDR-001\ntype: decision\nstatus: inferred\nlinks: []\ntitle: Local choice\nauthor: human\naccepted-by: []\n---\n\n# IDR-001 — Local choice\n"
	full := filepath.Join(root, rel)
	if err := os.WriteFile(full, []byte(decision), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Run(root); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(full)
	if err != nil || string(got) != decision {
		t.Fatalf("scaffold overwrote an existing decision: err=%v\n%s", err, got)
	}
	if _, err := os.Stat(filepath.Join(root, "docs", "decisions", "log.md")); !os.IsNotExist(err) {
		t.Fatalf("scaffold invented a legacy log while indexing an existing decision: %v", err)
	}
}

// AC-042: every initialized repository receives the PR form required by the
// scaffolded CI wall, rather than leaving authors to recreate it by hand.
func TestAC042_InitOutputIncludesAcceptanceBriefTemplate(t *testing.T) {
	root, _ := runInto(t)
	want, err := os.ReadFile(filepath.Join("..", "..", ".github", "pull_request_template.md"))
	if err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filepath.Join(root, ".github", "pull_request_template.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(want) {
		t.Error("scaffolded pull-request template differs from the repository template")
	}
}

func TestSanity_ScaffoldedEvidenceModelCarriersAgree(t *testing.T) {
	root, _ := runInto(t)
	requiredByFile := map[string][]string{
		"AGENTS.md":                        {"criterion → acceptance evidence", "Human proof in the acceptance brief"},
		"docs/README.md":                   {"acceptance evidence", "classified Go/JVM/Cucumber test reference", "Human acceptance brief"},
		"docs/capabilities/README.md":      {"positive/negative direction", "Test-type: Human", "@draft", "one supported reference", "does not execute tests"},
		"docs/decisions/README.md":         {"inventories every live carrier", "same change", "general obligation remains agent-enforced"},
		".github/pull_request_template.md": {"classified positive/negative Go, JVM, or Cucumber evidence", "JVM identity/type/direction attached to the same executable", "Human proof named in the acceptance brief", "per-criterion `@draft`", "legacy one-supported-reference rule"},
	}
	for rel, required := range requiredByFile {
		content, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
		if err != nil {
			t.Fatal(err)
		}
		for _, want := range required {
			if !strings.Contains(string(content), want) {
				t.Errorf("scaffolded %s is missing evidence-model carrier %q", rel, want)
			}
		}
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

// AC-038 negative: an ordinary .claude/skills directory is mirrored into
// exactly as before and produces no linked report — the skip is caused by
// the link, not by the path.
func TestAC038_OrdinaryMirrorDirectoryIsStillWritten(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".claude", "skills"), 0o755); err != nil {
		t.Fatal(err)
	}
	rep, err := Run(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(rep.Linked) != 0 {
		t.Fatalf("a real directory is not a link, got %v", rep.Linked)
	}
	if _, err := os.Stat(filepath.Join(root, ".claude", "skills", "clue-delta", "SKILL.md")); err != nil {
		t.Fatalf("the mirror was not written into a real directory: %v", err)
	}
}

// The two generated distribution trees must stay byte-identical. The
// skill-generator package separately holds both trees to their shared
// canonical render.
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
			t.Errorf("embedded %s differs from .agents/skills — run go generate ./internal/skills", rel)
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

func TestAC139_UnitPositive_ScaffoldedRoutingRecommendsByAcceptedContractBeforeCorpus(t *testing.T) {
	root, _ := runInto(t)
	data, err := os.ReadFile(filepath.Join(root, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	classify := strings.Index(content, "Recommended route: simple")
	readCorpus := strings.Index(content, "For a full change, read [`docs/README.md`]")
	if classify < 0 || readCorpus < 0 || classify >= readCorpus {
		t.Fatalf("AGENTS.md does not recommend a route before loading full-change corpus context:\n%s", content)
	}
	for _, want := range []string{
		"accepted contract unchanged",
		"defect correction restoring an unchanged criterion",
		"Paths and diff size may warn but never decide meaning",
		"Cliewen-Route: simple",
		"explicit user authorization and repository permission",
		"Release is not a Cliewen route",
	} {
		if !strings.Contains(content, want) {
			t.Errorf("AGENTS.md does not contain adaptive-routing rule %q", want)
		}
	}
}

func TestAC139_UnitNegative_ScaffoldedRoutingDoesNotExportLegacyTiersOrReleasePolicy(t *testing.T) {
	root, _ := runInto(t)
	data, err := os.ReadFile(filepath.Join(root, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	for _, forbidden := range []string{"**Plain", "**Light", "release route", "version cut"} {
		if strings.Contains(content, forbidden) {
			t.Errorf("scaffolded AGENTS.md exports retired or repository-local routing %q", forbidden)
		}
	}
}

func TestSanity_ScaffoldedRoutingRoutesDetailedHandoffToCanonicalSkill(t *testing.T) {
	root, _ := runInto(t)
	hub, err := os.ReadFile(filepath.Join(root, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(hub), "Use [`clue-delta`](.agents/skills/clue-delta/skill.md)'s full loop") {
		t.Fatal("AGENTS.md does not route a chosen full recommendation to clue-delta")
	}

	content := readSkillDirectory(t, filepath.Join("templates", "skills", "clue-delta"))
	for _, want := range []string{
		"Before marking a change ready",
		"Every review of an existing hosted PR is bound to its observed head SHA",
		"publish the finding there and leave it unresolved until a hosted commit contains the reviewed repair",
		"Any agent that edits an existing PR becomes the updater for that turn",
		"record its hosted head",
		"commit every intended edit",
		"clean agentic review pass",
		"git status --porcelain` to be empty",
		"head branch and SHA equal the current local branch and `HEAD`",
		"then mark the PR ready and perform the hosted check again immediately after",
		"If the head changed underneath the turn or a push is rejected as non-fast-forward",
		"If the PR merged or closed, stop without pushing",
		"Push is durability, never a signal",
		"The PR exists from first publication and starts as a draft",
		"Marking the PR ready for review is the explicit act that claims a candidate",
		"Stopping anywhere else is ordinary, not an exception",
	} {
		if !strings.Contains(content, want) {
			t.Errorf("clue-delta skill directory does not contain review-handoff rule %q", want)
		}
	}
	commitCandidate := strings.Index(content, "commit every intended edit")
	reviewCandidate := strings.Index(content, "clean agentic review pass")
	readyCandidate := strings.Index(content, "then mark the PR ready and perform the hosted check again immediately after")
	if commitCandidate < 0 || reviewCandidate <= commitCandidate || readyCandidate <= reviewCandidate {
		t.Error("clue-delta skill directory must commit and verify the candidate before agentic review, then mark the PR ready only after that binding")
	}
}

func TestSanity_ScaffoldedCanonicalSkillCarriesMergeHistoryBoundary(t *testing.T) {
	content := readSkillDirectory(t, filepath.Join("templates", "skills", "clue-delta"))
	for _, want := range []string{
		"human accepts the ready pull request with a merge commit",
		"disable squash and rebase-and-merge",
		"supported full-change adoption path",
	} {
		if !strings.Contains(content, want) {
			t.Errorf("scaffolded clue-delta skill is missing merge-history contract %q", want)
		}
	}
}

func readSkillDirectory(t *testing.T, directory string) string {
	t.Helper()
	var content strings.Builder
	err := filepath.WalkDir(directory, func(filePath string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil || entry.IsDir() {
			return walkErr
		}
		data, readErr := os.ReadFile(filePath)
		if readErr != nil {
			return readErr
		}
		content.Write(data)
		content.WriteByte('\n')
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return content.String()
}

// The CI caller's version pin comes from the embedded skills' stamp, and its
// reusable-workflow reference resolves to one of the two immutable forms
// ADR-038 admits — the emitting source commit, or the release tag a build
// without usable VCS metadata falls back to. Which one depends on how the
// test binary was built, so the assertion is on immutability, not on the
// commit form; requiring a SHA here would fail every build made outside a
// clean checkout of this repository. No placeholder may survive either way.
func TestUnit_WorkflowVersionSubstituted(t *testing.T) {
	root, _ := runInto(t)
	data, err := os.ReadFile(filepath.Join(root, ".github", "workflows", "clue.yml"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), versionPlaceholder) || strings.Contains(string(data), workflowRefPlaceholder) {
		t.Fatal("workflow still contains a scaffold placeholder")
	}
	version, err := PairVersion()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "clue-version: "+version) {
		t.Fatalf("workflow does not pin clue-version to %s", version)
	}
	if !regexp.MustCompile(`clue-validation\.yml@(?:[0-9a-f]{40}|v[0-9]+\.[0-9]+\.[0-9]+)\b`).Match(data) {
		t.Fatal("workflow does not pin the reusable validation unit to a commit or release tag")
	}
}

// extractGrepPattern pulls the single-quoted regex argument of a shipped
// `grep -E` invocation out of the emitted workflow and compiles it, so the
// behavioral test below judges the patterns actually shipped, not a copy.
func extractGrepPattern(t *testing.T, workflow, invocation string) *regexp.Regexp {
	t.Helper()
	m := regexp.MustCompile(invocation).FindStringSubmatch(workflow)
	if m == nil {
		t.Fatalf("clue.yml no longer contains a %q classifier", invocation)
	}
	return regexp.MustCompile(m[1])
}

// The scaffolded wall's plain/Cliewen split ships as inline grep patterns
// with no executable test of its own — only substring assertions. This runs
// the actual patterns emitted into clue.yml through the classifier's decision
// procedure, so a regex edit that reclassified a protected surface as plain
// (or a non-Markdown file as editorial) fails here instead of in an adopter's
// merged history. The empty and diff-failure branches are pure shell and stay
// covered by their fail-closed encoding here.
func TestUnit_ScaffoldWallClassifierClassifiesByShippedPatterns(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", ".github", "workflows", "clue-validation.yml"))
	if err != nil {
		t.Fatal(err)
	}
	workflow := string(data)

	notMarkdown := extractGrepPattern(t, workflow, `grep -Eqv '([^']*)'`)
	protected := extractGrepPattern(t, workflow, `grep -Eq '([^']*)'`)

	// Mirrors the emitted shell for a successful, non-empty diff: any file
	// that is not Markdown, any Markdown file on a protected path, and an
	// empty list are Cliewen changes; only Markdown outside protected paths
	// is plain.
	classify := func(files []string) string {
		if len(files) == 0 {
			return "cliewen" // shell: [ ! -s ] fails closed
		}
		for _, f := range files {
			if !notMarkdown.MatchString(f) { // grep -Eqv '\.md$': a non-.md file
				return "cliewen"
			}
		}
		for _, f := range files {
			if protected.MatchString(f) { // grep -Eq '^(docs/|...)'
				return "cliewen"
			}
		}
		return "plain"
	}

	cases := []struct {
		name  string
		files []string
		want  string
	}{
		{"guide markdown only is plain", []string{"guide/intro.md", "guide/nested/deep.md"}, "plain"},
		{"corpus markdown is protected", []string{"docs/decisions/PDR-011.md"}, "cliewen"},
		{"transient workspace is protected", []string{"changes/CH-099-x/proposal.md"}, "cliewen"},
		{"routing hub is protected", []string{"AGENTS.md"}, "cliewen"},
		{"changelog is protected", []string{"CHANGELOG.md"}, "cliewen"},
		{"skill source is protected", []string{".agents/skills/clue-delta/skill.md"}, "cliewen"},
		{"code is not editorial", []string{"cmd/clue/main.go"}, "cliewen"},
		{"config is not editorial", []string{"guide/.vitepress/config.mts"}, "cliewen"},
		{"mixed guide and code fails closed", []string{"guide/intro.md", "cmd/clue/main.go"}, "cliewen"},
		{"empty diff fails closed", nil, "cliewen"},
		// Documents the deliberately permissive path rule: unprotected
		// root Markdown passes as plain, and the human review boundary is
		// the backstop (PDR-011). Change this only with that decision.
		{"unprotected root markdown is plain", []string{"README.md"}, "plain"},
	}
	for _, tc := range cases {
		if got := classify(tc.files); got != tc.want {
			t.Errorf("%s: classify(%v) = %q, want %q", tc.name, tc.files, got, tc.want)
		}
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
	if !strings.Contains(text, "](G-001-first.md) · `proposed`\r\n") {
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
