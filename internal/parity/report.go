package parity

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/cliewen/cliewen/internal/corpus"
)

// Derived extraction reports (ADR-054, P-012 M-059). A durable extraction
// report states its criterion counts and mapping table inside one delimited
// region rendered from the pinned source manifest the region names, so a
// report cannot claim a population or a mapping the manifest never held.
// Rendering reads the manifest alone and never the live corpus: the report is
// history, and a region re-derived from the present tree would go stale on
// every later change to an unrelated criterion.
const (
	// RegionOpenPrefix and RegionEnd delimit the derived region. The opening
	// marker names the manifest, relative to the repository root, so the
	// report declares its own origin and no second index has to map reports
	// to manifests.
	RegionOpenPrefix = "<!-- clue:derived-from: "
	RegionOpenSuffix = " -->"
	RegionEnd        = "<!-- clue:derived-end -->"
)

// region locates one derived region inside a report's content by byte offset.
type region struct {
	manifest  string // repo-relative manifest path named by the opening marker
	bodyStart int    // first byte after the opening marker's line
	bodyEnd   int    // first byte of the closing marker's line
	body      string
}

// maskCode returns a copy of content, byte-for-byte the same length, with
// everything inside a fenced block or an inline code span replaced by spaces.
// Offsets into the mask are offsets into the original, so a caller can search
// the mask and slice the original. A document that documents the markers —
// this decision record, the skill, the design note — writes them in code
// spans, and a scanner that could not tell an example from a claim would make
// describing the contract impossible (the same rule the outward-reference
// scan already follows).
func maskCode(content string) string {
	out := []byte(content)
	fenced := false
	for lineStart := 0; lineStart < len(out); {
		lineEnd := strings.IndexByte(content[lineStart:], '\n')
		if lineEnd < 0 {
			lineEnd = len(out)
		} else {
			lineEnd += lineStart
		}
		line := content[lineStart:lineEnd]
		trimmed := strings.TrimLeft(line, " \t")
		isFence := strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~")
		switch {
		case isFence:
			fenced = !fenced
			blank(out, lineStart, lineEnd)
		case fenced:
			blank(out, lineStart, lineEnd)
		default:
			maskSpans(out, content, lineStart, lineEnd)
		}
		lineStart = lineEnd + 1
	}
	return string(out)
}

func blank(out []byte, from, to int) {
	for i := from; i < to; i++ {
		out[i] = ' '
	}
}

// maskSpans blanks each inline code span on one line. An unclosed backtick
// run opens no span, exactly as Markdown renders it.
func maskSpans(out []byte, content string, from, to int) {
	i := from
	for i < to {
		if content[i] != '`' {
			i++
			continue
		}
		run := 0
		for i+run < to && content[i+run] == '`' {
			run++
		}
		closeAt := strings.Index(content[i+run:to], strings.Repeat("`", run))
		if closeAt < 0 {
			i += run
			continue
		}
		end := i + run + closeAt + run
		blank(out, i, end)
		i = end
	}
}

// findRegion locates the single derived region in content. A file with no
// opening marker yields ok=false and no error: an ordinary document carries
// no obligation here.
func findRegion(content string) (region, bool, error) {
	masked := maskCode(content)
	open := strings.Index(masked, RegionOpenPrefix)
	if open < 0 {
		return region{}, false, nil
	}
	if strings.Contains(masked[open+len(RegionOpenPrefix):], RegionOpenPrefix) {
		return region{}, false, fmt.Errorf("more than one derived region: a report renders one manifest")
	}
	lineEnd := strings.Index(content[open:], "\n")
	if lineEnd < 0 {
		return region{}, false, fmt.Errorf("derived region opening marker is not a complete line")
	}
	lineEnd += open
	markerLine := strings.TrimRight(content[open:lineEnd], "\r")
	if !strings.HasSuffix(markerLine, RegionOpenSuffix) {
		return region{}, false, fmt.Errorf("derived region opening marker is malformed: %q", markerLine)
	}
	manifest := strings.TrimSpace(strings.TrimSuffix(markerLine[len(RegionOpenPrefix):], RegionOpenSuffix))
	if manifest == "" {
		return region{}, false, fmt.Errorf("derived region names no manifest")
	}
	end := strings.Index(masked, RegionEnd)
	if end < 0 || end < lineEnd {
		return region{}, false, fmt.Errorf("derived region has no %s marker", RegionEnd)
	}
	body := content[lineEnd+1 : end]
	return region{manifest: manifest, bodyStart: lineEnd + 1, bodyEnd: end, body: body}, true, nil
}

// outcome names how one criterion left the source mapping, in the order the
// rendered summary lists them.
type outcome struct {
	label  string
	detail string
	rank   int
}

const (
	rankProven = iota
	rankDeferred
	rankExcluded
)

// RenderRegion renders the derived region's body for manifest m, named at
// manifestPath. The result is deterministic: the same manifest always renders
// the same bytes, so a regenerated region either matches what is committed or
// names a real disagreement.
func RenderRegion(manifestPath string, m SourceManifest) string {
	byID := map[string][]SourceEntry{}
	var ids []string
	for _, e := range m.Entries {
		if _, seen := byID[e.ID]; !seen {
			ids = append(ids, e.ID)
		}
		byID[e.ID] = append(byID[e.ID], e)
	}
	sort.Strings(ids)

	outcomes := make(map[string]outcome, len(ids))
	counts := map[int]int{}
	for _, id := range ids {
		o := describe(byID[id])
		outcomes[id] = o
		counts[o.rank]++
	}

	var b strings.Builder
	b.WriteString("\n")
	fmt.Fprintf(&b, "Derived from `%s` — source revision `%s`, read from `%s`. Regenerate with `clue report`; never write inside this region by hand.\n\n",
		cell(manifestPath), cell(m.SourceRevision), cell(m.SourceLocation))
	b.WriteString("| Outcome | Criteria |\n|---|---|\n")
	fmt.Fprintf(&b, "| Proven in the target | %d |\n", counts[rankProven])
	fmt.Fprintf(&b, "| Deferred with a plan door | %d |\n", counts[rankDeferred])
	fmt.Fprintf(&b, "| Excluded from the migration | %d |\n", counts[rankExcluded])
	b.WriteString("\n| Criterion | Outcome | Detail |\n|---|---|---|\n")
	for _, id := range ids {
		o := outcomes[id]
		fmt.Fprintf(&b, "| `%s` | %s | %s |\n", cell(id), o.label, o.detail)
	}
	b.WriteString("\n")
	return b.String()
}

// describe states one criterion's outcome from its source entries. The
// manifest is already validated, so the three cases are exhaustive.
func describe(entries []SourceEntry) outcome {
	first := entries[0]
	switch {
	case first.Excluded:
		return outcome{label: "excluded", detail: cell(first.Reason), rank: rankExcluded}
	case first.Disposition != "":
		return outcome{
			label:  "deferred (" + string(first.Disposition) + ")",
			detail: fmt.Sprintf("door `%s` · source `%s` · %s", cell(first.PlanDoor), cell(first.SourceLocation), cell(first.Justification)),
			rank:   rankDeferred,
		}
	default:
		directions, locations := sourceEvidence(entries)
		return outcome{
			label:  "proven",
			detail: fmt.Sprintf("%s · %s · %s", cell(first.ProofClass), cell(strings.Join(directions, ", ")), cell(strings.Join(locations, ", "))),
			rank:   rankProven,
		}
	}
}

// cell makes an authored string safe to place in a Markdown table cell
// without changing what it says.
func cell(s string) string {
	s = strings.ReplaceAll(s, "|", `\|`)
	s = strings.ReplaceAll(s, "\r", " ")
	s = strings.ReplaceAll(s, "\n", " ")
	return strings.TrimSpace(s)
}

// CheckReports re-renders every derived region under root and reports each
// disagreement. It is the required half of ADR-054: a typed figure, a region
// left stale by a revised manifest, and a marker naming a manifest that is
// not there all fail here.
func CheckReports(root string) []corpus.Issue {
	var issues []corpus.Issue
	for _, top := range []string{"docs", "changes"} {
		dir := filepath.Join(root, top)
		info, err := os.Stat(dir)
		if err != nil || !info.IsDir() {
			continue
		}
		_ = filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() || !strings.HasSuffix(d.Name(), ".md") {
				return nil
			}
			rel := filepath.ToSlash(strings.TrimPrefix(p, root+string(filepath.Separator)))
			data, readErr := os.ReadFile(p)
			if readErr != nil {
				return nil
			}
			for _, msg := range checkReport(root, string(data)) {
				issues = append(issues, corpus.Issue{Path: rel, Msg: msg})
			}
			return nil
		})
	}
	sort.Slice(issues, func(i, j int) bool {
		if issues[i].Path != issues[j].Path {
			return issues[i].Path < issues[j].Path
		}
		return issues[i].Msg < issues[j].Msg
	})
	return issues
}

// checkReport returns every message one report's content earns.
func checkReport(root, content string) []string {
	r, ok, err := findRegion(content)
	if err != nil {
		return []string{err.Error()}
	}
	if !ok {
		return nil
	}
	m, err := LoadSourceManifest(filepath.Join(root, filepath.FromSlash(r.manifest)))
	if err != nil {
		return []string{fmt.Sprintf("derived region names a manifest it cannot read: %v", err)}
	}
	if RenderRegion(r.manifest, m) != r.body {
		return []string{fmt.Sprintf("derived region disagrees with %s — regenerate it with clue report, never by hand", r.manifest)}
	}
	return nil
}

// WriteReport renders the derived region of the report at reportPath from the
// manifest its marker names, replacing whatever the region held. Everything
// outside the region is left byte-for-byte alone: the report's prose is the
// author's, only its figures are generated.
func WriteReport(root, reportPath string) error {
	data, err := os.ReadFile(reportPath)
	if err != nil {
		return err
	}
	content := string(data)
	r, ok, err := findRegion(content)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("%s: no derived region — add %s<manifest>%s and %s", reportPath, RegionOpenPrefix, RegionOpenSuffix, RegionEnd)
	}
	m, err := LoadSourceManifest(filepath.Join(root, filepath.FromSlash(r.manifest)))
	if err != nil {
		return err
	}
	rendered := content[:r.bodyStart] + RenderRegion(r.manifest, m) + content[r.bodyEnd:]
	if rendered == content {
		return nil
	}
	return os.WriteFile(reportPath, []byte(rendered), 0o644)
}
