package corpus

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The extractor's reading, exercised where it lives. The row-level contract is
// held by AC-096, AC-097, and AC-098 in the scaffold package; these cases pin
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
