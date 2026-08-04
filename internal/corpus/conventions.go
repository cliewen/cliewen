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
		inFence, prevIsText := false, false
		for i, line := range lines {
			if fenceRe.MatchString(line) {
				inFence = !inFence
				prevIsText = false
				continue
			}
			if inFence {
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
				prevIsText = false
				continue
			}
			prevIsText = isText || listMarkerRe.MatchString(line)
		}
	}
	return issues
}

// isProseLine reports whether a line continues or starts running text, as
// opposed to opening a block of its own.
func isProseLine(line string) bool {
	t := strings.TrimSpace(line)
	switch {
	case t == "":
		return false
	case headingRe.MatchString(line), listMarkerRe.MatchString(line), underlineRe.MatchString(line), linkDefRe.MatchString(line):
		return false
	case strings.HasPrefix(t, "|"), strings.HasPrefix(t, ">"), strings.HasPrefix(t, "<"):
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
		inFence := false
		for i, line := range strings.Split(c.Contents[p], "\n") {
			if fenceRe.MatchString(line) {
				inFence = !inFence
				continue
			}
			if inFence {
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
		inFence := false
		for i, line := range strings.Split(c.Contents[p], "\n") {
			if fenceRe.MatchString(line) {
				inFence = !inFence
				continue
			}
			if inFence {
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
		col, inTable, inFence := -1, false, false
		for _, line := range strings.Split(a.Body, "\n") {
			if fenceRe.MatchString(line) {
				inFence = !inFence
				continue
			}
			t := strings.TrimSpace(line)
			if inFence || !strings.HasPrefix(t, "|") {
				col, inTable = -1, false
				continue
			}
			cells := tableCells(t)
			if col < 0 {
				for i, cell := range cells {
					if strings.EqualFold(cell, "status") {
						col, inTable = i, false
						break
					}
				}
				continue
			}
			if !inTable {
				// The delimiter row directly under the header.
				inTable = true
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

// tableCells splits a markdown table row into its cells, dropping the empty
// strings the leading and trailing pipes produce.
func tableCells(row string) []string {
	parts := strings.Split(strings.Trim(row, "|"), "|")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
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
