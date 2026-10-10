package corpus

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func readableFixture(t *testing.T, extra map[string]string) *Corpus {
	t.Helper()
	root := t.TempDir()
	files := map[string]string{
		"docs/goals/G-901.md":                   "---\nid: G-901\ntype: goal\nstatus: accepted\nlinks: []\ntitle: Understand backup status\n---\n\nA reader understands the status.\n",
		"docs/capabilities/CAP-901/README.md":   "---\nid: CAP-901\ntype: capability\nstatus: active\nlinks: [G-901]\ngoal: G-901\ntitle: Backup status\n---\n",
		"docs/capabilities/CAP-901/criteria.md": "---\nid: CAP-901-criteria\ntype: criteria\nstatus: active\nlinks: [CAP-901]\ntitle: Backup proof\n---\n\n```gherkin\n@AC-901\nScenario: A person sees the failed backup\n  Test-type: Human\n  Given a failed simulated backup\n  Then the person can explain the failure\n```\n",
		"docs/plans/P-901.md":                   "---\nid: P-901\ntype: plan\nstatus: active\nlinks: [G-901]\ntitle: Improve backup clarity\n---\n\n| ID | Milestone | Status |\n|---|---|---|\n| M-901 | Show the failure clearly | todo |\n",
	}
	for path, body := range extra {
		files[path] = body
	}
	for path, body := range files {
		full := filepath.Join(root, filepath.FromSlash(path))
		if e := os.MkdirAll(filepath.Dir(full), 0755); e != nil {
			t.Fatal(e)
		}
		if e := os.WriteFile(full, []byte(body), 0644); e != nil {
			t.Fatal(e)
		}
	}
	c, issues := Scan(root)
	if len(issues) != 0 {
		t.Fatalf("scan: %v", issues)
	}
	return c
}

func TestAC233_UnitPositive_NamesUseExistingArtifactScenarioAndMilestoneDeclarations(t *testing.T) {
	c := readableFixture(t, nil)
	before := c.Contents["docs/capabilities/CAP-901/criteria.md"]
	names := NewReferenceNames(c)
	for id, want := range map[string]string{"G-901": "Understand backup status (G-901)", "CAP-901": "Backup status (CAP-901)", "AC-901": "A person sees the failed backup (AC-901)", "M-901": "Show the failure clearly (M-901)"} {
		if got := names.Label(id); got != want {
			t.Fatalf("%s: %q", id, got)
		}
	}
	coverage := Coverage(c)
	if len(coverage) != 1 || coverage[0].Capability != "CAP-901" || coverage[0].State != "covered" {
		t.Fatalf("raw coverage changed: %v", coverage)
	}
	if c.Contents["docs/capabilities/CAP-901/criteria.md"] != before || c.ByID["G-901"][0].ID != "G-901" || c.localIssueReferences != nil {
		t.Fatal("naming changed raw corpus or diagnostic metadata")
	}
	if got := NewReferenceNames(nil).Label("G-901"); got != "unknown reference (G-901)" {
		t.Fatal(got)
	}
}

func TestAC233_UnitNegative_UnknownDuplicateAndUnsafeNamesAreNotGuessed(t *testing.T) {
	c := readableFixture(t, map[string]string{
		"docs/goals/G-901-duplicate.md":         "---\nid: G-901\ntype: goal\nstatus: accepted\nlinks: []\ntitle: Other goal\n---\n",
		"docs/goals/G-902.md":                   "---\nid: G-902\ntype: goal\nstatus: accepted\nlinks: []\ntitle: ''\n---\n",
		"docs/capabilities/CAP-902/criteria.md": "---\nid: CAP-902-criteria\ntype: criteria\nstatus: active\nlinks: []\ntitle: Duplicate proof\n---\n\n```gherkin\n@AC-901\nScenario: Other scenario\n  Test-type: Human\n```\n",
	})
	names := NewReferenceNames(c)
	for id, want := range map[string]string{"G-901": "ambiguous reference (G-901)", "G-902": "unnamed reference (G-902)", "G-999": "unknown reference (G-999)", "AC-901": "ambiguous reference (AC-901)", "clue:other/repo/G-901": "external reference, not resolved locally (clue:other/repo/G-901)"} {
		if got := names.Label(id); got != want {
			t.Fatalf("%s: %q", id, got)
		}
	}
	if got := HumanReference("Reader\n\tname\x1b[2J", "G-903"); got != "Reader name [2J (G-903)" {
		t.Fatalf("unsafe label %q", got)
	}
}

func TestAC234_UnitPositive_ProducingChecksNameTheirLocalDiagnosticSubjects(t *testing.T) {
	c := readableFixture(t, nil)
	criterion := c.ByID["CAP-901-criteria"][0]
	criterion.Body = strings.Replace(criterion.Body, "Test-type: Human", "Test-type: Unit", 1)
	issues := checkACTests(c)
	names := NewReferenceNames(c)
	found := false
	for _, issue := range issues {
		if !strings.Contains(issue.Msg, "AC-901 has no Unit positive") {
			continue
		}
		raw := issue.String()
		human := names.Issue(issue)
		if !strings.Contains(human, "A person sees the failed backup (AC-901) has no Unit positive") {
			t.Fatal(human)
		}
		if issue.String() != raw {
			t.Fatal("raw issue changed")
		}
		found = true
	}
	if !found {
		t.Fatal("missing proof finding")
	}
	link := localIssue(c, "docs/goals/G-901.md", "link G-999 resolves to no artifact", "G-999")
	if got := NewReferenceNames(c).Issue(link); !strings.Contains(got, "unknown reference (G-999)") {
		t.Fatal(got)
	}
	if raw := fmt.Sprint(link); raw != "docs/goals/G-901.md: link G-999 resolves to no artifact" {
		t.Fatal(raw)
	}
}

func TestAC234_UnitNegative_SourceRulesAndInsertedNameTokensAreNotResolvedAsLocalSubjects(t *testing.T) {
	c := readableFixture(t, map[string]string{"docs/decisions/ADR-032-local.md": "---\nid: ADR-032\ntype: decision\nstatus: inferred\nlinks: []\ntitle: An unrelated adopter decision\nauthor: agent\naccepted-by: []\n---\n"})
	c.ByID["CAP-901"][0].Title = "Actual capability"
	c.ByID["CAP-901-criteria"][0].Body = strings.Replace(c.ByID["CAP-901-criteria"][0].Body, "A person sees the failed backup", "A person sees CAP-901", 1)
	issue := localIssue(c, "docs/capabilities/CAP-901/criteria.md", "AC-901 has no Unit positive evidence for CAP-901 (ADR-032)", "AC-901", "CAP-901")
	human := NewReferenceNames(c).Issue(issue)
	if !strings.Contains(human, "A person sees CAP-901 (AC-901)") || !strings.Contains(human, "Actual capability (CAP-901)") || !strings.Contains(human, "(ADR-032)") || strings.Contains(human, "unrelated adopter decision") {
		t.Fatal(human)
	}
	unregistered := Issue{"docs/goals/G-901.md", "foreign rule ADR-032 and reference AC-901"}
	got := NewReferenceNames(c).Issue(unregistered)
	if !strings.Contains(got, unregistered.Msg) || strings.Contains(got, "A person sees") || strings.Contains(got, "unrelated adopter decision") {
		t.Fatal(got)
	}
	empty := localIssue(c, "unknown", "literal id follows", "")
	if got := NewReferenceNames(c).Issue(empty); got != "unknown: literal id follows" {
		t.Fatal(got)
	}
}

func TestAC233_UnitNegative_LegacyUnnamedMilestonesRemainKnownWithoutInventedNames(t *testing.T) {
	c := readableFixture(t, nil)
	c.ByID["P-901"][0].Body += "\n| M-902 | Legacy row without a named table header | todo |\n"
	if got := NewReferenceNames(c).Label("M-902"); got != "unnamed reference (M-902)" {
		t.Fatal(got)
	}
	if got := HumanReference("G-901", "G-901"); got != "unnamed reference (G-901)" {
		t.Fatal(got)
	}
	c.Artifacts = append(c.Artifacts, &Artifact{ID: "clue:other/repo/G-901", Title: "Pretended local title", Path: "fake.md"})
	if got := NewReferenceNames(c).Label("clue:other/repo/G-901"); strings.Contains(got, "Pretended") {
		t.Fatal(got)
	}
}
