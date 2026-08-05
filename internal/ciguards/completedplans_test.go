// Package ciguards holds the evidence for the CI steps that enforce rules
// `clue validate` deliberately cannot. The judge reads a repository state and
// never a transition (ADR-044), so a rule about what a change *did* is held by
// a machine that is allowed to have a base — and the machine that has one is
// the workflow. These tests are that machine's tests.
package ciguards

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

const completedPlan = `---
id: P-001
type: plan
status: completed
links: []
title: A finished campaign
---

# P-001

| ID | Milestone | Status | Evidence |
|---|---|---|---|
| M-001 | do it | ` + "`done`" + ` | CH-001 |
`

const activePlan = `---
id: P-002
type: plan
status: active
links: []
title: A running campaign
---

# P-002

| ID | Milestone | Status | Evidence |
|---|---|---|---|
| M-002 | do it | ` + "`todo`" + ` | |
`

// TestUnit_CompletedPlanGuardFreezesAFinishedCampaign proves the guard fires
// on the edit it exists to stop, and stays quiet for the three edits that look
// similar and are not: an active plan, a brand-new plan, and the digest that
// closes a campaign.
func TestUnit_CompletedPlanGuardFreezesAFinishedCampaign(t *testing.T) {
	for _, tc := range []struct {
		name     string
		change   func(t *testing.T, root string)
		wantFail bool
		diverged bool
	}{
		{
			name: "editing a completed plan fails",
			change: func(t *testing.T, root string) {
				write(t, root, "docs/plans/P-001-first.md", strings.Replace(completedPlan, "do it", "do it differently", 1))
			},
			wantFail: true,
		},
		{
			name: "deleting a completed plan fails",
			change: func(t *testing.T, root string) {
				if err := os.Remove(filepath.Join(root, "docs", "plans", "P-001-first.md")); err != nil {
					t.Fatal(err)
				}
			},
			wantFail: true,
		},
		{
			name: "editing an active plan passes",
			change: func(t *testing.T, root string) {
				write(t, root, "docs/plans/P-002-second.md", strings.Replace(activePlan, "`todo`", "`done`", 1))
			},
		},
		{
			// The case that would break the project if the guard read the
			// working tree instead of the base: closing a campaign sets
			// `completed` on a file that is `active` on the base.
			name: "the digest that closes a campaign passes",
			change: func(t *testing.T, root string) {
				write(t, root, "docs/plans/P-002-second.md", strings.Replace(activePlan, "status: active", "status: completed", 1))
			},
		},
		{
			name: "adding a plan passes",
			change: func(t *testing.T, root string) {
				write(t, root, "docs/plans/P-003-third.md", strings.Replace(activePlan, "P-002", "P-003", -1))
			},
		},
		{
			// The case a two-dot diff gets wrong. The base branch closed a
			// campaign after this branch left it; the branch never touched
			// that file, and a guard comparing against the branch tip would
			// fail it for somebody else's digest.
			name:     "a plan the base branch completed after the fork passes",
			diverged: true,
			change: func(t *testing.T, root string) {
				write(t, root, "docs/plans/P-003-third.md", strings.Replace(activePlan, "P-002", "P-003", -1))
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root, base := fixture(t)
			if tc.diverged {
				base = advanceBaseBranch(t, root)
			}
			tc.change(t, root)
			git(t, root, "add", "-A")
			git(t, root, "commit", "-m", "the change under test")
			head := strings.TrimSpace(run(t, root, "git", "rev-parse", "HEAD"))

			out, err := runGuard(t, root, base, head)
			if tc.wantFail && err == nil {
				t.Fatalf("expected the guard to fail; output: %s", out)
			}
			if !tc.wantFail && err != nil {
				t.Fatalf("expected the guard to pass, got %v; output: %s", err, out)
			}
			if tc.wantFail && !strings.Contains(out, "C-008") {
				t.Fatalf("a failure names the constraint it enforces; output: %s", out)
			}
		})
	}
}

// TestUnit_CompletedPlanGuardReadsTheStatusYAMLWrote proves the freeze holds
// however the base spelled the value. A quoted status is the same plan to
// `clue validate`, which unmarshals it, so a guard matching only the bare word
// would wave through exactly the edit it exists to stop — and say nothing.
func TestUnit_CompletedPlanGuardReadsTheStatusYAMLWrote(t *testing.T) {
	for _, form := range []string{`"completed"`, `'completed'`, `completed  `, `completed # closed by CH-001`} {
		t.Run(form, func(t *testing.T) {
			root := t.TempDir()
			git(t, root, "init", "-q", "-b", "main")
			git(t, root, "config", "user.email", "guard@example.invalid")
			git(t, root, "config", "user.name", "Guard Fixture")
			plan := strings.Replace(completedPlan, "status: completed", "status: "+form, 1)
			write(t, root, "docs/plans/P-001-first.md", plan)
			git(t, root, "add", "-A")
			git(t, root, "commit", "-m", "base")
			base := strings.TrimSpace(run(t, root, "git", "rev-parse", "HEAD"))

			write(t, root, "docs/plans/P-001-first.md", strings.Replace(plan, "do it", "do it differently", 1))
			git(t, root, "add", "-A")
			git(t, root, "commit", "-m", "the change under test")
			head := strings.TrimSpace(run(t, root, "git", "rev-parse", "HEAD"))

			out, err := runGuard(t, root, base, head)
			if err == nil {
				t.Fatalf("expected the guard to fail for %q; output: %s", form, out)
			}
			if !strings.Contains(out, "C-008") {
				t.Fatalf("a failure names the constraint it enforces; output: %s", out)
			}
		})
	}
}

// advanceBaseBranch moves the base branch on after the working branch forked
// from it, closing a campaign there — the ordinary situation once any digest
// lands while a pull request is open. It returns the new base branch tip, which
// is what a forge reports as `pull_request.base.sha`.
func advanceBaseBranch(t *testing.T, root string) string {
	t.Helper()
	fork := strings.TrimSpace(run(t, root, "git", "rev-parse", "HEAD"))
	git(t, root, "switch", "-q", "-c", "work")
	git(t, root, "switch", "-q", "main")
	write(t, root, "docs/plans/P-002-second.md", strings.Replace(activePlan, "status: active", "status: completed", 1))
	git(t, root, "add", "-A")
	git(t, root, "commit", "-q", "-m", "somebody else's digest closes P-002")
	tip := strings.TrimSpace(run(t, root, "git", "rev-parse", "HEAD"))
	git(t, root, "switch", "-q", "work")
	if fork == tip {
		t.Fatal("the base branch did not actually advance")
	}
	return tip
}

// fixture builds a repository whose base commit holds one completed and one
// active plan, and returns the root and the base revision.
func fixture(t *testing.T) (root, base string) {
	t.Helper()
	root = t.TempDir()
	git(t, root, "init", "-q", "-b", "main")
	git(t, root, "config", "user.email", "guard@example.invalid")
	git(t, root, "config", "user.name", "Guard Fixture")
	write(t, root, "docs/plans/P-001-first.md", completedPlan)
	write(t, root, "docs/plans/P-002-second.md", activePlan)
	git(t, root, "add", "-A")
	git(t, root, "commit", "-m", "base")
	return root, strings.TrimSpace(run(t, root, "git", "rev-parse", "HEAD"))
}

func write(t *testing.T, root, rel, content string) {
	t.Helper()
	full := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func git(t *testing.T, root string, args ...string) {
	t.Helper()
	run(t, root, "git", args...)
}

func run(t *testing.T, root, name string, args ...string) string {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s %s: %v\n%s", name, strings.Join(args, " "), err, out)
	}
	return string(out)
}

// runGuard runs the real script the workflow runs. Nothing is reimplemented
// here: a test of a paraphrase of the guard would prove nothing about the step
// that actually runs on a pull request.
func runGuard(t *testing.T, root, base, head string) (string, error) {
	t.Helper()
	cmd := exec.Command(shell(t), repoPath(t, ".github/scripts/completed-plans.sh"), base, head)
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// shell resolves a POSIX shell. On Windows that is Git's own bash, which sits
// beside the git executable — never the WSL bash on the system path, which
// would run in a different filesystem than the fixture.
func shell(t *testing.T) string {
	t.Helper()
	if runtime.GOOS != "windows" {
		path, err := exec.LookPath("bash")
		if err != nil {
			t.Fatalf("no bash on PATH: %v", err)
		}
		return path
	}
	gitExe, err := exec.LookPath("git")
	if err != nil {
		t.Fatalf("no git on PATH: %v", err)
	}
	// git.exe resolves to either <install>/cmd or <install>/mingw64/bin
	// depending on how Git for Windows was put on the path, so the install
	// root is found by ascending until a bin/bash.exe appears.
	for dir := filepath.Dir(gitExe); ; {
		bash := filepath.Join(dir, "bin", "bash.exe")
		if _, err := os.Stat(bash); err == nil {
			return bash
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatalf("Git bash not found above %s", gitExe)
		}
		dir = parent
	}
}

// repoPath resolves a repository-relative path from the package directory.
func repoPath(t *testing.T, rel string) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return filepath.ToSlash(filepath.Join(dir, filepath.FromSlash(rel)))
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("go.mod not found above the package directory")
		}
		dir = parent
	}
}

// TestSanity_WorkflowRunsTheCompletedPlanGuard keeps the tested script and the
// step that runs it from drifting apart: a guard nobody calls is not a guard.
func TestSanity_WorkflowRunsTheCompletedPlanGuard(t *testing.T) {
	data, err := os.ReadFile(repoPath(t, ".github/workflows/ci.yml"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{".github/scripts/completed-plans.sh", "BASE_SHA"} {
		if !strings.Contains(string(data), want) {
			t.Fatalf("CI does not run the completed-plan guard: %q missing from ci.yml", want)
		}
	}
	// The workflow invokes the script directly, so a committed mode without the
	// executable bit is `Permission denied` on every run — a guard that fails
	// open on the runner while every local test passes. Git's index carries the
	// mode; the working tree's does not survive a Windows checkout.
	mode := strings.Fields(run(t, repoRoot(t), "git", "ls-files", "-s", ".github/scripts/completed-plans.sh"))
	if len(mode) == 0 || mode[0] != "100755" {
		t.Fatalf("the guard is not committed executable: %v", mode)
	}
}

// repoRoot is the repository the tests run against.
func repoRoot(t *testing.T) string {
	t.Helper()
	return filepath.Dir(filepath.FromSlash(repoPath(t, "go.mod")))
}
