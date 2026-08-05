package migration

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/cliewen/cliewen/internal/carriers"
	"github.com/cliewen/cliewen/internal/corpus"
	"github.com/cliewen/cliewen/internal/importedchange"
	"github.com/cliewen/cliewen/internal/ledger"
	"github.com/cliewen/cliewen/internal/parity"
)

var (
	clueOnce sync.Once
	cluePath string
	clueErr  error
)

func fixtureClue(t *testing.T) string {
	t.Helper()
	clueOnce.Do(func() {
		wd, err := os.Getwd()
		if err != nil {
			clueErr = err
			return
		}
		root := filepath.Clean(filepath.Join(wd, "..", ".."))
		dir, err := os.MkdirTemp("", "cliewen-migration-fixture-*")
		if err != nil {
			clueErr = err
			return
		}
		name := "clue"
		if runtime.GOOS == "windows" {
			name += ".exe"
		}
		cluePath = filepath.Join(dir, name)
		cmd := exec.Command("go", "build", "-o", cluePath, "./cmd/clue")
		cmd.Dir = root
		clueErr = cmd.Run()
	})
	if clueErr != nil {
		t.Fatal(clueErr)
	}
	return cluePath
}

func runFixtureClue(t *testing.T, args ...string) (string, error) {
	t.Helper()
	cmd := exec.Command(fixtureClue(t), args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func mustFailFixtureClue(t *testing.T, want string, args ...string) {
	t.Helper()
	out, err := runFixtureClue(t, args...)
	if err == nil || !strings.Contains(out, want) {
		t.Fatalf("expected command failure containing %q, err=%v output=%s", want, err, out)
	}
}

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
	counter := component
	if prefix == "ARC" {
		counter = "99"
	}
	ledgerFile := fmt.Sprintf("counters:\n  G: 1\n  CAP: 901\n  IC: 901\n  %s: %s\nentries:\n  - id: G-001\n    kind: numeric\n    state: live\n    prefix: G\n    component: 1\n  - id: CAP-901\n    kind: numeric\n    state: live\n    prefix: CAP\n    component: 901\n  - id: IC-901\n    kind: numeric\n    state: live\n    prefix: IC\n    component: 901\n%s  - id: %s\n    kind: numeric\n    state: live\n    prefix: %s\n    component: %s\n    source-revision: %s\n    source-location: %s\n  - id: %s-criteria\n    kind: opaque\n    state: live\n", prefix, counter, archived, criterion, prefix, component, revision, location, prefix)
	writeFixtureFiles(t, root, map[string]string{
		"docs/README.md": "# fixture\n\n<!-- clue:index:start -->\n- [goals/](goals/) — fixture goal\n- [capabilities/](capabilities/) — fixture capability\n- [imported-changes/](imported-changes/) — fixture imported work\n<!-- clue:index:end -->\n",
		"docs/goals/README.md": "# goals\n\n<!-- clue:index:start -->\n- [G-001.md](G-001.md) — fixture goal\n<!-- clue:index:end -->\n",
		"docs/goals/G-001.md": "---\nid: G-001\ntype: goal\nstatus: accepted\nlinks: []\ntitle: Fixture goal\n---\n\n# Fixture goal\n",
		"docs/capabilities/README.md": "# capabilities\n\n<!-- clue:index:start -->\n- [fixture/](fixture/) — fixture capability\n<!-- clue:index:end -->\n",
		"docs/capabilities/fixture/README.md": "---\nid: CAP-901\ntype: capability\nstatus: active\nlinks: [G-001]\ntitle: Fixture capability\ngoal: G-001\n---\n\n# Fixture capability\n",
		"docs/capabilities/fixture/criteria.md": criteria,
		"docs/imported-changes/README.md": "# imported changes\n\n<!-- clue:index:start -->\n- [IC-901.md](IC-901.md) — fixture source work\n<!-- clue:index:end -->\n",
		"docs/imported-changes/IC-901.md": fmt.Sprintf("---\nid: IC-901\ntype: imported-change\nstatus: complete\nlinks: [CAP-901]\ntitle: Fixture pending source work\nsource-revision: %s\nsource-location: %s\n---\n\n# Fixture pending source work\n\n## Intent\n\nPreserve the source change.\n\n## Design rationale\n\nThe target retains the proof link.\n\n## Dependencies\n\nNone.\n\n## Proof links\n\n| Task | Criterion |\n| --- | --- |\n| Preserve source proof | %s |\n", revision, location, criterion),
		"AGENTS.md": "# fixture routing\n",
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

func writeSourceManifest(t *testing.T, root, revision, location, criterion string) string {
	t.Helper()
	path := filepath.Join(root, "source-manifest.yaml")
	content := fmt.Sprintf("source-revision: %s\nsource-location: %s\nentries:\n  - id: %s\n    proof-class: Integration\n    direction: positive\n    evidence-location: fixture_test.go\n  - id: %s\n    proof-class: Integration\n    direction: negative\n    evidence-location: fixture_test.go\n", revision, location, criterion, criterion)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func writeCarrierInventory(t *testing.T, root, revision, location, prefix string) string {
	t.Helper()
	paths := []string{"AGENTS.md", "docs/capabilities/README.md", "docs/capabilities/fixture/README.md"}
	ids := []string{prefix + "-INSTRUCTION", prefix + "-REGISTRY", prefix + "-LINK"}
	kinds := []carriers.Kind{carriers.KindInstruction, carriers.KindRegistry, carriers.KindLink}
	sources := []string{"openspec/AGENTS.md", "openspec/INDEX.md", "openspec/specs/fixture/spec.md"}
	var entries strings.Builder
	for i, target := range paths {
		fp, err := carriers.Fingerprint(filepath.Join(root, filepath.FromSlash(target)))
		if err != nil {
			t.Fatal(err)
		}
		fmt.Fprintf(&entries, "  - id: %s\n    kind: %s\n    source-path: %s\n    target-path: %s\n    fingerprint: %s\n", ids[i], kinds[i], sources[i], target, fp)
	}
	path := filepath.Join(root, "carrier-inventory.yaml")
	content := fmt.Sprintf("source-revision: %s\nsource-location: %s\ndeleted-paths:\n  - openspec/\n  - openspec/specs/fixture/spec.md\nentries:\n%s", revision, location, entries.String())
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
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
	numericRoot := fixtureTarget(t, "ARC", "ARC-002", numericRevision, numericLocation)
	mustValidateFixture(t, numericRoot)
	numericManifest := writeSourceManifest(t, numericRoot, numericRevision, numericLocation, "ARC-002")
	numericInventory := writeCarrierInventory(t, numericRoot, numericRevision, numericLocation, "NUMERIC")
	if out, err := runFixtureClue(t, "validate", numericRoot); err != nil {
		t.Fatalf("numeric target validation: %v\n%s", err, out)
	}
	if out, err := runFixtureClue(t, "parity", numericManifest, numericRoot); err != nil {
		t.Fatalf("numeric fixture parity: %v\n%s", err, out)
	}
	if out, err := runFixtureClue(t, "carriers", numericInventory, numericRoot); err != nil {
		t.Fatalf("numeric fixture carriers: %v\n%s", err, out)
	}

	numericLedger, err := ledger.Load(numericRoot)
	if err != nil {
		t.Fatal(err)
	}
	if id, err := numericLedger.NextNumeric("ARC"); err != nil || id != "ARC-100" {
		t.Fatalf("numeric fixture next id = %q, %v; want ARC-100 after archived ARC-099", id, err)
	}

	opaqueRevision, opaqueLocation := "opaque-identifier-fixture-v1", "fixtures/opaque-identifier"
	opaqueRoot := fixtureTarget(t, "OPA", "OPA-001", opaqueRevision, opaqueLocation)
	mustValidateFixture(t, opaqueRoot)
	opaqueManifest := writeSourceManifest(t, opaqueRoot, opaqueRevision, opaqueLocation, "OPA-001")
	opaqueInventory := writeCarrierInventory(t, opaqueRoot, opaqueRevision, opaqueLocation, "OPAQUE")
	if out, err := runFixtureClue(t, "validate", opaqueRoot); err != nil {
		t.Fatalf("opaque target validation: %v\n%s", err, out)
	}
	if out, err := runFixtureClue(t, "parity", opaqueManifest, opaqueRoot); err != nil {
		t.Fatalf("opaque fixture parity: %v\n%s", err, out)
	}
	if out, err := runFixtureClue(t, "carriers", opaqueInventory, opaqueRoot); err != nil {
		t.Fatalf("opaque fixture carriers: %v\n%s", err, out)
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

	c, _ := corpus.Scan(numericRoot)
	imported := c.ByID["IC-901"][0]
	links := importedchange.ParseProofLinks(imported.Body)
	if imported.Fields["source-revision"] != numericRevision || imported.Fields["source-location"] != numericLocation || len(links) != 1 || links[0].Criterion != "ARC-002" || !strings.Contains(imported.Body, "## Intent") || !strings.Contains(imported.Body, "## Design rationale") || !strings.Contains(imported.Body, "## Dependencies") {
		t.Fatalf("pending source work is not durably inspectable: %+v %v", imported.Fields, links)
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

	root := fixtureTarget(t, "ARC", "ARC-002", "rev-1", "fixture")
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

	manifest := filepath.Join(root, "source-manifest.yaml")
	writeFixtureFiles(t, root, map[string]string{
		"source-manifest.yaml": "source-revision: rev-1\nsource-location: fixture\nentries:\n  - id: ARC-999\n    proof-class: Unit\n    direction: positive\n    evidence-location: fixture_test.go\n",
	})
	mustFailFixtureClue(t, parity.ClassMissingCriterion, "parity", manifest, root)
	writeFixtureFiles(t, root, map[string]string{"source-manifest.yaml": "source-revision: rev-1\nsource-location: fixture\nentries: []\n"})
	mustFailFixtureClue(t, parity.ClassOrphanedTag, "parity", manifest, root)
	writeFixtureFiles(t, root, map[string]string{"source-manifest.yaml": "source-revision: rev-1\nsource-location: fixture\nentries:\n  - id: ARC-002\n    proof-class: Integration\n    direction: negative\n    evidence-location: moved_test.go\n"})
	mustFailFixtureClue(t, parity.ClassChangedEvidence, "parity", manifest, root)
	writeFixtureFiles(t, root, map[string]string{"source-manifest.yaml": "source-revision: rev-2\nsource-location: fixture\nentries:\n  - id: ARC-002\n    proof-class: Integration\n    direction: positive\n    evidence-location: fixture_test.go\n  - id: ARC-002\n    proof-class: Integration\n    direction: negative\n    evidence-location: fixture_test.go\n"})
	mustFailFixtureClue(t, parity.ClassStaleFingerprint, "parity", manifest, root)
	criteriaPath := filepath.Join(root, "docs", "capabilities", "fixture", "criteria.md")
	criteria, err := os.ReadFile(criteriaPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(criteriaPath, []byte(strings.Replace(string(criteria), "@ARC-002", "@ARC-002 @draft", 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	mustFailFixtureClue(t, parity.ClassUnjustifiedDisposition, "parity", manifest, root)

	stalePath := filepath.Join(root, "stale-carriers.yaml")
	writeFixtureFiles(t, root, map[string]string{"stale-carriers.yaml": "source-revision: rev-1\nsource-location: fixture\ndeleted-paths:\n  - openspec/specs/numeric-archive/spec.md\nentries: []\n"})
	mustFailFixtureClue(t, carriers.ClassStaleDeletedPath, "carriers", stalePath, root)
	lostPath := filepath.Join(root, "lost-carriers.yaml")
	writeFixtureFiles(t, root, map[string]string{"lost-carriers.yaml": "source-revision: rev-1\nsource-location: fixture\nentries:\n  - id: LOST\n    kind: instruction\n    source-path: openspec/AGENTS.md\n    target-path: docs/capabilities/fixture/README.md\n    fingerprint: " + strings.Repeat("0", 64) + "\n"})
	mustFailFixtureClue(t, carriers.ClassLostFingerprint, "carriers", lostPath, root)
	missingPath := filepath.Join(root, "missing-carriers.yaml")
	writeFixtureFiles(t, root, map[string]string{"missing-carriers.yaml": "source-revision: rev-1\nsource-location: fixture\nentries:\n  - id: MISSING\n    kind: diagram-asset\n    source-path: openspec/diagram.svg\n    target-path: docs/architecture/diagram.svg\n    fingerprint: " + strings.Repeat("0", 64) + "\n"})
	mustFailFixtureClue(t, carriers.ClassMissingAsset, "carriers", missingPath, root)
}
