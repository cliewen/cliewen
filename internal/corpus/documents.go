package corpus

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// IndexType is the document type a folder README declares. It is what makes
// the README a typed document rather than an artifact (PDR-065).
const IndexType = "index"

// MissingIndexHeader is the scan issue for a folder README with no
// frontmatter. clue migrate repairs it, so migration planning reads past it.
const MissingIndexHeader = "missing frontmatter (a folder README carries type: index and a title, PDR-065)"

// DeliveredMarkdown lists the Markdown files Cliewen materializes outside
// docs/ that an adopter keeps and edits. Each carries a document header when
// present (PDR-065). The pull-request template is deliberately absent: the
// forge pastes it verbatim into every new pull request. Generated skill files
// are checked separately, by ownership marker.
var DeliveredMarkdown = []string{
	"AGENTS.md",
	"CLAUDE.md",
	".clue/evidence/README.md",
	".clue/evidence/frameworks.md",
	".github/cliewen-wall.md",
}

// DocumentHeaderGap names what a non-artifact Markdown file's header lacks,
// or returns "" when it carries a non-empty type and title.
func DocumentHeaderGap(text string) string {
	fields, _, ok, err := parseFrontmatter(strings.ReplaceAll(text, "\r\n", "\n"))
	switch {
	case err != nil:
		return "frontmatter does not parse as YAML: " + err.Error()
	case !ok:
		return "missing frontmatter (expected type and title, PDR-065)"
	}
	var missing []string
	for _, key := range []string{"type", "title"} {
		if s, _ := fields[key].(string); strings.TrimSpace(s) == "" {
			missing = append(missing, key)
		}
	}
	if len(missing) > 0 {
		return "frontmatter has no " + strings.Join(missing, " or ") + " (PDR-065)"
	}
	return ""
}

// checkIndexHeader holds a folder README to its document header: type: index
// and a non-empty title.
func checkIndexHeader(rel string, fields map[string]any) []Issue {
	var issues []Issue
	if t, _ := fields["type"].(string); t != IndexType {
		issues = append(issues, Issue{rel, "folder README must declare type: index (PDR-065)"})
	}
	if t, _ := fields["title"].(string); strings.TrimSpace(t) == "" {
		issues = append(issues, Issue{rel, "folder README has no title (PDR-065)"})
	}
	return issues
}

// checkDeliveredHeaders requires a document header on every delivered
// Markdown file that is present and on every file of a Cliewen-managed skill.
// An absent file is the adopter's choice, and Markdown Cliewen did not
// deliver is never inspected (AC-214).
func checkDeliveredHeaders(c *Corpus) []Issue {
	var issues []Issue
	check := func(rel string) {
		data, err := os.ReadFile(filepath.Join(c.Root, filepath.FromSlash(rel)))
		if err != nil {
			return
		}
		if gap := DocumentHeaderGap(string(data)); gap != "" {
			issues = append(issues, Issue{rel, gap})
		}
	}
	for _, rel := range DeliveredMarkdown {
		check(rel)
	}
	for _, rel := range managedSkillMarkdown(c.Root) {
		check(rel)
	}
	return issues
}

// managedSkillMarkdown lists every Markdown file inside a skill directory
// whose manifest declares cliewen-skill: true, under .agents/skills and the
// .claude/skills mirror. A mirror that is itself a symlink resolves to files
// already listed, so it is skipped rather than reported twice.
func managedSkillMarkdown(root string) []string {
	var files []string
	for _, dir := range []string{skillsDir, ".claude/skills"} {
		abs := filepath.Join(root, filepath.FromSlash(dir))
		info, err := os.Lstat(abs)
		if err != nil || !info.IsDir() {
			continue
		}
		entries, err := os.ReadDir(abs)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			skillDir := filepath.Join(abs, e.Name())
			manifest, ambiguous, err := findSkillManifest(skillDir)
			if err != nil || ambiguous || manifest == "" {
				continue
			}
			data, err := os.ReadFile(filepath.Join(skillDir, manifest))
			if err != nil {
				continue
			}
			fields, _, ok, perr := parseFrontmatter(strings.ReplaceAll(string(data), "\r\n", "\n"))
			if perr != nil || !ok || fields[cliewenSkillMarker] != true {
				continue
			}
			_ = filepath.WalkDir(skillDir, func(p string, d fs.DirEntry, err error) error {
				if err != nil || d.IsDir() || !strings.HasSuffix(d.Name(), ".md") {
					return nil
				}
				rel, _ := filepath.Rel(root, p)
				files = append(files, filepath.ToSlash(rel))
				return nil
			})
		}
	}
	sort.Strings(files)
	return files
}
