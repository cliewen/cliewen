package skills

import (
	"io/fs"
	"os"
	"path"
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

func TestSanity_EveryMappingSourceHasAGeneratedCounterpart(t *testing.T) {
	mappingsDir := filepath.Join("source", "resources", "clue-extract", "mappings")

	rendered := map[string]bool{}
	for _, file := range mustRender(t) {
		rendered[file.relativePath] = true
	}

	found := 0
	err := filepath.WalkDir(mappingsDir, func(filePath string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		found++
		rel, relErr := filepath.Rel(mappingsDir, filePath)
		if relErr != nil {
			return relErr
		}
		want := path.Join("clue-extract", "mappings", filepath.ToSlash(rel))
		if !rendered[want] {
			t.Errorf("mapping source %s has no generated counterpart %s", filepath.ToSlash(rel), want)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if found == 0 {
		t.Fatal("no mapping sources found under source/resources/clue-extract/mappings")
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

// AC-054 is an extraction-guidance criterion: what it promises is that the
// canonical clue-extract skill instructs criterion-level phasing, so the
// rendered skill text is its evidence. The mechanism that guidance relies on
// is proven separately against the validator — AC-045 for the Human class and
// AC-046 for per-criterion @draft, both in internal/corpus.
func TestAC054_UnitPositive_ExtractionSupportsCriterionLevelPhasing(t *testing.T) {
	extract := ""
	for _, file := range mustRender(t) {
		if file.relativePath == "clue-extract/skill.md" {
			extract = string(file.content)
			break
		}
	}
	for _, want := range []string{
		"Whole-file draft phasing remains available",
		"tag each genuinely not-yet-proven criterion `@draft`",
		"`Test-type: Human` criterion is already proven by naming it in the pull request acceptance brief",
		"supported Go, JVM, or Cucumber evidence",
		"A capability is therefore not the smallest activation unit",
	} {
		if !strings.Contains(extract, want) {
			t.Errorf("clue-extract/skill.md does not carry criterion-level phasing rule %q", want)
		}
	}
}

func TestAC054_UnitNegative_ExtractionRejectsCapabilityOnlyPhasing(t *testing.T) {
	extract := ""
	for _, file := range mustRender(t) {
		if file.relativePath == "clue-extract/skill.md" {
			extract = string(file.content)
			break
		}
	}
	for _, stale := range []string{
		"Activation is per criteria file, not per criterion",
		"a capability is the smallest unit a phasing change can take",
		"Every active acceptance criterion has positive and negative tests, or its capability honestly stays `draft`",
	} {
		if strings.Contains(extract, stale) {
			t.Errorf("clue-extract/skill.md still carries capability-only phasing rule %q", stale)
		}
	}
}

func TestSanity_MethodologyContractChangesMoveEveryLiveCarrierTogether(t *testing.T) {
	for _, file := range mustRender(t) {
		if !strings.HasSuffix(file.relativePath, "/skill.md") {
			continue
		}
		content := string(file.content)
		for _, want := range []string{
			"A decision that changes a methodology contract inventories every live carrier",
			"updates that complete inventory in the same change",
			"Historical analyses, completed plans, and changelog entries remain pinned history",
			"that general obligation remains agent-enforced",
		} {
			if !strings.Contains(content, want) {
				t.Errorf("%s does not carry same-change methodology-carrier rule %q", file.relativePath, want)
			}
		}
	}
}

func TestSanity_VerifyRecognizesTheCompleteEvidenceContract(t *testing.T) {
	verify := ""
	for _, file := range mustRender(t) {
		if file.relativePath == "clue-verify/skill.md" {
			verify = string(file.content)
			break
		}
	}
	for _, want := range []string{
		"supported Go, JVM, or Cucumber evidence",
		"positive/negative direction",
		"`(single-direction)`",
		"a genuine `Human` criterion is named in the acceptance brief as its proof",
		"an individual not-yet-proven criterion carries `@draft`",
		"an unannotated legacy criterion has its one supported reference",
	} {
		if !strings.Contains(verify, want) {
			t.Errorf("clue-verify/skill.md does not recognize evidence-model case %q", want)
		}
	}
	for _, stale := range []string{
		"Every active acceptance criterion has positive and negative tests",
		"its capability honestly stays `draft`",
	} {
		if strings.Contains(verify, stale) {
			t.Errorf("clue-verify/skill.md still carries stale evidence-model rule %q", stale)
		}
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
			"If either side differs",
			`local stopping point such as "commit only"`,
			"not a completed or mergeable change",
			"Review fixes stay on the same branch and PR and repeat the complete updater handoff",
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

func TestAC040_ReviewResultsAreDurableAndCommitBound(t *testing.T) {
	for _, file := range mustRender(t) {
		if !strings.HasSuffix(file.relativePath, "/skill.md") || (!strings.HasPrefix(file.relativePath, "clue-delta/") && !strings.HasPrefix(file.relativePath, "clue-extract/") && !strings.HasPrefix(file.relativePath, "clue-verify/")) {
			continue
		}
		content := string(file.content)
		for _, want := range []string{
			"Every review of an existing hosted PR is bound to its observed head SHA",
			"A clean result applies only to that commit; every substantive edit invalidates it",
			"publish the finding there and leave it unresolved until a hosted commit contains the reviewed repair",
		} {
			if !strings.Contains(content, want) {
				t.Errorf("%s does not carry durable review state %q", file.relativePath, want)
			}
		}
	}
}

func TestAC040_ReviewWithoutResolvableHostStateFailsOpenly(t *testing.T) {
	verify := ""
	for _, file := range mustRender(t) {
		if file.relativePath == "clue-verify/skill.md" {
			verify = string(file.content)
		}
	}
	for _, want := range []string{
		"If the reviewer cannot publish a resolvable finding, report the PR as not merge-ready",
		"never claim a chat-only finding has equivalent protection",
		"the isolated reviewer itself remains read-only",
	} {
		if !strings.Contains(verify, want) {
			t.Errorf("clue-verify/skill.md does not expose the unenforced-review fallback %q", want)
		}
	}
}

func TestAC041_AnyEditorOwnsTheExactFastForwardHandoff(t *testing.T) {
	for _, file := range mustRender(t) {
		if !strings.HasSuffix(file.relativePath, "/skill.md") || (!strings.HasPrefix(file.relativePath, "clue-delta/") && !strings.HasPrefix(file.relativePath, "clue-extract/") && !strings.HasPrefix(file.relativePath, "clue-verify/")) {
			continue
		}
		content := string(file.content)
		for _, want := range []string{
			"Any agent that edits an existing PR becomes the updater for that turn",
			"record its hosted head",
			"push only a normal fast-forward update, never force",
			"After a PR is published, incorporate a newer accepted `main` by merging it into the PR branch with a normal push, never by rewriting hosted history",
			"Resolve satisfied review conversations only after the hosted head contains their reviewed repair",
			"the agent may review or help update an existing PR under the handoff above",
		} {
			if !strings.Contains(content, want) {
				t.Errorf("%s does not carry the exact updater handoff %q", file.relativePath, want)
			}
		}
	}
}

func TestAC041_ConcurrentOrClosedPRStateFailsSafely(t *testing.T) {
	verify := ""
	for _, file := range mustRender(t) {
		if file.relativePath == "clue-verify/skill.md" {
			verify = string(file.content)
		}
	}
	for _, want := range []string{
		"If the head changed or the push is rejected as non-fast-forward",
		"fetch and reconcile without overwriting remote work",
		"If the PR merged or closed, stop and report local work as unpublished",
	} {
		if !strings.Contains(verify, want) {
			t.Errorf("clue-verify/skill.md does not fail safely on contested PR state %q", want)
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
		"operative requirement or declared intent that is violated",
		"the concrete consequence",
		"Apply authoritative decisions and the repository's explicit lifecycle rules",
		"human-controlled merge does not require duplicate human code review",
		"lifecycle-successor evidence satisfies a requirement when the repository declares that transition",
		"lifecycle-correct state are not actionable defects by themselves",
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
