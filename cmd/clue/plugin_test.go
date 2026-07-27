package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/cliewen/cliewen/internal/skills"
)

// The managed lifecycle skills: committed repository files stamped with the
// version of the binary that wrote them (ADR-011, ADR-022). ADR-031 keeps
// every one of them out of the plugin, because a plugin's components are
// cache copies that `checkSkillVersions` cannot see.
var managedSkills = []string{"clue-analysis", "clue-delta", "clue-extract", "clue-plan", "clue-verify"}

type marketplaceManifest struct {
	Name    string `json:"name"`
	Owner   struct {
		Name string `json:"name"`
	} `json:"owner"`
	Plugins []struct {
		Name    string          `json:"name"`
		Source  json.RawMessage `json:"source"`
		Version string          `json:"version"`
	} `json:"plugins"`
}

func repoPath(parts ...string) string {
	return filepath.Join(append([]string{"..", ".."}, parts...)...)
}

func readMarketplace(t *testing.T) marketplaceManifest {
	t.Helper()
	data, err := os.ReadFile(repoPath(".claude-plugin", "marketplace.json"))
	if err != nil {
		t.Fatalf("marketplace manifest not found: %v — `/plugin marketplace add cliewen/cliewen` resolves this file", err)
	}
	var m marketplaceManifest
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("marketplace manifest is not valid JSON: %v — Claude Code would reject the marketplace outright", err)
	}
	return m
}

// AC-039: the marketplace resolves to a committed tree that ships exactly one
// skill, the bootstrap. The source is a relative path, so the plugin travels
// with the marketplace and cannot drift away from the repository it documents.
func TestAC039_PluginShipsOnlyTheBootstrapSkill(t *testing.T) {
	m := readMarketplace(t)
	if m.Name == "" || m.Owner.Name == "" {
		t.Error("marketplace manifest is missing a required name or owner — the catalog will not load")
	}
	if len(m.Plugins) != 1 {
		t.Fatalf("marketplace lists %d plugins, want exactly 1 — ADR-031 publishes one bootstrap", len(m.Plugins))
	}

	var source string
	if err := json.Unmarshal(m.Plugins[0].Source, &source); err != nil {
		t.Fatalf("plugin source is not a relative path: %v — a remote source would let the published plugin drift from this repository", err)
	}
	if !strings.HasPrefix(source, "./") || strings.Contains(source, "..") {
		t.Fatalf("plugin source %q must be a path under the marketplace root starting with ./", source)
	}

	tree := repoPath(filepath.FromSlash(strings.TrimPrefix(source, "./")))
	if _, err := os.Stat(filepath.Join(tree, ".claude-plugin", "plugin.json")); err != nil {
		t.Fatalf("plugin manifest missing at the source the marketplace names: %v", err)
	}

	entries, err := os.ReadDir(filepath.Join(tree, "skills"))
	if err != nil {
		t.Fatalf("plugin ships no skills directory: %v — the bootstrap is the whole plugin", err)
	}
	var shipped []string
	for _, e := range entries {
		if e.IsDir() {
			shipped = append(shipped, e.Name())
		}
	}
	if len(shipped) != 1 || shipped[0] != "setup" {
		t.Fatalf("plugin ships skills %v, want exactly [setup] — ADR-031 ships one thin bootstrap", shipped)
	}

	body, err := os.ReadFile(filepath.Join(tree, "skills", "setup", "SKILL.md"))
	if err != nil {
		t.Fatalf("bootstrap SKILL.md not found: %v", err)
	}
	// The bootstrap's own reason to exist: install, verify, then ask. A
	// version literal here would pin new users to a release nobody bumps,
	// and running `clue init` unasked scaffolds a repository the user has
	// not yet decided should become a Cliewen repository.
	for _, want := range []string{"install.sh", "install.ps1", "clue version", "clue init"} {
		if !strings.Contains(string(body), want) {
			t.Errorf("bootstrap does not mention %q — it must install through the published channels, confirm a release version, and offer init", want)
		}
	}
}

// AC-039 negative: nothing the plugin ships may be a managed lifecycle skill,
// carry the ownership marker, or pin a `clue` version — and the tree stays out
// of the generator's way.
func TestAC039_PluginCarriesNoManagedSkillAndNoVersionPin(t *testing.T) {
	m := readMarketplace(t)
	if m.Plugins[0].Version != "" {
		t.Error("marketplace pins a plugin version — ADR-031 omits it so the commit is the version and no second stamp can be forgotten")
	}

	tree := repoPath("plugins", "cliewen")

	manifest, err := os.ReadFile(filepath.Join(tree, ".claude-plugin", "plugin.json"))
	if err != nil {
		t.Fatalf("plugin manifest not found: %v", err)
	}
	var declared map[string]any
	if err := json.Unmarshal(manifest, &declared); err != nil {
		t.Fatalf("plugin manifest is not valid JSON: %v", err)
	}
	if _, pinned := declared["version"]; pinned {
		t.Error("plugin.json declares a version — ADR-031 omits it, leaving the commit as the version")
	}

	for _, name := range managedSkills {
		if _, err := os.Stat(filepath.Join(tree, "skills", name)); err == nil {
			t.Errorf("plugin ships the managed skill %q — it belongs in .agents/skills/, written by clue init and stamped with the binary's version, where checkSkillVersions can see it", name)
		}
	}

	// A version literal anywhere in the plugin tree would outlive the release
	// it names. `CLUE_VERSION` is the escape hatch the install scripts read,
	// and the bootstrap must not set it.
	pin := regexp.MustCompile(`(?i)CLUE_VERSION\s*=|@v[0-9]+\.[0-9]+\.[0-9]+|clue@[0-9]`)
	marker := regexp.MustCompile(`(?m)^cliewen-skill:`)
	err = filepath.WalkDir(tree, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil || entry.IsDir() {
			return walkErr
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(tree, path)
		if pin.Match(data) {
			t.Errorf("%s pins a clue version — the bootstrap installs the newest release, which is the one whose skills clue init writes", rel)
		}
		if marker.Match(data) {
			t.Errorf("%s carries the cliewen-skill ownership marker — that marker means a version-locked skill under .agents/skills/, which a cached plugin copy is not", rel)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walking the plugin tree: %v", err)
	}

	// Hand-authored on purpose: `go generate ./internal/skills` must neither
	// write this tree nor report it as an unexpected file.
	for _, owned := range skills.OutputRoots() {
		if strings.HasPrefix(filepath.ToSlash(filepath.Join("plugins", "cliewen")), owned) {
			t.Errorf("the plugin tree sits inside the generator-owned directory %q", owned)
		}
	}
}
