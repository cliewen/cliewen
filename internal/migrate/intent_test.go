package migrate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cliewen/cliewen/internal/corpus"
)

func intentNotice(plan MigrationPlan) (Notice, bool) {
	for _, notice := range plan.Notices {
		if notice.Migration == MigrationProductIntent {
			return notice, true
		}
	}
	return Notice{}, false
}

// AC-167: a repository that adopted Cliewen before the vision existed is told
// about the gap, gets the structure that asserts nothing, and stays appliable.
func TestAC167_UnitPositive_MigrationReportsAMissingVisionAndAddsOnlyStructure(t *testing.T) {
	root := migrationFixture(t, "")
	plan, err := Plan(root, Options{ReversalCost: "low"})
	if err != nil {
		t.Fatal(err)
	}
	notice, found := intentNotice(plan)
	if !found || notice.Path != corpus.VisionPath || !strings.Contains(notice.Message, "states no vision") {
		t.Fatalf("migration did not report the missing vision: %+v", plan.Notices)
	}
	// A notice blocks nothing: an established corpus does not go red for an
	// artifact its authors could not have written (ADR-067).
	for _, finding := range plan.Findings {
		if finding.Migration == MigrationProductIntent {
			t.Fatalf("the missing vision became a blocking finding: %+v", finding)
		}
	}
	folder := false
	for _, change := range plan.Changes {
		if change.Path == "docs/use-cases/README.md" && change.Migration == MigrationProductIntent {
			folder = true
		}
	}
	if !folder {
		t.Fatalf("migration did not plan the optional use-case folder: %+v", plan.Changes)
	}
	if err := Apply(root, plan); err != nil {
		t.Fatal(err)
	}
	// The whole point of the split: structure is written, meaning is not.
	if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(corpus.VisionPath))); !os.IsNotExist(err) {
		t.Fatal("migration wrote a vision; nothing in a repository proves why a product exists")
	}
	if _, err := os.Stat(filepath.Join(root, filepath.FromSlash("docs/use-cases/README.md"))); err != nil {
		t.Fatalf("migration did not write the use-case folder: %v", err)
	}
	index, err := os.ReadFile(filepath.Join(root, filepath.FromSlash("docs/README.md")))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(index), "use-cases/") {
		t.Fatalf("the corpus index was left without the folder this migration created:\n%s", index)
	}
}

// AC-167: a repository that already states a direction is told nothing, and
// re-running plans no second copy of the folder.
func TestAC167_UnitNegative_AVisionAndAFolderThatExistProduceNoNoticeOrChange(t *testing.T) {
	root := migrationFixture(t, "")
	for rel, content := range map[string]string{
		corpus.VisionPath:          "---\nid: VIS-001\ntype: vision\nstatus: active\nlinks: []\ntitle: Mine\n---\n\n# VIS-001\n\nMy own direction.\n",
		"docs/use-cases/README.md": "# Use cases\n\n<!-- clue:index:start -->\n<!-- clue:index:end -->\n",
	} {
		full := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	plan, err := Plan(root, Options{ReversalCost: "low"})
	if err != nil {
		t.Fatal(err)
	}
	if notice, found := intentNotice(plan); found {
		t.Fatalf("a corpus that states a direction was told it does not: %+v", notice)
	}
	for _, change := range plan.Changes {
		if change.Migration == MigrationProductIntent && change.Path == "docs/use-cases/README.md" {
			t.Fatalf("migration planned a second copy of an existing folder: %+v", change)
		}
	}
}
