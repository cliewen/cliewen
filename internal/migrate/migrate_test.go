package migrate

import (
	"bytes"
	"os"
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
	want := []string{MigrationReversalCost, MigrationStatusLifecycle, MigrationManagedCarriers}
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
