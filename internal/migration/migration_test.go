package migration

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cliewen/cliewen/internal/carriers"
	"github.com/cliewen/cliewen/internal/corpus"
	"github.com/cliewen/cliewen/internal/importedchange"
	"github.com/cliewen/cliewen/internal/ledger"
	"github.com/cliewen/cliewen/internal/parity"
)

func writeFixtureFiles(t *testing.T, root string, files map[string]string) {
	t.Helper()
	for name, content := range files {
		path := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func fixtureTarget(t *testing.T, prefix, criterion, revision, location string) string {
	t.Helper()
	root := t.TempDir()
	criteria := fmt.Sprintf("---\nid: %s-criteria\ntype: criteria\nstatus: active\nlinks: [CAP-901]\ntitle: Fixture criteria\nac-prefix: %s\n---\n\n```gherkin\nFeature: Fixture migration\n\n  @%s\n  Scenario: source proof survives migration\n    Test-type: Integration\n    Given a source fixture\n    When it is migrated\n    Then its proof is preserved\n```\n", prefix, prefix, criterion)
	component := strings.TrimLeft(criterion[strings.LastIndex(criterion, "-")+1:], "0")
	if component == "" {
		component = "0"
	}
	archived := ""
	if prefix == "ARC" {
		archived = fmt.Sprintf("  - id: %s-099\n    kind: numeric\n    state: retired\n    prefix: %s\n    component: 99\n    source-revision: %s\n    source-location: %s\n", prefix, prefix, revision, location)
	}
	ledgerFile := fmt.Sprintf("counters:\n  G: 1\n  CAP: 901\n  %s: %s\nentries:\n  - id: G-001\n    kind: numeric\n    state: live\n    prefix: G\n    component: 1\n  - id: CAP-901\n    kind: numeric\n    state: live\n    prefix: CAP\n    component: 901\n%s  - id: %s\n    kind: numeric\n    state: live\n    prefix: %s\n    component: %s\n    source-revision: %s\n    source-location: %s\n  - id: %s-criteria\n    kind: opaque\n    state: live\n", prefix, component, archived, criterion, prefix, component, revision, location, prefix)
	writeFixtureFiles(t, root, map[string]string{
		"docs/README.md": "# fixture\n\n<!-- clue:index:start -->\n- [goals/](goals/) — fixture goal\n- [capabilities/](capabilities/) — fixture capability\n<!-- clue:index:end -->\n",
		"docs/goals/README.md": "# goals\n\n<!-- clue:index:start -->\n- [G-001.md](G-001.md) — fixture goal\n<!-- clue:index:end -->\n",
		"docs/goals/G-001.md": "---\nid: G-001\ntype: goal\nstatus: accepted\nlinks: []\ntitle: Fixture goal\n---\n\n# Fixture goal\n",
		"docs/capabilities/README.md": "# capabilities\n\n<!-- clue:index:start -->\n- [fixture/](fixture/) — fixture capability\n<!-- clue:index:end -->\n",
		"docs/capabilities/fixture/README.md": "---\nid: CAP-901\ntype: capability\nstatus: active\nlinks: [G-001]\ntitle: Fixture capability\ngoal: G-001\n---\n\n# Fixture capability\n",
		"docs/capabilities/fixture/criteria.md": criteria,
		"fixture_test.go": "package fixture\n\nfunc Test" + strings.ReplaceAll(criterion, "-", "") + "_IntegrationPositive_preservesProof(t *testing.T) {}\nfunc Test" + strings.ReplaceAll(criterion, "-", "") + "_IntegrationNegative_rejectsLoss(t *testing.T) {}\n",
		".clue/id-ledger.yaml": ledgerFile,
	})
	return root
}

func mustValidateFixture(t *testing.T, root string) {
	t.Helper()
	c, scanIssues := corpus.Scan(root)
	if len(scanIssues) > 0 {
		t.Fatalf("fixture scan issues: %v", scanIssues)
	}
	if issues := corpus.Validate(c, corpus.Options{Version: "dev"}); len(issues) > 0 {
		t.Fatalf("fixture validation issues: %v", issues)
	}
}

func mustHaveFinding(t *testing.T, report parity.Report, class string) {
	t.Helper()
	for _, finding := range report.Findings {
		if finding.Class == class {
			return
		}
	}
	t.Fatalf("expected %s finding, got %v", class, report.Findings)
}

func mustHaveCarrierFinding(t *testing.T, report carriers.Report, class string) {
	t.Helper()
	for _, finding := range report.Findings {
		if finding.Class == class {
			return
		}
	}
	t.Fatalf("expected %s finding, got %v", class, report.Findings)
}

// TestAC123_UnitPositive_disposableFixturesProveComposedMigrationContract
// exercises the approved mutation phase against disposable source shapes.
// It deliberately never executes the source fixture's own test suite.
func TestAC123_UnitPositive_disposableFixturesProveComposedMigrationContract(t *testing.T) {
	numericRevision, numericLocation := "numeric-archive-fixture-v1", "fixtures/numeric-archive"
	numericRoot := fixtureTarget(t, "ARC", "ARC-100", numericRevision, numericLocation)
	mustValidateFixture(t, numericRoot)
	numericSource := parity.SourceManifest{SourceRevision: numericRevision, SourceLocation: numericLocation, Entries: []parity.SourceEntry{
		{ID: "ARC-100", ProofClass: "Integration", Direction: "positive", EvidenceLocation: "fixture_test.go"},
		{ID: "ARC-100", ProofClass: "Integration", Direction: "negative", EvidenceLocation: "fixture_test.go"},
	}}
	numericTarget, err := parity.DeriveTargetManifest(numericRoot)
	if err != nil {
		t.Fatal(err)
	}
	if report := parity.Compare(numericSource, numericTarget); report.Failed() {
		t.Fatalf("numeric fixture parity: %v", report.Findings)
	}

	numericLedger, err := ledger.Load(numericRoot)
	if err != nil {
		t.Fatal(err)
	}
	if err := numericLedger.ReserveOpaque("ARC-099", numericRevision, numericLocation); err == nil {
		t.Fatal("numeric fixture reused an archived identity")
	}

	opaqueRevision, opaqueLocation := "opaque-identifier-fixture-v1", "fixtures/opaque-identifier"
	opaqueRoot := fixtureTarget(t, "OPA", "OPA-001", opaqueRevision, opaqueLocation)
	mustValidateFixture(t, opaqueRoot)
	opaqueSource := parity.SourceManifest{SourceRevision: opaqueRevision, SourceLocation: opaqueLocation, Entries: []parity.SourceEntry{
		{ID: "OPA-001", ProofClass: "Integration", Direction: "positive", EvidenceLocation: "fixture_test.go"},
		{ID: "OPA-001", ProofClass: "Integration", Direction: "negative", EvidenceLocation: "fixture_test.go"},
	}}
	opaqueTarget, err := parity.DeriveTargetManifest(opaqueRoot)
	if err != nil {
		t.Fatal(err)
	}
	if report := parity.Compare(opaqueSource, opaqueTarget); report.Failed() {
		t.Fatalf("opaque fixture parity: %v", report.Findings)
	}

	opaqueLedger, err := ledger.Load(opaqueRoot)
	if err != nil {
		t.Fatal(err)
	}
	const opaqueID = "8f14e45f-ceea-467e-9a2b-a1c8b9d2f7a1"
	if err := opaqueLedger.ReserveOpaque(opaqueID, opaqueRevision, opaqueLocation); err != nil {
		t.Fatal(err)
	}
	if err := opaqueLedger.ReserveOpaque(opaqueID, opaqueRevision, opaqueLocation); err == nil {
		t.Fatal("opaque fixture reused a reserved source-owned identity")
	}

	carrierPath := "docs/capabilities/fixture/README.md"
	fingerprint, err := carriers.Fingerprint(filepath.Join(numericRoot, filepath.FromSlash(carrierPath)))
	if err != nil {
		t.Fatal(err)
	}
	inventory := carriers.Inventory{SourceRevision: numericRevision, SourceLocation: numericLocation, Entries: []carriers.Entry{{ID: "NUMERIC-README", Kind: carriers.KindInstruction, SourcePath: "openspec/AGENTS.md", TargetPath: carrierPath, Fingerprint: fingerprint}}}
	if report, err := carriers.Reconcile(inventory, numericRoot); err != nil || report.Failed() {
		t.Fatalf("numeric fixture carriers: report=%v err=%v", report.Findings, err)
	}

	links := importedchange.ParseProofLinks("## Proof links\n\n| Task | Criterion |\n| --- | --- |\n| Preserve archive | ARC-100 |\n")
	if len(links) != 1 || links[0].Criterion != "ARC-100" {
		t.Fatalf("pending source work is not inspectable: %v", links)
	}
}

// TestAC123_UnitNegative_disposableFixturesRejectRequiredFailurePaths keeps
// the mandatory failure reports separate from the clean composed fixture.
func TestAC123_UnitNegative_disposableFixturesRejectRequiredFailurePaths(t *testing.T) {
	source := parity.SourceManifest{SourceRevision: "rev-1", SourceLocation: "fixture", Entries: []parity.SourceEntry{{ID: "ARC-100", ProofClass: "Unit", Direction: "positive", EvidenceLocation: "fixture_test.go"}}}
	mustHaveFinding(t, parity.Compare(source, parity.TargetManifest{}), parity.ClassMissingCriterion)
	mustHaveFinding(t, parity.Compare(parity.SourceManifest{}, parity.TargetManifest{Entries: map[string]parity.TargetEntry{"ARC-100": {ID: "ARC-100", EvidenceLocations: []string{"fixture_test.go"}}}}), parity.ClassOrphanedTag)
	mustHaveFinding(t, parity.Compare(source, parity.TargetManifest{Entries: map[string]parity.TargetEntry{"ARC-100": {ID: "ARC-100", ProofClass: "Unit", Directions: []string{"negative"}, EvidenceLocations: []string{"moved_test.go"}}}}), parity.ClassChangedEvidence)
	mustHaveFinding(t, parity.Compare(source, parity.TargetManifest{Entries: map[string]parity.TargetEntry{"ARC-100": {ID: "ARC-100", ProofClass: "Unit", Directions: []string{"positive"}, EvidenceLocations: []string{"fixture_test.go"}, SourceRevision: "rev-2"}}}), parity.ClassStaleFingerprint)
	mustHaveFinding(t, parity.Compare(source, parity.TargetManifest{Entries: map[string]parity.TargetEntry{"ARC-100": {ID: "ARC-100", Draft: true}}}), parity.ClassUnjustifiedDisposition)

	root := fixtureTarget(t, "ARC", "ARC-100", "rev-1", "fixture")
	stale := carriers.Inventory{SourceRevision: "rev-1", SourceLocation: "fixture", DeletedPaths: []string{"openspec/specs/numeric-archive/spec.md"}}
	path := filepath.Join(root, "docs", "capabilities", "fixture", "README.md")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, append(data, []byte("\nSee [old](../../../openspec/specs/numeric-archive/spec.md).\n")...), 0o644); err != nil {
		t.Fatal(err)
	}
	staleReport, err := carriers.Reconcile(stale, root)
	if err != nil {
		t.Fatal(err)
	}
	mustHaveCarrierFinding(t, staleReport, carriers.ClassStaleDeletedPath)

	lost := carriers.Inventory{SourceRevision: "rev-1", SourceLocation: "fixture", Entries: []carriers.Entry{{ID: "LOST", Kind: carriers.KindInstruction, SourcePath: "openspec/AGENTS.md", TargetPath: "docs/capabilities/fixture/README.md", Fingerprint: strings.Repeat("0", 64)}}}
	lostReport, err := carriers.Reconcile(lost, root)
	if err != nil {
		t.Fatal(err)
	}
	mustHaveCarrierFinding(t, lostReport, carriers.ClassLostFingerprint)

	missing := carriers.Inventory{SourceRevision: "rev-1", SourceLocation: "fixture", Entries: []carriers.Entry{{ID: "MISSING", Kind: carriers.KindDiagramAsset, SourcePath: "openspec/diagram.svg", TargetPath: "docs/architecture/diagram.svg", Fingerprint: strings.Repeat("0", 64)}}}
	missingReport, err := carriers.Reconcile(missing, root)
	if err != nil {
		t.Fatal(err)
	}
	mustHaveCarrierFinding(t, missingReport, carriers.ClassMissingAsset)
}
