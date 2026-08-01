package corpus

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// External references are qualified by ADR-040: anything with an address is a
// full URL, an identity in another repository's corpus carries the clue:
// scheme, and a bare forge number is a failure because it asserts a namespace
// it never names. The judge enforces that form and nothing else — resolving a
// target needs the network, so it lives in a separate command and can never
// reach a validation verdict.

// bareForgeRe matches a forge number that names no repository. The lookbehind
// Go's regexp lacks is done by hand in scanBareForgeRefs: what survives here is
// the shape, and the caller decides whether the position is one where a
// reference can occur at all.
var bareForgeRe = regexp.MustCompile(`#([0-9]+)`)

// fenceRe opens or closes a fenced block. Info strings and indentation are
// irrelevant to the toggle, so only the marker matters.
var fenceRe = regexp.MustCompile("^\\s*(```|~~~)")

// clueRefRe matches the identity form: clue:owner/repo[@revision]/ID. The
// revision is optional because only acceptance-evidence pointers pin one; a
// plain cross-repository citation names the identity alone.
var clueRefRe = regexp.MustCompile(`\bclue:([A-Za-z0-9_.-]+)/([A-Za-z0-9_.-]+)(@[A-Za-z0-9]+)?/([A-Z][A-Za-z0-9-]*)\b`)

// checkExternalReferences reports every bare forge number in corpus prose.
//
// A reference is only sought where one can be written: fenced blocks, inline
// code spans, and link targets are content rather than citation, and a heading
// anchor merely shares the character. A colour literal carrying a hex letter is
// excluded by the same rule as an anchor; an all-digit one is read as a
// reference on purpose, because excluding runs of that length would silence the
// three-digit issue numbers adopters actually cite. Each exclusion is
// mechanical, so the same file yields the same verdict offline, on a pinned
// revision, in a year.
func checkExternalReferences(c *Corpus) []Issue {
	var issues []Issue
	for _, a := range c.Artifacts {
		for _, ref := range scanBareForgeRefs(a.Body) {
			// The line is part of the finding: a corpus file is long, the same
			// number can appear more than once in it, and an adopter repairing
			// a migration report should not have to search for the occurrence.
			issues = append(issues, Issue{a.Path, "line " + itoa(a.BodyLine+ref.line-1) + ": bare forge reference " + ref.text +
				" names no repository; write the full URL of what it points at"})
		}
	}
	return issues
}

type bareRef struct {
	text string
	line int
}

// scanBareForgeRefs returns the bare forge numbers in markdown body text.
func scanBareForgeRefs(body string) []bareRef {
	var found []bareRef
	inFence := false
	for i, line := range strings.Split(body, "\n") {
		if fenceRe.MatchString(line) {
			inFence = !inFence
			continue
		}
		if inFence {
			continue
		}
		stripped := stripNonReferenceSpans(line)
		for _, m := range bareForgeRe.FindAllStringSubmatchIndex(stripped, -1) {
			start, end := m[0], m[1]
			if !isReferencePosition(stripped, start, end) {
				continue
			}
			found = append(found, bareRef{text: stripped[start:end], line: i + 1})
		}
	}
	return found
}

// inlineCodeRe matches a backtick span. Colour literals and anchors belong in
// code spans by ordinary markdown practice, which is what makes this exclusion
// carry most of the ambiguity without a special case for either.
var inlineCodeRe = regexp.MustCompile("`[^`]*`")

// linkTargetRe matches the target half of a markdown link. A qualified URL
// ending in a fragment is the correct form, so its own '#' must never be read
// as a bare reference.
var linkTargetRe = regexp.MustCompile(`\]\([^)]*\)`)

// addressedLinkRe matches a complete markdown link whose target is an absolute
// URL. ADR-040 keeps link text free precisely so the familiar forge shorthand
// survives as a readable label — `[cliewen/cliewen#96](https://…)` names its
// target in the half that resolves, so the label is not a bare reference and
// flagging it would forbid the form the decision endorses.
var addressedLinkRe = regexp.MustCompile(`\[[^\]]*\]\([a-zA-Z][a-zA-Z0-9+.-]*:[^)]*\)`)

// autolinkRe matches an angle-bracket autolink, whose fragment is a target for
// the same reason.
var autolinkRe = regexp.MustCompile(`<[a-zA-Z][a-zA-Z0-9+.-]*:[^>\s]*>`)

// bareURLRe matches a URL written without markdown syntax. Its fragment is
// part of the address, not a citation.
var bareURLRe = regexp.MustCompile(`[a-zA-Z][a-zA-Z0-9+.-]*://[^\s)\]]+`)

// stripNonReferenceSpans blanks every span in which a '#' cannot be a
// reference, preserving length so positions stay meaningful.
func stripNonReferenceSpans(line string) string {
	// addressedLinkRe runs first: it consumes the whole link, so a label that
	// happens to carry a forge number is never left behind for the bare-
	// reference scan.
	for _, re := range []*regexp.Regexp{inlineCodeRe, addressedLinkRe, linkTargetRe, autolinkRe, bareURLRe, clueRefRe} {
		line = re.ReplaceAllStringFunc(line, blank)
	}
	return line
}

func blank(s string) string { return strings.Repeat(" ", len(s)) }

// isReferencePosition reports whether a '#N' match is a citation rather than a
// heading, an anchor, or a colour.
func isReferencePosition(line string, start, end int) bool {
	// A heading marker owns the start of its line.
	if strings.TrimSpace(line[:start]) == "" {
		return false
	}
	// A word character or path separator before '#' means the token belongs to
	// something larger — a URL path, an entity, an identifier.
	if start > 0 {
		switch p := line[start-1]; {
		case p == '/' || p == '&' || p == '-' || p == '_':
			return false
		case isWordByte(p):
			return false
		}
	}
	// A trailing word character means an anchor such as #4-resume-game, or a
	// colour literal whose digits run into letters.
	if end < len(line) {
		if n := line[end]; isWordByte(n) || n == '-' {
			return false
		}
	}
	// A colour literal carrying any hex letter is already excluded above, because
	// the digits run into a word character. What remains ambiguous is an
	// all-digit colour such as #777777, which is indistinguishable from an
	// issue number by shape alone.
	//
	// It is read as a reference, deliberately. Excluding three- and six-digit
	// runs would silence exactly the numbers adopters actually cite — the live
	// corpora here carry #202, #218, #224 — and a rule that cannot see the
	// common case is not worth the rare false positive it avoids. An all-digit
	// colour in prose belongs in a code span or a fence, which is ordinary
	// markdown practice and is where every colour in these corpora already
	// sits.
	return true
}

func isWordByte(b byte) bool {
	return b == '_' ||
		(b >= 'a' && b <= 'z') ||
		(b >= 'A' && b <= 'Z') ||
		(b >= '0' && b <= '9')
}

func itoa(n int) string { return strconv.Itoa(n) }

// BareReferenceIssues exposes the form findings to callers outside the judge —
// notably the migration that reports them to an adopter. It is the same scan
// `clue validate` runs, so a repository never sees two different answers about
// the same reference.
func BareReferenceIssues(c *Corpus) []Issue { return checkExternalReferences(c) }

// malformedClueRe matches a clue: token that opens the identity form but does
// not complete it. The judge must say so rather than ignore it: a pointer that
// silently does nothing is worse than none, because the criterion it sits on
// looks like it named its foreign proof.
var malformedClueRe = regexp.MustCompile(`\bclue:\S*`)

// checkForeignPointers reports clue: pointers that do not parse.
//
// A well-formed pointer is deliberately not evidence. It names a repository, a
// revision, and an identifier so a reader can go and look; the judge cannot
// see that run, so the pointer never becomes coverage and never becomes an
// imported verdict. Test-type: Human remains the proof of record.
func checkForeignPointers(c *Corpus) []Issue {
	var issues []Issue
	for _, a := range c.Artifacts {
		for _, line := range citableLines(a.Body) {
			for _, tok := range malformedClueRe.FindAllString(line.text, -1) {
				tok = strings.TrimRight(tok, ".,;:)]")
				if clueRefRe.MatchString(tok) {
					continue
				}
				issues = append(issues, Issue{a.Path, "line " + itoa(a.BodyLine+line.n-1) +
					": foreign pointer " + tok + " is malformed; name the repository, the revision, and the identifier as clue:owner/repo@revision/ID"})
			}
		}
	}
	return issues
}

// ForeignPointers returns every well-formed foreign-evidence pointer declared
// in the corpus, so a coverage report can name what is proven elsewhere
// without ever counting it as local proof.
//
// It reads the same lines the malformed check reads. A pointer inside a fence
// or a code span is documentation showing the form, not a claim that something
// was proven, and reporting it would put an example into a coverage summary.
func ForeignPointers(c *Corpus) []string {
	var out []string
	seen := map[string]bool{}
	for _, a := range c.Artifacts {
		for _, line := range citableLines(a.Body) {
			for _, m := range clueRefRe.FindAllString(line.text, -1) {
				if !seen[m] {
					seen[m] = true
					out = append(out, m)
				}
			}
		}
	}
	sort.Strings(out)
	return out
}

type citableLine struct {
	text string
	n    int
}

// citableLines yields the body lines on which a citation can occur, with fenced
// blocks dropped and code spans blanked. Both foreign-pointer paths read it, so
// they can never disagree about whether an example counts as a claim.
func citableLines(body string) []citableLine {
	var out []citableLine
	inFence := false
	for i, line := range strings.Split(body, "\n") {
		if fenceRe.MatchString(line) {
			inFence = !inFence
			continue
		}
		if inFence {
			continue
		}
		out = append(out, citableLine{text: inlineCodeRe.ReplaceAllStringFunc(line, blank), n: i + 1})
	}
	return out
}
