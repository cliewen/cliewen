// Package scaffold materializes the Cliewen convention into a repository:
// the /docs taxonomy, the AGENTS.md routing hub, the agent skills, and a
// CI workflow template (CAP-001). The templates are embedded so an
// installed binary is self-contained — no network, no checkout (ADR-018).
//
// One value is not embeddable: the emitted CI caller must name the upstream
// validation workflow at an immutable reference (ADR-038). It comes from Go's
// VCS build metadata, and for builds Go leaves unstamped from this package's
// own source checkout, which a source build still has on disk. That is a read
// of the emitting tree, never of the repository being scaffolded, and it falls
// back to the release tag rather than guess.
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
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"runtime/debug"
	"slices"
	"sort"
	"strings"
)

//go:embed all:templates
var templates embed.FS

// versionRe pulls the version stamp out of an embedded skill's
// frontmatter; the skills are the canonical carriers of the pair version
// (ADR-011), so the CI template inherits its pin from them.
var versionRe = regexp.MustCompile(`(?m)^version:\s*(\S+)`)

var workflowRefRe = regexp.MustCompile(`^[0-9a-f]{40}$`)

const versionPlaceholder = "__CLUE_VERSION__"

const workflowRefPlaceholder = "__CLUE_WORKFLOW_REF__"

// Report lists what a Run did, repo-relative with forward slashes.
type Report struct {
	Created        []string
	Skipped        []string // existed already — never overwritten
	Linked         []string // symlinked directories: nothing was written through them
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

// workflowReference returns the immutable reference an emitted caller uses
// for the upstream validation unit: the source commit that built this binary
// when that commit can be named honestly, and the protected release tag
// otherwise (ADR-038). Both forms are immutable; a wrong commit is worse than
// the tag, because an adopter's CI fails to resolve it with nothing in their
// repository to explain why. Every path below therefore prefers the tag over
// a revision it cannot vouch for.
func workflowReference(version string) string {
	tag := "v" + version
	if info, ok := debug.ReadBuildInfo(); ok {
		var revision string
		var modified bool
		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.revision":
				revision = setting.Value
			case "vcs.modified":
				modified = setting.Value == "true"
			}
		}
		// A dirty tree's revision names a commit whose validation workflow is
		// not the one this binary was built from, so it would hand an adopter
		// a reference to content they were never given.
		if modified {
			return tag
		}
		if workflowRefRe.MatchString(revision) {
			return revision
		}
	}
	if revision, ok := sourceCheckoutRevision(); ok {
		return revision
	}
	return tag
}

// sourceCheckoutRevision recovers the emitting commit for a build Go left
// without VCS metadata: `go test` does not stamp its test binaries, and a
// module built from the cache has no repository to stamp from. It reads only
// the checkout this package was compiled from — never the repository being
// scaffolded — and reports nothing unless that checkout is positively this
// project at a commit whose tracked content still matches it.
func sourceCheckoutRevision() (string, bool) {
	_, sourceFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", false
	}
	root := filepath.Dir(filepath.Dir(filepath.Dir(sourceFile)))
	if !isCliewenCheckout(root) {
		return "", false
	}
	output, err := exec.Command("git", "-C", root, "rev-parse", "HEAD").Output()
	if err != nil {
		return "", false
	}
	revision := strings.TrimSpace(string(output))
	if !workflowRefRe.MatchString(revision) {
		return "", false
	}
	// Uncommitted tracked changes mean HEAD does not describe the validation
	// workflow this binary would have an adopter call.
	if err := exec.Command("git", "-C", root, "diff", "--quiet", "HEAD").Run(); err != nil {
		return "", false
	}
	return revision, true
}

// isCliewenCheckout reports whether root is this project's own source tree.
// A module cache entry or an extracted source archive can sit inside an
// unrelated git repository, and that repository's HEAD is not a reference any
// adopter can resolve — so identity is established from the tree's own
// contents before git is consulted at all.
func isCliewenCheckout(root string) bool {
	if _, err := os.Stat(filepath.Join(root, ".github", "workflows", "clue-validation.yml")); err != nil {
		return false
	}
	mod, err := os.ReadFile(filepath.Join(root, "go.mod"))
	if err != nil {
		return false
	}
	return strings.Contains(string(mod), "module github.com/cliewen/cliewen")
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
	workflowRef := workflowReference(version)
	links := map[string]string{}
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
		data = []byte(strings.ReplaceAll(string(data), workflowRefPlaceholder, workflowRef))
		for _, target := range targetsFor(rel) {
			if werr := writeIfAbsent(root, target, data, rep, links); werr != nil {
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
	sort.Strings(rep.Linked)
	sort.Strings(rep.Indexed)
	sort.Strings(rep.MissingReadmes)
	return rep, nil
}

// Regen regenerates the taxonomy README index blocks under root without
// materializing anything (CAP-005) — the standalone exposure of the
// engine Run uses, per ADR-019. A root without a docs tree is an error:
// the command's entire purpose cannot apply, and clue init is the tool
// that materializes one.
func Regen(root string) (*Report, error) {
	if info, err := os.Stat(filepath.Join(root, "docs")); err != nil || !info.IsDir() {
		return nil, fmt.Errorf("no docs tree under %q — nothing to regenerate (`clue init` materializes one)", root)
	}
	rep := &Report{}
	if err := regenIndexes(root, rep); err != nil {
		return nil, err
	}
	sort.Strings(rep.Indexed)
	sort.Strings(rep.MissingReadmes)
	return rep, nil
}

// targetsFor maps a template path to the repo paths it materializes as.
// Skills go to .agents/skills (canonical) and are mirrored to
// .claude/skills with the Claude Code SKILL.md spelling; the github/
// prefix stands in for .github/ (a dotted folder could not be embedded).
//
// A successful run may emit no .claude/skills mirror at all: sharing an
// assistant's skills across checkouts makes that folder a symlink, and
// init never writes through one (see linkedAncestor).
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

// linkedAncestor returns the repo-relative directory under root that is a
// symlink and therefore blocks target, or "" when the path is clear. root
// itself is never inspected: a checkout reached through a link is normal.
// Lookups are memoized per directory in links, so a mirror of N files
// costs one ancestor walk rather than N.
func linkedAncestor(root, target string, links map[string]string) (string, error) {
	dir := path.Dir(target)
	if dir == "." || dir == "/" {
		return "", nil
	}
	if hit, ok := links[dir]; ok {
		return hit, nil
	}
	var prefix string
	for _, seg := range strings.Split(dir, "/") {
		prefix = path.Join(prefix, seg)
		info, err := os.Lstat(filepath.Join(root, filepath.FromSlash(prefix)))
		if err != nil {
			if os.IsNotExist(err) {
				break // nothing exists from here down, so no link can block us
			}
			return "", err
		}
		if info.Mode()&fs.ModeSymlink != 0 {
			links[dir] = prefix
			return prefix, nil
		}
	}
	links[dir] = ""
	return "", nil
}

func writeIfAbsent(root, target string, data []byte, rep *Report, links map[string]string) error {
	// A symlinked directory below root is a tree the user shares across
	// checkouts; init leaves it alone and reports the skip rather than
	// mutating a directory outside this repository.
	blocked, lerr := linkedAncestor(root, target, links)
	if lerr != nil {
		return lerr
	}
	if blocked != "" {
		if !slices.Contains(rep.Linked, blocked) {
			rep.Linked = append(rep.Linked, blocked)
		}
		return nil
	}
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
		if os.IsNotExist(err) {
			return nil // no docs tree: nothing to index
		}
		return err // a docs that is not a readable directory is a real failure
	}
	// Validate requires a README.md in every folder under docs,
	// recursively. Init does not invent the missing ones, but the report
	// names every gap so the first validate is not red without warning.
	_ = filepath.WalkDir(docs, func(p string, d fs.DirEntry, werr error) error {
		if werr != nil || !d.IsDir() {
			return nil
		}
		if _, serr := os.Stat(filepath.Join(p, "README.md")); serr != nil {
			if relDir, rerr := filepath.Rel(root, p); rerr == nil {
				rep.MissingReadmes = append(rep.MissingReadmes, path.Join(filepath.ToSlash(relDir), "README.md"))
			}
		}
		return nil
	})
	// Regeneration runs only at taxonomy depth — the same set of READMEs
	// checkIndexes judges (docs/README.md and docs/<folder>/README.md).
	readmes := []string{"docs/README.md"}
	for _, e := range entries {
		if e.IsDir() {
			readmes = append(readmes, path.Join("docs", e.Name(), "README.md"))
		}
	}
	for _, rel := range readmes {
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(rel))); err != nil {
			continue // already reported by the walk above
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
	// The file's own line endings are preserved: prose outside the
	// markers is never rewritten, and generated lines adopt the file's
	// style (CRLF checkouts stay CRLF).
	orig := string(raw)
	text := orig
	eol := "\n"
	if strings.Contains(orig, "\r\n") {
		eol = "\r\n"
	}
	start := strings.Index(text, indexStart)
	end := strings.Index(text, indexEnd)
	switch {
	case start < 0 && end < 0:
		// A pre-existing taxonomy README without markers would fail
		// validate ("index markers missing"), breaking the green-after-init
		// contract in existing repos. Append an empty block — prose stays
		// untouched, the generated entries land between the markers below.
		if !strings.HasSuffix(text, "\n") {
			text += eol
		}
		text += eol + indexStart + eol + indexEnd + eol
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

	// Keep existing lines that still cover a wanted entry, in order. A
	// line covers an entry the same way checkIndexes reads it: the exact
	// target, or — for a wanted subfolder README — any live link into
	// that subfolder, so a curated descendant entry survives here too.
	coverTarget := func(t string) string {
		if wanted[t] {
			return t
		}
		for w := range wanted {
			if !strings.HasSuffix(w, "/README.md") {
				continue
			}
			sub := path.Dir(w)
			if (t == sub || strings.HasPrefix(t, sub+"/")) && targetExists(root, path.Dir(rel), t) {
				return w
			}
		}
		return ""
	}
	var lines []string
	covered := map[string]bool{}
	for _, line := range strings.Split(text[start+len(indexStart):end], "\n") {
		line = strings.TrimSuffix(line, "\r")
		if strings.TrimSpace(line) == "" {
			continue
		}
		// Every link on the line counts, exactly as the validator reads
		// the block — a curated line may cover several entries at once.
		var covers []string
		for _, m := range indexLinkRe.FindAllStringSubmatch(line, -1) {
			if w := coverTarget(path.Clean(m[1])); w != "" {
				covers = append(covers, w)
			}
		}
		keep := false
		for _, w := range covers {
			if !covered[w] {
				keep = true
			}
		}
		if !keep {
			continue
		}
		for _, w := range covers {
			covered[w] = true
		}
		lines = append(lines, line)
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

	block := indexStart + eol
	if len(lines) > 0 {
		block += strings.Join(lines, eol) + eol
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

func targetExists(root, dir, t string) bool {
	_, err := os.Stat(filepath.Join(root, filepath.FromSlash(path.Join(dir, t))))
	return err == nil
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
