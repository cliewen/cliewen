package mergehistory

import (
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

type changeHistory struct {
	root           string
	proposal       string
	implementation string
	digest         string
}

func TestSanity_MergeModesDoNotPreserveTheSameArchive(t *testing.T) {
	for _, mode := range []struct {
		name                     string
		preservesOriginalCommits bool
	}{
		{name: "merge-commit", preservesOriginalCommits: true},
		{name: "squash-and-merge", preservesOriginalCommits: false},
		{name: "rebase-and-merge", preservesOriginalCommits: false},
	} {
		t.Run(mode.name, func(t *testing.T) {
			history := createChangeHistory(t)
			// All three modes start from the same diverged state: `main` moved on
			// while the change was under review. That is the realistic case, and
			// comparing the modes from anything else would be unfair — with an
			// undiverged `main`, a rebase is a no-op that trivially preserves the
			// commits it is meant to rewrite.
			advanceMain(t, history.root)

			var head string
			switch mode.name {
			case "merge-commit":
				runGit(t, history.root, "switch", "main")
				runGit(t, history.root, "merge", "--no-ff", "feature", "-m", "Merge CH-090 fixture")
				head = revision(t, history.root, "HEAD")
			case "squash-and-merge":
				runGit(t, history.root, "switch", "main")
				runGit(t, history.root, "merge", "--squash", "feature")
				runGit(t, history.root, "commit", "-m", "Squash CH-090 fixture")
				head = revision(t, history.root, "HEAD")
			case "rebase-and-merge":
				runGit(t, history.root, "switch", "feature")
				runGit(t, history.root, "rebase", "main")
				runGit(t, history.root, "switch", "main")
				runGit(t, history.root, "merge", "--ff-only", "feature")
				head = revision(t, history.root, "HEAD")
			default:
				t.Fatalf("unhandled integration mode %q", mode.name)
			}

			for _, commit := range []struct {
				name string
				sha  string
			}{
				{name: "proposal", sha: history.proposal},
				{name: "implementation", sha: history.implementation},
				{name: "digest", sha: history.digest},
			} {
				got := isAncestor(history.root, commit.sha, head)
				if got != mode.preservesOriginalCommits {
					t.Errorf("%s commit %s reachability = %v, want %v", mode.name, commit.name, got, mode.preservesOriginalCommits)
				}
			}

			if _, err := os.Stat(filepath.Join(history.root, "changes", "CH-090-fixture")); !errors.Is(err, fs.ErrNotExist) {
				t.Fatalf("digested change workspace still exists: %v", err)
			}
			if _, err := os.Stat(filepath.Join(history.root, "docs", "digest.md")); err != nil {
				t.Fatalf("durable digest is missing: %v", err)
			}
		})
	}
}

func createChangeHistory(t *testing.T) changeHistory {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skipf("git is not on PATH: %v", err)
	}
	root := t.TempDir()
	runGit(t, root, "init")
	runGit(t, root, "branch", "-M", "main")
	runGit(t, root, "config", "user.name", "Cliewen fixture")
	runGit(t, root, "config", "user.email", "fixture@example.invalid")

	writeFile(t, root, "docs/archive.md", "The durable corpus starts here.\n")
	commit(t, root, "baseline")
	runGit(t, root, "switch", "-c", "feature")

	writeFile(t, root, "changes/CH-090-fixture/proposal.md", "The accepted proposal.\n")
	writeFile(t, root, "changes/CH-090-fixture/tasks.md", "- [ ] Implement the fixture\n")
	writeFile(t, root, "changes/CH-090-fixture/open-questions.md", "None.\n")
	proposal := commit(t, root, "CH-090 fixture proposal")

	writeFile(t, root, "implementation.txt", "The implemented change.\n")
	implementation := commit(t, root, "CH-090 fixture implementation")

	runGit(t, root, "rm", "-r", "--", "changes/CH-090-fixture")
	writeFile(t, root, "docs/digest.md", "The proposal was digested into the corpus.\n")
	digest := commit(t, root, "CH-090 fixture digest")

	return changeHistory{root: root, proposal: proposal, implementation: implementation, digest: digest}
}

func advanceMain(t *testing.T, root string) {
	t.Helper()
	runGit(t, root, "switch", "main")
	writeFile(t, root, "unrelated-main-change.txt", "Accepted while the feature was under review.\n")
	commit(t, root, "unrelated accepted main change")
}

func writeFile(t *testing.T, root, name, contents string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(name))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}

func commit(t *testing.T, root, message string) string {
	t.Helper()
	runGit(t, root, "add", "--all")
	runGit(t, root, "commit", "-m", message)
	return revision(t, root, "HEAD")
}

func revision(t *testing.T, root, ref string) string {
	t.Helper()
	return strings.TrimSpace(runGit(t, root, "rev-parse", ref))
}

func isAncestor(root, ancestor, descendant string) bool {
	cmd := exec.Command("git", "merge-base", "--is-ancestor", ancestor, descendant)
	cmd.Dir = root
	return cmd.Run() == nil
}

func runGit(t *testing.T, root string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = root
	// The fixture must reproduce the three merge outcomes on any machine, so it
	// runs against no system or global configuration: a maintainer's global
	// commit.gpgsign, core.hooksPath, or merge defaults would otherwise decide
	// whether these commits can be created at all.
	cmd.Env = append(os.Environ(),
		"GIT_CONFIG_NOSYSTEM=1",
		"GIT_CONFIG_GLOBAL="+os.DevNull,
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s failed: %v\n%s", strings.Join(args, " "), err, output)
	}
	return string(output)
}
