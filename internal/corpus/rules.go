package corpus

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Options control which gates Validate applies.
type Options struct {
	// ForbidChanges fails when /changes contains any file: the
	// digest-before-merge gate. CI uses it; local runs during a change
	// loop do not.
	ForbidChanges bool
}

// statusVocab is the allowed status set per artifact type. The docs/README.md
// status table mirrors this map; the map is the consumer that keeps the
// `status` field honest (Foundation §7: every field must have a consumer).
var statusVocab = map[string][]string{
	"goal":           {"proposed", "accepted", "retired"},
	"plan":           {"draft", "active", "completed"},
	"capability":     {"draft", "active", "retired"},
	"criteria":       {"draft", "active", "retired"},
	"design":         {"draft", "active", "retired"},
	"decision":       {"inferred", "verified"},
	"constraint":     {"active", "retired"},
	"quality":        {"draft", "active", "retired"},
	"analysis":       {"verified"},
	"architecture":   {"draft", "verified"},
	"change":         {"open"},
	"tasks":          {"open"},
	"open-questions": {"open", "resolved"},
}

var (
	indexStart  = "<!-- clue:index:start -->"
	indexEnd    = "<!-- clue:index:end -->"
	milestoneRe = regexp.MustCompile(`\bM-\d+\b`)
	mdLinkRe    = regexp.MustCompile(`\]\(([^)#\s]+)\)`)
	externalRe  = regexp.MustCompile(`^[a-z][a-z0-9+.-]*:`) // http:, https:, mailto:, ...
)

// Validate applies the graph and layout rules to a scanned corpus.
func Validate(c *Corpus, opts Options) []Issue {
	var issues []Issue
	issues = append(issues, checkCoreFields(c)...)
	issues = append(issues, checkDuplicateIDs(c)...)
	issues = append(issues, checkStatusVocab(c)...)
	issues = append(issues, checkLinks(c)...)
	issues = append(issues, checkFolderReadmes(c)...)
	issues = append(issues, checkIndexes(c)...)
	if opts.ForbidChanges && c.HasChanges {
		issues = append(issues, Issue{"changes", "transient workspace present — digest before merge (main must never contain /changes)"})
	}
	sort.Slice(issues, func(i, j int) bool {
		if issues[i].Path != issues[j].Path {
			return issues[i].Path < issues[j].Path
		}
		return issues[i].Msg < issues[j].Msg
	})
	return issues
}

func checkCoreFields(c *Corpus) []Issue {
	var issues []Issue
	for _, a := range c.Artifacts {
		var missing []string
		for _, f := range CoreFields {
			v, present := a.Fields[f]
			if !present {
				missing = append(missing, f)
				continue
			}
			// links may legitimately be an empty list; the rest must be
			// non-empty strings.
			if f != "links" {
				if s, ok := v.(string); !ok || s == "" {
					missing = append(missing, f)
				}
			}
		}
		if len(missing) > 0 {
			issues = append(issues, Issue{a.Path, "missing or empty core field(s): " + strings.Join(missing, ", ")})
		}
	}
	return issues
}

func checkDuplicateIDs(c *Corpus) []Issue {
	var issues []Issue
	for id, as := range c.ByID {
		if len(as) > 1 {
			var paths []string
			for _, a := range as {
				paths = append(paths, a.Path)
			}
			sort.Strings(paths)
			issues = append(issues, Issue{paths[0], "duplicate id " + id + " (also in " + strings.Join(paths[1:], ", ") + ")"})
		}
	}
	return issues
}

func checkStatusVocab(c *Corpus) []Issue {
	var issues []Issue
	for _, a := range c.Artifacts {
		if a.Type == "" {
			continue // reported by checkCoreFields
		}
		allowed, known := statusVocab[a.Type]
		if !known {
			issues = append(issues, Issue{a.Path, "unknown type " + a.Type})
			continue
		}
		if a.Status == "" {
			continue // reported by checkCoreFields
		}
		ok := false
		for _, s := range allowed {
			if a.Status == s {
				ok = true
				break
			}
		}
		if !ok {
			issues = append(issues, Issue{a.Path, "status " + a.Status + " not allowed for type " + a.Type + " (allowed: " + strings.Join(allowed, ", ") + ")"})
		}
	}
	return issues
}

func checkLinks(c *Corpus) []Issue {
	// Milestones (M-xxx) live inside their plan file, not as separate
	// artifacts: harvest them from plan bodies.
	milestones := map[string]bool{}
	for _, a := range c.Artifacts {
		if a.Type == "plan" {
			for _, m := range milestoneRe.FindAllString(a.Body, -1) {
				milestones[m] = true
			}
		}
	}
	var issues []Issue
	for _, a := range c.Artifacts {
		for _, l := range a.Links {
			if milestoneRe.MatchString(l) {
				if !milestones[l] {
					issues = append(issues, Issue{a.Path, "link " + l + " not found in any plan"})
				}
				continue
			}
			if _, ok := c.ByID[l]; !ok {
				issues = append(issues, Issue{a.Path, "link " + l + " resolves to no artifact"})
			}
		}
	}
	return issues
}

func checkFolderReadmes(c *Corpus) []Issue {
	var issues []Issue
	docs := filepath.Join(c.Root, "docs")
	if info, err := os.Stat(docs); err != nil || !info.IsDir() {
		return nil
	}
	_ = filepath.WalkDir(docs, func(p string, d fs.DirEntry, err error) error {
		if err != nil || !d.IsDir() {
			return nil
		}
		if _, err := os.Stat(filepath.Join(p, "README.md")); err != nil {
			rel, _ := filepath.Rel(c.Root, p)
			issues = append(issues, Issue{filepath.ToSlash(rel), "folder has no README.md"})
		}
		return nil
	})
	return issues
}

// checkIndexes enforces the generated-index contract on the taxonomy
// READMEs (docs/README.md and each docs/<folder>/README.md): index
// markers must exist, every link inside the block must resolve to a
// file, and every sibling artifact or artifact-bearing subfolder must
// be referenced in the block.
func checkIndexes(c *Corpus) []Issue {
	var issues []Issue
	for _, rel := range c.MDFiles {
		parts := strings.Split(rel, "/")
		isTaxonomyReadme := rel == "docs/README.md" ||
			(len(parts) == 3 && parts[0] == "docs" && parts[2] == "README.md")
		if !isTaxonomyReadme {
			continue
		}
		text := c.Contents[rel]
		start := strings.Index(text, indexStart)
		end := strings.Index(text, indexEnd)
		if start < 0 || end < 0 || end < start {
			issues = append(issues, Issue{rel, "index markers missing or malformed (" + indexStart + " … " + indexEnd + ")"})
			continue
		}
		block := text[start:end]
		dir := path.Dir(rel)

		targets := map[string]bool{}
		for _, m := range mdLinkRe.FindAllStringSubmatch(block, -1) {
			t := m[1]
			if externalRe.MatchString(t) {
				continue
			}
			t = path.Clean(t)
			targets[t] = true
			if _, err := os.Stat(filepath.Join(c.Root, filepath.FromSlash(path.Join(dir, t)))); err != nil {
				issues = append(issues, Issue{rel, "index references missing file " + t})
			}
		}

		entries, err := os.ReadDir(filepath.Join(c.Root, filepath.FromSlash(dir)))
		if err != nil {
			continue
		}
		for _, e := range entries {
			name := e.Name()
			if e.IsDir() {
				if !dirHasMarkdown(filepath.Join(c.Root, filepath.FromSlash(dir), name)) {
					continue
				}
				covered := false
				for t := range targets {
					if t == name || strings.HasPrefix(t, name+"/") {
						covered = true
						break
					}
				}
				if !covered {
					issues = append(issues, Issue{rel, "index does not reference subfolder " + name + "/"})
				}
			} else if strings.HasSuffix(name, ".md") && name != "README.md" && !targets[name] {
				issues = append(issues, Issue{rel, "index does not reference sibling " + name})
			}
		}
	}
	return issues
}

func dirHasMarkdown(dir string) bool {
	found := false
	_ = filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() && strings.HasSuffix(d.Name(), ".md") {
			found = true
			return fs.SkipAll
		}
		return nil
	})
	return found
}
