package migrate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ownedDelivered is a header-less copy of each delivered file the adopter owns.
var ownedDelivered = map[string]string{
	"AGENTS.md":                    "# Our hub\n\nRun `clue latest --quiet` first.\n",
	"CLAUDE.md":                    "# Entry point\n\n@AGENTS.md\n",
	".clue/evidence/README.md":     "# Our evidence notes\n",
	".clue/evidence/frameworks.md": "# Our frameworks\n",
	".github/cliewen-wall.md":      "# Our wall checklist\n",
}

func writeFiles(t *testing.T, root string, files map[string]string) {
	t.Helper()
	for rel, body := range files {
		full := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func readFile(t *testing.T, root, rel string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// AC-216 positive: applying the plan gives each header-less folder README a
// type: index header titled from its first heading with its body unchanged,
// and reports every header-less delivered file the adopter owns by path with
// the header it needs, without writing any of them.
func TestAC216_UnitPositive_MigrateHeadsFolderReadmesAndReportsOwnedFiles(t *testing.T) {
	root := migrationFixture(t, "")
	writeFiles(t, root, ownedDelivered)
	before := map[string]string{}
	for _, rel := range []string{"docs/README.md", "docs/analysis/README.md"} {
		before[rel] = readFile(t, root, rel)
		if strings.HasPrefix(before[rel], "---") {
			t.Fatalf("fixture %s already carries frontmatter", rel)
		}
	}
	plan, err := Plan(root, Options{ReversalCost: "low"})
	if err != nil {
		t.Fatal(err)
	}
	reported := map[string]string{}
	for _, notice := range plan.Notices {
		if notice.Migration == MigrationDeliveredHeaders {
			reported[notice.Path] = notice.Message
		}
	}
	for rel := range ownedDelivered {
		message, ok := reported[rel]
		if !ok {
			t.Errorf("no MIG-020 notice for header-less %s", rel)
			continue
		}
		if !strings.Contains(message, "type: "+deliveredHeaderTypes[rel]) || !strings.Contains(message, "does not edit") {
			t.Errorf("notice for %s does not name the header it needs: %s", rel, message)
		}
	}
	if err := Apply(root, plan); err != nil {
		t.Fatal(err)
	}
	for rel, body := range before {
		got := readFile(t, root, rel)
		title := strings.TrimPrefix(strings.SplitN(body, "\n", 2)[0], "# ")
		if !strings.HasPrefix(got, "---\ntype: index\ntitle: "+title+"\n---\n\n") {
			t.Errorf("%s did not gain its type: index header titled %q:\n%s", rel, title, got)
		}
		if !strings.Contains(got, strings.SplitN(body, "<!--", 2)[0]) {
			t.Errorf("%s lost its prose:\n%s", rel, got)
		}
	}
	for rel, body := range ownedDelivered {
		if got := readFile(t, root, rel); got != body {
			t.Errorf("migration rewrote %s:\n%s", rel, got)
		}
	}
}

// AC-216 negative: a README or delivered file that already carries a header is
// neither changed nor reported, and a second run after applying plans no
// further header change.
func TestAC216_UnitNegative_HeadedFilesAreLeftAloneAndTheRepairIsIdempotent(t *testing.T) {
	root := migrationFixture(t, "")
	headed := map[string]string{
		"AGENTS.md":                    "---\ntype: agent-hub\ntitle: Our hub\n---\n\n# Our hub\n",
		"CLAUDE.md":                    "---\ntype: agent-hub\ntitle: Entry point\n---\n\n@AGENTS.md\n",
		".clue/evidence/README.md":     "---\ntype: evidence-guide\ntitle: Notes\n---\n\n# Notes\n",
		".clue/evidence/frameworks.md": "---\ntype: evidence-guide\ntitle: Frameworks\n---\n\n# Frameworks\n",
		".github/cliewen-wall.md":      "---\ntype: checklist\ntitle: Wall\n---\n\n# Wall\n",
	}
	writeFiles(t, root, headed)
	plan, err := Plan(root, Options{ReversalCost: "low"})
	if err != nil {
		t.Fatal(err)
	}
	for _, notice := range plan.Notices {
		if notice.Migration == MigrationDeliveredHeaders {
			t.Errorf("a headed delivered file was reported: %+v", notice)
		}
	}
	if err := Apply(root, plan); err != nil {
		t.Fatal(err)
	}
	again, err := Plan(root, Options{ReversalCost: "low"})
	if err != nil {
		t.Fatal(err)
	}
	for _, change := range again.Changes {
		if change.Migration == MigrationIndexHeaders || strings.Contains(change.Description, "type: index header") {
			t.Errorf("a README that already carries its header was planned again: %s %s", change.Path, change.Description)
		}
	}
	for rel, body := range headed {
		if got := readFile(t, root, rel); got != body {
			t.Errorf("migration changed headed %s:\n%s", rel, got)
		}
	}
}
