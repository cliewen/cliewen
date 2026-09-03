package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// intentCorpus writes a minimal corpus with a capability, an optional use
// case, and — when vision is non-empty — a vision at its fixed address.
func intentCorpus(t *testing.T, vision string) string {
	t.Helper()
	root := t.TempDir()
	files := map[string]string{
		"docs/README.md":                            "# Corpus\n\n<!-- clue:index:start -->\n- [goals/](goals/README.md)\n- [use-cases/](use-cases/README.md)\n- [capabilities/](capabilities/README.md)\n<!-- clue:index:end -->\n",
		"docs/goals/README.md":                      "# Goals\n\n<!-- clue:index:start -->\n- [G-001](G-001-first.md)\n<!-- clue:index:end -->\n",
		"docs/goals/G-001-first.md":                 "---\nid: G-001\ntype: goal\nstatus: accepted\nlinks: []\ntitle: First goal\n---\n\n# G-001\n",
		"docs/capabilities/README.md":               "# Capabilities\n\n<!-- clue:index:start -->\n- [CAP-001](CAP-001-thing/README.md)\n- [CAP-002](CAP-002-other/README.md)\n<!-- clue:index:end -->\n",
		"docs/capabilities/CAP-001-thing/README.md": "---\nid: CAP-001\ntype: capability\nstatus: active\nlinks: [G-001]\ntitle: A capability\ngoal: G-001\n---\n\n# CAP-001\n",
		"docs/capabilities/CAP-002-other/README.md": "---\nid: CAP-002\ntype: capability\nstatus: active\nlinks: [G-001]\ntitle: Another capability\ngoal: G-001\n---\n\n# CAP-002\n",
		"docs/use-cases/README.md":                  "# Use cases\n\n<!-- clue:index:start -->\n- [UC-001](UC-001-a-journey.md)\n<!-- clue:index:end -->\n",
		"docs/use-cases/UC-001-a-journey.md":        "---\nid: UC-001\ntype: use-case\nstatus: active\nlinks: [G-001, CAP-001]\ntitle: A journey\n---\n\n# UC-001\n\n## Actors\n\nSomeone.\n\n## Trigger\n\nThey ask.\n\n## Main flow\n\n1. It happens.\n\n## Outcome\n\nIt happened.\n",
	}
	if vision != "" {
		files["docs/vision.md"] = vision
		files["docs/README.md"] = strings.Replace(files["docs/README.md"], "<!-- clue:index:start -->\n", "<!-- clue:index:start -->\n- [vision.md](vision.md)\n", 1)
	}
	for rel, content := range files {
		full := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

// AC-164: the report states the corpus's direction and its journeys.
func TestAC164_UnitPositive_IntentReportNamesTheVisionAndItsProvenance(t *testing.T) {
	root := intentCorpus(t, "---\nid: VIS-001\ntype: vision\nstatus: active\nlinks: []\ntitle: A product\nprovenance: inferred\nreversal-cost: low\n---\n\n# VIS-001\n\nWhat it is for.\n")
	var out bytes.Buffer
	if code := runValidate([]string{"--intent", root}, &out); code != 0 {
		t.Fatalf("expected exit 0, got %d\n%s", code, out.String())
	}
	printed := out.String()
	for _, want := range []string{"vision: VIS-001", "active", "inferred", "use case: UC-001", "crosses CAP-001"} {
		if !strings.Contains(printed, want) {
			t.Fatalf("intent report does not state %q:\n%s", want, printed)
		}
	}
}

// AC-164: an absent vision is a state, and no ratio is ever printed.
func TestAC164_UnitNegative_IntentReportComputesNoCoverageFigure(t *testing.T) {
	root := intentCorpus(t, "")
	var out bytes.Buffer
	if code := runValidate([]string{"--intent", root}, &out); code != 0 {
		t.Fatalf("expected a corpus with no vision to pass, got %d\n%s", code, out.String())
	}
	printed := out.String()
	if !strings.Contains(printed, "vision: none") {
		t.Fatalf("intent report does not state the absence:\n%s", printed)
	}
	// A percentage over an optional artifact reads as a target, and the only
	// way to move it is to write use cases nobody needs (PDR-054).
	for _, forbidden := range []string{"%", "coverage", "0 of ", "1 of "} {
		if strings.Contains(strings.ToLower(printed), forbidden) {
			t.Fatalf("intent report computed a figure (%q):\n%s", forbidden, printed)
		}
	}
}

// AC-165: a capability's slice names the journeys that reach it.
func TestAC165_UnitPositive_ContextNamesTheUseCasesReachingTheRoot(t *testing.T) {
	root := intentCorpus(t, "")
	var out, errOut bytes.Buffer
	if code := runContext([]string{"CAP-001", root}, &out, &errOut); code != 0 {
		t.Fatalf("expected exit 0, got %d\n%s", code, errOut.String())
	}
	printed := out.String()
	if !strings.Contains(printed, "use case(s) naming this artifact") || !strings.Contains(printed, "UC-001 | A journey | docs/use-cases/UC-001-a-journey.md") {
		t.Fatalf("context did not name the use case:\n%s", printed)
	}
	// Names only: the use case's own body never reaches the slice.
	if strings.Contains(printed, "## Main flow") {
		t.Fatalf("context expanded the use case's content:\n%s", printed)
	}
}

// AC-165: an artifact no use case names gets no section at all.
func TestAC165_UnitNegative_ContextAddsNothingWhenNoUseCaseNamesTheRoot(t *testing.T) {
	root := intentCorpus(t, "")
	var out, errOut bytes.Buffer
	if code := runContext([]string{"CAP-002", root}, &out, &errOut); code != 0 {
		t.Fatalf("expected exit 0, got %d\n%s", code, errOut.String())
	}
	if strings.Contains(out.String(), "naming this artifact") {
		t.Fatalf("context named a use case for an artifact none of them links:\n%s", out.String())
	}
}
