package ciguards

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

const boundaryStepName = "Report whether the merge boundary is enforced"

// A rule set satisfying PDR-021, pretty-printed so the probe's whitespace
// handling is exercised rather than assumed from GitHub's compact output.
const enforcedRules = `[
  {"type": "deletion", "ruleset_source_type": "Repository", "ruleset_source": "o/r", "ruleset_id": 1},
  {"type": "non_fast_forward", "ruleset_source_type": "Repository", "ruleset_source": "o/r", "ruleset_id": 1},
  {"type": "pull_request", "parameters": {"required_approving_review_count": 0, "required_review_thread_resolution": true, "allowed_merge_methods": [ "merge" ]}, "ruleset_id": 1},
  {"type": "required_status_checks", "parameters": {"strict_required_status_checks_policy": true, "required_status_checks": [{"context": "validate / validate", "integration_id": 15368}]}, "ruleset_id": 1}
]`

type boundaryStep struct {
	Name string `yaml:"name"`
	If   string `yaml:"if"`
	Run  string `yaml:"run"`
}

func readBoundaryStep(t *testing.T) boundaryStep {
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
		if step.Name == boundaryStepName {
			return step
		}
	}
	t.Fatalf("clue-validation.yml has no %q step", boundaryStepName)
	return boundaryStep{}
}

type probeResult struct {
	output  string
	request string
}

// runProbe runs the workflow's own step text, under the shell flags GitHub
// uses, with only `curl` replaced — the host is the one thing a test cannot
// reach, and everything the step decides from its answer runs for real.
func runProbe(t *testing.T, rules string, fail bool) probeResult {
	t.Helper()
	tmp := filepath.ToSlash(t.TempDir())
	stub := `curl() {
  printf '%s\n' "$*" > "$RUNNER_TEMP/request"
  if [ -n "$STUB_FAIL" ]; then
    echo "curl: (22) The requested URL returned error: 403" >&2
    return 22
  fi
  printf '%s' "$STUB_RULES"
}
`
	script := filepath.Join(tmp, "probe.sh")
	if err := os.WriteFile(script, []byte(stub+readBoundaryStep(t).Run), 0o644); err != nil {
		t.Fatal(err)
	}
	failFlag := ""
	if fail {
		failFlag = "1"
	}
	cmd := exec.Command(shell(t), "--noprofile", "--norc", "-eo", "pipefail", script)
	cmd.Env = append(os.Environ(),
		"GITHUB_BASE_REF=main",
		"GITHUB_REF_NAME=feature",
		"GITHUB_API_URL=https://api.github.com",
		"GITHUB_REPOSITORY=o/r",
		"GH_TOKEN=secret-token",
		"RUNNER_TEMP="+tmp,
		"STUB_RULES="+rules,
		"STUB_FAIL="+failFlag,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("the probe failed the step, which it must never do: %v\n%s", err, out)
	}
	request, _ := os.ReadFile(filepath.Join(tmp, "request"))
	if strings.Contains(string(out), "secret-token") {
		t.Fatalf("the probe printed the workflow token:\n%s", out)
	}
	return probeResult{output: string(out), request: string(request)}
}

// AC-197 positive: a boundary GitHub reports in full is recorded as observed,
// on the pull request's base branch, with the unseen bypass list named.
func TestAC197_UnitPositive_EnforcedBoundaryIsReportedAsObserved(t *testing.T) {
	got := runProbe(t, enforcedRules, false)
	if !strings.Contains(got.output, "::notice title=Merge boundary observed::") {
		t.Fatalf("an enforced boundary was not reported as observed:\n%s", got.output)
	}
	if strings.Contains(got.output, "::warning") {
		t.Fatalf("an enforced boundary produced a warning:\n%s", got.output)
	}
	if !strings.Contains(got.output, "bypass list is not visible") {
		t.Errorf("the observed notice does not say the bypass list went unseen:\n%s", got.output)
	}
	if !strings.Contains(got.request, "/repos/o/r/rules/branches/main") {
		t.Errorf("the probe did not ask about the pull request's base branch: %q", got.request)
	}
}

// AC-197 negative: losing any single rule withholds the observed notice and
// names that rule, so the notice can never cover a partial boundary.
func TestAC197_UnitNegative_AnyMissingRuleWithholdsTheObservedNotice(t *testing.T) {
	for _, tc := range []struct {
		name, from, to, wantMissing string
	}{
		{"deletion allowed", `"type": "deletion"`, `"type": "creation"`, "branch deletion blocked"},
		{"force pushes allowed", `"type": "non_fast_forward"`, `"type": "update"`, "force pushes blocked"},
		{"no pull request rule", `"type": "pull_request"`, `"type": "required_signatures"`, "pull requests required"},
		{"squash also allowed", `[ "merge" ]`, `[ "merge", "squash" ]`, "merge commits as the only merge method"},
		{"conversations may stay open", `"required_review_thread_resolution": true`, `"required_review_thread_resolution": false`, "review conversations resolved before merge"},
		{"a different check is required", `"context": "validate / validate"`, `"context": "lint"`, "the validate check required"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rules := strings.Replace(enforcedRules, tc.from, tc.to, 1)
			if rules == enforcedRules {
				t.Fatalf("fixture edit %q matched nothing, so this case would pass vacuously", tc.from)
			}
			got := runProbe(t, rules, false)
			if strings.Contains(got.output, "Merge boundary observed") {
				t.Fatalf("a partial boundary was reported as observed:\n%s", got.output)
			}
			if !strings.Contains(got.output, "::warning title=Merge boundary not enforced::") || !strings.Contains(got.output, tc.wantMissing) {
				t.Fatalf("the warning does not name %q:\n%s", tc.wantMissing, got.output)
			}
		})
	}
}

// AC-198 positive: the rules this repository actually had when M-102 was
// added — deletion and force pushes blocked, nothing else — warn in plain
// terms, point to the checklist, and leave the job's result alone.
func TestAC198_UnitPositive_UnenforcedBoundaryWarnsWithoutFailing(t *testing.T) {
	rules := `[{"type":"deletion","ruleset_source_type":"Repository","ruleset_source":"o/r","ruleset_id":2},{"type":"non_fast_forward","ruleset_source_type":"Repository","ruleset_source":"o/r","ruleset_id":2}]`
	got := runProbe(t, rules, false)
	if !strings.Contains(got.output, "::warning title=Merge boundary not enforced::") {
		t.Fatalf("an unenforced boundary produced no warning:\n%s", got.output)
	}
	for _, want := range []string{"pull requests required", "merge commits as the only merge method", "review conversations resolved before merge", "the validate check required", ".github/cliewen-wall.md", "https://cliewen.dev/ci-wall", "Classic branch protection is not visible"} {
		if !strings.Contains(got.output, want) {
			t.Errorf("the warning does not say %q:\n%s", want, got.output)
		}
	}
	for _, present := range []string{"branch deletion blocked", "force pushes blocked"} {
		if strings.Contains(got.output, present) {
			t.Errorf("the warning names %q as missing although GitHub reported it:\n%s", present, got.output)
		}
	}
	empty := runProbe(t, "[]", false)
	if !strings.Contains(empty.output, "::warning title=Merge boundary not enforced::") || !strings.Contains(empty.output, "branch deletion blocked") {
		t.Errorf("a branch with no rules at all was not reported as unenforced:\n%s", empty.output)
	}
}

// AC-198 negative: an unanswered request, or an answer that is not a rule
// list, is unknown — never an observed boundary and never a list of missing
// rules the probe did not actually see.
func TestAC198_UnitNegative_UnanswerableHostIsUnknownNeverObserved(t *testing.T) {
	for _, tc := range []struct {
		name  string
		rules string
		fail  bool
		want  string
	}{
		{"request refused", enforcedRules, true, "403"},
		{"answer is not a rule list", `{"message":"Not Found"}`, false, "did not return a rule list"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := runProbe(t, tc.rules, tc.fail)
			if !strings.Contains(got.output, "::warning title=Merge boundary unknown::") || !strings.Contains(got.output, tc.want) {
				t.Fatalf("an unanswerable host was not reported as unknown:\n%s", got.output)
			}
			if !strings.Contains(got.output, "not a pass") {
				t.Errorf("the unknown warning does not say it is not a pass:\n%s", got.output)
			}
			for _, claim := range []string{"Merge boundary observed", "Merge boundary not enforced"} {
				if strings.Contains(got.output, claim) {
					t.Errorf("an unanswerable host produced the claim %q:\n%s", claim, got.output)
				}
			}
		})
	}
}

// TestSanity_BoundaryProbeRunsAfterFailuresAndTakesNoExpressionInput keeps the
// report present on a red run, where it matters most, and keeps untrusted
// event data out of the script text GitHub would otherwise interpolate.
func TestSanity_BoundaryProbeRunsAfterFailuresAndTakesNoExpressionInput(t *testing.T) {
	step := readBoundaryStep(t)
	if !strings.Contains(step.If, "!cancelled()") {
		t.Errorf("the probe is skipped when an earlier step fails: if=%q", step.If)
	}
	if strings.Contains(step.Run, "${{") {
		t.Error("the probe interpolates an expression into its script; pass values through env instead")
	}
}
