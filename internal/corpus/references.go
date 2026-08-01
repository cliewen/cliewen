package corpus

import (
	"regexp"
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
// code spans, and link targets are content rather than citation, and heading
// anchors and colour literals merely share the character. Each exclusion is
// mechanical, so the same file yields the same verdict offline, on a pinned
// revision, in a year.
func checkExternalReferences(c *Corpus) []Issue {
	var issues []Issue
	for _, a := range c.Artifacts {
		for _, ref := range scanBareForgeRefs(a.Body) {
			// The line is part of the finding: a corpus file is long, the same
			// number can appear more than once in it, and an adopter repairing
			// a migration report should not have to search for the occurrence.
			issues = append(issues, Issue{a.Path, "line " + itoa(ref.line) + ": bare forge reference " + ref.text +
				" names no repository; write the full URL of what it points at (ADR-040)"})
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
	// Three or six digits are the shapes a hex colour takes. Requiring such a
	// literal to sit in a code span is ordinary markdown practice and is what
	// keeps this from guessing: outside one, the digits are read as a
	// reference and the author qualifies or fences them.
	return true
}

func isWordByte(b byte) bool {
	return b == '_' ||
		(b >= 'a' && b <= 'z') ||
		(b >= 'A' && b <= 'Z') ||
		(b >= '0' && b <= '9')
}

func itoa(n int) string { return strconv.Itoa(n) }
