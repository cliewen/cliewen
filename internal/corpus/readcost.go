package corpus

import (
	"sort"
	"strings"
)

// DefaultContextSliceBudget is the maximum number of durable artifacts a
// default (one-hop) context slice should print before validate reports it as a
// read-cost backlog item. It is a reporting threshold, never a validation
// failure: a reader can still widen a slice when the work needs it.
const DefaultContextSliceBudget = 8

// MultiDocumentArtifact is one durable artifact whose Markdown body holds
// more than one rendered H1 document.
type MultiDocumentArtifact struct {
	Artifact  *Artifact
	Documents int
}

// ContextSliceOverBudget is one durable identity whose default context slice
// prints more artifacts than DefaultContextSliceBudget permits.
type ContextSliceOverBudget struct {
	Identity  string
	Artifacts int
}

// MultiDocumentBacklog reports durable, non-completed artifacts whose bodies
// contain more than one H1 outside fenced examples. An H1 begins a document in
// the rendered corpus; a second H1 therefore has a distinct primary reader to
// serve. Completed plans remain historical records and are excluded by state.
func MultiDocumentBacklog(c *Corpus) []MultiDocumentArtifact {
	var out []MultiDocumentArtifact
	for _, a := range c.Artifacts {
		if !strings.HasPrefix(a.Path, "docs/") || a.Status == "completed" {
			continue
		}
		if documents := renderedDocumentCount(a.Body); documents > 1 {
			out = append(out, MultiDocumentArtifact{Artifact: a, Documents: documents})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Artifact.Path < out[j].Artifact.Path })
	return out
}

// ContextSliceBudgetBacklog reports every durable identity whose ordinary
// one-hop context read crosses the stated artifact budget. Invalid or
// ambiguous identities are the validator's existing concern, so they are not
// rendered as a second read-cost finding. Identities owned by a completed plan
// are excluded on the same ground as MultiDocumentBacklog: their links are part
// of a historical record C-008 makes unavailable for rewriting, so reporting
// them would name a backlog no change is permitted to repair.
func ContextSliceBudgetBacklog(c *Corpus) []ContextSliceOverBudget {
	owners := contextOwners(c)
	ids := make([]string, 0, len(owners))
	for id, artifacts := range owners {
		if len(artifacts) != 1 || artifacts[0].Status == "completed" {
			continue
		}
		if strings.HasPrefix(artifacts[0].Path, "docs/") {
			ids = append(ids, id)
		}
	}
	sort.Strings(ids)

	var out []ContextSliceOverBudget
	for _, id := range ids {
		root, err := owners.resolve(id)
		if err != nil {
			continue
		}
		if artifacts := defaultContextSliceSize(owners, root); artifacts > DefaultContextSliceBudget {
			out = append(out, ContextSliceOverBudget{Identity: id, Artifacts: artifacts})
		}
	}
	return out
}

// defaultContextSliceSize is the cost validate reports: the root and the
// unique resolvable artifacts it links directly. It intentionally does not
// call Context, whose frontier count walks the remaining graph in order to
// describe it to a reader; the validator needs only the bounded slice itself.
func defaultContextSliceSize(owners contextIndex, root *Artifact) int {
	seen := map[string]bool{root.Path: true}
	for _, link := range root.Links {
		if target, err := owners.resolve(link); err == nil {
			seen[target.Path] = true
		}
	}
	return len(seen)
}

// renderedDocumentCount counts H1 headings in an artifact body while skipping
// fenced examples. A body with no H1 is still one document for the purpose of
// the convention, but cannot be a multi-document artifact. Both heading forms
// Markdown renders as an H1 count: an ATX `# ` line and a setext title over a
// `=` underline.
func renderedDocumentCount(body string) int {
	count := 0
	fence := ""
	prevIsSetextTitle := false
	for _, line := range strings.Split(body, "\n") {
		trimmed := strings.TrimLeft(line, " ")
		indented := len(line)-len(trimmed) > 3

		if marker, closes := fenceMarker(trimmed); marker != "" && !indented {
			switch {
			case fence == "":
				fence = marker
			case closes && marker[0] == fence[0] && len(marker) >= len(fence):
				fence = ""
			}
			prevIsSetextTitle = false
			continue
		}
		if fence != "" || indented {
			prevIsSetextTitle = false
			continue
		}

		switch {
		case strings.HasPrefix(trimmed, "# "):
			count++
			prevIsSetextTitle = false
		case prevIsSetextTitle && isSetextH1Underline(trimmed):
			count++
			prevIsSetextTitle = false
		default:
			prevIsSetextTitle = canCarrySetextUnderline(trimmed)
		}
	}
	return count
}

// fenceMarker returns the leading run of fence characters on a line, and
// whether the line could close a fence — a closer carries no info string.
// Tracking the character and the run length keeps a longer outer fence from
// being closed by the shorter fence it demonstrates, which would invert the
// in-fence state and hide every heading after it.
func fenceMarker(trimmed string) (marker string, closes bool) {
	if !strings.HasPrefix(trimmed, "```") && !strings.HasPrefix(trimmed, "~~~") {
		return "", false
	}
	run := 0
	for run < len(trimmed) && trimmed[run] == trimmed[0] {
		run++
	}
	return trimmed[:run], strings.TrimSpace(trimmed[run:]) == ""
}

// isSetextH1Underline reports whether a line is a run of `=`, which underlines
// the preceding line into a rendered H1. `-` underlines an H2 and is also the
// frontmatter delimiter, so it is deliberately not recognized here.
func isSetextH1Underline(trimmed string) bool {
	trimmed = strings.TrimRight(trimmed, " ")
	return trimmed != "" && strings.Trim(trimmed, "=") == ""
}

// canCarrySetextUnderline reports whether a line is ordinary text that a `=`
// underline would turn into an H1. A blank line, a heading, and the block
// starters that open their own construct cannot become a setext title.
func canCarrySetextUnderline(trimmed string) bool {
	if trimmed == "" {
		return false
	}
	switch trimmed[0] {
	case '#', '>', '-', '*', '+', '|', '=':
		return false
	}
	return true
}
