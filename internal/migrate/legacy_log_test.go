package migrate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAC144_UnitNegative_LegacyDecisionRowsBlockWithoutGuessing(t *testing.T) {
	root := migrationFixture(t, "")
	dir := filepath.Join(root, "docs", "decisions")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	log := "---\nid: LOG-001\ntype: log\nstatus: active\nlinks: []\ntitle: Decision log\n---\n\n| Date | Decision | Why | Change/PR |\n|---|---|---|---|\n| 2026-01-01 | Keep the cache local | Fast | CH-001 |\n| 2026-01-02 | Publish signed builds | Trust | CH-002 |\n"
	logPath := filepath.Join(dir, "log.md")
	if err := os.WriteFile(logPath, []byte(log), 0o644); err != nil {
		t.Fatal(err)
	}

	plan, err := Plan(root, Options{ReversalCost: "low"})
	if err != nil {
		t.Fatal(err)
	}
	var inventory []Finding
	for _, finding := range plan.Findings {
		if finding.Migration == MigrationLegacyDecisionLog {
			inventory = append(inventory, finding)
		}
	}
	if len(inventory) != 2 {
		t.Fatalf("legacy inventory has %d rows, want 2: %+v", len(inventory), inventory)
	}
	for _, want := range []string{"Keep the cache local", "Publish signed builds", "ADR, PDR, or IDR", "never guesses"} {
		joined := inventory[0].Message + "\n" + inventory[1].Message
		if !strings.Contains(joined, want) {
			t.Fatalf("legacy inventory does not contain %q: %s", want, joined)
		}
	}
	if err := Apply(root, plan); err == nil {
		t.Fatal("apply with unclassified legacy rows succeeded")
	}
	got, err := os.ReadFile(logPath)
	if err != nil || string(got) != log {
		t.Fatalf("blocked migration changed the legacy log: err=%v\n%s", err, got)
	}
}

func TestAC144_UnitPositive_ConvertedCorpusHasNoLegacyMigrationWork(t *testing.T) {
	root := migrationFixture(t, "")
	plan, err := Plan(root, Options{ReversalCost: "low"})
	if err != nil {
		t.Fatal(err)
	}
	for _, finding := range plan.Findings {
		if finding.Migration == MigrationLegacyDecisionLog {
			t.Fatalf("converted corpus received legacy-log finding: %+v", finding)
		}
	}
	for _, change := range plan.Changes {
		if change.Migration == MigrationLegacyDecisionLog {
			t.Fatalf("legacy-log migration guessed a write: %+v", change)
		}
	}
}
