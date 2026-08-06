package parity

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const fixtureManifest = `source-revision: abc123
source-location: https://example.invalid/source
entries:
  - id: AC-900
    proof-class: Integration
    direction: positive
    evidence-location: fixture_test.go
  - id: AC-901
    excluded: true
    reason: not carried forward
  - id: AC-902
    disposition: draft
    justification: the source class-level tag proves nothing attributable
    disposition-source-location: legacy/spec.md:12
    plan-door: M-060
`

// reportWithRegion builds a report whose region holds body, so a test can
// commit either a rendered region or a typed one.
func reportWithRegion(manifest, body string) string {
	return "# AN-900 — fixture extraction report\n\nProse the author owns.\n\n" +
		RegionOpenPrefix + manifest + RegionOpenSuffix + "\n" +
		body +
		RegionEnd + "\n\nMore prose the author owns.\n"
}

// fixtureReportTree writes a manifest and a report carrying body, returning
// the root and the report's absolute path.
func fixtureReportTree(t *testing.T, body string) (root, report string) {
	t.Helper()
	root = writeFiles(t, map[string]string{
		"docs/analysis/AN-900-source-manifest.yaml": fixtureManifest,
		"docs/analysis/AN-900-extraction.md":        reportWithRegion("docs/analysis/AN-900-source-manifest.yaml", body),
	})
	return root, filepath.Join(root, "docs", "analysis", "AN-900-extraction.md")
}

func renderedFixture(t *testing.T, root string) string {
	t.Helper()
	m, err := LoadSourceManifest(filepath.Join(root, "docs", "analysis", "AN-900-source-manifest.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	return RenderRegion("docs/analysis/AN-900-source-manifest.yaml", m)
}

// TestAC126_UnitPositive_renderedRegionPasses proves a region rendered from
// the manifest it names passes the required check, that the rendering states
// the derived counts and one row per criterion, and that rendering twice
// produces the same bytes.
func TestAC126_UnitPositive_renderedRegionPasses(t *testing.T) {
	root := writeFiles(t, map[string]string{
		"docs/analysis/AN-900-source-manifest.yaml": fixtureManifest,
	})
	body := renderedFixture(t, root)
	if body != renderedFixture(t, root) {
		t.Fatal("expected rendering to be deterministic")
	}
	for _, want := range []string{
		"| Proven in the target | 1 |",
		"| Deferred with a plan door | 1 |",
		"| Excluded from the migration | 1 |",
		"| `AC-900` | proven | Integration · positive · fixture_test.go |",
		"| `AC-901` | excluded | not carried forward |",
		"door `M-060`",
		"source revision `abc123`",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected the rendered region to state %q, got:\n%s", want, body)
		}
	}
	root, _ = fixtureReportTree(t, body)
	if issues := CheckReports(root); len(issues) > 0 {
		t.Fatalf("expected a rendered region to pass, got %v", issues)
	}
}

// TestAC126_UnitNegative_typedFigureFails proves a hand-typed count inside
// the region fails: the report's figures cannot claim a population the
// manifest never held.
func TestAC126_UnitNegative_typedFigureFails(t *testing.T) {
	root := writeFiles(t, map[string]string{
		"docs/analysis/AN-900-source-manifest.yaml": fixtureManifest,
	})
	typed := strings.Replace(renderedFixture(t, root), "| Proven in the target | 1 |", "| Proven in the target | 417 |", 1)
	root, _ = fixtureReportTree(t, typed)
	issues := CheckReports(root)
	if len(issues) != 1 || !strings.Contains(issues[0].Msg, "disagrees with") {
		t.Fatalf("expected a typed figure to fail, got %v", issues)
	}
	if issues[0].Path != "docs/analysis/AN-900-extraction.md" {
		t.Fatalf("expected the issue to name the report, got %q", issues[0].Path)
	}
}

// TestAC126_UnitNegative_staleRegionFails proves a region left behind by a
// revised manifest fails, so a report cannot keep describing a source state
// that has been re-read.
func TestAC126_UnitNegative_staleRegionFails(t *testing.T) {
	root := writeFiles(t, map[string]string{
		"docs/analysis/AN-900-source-manifest.yaml": fixtureManifest,
	})
	body := renderedFixture(t, root)
	revised := strings.Replace(fixtureManifest, "source-revision: abc123", "source-revision: def456", 1)
	root = writeFiles(t, map[string]string{
		"docs/analysis/AN-900-source-manifest.yaml": revised,
		"docs/analysis/AN-900-extraction.md":        reportWithRegion("docs/analysis/AN-900-source-manifest.yaml", body),
	})
	if issues := CheckReports(root); len(issues) != 1 {
		t.Fatalf("expected a stale region to fail, got %v", issues)
	}
}

// TestAC126_UnitNegative_unreadableManifestFails proves a region naming a
// manifest that is not there fails rather than being silently exempt.
func TestAC126_UnitNegative_unreadableManifestFails(t *testing.T) {
	root := writeFiles(t, map[string]string{
		"docs/analysis/AN-900-extraction.md": reportWithRegion("docs/analysis/absent.yaml", "\ntyped\n"),
	})
	issues := CheckReports(root)
	if len(issues) != 1 || !strings.Contains(issues[0].Msg, "cannot read") {
		t.Fatalf("expected an absent manifest to fail, got %v", issues)
	}
}

// TestAC126_UnitNegative_malformedRegionFails proves an unterminated or
// duplicated region is reported rather than skipped, so the marker cannot be
// used to hide typed figures.
func TestAC126_UnitNegative_malformedRegionFails(t *testing.T) {
	open := RegionOpenPrefix + "docs/analysis/AN-900-source-manifest.yaml" + RegionOpenSuffix + "\n"
	cases := map[string]string{
		"unterminated": "# report\n\n" + open + "typed\n",
		"duplicated":   "# report\n\n" + open + "\n" + RegionEnd + "\n" + open + "\n" + RegionEnd + "\n",
		"nameless":     "# report\n\n" + RegionOpenPrefix + RegionOpenSuffix + "\n" + RegionEnd + "\n",
	}
	for name, content := range cases {
		root := writeFiles(t, map[string]string{
			"docs/analysis/AN-900-source-manifest.yaml": fixtureManifest,
			"docs/analysis/AN-900-extraction.md":        content,
		})
		if issues := CheckReports(root); len(issues) != 1 {
			t.Fatalf("%s: expected a malformed region to fail, got %v", name, issues)
		}
	}
}

// TestAC126_UnitPositive_documentWithoutRegionIsUnaffected proves an ordinary
// corpus document carries no obligation here: the contract reaches extraction
// reports, not every analysis record.
func TestAC126_UnitPositive_documentWithoutRegionIsUnaffected(t *testing.T) {
	root := writeFiles(t, map[string]string{
		"docs/analysis/AN-901-spike.md": "# AN-901\n\nOrdinary prose with counts: 12 criteria were read.\n",
	})
	if issues := CheckReports(root); len(issues) > 0 {
		t.Fatalf("expected a report-less document to pass, got %v", issues)
	}
}

// TestAC127_UnitPositive_regionIsRegenerated proves clue report renders the
// region from the manifest and leaves the author's prose byte-for-byte alone.
func TestAC127_UnitPositive_regionIsRegenerated(t *testing.T) {
	root, report := fixtureReportTree(t, "\n| Proven in the target | 417 |\n")
	if err := WriteReport(root, report); err != nil {
		t.Fatal(err)
	}
	after, err := os.ReadFile(report)
	if err != nil {
		t.Fatal(err)
	}
	content := string(after)
	if !strings.Contains(content, "| Proven in the target | 1 |") {
		t.Fatalf("expected the region to be rendered from the manifest, got:\n%s", content)
	}
	if strings.Contains(content, "417") {
		t.Fatal("expected the typed figure to be replaced")
	}
	for _, prose := range []string{"Prose the author owns.", "More prose the author owns."} {
		if !strings.Contains(content, prose) {
			t.Fatalf("expected prose outside the region to survive, got:\n%s", content)
		}
	}
	if issues := CheckReports(root); len(issues) > 0 {
		t.Fatalf("expected the regenerated report to pass, got %v", issues)
	}
}

// TestAC127_UnitNegative_reportWithoutRegionIsRefused proves the writer
// refuses a report that declares no region rather than inventing one: which
// manifest a report is derived from is the author's statement.
func TestAC127_UnitNegative_reportWithoutRegionIsRefused(t *testing.T) {
	root := writeFiles(t, map[string]string{
		"docs/analysis/AN-900-source-manifest.yaml": fixtureManifest,
		"docs/analysis/AN-900-extraction.md":        "# AN-900\n\nNo region here.\n",
	})
	report := filepath.Join(root, "docs", "analysis", "AN-900-extraction.md")
	before, err := os.ReadFile(report)
	if err != nil {
		t.Fatal(err)
	}
	if err := WriteReport(root, report); err == nil {
		t.Fatal("expected a report with no derived region to be refused")
	}
	after, err := os.ReadFile(report)
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Fatal("expected a refused report to be left untouched")
	}
}
