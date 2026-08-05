package corpus

import (
	"os"
	"path"
	"regexp"
	"sort"
	"strings"
)

// IndexRowIdentity is what a taxonomy index row states about the artifact it
// links: the record's own id and title, and the status its frontmatter
// declares.
type IndexRowIdentity struct {
	ID     string
	Title  string
	Status string
}

// RowIdentity reads the identity an index row states for the artifact file at
// p. The values come from parsed frontmatter and never from raw lines: a title
// whose value contains a colon is YAML-quoted in the file, and a row spells the
// value, not the quoting.
//
// All three of id, title, and status must be readable, or this returns false
// and the caller emits the plain link. A row is one shape or the other, never a
// third carrying an empty status badge — an artifact missing a core field is
// the judge's to report, not index generation's to half-render.
func RowIdentity(p string) (IndexRowIdentity, bool) {
	raw, err := os.ReadFile(p)
	if err != nil {
		return IndexRowIdentity{}, false
	}
	fields, _, ok, err := parseFrontmatter(string(raw))
	if err != nil || !ok {
		return IndexRowIdentity{}, false
	}
	id, _ := fields["id"].(string)
	title, _ := fields["title"].(string)
	status, _ := fields["status"].(string)
	if id == "" || title == "" || status == "" {
		return IndexRowIdentity{}, false
	}
	return IndexRowIdentity{ID: id, Title: title, Status: status}, true
}

// descriptionLimit bounds the seeded sentence. A row is one line, because a
// continuation line is not an entry checkIndexes recognizes, and an unbounded
// first sentence would push the identity a reader scans for off the screen.
const descriptionLimit = 200

// RowDescription reads the sentence an appended index row seeds for the
// artifact file at p: a lede paragraph directly beneath the H1 where one
// exists, otherwise the first sentence of the first paragraph under the first
// heading.
//
// It is a seed and never an assertion (ADR-046). The extractor cannot tell a
// purpose statement from a problem statement, because both are declarative
// prose in the same position: a capability README opens with what it delivers,
// while a decision opens under "Context and problem statement" and states, in
// the present tense, the defect it removed. The author corrects the row, and
// regeneration never rewrites it.
//
// Returns false when no prose sentence can be read. The caller then emits the
// shorter row: a row is one shape or the other, never a third carrying an
// empty tail.
func RowDescription(p string) (string, bool) {
	raw, err := os.ReadFile(p)
	if err != nil {
		return "", false
	}
	_, body, ok, err := parseFrontmatter(string(raw))
	if err != nil || !ok {
		return "", false
	}
	return describeBody(body)
}

// describeBody is RowDescription's reading, split out so it can be exercised
// on prose directly.
func describeBody(body string) (string, bool) {
	lines := strings.Split(body, "\n")
	i := 0
	for i < len(lines) && !strings.HasPrefix(strings.TrimSpace(lines[i]), "# ") {
		i++
	}
	if i < len(lines) {
		i++ // start after the H1; a body without one is read from the top
	} else {
		i = 0
	}
	fenced := false
	for ; i < len(lines); i++ {
		s := strings.TrimSpace(lines[i])
		if strings.HasPrefix(s, "```") || strings.HasPrefix(s, "~~~") {
			fenced = !fenced
			continue
		}
		if fenced || s == "" {
			continue
		}
		// A heading is skipped rather than read: passing one is what turns a
		// missing lede into the fallback reading, and both land here.
		if strings.HasPrefix(s, "#") {
			continue
		}
		// A table row, list item, blockquote, or HTML comment is structure
		// rather than the paragraph this is looking for.
		if strings.HasPrefix(s, "|") || strings.HasPrefix(s, ">") || strings.HasPrefix(s, "<!--") ||
			strings.HasPrefix(s, "- ") || strings.HasPrefix(s, "* ") || strings.HasPrefix(s, "+ ") {
			continue
		}
		if sentence, ok := firstSentence(stripLinks(s)); ok {
			return sentence, true
		}
	}
	return "", false
}

// sentenceEndRe finds a sentence terminator followed by a space. Splitting on
// the terminator alone would cut inside a version number or a filename.
var sentenceEndRe = regexp.MustCompile(`[.!?]\s`)

// mdLinkTextRe matches an inline markdown link so it can be reduced to its
// text.
var mdLinkTextRe = regexp.MustCompile(`\[([^\]]*)\]\([^)\s]*\)`)

// stripLinks reduces every inline markdown link to its label, including one
// written inside a code span.
//
// A seeded description carrying a live link would make its row cover a second
// target: regenIndex credits every link on a line, so a prose reference to
// another artifact would satisfy that artifact's index entry and let a real row
// for it go missing. It would also take the row out of both index-row
// populations, which read only a row carrying exactly one link.
//
// A code span is no exception, even though the link inside it is literal text
// to a markdown reader. checkIndexes reads the block with a link pattern that
// knows nothing about spans, and an artifact quoting the row format it defines
// writes a placeholder target — `<file>` — that resolves to nothing. Preserving
// it would make index generation emit a corpus the judge then rejects, turning
// a Cliewen defect into the adopter's red build, which is the outcome ADR-041
// and ADR-046 both refuse. Faithfully quoting an example is worth less than a
// generator that cannot produce an unresolvable link.
func stripLinks(s string) string {
	return mdLinkTextRe.ReplaceAllString(s, "$1")
}

// firstSentence takes the opening sentence of a paragraph and bounds it.
// Terminators inside backticks are ignored, because a code span carrying a
// path or a command routinely contains a period that ends nothing.
func firstSentence(s string) (string, bool) {
	masked := maskCodeSpans(s)
	if m := sentenceEndRe.FindStringIndex(masked); m != nil {
		s = strings.TrimSpace(s[:m[0]+1])
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return "", false
	}
	if len(s) > descriptionLimit {
		cut := strings.LastIndex(s[:descriptionLimit], " ")
		if cut <= 0 {
			cut = descriptionLimit
		}
		s = strings.TrimRight(s[:cut], " ,;:") + "…"
	}
	return s, true
}

// maskCodeSpans blanks the interior of backtick spans so sentence detection
// cannot terminate inside one, preserving offsets so the caller can slice the
// original string.
func maskCodeSpans(s string) string {
	out := []byte(s)
	inSpan := false
	for i := 0; i < len(out); i++ {
		if out[i] == '`' {
			inSpan = !inSpan
			continue
		}
		if inSpan {
			out[i] = 'x'
		}
	}
	return string(out)
}

// IndexRowFiller is one index row whose label only restates its own link — the
// target's filename with the extension removed. It carries no title and no
// status, so it tells a reader nothing the link did not already say.
type IndexRowFiller struct {
	Readme string // repo-relative taxonomy README carrying the row
	Target string // the row's link target, relative to that README
	Label  string // what the row's label actually says
}

// indexRowLinkRe captures a markdown link's label and target together;
// mdLinkRe deliberately captures only the target and is left as it is.
var indexRowLinkRe = regexp.MustCompile(`\[([^\]]*)\]\(([^)#\s]+)\)`)

// IndexRowBacklog reports generated filler rows across the taxonomy READMEs.
//
// It is a counted population and never an Issue (ADR-041): these rows are the
// tool's own former output sitting in a file the adopter owns, so the judge
// names them and does not fail on them, following the shape ADR-035 set for
// unverified meaning and ADR-017 for the constraint backlog.
//
// Only a row linking one sibling artifact is read. A row pointing at a
// subfolder README states a section rather than a record and has no
// frontmatter identity to carry, and a deliberately curated row covering
// several targets at once is left alone.
func IndexRowBacklog(c *Corpus) []IndexRowFiller {
	var out []IndexRowFiller
	for _, rel := range c.MDFiles {
		if !isTaxonomyReadme(rel) {
			continue
		}
		text := c.Contents[rel]
		start := strings.Index(text, indexStart)
		end := strings.Index(text, indexEnd)
		if start < 0 || end < 0 || end < start {
			continue // checkIndexes already reports missing or malformed markers
		}
		for _, line := range strings.Split(text[start+len(indexStart):end], "\n") {
			line = strings.TrimSuffix(line, "\r")
			links := indexRowLinkRe.FindAllStringSubmatch(line, -1)
			if len(links) != 1 {
				continue
			}
			label, raw := links[0][1], links[0][2]
			if externalRe.MatchString(raw) {
				continue
			}
			target := path.Clean(raw)
			if strings.Contains(target, "/") {
				continue // a subfolder README states a section, not a record
			}
			if label != strings.TrimSuffix(target, ".md") {
				continue
			}
			out = append(out, IndexRowFiller{Readme: rel, Target: target, Label: label})
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Readme != out[j].Readme {
			return out[i].Readme < out[j].Readme
		}
		return out[i].Target < out[j].Target
	})
	return out
}

// IndexRowUndescribed is one index row that states its record — id, title, and
// status — but says nothing about what the artifact contains.
type IndexRowUndescribed struct {
	Readme string // repo-relative taxonomy README carrying the row
	Target string // the row's link target, relative to that README
}

// rowTailRe matches a stated-record row and captures whatever follows its
// status badge. A row with no badge is not this population's business: it
// states no record, so IndexRowBacklog already counts it and the two
// populations stay disjoint.
var rowTailRe = regexp.MustCompile("\\)\\s*·\\s*`[^`]*`\\s*(.*)$")

// IndexDescriptionBacklog reports rows that state their record but carry no
// description across the taxonomy READMEs.
//
// Like IndexRowBacklog this is a counted population and never an Issue
// (ADR-046): every such row was written by the generator before a description
// was seeded, in a file the adopter owns, and no command clears the count
// because regeneration preserves any row whose target still exists.
func IndexDescriptionBacklog(c *Corpus) []IndexRowUndescribed {
	var out []IndexRowUndescribed
	for _, rel := range c.MDFiles {
		if !isTaxonomyReadme(rel) {
			continue
		}
		text := c.Contents[rel]
		start := strings.Index(text, indexStart)
		end := strings.Index(text, indexEnd)
		if start < 0 || end < 0 || end < start {
			continue // checkIndexes already reports missing or malformed markers
		}
		for _, line := range strings.Split(text[start+len(indexStart):end], "\n") {
			line = strings.TrimSuffix(line, "\r")
			links := indexRowLinkRe.FindAllStringSubmatch(line, -1)
			if len(links) != 1 {
				continue // a curated row covering several targets is left alone
			}
			raw := links[0][2]
			if externalRe.MatchString(raw) {
				continue
			}
			target := path.Clean(raw)
			if strings.Contains(target, "/") {
				continue // a subfolder README states a section, not a record
			}
			m := rowTailRe.FindStringSubmatch(line)
			if m == nil {
				continue // no status badge: not a row that states its record
			}
			if strings.TrimSpace(strings.TrimLeft(strings.TrimSpace(m[1]), "—-")) != "" {
				continue
			}
			out = append(out, IndexRowUndescribed{Readme: rel, Target: target})
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Readme != out[j].Readme {
			return out[i].Readme < out[j].Readme
		}
		return out[i].Target < out[j].Target
	})
	return out
}
