// Package importedchange implements the durable imported-change record
// (ADR-050, P-011 M-054): the artifact that replaces the lossy
// milestone-row-plus-draft-capability treatment of a source repository's
// pending change during brownfield extraction. A record's proof-links table
// is the piece the old mapping discarded outright — the trace from a source
// task to the target criterion it proves — and this package parses that
// table and checks whether a record's claimed completion is honest, without
// inventing a second reading of the criteria declarations clue validate
// already trusts (internal/corpus.AcceptanceEvidence remains the one
// authority on whether a criterion exists, is @draft, or is retired).
package importedchange

import (
	"regexp"
	"sort"
	"strings"
)

// Status is the imported-change lifecycle (ADR-050): in-progress declares
// work still pending and may name proof links that are not yet satisfied;
// complete asserts every named proof link is already proven. There is no
// retired value — the type is durable, not transient (ADR-034's exception).
type Status string

const (
	StatusInProgress Status = "in-progress"
	StatusComplete   Status = "complete"
)

// ProofLink is one row of a record's proof-links table: a source task and
// the target criterion ID it claims to prove.
type ProofLink struct {
	Task      string
	Criterion string
}

var (
	proofLinksHeadingRe = regexp.MustCompile(`(?im)^#+\s*proof links\s*$`)
	tableDelimRe        = regexp.MustCompile(`^\s*\|?\s*:?-+:?\s*(\|\s*:?-+:?\s*)+\|?\s*$`)
)

// ParseProofLinks extracts the Task and Criterion columns from the markdown
// table under the body's "## Proof links" heading (any heading level). It
// returns nil when the body has no such heading or the heading's table
// declares neither column — checkImportedChanges in internal/corpus/rules.go
// reports that absence as its own issue rather than this package guessing a
// default.
func ParseProofLinks(body string) []ProofLink {
	// A record may show a proof-links table as an example in a fenced code
	// block. Examples are not the record's own proof, so remove fenced content
	// before looking for the heading or parsing rows.
	body = outsideFencedCode(body)
	loc := proofLinksHeadingRe.FindStringIndex(body)
	if loc == nil {
		return nil
	}
	rest := body[loc[1]:]
	// Stop at the next heading of any level, so a later section's table
	// (e.g. an unrelated "## Dependencies" list) is never read as proof
	// links.
	if next := regexp.MustCompile(`(?m)^#+\s`).FindStringIndex(rest); next != nil {
		rest = rest[:next[0]]
	}

	lines := strings.Split(rest, "\n")
	for i, line := range lines {
		t := strings.TrimSpace(line)
		if t == "" || !strings.Contains(t, "|") || tableDelimRe.MatchString(t) {
			continue
		}
		cells := tableCells(t)
		taskCol, criterionCol := -1, -1
		for j, cell := range cells {
			switch strings.ToLower(strings.TrimSpace(cell)) {
			case "task":
				taskCol = j
			case "criterion":
				criterionCol = j
			}
		}
		if taskCol < 0 || criterionCol < 0 || i+1 == len(lines) {
			continue
		}
		delimiter := strings.TrimSpace(lines[i+1])
		if !tableDelimRe.MatchString(delimiter) || len(tableCells(delimiter)) != len(cells) {
			continue
		}

		var links []ProofLink
		for _, row := range lines[i+2:] {
			row = strings.TrimSpace(row)
			if row == "" || !strings.Contains(row, "|") {
				break
			}
			cells = tableCells(row)
			if taskCol >= len(cells) || criterionCol >= len(cells) {
				continue
			}
			task := strings.Trim(strings.TrimSpace(cells[taskCol]), "`")
			criterion := strings.Trim(strings.TrimSpace(cells[criterionCol]), "`")
			if task == "" && criterion == "" {
				continue
			}
			links = append(links, ProofLink{Task: task, Criterion: criterion})
		}
		return links
	}
	return nil
}

// outsideFencedCode preserves the document's line structure while removing
// fenced code blocks. A fence closes only with the same marker character, at
// least the opening length, and no info string, so a nested example cannot
// make its following lines look like proof links.
func outsideFencedCode(doc string) string {
	lines := strings.Split(doc, "\n")
	var fence string
	for i, line := range lines {
		marker := fenceMarker(line)
		if fence == "" {
			if marker != "" {
				fence = marker
				lines[i] = ""
			}
			continue
		}
		if marker != "" && marker[0] == fence[0] && len(marker) >= len(fence) && closesFence(line, marker) {
			fence = ""
		}
		lines[i] = ""
	}
	return strings.Join(lines, "\n")
}

func fenceMarker(line string) string {
	indent := 0
	for indent < len(line) && line[indent] == ' ' {
		indent++
	}
	if indent > 3 || indent == len(line) || line[indent] == '\t' {
		return ""
	}
	trimmed := line[indent:]
	if len(trimmed) < 3 || (trimmed[0] != '`' && trimmed[0] != '~') {
		return ""
	}
	run := 1
	for run < len(trimmed) && trimmed[run] == trimmed[0] {
		run++
	}
	if run < 3 {
		return ""
	}
	return trimmed[:run]
}

func closesFence(line, marker string) bool {
	trimmed := strings.TrimLeft(line, " \t")
	return strings.TrimSpace(trimmed[len(marker):]) == ""
}

// tableCells splits a markdown table row into its cells, skipping pipes
// inside a code span, and drops the empty cells the outer pipes produce.
func tableCells(row string) []string {
	var cells []string
	var cur strings.Builder
	open := 0
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
	if len(cells) > 0 && cells[0] == "" {
		cells = cells[1:]
	}
	if len(cells) > 0 && cells[len(cells)-1] == "" {
		cells = cells[:len(cells)-1]
	}
	return cells
}

// UnprovenLinks returns, sorted and deduplicated, every criterion ID a
// complete record's proof-links table names that provable does not mark
// true — a criterion that does not exist, is @draft, or is retired. The
// caller (checkImportedChanges) builds provable from the same declaration
// harvest checkACTests and clue parity already trust, so this package never
// reads a criteria.md itself. Call only when the record's status is
// complete: an in-progress record is allowed unproven links by design.
func UnprovenLinks(links []ProofLink, provable map[string]bool) []string {
	seen := map[string]bool{}
	var bad []string
	for _, l := range links {
		if l.Criterion == "" || provable[l.Criterion] {
			continue
		}
		if seen[l.Criterion] {
			continue
		}
		seen[l.Criterion] = true
		bad = append(bad, l.Criterion)
	}
	sort.Strings(bad)
	return bad
}
