package migrate

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/cliewen/cliewen/internal/corpus"
	"github.com/cliewen/cliewen/internal/scaffold"
)

func TestAC064_UnitPositive_MigrationIsPreviewableIdempotentAndCoordinated(t *testing.T) {
	root := migrationFixture(t, "")

	workflowPath := filepath.Join(root, ".github", "workflows", "clue.yml")
	workflow, err := os.ReadFile(workflowPath)
	if err != nil {
		t.Fatal(err)
	}
	wantCarriers, err := scaffold.ManagedCarrierFiles()
	if err != nil {
		t.Fatal(err)
	}
	wantWorkflow := string(wantCarriers[".github/workflows/clue.yml"])
	wantUses := callerUsesRe.FindStringSubmatch(wantWorkflow)
	if len(wantUses) == 0 {
		t.Fatal("embedded caller has no upstream reference")
	}
	version, err := scaffold.PairVersion()
	if err != nil {
		t.Fatal(err)
	}
	workflowText := string(workflow)
	workflowText = strings.Replace(workflowText, "@"+wantUses[4], "@v0.9.0", 1)
	workflowText = strings.Replace(workflowText, "clue-version: "+version, "clue-version: 0.9.0", 1)
	workflowText = strings.Replace(workflowText, `runner: '["ubuntu-latest"]'`, `runner: '["self-hosted"]'`, 1)
	if err := os.WriteFile(workflowPath, []byte(workflowText), 0o644); err != nil {
		t.Fatal(err)
	}
	oldSkill, err := os.ReadFile(filepath.Join(root, ".agents", "skills", "clue-analysis", "skill.md"))
	if err != nil {
		t.Fatal(err)
	}
	mirrorPath := filepath.Join(root, ".claude", "skills", "clue-analysis", "SKILL.md")
	if err := os.WriteFile(mirrorPath, oldSkill, 0o644); err != nil {
		t.Fatal(err)
	}

	preview, err := Plan(root, Options{ReversalCost: "low"})
	if err != nil {
		t.Fatal(err)
	}
	if len(preview.Findings) != 0 {
		t.Fatalf("preview findings: %+v", preview.Findings)
	}
	if len(preview.Changes) < 2 {
		t.Fatalf("preview planned %d changes, want corpus and caller changes", len(preview.Changes))
	}
	var mirrorPlanned bool
	for _, change := range preview.Changes {
		if change.Path == ".claude/skills/clue-analysis/SKILL.md" {
			mirrorPlanned = true
		}
	}
	if !mirrorPlanned {
		t.Fatal("recognized old Claude mirror was not included in the migration plan")
	}
	docBefore, err := os.ReadFile(filepath.Join(root, "docs", "analysis", "AN-001.md"))
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(docBefore, []byte("reversal-cost:")) || !bytes.Contains(docBefore, []byte("status: verified")) {
		t.Fatal("preview changed the corpus before apply")
	}
	if err := Apply(root, preview); err != nil {
		t.Fatal(err)
	}

	docAfter, err := os.ReadFile(filepath.Join(root, "docs", "analysis", "AN-001.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(docAfter, []byte("status: active")) || !bytes.Contains(docAfter, []byte("reversal-cost: low")) || !bytes.Contains(docAfter, []byte("The body remains byte-for-byte user content.")) {
		t.Fatalf("migration did not apply the expected corpus edits:\n%s", docAfter)
	}
	workflowAfter, err := os.ReadFile(workflowPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(workflowAfter, []byte(`runner: '["self-hosted"]'`)) || !bytes.Contains(workflowAfter, []byte("clue-version: "+version)) || !bytes.Contains(workflowAfter, []byte("@"+wantUses[4])) {
		t.Fatalf("migration did not preserve caller choices and update its owned fields:\n%s", workflowAfter)
	}
	mirrorAfter, err := os.ReadFile(mirrorPath)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(mirrorAfter, oldSkill) || !bytes.Contains(mirrorAfter, []byte("version: "+version)) {
		t.Fatal("migration did not refresh the recognized Claude mirror")
	}

	repeat, err := Plan(root, Options{ReversalCost: "low"})
	if err != nil {
		t.Fatal(err)
	}
	if len(repeat.Changes) != 0 || len(repeat.Findings) != 0 {
		t.Fatalf("migration is not idempotent: changes=%+v findings=%+v", repeat.Changes, repeat.Findings)
	}

	c, issues := corpus.Scan(root)
	if len(issues) != 0 {
		t.Fatalf("fixture scan issues after migration: %v", issues)
	}
	if issues := corpus.Validate(c, corpus.Options{}); len(issues) != 0 {
		t.Fatalf("fixture is not green after migration: %v", issues)
	}
}

// TestUnit_CarrierPreviewNamesTheReleaseItWrites pins the direction of the
// MIG-003 preview line.
//
// The preview used to read "replace generated carrier clue-plan/skill.md with
// release 0.11.2 content" while writing the target's bytes — it named the
// release it recognized, which is the one release those bytes are certainly
// not. An adopter upgrading from 0.11.2 would read that as a downgrade and be
// right to refuse it. The recognized release is still worth saying, because it
// is why the replacement is safe; it just is not the source of the content.
func TestUnit_CarrierPreviewNamesTheReleaseItWrites(t *testing.T) {
	root := migrationFixture(t, "")
	repoRoot := repositoryRoot(t)

	const legacy = "v0.11.2"
	if !tagExists(t, repoRoot, legacy) {
		t.Skipf("%s not present in this checkout", legacy)
	}
	old, err := gitShow(t, repoRoot, legacy+":.agents/skills/clue-analysis/skill.md")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".agents", "skills", "clue-analysis", "skill.md"), old, 0o644); err != nil {
		t.Fatal(err)
	}
	target, err := scaffold.PairVersion()
	if err != nil {
		t.Fatal(err)
	}

	preview, err := Plan(root, Options{ReversalCost: "low"})
	if err != nil {
		t.Fatal(err)
	}
	var got string
	for _, change := range preview.Changes {
		if change.Path == ".agents/skills/clue-analysis/skill.md" {
			got = change.Description
		}
	}
	if got == "" {
		t.Fatalf("recognized 0.11.2 carrier was not planned for replacement: %+v", preview.Changes)
	}
	if !strings.Contains(got, "with release "+target+" content") {
		t.Errorf("preview does not name the release it writes:\n  got:  %s\n  want it to contain: with release %s content", got, target)
	}
	if !strings.Contains(got, "recognized as release 0.11.2") {
		t.Errorf("preview dropped the release it recognized:\n  got: %s", got)
	}
	if strings.Contains(got, "with release 0.11.2 content") {
		t.Errorf("preview still names the recognized release as the source of the bytes:\n  got: %s", got)
	}
}

// TestUnit_AddedCarrierPreviewNamesTheReleaseItWrites covers the other half of
// the same defect, where it was worse.
//
// When a release introduces a carrier its predecessor never shipped, the plan
// adds it and used to describe it as "from release <old> content" — naming, as
// the source of the bytes, the one release that is known not to contain the
// file at all. clue-extract/mappings/madr.md first shipped in 0.7.0, so a
// 0.6.0 repository reaches exactly this branch.
func TestUnit_AddedCarrierPreviewNamesTheReleaseItWrites(t *testing.T) {
	root := migrationFixture(t, "")
	repoRoot := repositoryRoot(t)

	const legacy = "v0.6.0"
	if !tagExists(t, repoRoot, legacy) {
		t.Skipf("%s not present in this checkout", legacy)
	}
	old, err := gitShow(t, repoRoot, legacy+":.agents/skills/clue-extract/skill.md")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".agents", "skills", "clue-extract", "skill.md"), old, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(root, ".agents", "skills", "clue-extract", "mappings", "madr.md")); err != nil {
		t.Fatal(err)
	}
	target, err := scaffold.PairVersion()
	if err != nil {
		t.Fatal(err)
	}

	preview, err := Plan(root, Options{ReversalCost: "low"})
	if err != nil {
		t.Fatal(err)
	}
	var got string
	for _, change := range preview.Changes {
		if change.Path == ".agents/skills/clue-extract/mappings/madr.md" {
			got = change.Description
		}
	}
	if got == "" {
		t.Fatalf("carrier absent from 0.6.0 was not planned as an add: changes=%+v findings=%+v", preview.Changes, preview.Findings)
	}
	if !strings.Contains(got, "release "+target+" content") {
		t.Errorf("preview does not name the release it writes:\n  got:  %s\n  want it to contain: release %s content", got, target)
	}
	if strings.Contains(got, "release 0.6.0 content") {
		t.Errorf("preview names a release that never published this file as the source of its bytes:\n  got: %s", got)
	}
}

func TestAC064_UnitNegative_MigrationRefusesAmbiguousAndModifiedInputs(t *testing.T) {
	root := migrationFixture(t, `---
id: AN-001
type: experiment
status: verified
links: []
title: Ambiguous analysis
provenance: inferred
---
`)

	carrierPath := filepath.Join(root, ".agents", "skills", "clue-delta", "skill.md")
	carrier, err := os.ReadFile(carrierPath)
	if err != nil {
		t.Fatal(err)
	}
	carrier = append(carrier, []byte("\nlocal edit\n")...)
	if err := os.WriteFile(carrierPath, carrier, 0o644); err != nil {
		t.Fatal(err)
	}
	docPath := filepath.Join(root, "docs", "analysis", "AN-001.md")
	docBefore, err := os.ReadFile(docPath)
	if err != nil {
		t.Fatal(err)
	}

	plan, err := Plan(root, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Findings) < 2 {
		t.Fatalf("want missing semantic resolution and local-edit findings, got %+v", plan.Findings)
	}
	if err := Apply(root, plan); err == nil {
		t.Fatal("Apply accepted a plan with unresolved findings")
	}
	docAfter, err := os.ReadFile(docPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(docBefore, docAfter) || bytes.Contains(docAfter, []byte("reversal-cost:")) {
		t.Fatal("failed migration changed an ambiguous corpus file")
	}
	carrierAfter, err := os.ReadFile(carrierPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(carrier, carrierAfter) {
		t.Fatal("failed migration changed a locally modified carrier")
	}
}

func TestAC064_UnitNegative_MigrationRejectsChangedSourceAfterPreview(t *testing.T) {
	root := migrationFixture(t, "")
	plan, err := Plan(root, Options{ReversalCost: "low"})
	if err != nil {
		t.Fatal(err)
	}
	docPath := filepath.Join(root, "docs", "analysis", "AN-001.md")
	doc, err := os.ReadFile(docPath)
	if err != nil {
		t.Fatal(err)
	}
	doc = append(doc, []byte("local change after preview\n")...)
	if err := os.WriteFile(docPath, doc, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Apply(root, plan); err == nil || !strings.Contains(err.Error(), "changed after planning") {
		t.Fatalf("Apply did not reject a changed source: %v", err)
	}
	if got, err := os.ReadFile(docPath); err != nil || !bytes.Equal(got, doc) {
		t.Fatal("changed source was overwritten after preview")
	}
}

func TestAC064_UnitPositive_MigrationRegistryIsOrdered(t *testing.T) {
	registry := Registry()
	want := []string{MigrationReversalCost, MigrationStatusLifecycle, MigrationManagedCarriers, MigrationQualifiedReferences, MigrationClaudeEntryPoint, MigrationHubReleaseCheck, MigrationPromotedConstraints}
	if len(registry) != len(want) {
		t.Fatalf("registry has %d entries, want %d", len(registry), len(want))
	}
	for i, definition := range registry {
		if definition.ID != want[i] || definition.Description == "" {
			t.Fatalf("registry[%d] = %+v, want ID %s with a description", i, definition, want[i])
		}
	}
	registry[0].ID = "changed"
	if Registry()[0].ID != MigrationReversalCost {
		t.Fatal("Registry exposed mutable package state")
	}
}

func TestAC064_UnitNegative_MigrationRejectsUnsupportedSyntaxAndOptions(t *testing.T) {
	root := migrationFixture(t, "")
	plan, err := Plan(root, Options{ReversalCost: "guess"})
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Findings) != 1 || !strings.Contains(plan.Findings[0].Message, "low or high") {
		t.Fatalf("invalid routing choice was not rejected: %+v", plan.Findings)
	}

	if after, changes, findings := migrateArtifact("plain.md", []byte("plain text\n"), "low"); len(changes) != 0 || len(findings) != 0 || !bytes.Equal(after, []byte("plain text\n")) {
		t.Fatalf("plain markdown should be outside the migration: changes=%v findings=%v", changes, findings)
	}
	if _, _, findings := migrateArtifact("broken.md", []byte("---\nid: AN-001\n"), "low"); len(findings) != 1 {
		t.Fatalf("unclosed frontmatter was not reported: %+v", findings)
	}
	quotedStatus := []byte("---\nid: AN-001\ntype: analysis\nstatus: \"verified\"\nlinks: []\ntitle: A\n---\n")
	if after, _, findings := migrateArtifact("quoted-status.md", quotedStatus, "low"); len(findings) != 1 || !bytes.Equal(after, quotedStatus) {
		t.Fatalf("ambiguous status syntax was not refused: findings=%v after=%q", findings, after)
	}
	invalidCost := []byte("---\nid: AN-001\ntype: analysis\nstatus: active\nlinks: []\ntitle: A\nprovenance: inferred\nreversal-cost: cheap\n---\n")
	if after, _, findings := migrateArtifact("invalid-cost.md", invalidCost, "low"); len(findings) != 1 || !bytes.Equal(after, invalidCost) {
		t.Fatalf("invalid reversal-cost was not refused: findings=%v after=%q", findings, after)
	}
	decisionFields := []byte("---\nid: ADR-001\ntype: decision\nstatus: inferred\nlinks: []\ntitle: A\nprovenance: inferred\nreversal-cost: low\n---\n")
	if after, _, findings := migrateArtifact("decision-fields.md", decisionFields, "low"); len(findings) != 2 || !bytes.Equal(after, decisionFields) {
		t.Fatalf("decision-only provenance fields were not refused: findings=%v after=%q", findings, after)
	}
	crlf := []byte("---\r\nid: AN-001\r\ntype: analysis\r\nstatus: verified\r\nlinks: []\r\ntitle: A\r\nprovenance: inferred\r\n---\r\nbody\r\n")
	after, _, findings := migrateArtifact("crlf.md", crlf, "low")
	if len(findings) != 0 || !bytes.Contains(after, []byte("reversal-cost: low\r\n")) || !bytes.Contains(after, []byte("status: active\r\n")) {
		t.Fatalf("CRLF frontmatter was not migrated without a finding: findings=%v\n%s", findings, after)
	}
}

func TestAC064_UnitPositive_MigrationEditsOnlyTopLevelFields(t *testing.T) {
	// An indented key belongs to a nested mapping or to block scalar content.
	// It is not the artifact's field, and editing it would both rewrite user
	// content and leave the real field unmigrated.
	nested := []byte("---\nid: AN-001\ntype: analysis\nsource:\n  provenance: inferred\n  status: verified\n  note: keep\nstatus: verified\nlinks: []\ntitle: A\nprovenance: inferred\n---\n\nbody\n")
	want := []byte("---\nid: AN-001\ntype: analysis\nsource:\n  provenance: inferred\n  status: verified\n  note: keep\nstatus: active\nlinks: []\ntitle: A\nprovenance: inferred\nreversal-cost: low\n---\n\nbody\n")
	after, changes, findings := migrateArtifact("docs/analysis/AN-001.md", nested, "low")
	if len(findings) != 0 || len(changes) != 2 {
		t.Fatalf("nested keys blocked the real fields: changes=%v findings=%v", changes, findings)
	}
	if !bytes.Equal(after, want) {
		t.Fatalf("migration did not edit the top-level fields:\n got %s\nwant %s", after, want)
	}
}

func TestAC064_UnitPositive_MigrationCoversREADMEsThatAreArtifacts(t *testing.T) {
	// A folder README is prose, but a capability's README carries frontmatter
	// and the validator judges it like any other artifact. Skipping every file
	// named README.md let `clue migrate` report a clean plan while
	// `clue validate` still failed on the very field MIG-001 exists to add.
	root := migrationFixture(t, "")
	write := func(rel, content string) {
		t.Helper()
		full := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	capability := "docs/capabilities/CAP-009-example/README.md"
	write(capability, "---\nid: CAP-009\ntype: capability\nstatus: active\nlinks: []\ntitle: Example capability\ngoal: G-101\nprovenance: inferred\n---\n\nProse body.\n")
	// Prose README, and one whose leading `---` never closes: corpus.Scan
	// treats both as prose, so neither may become a change or a finding.
	write("docs/capabilities/README.md", "# Capabilities\n\nNo frontmatter here.\n")
	write("docs/capabilities/CAP-009-example/notes/README.md", "---\nnot a frontmatter block\n\nprose\n")

	plan, err := Plan(root, Options{ReversalCost: "high"})
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Findings) != 0 {
		t.Fatalf("prose READMEs produced findings: %+v", plan.Findings)
	}
	var planned *Change
	for i, change := range plan.Changes {
		if change.Path == capability {
			planned = &plan.Changes[i]
		}
		if strings.HasSuffix(change.Path, "capabilities/README.md") || strings.Contains(change.Path, "notes/README.md") {
			t.Fatalf("a prose README was migrated: %+v", change)
		}
	}
	if planned == nil {
		t.Fatalf("capability README with inferred provenance was not migrated: %+v", plan.Changes)
	}
	if !bytes.Contains(planned.After, []byte("reversal-cost: high")) || !bytes.Contains(planned.After, []byte("Prose body.")) {
		t.Fatalf("capability README migration is wrong:\n%s", planned.After)
	}
}

func TestAC064_UnitNegative_MigrationRefusesFieldsItCannotAnchor(t *testing.T) {
	// A key repeated at the top level has no single owner, so there is no safe
	// line to rewrite. Refuse rather than guess at the first or last one.
	duplicate := []byte("---\nid: AN-001\ntype: analysis\nstatus: verified\nstatus: verified\nlinks: []\ntitle: A\n---\n")
	if after, changes, findings := migrateArtifact("duplicate.md", duplicate, "low"); len(findings) == 0 || len(changes) != 0 || !bytes.Equal(after, duplicate) {
		t.Fatalf("a repeated top-level field was migrated: changes=%v findings=%v\n%s", changes, findings, after)
	}
}

func TestAC064_UnitPositive_MigrationPreservesEveryUntouchedLineEnding(t *testing.T) {
	// The migration rewrites the fields it owns, not the file's line endings.
	// A mixed-ending artifact must come back with each untouched line — prose
	// included — terminated exactly as it was.
	before := []byte("---\nid: AN-001\ntype: analysis\nstatus: verified\r\nlinks: []\ntitle: A\nprovenance: inferred\r\n---\n\nprose with LF\nprose with CRLF\r\nfinal LF\n")
	want := []byte("---\nid: AN-001\ntype: analysis\nstatus: active\r\nlinks: []\ntitle: A\nprovenance: inferred\r\nreversal-cost: low\r\n---\n\nprose with LF\nprose with CRLF\r\nfinal LF\n")

	after, changes, findings := migrateArtifact("docs/analysis/AN-001.md", before, "low")
	if len(findings) != 0 {
		t.Fatalf("mixed line endings were reported as a finding: %+v", findings)
	}
	if len(changes) != 2 {
		t.Fatalf("changes = %v, want the reversal-cost insert and the status replacement", changes)
	}
	if !bytes.Equal(after, want) {
		t.Fatalf("line endings were not preserved:\n got %q\nwant %q", after, want)
	}
}

func TestAC064_UnitPositive_LegacyReleasesAreReportedOldestFirst(t *testing.T) {
	// The manifest's own order names the release, so bytes shared by several
	// releases resolve to the earliest one regardless of how the version
	// strings would compare. A lexicographic sort would answer 0.10.0 here.
	shared := []byte("shared carrier bytes\n")
	digest := sha256.Sum256(shared)
	rel := "clue-test/skill.md"
	restore := legacyDigests
	t.Cleanup(func() { legacyDigests = restore })
	legacyDigests = []releaseManifest{
		{Version: "0.9.0", Files: map[string]string{rel: hex.EncodeToString(digest[:])}},
		{Version: "0.10.0", Files: map[string]string{rel: hex.EncodeToString(digest[:])}},
	}

	if got := legacyVersion(rel, shared); got != "0.9.0" {
		t.Fatalf("legacyVersion = %q, want the earliest release 0.9.0", got)
	}
	if got := legacyVersion(rel, []byte("locally edited\n")); got != "" {
		t.Fatalf("legacyVersion = %q, want \"\" for bytes no release published", got)
	}
	if got := releaseDigest("0.10.0", rel); got == "" {
		t.Fatal("releaseDigest did not find a release the manifest lists")
	}
	if got := releaseDigest("0.10.0", "clue-absent/skill.md"); got != "" {
		t.Fatalf("releaseDigest = %q for a file the release never shipped", got)
	}
}

func TestAC064_UnitNegative_MigrationRejectsCopiedOrAmbiguousCaller(t *testing.T) {
	want, err := scaffold.ManagedCarrierFiles()
	if err != nil {
		t.Fatal(err)
	}
	version, err := scaffold.PairVersion()
	if err != nil {
		t.Fatal(err)
	}
	embedded := want[".github/workflows/clue.yml"]
	current := string(embedded)
	for _, input := range []string{
		current + "\njobs:\n  extra:\n    steps:\n      - run: echo copied\n",
		strings.Replace(current, "clue-version:", "clue-version:\n  duplicate:", 1),
	} {
		if _, _, ok, message := updateCaller([]byte(input), embedded, version); ok || message == "" {
			t.Fatalf("ambiguous caller was accepted: %q", message)
		}
	}
	if _, _, ok, message := updateCaller(embedded, []byte("name: broken\n"), version); ok || !strings.Contains(message, "embedded") {
		t.Fatalf("invalid embedded caller was accepted: %q", message)
	}
}

func TestAC064_UnitNegative_MigrationDoesNotWriteThroughManagedSymlink(t *testing.T) {
	root := migrationFixture(t, "")
	mirror := filepath.Join(root, ".claude", "skills")
	if err := os.RemoveAll(mirror); err != nil {
		t.Fatal(err)
	}
	external := t.TempDir()
	if err := os.Symlink(external, mirror); err != nil {
		t.Skipf("symlinks are unavailable in this test environment: %v", err)
	}

	plan, err := Plan(root, Options{ReversalCost: "low"})
	if err != nil {
		t.Fatal(err)
	}
	for _, change := range plan.Changes {
		if strings.HasPrefix(change.Path, ".claude/") {
			t.Fatalf("migration planned a write through a managed symlink: %+v", change)
		}
	}
	if len(plan.Notices) == 0 {
		t.Fatal("symlinked managed mirror was not reported as a notice")
	}
	if entries, err := os.ReadDir(external); err != nil {
		t.Fatal(err)
	} else if len(entries) != 0 {
		t.Fatalf("planning wrote into the symlink target: %v", entries)
	}
}

func migrationFixture(t *testing.T, artifact string) string {
	t.Helper()
	root := t.TempDir()
	write := func(rel string, data []byte) {
		t.Helper()
		full := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, data, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	fixtureRoot := filepath.Join(repositoryRoot(t), "internal", "migrate", "testdata", "pre-contract")
	for _, rel := range []string{"docs/README.md", "docs/analysis/README.md"} {
		data, err := os.ReadFile(filepath.Join(fixtureRoot, filepath.FromSlash(rel)))
		if err != nil {
			t.Fatal(err)
		}
		write(rel, data)
	}
	if artifact == "" {
		var err error
		artifactBytes, err := os.ReadFile(filepath.Join(fixtureRoot, "docs", "analysis", "AN-001.md"))
		if err != nil {
			t.Fatal(err)
		}
		artifact = string(artifactBytes)
	}
	write("docs/analysis/AN-001.md", []byte(artifact))

	carriers, err := scaffold.ManagedCarrierFiles()
	if err != nil {
		t.Fatal(err)
	}
	for rel, data := range carriers {
		write(rel, data)
	}
	oldSkill, err := os.ReadFile(filepath.Join(fixtureRoot, ".agents", "skills", "clue-analysis", "skill.md"))
	if err != nil {
		t.Fatal(err)
	}
	write(".agents/skills/clue-analysis/skill.md", oldSkill)
	return root
}

func repositoryRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot locate migration fixture")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

// TestAC064_UnitNegative_BareReferencesAreReportedNeverRepaired covers the
// report-only migration: a bare forge reference is named with its file and
// line, and nothing is written. Repairing one would mean guessing which
// repository was meant, which is how a confidently wrong reference becomes a
// confidently wrong qualified one.
func TestAC064_UnitNegative_BareReferencesAreReportedNeverRepaired(t *testing.T) {
	root := t.TempDir()
	docs := filepath.Join(root, "docs", "analysis")
	if err := os.MkdirAll(docs, 0o755); err != nil {
		t.Fatal(err)
	}
	body := "---\nid: AN-001\ntype: analysis\nstatus: active\nlinks: []\ntitle: t\nprovenance: inferred\nreversal-cost: low\n---\n\nTracked by issue #202 upstream.\n"
	file := filepath.Join(docs, "AN-001-x.md")
	if err := os.WriteFile(file, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}

	plan, err := Plan(root, Options{})
	if err != nil {
		t.Fatal(err)
	}
	var found *Finding
	for i := range plan.Findings {
		if plan.Findings[i].Migration == MigrationQualifiedReferences {
			found = &plan.Findings[i]
		}
	}
	if found == nil {
		t.Fatal("expected the bare reference to be reported")
	}
	if !strings.Contains(found.Message, "#202") || !strings.Contains(found.Message, "line ") {
		t.Fatalf("expected the number and its line, got %q", found.Message)
	}
	for _, c := range plan.Changes {
		if c.Migration == MigrationQualifiedReferences {
			t.Fatal("a bare reference must never be repaired mechanically")
		}
	}
	after, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != body {
		t.Fatal("the file was modified by a report-only migration")
	}
}

// TestAC064_UnitNegative_ManifestEntriesMatchTheReleaseTheyName checks the one
// thing about legacyDigests that is true at every moment.
//
// The manifest is hand-maintained and append-only: a generated carrier is
// replaceable only when its bytes match the release that wrote it. Each entry
// is therefore a claim about a published tag, and the claim is verifiable —
// the tag is in this repository.
//
// The first version of this guard compared the entry for the *stamped* version
// against the working tree instead, and that was wrong in a way worth recording.
// Between releases the working tree legitimately diverges from the last release,
// so an ordinary skill edit failed the test, and its advice — update the entry —
// would have rewritten the record of what v0.11.2 actually shipped. A guard that
// tells you to falsify the data it protects is worse than no guard.
//
// The missing-entry half genuinely only holds while a release is being cut, so
// it moved to .github/scripts/release-gates.sh, which knows it is looking at a
// release because the stamp changed.
func TestAC080_UnitPositive_RecognizedReleasedSetAddsNewManagedSkill(t *testing.T) {
	root := migrationFixture(t, "")
	repoRoot := repositoryRoot(t)
	const previous = "0.12.0"
	if !tagExists(t, repoRoot, "v"+previous) {
		t.Skipf("v%s not present in this checkout", previous)
	}

	for _, release := range legacyDigests {
		if release.Version != previous {
			continue
		}
		for rel := range release.Files {
			data, err := gitShow(t, repoRoot, "v"+previous+":.agents/skills/"+rel)
			if err != nil {
				t.Fatal(err)
			}
			path := filepath.Join(root, ".agents", "skills", filepath.FromSlash(rel))
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, data, 0o644); err != nil {
				t.Fatal(err)
			}
		}
	}
	if err := os.RemoveAll(filepath.Join(root, ".agents", "skills", "clue-upgrade")); err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll(filepath.Join(root, ".claude", "skills", "clue-upgrade")); err != nil {
		t.Fatal(err)
	}

	preview, err := Plan(root, Options{ReversalCost: "low"})
	if err != nil {
		t.Fatal(err)
	}
	if len(preview.Findings) != 0 {
		t.Fatalf("recognized released set must stay migratable, got findings: %+v", preview.Findings)
	}
	var added *Change
	for i := range preview.Changes {
		if preview.Changes[i].Path == ".agents/skills/clue-upgrade/skill.md" {
			added = &preview.Changes[i]
			break
		}
	}
	if added == nil {
		t.Fatalf("migration did not add the new managed skill: %+v", preview.Changes)
	}
	if !strings.Contains(added.Description, "with release "+preview.Target+" content") || !strings.Contains(added.Description, "recognized as release "+previous) {
		t.Fatalf("new-carrier preview does not name both release roles: %q", added.Description)
	}
	if err := Apply(root, preview); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, ".agents", "skills", "clue-upgrade", "skill.md")); err != nil {
		t.Fatalf("migration did not create the whole new managed skill directory: %v", err)
	}
	repeat, err := Plan(root, Options{ReversalCost: "low"})
	if err != nil {
		t.Fatal(err)
	}
	if len(repeat.Changes) != 0 || len(repeat.Findings) != 0 {
		t.Fatalf("migrated release is not a no-op: changes=%+v findings=%+v", repeat.Changes, repeat.Findings)
	}
}

func TestAC080_UnitNegative_PartialManagedSetCannotAuthorizeNewSkill(t *testing.T) {
	root := migrationFixture(t, "")
	repoRoot := repositoryRoot(t)
	const previous = "0.12.0"
	if !tagExists(t, repoRoot, "v"+previous) {
		t.Skipf("v%s not present in this checkout", previous)
	}

	for _, release := range legacyDigests {
		if release.Version != previous {
			continue
		}
		for rel := range release.Files {
			data, err := gitShow(t, repoRoot, "v"+previous+":.agents/skills/"+rel)
			if err != nil {
				t.Fatal(err)
			}
			path := filepath.Join(root, ".agents", "skills", filepath.FromSlash(rel))
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, data, 0o644); err != nil {
				t.Fatal(err)
			}
		}
	}
	if err := os.Remove(filepath.Join(root, ".agents", "skills", "clue-plan", "skill.md")); err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll(filepath.Join(root, ".agents", "skills", "clue-upgrade")); err != nil {
		t.Fatal(err)
	}

	preview, err := Plan(root, Options{ReversalCost: "low"})
	if err != nil {
		t.Fatal(err)
	}
	if len(preview.Findings) == 0 {
		t.Fatal("partial managed set must block the migration")
	}
	for _, change := range preview.Changes {
		if change.Path == ".agents/skills/clue-upgrade/skill.md" {
			t.Fatalf("partial managed set planned an unsafe new skill addition: %+v", change)
		}
	}
	for _, finding := range preview.Findings {
		if finding.Path == ".agents/skills/clue-upgrade/skill.md" && strings.Contains(finding.Message, "managed carrier is missing") {
			return
		}
	}
	t.Fatalf("new skill was not held as a finding: %+v", preview.Findings)
}

func TestAC064_UnitNegative_ManifestEntriesMatchTheReleaseTheyName(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot locate the test file")
	}
	root := filepath.Join(filepath.Dir(thisFile), "..", "..")

	checked := 0
	for _, entry := range legacyDigests {
		tag := "v" + entry.Version
		if !tagExists(t, root, tag) {
			// A shallow clone or a fork without tags cannot answer, and
			// failing there would punish the checkout rather than the data.
			t.Logf("%s not present in this checkout; skipping", tag)
			continue
		}
		checked++
		for rel, want := range entry.Files {
			data, err := gitShow(t, root, tag+":.agents/skills/"+rel)
			if err != nil {
				t.Errorf("the %s manifest entry names %s, which %s does not contain: %v", entry.Version, rel, tag, err)
				continue
			}
			sum := sha256.Sum256(data)
			if got := hex.EncodeToString(sum[:]); got != want {
				t.Errorf("%s digest in the %s manifest entry does not match what %s published.\n  manifest: %s\n  %s:  %s",
					rel, entry.Version, tag, want, tag, got)
			}
		}
	}
	if checked == 0 {
		t.Skip("no release tags in this checkout")
	}
}

func tagExists(t *testing.T, root, tag string) bool {
	t.Helper()
	cmd := exec.Command("git", "-C", root, "rev-parse", "-q", "--verify", "refs/tags/"+tag)
	return cmd.Run() == nil
}

func gitShow(t *testing.T, root, spec string) ([]byte, error) {
	t.Helper()
	var out, errOut bytes.Buffer
	cmd := exec.Command("git", "-C", root, "show", spec)
	cmd.Stdout = &out
	// Without this the caller reports "exit status 128" for every failure,
	// which names neither the path nor the reason. git already says both.
	cmd.Stderr = &errOut
	if err := cmd.Run(); err != nil {
		if msg := strings.TrimSpace(errOut.String()); msg != "" {
			return nil, fmt.Errorf("%s: %s", err, msg)
		}
		return nil, err
	}
	return out.Bytes(), nil
}
