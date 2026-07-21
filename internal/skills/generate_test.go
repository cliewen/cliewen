package skills

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAC028_GenerationProducesMatchingStandaloneSkillTrees(t *testing.T) {
	root := t.TempDir()
	if err := Write(root); err != nil {
		t.Fatal(err)
	}
	drifts, err := Check(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(drifts) != 0 {
		t.Fatalf("freshly generated skills drifted: %v", drifts)
	}

	for _, file := range mustRender(t) {
		agentPath := filepath.Join(root, ".agents", "skills", filepath.FromSlash(file.relativePath))
		templatePath := filepath.Join(root, "internal", "scaffold", "templates", "skills", filepath.FromSlash(file.relativePath))
		agent, agentErr := os.ReadFile(agentPath)
		if agentErr != nil {
			t.Fatal(agentErr)
		}
		embedded, embeddedErr := os.ReadFile(templatePath)
		if embeddedErr != nil {
			t.Fatal(embeddedErr)
		}
		if string(agent) != string(embedded) {
			t.Fatalf("%s differs between generated output trees", file.relativePath)
		}
		if filepath.Base(file.relativePath) == "skill.md" && !strings.Contains(string(agent), "\ncliewen-skill: true\n") {
			t.Fatalf("%s carries no Cliewen ownership marker", file.relativePath)
		}
	}
}

func TestAC028_DriftIsRejected(t *testing.T) {
	tests := map[string]func(*testing.T, string){
		"changed": func(t *testing.T, root string) {
			target := filepath.Join(root, ".agents", "skills", "clue-delta", "skill.md")
			if err := os.WriteFile(target, []byte("edited generated output\n"), 0o644); err != nil {
				t.Fatal(err)
			}
		},
		"missing": func(t *testing.T, root string) {
			target := filepath.Join(root, ".agents", "skills", "clue-delta", "skill.md")
			if err := os.Remove(target); err != nil {
				t.Fatal(err)
			}
		},
		"unexpected": func(t *testing.T, root string) {
			target := filepath.Join(root, ".agents", "skills", "clue-delta", "manual.md")
			if err := os.WriteFile(target, []byte("not generated\n"), 0o644); err != nil {
				t.Fatal(err)
			}
		},
		"changed in template tree": func(t *testing.T, root string) {
			target := filepath.Join(root, "internal", "scaffold", "templates", "skills", "clue-delta", "skill.md")
			if err := os.WriteFile(target, []byte("edited generated output\n"), 0o644); err != nil {
				t.Fatal(err)
			}
		},
	}

	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			root := t.TempDir()
			if err := Write(root); err != nil {
				t.Fatal(err)
			}
			mutate(t, root)
			drifts, err := Check(root)
			if err != nil {
				t.Fatal(err)
			}
			if len(drifts) == 0 {
				t.Fatal("expected generated skill drift to be rejected")
			}
			if !strings.Contains(drifts[0].Path, "clue-delta") {
				t.Fatalf("drift did not name the affected skill: %v", drifts)
			}
		})
	}
}

func TestSanity_CommittedSkillsMatchCanonicalSources(t *testing.T) {
	root := filepath.Join("..", "..")
	drifts, err := Check(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, drift := range drifts {
		t.Error(drift)
	}
}

func TestUnit_ReviewBoundaryRequiresExactHostedHandoff(t *testing.T) {
	rendered := map[string]string{}
	for _, file := range mustRender(t) {
		rendered[file.relativePath] = string(file.content)
	}

	for _, name := range []string{"clue-delta/skill.md", "clue-extract/skill.md", "clue-verify/skill.md"} {
		content := rendered[name]
		for _, want := range []string{
			"authorization and protected-integration boundary",
			"not a demand for duplicate human code review",
			"only a human-controlled PR merge accepts it",
			"A PR alone displays hosted CI but does not enforce it",
			"branch protection makes its required status check a merge precondition",
			"commit every intended edit",
			"`git status --porcelain` to be empty",
			"head branch and SHA equal the current local branch and `HEAD`",
			"if either side differs",
			`local stopping point such as "commit only"`,
			"not a completed or mergeable change",
			"Review fixes stay on the same branch and PR: commit them, rerun local verification and the automatic agentic review loop against that commit, require a clean worktree, push the reviewed commit there",
		} {
			if !strings.Contains(content, want) {
				t.Errorf("%s does not contain review-handoff rule %q", name, want)
			}
		}
	}

	verify := rendered["clue-verify/skill.md"]
	for _, want := range []string{
		"Every intended edit, including each review fix, is committed and `git status --porcelain` is empty",
		"After publishing, the current branch is the ready hosted PR's head branch",
		"reported verification ran against that commit",
	} {
		if !strings.Contains(verify, want) {
			t.Errorf("clue-verify/skill.md does not contain hosted verification item %q", want)
		}
	}
}

func TestUnit_AgenticReviewLoopConvergesOnCurrentCommit(t *testing.T) {
	rendered := map[string]string{}
	for _, file := range mustRender(t) {
		rendered[file.relativePath] = string(file.content)
	}

	verify := rendered["clue-verify/skill.md"]
	for _, want := range []string{
		"never ask the human to clear context or initiate a separate review",
		"start a new read-only reviewer without the implementation conversation",
		"recover a full change's proposal from branch history",
		"label it `in-context fallback`",
		"only actionable findings about correctness, intent mismatch, regressions, security, missing evidence, or unjustified complexity",
		"a previous clean result applies only to the commit it reviewed",
		"Do not publish with unresolved findings or without a clean pass",
		"Report the final review mode and reviewed commit",
	} {
		if !strings.Contains(verify, want) {
			t.Errorf("clue-verify/skill.md does not contain agentic-review rule %q", want)
		}
	}
	commitCandidate := strings.Index(verify, "commit the complete candidate")
	verifyCandidate := strings.Index(verify, "run the applicable local checks against that commit")
	reviewCandidate := strings.Index(verify, "start a new read-only reviewer")
	if commitCandidate < 0 || verifyCandidate <= commitCandidate || reviewCandidate <= verifyCandidate {
		t.Error("clue-verify must commit the candidate, verify that commit, then start agentic review")
	}

	for _, name := range []string{"clue-delta/skill.md", "clue-extract/skill.md"} {
		if !strings.Contains(rendered[name], "automatic agentic review loop") {
			t.Errorf("%s does not invoke the automatic agentic review loop", name)
		}
	}
}

func mustRender(t *testing.T) []renderedFile {
	t.Helper()
	files, err := render()
	if err != nil {
		t.Fatal(err)
	}
	return files
}
