package corpus

import (
	"os"
	"path"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"
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

// orderedItemRe matches an ordered list item, which is a list item as much as a
// bulleted one is.
var orderedItemRe = regexp.MustCompile(`^\d+[.)]\s`)

// setextUnderlineRe matches the run of = or - that makes the line above it a
// heading. A body may title itself that way instead of with an ATX #.
var setextUnderlineRe = regexp.MustCompile(`^(=+|-+)$`)

// thematicBreakRe matches a horizontal rule, which carries no prose at all.
var thematicBreakRe = regexp.MustCompile(`^(-{3,}|\*{3,}|_{3,})$`)

// setextH1Re matches the run of = that makes the line above it the title. The
// full match matters: a following line merely starting with = — `=1 is the
// count.` — underlines nothing.
var setextH1Re = regexp.MustCompile(`^=+$`)

// htmlBlockRe matches a line that is nothing but an HTML tag, which is the
// shape that opens a block. Requiring the tag to end the line keeps prose that
// merely contains inline HTML — `<b>bold</b> text is prose.` — readable, and
// requiring a tag name keeps an autolink (`<https://…>`) or a comparison
// (`<20% of runs`) readable too. A declaration or comment is handled by the
// `<!` test beside it.
var htmlBlockRe = regexp.MustCompile(`^</?[a-zA-Z][a-zA-Z0-9-]*[^>]*>$`)

// describeBody is RowDescription's reading, split out so it can be exercised
// on prose directly.
func describeBody(body string) (string, bool) {
	lines := strings.Split(body, "\n")

	// Find the H1 and read after it. A setext H1 — the title underlined with
	// = — counts, or a body that titles itself that way would have its own
	// title read back as its description.
	//
	// Whether the title was found is tracked explicitly rather than inferred
	// from the cursor: a body whose last line is its H1 leaves the cursor at
	// the end too, and restarting from the top there would seed prose written
	// above the title.
	//
	// The search tracks fences, because a `#` line inside one is a shell
	// comment or sample markdown and not this artifact's title. Reading it as
	// the title both seeds the next line of the sample and flips fence parity,
	// which then hides the real prose for the rest of the body.
	i, titled, fenced := 0, false, false
	for i < len(lines) {
		s := strings.TrimSpace(lines[i])
		if strings.HasPrefix(s, "```") || strings.HasPrefix(s, "~~~") {
			fenced = !fenced
			i++
			continue
		}
		if !fenced {
			if strings.HasPrefix(s, "# ") {
				i, titled = i+1, true
				break
			}
			if s != "" && i+1 < len(lines) && setextH1Re.MatchString(strings.TrimSpace(lines[i+1])) {
				i, titled = i+2, true
				break
			}
		}
		i++
	}
	if !titled {
		i, fenced = 0, false // no heading anywhere: the body is read from the top
	}

	for ; i < len(lines); i++ {
		raw := lines[i]
		s := strings.TrimSpace(raw)
		if strings.HasPrefix(s, "```") || strings.HasPrefix(s, "~~~") {
			fenced = !fenced
			continue
		}
		if fenced || s == "" {
			continue
		}
		// An indented code block is code, and its four leading spaces are the
		// only thing that says so.
		if strings.HasPrefix(raw, "    ") || strings.HasPrefix(raw, "\t") {
			continue
		}
		// An ATX heading is skipped rather than read: passing one is what turns
		// a missing lede into the fallback reading, and both land here. A setext
		// heading is the same thing written differently, so the underline and
		// the line it underlines are both skipped.
		if strings.HasPrefix(s, "#") || setextUnderlineRe.MatchString(s) {
			continue
		}
		if i+1 < len(lines) && setextUnderlineRe.MatchString(setextUnderline(lines[i+1])) {
			continue
		}
		if descriptionStructure(s, lines, i) {
			continue
		}

		// Markdown paragraphs commonly wrap at the author's margin. Join the
		// complete paragraph before taking its first sentence; otherwise a line
		// break without punctuation would seed a fragment and discard the rest
		// of the sentence.
		paragraph := []string{s}
		for j := i + 1; titled && j < len(lines); j++ {
			nextRaw := lines[j]
			next := strings.TrimSpace(nextRaw)
			if next == "" || strings.HasPrefix(next, "```") || strings.HasPrefix(next, "~~~") ||
				strings.HasPrefix(nextRaw, "    ") || strings.HasPrefix(nextRaw, "\t") ||
				descriptionStructure(next, lines, j) {
				break
			}
			paragraph = append(paragraph, next)
		}
		if sentence, ok := firstSentence(stripLinks(strings.Join(paragraph, " "))); ok {
			return sentence, true
		}
	}
	return "", false
}

// descriptionStructure reports Markdown block structure that cannot serve as
// the opening prose paragraph. The caller handles indentation and fences
// separately because their raw prefixes determine whether the line is code.
func descriptionStructure(s string, lines []string, i int) bool {
	if strings.HasPrefix(s, "#") || setextUnderlineRe.MatchString(s) {
		return true
	}
	if i+1 < len(lines) && setextUnderlineRe.MatchString(strings.TrimSpace(lines[i+1])) {
		return true
	}
	return strings.HasPrefix(s, "|") || strings.HasPrefix(s, ">") ||
		strings.HasPrefix(s, "<!") || htmlBlockRe.MatchString(s) ||
		strings.HasPrefix(s, "- ") || strings.HasPrefix(s, "* ") || strings.HasPrefix(s, "+ ") ||
		orderedItemRe.MatchString(s) || thematicBreakRe.MatchString(s)
}

// setextUnderline normalizes a line for the underline test.
func setextUnderline(line string) string { return strings.TrimSpace(line) }

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
//
// It returns false rather than a sentence it cannot make safe. A seed is worth
// only what a row is worth, and the shorter row AC-097 already blesses is
// always available, so anything that could make index generation emit a corpus
// the judge rejects is declined instead of repaired.
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
		// Prose with no ASCII space inside the bound — ideographic writing, for
		// one — reaches the fallback cut, which is a byte offset. Backing off to
		// a rune start keeps the README valid UTF-8.
		for cut > 0 && !utf8.RuneStart(s[cut]) {
			cut--
		}
		if cut == 0 {
			return "", false
		}
		s = strings.TrimRight(s[:cut], " ,;:") + "…"
	}
	if !seedIsSafe(s) {
		return "", false
	}
	return s, true
}

// seedIsSafe rejects a candidate description that would corrupt the block it is
// written into. Both classes are closed by refusing the seed rather than by
// chasing the shapes that produce it, because a reduction that merely looks
// right is how the first attempt shipped a link inside a code span.
//
//   - Residual link syntax. stripLinks reduces a well-formed link, but a label
//     carrying its own brackets defeats any single pattern, and checkIndexes
//     reads the leftover `](target)` anyway. A live target then either resolves
//     to nothing, which fails the judge, or resolves to a sibling, which lets
//     the row cover a second artifact so that artifact's own row can go missing
//     unnoticed.
//   - A comment or index marker. regenIndex finds the block by searching for
//     the marker text, so a seed carrying one truncates the block on the next
//     run, duplicates rows, and strands prose outside the markers — data loss in
//     a file the adopter owns.
func seedIsSafe(s string) bool {
	if strings.Contains(s, "](") || mdLinkTextRe.MatchString(s) {
		return false
	}
	return !strings.Contains(s, "<!--") && !strings.Contains(s, "-->")
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
// status badge. A row with no badge is not this population's business: the
// generator always writes one, so a row without it is adopter prose in a shape
// no release produced. Where its label also restates the target's filename,
// IndexRowBacklog counts it; where it does not, no population counts it, which
// is deliberate and documented on IndexDescriptionBacklog below.
//
// Disjointness does not rest on the badge — a filler row can be hand-given one
// while ADR-041's backlog is worked down. It rests on the label-versus-filename
// test inside the loop.
var rowTailRe = regexp.MustCompile("\\)\\s*·\\s*`[^`]*`\\s*(.*)$")

// IndexDescriptionBacklog reports rows that state their record but carry no
// description across the taxonomy READMEs.
//
// Like IndexRowBacklog this is a counted population and never an Issue
// (ADR-046): almost every such row was written by the generator before a
// description was seeded, in a file the adopter owns, and no command clears the
// count because regeneration preserves any row whose target still exists. The
// rest are artifacts whose body holds no readable opening sentence, where the
// generator still emits the shorter row by design.
//
// A row that states a record without the status badge the generator always
// writes is adopter prose in a shape no release produced, so it is left
// uncounted. ADR-041 drew that line first: the count names generated output and
// never grades a curated row.
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
			// A row whose label restates its own filename belongs to ADR-041's
			// filler population, badge or no badge. Working that backlog down by
			// hand-adding a status to an old filler row is the natural
			// intermediate state, and reading the badge alone would count that
			// one row in both populations at once.
			if links[0][1] == strings.TrimSuffix(target, ".md") {
				continue
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
