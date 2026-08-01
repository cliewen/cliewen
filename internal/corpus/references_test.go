package corpus

import "testing"

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
	// The reported line is offset by where the body starts in the file, so a
	// finding points at the reference rather than several lines above it.
	body := "first\nsecond\nthe bare PR #12 lives here\n"
	refs := scanBareForgeRefs(body)
	if len(refs) != 1 {
		t.Fatalf("expected one reference, got %d", len(refs))
	}
	if refs[0].line != 3 {
		t.Fatalf("expected the third body line, got %d", refs[0].line)
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

func TestAC067_UnitNegative_MalformedPointersDoNotParse(t *testing.T) {
	cases := []struct{ name, text string }{
		{"no repository", "clue:tank-royale/BR-001"},
		{"no identifier", "clue:robocode-dev/tank-royale@384d27d5/"},
		{"lowercase identifier", "clue:robocode-dev/tank-royale/br-001"},
		{"no scheme", "robocode-dev/tank-royale/BR-001"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if m := clueRefRe.FindStringSubmatch(c.text); m != nil {
				t.Fatalf("expected no parse, got %v", m)
			}
		})
	}
}
