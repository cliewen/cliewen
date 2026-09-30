package corpus

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cliewen/cliewen/internal/evidence"
)

func exportedCorpus(t *testing.T) (string, *Corpus) {
	t.Helper()
	files := capFiles("active")
	files["docs/capabilities/CAP-101-x/criteria.md"] = strings.Replace(files["docs/capabilities/CAP-101-x/criteria.md"], "    Given", "    Test-type: Unit\n    Given", 1)
	files["backend/test.unknown"] = "backend executable"
	files["frontend/test.other"] = "frontend executable"
	root := writeCorpus(t, files)
	var producers []evidence.Producer
	for _, tc := range []struct{ id, path, direction string }{{"backend", "backend/test.unknown", "positive"}, {"frontend", "frontend/test.other", "negative"}} {
		p := evidence.Producer{ID: tc.id, Include: []string{tc.path}, References: []evidence.Reference{{ID: "AC-101", Path: tc.path, Subject: "same test name", Type: "Unit", Direction: tc.direction}}}
		var err error
		p.Inputs, err = evidence.Snapshot(root, p)
		if err != nil {
			t.Fatal(err)
		}
		producers = append(producers, p)
	}
	if err := evidence.Write(root, evidence.Manifest{Version: 1, Producers: producers}); err != nil {
		t.Fatal(err)
	}
	c, issues := Scan(root)
	if len(issues) > 0 {
		t.Fatal(issues)
	}
	return root, c
}

func TestAC203_UnitPositive_CommonEvidenceFeedsValidationAndCoverage(t *testing.T) {
	_, c := exportedCorpus(t)
	if issues := Validate(c, Options{}); len(issues) != 0 {
		t.Fatal(issues)
	}
	declared, refs, issues := AcceptanceEvidence(c)
	if len(issues) > 0 || declared["AC-101"].TestType != "Unit" || len(refs["AC-101"]) != 2 {
		t.Fatalf("lost common evidence: %v %v", refs, issues)
	}
	if refs["AC-101"][0].Producer == refs["AC-101"][1].Producer {
		t.Fatal("lost producer identities")
	}
	if got := Coverage(c); len(got) != 1 || got[0].State != "covered" {
		t.Fatalf("coverage: %v", got)
	}
}

func TestAC203_UnitNegative_UnknownAndRetiredExportReferencesFail(t *testing.T) {
	for _, id := range []string{"AC-999", "AC-101"} {
		t.Run(id, func(t *testing.T) {
			root, c := exportedCorpus(t)
			m, err := evidence.Load(root)
			if err != nil {
				t.Fatal(err)
			}
			m.Producers[0].References[0].ID = id
			if err := evidence.Write(root, m); err != nil {
				t.Fatal(err)
			}
			if id == "AC-101" {
				for _, a := range c.Artifacts {
					if a.Type == "criteria" {
						a.Body = "@AC-101 @retired\nScenario: retired\n"
					}
				}
			}
			issues := checkACTests(c)
			if id == "AC-999" {
				assertIssue(t, issues, "which no criteria.md declares")
			} else {
				assertIssue(t, issues, "references retired AC-101")
			}
		})
	}
}

func TestAC206_UnitNegative_InvalidExportDoesNotLeaveCoveredClaims(t *testing.T) {
	root, c := exportedCorpus(t)
	if err := os.WriteFile(filepath.Join(root, "backend/test.unknown"), []byte("changed source"), 0644); err != nil {
		t.Fatal(err)
	}
	_, refs, issues := AcceptanceEvidence(c)
	if len(refs) != 0 || len(issues) == 0 {
		t.Fatalf("stale references leaked: %v %v", refs, issues)
	}
	if got := Coverage(c); got[0].State != "gap" {
		t.Fatalf("stale coverage: %v", got)
	}
}
