package carriers

import (
	"os"
	"path/filepath"
	"testing"
)

func writeFiles(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for p, content := range files {
		full := filepath.Join(root, filepath.FromSlash(p))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

const workflowContent = "name: ci\non: [push]\n"

func mustFingerprint(t *testing.T, root, targetPath string) string {
	t.Helper()
	fp, err := Fingerprint(filepath.Join(root, filepath.FromSlash(targetPath)))
	if err != nil {
		t.Fatal(err)
	}
	return fp
}

// TestAC118_UnitPositive_cleanRunPasses proves a clean reconciliation run —
// every mapped fingerprint matches, no deleted path is still referenced, no
// target is missing — produces no findings, and that reconciling twice
// against the same inputs is deterministic.
func TestAC118_UnitPositive_cleanRunPasses(t *testing.T) {
	root := writeFiles(t, map[string]string{
		".github/workflows/ci.yml": workflowContent,
		"docs/decisions/README.md": "# decisions\n",
	})
	fp := mustFingerprint(t, root, ".github/workflows/ci.yml")
	inv := Inventory{
		SourceRevision: "rev-1",
		SourceLocation: "source",
		Entries: []Entry{
			{ID: "CARRIER-1", Kind: KindWorkflow, SourcePath: "openspec/.github/workflows/ci.yml", TargetPath: ".github/workflows/ci.yml", Fingerprint: fp},
		},
	}
	r1, err := Reconcile(inv, root)
	if err != nil {
		t.Fatal(err)
	}
	if r1.Failed() {
		t.Fatalf("expected a clean report, got %v", r1.Findings)
	}
	r2, err := Reconcile(inv, root)
	if err != nil {
		t.Fatal(err)
	}
	if len(r1.Findings) != len(r2.Findings) {
		t.Fatalf("expected deterministic reports, got %v and %v", r1.Findings, r2.Findings)
	}
}

// TestAC118_UnitNegative_dirtyRunFails proves an altered target still fails
// after a clean baseline, so a clean run is not a permanent exemption.
func TestAC118_UnitNegative_dirtyRunFails(t *testing.T) {
	root := writeFiles(t, map[string]string{
		".github/workflows/ci.yml": workflowContent,
	})
	fp := mustFingerprint(t, root, ".github/workflows/ci.yml")
	inv := Inventory{
		SourceRevision: "rev-1",
		SourceLocation: "source",
		Entries: []Entry{
			{ID: "CARRIER-1", Kind: KindWorkflow, SourcePath: "openspec/.github/workflows/ci.yml", TargetPath: ".github/workflows/ci.yml", Fingerprint: fp},
		},
	}
	if err := os.WriteFile(filepath.Join(root, filepath.FromSlash(".github/workflows/ci.yml")), []byte("name: ci\non: [push, pull_request]\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	r, err := Reconcile(inv, root)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Failed() {
		t.Fatal("expected a failing report after the mapped target drifted")
	}
}

// TestAC119_UnitPositive_blockedEntryIsNotReconciled proves a blocked entry
// with no target is not reconciled against anything and produces no
// finding — its presence alone is the record of a known gap.
func TestAC119_UnitPositive_blockedEntryIsNotReconciled(t *testing.T) {
	root := writeFiles(t, map[string]string{"docs/decisions/README.md": "# decisions\n"})
	inv := Inventory{
		SourceRevision: "rev-1",
		SourceLocation: "source",
		Entries: []Entry{
			{ID: "CARRIER-2", Kind: KindFreshnessInput, SourcePath: "openspec/config.yaml", Blocked: true, Reason: "no target chosen yet"},
		},
	}
	r, err := Reconcile(inv, root)
	if err != nil {
		t.Fatal(err)
	}
	if r.Failed() {
		t.Fatalf("expected no finding for a blocked entry, got %v", r.Findings)
	}
}

// TestAC119_UnitNegative_incompleteEntryRejectedAtLoad proves an inventory
// entry naming neither a mapped target nor a blocked marker is rejected at
// load — every carrier the rehearsal finds gets exactly one outcome.
func TestAC119_UnitNegative_incompleteEntryRejectedAtLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), "inventory.yaml")
	manifest := "source-revision: rev-1\nsource-location: source\nentries:\n  - id: CARRIER-3\n    kind: link\n    source-path: openspec/README.md\n"
	if err := os.WriteFile(path, []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadInventory(path); err == nil {
		t.Fatal("expected an entry with neither a target nor blocked to be rejected")
	}
}

// TestAC119_UnitNegative_combinedTargetAndBlockedRejectedAtLoad proves an
// entry cannot combine target-path/fingerprint with blocked.
func TestAC119_UnitNegative_combinedTargetAndBlockedRejectedAtLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), "inventory.yaml")
	manifest := "source-revision: rev-1\nsource-location: source\nentries:\n  - id: CARRIER-4\n    kind: link\n    source-path: openspec/README.md\n    target-path: README.md\n    fingerprint: abc\n    blocked: true\n    reason: mixed\n"
	if err := os.WriteFile(path, []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadInventory(path); err == nil {
		t.Fatal("expected an entry combining a mapped target with blocked to be rejected")
	}
}

// TestAC119_UnitNegative_unsupportedKindRejectedAtLoad proves the six-kind
// vocabulary is enforced at load.
func TestAC119_UnitNegative_unsupportedKindRejectedAtLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), "inventory.yaml")
	manifest := "source-revision: rev-1\nsource-location: source\nentries:\n  - id: CARRIER-5\n    kind: dashboard\n    source-path: openspec/README.md\n    target-path: README.md\n    fingerprint: abc\n"
	if err := os.WriteFile(path, []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadInventory(path); err == nil {
		t.Fatal("expected an unsupported kind to be rejected")
	}
}

// TestAC119_UnitNegative_duplicateIDRejectedAtLoad proves two entries
// sharing an ID are rejected — unlike parity, a carrier row has no
// several-proof-outcomes-per-ID shape.
func TestAC119_UnitNegative_duplicateIDRejectedAtLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), "inventory.yaml")
	manifest := "source-revision: rev-1\nsource-location: source\nentries:\n" +
		"  - id: CARRIER-6\n    kind: link\n    source-path: a.md\n    target-path: a.md\n    fingerprint: abc\n" +
		"  - id: CARRIER-6\n    kind: link\n    source-path: b.md\n    target-path: b.md\n    fingerprint: def\n"
	if err := os.WriteFile(path, []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadInventory(path); err == nil {
		t.Fatal("expected a duplicate id to be rejected")
	}
}

// TestAC120_UnitPositive_matchingFingerprintPasses proves a mapped entry
// whose current target content matches the pinned fingerprint produces no
// finding.
func TestAC120_UnitPositive_matchingFingerprintPasses(t *testing.T) {
	root := writeFiles(t, map[string]string{"AGENTS.md": "# agents\n"})
	fp := mustFingerprint(t, root, "AGENTS.md")
	inv := Inventory{
		SourceRevision: "rev-1", SourceLocation: "source",
		Entries: []Entry{{ID: "CARRIER-7", Kind: KindInstruction, SourcePath: "openspec/AGENTS.md", TargetPath: "AGENTS.md", Fingerprint: fp}},
	}
	r, err := Reconcile(inv, root)
	if err != nil {
		t.Fatal(err)
	}
	if r.Failed() {
		t.Fatalf("expected no finding for a matching fingerprint, got %v", r.Findings)
	}
}

// TestAC120_UnitNegative_lostFingerprintFails proves a mapped entry whose
// target content no longer matches the pinned fingerprint fails as a lost
// fingerprint.
func TestAC120_UnitNegative_lostFingerprintFails(t *testing.T) {
	root := writeFiles(t, map[string]string{"AGENTS.md": "# agents\n"})
	inv := Inventory{
		SourceRevision: "rev-1", SourceLocation: "source",
		Entries: []Entry{{ID: "CARRIER-7", Kind: KindInstruction, SourcePath: "openspec/AGENTS.md", TargetPath: "AGENTS.md", Fingerprint: "0000000000000000000000000000000000000000000000000000000000000000"}},
	}
	r, err := Reconcile(inv, root)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Failed() || r.Findings[0].Class != ClassLostFingerprint {
		t.Fatalf("expected a lost-fingerprint finding, got %v", r.Findings)
	}
}

// TestAC121_UnitPositive_presentTargetIsNotMissing proves a mapped entry
// whose target path exists produces no missing-asset finding.
func TestAC121_UnitPositive_presentTargetIsNotMissing(t *testing.T) {
	root := writeFiles(t, map[string]string{"docs/decisions/README.md": "# decisions\n"})
	fp := mustFingerprint(t, root, "docs/decisions/README.md")
	inv := Inventory{
		SourceRevision: "rev-1", SourceLocation: "source",
		Entries: []Entry{{ID: "CARRIER-8", Kind: KindRegistry, SourcePath: "openspec/decisions/README.md", TargetPath: "docs/decisions/README.md", Fingerprint: fp}},
	}
	r, err := Reconcile(inv, root)
	if err != nil {
		t.Fatal(err)
	}
	if r.Failed() {
		t.Fatalf("expected no finding for a present target, got %v", r.Findings)
	}
}

// TestAC121_UnitNegative_missingTargetFails proves a mapped entry whose
// target path does not exist in the reconciled corpus fails as a missing
// asset.
func TestAC121_UnitNegative_missingTargetFails(t *testing.T) {
	root := writeFiles(t, map[string]string{"docs/decisions/README.md": "# decisions\n"})
	inv := Inventory{
		SourceRevision: "rev-1", SourceLocation: "source",
		Entries: []Entry{{ID: "CARRIER-8", Kind: KindDiagramAsset, SourcePath: "openspec/diagram.svg", TargetPath: "docs/architecture/diagram.svg", Fingerprint: "deadbeef"}},
	}
	r, err := Reconcile(inv, root)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Failed() || r.Findings[0].Class != ClassMissingAsset {
		t.Fatalf("expected a missing-asset finding, got %v", r.Findings)
	}
}

const criteriaReferencingDeletedPath = `---
id: CARRIERTEST-criteria
type: criteria
status: active
links: []
title: Carrier fixture criteria
---

See [the old mapping](openspec/specs/foo/spec.md) for context.
`

// TestAC122_UnitPositive_liveLinkIsNotStale proves a link that does not
// resolve to a deleted path produces no finding.
func TestAC122_UnitPositive_liveLinkIsNotStale(t *testing.T) {
	root := writeFiles(t, map[string]string{
		"docs/capabilities/CARRIERTEST/criteria.md": "---\nid: CARRIERTEST2-criteria\ntype: criteria\nstatus: active\nlinks: []\ntitle: fixture\n---\n\nSee [design](design.md) for context.\n",
	})
	inv := Inventory{SourceRevision: "rev-1", SourceLocation: "source", DeletedPaths: []string{"openspec/specs/foo/spec.md"}}
	r, err := Reconcile(inv, root)
	if err != nil {
		t.Fatal(err)
	}
	if r.Failed() {
		t.Fatalf("expected no finding for a link that does not resolve to a deleted path, got %v", r.Findings)
	}
}

// TestAC122_UnitNegative_staleDeletedPathReferenceFails proves a Markdown
// link still resolving to a path the inventory names as deleted fails as a
// stale deleted-path reference.
func TestAC122_UnitNegative_staleDeletedPathReferenceFails(t *testing.T) {
	root := writeFiles(t, map[string]string{
		"docs/capabilities/CARRIERTEST/criteria.md": criteriaReferencingDeletedPath,
	})
	inv := Inventory{SourceRevision: "rev-1", SourceLocation: "source", DeletedPaths: []string{"openspec/specs/foo/spec.md"}}
	r, err := Reconcile(inv, root)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Failed() || r.Findings[0].Class != ClassStaleDeletedPath {
		t.Fatalf("expected a stale-deleted-path finding, got %v", r.Findings)
	}
}
