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

	stampedOnce sync.Once
	stampedPath string
	stampedErr  error

	// fixtureBuildDirs collects the temporary build directories so TestMain can
	// remove them: the builds are package-scoped by sync.Once, so no individual
	// test can own their lifetime with t.Cleanup.
	fixtureBuildMu   sync.Mutex
	fixtureBuildDirs []string
)

func TestMain(m *testing.M) {
	code := m.Run()
	fixtureBuildMu.Lock()
	for _, dir := range fixtureBuildDirs {
		os.RemoveAll(dir)
	}
	fixtureBuildMu.Unlock()
	os.Exit(code)
}

// assessmentScalePin is the release stamp the pinned-release evidence runs
// under. It belongs to the fixture rather than to this repository: the binary
// and the fixture's installed skills both carry it, so the drift comparison is
// exercised without the test churning on every real release.
const assessmentScalePin = "1.4.0"

// buildFixtureClue builds the public command once, optionally stamping the
// release version the release workflow injects (`-X main.version=...`). A
// source build reports "dev", which makes `clue validate` skip the skill-drift
// comparison altogether; only a stamped build runs it.
func buildFixtureClue(version string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	root := filepath.Clean(filepath.Join(wd, "..", ".."))
	dir, err := os.MkdirTemp("", "cliewen-migration-fixture-*")
	if err != nil {
		return "", err
	}
	fixtureBuildMu.Lock()
	fixtureBuildDirs = append(fixtureBuildDirs, dir)
	fixtureBuildMu.Unlock()
	name := "clue"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	path := filepath.Join(dir, name)
	args := []string{"build", "-o", path}
	if version != "" {
		args = append(args, "-ldflags", "-X main.version="+version)
	}
	cmd := exec.Command("go", append(args, "./cmd/clue")...)
	cmd.Dir = root
	// The stamped build passes -ldflags, so a compiler diagnostic is the first
	// thing a reader needs; a bare exit status would hide it.
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("building the fixture command (version %q): %w\n%s", version, err, out)
	}
	return path, nil
}

func fixtureClue(t *testing.T) string {
	t.Helper()
	clueOnce.Do(func() { cluePath, clueErr = buildFixtureClue("") })
	if clueErr != nil {
		t.Fatal(clueErr)
	}
	return cluePath
}

// pinnedReleaseClue is the same command built as a pinned release rather than
// from source, so a run through it is a pinned-release run (AC-129).
func pinnedReleaseClue(t *testing.T) string {
	t.Helper()
	stampedOnce.Do(func() { stampedPath, stampedErr = buildFixtureClue(assessmentScalePin) })
	if stampedErr != nil {
		t.Fatal(stampedErr)
	}
	return stampedPath
}

func runPinnedClue(t *testing.T, args ...string) (string, error) {
	t.Helper()
	out, err := exec.Command(pinnedReleaseClue(t), args...).CombinedOutput()
	return string(out), err
}

func mustFailPinnedClue(t *testing.T, want string, args ...string) {
	t.Helper()
	out, err := runPinnedClue(t, args...)
	if err == nil || !strings.Contains(out, want) {
		t.Fatalf("expected pinned-release failure containing %q, err=%v output=%s", want, err, out)
	}
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
		"docs/README.md":                        "# fixture\n\n<!-- clue:index:start -->\n- [goals/](goals/) — fixture goal\n- [capabilities/](capabilities/) — fixture capability\n- [imported-changes/](imported-changes/) — fixture imported work\n<!-- clue:index:end -->\n",
		"docs/goals/README.md":                  "# goals\n\n<!-- clue:index:start -->\n- [G-001.md](G-001.md) — fixture goal\n<!-- clue:index:end -->\n",
		"docs/goals/G-001.md":                   "---\nid: G-001\ntype: goal\nstatus: accepted\nlinks: []\ntitle: Fixture goal\n---\n\n# Fixture goal\n",
		"docs/capabilities/README.md":           "# capabilities\n\n<!-- clue:index:start -->\n- [fixture/](fixture/) — fixture capability\n<!-- clue:index:end -->\n",
		"docs/capabilities/fixture/README.md":   "---\nid: CAP-901\ntype: capability\nstatus: active\nlinks: [G-001]\ntitle: Fixture capability\ngoal: G-001\n---\n\n# Fixture capability\n",
		"docs/capabilities/fixture/criteria.md": criteria,
		"docs/imported-changes/README.md":       "# imported changes\n\n<!-- clue:index:start -->\n- [IC-901.md](IC-901.md) — fixture source work\n<!-- clue:index:end -->\n",
		"docs/imported-changes/IC-901.md":       fmt.Sprintf("---\nid: IC-901\ntype: imported-change\nstatus: complete\nlinks: [CAP-901]\ntitle: Fixture pending source work\nsource-revision: %s\nsource-location: %s\n---\n\n# Fixture pending source work\n\n## Intent\n\nPreserve the source change.\n\n## Design rationale\n\nThe target retains the proof link.\n\n## Dependencies\n\nNone.\n\n## Proof links\n\n| Task | Criterion |\n| --- | --- |\n| Preserve source proof | %s |\n", revision, location, criterion),
		"AGENTS.md":                             "# fixture routing\n",
		"tests/fixture_test.go":                 "package fixture\n\nfunc Test" + strings.ReplaceAll(criterion, "-", "") + "_IntegrationPositive_preservesProof(t *testing.T) {}\nfunc Test" + strings.ReplaceAll(criterion, "-", "") + "_IntegrationNegative_rejectsLoss(t *testing.T) {}\n",
		".clue/id-ledger.yaml":                  ledgerFile,
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

func mustValidateAfterOpaqueReservation(root string) error {
	c, scanIssues := corpus.Scan(root)
	if len(scanIssues) > 0 {
		return fmt.Errorf("fixture scan issues: %v", scanIssues)
	}
	if issues := corpus.Validate(c, corpus.Options{Version: "dev"}); len(issues) > 0 {
		return fmt.Errorf("fixture validation issues: %v", issues)
	}
	return nil
}

func writeSourceManifest(t *testing.T, root, revision, location, criterion string) string {
	t.Helper()
	path := filepath.Join(root, "source-manifest.yaml")
	content := fmt.Sprintf("source-revision: %s\nsource-location: %s\nentries:\n  - id: %s\n    proof-class: Integration\n    direction: positive\n    evidence-location: tests/fixture_test.go\n  - id: %s\n    proof-class: Integration\n    direction: negative\n    evidence-location: tests/fixture_test.go\n", revision, location, criterion, criterion)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func fixtureSource(t *testing.T, prefix, criterion, revision, location string) (string, string) {
	t.Helper()
	root := t.TempDir()
	writeFixtureFiles(t, root, map[string]string{
		"openspec/config.yaml":                         "revision: " + revision + "\n",
		"openspec/specs/fixture/spec.md":               "# Fixture source\n\n### Requirement: preserve source proof\n\n#### Scenario: fixture [" + criterion + "]\n\nTest-type: Integration\n",
		"openspec/changes/preserve-source/proposal.md": "# Preserve source work\n",
		"openspec/changes/preserve-source/design.md":   "# Source rationale\n",
		"openspec/changes/preserve-source/tasks.md":    "- [ ] Preserve " + criterion + "\n",
		"openspec/AGENTS.md":                           "# source routing\n",
		"openspec/INDEX.md":                            "# source registry\n",
		"tests/fixture_test.go":                        "package source\n\nfunc Test" + strings.ReplaceAll(criterion, "-", "") + "_IntegrationPositive_sourceProof(t *testing.T) {}\nfunc Test" + strings.ReplaceAll(criterion, "-", "") + "_IntegrationNegative_sourceLoss(t *testing.T) {}\n",
	})
	manifest := writeSourceManifest(t, root, revision, location, criterion)
	writeFixtureFiles(t, root, map[string]string{"carrier-inventory.yaml": fmt.Sprintf("source-revision: %s\nsource-location: %s\ndeleted-paths:\n  - openspec/\nentries:\n  - id: %s-INSTRUCTION\n    kind: instruction\n    source-path: openspec/AGENTS.md\n    blocked: true\n    reason: Target is not authorized during report-only rehearsal.\n  - id: %s-REGISTRY\n    kind: registry\n    source-path: openspec/INDEX.md\n    blocked: true\n    reason: Target is not authorized during report-only rehearsal.\n  - id: %s-LINK\n    kind: link\n    source-path: openspec/specs/fixture/spec.md\n    blocked: true\n    reason: Target is not authorized during report-only rehearsal.\n", revision, location, prefix, prefix, prefix)})
	return root, manifest
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

const (
	assessmentScaleCriteria = 240
	assessmentScaleArchived = 40
	assessmentScaleInFlight = 4

	// The retired criteria and the in-flight imported changes start above
	// their live neighbours. Every size-dependent ledger counter and expected
	// allocation is derived from these bases and the sizes above — the fixed
	// G and CAP identities stay literal — so changing a size keeps the ledger
	// consistent with what the fixture asserts. A size that grew into the next
	// base would instead reissue a live identity, which is what
	// assertAssessmentScaleRanges refuses to build.
	assessmentScaleArchivedBase = 400
	assessmentScaleInFlightBase = 900

	// assessmentScaleDrafted is the in-flight record whose proof-linked
	// criterion the negative run drafts. The fixture proof-links IC-<base+i>
	// to SCL-<i>, so the pair is derived from this one number.
	assessmentScaleDrafted = 2
)

// assertAssessmentScaleRanges refuses to build a fixture whose size constants
// have grown into the next identity band: the ledger would then carry the same
// SCL component twice, once live and once retired, and the fixture's own
// contradiction would surface as an unrelated validate failure.
func assertAssessmentScaleRanges(t *testing.T) {
	t.Helper()
	if assessmentScaleCriteria >= assessmentScaleArchivedBase {
		t.Fatalf("assessmentScaleCriteria (%d) reaches the retired band at %d", assessmentScaleCriteria, assessmentScaleArchivedBase)
	}
	if assessmentScaleArchivedBase+assessmentScaleArchived >= assessmentScaleInFlightBase {
		t.Fatalf("the retired band (%d..%d) reaches the in-flight band at %d", assessmentScaleArchivedBase+1, assessmentScaleArchivedBase+assessmentScaleArchived, assessmentScaleInFlightBase)
	}
	if assessmentScaleInFlight < assessmentScaleDrafted {
		t.Fatalf("assessmentScaleInFlight (%d) does not reach the drafted record at offset %d", assessmentScaleInFlight, assessmentScaleDrafted)
	}
}

// assessmentScaleTarget is deliberately generated in a temporary directory:
// its size is evidence for the composed contract, not a second corpus to
// maintain in the repository. It composes the two phases the ordered path
// keeps apart — the rehearsal writes the source pins first, and the approved
// mutation materializes the target afterwards.
func assessmentScaleTarget(t *testing.T, revision, location string) (string, string) {
	t.Helper()
	root := t.TempDir()
	manifest := assessmentScaleSource(t, revision, location)
	materializeAssessmentScaleTarget(t, root, revision, location)
	return root, manifest
}

// assessmentScaleSource writes only what the report-only rehearsal produces:
// the source tree, its pinned source manifest, and a carrier inventory whose
// entries are all blocked because no target exists yet. It returns the
// manifest path.
func assessmentScaleSource(t *testing.T, revision, location string) string {
	t.Helper()
	var manifest strings.Builder
	fmt.Fprintf(&manifest, "source-revision: %s\nsource-location: %s\nentries:\n", revision, location)
	for i := 1; i <= assessmentScaleCriteria; i++ {
		id := fmt.Sprintf("SCL-%03d", i)
		fmt.Fprintf(&manifest, "  - id: %s\n    proof-class: Integration\n    direction: positive\n    evidence-location: tests/scale_test.go\n  - id: %s\n    proof-class: Integration\n    direction: negative\n    evidence-location: tests/scale_test.go\n", id, id)
	}
	sourceRoot := t.TempDir()
	writeFixtureFiles(t, sourceRoot, map[string]string{
		"source-manifest.yaml":   manifest.String(),
		"carrier-inventory.yaml": fmt.Sprintf("source-revision: %s\nsource-location: %s\ndeleted-paths:\n  - openspec/\nentries:\n  - id: SCALE-INSTRUCTION\n    kind: instruction\n    source-path: openspec/AGENTS.md\n    blocked: true\n    reason: Target is not authorized during report-only rehearsal.\n", revision, location),
		"openspec/AGENTS.md":     "# source routing\n",
		"tests/scale_test.go":    "package source\n",
	})
	return filepath.Join(sourceRoot, "source-manifest.yaml")
}

// materializeAssessmentScaleTarget is the mutate phase: it writes the target
// corpus the rehearsal's pins are later verified against.
func materializeAssessmentScaleTarget(t *testing.T, root, revision, location string) {
	t.Helper()
	assertAssessmentScaleRanges(t)
	var criteria, tests, ledgerFile, importedIndex strings.Builder
	fmt.Fprint(&criteria, "---\nid: SCL-criteria\ntype: criteria\nstatus: active\nlinks: [CAP-901]\ntitle: Assessment-scale fixture criteria\nac-prefix: SCL\n---\n\n```gherkin\nFeature: Assessment-scale migration\n")
	fmt.Fprintf(&ledgerFile, "counters:\n  G: 1\n  CAP: 901\n  IC: %d\n  SCL: %d\nentries:\n  - id: G-001\n    kind: numeric\n    state: live\n    prefix: G\n    component: 1\n  - id: CAP-901\n    kind: numeric\n    state: live\n    prefix: CAP\n    component: 901\n  - id: SCL-criteria\n    kind: opaque\n    state: live\n", assessmentScaleInFlightBase+assessmentScaleInFlight, assessmentScaleArchivedBase+assessmentScaleArchived)
	fmt.Fprint(&importedIndex, "# imported changes\n\n<!-- clue:index:start -->\n")
	for i := 1; i <= assessmentScaleCriteria; i++ {
		id := fmt.Sprintf("SCL-%03d", i)
		fmt.Fprintf(&criteria, "\n  @%s\n  Scenario: source proof %d survives migration\n    Test-type: Integration\n    Given source proof %d\n    When it is migrated\n    Then it remains classified\n", id, i, i)
		fmt.Fprintf(&tests, "func TestSCL%03d_IntegrationPositive_preservesProof(t *testing.T) {}\nfunc TestSCL%03d_IntegrationNegative_rejectsLoss(t *testing.T) {}\n", i, i)
		fmt.Fprintf(&ledgerFile, "  - id: %s\n    kind: numeric\n    state: live\n    prefix: SCL\n    component: %d\n    source-revision: %s\n    source-location: %s\n", id, i, revision, location)
	}
	fmt.Fprint(&criteria, "```\n")
	for i := 1; i <= assessmentScaleArchived; i++ {
		component := assessmentScaleArchivedBase + i
		fmt.Fprintf(&ledgerFile, "  - id: SCL-%03d\n    kind: numeric\n    state: retired\n    prefix: SCL\n    component: %d\n    source-revision: %s\n    source-location: %s\n", component, component, revision, location)
	}
	for i := 1; i <= assessmentScaleInFlight; i++ {
		id := fmt.Sprintf("IC-%03d", assessmentScaleInFlightBase+i)
		fmt.Fprintf(&ledgerFile, "  - id: %s\n    kind: numeric\n    state: live\n    prefix: IC\n    component: %d\n", id, assessmentScaleInFlightBase+i)
		fmt.Fprintf(&importedIndex, "- [%s.md](%s.md) — fixture source work\n", id, id)
	}
	writeFixtureFiles(t, root, map[string]string{
		"docs/README.md":                        "# fixture\n\n<!-- clue:index:start -->\n- [goals/](goals/) — fixture goal\n- [capabilities/](capabilities/) — fixture capability\n- [imported-changes/](imported-changes/) — fixture imported work\n<!-- clue:index:end -->\n",
		"docs/goals/README.md":                  "# goals\n\n<!-- clue:index:start -->\n- [G-001.md](G-001.md) — fixture goal\n<!-- clue:index:end -->\n",
		"docs/goals/G-001.md":                   "---\nid: G-001\ntype: goal\nstatus: accepted\nlinks: []\ntitle: Fixture goal\n---\n\n# Fixture goal\n",
		"docs/capabilities/README.md":           "# capabilities\n\n<!-- clue:index:start -->\n- [fixture/](fixture/) — fixture capability\n<!-- clue:index:end -->\n",
		"docs/capabilities/fixture/README.md":   "---\nid: CAP-901\ntype: capability\nstatus: active\nlinks: [G-001]\ntitle: Fixture capability\ngoal: G-001\n---\n\n# Fixture capability\n",
		"docs/capabilities/fixture/criteria.md": criteria.String(),
		"docs/imported-changes/README.md":       importedIndex.String() + "<!-- clue:index:end -->\n",
		"tests/scale_test.go":                   "package fixture\n\nimport \"testing\"\n\n" + tests.String(),
		".clue/id-ledger.yaml":                  ledgerFile.String(),
		"AGENTS.md":                             "# assessment-scale fixture routing\n",
		".github/workflows/validate.yml":        "name: validate\non: [push]\njobs:\n  validate:\n    runs-on: ubuntu-latest\n    steps:\n      - run: clue validate .\n",
	})
	for i := 1; i <= assessmentScaleInFlight; i++ {
		id := fmt.Sprintf("IC-%03d", assessmentScaleInFlightBase+i)
		writeFixtureFiles(t, root, map[string]string{fmt.Sprintf("docs/imported-changes/%s.md", id): fmt.Sprintf("---\nid: %s\ntype: imported-change\nstatus: in-progress\nlinks: [CAP-901]\ntitle: Assessment-scale in-flight source work %d\nsource-revision: %s\nsource-location: %s\n---\n\n# Assessment-scale in-flight source work %d\n\n## Intent\n\nPreserve concurrent source work.\n\n## Design rationale\n\nThe target keeps its proof link inspectable.\n\n## Dependencies\n\nNone.\n\n## Proof links\n\n| Task | Criterion |\n| --- | --- |\n| Preserve source proof | SCL-%03d |\n", id, i, revision, location, i, i)})
	}
	writeFixtureSkills(t, root, assessmentScalePin)
}

// writeFixtureSkills installs the managed skills a migrated repository
// carries, stamped with the release the run is pinned to. The stamp is what
// `clue validate` compares a released binary against; a source build has no
// release to drift from and skips the comparison entirely.
func writeFixtureSkills(t *testing.T, root, version string) {
	t.Helper()
	files := map[string]string{}
	for _, name := range []string{"clue-delta", "clue-extract", "clue-verify"} {
		files[".agents/skills/"+name+"/skill.md"] = fmt.Sprintf("---\ncliewen-skill: true\nversion: %s\n---\n\n# %s\n", version, name)
	}
	writeFixtureFiles(t, root, files)
}

// assessmentScaleInventory pins one carrier per kind the target actually
// holds, so reconciliation at scale exercises the same kind vocabulary the
// smaller AC-123 inventory does rather than labelling every row
// `instruction`.
func assessmentScaleInventory(t *testing.T, root, revision, location string) string {
	t.Helper()
	paths := []string{"AGENTS.md", ".github/workflows/validate.yml", "docs/capabilities/README.md", "docs/imported-changes/README.md", "docs/capabilities/fixture/README.md"}
	kinds := []carriers.Kind{carriers.KindInstruction, carriers.KindWorkflow, carriers.KindRegistry, carriers.KindRegistry, carriers.KindLink}
	var entries strings.Builder
	for i, target := range paths {
		fp, err := carriers.Fingerprint(filepath.Join(root, filepath.FromSlash(target)))
		if err != nil {
			t.Fatal(err)
		}
		fmt.Fprintf(&entries, "  - id: SCALE-%d\n    kind: %s\n    source-path: openspec/carrier-%d.md\n    target-path: %s\n    fingerprint: %s\n", i+1, kinds[i], i+1, target, fp)
	}
	path := filepath.Join(root, "carrier-inventory.yaml")
	writeFixtureFiles(t, root, map[string]string{"carrier-inventory.yaml": fmt.Sprintf("source-revision: %s\nsource-location: %s\ndeleted-paths:\n  - openspec/\nentries:\n%s", revision, location, entries.String())})
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

// TestAC128_UnitPositive_assessmentScaleFixtureProvesComposedMigrationContract
// holds the complete contract at the originating assessment's order of
// magnitude. The temporary source is never executed as a test suite.
func TestAC128_UnitPositive_assessmentScaleFixtureProvesComposedMigrationContract(t *testing.T) {
	revision, location := "assessment-scale-fixture-v1", "fixtures/assessment-scale"
	targetRoot, manifest := assessmentScaleTarget(t, revision, location)
	sourceRoot := filepath.Dir(manifest)
	// The rehearsal runs before the target is authorized, so the only entry
	// shape it can carry is `blocked`. Pairing the clean run with a premature
	// mapping keeps the step from passing vacuously: an inventory that claims
	// a target the rehearsal has not written yet must fail.
	if out, err := runFixtureClue(t, "carriers", filepath.Join(sourceRoot, "carrier-inventory.yaml"), sourceRoot); err != nil {
		t.Fatalf("report-only carrier rehearsal: %v\n%s", err, out)
	}
	prematurePath := filepath.Join(sourceRoot, "premature-carriers.yaml")
	writeFixtureFiles(t, sourceRoot, map[string]string{filepath.Base(prematurePath): "source-revision: " + revision + "\nsource-location: " + location + "\nentries:\n  - id: SCALE-INSTRUCTION\n    kind: instruction\n    source-path: openspec/AGENTS.md\n    target-path: AGENTS.md\n    fingerprint: " + strings.Repeat("0", 64) + "\n"})
	mustFailFixtureClue(t, carriers.ClassMissingAsset, "carriers", prematurePath, sourceRoot)
	mustValidateFixture(t, targetRoot)
	if out, err := runFixtureClue(t, "validate", targetRoot); err != nil {
		t.Fatalf("assessment-scale target validation: %v\n%s", err, out)
	}
	if out, err := runFixtureClue(t, "parity", manifest, targetRoot); err != nil {
		t.Fatalf("assessment-scale parity: %v\n%s", err, out)
	}
	inventory := assessmentScaleInventory(t, targetRoot, revision, location)
	if out, err := runFixtureClue(t, "carriers", inventory, targetRoot); err != nil {
		t.Fatalf("assessment-scale carrier reconciliation: %v\n%s", err, out)
	}

	// Allocation is asked on a throwaway instance: NextNumeric reserves the
	// ID it returns and advances the stored counter, so asking it on the
	// instance this test later saves would persist four reservations the
	// fixture never intended and make the assertion unrepeatable.
	allocation, err := ledger.Load(targetRoot)
	if err != nil {
		t.Fatal(err)
	}
	nextAllocation := map[string]string{
		"G":   "G-002",
		"CAP": "CAP-902",
		"IC":  fmt.Sprintf("IC-%03d", assessmentScaleInFlightBase+assessmentScaleInFlight+1),
		"SCL": fmt.Sprintf("SCL-%03d", assessmentScaleArchivedBase+assessmentScaleArchived+1),
	}
	for prefix, want := range nextAllocation {
		got, err := allocation.NextNumeric(prefix)
		if err != nil || got != want {
			t.Fatalf("next %s = %q, %v; want %s", prefix, got, err, want)
		}
	}

	l, err := ledger.Load(targetRoot)
	if err != nil {
		t.Fatal(err)
	}
	// Reservation must refuse an identity the ledger already carries in any
	// state, not only one reserved a moment ago in memory: an existing live
	// opaque ID and one of the forty retired numeric IDs are both refused.
	if err := l.ReserveOpaque("SCL-criteria", revision, location); err == nil {
		t.Fatal("existing live opaque identity was reissued")
	}
	if err := l.ReserveOpaque(fmt.Sprintf("SCL-%03d", assessmentScaleArchivedBase+1), revision, location); err == nil {
		t.Fatal("retired identity was reissued")
	}
	if err := l.ReserveOpaque("assessment-scale-source-opaque-id", revision, location); err != nil {
		t.Fatal(err)
	}
	if err := l.ReserveOpaque("assessment-scale-source-opaque-id", revision, location); err == nil {
		t.Fatal("opaque source identity was reused")
	}
	if err := l.Save(); err != nil {
		t.Fatal(err)
	}
	// The reservation is permanent because it is on disk, not because it is
	// in this process: reload and re-attempt it.
	reloaded, err := ledger.Load(targetRoot)
	if err != nil {
		t.Fatal(err)
	}
	if err := reloaded.ReserveOpaque("assessment-scale-source-opaque-id", revision, location); err == nil {
		t.Fatal("opaque source identity was reissued after a ledger round trip")
	}
	if err := mustValidateAfterOpaqueReservation(targetRoot); err != nil {
		t.Fatal(err)
	}
}

// TestAC128_UnitNegative_assessmentScaleFixtureRejectsEveryFailureClass
// runs the public failure paths against the same order-of-magnitude shape,
// rather than inferring scale behaviour from the smaller AC-123 fixture.
func TestAC128_UnitNegative_assessmentScaleFixtureRejectsEveryFailureClass(t *testing.T) {
	revision, location := "assessment-scale-fixture-v1", "fixtures/assessment-scale"
	targetRoot, manifest := assessmentScaleTarget(t, revision, location)
	pristine, err := os.ReadFile(manifest)
	if err != nil {
		t.Fatal(err)
	}
	writeManifest := func(content string) {
		t.Helper()
		if err := os.WriteFile(manifest, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	writeManifest(strings.ReplaceAll(string(pristine), "SCL-001", "SCL-999"))
	mustFailFixtureClue(t, parity.ClassMissingCriterion, "parity", manifest, targetRoot)
	writeManifest("source-revision: " + revision + "\nsource-location: " + location + "\nentries: []\n")
	mustFailFixtureClue(t, parity.ClassOrphanedTag, "parity", manifest, targetRoot)
	writeManifest(strings.ReplaceAll(string(pristine), "tests/scale_test.go", "tests/moved_test.go"))
	mustFailFixtureClue(t, parity.ClassChangedEvidence, "parity", manifest, targetRoot)
	writeManifest(strings.Replace(string(pristine), revision, "assessment-scale-fixture-v2", 1))
	mustFailFixtureClue(t, parity.ClassStaleFingerprint, "parity", manifest, targetRoot)
	writeManifest(string(pristine))
	criteriaPath := filepath.Join(targetRoot, "docs", "capabilities", "fixture", "criteria.md")
	criteria, err := os.ReadFile(criteriaPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(criteriaPath, []byte(strings.Replace(string(criteria), "@SCL-001", "@SCL-001 @draft", 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	mustFailFixtureClue(t, parity.ClassUnjustifiedDisposition, "parity", manifest, targetRoot)
	// A deferral is unaccountable when the target corpus declares no plan
	// door at all, which is this fixture's shape (ADR-053). The criterion
	// still carries the @draft tag the step above added, so the disposition
	// matches the target and the unjustified class must stay silent — that
	// absence is asserted below, because otherwise this step would keep
	// passing on the wrong class if the tag were ever restored.
	provenProof := `  - id: SCL-001
    proof-class: Integration
    direction: positive
    evidence-location: tests/scale_test.go
  - id: SCL-001
    proof-class: Integration
    direction: negative
    evidence-location: tests/scale_test.go
`
	deferredToAbsentDoor := `  - id: SCL-001
    disposition: draft
    justification: The source reference is not re-derived at this scale yet.
    disposition-source-location: openspec/specs/fixture/spec.md:1
    plan-door: M-999
`
	deferred := strings.Replace(string(pristine), provenProof, deferredToAbsentDoor, 1)
	if deferred == string(pristine) {
		t.Fatal("deferral rewrite matched nothing; the manifest shape changed")
	}
	writeManifest(deferred)
	out, err := runFixtureClue(t, "parity", manifest, targetRoot)
	if err == nil || !strings.Contains(out, parity.ClassUnaccountableDisposition) || strings.Contains(out, parity.ClassUnjustifiedDisposition) {
		t.Fatalf("expected only %q, err=%v output=%s", parity.ClassUnaccountableDisposition, err, out)
	}

	stalePath := filepath.Join(targetRoot, "stale-carriers.yaml")
	staleTarget := filepath.Join(targetRoot, "docs", "capabilities", "fixture", "README.md")
	staleBody, err := os.ReadFile(staleTarget)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(staleTarget, append(staleBody, []byte("\nSee [old](../../../openspec/specs/scale/spec.md).\n")...), 0o644); err != nil {
		t.Fatal(err)
	}
	writeFixtureFiles(t, targetRoot, map[string]string{filepath.Base(stalePath): "source-revision: " + revision + "\nsource-location: " + location + "\ndeleted-paths:\n  - openspec/specs/scale/spec.md\nentries: []\n"})
	mustFailFixtureClue(t, carriers.ClassStaleDeletedPath, "carriers", stalePath, targetRoot)

	lostPath := filepath.Join(targetRoot, "lost-carriers.yaml")
	missingPath := filepath.Join(targetRoot, "missing-carriers.yaml")
	writeFixtureFiles(t, targetRoot, map[string]string{
		filepath.Base(lostPath):    "source-revision: " + revision + "\nsource-location: " + location + "\nentries:\n  - id: LOST\n    kind: instruction\n    source-path: openspec/AGENTS.md\n    target-path: docs/capabilities/fixture/README.md\n    fingerprint: " + strings.Repeat("0", 64) + "\n",
		filepath.Base(missingPath): "source-revision: " + revision + "\nsource-location: " + location + "\nentries:\n  - id: MISSING\n    kind: diagram-asset\n    source-path: openspec/diagram.svg\n    target-path: docs/architecture/diagram.svg\n    fingerprint: " + strings.Repeat("0", 64) + "\n",
	})
	mustFailFixtureClue(t, carriers.ClassLostFingerprint, "carriers", lostPath, targetRoot)
	mustFailFixtureClue(t, carriers.ClassMissingAsset, "carriers", missingPath, targetRoot)
}

// TestAC129_UnitPositive_orderedPinnedReleasePathHoldsAtAssessmentScale runs
// the migration path in the order the extraction contract requires — rehearsal
// pins written while no target exists, then the approved mutation verified
// against those same unmodified pins — through a binary built as a pinned
// release rather than from source. The source fixture's own suite is never
// executed, and no production adopter is involved.
func TestAC129_UnitPositive_orderedPinnedReleasePathHoldsAtAssessmentScale(t *testing.T) {
	revision, location := "assessment-scale-fixture-v1", "fixtures/assessment-scale"
	if out, err := runPinnedClue(t, "version"); err != nil || !strings.Contains(out, assessmentScalePin) {
		t.Fatalf("pinned-release build reports %q, %v; want %s", out, err, assessmentScalePin)
	}

	// Rehearsal: the source pins exist and the target does not. A parity or
	// carrier verdict claimed here would be a claim about a corpus nobody has
	// written, so both must fail rather than pass vacuously.
	manifest := assessmentScaleSource(t, revision, location)
	sourceRoot := filepath.Dir(manifest)
	targetRoot := t.TempDir()
	// The count is asserted alongside the class so a partially written target
	// cannot satisfy this step: every criterion the source declared is
	// missing, which is what "the target has not been written yet" means.
	out, err := runPinnedClue(t, "parity", manifest, targetRoot)
	if err == nil || strings.Count(out, parity.ClassMissingCriterion) != assessmentScaleCriteria {
		t.Fatalf("parity before mutation: err=%v, %d %s findings, want %d\n%s", err, strings.Count(out, parity.ClassMissingCriterion), parity.ClassMissingCriterion, assessmentScaleCriteria, out)
	}
	// The rehearsal's own inventory is clean, because every entry it carries
	// is blocked: that is the only shape a rehearsal can produce, and it is
	// what keeps this step from being a broken-inventory result. An inventory
	// that instead maps a target — the shape only the mutate phase may write —
	// fails here and passes only after mutation.
	if out, err := runPinnedClue(t, "carriers", filepath.Join(sourceRoot, "carrier-inventory.yaml"), sourceRoot); err != nil {
		t.Fatalf("report-only carrier rehearsal: %v\n%s", err, out)
	}
	prematurePath := filepath.Join(sourceRoot, "premature-carriers.yaml")
	writeFixtureFiles(t, sourceRoot, map[string]string{filepath.Base(prematurePath): "source-revision: " + revision + "\nsource-location: " + location + "\nentries:\n  - id: SCALE-INSTRUCTION\n    kind: instruction\n    source-path: openspec/AGENTS.md\n    target-path: AGENTS.md\n    fingerprint: " + strings.Repeat("0", 64) + "\n"})
	mustFailPinnedClue(t, carriers.ClassMissingAsset, "carriers", prematurePath, targetRoot)
	pins, err := os.ReadFile(manifest)
	if err != nil {
		t.Fatal(err)
	}

	// Mutation, which a human authorizes; the merge of this change is that
	// acceptance for the fixture.
	materializeAssessmentScaleTarget(t, targetRoot, revision, location)
	if out, err := runPinnedClue(t, "validate", targetRoot); err != nil {
		t.Fatalf("pinned-release target validation: %v\n%s", err, out)
	}
	// The manifest that proves the migration is the one the rehearsal wrote,
	// byte for byte — a manifest rewritten to match the target afterwards
	// would prove nothing about what the source held. This is a regression
	// guard rather than a demonstration: the mutation writes only under the
	// target root today, so the check holds by construction and exists to fail
	// if a later change lets the mutate phase reach back into the pins.
	after, err := os.ReadFile(manifest)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != string(pins) {
		t.Fatal("the rehearsal's pinned manifest changed during mutation")
	}
	if out, err := runPinnedClue(t, "parity", manifest, targetRoot); err != nil {
		t.Fatalf("parity against the rehearsal's own pins: %v\n%s", err, out)
	}
	inventory := assessmentScaleInventory(t, targetRoot, revision, location)
	if out, err := runPinnedClue(t, "carriers", inventory, targetRoot); err != nil {
		t.Fatalf("carrier reconciliation after mutation: %v\n%s", err, out)
	}

	// Source-work preservation: every in-flight record the migration imported
	// stays readable on its own terms after the source is gone.
	c, scanIssues := corpus.Scan(targetRoot)
	if len(scanIssues) > 0 {
		t.Fatalf("target scan issues: %v", scanIssues)
	}
	declared, _ := corpus.AcceptanceEvidence(c)
	for i := 1; i <= assessmentScaleInFlight; i++ {
		id := fmt.Sprintf("IC-%03d", assessmentScaleInFlightBase+i)
		records := c.ByID[id]
		if len(records) != 1 {
			t.Fatalf("imported change %s: got %d records, want 1", id, len(records))
		}
		record := records[0]
		if record.Fields["source-revision"] != revision || record.Fields["source-location"] != location {
			t.Fatalf("imported change %s lost its pinned origin: %+v", id, record.Fields)
		}
		for _, section := range []string{"## Intent", "## Design rationale", "## Dependencies"} {
			if !strings.Contains(record.Body, section) {
				t.Fatalf("imported change %s has no %q section", id, section)
			}
		}
		links := importedchange.ParseProofLinks(record.Body)
		if len(links) != 1 || links[0].Task == "" {
			t.Fatalf("imported change %s has no inspectable proof link: %v", id, links)
		}
		d, ok := declared[links[0].Criterion]
		if !ok || d.Draft || d.Retired {
			t.Fatalf("imported change %s proof-links %s, which the target cannot back", id, links[0].Criterion)
		}
	}

	// A record may close only once the work it names is proven, and at this
	// size it is: promoting one to complete keeps the corpus green.
	completePath := filepath.Join(targetRoot, "docs", "imported-changes", fmt.Sprintf("IC-%03d.md", assessmentScaleInFlightBase+1))
	record, err := os.ReadFile(completePath)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(completePath, []byte(strings.Replace(string(record), "status: in-progress", "status: complete", 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	if out, err := runPinnedClue(t, "validate", targetRoot); err != nil {
		t.Fatalf("completing a proven imported change: %v\n%s", err, out)
	}
}

// TestAC129_UnitNegative_orderedPinnedReleasePathRejectsItsViolations proves
// each half of the positive path is load-bearing: the release stamp, the
// ordering, and the imported-change closing gate.
func TestAC129_UnitNegative_orderedPinnedReleasePathRejectsItsViolations(t *testing.T) {
	revision, location := "assessment-scale-fixture-v1", "fixtures/assessment-scale"
	targetRoot, manifest := assessmentScaleTarget(t, revision, location)

	// Drift: skills stamped with a release other than the binary's. Without
	// this the positive run's pinned stamp would be decoration, since a
	// source build skips the comparison rather than passing it.
	writeFixtureSkills(t, targetRoot, "1.3.0")
	mustFailPinnedClue(t, "drift", "validate", targetRoot)
	writeFixtureSkills(t, targetRoot, assessmentScalePin)
	if out, err := runPinnedClue(t, "validate", targetRoot); err != nil {
		t.Fatalf("restored pinned skills: %v\n%s", err, out)
	}

	// Ordering: the same pins that pass against the mutated target give no
	// verdict against a target the mutation has not written.
	unwritten := t.TempDir()
	mustFailPinnedClue(t, parity.ClassMissingCriterion, "parity", manifest, unwritten)
	mustFailPinnedClue(t, carriers.ClassMissingAsset, "carriers", assessmentScaleInventory(t, targetRoot, revision, location), unwritten)

	// Source work: a record may not declare itself complete over work the
	// target cannot back, whatever the corpus's size.
	criteriaPath := filepath.Join(targetRoot, "docs", "capabilities", "fixture", "criteria.md")
	criteria, err := os.ReadFile(criteriaPath)
	if err != nil {
		t.Fatal(err)
	}
	// The drafted criterion and the record that proof-links it are the same
	// pair the fixture wrote, so this step stays a proof-backing violation
	// rather than a missing-file failure when the sizes change.
	draftedTag := fmt.Sprintf("@SCL-%03d", assessmentScaleDrafted)
	if err := os.WriteFile(criteriaPath, []byte(strings.Replace(string(criteria), draftedTag, draftedTag+" @draft", 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	completePath := filepath.Join(targetRoot, "docs", "imported-changes", fmt.Sprintf("IC-%03d.md", assessmentScaleInFlightBase+assessmentScaleDrafted))
	record, err := os.ReadFile(completePath)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(completePath, []byte(strings.Replace(string(record), "status: in-progress", "status: complete", 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	mustFailPinnedClue(t, "imported-change is complete", "validate", targetRoot)
}

// TestAC123_UnitPositive_disposableFixturesProveComposedMigrationContract
// exercises the approved mutation phase against disposable source shapes.
// It deliberately never executes the source fixture's own test suite.
func TestAC123_UnitPositive_disposableFixturesProveComposedMigrationContract(t *testing.T) {
	numericRevision, numericLocation := "numeric-archive-fixture-v1", "fixtures/numeric-archive"
	numericSourceRoot, numericManifest := fixtureSource(t, "NUMERIC", "ARC-002", numericRevision, numericLocation)
	numericRoot := fixtureTarget(t, "ARC", "ARC-002", numericRevision, numericLocation)
	mustValidateFixture(t, numericRoot)
	numericInventory := writeCarrierInventory(t, numericRoot, numericRevision, numericLocation, "NUMERIC")
	if _, err := os.Stat(filepath.Join(numericSourceRoot, "tests", "fixture_test.go")); err != nil {
		t.Fatal(err)
	}
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
	opaqueSourceRoot, opaqueManifest := fixtureSource(t, "OPAQUE", "OPA-001", opaqueRevision, opaqueLocation)
	opaqueRoot := fixtureTarget(t, "OPA", "OPA-001", opaqueRevision, opaqueLocation)
	mustValidateFixture(t, opaqueRoot)
	opaqueInventory := writeCarrierInventory(t, opaqueRoot, opaqueRevision, opaqueLocation, "OPAQUE")
	if _, err := os.Stat(filepath.Join(opaqueSourceRoot, "tests", "fixture_test.go")); err != nil {
		t.Fatal(err)
	}
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
	if err := opaqueLedger.Save(); err != nil {
		t.Fatal(err)
	}
	if err := mustValidateAfterOpaqueReservation(opaqueRoot); err != nil {
		t.Fatal(err)
	}
	opaqueLedger, err = ledger.Load(opaqueRoot)
	if err != nil {
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
