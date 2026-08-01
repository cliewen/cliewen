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
	c := corpusWithBody(t, "proven by clue:robocode-dev/tank-royale@384d27d5/BR-001 there")
	if issues := checkForeignPointers(c); len(issues) != 0 {
		t.Fatalf("a pointer naming repository, revision and identifier must pass, got %v", issues)
	}
	got := ForeignPointers(c)
	if len(got) != 1 || got[0] != "clue:robocode-dev/tank-royale@384d27d5/BR-001" {
		t.Fatalf("expected the pointer reported once, got %v", got)
	}
}

func TestAC067_UnitPositive_CrossRepositoryIdentityNeedsNoRevision(t *testing.T) {
	// A plain cross-repository citation names an identity; only an
	// acceptance-evidence pointer pins the state that was proven.
	c := corpusWithBody(t, "see clue:robocode-dev/tank-royale/CAP-001 upstream")
	if issues := checkForeignPointers(c); len(issues) != 0 {
		t.Fatalf("a revisionless identity must pass, got %v", issues)
	}
	if got := ForeignPointers(c); len(got) != 1 || got[0] != "clue:robocode-dev/tank-royale/CAP-001" {
		t.Fatalf("expected the identity reported, got %v", got)
	}
}

func TestAC067_UnitNegative_AnExampleInAFenceIsNotAClaim(t *testing.T) {
	// Documentation showing the form must not appear in a coverage summary as
	// though something had been proven.
	c := corpusWithBody(t, "The form is:\n\n```\nclue:robocode-dev/tank-royale@384d27d5/BR-001\n```\n\nand `clue:owner/repo/ID` inline.")
	if got := ForeignPointers(c); len(got) != 0 {
		t.Fatalf("a fenced or code-span example must not be reported as a claim, got %v", got)
	}
	if issues := checkForeignPointers(c); len(issues) != 0 {
		t.Fatalf("an example must not be judged malformed either, got %v", issues)
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

func TestAC067_UnitNegative_OrdinaryProseIsNotABrokenPointer(t *testing.T) {
	// "clue" is an ordinary English noun and the scheme word at once. A
	// sentence that happens to use it must not fail validation, and neither
	// must an index marker written outside a code span.
	for _, prose := range []string{
		"The clue: it was the deployment order all along.",
		"Write clue: followed by the repository to cite one.",
		"<!-- clue:index:start -->",
	} {
		t.Run(prose, func(t *testing.T) {
			if issues := checkForeignPointers(corpusWithBody(t, prose)); len(issues) != 0 {
				t.Fatalf("ordinary prose must not be judged a malformed pointer, got %v", issues)
			}
		})
	}
}

func TestAC067_UnitNegative_AMalformedPointerIsNamedInFull(t *testing.T) {
	// The colon belongs to the scheme: trimming it made the diagnostic name a
	// word rather than the token that was actually found.
	issues := checkForeignPointers(corpusWithBody(t, "proven by clue:tank-royale/BR-001."))
	if len(issues) != 1 {
		t.Fatalf("expected one finding, got %v", issues)
	}
	if !strings.Contains(issues[0].Msg, "clue:tank-royale/BR-001") {
		t.Fatalf("expected the token named in full, got %q", issues[0].Msg)
	}
}

func TestAC067_UnitNegative_ASchemeWithNoNamespaceIsMalformed(t *testing.T) {
	// The shape a reader is most likely to write by hand names no repository
	// at all. Reading it as prose would let a criterion look as though it had
	// cited its foreign proof while the judge said nothing.
	issues := checkForeignPointers(corpusWithBody(t, "proven by clue:BR-001 upstream"))
	if len(issues) != 1 {
		t.Fatalf("expected the namespace-less pointer to fail, got %v", issues)
	}
	if !strings.Contains(issues[0].Msg, "clue:BR-001") {
		t.Fatalf("expected the token named, got %q", issues[0].Msg)
	}
}

func TestAC066_UnitPositive_TheChecksAreWiredIntoTheJudge(t *testing.T) {
	// Both criteria say "when the user runs clue validate". Without this, the
	// checks could be unregistered and every other test would stay green.
	c := corpusWithBody(t, "a bare PR #96 and a broken clue:BR-001 pointer")
	var found int
	for _, issue := range Validate(c, Options{}) {
		if strings.Contains(issue.Msg, "bare forge reference #96") ||
			strings.Contains(issue.Msg, "clue:BR-001 is malformed") {
			found++
		}
	}
	if found != 2 {
		t.Fatalf("expected both checks to run inside Validate, matched %d of 2", found)
	}
}

func TestAC067_UnitNegative_CapitalisedProseIsNotAPointer(t *testing.T) {
	// The identity half requires the digits a corpus ID carries, so a
	// capitalised word after the scheme is prose, not a broken pointer.
	for _, prose := range []string{
		"The clue:Index was a red herring.",
		"See https://example.com/search/clue:Overview for context.",
	} {
		t.Run(prose, func(t *testing.T) {
			if issues := checkForeignPointers(corpusWithBody(t, prose)); len(issues) != 0 {
				t.Fatalf("expected no finding, got %v", issues)
			}
		})
	}
}

func TestAC066_UnitNegative_ADoubleBacktickSpanIsStillContent(t *testing.T) {
	// A span opens and closes with runs of the same length, so a literal that
	// itself contains a backtick is one span rather than three.
	for _, body := range []string{
		"See ``the #96 form`` above.",
		"Use ``a `nested` #42 example`` here.",
	} {
		t.Run(body, func(t *testing.T) {
			if got := scanBody(t, body); len(got) != 0 {
				t.Fatalf("expected no finding inside a code span, got %v", got)
			}
		})
	}
	// A lone run with no closing partner is literal text, not an open span
	// swallowing the rest of the line.
	if got := scanBody(t, "``unclosed and then a bare PR #96 here"); len(got) != 1 {
		t.Fatalf("expected the bare reference outside an unclosed run, got %v", got)
	}
}
