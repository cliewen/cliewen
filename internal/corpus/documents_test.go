package corpus

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// AC-213 positive: folder READMEs that carry a type: index header pass, and
// none of them becomes an artifact that owes an id, status, or links. A
// capability's README keeps its identity and stays an artifact.
func TestAC213_UnitPositive_TypedFolderReadmesPassAndStayNonArtifacts(t *testing.T) {
	files := with(validFiles, map[string]string{
		"docs/README.md":                            "---\ntype: index\ntitle: Corpus\n---\n\n# Corpus\n\n<!-- clue:index:start -->\n- [goals/](goals/README.md)\n- [plans/](plans/README.md)\n- [capabilities/](capabilities/README.md)\n<!-- clue:index:end -->\n",
		"docs/capabilities/README.md":               "---\ntype: index\ntitle: Capabilities\n---\n\n# Capabilities\n\n<!-- clue:index:start -->\n- [CAP-001](CAP-001-thing/README.md)\n<!-- clue:index:end -->\n",
		"docs/capabilities/CAP-001-thing/README.md": "---\nid: CAP-001\ntype: capability\nstatus: draft\nlinks: [G-001]\ntitle: Thing\ngoal: G-001\n---\n\n# CAP-001\n",
	})
	c, issues := Scan(writeCorpus(t, files))
	issues = append(issues, Validate(c, Options{})...)
	if len(issues) != 0 {
		t.Fatalf("a corpus of typed folder READMEs should be clean, got %v", issues)
	}
	for _, a := range c.Artifacts {
		if isTaxonomyReadme(a.Path) {
			t.Errorf("folder README %s was scanned as an artifact", a.Path)
		}
	}
	if len(c.ByID["CAP-001"]) != 1 {
		t.Fatalf("a capability README with an id must stay an artifact, got %v", c.ByID["CAP-001"])
	}
}

// AC-213 negative: a header-less corpus file, a folder README of the wrong
// type or without a title, and an artifact without an id are each named.
func TestAC213_UnitNegative_MissingOrWrongHeadersAreNamed(t *testing.T) {
	cases := []struct{ name, path, content, want string }{
		{"readme without frontmatter", "docs/goals/README.md", "# Goals\n\n<!-- clue:index:start -->\n- [G-001](G-001-first.md)\n<!-- clue:index:end -->\n", "docs/goals/README.md: " + MissingIndexHeader},
		{"readme of another type", "docs/goals/README.md", "---\ntype: guide\ntitle: Goals\n---\n\n# Goals\n\n<!-- clue:index:start -->\n- [G-001](G-001-first.md)\n<!-- clue:index:end -->\n", "docs/goals/README.md: folder README must declare type: index"},
		{"readme without a title", "docs/goals/README.md", "---\ntype: index\ntitle: \"\"\n---\n\n# Goals\n\n<!-- clue:index:start -->\n- [G-001](G-001-first.md)\n<!-- clue:index:end -->\n", "docs/goals/README.md: folder README has no title"},
		{"artifact without frontmatter", "docs/goals/G-002-bare.md", "# G-002\n", "docs/goals/G-002-bare.md: missing frontmatter (expected id"},
		{"non-README file with only a document header", "docs/goals/notes.md", "---\ntype: guide\ntitle: Notes\n---\n\n# Notes\n", "docs/goals/notes.md: missing or empty core field(s): id"},
		{"deep README without frontmatter", "docs/goals/archive/README.md", "# Archive\n", "docs/goals/archive/README.md: " + MissingIndexHeader},
		{"change workspace without frontmatter", "changes/CH-001-x/notes.md", "# Notes\n", "changes/CH-001-x/notes.md: missing frontmatter"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assertIssue(t, run(t, with(validFiles, map[string]string{tc.path: tc.content}), false), tc.want)
		})
	}
}

// deliveredFiles is every delivered file with a valid document header, plus a
// managed skill whose entry point and reference both carry one.
var deliveredFiles = map[string]string{
	"AGENTS.md":                          "---\ntype: agent-hub\ntitle: Agent routing hub\n---\n\n# Agent routing hub\n",
	"CLAUDE.md":                          "---\ntype: agent-hub\ntitle: Claude Code entry point\n---\n\n@AGENTS.md\n",
	".clue/evidence/README.md":           "---\ntype: evidence-guide\ntitle: Evidence\n---\n\n# Evidence\n",
	".clue/evidence/frameworks.md":       "---\ntype: evidence-guide\ntitle: Frameworks\n---\n\n# Frameworks\n",
	".github/cliewen-wall.md":            "---\ntype: checklist\ntitle: Wall\n---\n\n# Wall\n",
	".agents/skills/clue-delta/skill.md": "---\ncliewen-skill: true\nversion: 0.1.0\ntype: skill\ntitle: clue-delta\n---\n\n# clue-delta\n",
	".agents/skills/clue-delta/references/change-loop.md": "---\ntype: skill-reference\ntitle: Change loop\n---\n\n## Change loop\n",
}

func deliveredIssues(t *testing.T, files map[string]string) []Issue {
	t.Helper()
	c, issues := Scan(writeCorpus(t, with(validFiles, files)))
	return append(issues, checkDeliveredHeaders(c)...)
}

// AC-214 positive: delivered files and managed skill files with a type and
// title pass, an absent delivered file is not reported, and neither the
// pull-request template nor the adopter's own Markdown is inspected.
func TestAC214_UnitPositive_HeadedDeliveredFilesAndUntouchedAdopterMarkdownPass(t *testing.T) {
	files := with(deliveredFiles, map[string]string{
		".github/pull_request_template.md":    "## Summary\n",
		"src/notes.md":                        "# Adopter notes without a header\n",
		".agents/skills/their-skill/SKILL.md": "---\nname: their-skill\n---\n\n# Theirs\n",
		".agents/skills/their-skill/ref.md":   "# Third-party reference\n",
	})
	if issues := deliveredIssues(t, files); len(issues) != 0 {
		t.Fatalf("headed delivered files should pass, got %v", issues)
	}
	partial := with(deliveredFiles, nil)
	delete(partial, ".github/cliewen-wall.md")
	delete(partial, ".clue/evidence/frameworks.md")
	if issues := deliveredIssues(t, partial); len(issues) != 0 {
		t.Fatalf("an absent delivered file must not be reported, got %v", issues)
	}
}

// AC-214 negative: each delivered file, and each file of a managed skill, that
// is present without a type or title is named.
func TestAC214_UnitNegative_DeliveredFileWithoutTypeOrTitleIsNamed(t *testing.T) {
	cases := map[string]string{
		"AGENTS.md":                          "# Agent routing hub\n",
		"CLAUDE.md":                          "---\ntitle: Claude Code entry point\n---\n\n@AGENTS.md\n",
		".clue/evidence/README.md":           "---\ntype: evidence-guide\n---\n\n# Evidence\n",
		".clue/evidence/frameworks.md":       "# Frameworks\n",
		".github/cliewen-wall.md":            "---\ntype: \"\"\ntitle: Wall\n---\n",
		".agents/skills/clue-delta/skill.md": "---\ncliewen-skill: true\nversion: 0.1.0\n---\n\n# clue-delta\n",
		".agents/skills/clue-delta/references/change-loop.md": "## Change loop\n",
	}
	for rel, content := range cases {
		t.Run(rel, func(t *testing.T) {
			issues := deliveredIssues(t, with(deliveredFiles, map[string]string{rel: content}))
			assertIssue(t, issues, rel+": ")
			if len(issues) != 1 {
				t.Fatalf("only %s should be reported, got %v", rel, issues)
			}
		})
	}
}

// AC-214 negative: a .claude/skills mirror that is a symlink to .agents/skills
// reports a header-less file once, under its real path.
func TestAC214_UnitNegative_SymlinkedSkillMirrorIsNotReportedTwice(t *testing.T) {
	root := writeCorpus(t, with(validFiles, with(deliveredFiles, map[string]string{
		".agents/skills/clue-delta/references/change-loop.md": "## Change loop\n",
	})))
	if err := os.MkdirAll(filepath.Join(root, ".claude"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join("..", ".agents", "skills"), filepath.Join(root, ".claude", "skills")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	c, _ := Scan(root)
	issues := checkDeliveredHeaders(c)
	if len(issues) != 1 || !strings.HasPrefix(issues[0].Path, ".agents/") {
		t.Fatalf("expected one issue under .agents, got %v", issues)
	}
}

// c024Exceptions are the tracked Markdown files C-024 names as deliberately
// header-less: the pull-request templates, which the forge pastes verbatim
// into every new pull request, and the fixtures that model an adopter
// repository from before document headers existed.
var c024Exceptions = []string{
	".github/pull_request_template.md",
	"internal/scaffold/templates/github/pull_request_template.md",
	"internal/migrate/testdata/pre-contract/",
}

// Sanity (C-024): every Markdown file this repository tracks carries
// frontmatter with a type and a title, except the files C-024 names.
func TestSanity_EveryTrackedMarkdownFileCarriesFrontmatter(t *testing.T) {
	root, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(root, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(root)
		if parent == root {
			t.Fatal("repo root with go.mod not found")
		}
		root = parent
	}
	cmd := exec.Command("git", "ls-files", "-z", "--", "*.md")
	cmd.Dir = root
	out, err := cmd.Output()
	if err != nil {
		t.Skipf("git ls-files unavailable: %v", err)
	}
	checked := 0
	for _, rel := range strings.Split(strings.TrimRight(string(out), "\x00"), "\x00") {
		if rel == "" || c024Excepted(rel) {
			continue
		}
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
		if err != nil {
			continue // deleted in the working tree but still in the index
		}
		checked++
		if gap := DocumentHeaderGap(string(data)); gap != "" {
			t.Errorf("%s: %s (C-024)", rel, gap)
		}
	}
	if checked == 0 {
		t.Fatal("no tracked Markdown files were checked")
	}
}

func c024Excepted(rel string) bool {
	for _, e := range c024Exceptions {
		if rel == e || (strings.HasSuffix(e, "/") && strings.HasPrefix(rel, e)) {
			return true
		}
	}
	return false
}
