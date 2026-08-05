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

	taskCol, criterionCol := -1, -1
	var links []ProofLink
	for _, line := range strings.Split(rest, "\n") {
		t := strings.TrimSpace(line)
		if t == "" {
			continue
		}
		if tableDelimRe.MatchString(t) {
			continue // the header was already read on the previous line
		}
		if !strings.Contains(t, "|") {
			continue
		}
		cells := tableCells(t)
		if taskCol < 0 && criterionCol < 0 {
			for i, cell := range cells {
				switch strings.ToLower(strings.TrimSpace(cell)) {
				case "task":
					taskCol = i
				case "criterion":
					criterionCol = i
				}
			}
			continue
		}
		if taskCol < 0 || criterionCol < 0 || taskCol >= len(cells) || criterionCol >= len(cells) {
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
