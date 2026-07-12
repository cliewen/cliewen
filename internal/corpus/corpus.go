// Package corpus scans a Cliewen repository for markdown artifacts and
// validates the frontmatter graph. Frontmatter is the machine-readable
// layer — id, type, status, links, title — the markdown body is for
// humans (Foundation §7).
package corpus

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// CoreFields is the common frontmatter core every artifact must carry.
var CoreFields = []string{"id", "type", "status", "links", "title"}

// Artifact is one markdown file carrying frontmatter identity. The ID is
// the identity; the path is only the current address (Foundation §4).
type Artifact struct {
	ID     string
	Type   string
	Status string
	Links  []string
	Title  string
	Path   string // repo-relative, forward slashes
	Body   string
	Fields map[string]any // raw frontmatter, for presence checks
}

// Issue is one validation finding.
type Issue struct {
	Path string
	Msg  string
}

func (i Issue) String() string { return i.Path + ": " + i.Msg }

// Corpus is a scanned repository: the artifact graph plus the file
// inventory the layout rules need.
type Corpus struct {
	Root       string
	Artifacts  []*Artifact
	ByID       map[string][]*Artifact // >1 entry means duplicate identity
	MDFiles    []string               // all .md under docs/ and changes/
	Contents   map[string]string      // path -> file text (LF-normalized)
	HasChanges bool                   // any file under changes/
}

// Scan walks docs/ and changes/ under root, parsing frontmatter into the
// artifact graph. Parse-level problems are returned as issues; graph
// rules are applied separately by Validate.
func Scan(root string) (*Corpus, []Issue) {
	c := &Corpus{
		Root:     root,
		ByID:     map[string][]*Artifact{},
		Contents: map[string]string{},
	}
	var issues []Issue

	for _, top := range []string{"docs", "changes"} {
		dir := filepath.Join(root, top)
		if info, err := os.Stat(dir); err != nil || !info.IsDir() {
			continue
		}
		err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			if top == "changes" {
				c.HasChanges = true
			}
			if !strings.HasSuffix(d.Name(), ".md") {
				return nil
			}
			rel, _ := filepath.Rel(root, path)
			rel = filepath.ToSlash(rel)
			c.MDFiles = append(c.MDFiles, rel)

			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			text := strings.ReplaceAll(string(data), "\r\n", "\n")
			c.Contents[rel] = text

			fields, body, ok, perr := parseFrontmatter(text)
			if perr != nil {
				issues = append(issues, Issue{rel, "frontmatter does not parse as YAML: " + perr.Error()})
				return nil
			}
			if !ok {
				// README.md files are exempt from carrying frontmatter;
				// everything else in docs/ and changes/ is an artifact.
				if d.Name() != "README.md" {
					issues = append(issues, Issue{rel, "missing frontmatter (expected id, type, status, links, title)"})
				}
				return nil
			}
			a := &Artifact{Path: rel, Body: body, Fields: fields}
			a.ID, _ = fields["id"].(string)
			a.Type, _ = fields["type"].(string)
			a.Status, _ = fields["status"].(string)
			a.Title, _ = fields["title"].(string)
			a.Links, _ = stringList(fields["links"])
			c.Artifacts = append(c.Artifacts, a)
			if a.ID != "" {
				c.ByID[a.ID] = append(c.ByID[a.ID], a)
			}
			return nil
		})
		if err != nil {
			issues = append(issues, Issue{top, "walk failed: " + err.Error()})
		}
	}
	sort.Slice(c.Artifacts, func(i, j int) bool { return c.Artifacts[i].Path < c.Artifacts[j].Path })
	return c, issues
}

// parseFrontmatter splits a leading `---` YAML block from the body.
// ok is false when the file has no frontmatter at all; err is non-nil
// when a block exists but is not valid YAML.
func parseFrontmatter(text string) (fields map[string]any, body string, ok bool, err error) {
	if !strings.HasPrefix(text, "---\n") {
		return nil, text, false, nil
	}
	rest := text[len("---\n"):]
	end := strings.Index(rest, "\n---\n")
	if end < 0 {
		if strings.HasSuffix(rest, "\n---") {
			end = len(rest) - len("\n---")
		} else {
			return nil, text, false, nil
		}
	}
	if uerr := yaml.Unmarshal([]byte(rest[:end]), &fields); uerr != nil {
		return nil, text, false, uerr
	}
	body = rest[min(end+len("\n---\n"), len(rest)):]
	return fields, body, true, nil
}

// stringList coerces a frontmatter value into a list of strings.
func stringList(v any) ([]string, bool) {
	switch t := v.(type) {
	case nil:
		return nil, false
	case string:
		if t == "" {
			return nil, true
		}
		return []string{t}, true
	case []any:
		out := make([]string, 0, len(t))
		for _, e := range t {
			out = append(out, fmt.Sprint(e))
		}
		return out, true
	}
	return nil, false
}
