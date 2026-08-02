package corpus

import (
	"os"
	"path"
	"regexp"
	"sort"
	"strings"
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
