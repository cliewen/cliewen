package scaffold

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cliewen/cliewen/internal/corpus"
	"gopkg.in/yaml.v3"
)

const prTemplate = ".github/pull_request_template.md"

// headerOf returns a file's frontmatter fields, or nil when it has none.
func headerOf(t *testing.T, text string) map[string]any {
	t.Helper()
	text = strings.ReplaceAll(text, "\r\n", "\n")
	if !strings.HasPrefix(text, "---\n") {
		return nil
	}
	end := strings.Index(text[4:], "\n---\n")
	if end < 0 {
		return nil
	}
	var fields map[string]any
	if err := yaml.Unmarshal([]byte(text[4:4+end]), &fields); err != nil {
		t.Fatalf("frontmatter does not parse: %v\n%s", err, text)
	}
	return fields
}

// writtenMarkdown lists every Markdown file under root, the way an adopter
// would find them after init, following no symlink.
func writtenMarkdown(t *testing.T, root string) []string {
	t.Helper()
	var files []string
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".md") {
			rel, _ := filepath.Rel(root, p)
			files = append(files, filepath.ToSlash(rel))
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return files
}

// AC-215 positive: every Markdown file init and scaffold write, generated
// skills and their mirror included, opens with a type and a title, folder
// READMEs as type: index, and regenerating an index keeps that header.
func TestAC215_UnitPositive_InitWritesEveryMarkdownFileWithFrontmatter(t *testing.T) {
	root, _ := runInto(t)
	if _, err := Run(root); err != nil {
		t.Fatal(err)
	}
	files := writtenMarkdown(t, root)
	if len(files) < 20 {
		t.Fatalf("init wrote suspiciously few Markdown files: %v", files)
	}
	sawSkill, sawMirror := false, false
	for _, rel := range files {
		if rel == prTemplate {
			continue
		}
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
		if err != nil {
			t.Fatal(err)
		}
		fields := headerOf(t, string(data))
		typ, _ := fields["type"].(string)
		title, _ := fields["title"].(string)
		if typ == "" || strings.TrimSpace(title) == "" {
			t.Errorf("%s has no type and title in its frontmatter", rel)
		}
		if strings.HasSuffix(rel, "/README.md") && strings.Count(rel, "/") <= 2 && strings.HasPrefix(rel, "docs/") && typ != corpus.IndexType {
			t.Errorf("folder README %s declares type %q, want index", rel, typ)
		}
		sawSkill = sawSkill || strings.HasPrefix(rel, ".agents/skills/")
		sawMirror = sawMirror || strings.HasPrefix(rel, ".claude/skills/")
	}
	if !sawSkill || !sawMirror {
		t.Fatalf("expected generated skills and their mirror among the checked files: %v", files)
	}

	// Regeneration rewrites the index block of a README that already carries
	// a header and leaves the header byte for byte.
	readme := filepath.Join(root, "docs", "goals", "README.md")
	before, err := os.ReadFile(readme)
	if err != nil {
		t.Fatal(err)
	}
	header := string(before[:strings.Index(string(before), "\n---\n")+5])
	artifact := "---\nid: G-001\ntype: goal\nstatus: proposed\nlinks: []\ntitle: First goal\n---\n\n# G-001\n"
	if err := os.WriteFile(filepath.Join(root, "docs", "goals", "G-001-first.md"), []byte(artifact), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Regen(root); err != nil {
		t.Fatal(err)
	}
	after, err := os.ReadFile(readme)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(string(after), header) || !strings.Contains(string(after), "G-001-first.md") {
		t.Fatalf("regeneration changed the header or skipped the new row:\n%s", after)
	}
	if strings.Count(string(after), "\n---\n") != strings.Count(header, "\n---\n") {
		t.Fatalf("regeneration added a second header:\n%s", after)
	}
}

// AC-215 positive: a folder README without frontmatter gains its type: index
// header, titled from its first heading, with the prose after it unchanged.
func TestAC215_UnitPositive_ScaffoldGivesAHeaderlessFolderReadmeItsHeader(t *testing.T) {
	root, _ := runInto(t)
	readme := filepath.Join(root, "docs", "goals", "README.md")
	const prose = "# Goals, our way\n\nWhy we keep goals.\n"
	if err := os.WriteFile(readme, []byte(prose), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Regen(root); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(readme)
	if err != nil {
		t.Fatal(err)
	}
	want := "---\ntype: index\ntitle: \"Goals, our way\"\n---\n\n" + prose
	if !strings.HasPrefix(string(got), want) {
		t.Fatalf("expected the header before unchanged prose, got:\n%s", got)
	}
}

// AC-215 negative: the pull-request template is the one Markdown file init
// writes without frontmatter, because the forge pastes it into every pull
// request body verbatim.
func TestAC215_UnitNegative_PullRequestTemplateCarriesNoFrontmatter(t *testing.T) {
	root, _ := runInto(t)
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(prTemplate)))
	if err != nil {
		t.Fatal(err)
	}
	if strings.HasPrefix(string(data), "---") {
		t.Fatalf("the pull-request template must not open with frontmatter:\n%s", data)
	}
	if got := WithIndexHeader("---\ntype: guide\ntitle: Kept\n---\n\n# Kept\n", "docs/goals/README.md"); got != "---\ntype: guide\ntitle: Kept\n---\n\n# Kept\n" {
		t.Fatalf("a file that already opens frontmatter must be returned unchanged, got:\n%s", got)
	}
}
