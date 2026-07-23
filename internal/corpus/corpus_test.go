package corpus

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// validFiles is a minimal corpus that must pass every rule.
var validFiles = map[string]string{
	"docs/README.md":               "# Corpus\n\n<!-- clue:index:start -->\n- [goals/](goals/README.md)\n- [plans/](plans/README.md)\n<!-- clue:index:end -->\n",
	"docs/goals/README.md":         "# Goals\n\n<!-- clue:index:start -->\n- [G-001](G-001-first.md)\n<!-- clue:index:end -->\n",
	"docs/goals/G-001-first.md":    "---\nid: G-001\ntype: goal\nstatus: accepted\nlinks: []\ntitle: First goal\n---\n\n# G-001\n",
	"docs/plans/README.md":         "# Plans\n\n<!-- clue:index:start -->\n- [P-001](P-001-baseline.md)\n<!-- clue:index:end -->\n",
	"docs/plans/P-001-baseline.md": "---\nid: P-001\ntype: plan\nstatus: active\nlinks: [G-001]\ntitle: Baseline\n---\n\n| M-001 | do it | todo |\n",
}

func writeCorpus(t *testing.T, files map[string]string) string {
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

func with(base map[string]string, extra map[string]string) map[string]string {
	out := map[string]string{}
	for k, v := range base {
		out[k] = v
	}
	for k, v := range extra {
		out[k] = v
	}
	return out
}

func run(t *testing.T, files map[string]string, forbid bool) []Issue {
	t.Helper()
	c, issues := Scan(writeCorpus(t, files))
	return append(issues, Validate(c, Options{ForbidChanges: forbid})...)
}

func assertIssue(t *testing.T, issues []Issue, substr string) {
	t.Helper()
	for _, i := range issues {
		if strings.Contains(i.String(), substr) {
			return
		}
	}
	t.Fatalf("expected an issue containing %q, got %v", substr, issues)
}

// AC-004: a valid corpus exits clean.
func TestAC004_ValidCorpusHasNoIssues(t *testing.T) {
	if issues := run(t, validFiles, true); len(issues) != 0 {
		t.Fatalf("expected no issues, got %v", issues)
	}
}

// AC-005: a missing core field is reported with file and field.
func TestAC005_MissingCoreFieldReported(t *testing.T) {
	files := with(validFiles, map[string]string{
		"docs/goals/G-001-first.md": "---\nid: G-001\ntype: goal\nlinks: []\ntitle: First goal\n---\n",
	})
	issues := run(t, files, false)
	assertIssue(t, issues, "G-001-first.md")
	assertIssue(t, issues, "status")
}

func TestAC005_MissingFrontmatterReported(t *testing.T) {
	files := with(validFiles, map[string]string{
		"docs/goals/G-002-bare.md": "# No frontmatter here\n",
		"docs/goals/README.md":     "# Goals\n\n<!-- clue:index:start -->\n- [G-001](G-001-first.md)\n- [G-002](G-002-bare.md)\n<!-- clue:index:end -->\n",
	})
	assertIssue(t, run(t, files, false), "missing frontmatter")
}

// AC-006: a link to an unknown ID is reported with file and ID.
func TestAC006_UnresolvedLinkReported(t *testing.T) {
	files := with(validFiles, map[string]string{
		"docs/goals/G-001-first.md": "---\nid: G-001\ntype: goal\nstatus: accepted\nlinks: [G-999]\ntitle: First goal\n---\n",
	})
	issues := run(t, files, false)
	assertIssue(t, issues, "G-001-first.md")
	assertIssue(t, issues, "G-999")
}

func TestAC006_MilestoneLinksResolveViaPlanBody(t *testing.T) {
	files := with(validFiles, map[string]string{
		"docs/goals/G-001-first.md": "---\nid: G-001\ntype: goal\nstatus: accepted\nlinks: [M-001]\ntitle: First goal\n---\n",
	})
	if issues := run(t, files, false); len(issues) != 0 {
		t.Fatalf("M-001 exists in P-001 body; expected no issues, got %v", issues)
	}
	files["docs/goals/G-001-first.md"] = "---\nid: G-001\ntype: goal\nstatus: accepted\nlinks: [M-999]\ntitle: First goal\n---\n"
	assertIssue(t, run(t, files, false), "M-999")
}

func TestUnit_DuplicateIDReported(t *testing.T) {
	files := with(validFiles, map[string]string{
		"docs/goals/G-001-copy.md": validFiles["docs/goals/G-001-first.md"],
		"docs/goals/README.md":     "# Goals\n\n<!-- clue:index:start -->\n- [G-001](G-001-first.md)\n- [copy](G-001-copy.md)\n<!-- clue:index:end -->\n",
	})
	assertIssue(t, run(t, files, false), "duplicate id G-001")
}

func TestUnit_StatusVocabEnforced(t *testing.T) {
	files := with(validFiles, map[string]string{
		"docs/goals/G-001-first.md": "---\nid: G-001\ntype: goal\nstatus: wip\nlinks: []\ntitle: First goal\n---\n",
	})
	assertIssue(t, run(t, files, false), "status wip not allowed for type goal")
}

// ADR-026: a type the validator does not recognize is an adopter extension,
// validated against the default lifecycle rather than rejected.
func TestUnit_AdopterTypeValidatesAgainstDefault(t *testing.T) {
	files := with(validFiles, map[string]string{
		"docs/goals/G-001-first.md": "---\nid: G-001\ntype: risk\nstatus: active\nlinks: []\ntitle: An adopter-defined type\n---\n",
	})
	if issues := run(t, files, false); len(issues) != 0 {
		t.Fatalf("an adopter type with a default status is valid; expected no issues, got %v", issues)
	}

	files["docs/goals/G-001-first.md"] = "---\nid: G-001\ntype: risk\nstatus: accepted\nlinks: []\ntitle: An adopter-defined type\n---\n"
	assertIssue(t, run(t, files, false), "status accepted not allowed for type risk (allowed: draft, active, retired)")
}

func TestUnit_LogStatusVocab(t *testing.T) {
	logFiles := map[string]string{
		"docs/README.md":           "# Corpus\n\n<!-- clue:index:start -->\n- [goals/](goals/README.md)\n- [plans/](plans/README.md)\n- [decisions/](decisions/README.md)\n<!-- clue:index:end -->\n",
		"docs/decisions/README.md": "# Decisions\n\n<!-- clue:index:start -->\n- [log](log.md)\n<!-- clue:index:end -->\n",
		"docs/decisions/log.md":    "---\nid: LOG-001\ntype: log\nstatus: active\nlinks: []\ntitle: Decision log\n---\n\n| Date | Decision | Why | Change/PR |\n",
	}
	if issues := run(t, with(validFiles, logFiles), false); len(issues) != 0 {
		t.Fatalf("an active log is valid; expected no issues, got %v", issues)
	}

	logFiles["docs/decisions/log.md"] = strings.Replace(logFiles["docs/decisions/log.md"], "status: active", "status: open", 1)
	assertIssue(t, run(t, with(validFiles, logFiles), false), "status open not allowed for type log")
}

func TestUnit_FolderWithoutReadme(t *testing.T) {
	files := with(validFiles, map[string]string{
		"docs/quality/QS-001-fast.md": "---\nid: QS-001\ntype: quality\nstatus: active\nlinks: []\ntitle: Fast\n---\n",
		"docs/README.md":              "# Corpus\n\n<!-- clue:index:start -->\n- [goals/](goals/README.md)\n- [plans/](plans/README.md)\n- [quality/](quality/QS-001-fast.md)\n<!-- clue:index:end -->\n",
	})
	assertIssue(t, run(t, files, false), "docs/quality: folder has no README.md")
}

// AC-007: index drift — a block link to a missing file, or a sibling
// artifact the block does not reference.
func TestAC007_IndexDriftReported(t *testing.T) {
	files := with(validFiles, map[string]string{
		"docs/goals/README.md": "# Goals\n\n<!-- clue:index:start -->\n- [G-001](G-001-first.md)\n- [gone](G-777-gone.md)\n<!-- clue:index:end -->\n",
	})
	assertIssue(t, run(t, files, false), "index references missing file G-777-gone.md")

	files = with(validFiles, map[string]string{
		"docs/goals/README.md": "# Goals\n\n<!-- clue:index:start -->\n<!-- clue:index:end -->\n",
	})
	assertIssue(t, run(t, files, false), "index does not reference sibling G-001-first.md")
}

func TestAC007_MissingIndexMarkersReported(t *testing.T) {
	files := with(validFiles, map[string]string{
		"docs/goals/README.md": "# Goals — hand-written, no markers\n",
	})
	assertIssue(t, run(t, files, false), "index markers missing")
}

// AC-008: the digest-before-merge gate.
func TestAC008_ForbidChangesGate(t *testing.T) {
	files := with(validFiles, map[string]string{
		"changes/CH-009-x/proposal.md": "---\nid: CH-009\ntype: change\nstatus: open\nlinks: [P-001]\ntitle: X\n---\n",
	})
	if issues := run(t, files, false); len(issues) != 0 {
		t.Fatalf("changes/ allowed without the gate; got %v", issues)
	}
	assertIssue(t, run(t, files, true), "digest before merge")
}

func TestUnit_CRLFFrontmatterParses(t *testing.T) {
	files := with(validFiles, map[string]string{
		"docs/goals/G-001-first.md": strings.ReplaceAll(validFiles["docs/goals/G-001-first.md"], "\n", "\r\n"),
	})
	if issues := run(t, files, false); len(issues) != 0 {
		t.Fatalf("CRLF corpus should validate; got %v", issues)
	}
}

// AC-018: the optional provenance field is linted (ADR-010).
func TestAC018_ProvenanceVocabularyLinted(t *testing.T) {
	files := with(validFiles, map[string]string{
		"docs/goals/G-001-first.md": "---\nid: G-001\ntype: goal\nstatus: accepted\nlinks: []\ntitle: First goal\nprovenance: guessed\n---\n",
	})
	assertIssue(t, run(t, files, false), "provenance must be inferred or verified")

	files["docs/goals/G-001-first.md"] = "---\nid: G-001\ntype: goal\nstatus: accepted\nlinks: []\ntitle: First goal\nprovenance: inferred\n---\n"
	if issues := run(t, files, false); len(issues) != 0 {
		t.Fatalf("provenance: inferred is valid; expected no issues, got %v", issues)
	}
}

func TestAC018_DecisionsMustNotCarryProvenanceField(t *testing.T) {
	files := with(validFiles, map[string]string{
		"docs/README.md":              "# Corpus\n\n<!-- clue:index:start -->\n- [goals/](goals/README.md)\n- [plans/](plans/README.md)\n- [decisions/](decisions/README.md)\n<!-- clue:index:end -->\n",
		"docs/decisions/README.md":    "# Decisions\n\n<!-- clue:index:start -->\n- [ADR-001](ADR-001-x.md)\n<!-- clue:index:end -->\n",
		"docs/decisions/ADR-001-x.md": "---\nid: ADR-001\ntype: decision\nstatus: inferred\nlinks: []\ntitle: X\nprovenance: inferred\n---\n",
	})
	assertIssue(t, run(t, files, false), "decisions carry provenance in status")
}

// AC-023: constraints carry a non-empty source and an enforcement class
// from machine|agent|human.
func TestAC023_ConstraintRegisterFieldsLinted(t *testing.T) {
	constraintFiles := func(frontmatter string) map[string]string {
		return with(validFiles, map[string]string{
			"docs/README.md":                 "# Corpus\n\n<!-- clue:index:start -->\n- [goals/](goals/README.md)\n- [plans/](plans/README.md)\n- [constraints/](constraints/README.md)\n<!-- clue:index:end -->\n",
			"docs/constraints/README.md":     "# Constraints\n\n<!-- clue:index:start -->\n- [C-001](C-001-rule.md)\n<!-- clue:index:end -->\n",
			"docs/constraints/C-001-rule.md": frontmatter,
		})
	}

	valid := "---\nid: C-001\ntype: constraint\nstatus: active\nlinks: []\ntitle: A rule\nsource: AGENTS.md rule 5\nenforcement: agent\n---\n"
	if issues := run(t, constraintFiles(valid), false); len(issues) != 0 {
		t.Fatalf("a sourced, agent-enforced constraint is valid; expected no issues, got %v", issues)
	}

	noSource := "---\nid: C-001\ntype: constraint\nstatus: active\nlinks: []\ntitle: A rule\nenforcement: agent\n---\n"
	assertIssue(t, run(t, constraintFiles(noSource), false), "constraint missing or empty source field")

	noEnforcement := "---\nid: C-001\ntype: constraint\nstatus: active\nlinks: []\ntitle: A rule\nsource: AGENTS.md rule 5\n---\n"
	assertIssue(t, run(t, constraintFiles(noEnforcement), false), "constraint missing or empty enforcement field")

	badVocab := "---\nid: C-001\ntype: constraint\nstatus: active\nlinks: []\ntitle: A rule\nsource: AGENTS.md rule 5\nenforcement: hope\n---\n"
	assertIssue(t, run(t, constraintFiles(badVocab), false), "enforcement hope not allowed")
}

// The dogfood test: this repository's own corpus must always be valid.
func TestSanity_RepoCorpusIsValid(t *testing.T) {
	root, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(root, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(root)
		if parent == root {
			t.Fatal("repo root with go.mod not found")
		}
		root = parent
	}
	c, issues := Scan(root)
	issues = append(issues, Validate(c, Options{ForbidChanges: false})...)
	if len(issues) != 0 {
		t.Fatalf("the repo's own corpus has issues:\n%v", issues)
	}
	if len(c.Artifacts) == 0 {
		t.Fatal("expected artifacts in the repo corpus")
	}
}

// testBOM is built from the code point so no literal BOM lands in this
// source file (Go rejects U+FEFF anywhere but byte 0).
var testBOM = string(rune(0xFEFF))

func assertNoIssue(t *testing.T, issues []Issue, substr string) {
	t.Helper()
	for _, i := range issues {
		if strings.Contains(i.String(), substr) {
			t.Fatalf("unexpected issue containing %q: %v", substr, i)
		}
	}
}

// AC-034: a UTF-8 BOM at the start of a corpus file is reported by name.
func TestAC034_BOMAtStartReported(t *testing.T) {
	files := with(validFiles, map[string]string{
		"docs/goals/G-001-first.md": testBOM + "---\nid: G-001\ntype: goal\nstatus: accepted\nlinks: []\ntitle: First goal\n---\n\n# G-001\n",
	})
	issues := run(t, files, false)
	assertIssue(t, issues, "G-001-first.md")
	assertIssue(t, issues, "byte-order mark")
}

// AC-034: an embedded BOM is reported too.
func TestAC034_EmbeddedBOMReported(t *testing.T) {
	files := with(validFiles, map[string]string{
		"docs/goals/G-001-first.md": "---\nid: G-001\ntype: goal\nstatus: accepted\nlinks: []\ntitle: First goal\n---\n\n# G-001\n\nText" + testBOM + " here.\n",
	})
	issues := run(t, files, false)
	assertIssue(t, issues, "byte-order mark")
}

// AC-034 negative: a BOM-free corpus raises no byte-order-mark issue.
func TestAC034_BOMFreeCorpusClean(t *testing.T) {
	assertNoIssue(t, run(t, validFiles, false), "byte-order mark")
}

// AC-035: a body opening with another frontmatter fence is a leftover
// second frontmatter, with or without a BOM hiding it.
func TestAC035_SecondFrontmatterReported(t *testing.T) {
	files := with(validFiles, map[string]string{
		"docs/goals/G-001-first.md": "---\nid: G-001\ntype: goal\nstatus: accepted\nlinks: []\ntitle: First goal\n---\n\n---\nstatus: accepted\ndate: 2026-01-01\n---\n\n# G-001\n",
	})
	issues := run(t, files, false)
	assertIssue(t, issues, "G-001-first.md")
	assertIssue(t, issues, "second frontmatter")
}

func TestAC035_BOMHiddenSecondFrontmatterReported(t *testing.T) {
	files := with(validFiles, map[string]string{
		"docs/goals/G-001-first.md": "---\nid: G-001\ntype: goal\nstatus: accepted\nlinks: []\ntitle: First goal\n---\n\n" + testBOM + "---\nstatus: accepted\ndate: 2026-01-01\n---\n\n# G-001\n",
	})
	issues := run(t, files, false)
	assertIssue(t, issues, "second frontmatter")
}

// AC-035 negative: thematic breaks are not frontmatter fences — an
// unclosed break, a pair enclosing ordinary markdown, or an empty pair
// stays clean wherever it appears in the body.
func TestAC035_ThematicBreaksClean(t *testing.T) {
	tests := map[string]string{
		"opens body":                "---\nid: G-001\ntype: goal\nstatus: accepted\nlinks: []\ntitle: First goal\n---\n\n---\n\n# G-001\n",
		"later in body":             "---\nid: G-001\ntype: goal\nstatus: accepted\nlinks: []\ntitle: First goal\n---\n\n# G-001\n\nAbove the break.\n\n---\n\nBelow the break.\n",
		"opens body, another later": "---\nid: G-001\ntype: goal\nstatus: accepted\nlinks: []\ntitle: First goal\n---\n\n---\n\n# G-001\n\nBetween the breaks.\n\n---\n\nBelow the break.\n",
		"empty between breaks":      "---\nid: G-001\ntype: goal\nstatus: accepted\nlinks: []\ntitle: First goal\n---\n\n---\n\n---\n\n# G-001\n",
	}
	for name, content := range tests {
		t.Run(name, func(t *testing.T) {
			files := with(validFiles, map[string]string{
				"docs/goals/G-001-first.md": content,
			})
			assertNoIssue(t, run(t, files, false), "second frontmatter")
		})
	}
}
