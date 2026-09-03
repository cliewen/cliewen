package corpus

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode/utf8"
)

// The extractor's reading, exercised where it lives. The row-level contract is
// held by AC-096, AC-097, and AC-160 in the scaffold package; these cases pin
// the prose rules those criteria depend on, including the ones a corpus is
// unlikely to produce on demand.
func TestUnit_DescribeBodyReadsTheSeedSentence(t *testing.T) {
	for _, tc := range []struct {
		name string
		body string
		want string // "" means no sentence could be read
	}{
		{
			name: "lede beneath the H1 wins over prose under a later heading",
			body: "# Title\n\nThe lede sentence. A second one.\n\n## Context\n\nProse that must lose.\n",
			want: "The lede sentence.",
		},
		{
			name: "no lede falls back to the first paragraph under the first heading",
			body: "# Title\n\n## Context\n\nThe fallback sentence.\n",
			want: "The fallback sentence.",
		},
		{
			name: "structure is skipped rather than read",
			body: "# Title\n\n## Context\n\n| a | b |\n|---|---|\n\n- item\n* item\n+ item\n\n> quote\n\n<!-- comment -->\n\nThe prose.\n",
			want: "The prose.",
		},
		{
			name: "a fenced block is skipped even when it holds sentences",
			body: "# Title\n\n```go\nfmt.Println(\"a sentence. inside code\")\n```\n\nThe prose.\n",
			want: "The prose.",
		},
		{
			name: "a tilde fence is skipped too",
			body: "# Title\n\n~~~\nnot prose. really\n~~~\n\nThe prose.\n",
			want: "The prose.",
		},
		{
			name: "a terminator inside a code span does not end the sentence",
			body: "# Title\n\nRun `go test ./...` then read the result. Next sentence.\n",
			want: "Run `go test ./...` then read the result.",
		},
		{
			name: "a body with no H1 is read from the top",
			body: "Just prose, no heading.\n",
			want: "Just prose, no heading.",
		},
		{
			name: "a link is reduced to its label",
			body: "# Title\n\nA reader reaches [the thread](../g/G-001.md) in one hop.\n",
			want: "A reader reaches the thread in one hop.",
		},
		{
			name: "a link inside a code span is reduced too",
			body: "# Title\n\nThe row reads `- [<id>](<file>)` today.\n",
			want: "The row reads `- <id>` today.",
		},
		{
			name: "a sentence with no terminator is taken whole",
			body: "# Title\n\nNo terminator here\n",
			want: "No terminator here",
		},
		{
			name: "nothing readable yields nothing",
			body: "# Title\n\n## Context\n\n- only a list\n",
			want: "",
		},
		// Structure the first implementation read as prose. An ordered list item
		// is a list item, a setext heading is a heading, and an indented block is
		// code, however little the leading spaces look like syntax.
		{
			name: "an ordered list item is a list item",
			body: "# Title\n\n1. first item. second\n\nThe prose.\n",
			want: "The prose.",
		},
		{
			name: "an ordered list item written with a paren is one too",
			body: "# Title\n\n1) first item. second\n\nThe prose.\n",
			want: "The prose.",
		},
		{
			name: "a setext heading is a heading, underline included",
			body: "# Title\n\nContext\n-------\n\nThe prose.\n",
			want: "The prose.",
		},
		{
			name: "a setext H1 is the title, not the description",
			body: "Title\n=====\n\nThe prose.\n",
			want: "The prose.",
		},
		{
			name: "an indented code block is code",
			body: "# Title\n\n    indented code. not prose\n\nThe prose.\n",
			want: "The prose.",
		},
		{
			name: "a tab-indented code block is code",
			body: "# Title\n\n\ttabbed code. not prose\n\nThe prose.\n",
			want: "The prose.",
		},
		{
			name: "an HTML block is not prose",
			body: "# Title\n\n<div class=\"x\">\n\nThe prose.\n",
			want: "The prose.",
		},
		{
			name: "a closing HTML tag is not prose either",
			body: "# Title\n\n</section>\n\nThe prose.\n",
			want: "The prose.",
		},
		// The HTML test is a tag, not a bare "<": a paragraph opening with an
		// autolink or a comparison is prose and keeps its description.
		{
			name: "a lede opening with an autolink survives",
			body: "# Title\n\n<https://example.com> is the source of record.\n",
			want: "<https://example.com> is the source of record.",
		},
		{
			name: "a lede opening with a comparison survives",
			body: "# Title\n\n<20% of runs reach the slow path.\n",
			want: "<20% of runs reach the slow path.",
		},
		// A body whose last line is its H1 leaves the cursor at the end. Reading
		// from the top there would seed prose written above the title.
		{
			name: "prose above a trailing H1 is not the description",
			body: "Prose above the title.\n\n# T",
			want: "",
		},
		// A `#` inside a fence is a shell comment or sample markdown, not this
		// artifact's title. Reading it as one seeds the next line of the sample
		// and flips fence parity, hiding the real prose for the rest of the body.
		{
			name: "a hash inside a fence before the heading is not the title",
			body: "```bash\n# install the tool\nclue init .\n```\n\n## Context\n\nThe real prose.\n",
			want: "The real prose.",
		},
		{
			name: "a fence before the real H1 does not swallow the lede",
			body: "```\n# not a title\n```\n\n# G-002 — Real title\n\nThe real lede.\n",
			want: "The real lede.",
		},
		{
			name: "a wrapped lede is read as one paragraph",
			body: "# Title\n\nThe first half of the sentence\ncontinues on the next line. A second sentence.\n",
			want: "The first half of the sentence continues on the next line.",
		},
		{
			name: "a line merely starting with = underlines nothing",
			body: "Not a title\n=1 is the count.\n",
			want: "Not a title",
		},
		{
			name: "a doctype is a declaration, not prose",
			body: "# Title\n\n<!DOCTYPE html>\n\nThe prose.\n",
			want: "The prose.",
		},
		{
			name: "prose merely containing inline HTML is still prose",
			body: "# Title\n\n<b>bold</b> text is prose.\n",
			want: "<b>bold</b> text is prose.",
		},
		{
			name: "a thematic break carries nothing",
			body: "# Title\n\n***\n\nThe prose.\n",
			want: "The prose.",
		},
		// A seed that cannot be made safe is declined, because the shorter row is
		// always available and an unsafe row breaks the corpus.
		{
			name: "a link label carrying its own brackets defeats reduction, so the seed is declined",
			body: "# Title\n\nA reader follows [the [Unreleased] section](../../CHANGELOG.md) to see what shipped.\n",
			want: "",
		},
		{
			name: "a bare link fragment in prose is declined",
			body: "# Title\n\nThe path reads ](file.md) when written by hand.\n",
			want: "",
		},
		{
			name: "an index marker is declined rather than written into the block",
			body: "# Title\n\nEverything before the <!-- clue:index:end --> line is regenerated.\n",
			want: "",
		},
		{
			name: "an empty body yields nothing",
			body: "",
			want: "",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := describeBody(tc.body)
			if tc.want == "" {
				if ok {
					t.Fatalf("expected no sentence, got %q", got)
				}
				return
			}
			if !ok {
				t.Fatalf("expected %q, got nothing", tc.want)
			}
			if got != tc.want {
				t.Fatalf("want %q, got %q", tc.want, got)
			}
		})
	}
}

// A sentence past the bound is cut on a word boundary and marked, so the
// identity a reader scans for cannot be pushed off the line.
func TestUnit_DescribeBodyBoundsTheSentence(t *testing.T) {
	long := strings.TrimSpace(strings.Repeat("elaboration ", 40)) + " tail."
	got, ok := describeBody("# Title\n\n" + long + "\n")
	if !ok {
		t.Fatal("expected a sentence")
	}
	if len(got) > descriptionLimit+len("…") {
		t.Fatalf("sentence exceeds the bound at %d bytes: %q", len(got), got)
	}
	if !strings.HasSuffix(got, "…") {
		t.Fatalf("a cut sentence is marked, got %q", got)
	}
	if strings.Contains(got, "tail") {
		t.Fatalf("the cut drops the tail, got %q", got)
	}
	if strings.HasSuffix(strings.TrimSuffix(got, "…"), "elaborat") {
		t.Fatalf("the cut falls on a word boundary, got %q", got)
	}
}

// Ideographic prose carries no ASCII space, so it reaches the byte-offset
// fallback cut. The result must still be valid UTF-8: a README sliced through
// the middle of a rune is a corrupt file, not a short description.
func TestUnit_DescribeBodyCutsOnARuneBoundary(t *testing.T) {
	body := "# Title\n\n" + strings.Repeat("説明", 200) + "。\n"
	got, ok := describeBody(body)
	if !ok {
		t.Fatal("expected a bounded sentence")
	}
	if !utf8.ValidString(got) {
		t.Fatalf("the cut produced invalid UTF-8: %q", got)
	}
	if !strings.HasSuffix(got, "…") {
		t.Fatalf("a cut sentence is marked, got %q", got)
	}
}

// RowDescription reads through frontmatter, and reports nothing rather than
// guessing when the file cannot be read or carries none.
func TestUnit_RowDescriptionReadsThroughFrontmatter(t *testing.T) {
	dir := t.TempDir()
	withFM := filepath.Join(dir, "a.md")
	if err := os.WriteFile(withFM, []byte("---\nid: G-001\ntype: goal\nstatus: accepted\nlinks: []\ntitle: A\n---\n\n# A\n\nThe sentence.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	got, ok := RowDescription(withFM)
	if !ok || got != "The sentence." {
		t.Fatalf("want %q, got %q (ok=%v)", "The sentence.", got, ok)
	}
	noFM := filepath.Join(dir, "b.md")
	if err := os.WriteFile(noFM, []byte("# B\n\nProse with no frontmatter.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got, ok := RowDescription(noFM); ok {
		t.Fatalf("a file with no frontmatter is not an artifact; got %q", got)
	}
	if got, ok := RowDescription(filepath.Join(dir, "missing.md")); ok {
		t.Fatalf("an unreadable file yields nothing; got %q", got)
	}
}
