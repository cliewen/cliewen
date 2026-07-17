// Package scaffold materializes the Cliewen convention into a repository:
// the /docs taxonomy, the AGENTS.md routing hub, the agent skills, and a
// CI workflow template (CAP-001). The templates are embedded so an
// installed binary is self-contained — no network, no checkout (ADR-018).
//
// The template tree lives here rather than in a repo-root dotted folder
// because the Go toolchain ignores directories starting with "." or "_",
// which puts them out of go:embed's reach.
package scaffold

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

//go:embed all:templates
var templates embed.FS

// versionRe pulls the version stamp out of an embedded skill's
// frontmatter; the skills are the canonical carriers of the pair version
// (ADR-011), so the CI template inherits its pin from them.
var versionRe = regexp.MustCompile(`(?m)^version:\s*(\S+)`)

const versionPlaceholder = "__CLUE_VERSION__"

// Report lists what a Run did, repo-relative with forward slashes.
type Report struct {
	Created        []string
	Skipped        []string // existed already — never overwritten
	Indexed        []string // README index blocks regenerated on this run
	MissingReadmes []string // pre-existing docs folders without the README validate requires
}

// PairVersion is the version stamp the embedded skills carry.
func PairVersion() (string, error) {
	data, err := templates.ReadFile("templates/skills/clue-delta/skill.md")
	if err != nil {
		return "", err
	}
	m := versionRe.FindSubmatch(data)
	if m == nil {
		return "", fmt.Errorf("embedded skills carry no version stamp")
	}
	return string(m[1]), nil
}

// Run emits the convention into root. Existing files are never touched
// except taxonomy README index blocks, which are regenerated between
// their clue:index markers (prose outside the markers is preserved).
func Run(root string) (*Report, error) {
	version, err := PairVersion()
	if err != nil {
		return nil, err
	}
	rep := &Report{}
	err = fs.WalkDir(templates, "templates", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel := strings.TrimPrefix(p, "templates/")
		data, rerr := templates.ReadFile(p)
		if rerr != nil {
			return rerr
		}
		data = []byte(strings.ReplaceAll(string(data), versionPlaceholder, version))
		for _, target := range targetsFor(rel) {
			if werr := writeIfAbsent(root, target, data, rep); werr != nil {
				return werr
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if err := regenIndexes(root, rep); err != nil {
		return nil, err
	}
	sort.Strings(rep.Created)
	sort.Strings(rep.Skipped)
	sort.Strings(rep.Indexed)
	sort.Strings(rep.MissingReadmes)
	return rep, nil
}

// targetsFor maps a template path to the repo paths it materializes as.
// Skills go to .agents/skills (canonical) and are mirrored to
// .claude/skills with the Claude Code SKILL.md spelling; the github/
// prefix stands in for .github/ (a dotted folder could not be embedded).
func targetsFor(rel string) []string {
	switch {
	case strings.HasPrefix(rel, "github/"):
		return []string{"." + rel}
	case strings.HasPrefix(rel, "skills/"):
		mirror := strings.TrimPrefix(rel, "skills/")
		if path.Base(mirror) == "skill.md" {
			mirror = path.Dir(mirror) + "/SKILL.md"
		}
		return []string{".agents/" + rel, ".claude/skills/" + mirror}
	default:
		return []string{rel}
	}
}

func writeIfAbsent(root, target string, data []byte, rep *Report) error {
	full := filepath.Join(root, filepath.FromSlash(target))
	if _, err := os.Stat(full); err == nil {
		rep.Skipped = append(rep.Skipped, target)
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(full, data, 0o644); err != nil {
		return err
	}
	rep.Created = append(rep.Created, target)
	return nil
}

// Index regeneration: the taxonomy READMEs (docs/README.md and each
// docs/<folder>/README.md) must reference every sibling artifact and
// artifact-bearing subfolder inside their clue:index block — the same
// contract clue validate checks. Regeneration keeps existing entries
// whose targets still exist (hand-written descriptions survive), drops
// entries whose targets are gone, and appends plain entries for anything
// missing. Prose outside the markers is never touched.
//
// An entry survives only as its single line, and only when its link is
// one validate's index rule recognizes (no anchor, no continuation
// lines) — anything else is replaced by a plain generated entry, the
// same reading checkIndexes applies.
const (
	indexStart = "<!-- clue:index:start -->"
	indexEnd   = "<!-- clue:index:end -->"
)

var indexLinkRe = regexp.MustCompile(`\]\(([^)#\s]+)\)`)

func regenIndexes(root string, rep *Report) error {
	docs := filepath.Join(root, "docs")
	entries, err := os.ReadDir(docs)
	if err != nil {
		return nil // no docs tree: nothing to index
	}
	readmes := []string{"docs/README.md"}
	for _, e := range entries {
		if e.IsDir() {
			readmes = append(readmes, path.Join("docs", e.Name(), "README.md"))
		}
	}
	for _, rel := range readmes {
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(rel))); err != nil {
			// A pre-existing folder without the README validate requires:
			// init does not invent one, but the report names the gap so
			// the first validate is not red without warning.
			rep.MissingReadmes = append(rep.MissingReadmes, rel)
			continue
		}
		changed, err := regenIndex(root, rel)
		if err != nil {
			return err
		}
		if changed {
			rep.Indexed = append(rep.Indexed, rel)
		}
	}
	return nil
}

func regenIndex(root, rel string) (bool, error) {
	full := filepath.Join(root, filepath.FromSlash(rel))
	raw, err := os.ReadFile(full)
	if err != nil {
		return false, nil // folder without README: validate reports it, init does not invent one
	}
	orig := strings.ReplaceAll(string(raw), "\r\n", "\n")
	text := orig
	start := strings.Index(text, indexStart)
	end := strings.Index(text, indexEnd)
	switch {
	case start < 0 && end < 0:
		// A pre-existing taxonomy README without markers would fail
		// validate ("index markers missing"), breaking the green-after-init
		// contract in existing repos. Append an empty block — prose stays
		// untouched, the generated entries land between the markers below.
		if !strings.HasSuffix(text, "\n") {
			text += "\n"
		}
		text += "\n" + indexStart + "\n" + indexEnd + "\n"
		start = strings.Index(text, indexStart)
		end = strings.Index(text, indexEnd)
	case start < 0 || end < 0 || end < start:
		// A lone or reversed marker is ambiguous: guessing at the block's
		// bounds could swallow prose between the stray marker and the end
		// of the file. Error loudly instead — the user fixes the markers,
		// init never guesses.
		return false, fmt.Errorf("%s: index markers malformed (lone or reversed %s … %s) — fix the markers by hand and re-run", rel, indexStart, indexEnd)
	}

	wanted, err := indexTargets(root, path.Dir(rel))
	if err != nil {
		return false, err
	}

	// Keep existing lines that still point at a wanted target, in order.
	var lines []string
	covered := map[string]bool{}
	for _, line := range strings.Split(text[start+len(indexStart):end], "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		m := indexLinkRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		t := path.Clean(m[1])
		if wanted[t] && !covered[t] {
			covered[t] = true
			lines = append(lines, line)
		}
	}
	var missing []string
	for t := range wanted {
		if !covered[t] {
			missing = append(missing, t)
		}
	}
	sort.Strings(missing)
	for _, t := range missing {
		label := strings.TrimSuffix(t, ".md")
		if strings.HasSuffix(t, "/README.md") {
			label = path.Dir(t) + "/"
		}
		lines = append(lines, "- ["+label+"]("+t+")")
	}

	block := indexStart + "\n"
	if len(lines) > 0 {
		block += strings.Join(lines, "\n") + "\n"
	}
	next := text[:start] + block + text[end:]
	if next == orig {
		return false, nil
	}
	return true, os.WriteFile(full, []byte(next), 0o644)
}

// indexTargets lists what a taxonomy README's index must reference:
// sibling .md artifacts and subfolders that contain markdown.
func indexTargets(root, dir string) (map[string]bool, error) {
	wanted := map[string]bool{}
	entries, err := os.ReadDir(filepath.Join(root, filepath.FromSlash(dir)))
	if err != nil {
		return wanted, err
	}
	for _, e := range entries {
		if e.IsDir() {
			if dirHasMarkdown(filepath.Join(root, filepath.FromSlash(dir), e.Name())) {
				wanted[e.Name()+"/README.md"] = true
			}
		} else if strings.HasSuffix(e.Name(), ".md") && e.Name() != "README.md" {
			wanted[e.Name()] = true
		}
	}
	return wanted, nil
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
