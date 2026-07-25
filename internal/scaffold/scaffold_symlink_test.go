package scaffold

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// symlinkOrSkip creates newname -> oldname, or skips the test on a host
// that cannot make symlinks at all (an unprivileged Windows shell). The
// behavior under test is then untestable here, not broken.
func symlinkOrSkip(t *testing.T, oldname, newname string) {
	t.Helper()
	if err := os.Symlink(oldname, newname); err != nil {
		t.Skip("this host cannot create a symlink: ", err)
	}
}

// entryCount is how many names a directory holds, following the link that
// leads to it — the point of the assertion is what the shared tree
// contains, not what the repository says about it.
func entryCount(t *testing.T, dir string) int {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	return len(entries)
}

func mustExist(t *testing.T, root, rel string) {
	t.Helper()
	if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(rel))); err != nil {
		t.Fatalf("expected %s to be created: %v", rel, err)
	}
}

// AC-038: a .claude/skills that is a symlink to a skills tree shared
// across checkouts is left alone — nothing is written through it — while
// the canonical skills and the rest of the convention still land.
func TestAC038_SymlinkedMirrorIsNotWrittenThrough(t *testing.T) {
	root := t.TempDir()
	shared := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".claude"), 0o755); err != nil {
		t.Fatal(err)
	}
	symlinkOrSkip(t, shared, filepath.Join(root, ".claude", "skills"))

	rep, err := Run(root)
	if err != nil {
		t.Fatal(err)
	}
	if n := entryCount(t, shared); n != 0 {
		t.Fatalf("init wrote %d entries into the shared tree behind the link", n)
	}
	if !slices.Contains(rep.Linked, ".claude/skills") {
		t.Fatalf("the report does not name the symlinked mirror: %v", rep.Linked)
	}
	for _, rel := range rep.Created {
		if strings.HasPrefix(rel, ".claude/") {
			t.Fatalf("nothing under the link may be reported as created, got %s", rel)
		}
	}
	mustExist(t, root, ".agents/skills/clue-delta/skill.md")
	mustExist(t, root, "AGENTS.md")
	mustExist(t, root, "docs/README.md")
	mustExist(t, root, ".github/workflows/clue.yml")
}

// AC-038: the rule is about any ancestor below the root, so a symlinked
// .claude blocks the mirror just as .claude/skills does, and is reported
// at the directory that actually is the link.
func TestAC038_SymlinkedAncestorIsDetected(t *testing.T) {
	root := t.TempDir()
	shared := t.TempDir()
	symlinkOrSkip(t, shared, filepath.Join(root, ".claude"))

	rep, err := Run(root)
	if err != nil {
		t.Fatal(err)
	}
	if n := entryCount(t, shared); n != 0 {
		t.Fatalf("init wrote %d entries into the shared tree behind the link", n)
	}
	if !slices.Contains(rep.Linked, ".claude") {
		t.Fatalf("the report names %v, not the linked ancestor .claude", rep.Linked)
	}
	mustExist(t, root, ".agents/skills/clue-delta/skill.md")
}

// AC-038: the canonical location gets the same treatment as the mirror —
// no path is special-cased.
func TestAC038_SymlinkedCanonicalSkillsIsSkipped(t *testing.T) {
	root := t.TempDir()
	shared := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".agents"), 0o755); err != nil {
		t.Fatal(err)
	}
	symlinkOrSkip(t, shared, filepath.Join(root, ".agents", "skills"))

	rep, err := Run(root)
	if err != nil {
		t.Fatal(err)
	}
	if n := entryCount(t, shared); n != 0 {
		t.Fatalf("init wrote %d entries into the shared tree behind the link", n)
	}
	if !slices.Contains(rep.Linked, ".agents/skills") {
		t.Fatalf("the report does not name the symlinked canonical folder: %v", rep.Linked)
	}
	mustExist(t, root, ".claude/skills/clue-delta/SKILL.md")
}

// AC-038: a link whose target is gone is still a link. Skipping it keeps
// a dangling mirror a reported fact rather than a hard failure.
func TestAC038_DanglingSymlinkIsSkippedNotFollowed(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".claude"), 0o755); err != nil {
		t.Fatal(err)
	}
	symlinkOrSkip(t, filepath.Join(t.TempDir(), "gone"), filepath.Join(root, ".claude", "skills"))

	rep, err := Run(root)
	if err != nil {
		t.Fatalf("a dangling link must be reported, not fail the run: %v", err)
	}
	if !slices.Contains(rep.Linked, ".claude/skills") {
		t.Fatalf("the report does not name the dangling link: %v", rep.Linked)
	}
	mustExist(t, root, ".agents/skills/clue-delta/skill.md")
}

// AC-038 negative: the root itself is never inspected. A checkout reached
// through a link is ordinary and must be initialized in full.
func TestAC038_RootReachedThroughASymlinkIsFullyInitialized(t *testing.T) {
	real := t.TempDir()
	link := filepath.Join(t.TempDir(), "checkout")
	symlinkOrSkip(t, real, link)

	rep, err := Run(link)
	if err != nil {
		t.Fatal(err)
	}
	if len(rep.Linked) != 0 {
		t.Fatalf("a symlinked root is not a blocked ancestor, got %v", rep.Linked)
	}
	mustExist(t, real, ".claude/skills/clue-delta/SKILL.md")
	mustExist(t, real, ".agents/skills/clue-delta/skill.md")
}

// AC-038: init stays idempotent across the new category — a second run
// reports the same link and still writes nothing through it.
func TestAC038_SecondRunReportsTheSameLink(t *testing.T) {
	root := t.TempDir()
	shared := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".claude"), 0o755); err != nil {
		t.Fatal(err)
	}
	symlinkOrSkip(t, shared, filepath.Join(root, ".claude", "skills"))

	first, err := Run(root)
	if err != nil {
		t.Fatal(err)
	}
	second, err := Run(root)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(first.Linked, second.Linked) {
		t.Fatalf("the linked set changed between runs: %v then %v", first.Linked, second.Linked)
	}
	if len(second.Created) != 0 {
		t.Fatalf("a second run created %v", second.Created)
	}
	if n := entryCount(t, shared); n != 0 {
		t.Fatalf("the second run wrote %d entries into the shared tree", n)
	}
}
