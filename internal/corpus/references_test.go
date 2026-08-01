package corpus

import (
	"strings"
	"testing"
)

// scanBody is the unit under test for the form rule: it decides where a '#'
// can be a citation at all.
func scanBody(t *testing.T, body string) []string {
	t.Helper()
	var out []string
	for _, r := range scanBareForgeRefs(body) {
		out = append(out, r.text)
	}
	return out
}

func TestAC066_UnitPositive_BareForgeReferenceIsTheOnlyFinding(t *testing.T) {
	body := "A bare citation in PR #96 names no repository.\n"
	got := scanBody(t, body)
	if len(got) != 1 || got[0] != "#96" {
		t.Fatalf("expected the bare reference to be found once, got %v", got)
	}
}

func TestAC066_UnitPositive_LineNamesWhatAReaderOpens(t *testing.T) {
	// The finding must name the line in the file, not the offset within the
	// body: frontmatter sits above it, so an unoffset number points several
	// lines too high and sends a reader to the wrong place.
	c := corpusWithBody(t, "first\nsecond\nthe bare PR #12 lives here")
	issues := checkExternalReferences(c)
	if len(issues) != 1 {
		t.Fatalf("expected one finding, got %v", issues)
	}
	// The fixture closes its frontmatter on line 9 and leaves line 10 blank,
	// so the body's third line is the file's thirteenth.
	if !strings.Contains(issues[0].Msg, "line 13:") {
		t.Fatalf("expected the file line, got %q", issues[0].Msg)
	}
}

func TestAC066_UnitNegative_QualifiedAndNonReferenceFormsAreNotFindings(t *testing.T) {
	cases := []struct{ name, body string }{
		{"full URL", "See https://github.com/cliewen/cliewen/pull/96 for the merge.\n"},
		{"markdown target", "See [the merge](https://github.com/cliewen/cliewen/pull/96).\n"},
		{"labelled link keeps the shorthand", "See [cliewen/cliewen#96](https://github.com/cliewen/cliewen/pull/96).\n"},
		{"autolink", "See <https://github.com/cliewen/cliewen/pull/96>.\n"},
		{"heading anchor", "Jump to [the step](#4-resume-game) below.\n"},
		{"colour literal in a code span", "Use `#777777` for the line colour.\n"},
		{"fenced block", "```mermaid\nUpdateRelStyle(a, b, $lineColor=\"#777777\")\n```\n"},
		{"heading marker", "# 42 things\n"},
		{"clue identity", "Proven by clue:robocode-dev/tank-royale@384d27d5/BR-001 upstream.\n"},
		{"corpus ID and relative path are untouched", "See [ADR-040](../decisions/ADR-040-qualified-external-references.md).\n"},
		{"number inside a word", "The colour C#5 is not a citation.\n"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := scanBody(t, c.body); len(got) != 0 {
				t.Fatalf("expected no finding, got %v", got)
			}
		})
	}
}

func TestAC067_UnitPositive_ForeignEvidencePointerIsWellFormed(t *testing.T) {
	m := clueRefRe.FindStringSubmatch("proven by clue:robocode-dev/tank-royale@384d27d5/BR-001 there")
	if m == nil {
		t.Fatal("expected the pointer to parse")
	}
	if m[1] != "robocode-dev" || m[2] != "tank-royale" || m[3] != "@384d27d5" || m[4] != "BR-001" {
		t.Fatalf("expected owner, repo, revision and identifier, got %q", m[1:])
	}
}

func TestAC067_UnitPositive_CrossRepositoryIdentityNeedsNoRevision(t *testing.T) {
	// A plain cross-repository citation names an identity; only an
	// acceptance-evidence pointer pins the state that was proven.
	m := clueRefRe.FindStringSubmatch("see clue:robocode-dev/tank-royale/CAP-001 upstream")
	if m == nil || m[3] != "" || m[4] != "CAP-001" {
		t.Fatalf("expected a revisionless identity, got %v", m)
	}
}

func TestAC067_UnitNegative_MalformedPointersFailValidation(t *testing.T) {
	for _, c := range []struct{ name, text string }{
		{"no repository", "clue:tank-royale/BR-001"},
		{"no identifier", "clue:robocode-dev/tank-royale@384d27d5/"},
		{"lowercase identifier", "clue:robocode-dev/tank-royale/br-001"},
	} {
		t.Run(c.name, func(t *testing.T) {
			c2 := corpusWithBody(t, "proven by "+c.text+" upstream")
			issues := checkForeignPointers(c2)
			if len(issues) != 1 {
				t.Fatalf("expected the malformed pointer to fail, got %v", issues)
			}
			if !strings.Contains(issues[0].Msg, "malformed") || !strings.Contains(issues[0].Msg, "line ") {
				t.Fatalf("expected a malformed diagnostic naming the line, got %q", issues[0].Msg)
			}
		})
	}
}

func TestAC067_UnitPositive_WellFormedPointerPassesAndIsReportable(t *testing.T) {
	c := corpusWithBody(t, "proven by clue:robocode-dev/tank-royale@384d27d5/BR-001 upstream")
	if issues := checkForeignPointers(c); len(issues) != 0 {
		t.Fatalf("a well-formed pointer must not fail, got %v", issues)
	}
	got := ForeignPointers(c)
	if len(got) != 1 || got[0] != "clue:robocode-dev/tank-royale@384d27d5/BR-001" {
		t.Fatalf("expected the pointer reportable as named-but-unproven, got %v", got)
	}
	// The pointer is a citation, never local proof: it must not reach the
	// classified-evidence path at all.
	if len(checkACTests(c)) != 0 {
		t.Fatal("a foreign pointer must never enter classified evidence")
	}
}

// corpusWithBody builds a one-artifact corpus whose frontmatter is real, so
// the reported line is the line a reader opens to rather than a body offset.
func corpusWithBody(t *testing.T, body string) *Corpus {
	t.Helper()
	text := "---\nid: AN-001\ntype: analysis\nstatus: active\nlinks: []\ntitle: t\nprovenance: inferred\nreversal-cost: low\n---\n\n" + body + "\n"
	fields, parsed, ok, err := parseFrontmatter(text)
	if err != nil || !ok {
		t.Fatalf("fixture frontmatter did not parse: %v", err)
	}
	a := &Artifact{Path: "docs/analysis/AN-001-x.md", Body: parsed, Fields: fields,
		BodyLine: strings.Count(text[:len(text)-len(parsed)], "\n") + 1}
	a.ID, _ = fields["id"].(string)
	a.Type, _ = fields["type"].(string)
	a.Status, _ = fields["status"].(string)
	return &Corpus{Artifacts: []*Artifact{a}, ByID: map[string][]*Artifact{a.ID: {a}}}
}
