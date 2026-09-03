package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeGoal puts one goal artifact in the goals folder.
func writeGoal(t *testing.T, root, id, status, title string) {
	t.Helper()
	body := "---\nid: " + id + "\ntype: goal\nstatus: " + status + "\nlinks: []\ntitle: " + title + "\n---\n\n# " + id + "\n\nA sentence.\n"
	if err := os.WriteFile(filepath.Join(root, "docs", "goals", id+"-slug.md"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

// writeIndex replaces the goals index block with the given rows.
func writeIndex(t *testing.T, root string, rows ...string) string {
	t.Helper()
	readme := filepath.Join(root, "docs", "goals", "README.md")
	body := "# Goals\n\n<!-- clue:index:start -->\n" + strings.Join(rows, "\n") + "\n<!-- clue:index:end -->\n"
	if err := os.WriteFile(readme, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return readme
}

func regenAndRead(t *testing.T, root, readme string) string {
	t.Helper()
	if _, err := Regen(root); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(readme)
	if err != nil {
		t.Fatal(err)
	}
	return string(raw)
}

// AC-159: a kept row's badge is generated content, so regeneration brings it
// back in step with the artifact while the author's description survives
// beside it untouched.
func TestAC159_UnitPositive_AStaleBadgeIsRefreshed(t *testing.T) {
	root, _ := runInto(t)
	writeGoal(t, root, "G-001", "accepted", "First goal")
	readme := writeIndex(t, root, "- [G-001 — First goal](G-001-slug.md) · `proposed` — What an author decided this actually means.")
	got := regenAndRead(t, root, readme)
	want := "- [G-001 — First goal](G-001-slug.md) · `accepted` — What an author decided this actually means."
	if !strings.Contains(got, want) {
		t.Fatalf("the badge must follow the artifact and the description must survive, got:\n%s", got)
	}
}

// A constraint's badge is its enforcement, not its status (IDR-001), so the
// refresh must read the same value the appended row would have carried.
func TestAC159_UnitPositive_AConstraintKeepsAnEnforcementBadge(t *testing.T) {
	root, _ := runInto(t)
	dir := filepath.Join(root, "docs", "constraints")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	body := "---\nid: C-001\ntype: constraint\nstatus: active\nlinks: []\ntitle: A rule\nsource: G-001\nenforcement: machine\n---\n\n# C-001\n\nA sentence.\n"
	if err := os.WriteFile(filepath.Join(dir, "C-001-rule.md"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	readme := filepath.Join(dir, "README.md")
	row := "- [C-001 — A rule](C-001-rule.md) · `human` — What the author wrote."
	if err := os.WriteFile(readme, []byte("# Constraints\n\n<!-- clue:index:start -->\n"+row+"\n<!-- clue:index:end -->\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	got := regenAndRead(t, root, readme)
	if !strings.Contains(got, "· `machine` — What the author wrote.") {
		t.Fatalf("a constraint's badge must refresh to its enforcement, got:\n%s", got)
	}
	if strings.Contains(got, "`active`") {
		t.Fatalf("a constraint's badge must never become its status, got:\n%s", got)
	}
}

// The refresh is bounded: it never adds a badge an author removed, never
// guesses which artifact owns a multi-target line's badge, and never rewrites
// a row whose artifact it cannot read.
func TestAC159_UnitNegative_AnAmbiguousOrBadgelessRowIsLeftAlone(t *testing.T) {
	t.Run("no badge gains none", func(t *testing.T) {
		root, _ := runInto(t)
		writeGoal(t, root, "G-001", "accepted", "First goal")
		row := "- [G-001 — First goal](G-001-slug.md) — An author who wanted no badge."
		readme := writeIndex(t, root, row)
		if got := regenAndRead(t, root, readme); !strings.Contains(got, row) || strings.Contains(got, "`accepted`") {
			t.Fatalf("a row without a badge must not gain one, got:\n%s", got)
		}
	})

	t.Run("multi-target line is untouched", func(t *testing.T) {
		root, _ := runInto(t)
		writeGoal(t, root, "G-001", "accepted", "First goal")
		writeGoal(t, root, "G-002", "accepted", "Second goal")
		row := "- [G-001 — First goal](G-001-slug.md) · `proposed` — Covered together with [G-002](G-002-slug.md)."
		readme := writeIndex(t, root, row)
		if got := regenAndRead(t, root, readme); !strings.Contains(got, row) {
			t.Fatalf("a line covering two artifacts owns no single badge and must be left alone, got:\n%s", got)
		}
	})

	t.Run("unreadable frontmatter keeps its row", func(t *testing.T) {
		root, _ := runInto(t)
		if err := os.WriteFile(filepath.Join(root, "docs", "goals", "G-001-slug.md"), []byte("no frontmatter here\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		row := "- [G-001 — First goal](G-001-slug.md) · `proposed` — A row whose artifact cannot be read."
		readme := writeIndex(t, root, row)
		if got := regenAndRead(t, root, readme); !strings.Contains(got, row) {
			t.Fatalf("an unreadable artifact must leave its row exactly as it was, got:\n%s", got)
		}
	})
}

// B1 regression: the badge is found by its position after the row's own
// link. A description may itself contain a link, and a scan from the left
// would treat whatever followed that link as the badge — overwriting the
// author's prose in exactly the case the badgeless rule exists to protect.
func TestAC159_UnitNegative_ADescriptionLinkIsNotMistakenForTheBadge(t *testing.T) {
	root, _ := runInto(t)
	writeGoal(t, root, "G-001", "accepted", "First goal")
	row := "- [G-001 — First goal](G-001-slug.md) — see [the rule](../decisions/ADR-046-x.md) · `verified` for why."
	readme := writeIndex(t, root, row)
	if got := regenAndRead(t, root, readme); !strings.Contains(got, row) {
		t.Fatalf("a link inside the description must never be read as the row's badge, got:\n%s", got)
	}
}

// A row naming its own target twice describes one artifact, not two, so the
// multi-artifact guard must not swallow it.
func TestAC159_UnitPositive_ARowNamingOneTargetTwiceStillRefreshes(t *testing.T) {
	root, _ := runInto(t)
	writeGoal(t, root, "G-001", "accepted", "First goal")
	readme := writeIndex(t, root, "- [G-001 — First goal](G-001-slug.md) · `proposed` — also at [here](G-001-slug.md).")
	got := regenAndRead(t, root, readme)
	if !strings.Contains(got, "· `accepted` — also at [here](G-001-slug.md).") {
		t.Fatalf("one artifact named twice is still one artifact, got:\n%s", got)
	}
}

// A row may spell its own link any way that still resolves. Anchoring on the
// literal text left a "./" row silently unrefreshed while every carrier said
// there were exactly three shapes regeneration skips.
func TestAC159_UnitPositive_ARelativeSpellingStillAnchors(t *testing.T) {
	root, _ := runInto(t)
	writeGoal(t, root, "G-001", "accepted", "First goal")
	readme := writeIndex(t, root, "- [G-001 — First goal](./G-001-slug.md) · `proposed` — Written with a leading dot.")
	got := regenAndRead(t, root, readme)
	if !strings.Contains(got, "](./G-001-slug.md) · `accepted` — Written with a leading dot.") {
		t.Fatalf("a row spelled ./x.md must anchor like any other, got:\n%s", got)
	}
}
