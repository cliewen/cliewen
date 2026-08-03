package scaffold

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// AC-083 positive: the materialized hub asks the agent to check whether the
// repository is behind. The hub is the one file every assistant reads when a
// session starts, which is what makes it the channel — `clue latest` otherwise
// only helps someone who already knows the command exists, and `clue-upgrade`
// only runs when asked to upgrade, which is the request nobody makes when
// nothing told them there was anything to upgrade to (PDR-023).
func TestAC083_UnitPositive_InitEmitsAHubThatAsksWhetherWeAreBehind(t *testing.T) {
	root, _ := runInto(t)
	hub, err := os.ReadFile(filepath.Join(root, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	content := string(hub)
	if !strings.Contains(content, "clue latest --quiet") {
		t.Fatalf("the emitted hub never names the quiet release check, so no session start asks:\n%s", content)
	}
	// The answer routes to the skill that holds the human choice, not to an
	// installation command: one text shipped to every repository cannot name
	// the platform each agent is running on, and the upgrade is a reviewed
	// change a human accepts.
	if !strings.Contains(content, "clue-upgrade") {
		t.Errorf("the hub does not route a non-empty answer to the upgrade skill:\n%s", content)
	}
	for _, forbidden := range []string{"install.sh", "install.ps1", "go install"} {
		if strings.Contains(content, forbidden) {
			t.Errorf("the hub reproduces the installation command %q, which cannot be right for every machine", forbidden)
		}
	}
}

// AC-083 negative: nothing is emitted for any assistant's runtime behaviour.
// Cliewen is agent-agnostic, and a scaffold that ships executable
// configuration for one vendor has taken a position on an adopter's tooling
// that no methodology should take (PDR-023). The inert `CLAUDE.md` pointer
// PDR-022 permits is prose and is deliberately not caught by this.
func TestAC083_UnitNegative_InitEmitsNoVendorConfiguration(t *testing.T) {
	root, rep := runInto(t)
	for _, rel := range rep.Created {
		if strings.HasSuffix(rel, ".json") {
			t.Errorf("init emitted %s; no assistant's configuration file may ship with the scaffold", rel)
		}
	}
	// Belt and braces: walk what actually landed, in case a future template
	// lands somewhere the report does not name.
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, rerr := filepath.Rel(root, p)
		if rerr != nil {
			return rerr
		}
		if strings.HasSuffix(rel, ".json") {
			t.Errorf("init materialized %s; the scaffold emits no vendor configuration", filepath.ToSlash(rel))
		}
		data, rerr := os.ReadFile(p)
		if rerr != nil {
			return rerr
		}
		// A hook declaration is the shape this rule exists to keep out,
		// whatever file it might be smuggled into.
		if strings.Contains(string(data), `"hooks"`) {
			t.Errorf("%s declares hooks; the scaffold takes no position on what an assistant runs", filepath.ToSlash(rel))
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
