package corpus

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// The convention register's rules about corpus *state*, promoted from prose to
// checks by CH-110. Each function here is one constraint's detectable subset,
// and every one of them answers from the files in front of it: ADR-044 keeps
// the judge out of history, so a rule about what a change did is enforced by a
// machine that is allowed to have a base, never by this package.

var (
	// A markdown line that opens its own block rather than continuing the
	// paragraph above it (C-001).
	headingRe    = regexp.MustCompile(`^ {0,3}#{1,6}(\s|$)`)
	listMarkerRe = regexp.MustCompile(`^\s*([-*+]|\d+[.)])(\s|$)`)
	underlineRe  = regexp.MustCompile(`^\s*(=+|-+|\*{3,}|_{3,})\s*$`)
	linkDefRe    = regexp.MustCompile(`^\s*\[[^\]]+\]:\s`)

	// A table's delimiter row, with or without the outer pipes. It is what
	// makes the line above it a header rather than a paragraph, so a table
	// written without leading pipes is still a table (C-001).
	tableDelimRe = regexp.MustCompile(`^\s*\|?\s*:?-+:?\s*(\|\s*:?-+:?\s*)+\|?\s*$`)
	// A fence's opening marker, captured so a longer fence can hold a shorter
	// one as an example without the inner one closing the outer.
	fenceOpenRe = regexp.MustCompile("^\\s*(`{3,}|~{3,})")
	// An HTML block's opener: a tag, a comment, or a processing instruction.
	// It deliberately does not match `<https://example.com>` — an autolink is
	// prose, and reading one as a block would exempt the paragraph it sits in.
	htmlOpenRe    = regexp.MustCompile(`^\s*</?[A-Za-z][A-Za-z0-9-]*(\s|/?>|$)`)
	htmlCommentRe = regexp.MustCompile(`^\s*<!--`)

	skippedTaskRe = regexp.MustCompile(`^\s*[-*+] \[-\]\s*(.*)$`)
	planItemRe    = regexp.MustCompile(`^[PM]-\d+$`)
	imageLinkRe   = regexp.MustCompile(`!\[[^\]]*\]\(`)
	codeSpanRe    = regexp.MustCompile("`[^`]*`")

	imageExts = []string{".png", ".jpg", ".jpeg", ".gif", ".svg", ".webp", ".bmp", ".tif", ".tiff", ".ico"}

	// MilestoneStatuses is the vocabulary a plan table's status cell may use
	// (C-010, decision log 2026-08-04).
	MilestoneStatuses = []string{"todo", "doing", "done", "dropped"}
)

// checkProseLayout flags a paragraph or list item broken across lines (C-001).
// Wrapping is the reader's IDE concern, and a wrapped paragraph makes every
// later edit a multi-line diff over text that did not change.
//
// The rule is one sentence long: two text lines in a row are one paragraph
// someone broke. Everything else on the page opens its own block — headings,
// list markers, table rows, fences, blockquotes, HTML, link definitions — and
// a line inside a fence or a frontmatter block is not prose at all.
func checkProseLayout(c *Corpus) []Issue {
	var issues []Issue
	for _, p := range c.MDFiles {
		lines := strings.Split(bodyOf(c.Contents[p]), "\n")
		var blocks blockScanner
		prevIsText := false
		for i, line := range lines {
			// A table's delimiter row identifies the line above it as a header
			// rather than a paragraph this one continues, so the retraction has
			// to reach back one line.
			if !blocks.inVerbatim() && tableDelimRe.MatchString(line) {
				blocks.openTable()
			}
			if blocks.next(line) {
				prevIsText = false
				continue
			}
			// A list marker opens running text without being a continuation
			// itself, so it can be broken across lines exactly as a paragraph
			// can — and is the more common way to do it by accident.
			isText := isProseLine(line)
			if isText && prevIsText {
				issues = append(issues, Issue{p, lineRef(i+1) + ": paragraph continues on a new line — one line per paragraph and per list item (C-001)"})
				// One finding per paragraph: a whole wrapped page is one
				// problem to fix, not fifty.
				isText = false
			}
			prevIsText = isText || listMarkerRe.MatchString(line)
		}
	}
	return issues
}

// blockScanner tracks the markdown blocks whose interior line breaks are
// structure rather than wrapping: fenced code, indented code, HTML, and
// tables. A fence records its own marker and run length, so a longer fence
// holds a shorter one as an example without the inner one closing it — a page
// documenting markdown is exactly where these checks would otherwise invert
// and skip the rest of the file in silence.
type blockScanner struct {
	fence     string // the open fence's marker run, empty when outside one
	inComment bool   // an HTML comment, which ends at its own `-->`
	inHTML    bool   // an HTML tag block, which ends at a blank line
	inIndent  bool
	inTable   bool
	inList    bool // a list is open, so an indented line is a continuation
	prevBlank bool
}

// inVerbatim reports whether the scanner is inside a block whose lines are not
// read as markdown at all.
func (b *blockScanner) inVerbatim() bool { return b.fence != "" || b.inIndent }

func (b *blockScanner) openTable() { b.inTable = true }

// next advances the scanner over one line and reports whether that line belongs
// to a block rather than to running prose. It remembers whether the previous
// line was blank itself: that fact decides whether an indented line opens a
// code block, and a caller that got it wrong would silently stop checking.
func (b *blockScanner) next(line string) bool {
	inBlock := b.consume(line)
	b.prevBlank = strings.TrimSpace(line) == ""
	return inBlock
}

func (b *blockScanner) consume(line string) bool {
	if b.fence != "" {
		if m := fenceOpenRe.FindStringSubmatch(line); m != nil && closesFence(b.fence, m[1]) {
			b.fence = ""
		}
		return true
	}
	// A comment ends where it says it ends. Closing it at the next blank line
	// instead would exempt everything under an index marker — which is every
	// generated index row in every corpus this ships to.
	if b.inComment {
		if strings.Contains(line, "-->") {
			b.inComment = false
		}
		return true
	}
	if m := fenceOpenRe.FindStringSubmatch(line); m != nil {
		b.fence = m[1]
		return true
	}
	if strings.TrimSpace(line) == "" {
		// A blank line ends every block that has no closing marker of its own,
		// and ends the list whose items an indented line would continue.
		b.inHTML, b.inIndent, b.inTable = false, false, false
		return false
	}
	if b.inHTML || b.inTable || b.inIndent {
		return true
	}
	if htmlCommentRe.MatchString(line) {
		b.inComment = !strings.Contains(line, "-->")
		return true
	}
	if htmlOpenRe.MatchString(line) {
		b.inHTML = true
		return true
	}
	if strings.HasPrefix(strings.TrimSpace(line), "|") {
		b.inTable = true
		return true
	}
	// Indented code opens after a blank line — but not inside a list, where an
	// indented line is the item's own continuation or a nested item. CommonMark
	// agrees: at four spaces under a `- ` item, that text is prose, and prose
	// broken across lines is exactly what C-001 is about.
	//
	// This runs before the list-marker test on purpose. An indented example of
	// a list item is a list item in shape, and reading it as one is how a
	// documented example of what not to write fails as the thing itself.
	if b.prevBlank && strings.HasPrefix(line, "    ") && !b.inList {
		b.inIndent = true
		return true
	}
	if listMarkerRe.MatchString(line) {
		b.inList = true
		return false
	}
	if !strings.HasPrefix(line, " ") {
		b.inList = false
	}
	return false
}

// closesFence reports whether a fence marker closes an open one: same
// character, and at least as long.
func closesFence(open, marker string) bool {
	return len(marker) >= len(open) && marker[0] == open[0]
}

// isProseLine reports whether a line continues or starts running text, as
// opposed to opening a block of its own. Block interiors never reach it — the
// scanner has already consumed them.
func isProseLine(line string) bool {
	t := strings.TrimSpace(line)
	switch {
	case t == "":
		return false
	case headingRe.MatchString(line), listMarkerRe.MatchString(line), underlineRe.MatchString(line), linkDefRe.MatchString(line):
		return false
	case strings.HasPrefix(t, ">"):
		return false
	}
	return true
}

// bodyOf returns the text after a leading frontmatter block, keeping the line
// numbering by blanking the block rather than removing it.
func bodyOf(text string) string {
	if !strings.HasPrefix(text, "---\n") {
		return text
	}
	end := strings.Index(text[4:], "\n---")
	if end < 0 {
		return text
	}
	head := text[:4+end+len("\n---")]
	return strings.Repeat("\n", strings.Count(head, "\n")) + text[len(head):]
}

// lineRef renders the 1-based line a reader opens to.
func lineRef(line int) string { return "line " + strconv.Itoa(line) }

// checkSkippedTasks flags a `[-]` task carrying no reason (C-003). The prose
// is the whole point of the mark: an unexplained skip is indistinguishable
// from an abandoned one, and the digest precondition cannot be judged without
// it. Whether the prose is a genuine reason stays with the reader.
func checkSkippedTasks(c *Corpus) []Issue {
	var issues []Issue
	for _, p := range c.MDFiles {
		if !strings.HasPrefix(p, "changes/") || !strings.HasSuffix(p, "/tasks.md") {
			continue
		}
		var blocks blockScanner
		for i, line := range strings.Split(c.Contents[p], "\n") {
			if blocks.next(line) {
				continue
			}
			m := skippedTaskRe.FindStringSubmatch(line)
			if m != nil && strings.TrimSpace(m[1]) == "" {
				issues = append(issues, Issue{p, lineRef(i+1) + ": skipped task carries no reason (C-003)"})
			}
		}
	}
	return issues
}

// checkProposalPlanItem requires a proposal to name the plan item it serves or
// to declare itself plan-less (C-005). A change that serves nothing planned is
// a legitimate change; one that never says so is an unanswered question about
// why the work is being done.
func checkProposalPlanItem(c *Corpus) []Issue {
	var issues []Issue
	for _, a := range c.Artifacts {
		if !strings.HasPrefix(a.Path, "changes/") || !strings.HasSuffix(a.Path, "/proposal.md") {
			continue
		}
		if strings.Contains(a.Body, "plan-less") {
			continue
		}
		planned := false
		for _, l := range a.Links {
			if planItemRe.MatchString(l) {
				planned = true
				break
			}
		}
		if !planned {
			issues = append(issues, Issue{a.Path, "proposal names no plan item: link a P-xxx or M-xxx, or declare the change plan-less (C-005)"})
		}
	}
	return issues
}

// checkInlineDiagrams flags images in the corpus (C-007). A diagram is a
// Mermaid block in the markdown it illustrates: versionable, diffable, and
// rendered wherever the corpus is read. A binary or externally hosted picture
// is none of those, and it goes stale without ever showing a diff.
func checkInlineDiagrams(c *Corpus) []Issue {
	var issues []Issue
	for _, p := range c.MDFiles {
		if !strings.HasPrefix(p, "docs/") {
			continue
		}
		var blocks blockScanner
		for i, line := range strings.Split(c.Contents[p], "\n") {
			if blocks.next(line) {
				continue
			}
			if imageLinkRe.MatchString(codeSpanRe.ReplaceAllString(line, "")) {
				issues = append(issues, Issue{p, lineRef(i+1) + ": image link — diagrams in the corpus are inline Mermaid (C-007)"})
			}
		}
	}
	for _, p := range c.OtherFiles {
		if !strings.HasPrefix(p, "docs/") {
			continue
		}
		lower := strings.ToLower(p)
		for _, ext := range imageExts {
			if strings.HasSuffix(lower, ext) {
				issues = append(issues, Issue{p, "image file under docs/ — diagrams in the corpus are inline Mermaid (C-007)"})
				break
			}
		}
	}
	return issues
}

// checkMilestoneStatus reads the status column of a plan's milestone table and
// holds it to one vocabulary (C-010). The column is found by its header, so a
// plan table that declares no status is not one this rule reads.
func checkMilestoneStatus(c *Corpus) []Issue {
	var issues []Issue
	for _, a := range c.Artifacts {
		if a.Type != "plan" {
			continue
		}
		col := -1
		var blocks blockScanner
		var header []string
		for _, line := range strings.Split(a.Body, "\n") {
			if blocks.fence != "" || fenceOpenRe.MatchString(line) {
				blocks.next(line)
				col, header = -1, nil
				continue
			}
			t := strings.TrimSpace(line)
			// A delimiter row is what makes the line above it a header, with or
			// without outer pipes — the same rule the prose lint applies, so
			// the two checks cannot disagree about what a table is.
			if tableDelimRe.MatchString(t) {
				col = -1
				for i, cell := range header {
					if strings.EqualFold(cell, "status") {
						col = i
						break
					}
				}
				continue
			}
			if !strings.Contains(t, "|") {
				col, header = -1, nil
				continue
			}
			cells := tableCells(t)
			if col < 0 {
				header = cells
				continue
			}
			if col >= len(cells) {
				continue
			}
			v := strings.Trim(strings.TrimSpace(cells[col]), "`")
			if v == "" || contains(MilestoneStatuses, v) {
				continue
			}
			name := "row"
			if len(cells) > 0 && strings.TrimSpace(cells[0]) != "" {
				name = strings.TrimSpace(cells[0])
			}
			issues = append(issues, Issue{a.Path, name + ": milestone status " + v + " is not one of " + strings.Join(MilestoneStatuses, ", ") + " (C-010)"})
		}
	}
	return issues
}

// tableCells splits a markdown table row into its cells. Only an unescaped
// pipe outside a code span divides one: `\\|` is the documented way to put a
// pipe in a cell, and this project's plan tables carry long prose cells full of
// code spans. Splitting on every pipe would shift every column after the first
// such cell, which reads a neighbouring cell as the status and misses the real
// one — a false failure and a silent miss from the same bug.
func tableCells(row string) []string {
	var cells []string
	var cur strings.Builder
	open := 0 // the length of the open code span's backtick run, 0 when outside
	runes := []rune(strings.TrimSpace(row))
	for i := 0; i < len(runes); i++ {
		switch c := runes[i]; {
		case c == '\\' && i+1 < len(runes):
			cur.WriteRune(c)
			i++
			cur.WriteRune(runes[i])
		case c == '`':
			run := 1
			for i+run < len(runes) && runes[i+run] == '`' {
				run++
			}
			switch {
			case open == 0:
				open = run
			case run == open:
				open = 0
			}
			for n := 0; n < run; n++ {
				cur.WriteRune('`')
			}
			i += run - 1
		case c == '|' && open == 0:
			cells = append(cells, strings.TrimSpace(cur.String()))
			cur.Reset()
		default:
			cur.WriteRune(c)
		}
	}
	cells = append(cells, strings.TrimSpace(cur.String()))
	// The outer pipes produce an empty cell at each end; dropping them keeps
	// column indexes the same whether or not a table is written with them.
	if len(cells) > 0 && cells[0] == "" {
		cells = cells[1:]
	}
	if len(cells) > 0 && cells[len(cells)-1] == "" {
		cells = cells[:len(cells)-1]
	}
	return cells
}

func contains(list []string, v string) bool {
	for _, s := range list {
		if s == v {
			return true
		}
	}
	return false
}

// typeFields are the frontmatter extensions a type carries beyond the linted
// core (C-009). A field listed here must be present; `accepted-by` may be an
// empty list, because an unsigned decision declares itself unsigned and that
// declaration is the record ADR-029 wants.
var typeFields = map[string][]string{
	"decision":   {"author", "accepted-by"},
	"capability": {"goal"},
}

// checkTypeFields requires each type's own frontmatter extensions (C-009). A
// type the validator does not recognize carries the core alone: an adopter's
// own artifact types are theirs to extend, not this rule's to reject.
func checkTypeFields(c *Corpus) []Issue {
	var issues []Issue
	for _, a := range c.Artifacts {
		fields, ok := typeFields[a.Type]
		if !ok {
			continue
		}
		var missing []string
		for _, f := range fields {
			v, present := a.Fields[f]
			if !present {
				missing = append(missing, f)
				continue
			}
			if s, isString := v.(string); isString && strings.TrimSpace(s) == "" {
				missing = append(missing, f)
			}
		}
		if len(missing) > 0 {
			sort.Strings(missing)
			issues = append(issues, Issue{a.Path, a.Type + " missing or empty field(s): " + strings.Join(missing, ", ") + " (C-009)"})
		}
	}
	return issues
}
