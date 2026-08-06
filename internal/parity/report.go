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
// everything inside a fenced block, an indented code block, or an inline code
// span replaced by spaces, and the offset of an opening fence that is never
// closed, or -1 when every fence closes. Offsets into the mask are offsets
// into the original, so a caller can search the mask and slice the original.
// A document that documents the markers — this decision record, the skill,
// the design note — writes them in code spans, in fenced blocks, or indented,
// and a scanner that could not tell an example from a claim would make
// describing the contract impossible (the same rule the outward-reference
// scan already follows).
func maskCode(content string) (string, int) {
	out := []byte(content)
	// An open fence is closed only by a bare fence of the same character and
	// at least the same length, as CommonMark defines it. Closing on any fence
	// line would let a tilde fence inside a backtick block, or an inner fence
	// inside a longer outer one, mask the whole rest of the file — and a
	// masked file has no region, so the required check would pass a report
	// whose figures were typed. This check has to fail closed.
	var openChar byte
	openLen, openAt := 0, -1
	// An indented code block is the other way a record writes an example. It
	// opens on a line indented four columns after a blank one and runs until
	// a line that is neither blank nor indented, which is what CommonMark
	// reads there too.
	indented, prevBlank := false, true
	for lineStart := 0; lineStart < len(out); {
		lineEnd := strings.IndexByte(content[lineStart:], '\n')
		if lineEnd < 0 {
			lineEnd = len(out)
		} else {
			lineEnd += lineStart
		}
		line := strings.TrimRight(content[lineStart:lineEnd], "\r")
		blankLine := strings.TrimLeft(line, " \t") == ""
		fenceChar, fenceLen := fenceOf(line)
		switch {
		case openLen > 0:
			blank(out, lineStart, lineEnd)
			if fenceChar == openChar && fenceLen >= openLen && isBareFence(line, openChar) {
				openChar, openLen, openAt = 0, 0, -1
			}
		case indented && (blankLine || indentWidth(line) >= 4):
			blank(out, lineStart, lineEnd)
		case fenceLen > 0:
			openChar, openLen, openAt = fenceChar, fenceLen, lineStart
			indented = false
			blank(out, lineStart, lineEnd)
		case prevBlank && !blankLine && indentWidth(line) >= 4:
			indented = true
			blank(out, lineStart, lineEnd)
		default:
			indented = false
			maskSpans(out, content, lineStart, lineEnd)
		}
		prevBlank = blankLine
		lineStart = lineEnd + 1
	}
	return string(out), openAt
}

// fenceOf reports the fence character and run length a line opens or closes
// with, or a zero length when the line is not a fence. A line indented four
// columns is an indented code block rather than a fence, so an example that
// shows a fence cannot open one.
func fenceOf(line string) (byte, int) {
	if indentWidth(line) >= 4 {
		return 0, 0
	}
	trimmed := strings.TrimLeft(line, " \t")
	if trimmed == "" {
		return 0, 0
	}
	c := trimmed[0]
	if c != '`' && c != '~' {
		return 0, 0
	}
	n := 0
	for n < len(trimmed) && trimmed[n] == c {
		n++
	}
	if n < 3 {
		return 0, 0
	}
	return c, n
}

// isBareFence reports whether a line is nothing but a run of c, which is what
// a closing fence is. A line carrying an info string — the "```markdown" that
// opens a nested example — opens a block and never closes one.
func isBareFence(line string, c byte) bool {
	trimmed := strings.TrimRight(strings.TrimLeft(line, " \t"), " \t")
	if trimmed == "" {
		return false
	}
	for i := 0; i < len(trimmed); i++ {
		if trimmed[i] != c {
			return false
		}
	}
	return true
}

// indentWidth counts the columns a line is indented, expanding a tab the way
// Markdown does.
func indentWidth(line string) int {
	width := 0
	for i := 0; i < len(line); i++ {
		switch line[i] {
		case ' ':
			width++
		case '\t':
			width += 4
		default:
			return width
		}
	}
	return width
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
	masked, unterminatedFenceAt := maskCode(content)
	// A fence that is never closed masks every line after it, so a report
	// carrying one would hide its own region and pass the required check in
	// silence. The exemption exists for an example, not for an accident:
	// name the unclosed fence rather than let it switch the check off.
	if unterminatedFenceAt >= 0 {
		hidden := content[unterminatedFenceAt:]
		if strings.Contains(hidden, RegionOpenPrefix) || strings.Contains(hidden, RegionEnd) {
			return region{}, false, fmt.Errorf("an unterminated code fence hides a derived region marker: close the fence, since an unclosed one would exempt this file from the check")
		}
	}
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
		_ = filepath.WalkDir(dir, func(p string, d fs.DirEntry, walkErr error) error {
			if d == nil || d.IsDir() || !strings.HasSuffix(d.Name(), ".md") {
				return nil
			}
			rel := relativePath(root, p)
			// A file the judge cannot read is reported, never skipped: a
			// required check that exempts what it failed to open is a check
			// that can be turned off by making a file unreadable.
			if walkErr != nil {
				issues = append(issues, corpus.Issue{Path: rel, Msg: fmt.Sprintf("cannot be read for a derived region: %v", walkErr)})
				return nil
			}
			data, readErr := os.ReadFile(p)
			if readErr != nil {
				issues = append(issues, corpus.Issue{Path: rel, Msg: fmt.Sprintf("cannot be read for a derived region: %v", readErr)})
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

// relativePath states a scanned file the way every other issue in the same
// verdict states one: relative to root, forward slashes, whatever shape the
// caller wrote root in.
func relativePath(root, p string) string {
	rel, err := filepath.Rel(root, p)
	if err != nil {
		return filepath.ToSlash(p)
	}
	return filepath.ToSlash(rel)
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
		return []string{fmt.Sprintf("derived region names a manifest it cannot read: %v (a documented example belongs in a code span or a fenced block, which is not read as a region)", err)}
	}
	// Line endings are the checkout's, not the report's: a Windows working
	// tree carrying CRLF must not read as a disagreement about figures.
	if normalizeNewlines(RenderRegion(r.manifest, m)) != normalizeNewlines(r.body) {
		return []string{fmt.Sprintf("derived region disagrees with %s — regenerate it with clue report, never by hand", r.manifest)}
	}
	return nil
}

func normalizeNewlines(s string) string { return strings.ReplaceAll(s, "\r\n", "\n") }

// usesCRLF reports whether a file's first line ending is a CRLF, which is
// what a Windows checkout of a repository without an end-of-line policy has.
func usesCRLF(content string) bool {
	i := strings.IndexByte(content, '\n')
	return i > 0 && content[i-1] == '\r'
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
	region := RenderRegion(r.manifest, m)
	// A report the check already accepts is left untouched, and a checkout
	// using CRLF is rendered in its own line endings. Writing LF into a CRLF
	// file would rewrite a green report, leave it mixed, and dirty a worktree
	// at exactly the moment publication requires a clean one.
	if normalizeNewlines(region) == normalizeNewlines(r.body) {
		return nil
	}
	if usesCRLF(content) {
		region = strings.ReplaceAll(normalizeNewlines(region), "\n", "\r\n")
	}
	return os.WriteFile(reportPath, []byte(content[:r.bodyStart]+region+content[r.bodyEnd:]), 0o644)
}
