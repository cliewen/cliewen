package migrate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// entryPointNotice returns the MIG-005 notice in a plan, or nil.
func entryPointNotice(t *testing.T, root string) *Notice {
	t.Helper()
	plan, err := Plan(root, Options{ReversalCost: "low"})
	if err != nil {
		t.Fatal(err)
	}
	for i, notice := range plan.Notices {
		if notice.Migration == MigrationClaudeEntryPoint {
			return &plan.Notices[i]
		}
	}
	return nil
}

// AC-072 positive: an adopter whose Claude Code entry point never reaches the
// routing hub is told so, in both shapes the gap takes — no file at all, and a
// file that mentions the hub without importing it.
func TestAC072_UnitPositive_MigrateReportsAnEntryPointThatNeverReachesTheHub(t *testing.T) {
	for _, tc := range []struct {
		name    string
		content string // "" means write no entry point
		want    string
	}{
		{name: "absent", want: "run `clue init`"},
		{name: "no import at all", content: "# Ours\n\nRun `make it` before pushing.\n", want: "never imports AGENTS.md"},
		{
			name: "mentioned but not imported",
			// Claude Code skips imports inside code spans, so this loads nothing
			// while reading, to a human, exactly like a working pointer.
			content: "# Ours\n\nThe rules live in `@AGENTS.md`.\n",
			want:    "never imports AGENTS.md",
		},
		{
			name:    "shown as an example in a fence",
			content: "# Ours\n\nTo route an agent here you would write:\n\n```markdown\n@AGENTS.md\n```\n\nWe have not done that.\n",
			want:    "never imports AGENTS.md",
		},
		{
			// A fence wrapping a fence: the inner ``` does not close the
			// outer ````, so everything between them is still an example.
			// Reading the inner block as prose would call this repository
			// routed when nothing loads.
			name:    "an example that shows a fence inside a fence",
			content: "# Ours\n\nAn entry point looks like this:\n\n````markdown\n```\n@AGENTS.md\n```\n````\n\nOurs does not.\n",
			want:    "never imports AGENTS.md",
		},
		{
			// A closing fence carries no info string, so the ```go line is
			// content rather than the end of the block. Reading it as the
			// end would hand the rest of the example back as prose and call
			// this repository routed when nothing loads.
			name:    "an example whose inner fence names a language",
			content: "# Ours\n\nAn entry point looks like this:\n\n```\n```go\n@AGENTS.md\n```\n\nOurs does not.\n",
			want:    "never imports AGENTS.md",
		},
		{
			// The path token must stand alone. Claude Code would try to load
			// `AGENTS.mdx`, and `x@AGENTS.md` is not a path it reads at all.
			name:    "a longer path that only starts the same way",
			content: "# Ours\n\n@AGENTS.mdx\n",
			want:    "never imports AGENTS.md",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := migrationFixture(t, "")
			if tc.content != "" {
				writeEntryPoint(t, root, "CLAUDE.md", tc.content)
			}
			notice := entryPointNotice(t, root)
			if notice == nil {
				t.Fatalf("no MIG-005 notice for the %s case", tc.name)
			}
			if !strings.Contains(notice.Message, tc.want) {
				t.Errorf("notice does not name the remedy %q: %s", tc.want, notice.Message)
			}
			assertEntryPointUntouched(t, root, tc.content)
		})
	}
}

// AC-072 negative: an entry point that does reach the hub is silent. Without
// this, the notice would be noise every adopter learns to ignore, which is
// worse than not reporting at all.
func TestAC072_UnitNegative_ARoutedEntryPointProducesNoNotice(t *testing.T) {
	for _, tc := range []struct {
		name string
		rel  string
		body string
	}{
		{name: "root import", rel: "CLAUDE.md", body: "# Entry point\n\n@AGENTS.md\n\n## Mine\n\nUse plan mode under `src/billing/`.\n"},
		{name: "relative spelling", rel: "CLAUDE.md", body: "@./AGENTS.md\n"},
		// Claude Code reads an import as a path token anywhere in the prose,
		// not only from a line that is nothing else. A repository routed this
		// way is routed, and telling its adopter to add the import they
		// already have is the mirror of the failure this migration catches.
		{name: "named in a sentence", rel: "CLAUDE.md", body: "# Ours\n\nSee @AGENTS.md for the routing rules, then read the skill.\n"},
		{name: "named in a list item", rel: "CLAUDE.md", body: "# Ours\n\n- routing: @AGENTS.md\n"},
		// A code span earlier on the line is skipped by the parser; the
		// import after it still loads.
		{name: "after a code span", rel: "CLAUDE.md", body: "Run `make it` before pushing, and follow @AGENTS.md\n"},
		// An import resolves against the directory of the file holding it, so
		// the parent hop is what reaches the hub from `.claude/`.
		{name: "the .claude location", rel: ".claude/CLAUDE.md", body: "@../AGENTS.md\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := migrationFixture(t, "")
			writeEntryPoint(t, root, tc.rel, tc.body)
			if notice := entryPointNotice(t, root); notice != nil {
				t.Fatalf("a routed entry point was reported anyway: %s", notice.Message)
			}
		})
	}

	// Claude Code loads both locations when both exist, so one of them
	// carrying the import is enough. Reporting the other would be a notice
	// about a repository that is already routed.
	t.Run("one location routes while the other does not", func(t *testing.T) {
		root := migrationFixture(t, "")
		writeEntryPoint(t, root, "CLAUDE.md", "# Ours\n\nNo import here.\n")
		writeEntryPoint(t, root, ".claude/CLAUDE.md", "@../AGENTS.md\n")
		if notice := entryPointNotice(t, root); notice != nil {
			t.Fatalf("a routed repository was reported anyway: %s", notice.Message)
		}
	})

	// The other documented bridge: the entry point is the hub. Reading through
	// the link finds AGENTS.md, which imports no copy of itself, so a check
	// that only looked at content would report a repository that is routed.
	t.Run("a symlink to the hub", func(t *testing.T) {
		root := migrationFixture(t, "")
		writeEntryPoint(t, root, "AGENTS.md", "# Hub\n\nRouting lives here.\n")
		symlinkOrSkip(t, filepath.Join(root, "AGENTS.md"), filepath.Join(root, "CLAUDE.md"))
		if notice := entryPointNotice(t, root); notice != nil {
			t.Fatalf("an entry point symlinked to the hub was reported anyway: %s", notice.Message)
		}
	})
}

// AC-072 positive: an import is resolved from the file that holds it, so the
// spelling that routes at the root loads nothing under `.claude/` — it names
// `.claude/AGENTS.md`, a file nothing emits. Accepting it would be the failure
// this migration exists to catch: a report of routing that never arrives.
func TestAC072_UnitPositive_AnImportThatResolvesNowhereIsReported(t *testing.T) {
	root := migrationFixture(t, "")
	writeEntryPoint(t, root, ".claude/CLAUDE.md", "@AGENTS.md\n")
	notice := entryPointNotice(t, root)
	if notice == nil {
		t.Fatal("an import resolving to a file that does not exist was read as routing")
	}
	if notice.Path != ".claude/CLAUDE.md" {
		t.Errorf("the notice names %q rather than the file that carries the broken import", notice.Path)
	}
	if !strings.Contains(notice.Message, "@../AGENTS.md") {
		t.Errorf("the notice does not name the spelling that reaches the hub from there: %s", notice.Message)
	}
}

// symlinkOrSkip creates newname -> oldname, or skips on a host that cannot
// make symlinks at all (an unprivileged Windows shell). The behavior is then
// untestable here, not broken.
func symlinkOrSkip(t *testing.T, oldname, newname string) {
	t.Helper()
	if err := os.Symlink(oldname, newname); err != nil {
		t.Skip("this host cannot create a symlink: ", err)
	}
}

// AC-072: the migration reports and never repairs, in preview or in apply.
// An entry point is adopter-owned prose; a migration that rewrote it would
// discard the Claude-specific instructions the file exists to hold.
func TestAC072_UnitNegative_MigrateNeverWritesTheEntryPoint(t *testing.T) {
	root := migrationFixture(t, "")
	own := "# Ours\n\nNothing about the hub here.\n"
	writeEntryPoint(t, root, "CLAUDE.md", own)

	plan, err := Plan(root, Options{ReversalCost: "low"})
	if err != nil {
		t.Fatal(err)
	}
	for _, change := range plan.Changes {
		if change.Path == "CLAUDE.md" || change.Path == ".claude/CLAUDE.md" {
			t.Fatalf("migration planned a write to an adopter-owned entry point: %+v", change)
		}
	}
	if len(plan.Findings) > 0 {
		t.Fatalf("a missing bridge blocked the whole migration: %+v", plan.Findings)
	}
	if err := Apply(root, plan); err != nil {
		t.Fatal(err)
	}
	assertEntryPointUntouched(t, root, own)

	// The notice survives apply: nothing repaired it, so it is still true.
	if notice := entryPointNotice(t, root); notice == nil {
		t.Error("the notice disappeared after apply, though no file was repaired")
	}
}

func writeEntryPoint(t *testing.T, root, rel, body string) {
	t.Helper()
	full := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

// assertEntryPointUntouched checks that planning left the file exactly as it
// was — including still absent when want is empty.
func assertEntryPointUntouched(t *testing.T, root, want string) {
	t.Helper()
	got, err := os.ReadFile(filepath.Join(root, "CLAUDE.md"))
	if want == "" {
		if err == nil {
			t.Errorf("migration materialized an entry point it must only report:\n%s", got)
		}
		return
	}
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != want {
		t.Errorf("an adopter-owned entry point was rewritten:\n%s", got)
	}
}
