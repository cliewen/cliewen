package migrate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// hubNotice returns the MIG-006 notice in a plan, or nil.
func hubNotice(t *testing.T, root string) *Notice {
	t.Helper()
	plan, err := Plan(root, Options{ReversalCost: "low"})
	if err != nil {
		t.Fatal(err)
	}
	for i, notice := range plan.Notices {
		if notice.Migration == MigrationHubReleaseCheck {
			return &plan.Notices[i]
		}
	}
	return nil
}

// writeHub puts a routing hub under root, or removes it when body is empty.
func writeHub(t *testing.T, root, body string) {
	t.Helper()
	full := filepath.Join(root, "AGENTS.md")
	if body == "" {
		if err := os.Remove(full); err != nil && !os.IsNotExist(err) {
			t.Fatal(err)
		}
		return
	}
	if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

// assertHubUntouched checks that planning left the hub exactly as it was —
// including still absent when want is empty.
func assertHubUntouched(t *testing.T, root, want string) {
	t.Helper()
	got, err := os.ReadFile(filepath.Join(root, "AGENTS.md"))
	if want == "" {
		if err == nil {
			t.Errorf("migration materialized a hub it must only report:\n%s", got)
		}
		return
	}
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != want {
		t.Errorf("an adopter-owned hub was rewritten:\n%s", got)
	}
}

// AC-084 positive: an adopter whose sessions never learn a release is
// available is told so, in both shapes the gap takes. Existing adopters never
// re-run `init`, so without this report the emitted line would reach only
// repositories onboarded after it shipped (PDR-023).
func TestAC084_UnitPositive_MigrateReportsAHubThatNeverAsks(t *testing.T) {
	for _, tc := range []struct {
		name string
		body string // "" means remove the hub
		want string
	}{
		{name: "absent", want: "run `clue init`"},
		{
			name: "no mention at all",
			body: "# Agent routing hub\n\nThis repo runs Cliewen. Classify the work first.\n",
			want: "never asks whether this repository is behind",
		},
		{
			// A hub that names the upgrade skill still never tells the agent
			// to find out whether there is anything to upgrade to.
			name: "the skill without the check",
			body: "# Agent routing hub\n\nUse `clue-upgrade` when you are asked to upgrade.\n",
			want: "never asks whether this repository is behind",
		},
		{
			// A longer subcommand that only starts the same way is a
			// different command, and an agent would run that one.
			name: "a subcommand that only starts the same way",
			body: "# Agent routing hub\n\nRun `clue latest-notes` at some point.\n",
			want: "never asks whether this repository is behind",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := migrationFixture(t, "")
			writeHub(t, root, tc.body)
			notice := hubNotice(t, root)
			if notice == nil {
				t.Fatalf("no MIG-006 notice for the %s case", tc.name)
			}
			if !strings.Contains(notice.Message, tc.want) {
				t.Errorf("notice does not name the remedy %q: %s", tc.want, notice.Message)
			}
			assertHubUntouched(t, root, tc.body)
		})
	}
}

// AC-084 negative: a hub that does ask is silent. Without this, the notice
// would be noise every adopter learns to ignore, which is worse than not
// reporting at all.
func TestAC084_UnitNegative_AHubThatAsksProducesNoNotice(t *testing.T) {
	for _, tc := range []struct {
		name string
		body string
	}{
		{
			// A command is written in backticks by every convention this
			// project follows, so a code span must count as naming it.
			name: "the emitted shape",
			body: "# Agent routing hub\n\nWhen you start, run `clue latest --quiet`; route a non-empty answer to `clue-upgrade`.\n",
		},
		{
			// The loud form is not what the template says and it still asks.
			name: "the loud form",
			body: "# Agent routing hub\n\nRun `clue latest` before you start.\n",
		},
		{
			name: "named in plain prose",
			body: "# Agent routing hub\n\nStart by running clue latest --quiet.\n",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := migrationFixture(t, "")
			writeHub(t, root, tc.body)
			if notice := hubNotice(t, root); notice != nil {
				t.Fatalf("a hub that asks was reported anyway: %s", notice.Message)
			}
			assertHubUntouched(t, root, tc.body)
		})
	}
}

// AC-084: the report never becomes a repair, in preview or in apply, and it
// never blocks the carrier upgrade beside it. The hub is the adopter's own
// routing prose and their repo-local conventions live in it, so rewriting it
// would take their words with the line.
func TestAC084_UnitNegative_MigrateNeverWritesTheHub(t *testing.T) {
	root := migrationFixture(t, "")
	own := "# Ours\n\nOur own routing, with no release check in it.\n"
	writeHub(t, root, own)

	plan, err := Plan(root, Options{ReversalCost: "low"})
	if err != nil {
		t.Fatal(err)
	}
	for _, change := range plan.Changes {
		if change.Path == "AGENTS.md" {
			t.Fatalf("migration planned a write to the adopter-owned hub: %+v", change)
		}
	}
	if len(plan.Findings) > 0 {
		t.Fatalf("a missing release check blocked the whole migration: %+v", plan.Findings)
	}
	if err := Apply(root, plan); err != nil {
		t.Fatal(err)
	}
	assertHubUntouched(t, root, own)

	// The notice survives apply: nothing repaired it, so it is still true.
	if notice := hubNotice(t, root); notice == nil {
		t.Error("the notice disappeared after apply, though no file was repaired")
	}
}
