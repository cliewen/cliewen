package ciguards

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

const scopeStepName = "Detect Cliewen change"

func readScopeStep(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile(repoPath(t, ".github/workflows/clue-validation.yml"))
	if err != nil {
		t.Fatal(err)
	}
	var workflow struct {
		Jobs map[string]struct {
			Steps []boundaryStep `yaml:"steps"`
		} `yaml:"jobs"`
	}
	if err := yaml.Unmarshal(data, &workflow); err != nil {
		t.Fatal(err)
	}
	for _, step := range workflow.Jobs["validate"].Steps {
		if step.Name == scopeStepName {
			return step.Run
		}
	}
	t.Fatalf("clue-validation.yml has no %q step", scopeStepName)
	return ""
}

// runScope runs the shipped workflow's scope step against a throwaway
// repository whose history proposes a change and whose head commit carries
// message, and returns the outputs the step wrote.
func runScope(t *testing.T, message string) map[string]string {
	t.Helper()
	repo := t.TempDir()
	git := func(args ...string) string {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = repo
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@example.com",
			"GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@example.com",
		)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
		return strings.TrimSpace(string(out))
	}
	write := func(rel, body string) {
		t.Helper()
		full := filepath.Join(repo, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	git("init", "-q")
	write("README.md", "base\n")
	git("add", "-A")
	git("commit", "-q", "-m", "base")
	base := git("rev-parse", "HEAD")
	write("changes/CH-001-example/proposal.md", "proposal\n")
	git("add", "-A")
	git("commit", "-q", "-m", "propose")
	write("main.go", "package main\n")
	git("add", "-A")
	git("commit", "-q", "-m", message)
	head := git("rev-parse", "HEAD")

	tmp := t.TempDir()
	output := filepath.Join(tmp, "output")
	script := filepath.Join(tmp, "scope.sh")
	if err := os.WriteFile(script, []byte(readScopeStep(t)), 0o644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(shell(t), "--noprofile", "--norc", "-eo", "pipefail", script)
	cmd.Dir = repo
	cmd.Env = append(os.Environ(),
		"BASE_SHA="+base,
		"HEAD_SHA="+head,
		"GITHUB_SHA="+head,
		"GITHUB_OUTPUT="+output,
		"RUNNER_TEMP="+filepath.ToSlash(tmp),
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("the scope step failed: %v\n%s", err, out)
	}
	data, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	outputs := map[string]string{}
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		if name, value, ok := strings.Cut(line, "="); ok {
			outputs[name] = value
		}
	}
	return outputs
}

// AC-211 positive: an override is complete in the current spelling, in the
// legacy spelling adopter history still carries, and in a mix of the two, so
// a commit written under older skills keeps passing after an upgrade.
func TestAC211_UnitPositive_ShippedWorkflowReadsOverrideInEitherSpelling(t *testing.T) {
	for _, tc := range []struct{ name, message string }{
		{"current", "Fix\n\nCliewen-Route: direct\nCliewen-Recommendation: tracked\nCliewen-Override: user chose direct; criterion risk accepted\n"},
		{"legacy", "Fix\n\nCliewen-Route: simple\nCliewen-Recommendation: full\nCliewen-Override: user chose simple; criterion risk accepted\n"},
		{"mixed", "Fix\n\nCliewen-Route: direct\nCliewen-Recommendation: full\nCliewen-Override: user chose direct; risk\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := runScope(t, tc.message)["tracked"]; got != "false" {
				t.Fatalf("a complete %s override still selected tracked-route bookkeeping: tracked=%q", tc.name, got)
			}
		})
	}
}

// AC-211 negative: anything short of a complete override keeps the proposal's
// tracked-route bookkeeping, so the dual spelling never widens what counts as
// an override.
func TestAC211_UnitNegative_IncompleteOrUnknownOverrideStaysTracked(t *testing.T) {
	for _, tc := range []struct{ name, message string }{
		{"no trailers", "Implement\n"},
		{"missing recommendation", "Fix\n\nCliewen-Route: direct\nCliewen-Override: user chose direct; risk\n"},
		{"missing route", "Fix\n\nCliewen-Recommendation: tracked\nCliewen-Override: user chose direct; risk\n"},
		{"empty risk", "Fix\n\nCliewen-Route: direct\nCliewen-Recommendation: tracked\nCliewen-Override: \n"},
		{"unknown route", "Fix\n\nCliewen-Route: quick\nCliewen-Recommendation: tracked\nCliewen-Override: user chose quick; risk\n"},
		{"unknown recommendation", "Fix\n\nCliewen-Route: direct\nCliewen-Recommendation: planned\nCliewen-Override: user chose direct; risk\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := runScope(t, tc.message)["tracked"]; got != "true" {
				t.Fatalf("an incomplete override (%s) skipped tracked-route bookkeeping: tracked=%q", tc.name, got)
			}
		})
	}
}
