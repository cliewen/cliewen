package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

type ciWorkflow struct {
	Jobs map[string]ciJob `yaml:"jobs"`
}

type ciJob struct {
	Name   string            `yaml:"name"`
	Uses   string            `yaml:"uses"`
	RunsOn string            `yaml:"runs-on"`
	With   map[string]string `yaml:"with"`
	Steps  []ciStep          `yaml:"steps"`
}

type ciStep struct {
	Name string `yaml:"name"`
	Uses string `yaml:"uses"`
	Run  string `yaml:"run"`
}

var reusableWorkflowRef = regexp.MustCompile(`^cliewen/cliewen/\.github/workflows/clue-validation\.yml@[0-9a-f]{40}$`)

func readCIWorkflow(t *testing.T, rel string) (string, ciWorkflow) {
	t.Helper()
	data, err := os.ReadFile(filepath.FromSlash(rel))
	if err != nil {
		t.Fatal(err)
	}
	var workflow ciWorkflow
	if err := yaml.Unmarshal(data, &workflow); err != nil {
		t.Fatalf("%s is not valid workflow YAML: %v", rel, err)
	}
	return string(data), workflow
}

func assertThinCaller(t *testing.T, raw string, workflow ciWorkflow) error {
	t.Helper()
	job, ok := workflow.Jobs["validate"]
	if !ok || len(workflow.Jobs) != 1 {
		return fmt.Errorf("caller must contain exactly one validate job")
	}
	if job.Name != "validate" {
		return fmt.Errorf("caller job name is %q, want validate", job.Name)
	}
	if !reusableWorkflowRef.MatchString(job.Uses) {
		return fmt.Errorf("caller uses %q, want the reusable workflow at a full commit SHA", job.Uses)
	}
	if len(job.Steps) != 0 {
		return fmt.Errorf("caller embeds %d steps instead of delegating the validation unit", len(job.Steps))
	}
	wantInputs := []string{"clue-install-directory", "clue-source", "clue-version", "runner"}
	gotInputs := make([]string, 0, len(job.With))
	for input := range job.With {
		gotInputs = append(gotInputs, input)
	}
	slices.Sort(gotInputs)
	if !slices.Equal(gotInputs, wantInputs) {
		return fmt.Errorf("caller inputs are %v, want %v", gotInputs, wantInputs)
	}
	if !regexp.MustCompile(`(?m)^\s*uses:\s*actions/`).MatchString(raw) && strings.Contains(raw, "actions/") {
		return fmt.Errorf("caller mentions an action without owning an action step")
	}
	return nil
}

func assertReusableWorkflow(t *testing.T, raw string, workflow ciWorkflow) error {
	t.Helper()
	if !strings.Contains(raw, "workflow_call:") {
		return fmt.Errorf("validation unit is not callable through workflow_call")
	}
	job, ok := workflow.Jobs["validate"]
	if !ok || job.Name != "validate" {
		return fmt.Errorf("validation unit must expose a stable validate job")
	}
	if !strings.Contains(job.RunsOn, "fromJSON(inputs.runner)") {
		return fmt.Errorf("validation unit does not expose the caller runner choice")
	}
	for _, input := range []string{"runner", "clue-version", "clue-source", "clue-install-directory"} {
		if !strings.Contains(raw, input+":") {
			return fmt.Errorf("validation unit does not declare input %q", input)
		}
	}
	actionRef := regexp.MustCompile(`(?m)^\s*-?\s*uses:\s*actions/[^@\s]+@([^\s#]+)`)
	actionRefs := actionRef.FindAllStringSubmatch(raw, -1)
	if len(actionRefs) == 0 {
		return fmt.Errorf("validation unit has no immutable actions reference")
	}
	for _, ref := range actionRefs {
		if !workflowRefRe.MatchString(ref[1]) {
			return fmt.Errorf("validation unit contains a mutable actions reference %q", ref[1])
		}
	}
	for _, fragment := range []string{
		"Detect Cliewen change",
		"Check whether the wall is armed",
		"Require a completed acceptance brief",
		"clue validate --forbid-changes",
		"grep -Eqv '\\.md$'",
		"steps.scope.outputs.cliewen == 'true'",
	} {
		if !strings.Contains(raw, fragment) {
			return fmt.Errorf("validation unit is missing %q", fragment)
		}
	}
	if strings.Contains(raw, "/usr/local/bin") {
		return fmt.Errorf("validation unit hardcodes a root-only install path")
	}
	return nil
}

func callerChoices(job ciJob) map[string]string {
	choices := make(map[string]string, len(job.With))
	for key, value := range job.With {
		choices[key] = value
	}
	return choices
}

func equalChoices(left, right map[string]string) bool {
	if len(left) != len(right) {
		return false
	}
	for key, value := range left {
		if right[key] != value {
			return false
		}
	}
	return true
}

// AC-062 positive: init delegates the complete validation unit and leaves
// only the demonstrated adopter inputs in its caller.
func TestAC062_UnitPositive_InitEmitsThinImmutableCaller(t *testing.T) {
	root, _ := runInto(t)
	callerPath := filepath.Join(root, ".github", "workflows", "clue.yml")
	raw, workflow := readCIWorkflow(t, callerPath)
	if err := assertThinCaller(t, raw, workflow); err != nil {
		t.Fatal(err)
	}
	job := workflow.Jobs["validate"]
	version, err := PairVersion()
	if err != nil {
		t.Fatal(err)
	}
	if job.With["clue-version"] != version {
		t.Fatalf("caller pins clue version %q, want %q", job.With["clue-version"], version)
	}
	if job.With["clue-source"] != "vendored" || job.With["clue-install-directory"] != "" {
		t.Fatalf("fresh caller changed default acquisition choices: %v", job.With)
	}
}

// AC-062 negative: a copied action or validation step is not a thin caller.
func TestAC062_UnitNegative_EmbeddedWallIsRejected(t *testing.T) {
	bad := `name: clue
on: push
jobs:
  validate:
    name: validate
    uses: cliewen/cliewen/.github/workflows/clue-validation.yml@main
    with:
      runner: '["ubuntu-latest"]'
      clue-version: 0.10.0
      clue-source: vendored
      clue-install-directory: ''
    steps:
      - name: copied validation
        run: clue validate --forbid-changes
`
	var workflow ciWorkflow
	if err := yaml.Unmarshal([]byte(bad), &workflow); err != nil {
		t.Fatal(err)
	}
	if err := assertThinCaller(t, bad, workflow); err == nil {
		t.Fatal("embedded or mutable caller unexpectedly satisfied the thin-caller contract")
	}
}

// AC-063 positive: the fixture changes only the upstream reference, while
// the enterprise fixture demonstrates the two caller-owned policy choices.
func TestAC063_UnitPositive_UpstreamUpdatePreservesCallerChoices(t *testing.T) {
	oldRaw, oldWorkflow := readCIWorkflow(t, filepath.Join("testdata", "ci-boundary", "default-old.yml"))
	newRaw, newWorkflow := readCIWorkflow(t, filepath.Join("testdata", "ci-boundary", "default-new.yml"))
	if err := assertThinCaller(t, oldRaw, oldWorkflow); err != nil {
		t.Fatal(err)
	}
	if err := assertThinCaller(t, newRaw, newWorkflow); err != nil {
		t.Fatal(err)
	}
	oldJob := oldWorkflow.Jobs["validate"]
	newJob := newWorkflow.Jobs["validate"]
	if oldJob.Uses == newJob.Uses {
		t.Fatal("fixture versions do not change the upstream reference")
	}
	if !equalChoices(callerChoices(oldJob), callerChoices(newJob)) {
		t.Fatalf("upstream update changed caller choices: old=%v new=%v", oldJob.With, newJob.With)
	}

	enterpriseRaw, enterpriseWorkflow := readCIWorkflow(t, filepath.Join("testdata", "ci-boundary", "enterprise.yml"))
	if err := assertThinCaller(t, enterpriseRaw, enterpriseWorkflow); err != nil {
		t.Fatal(err)
	}
	enterprise := enterpriseWorkflow.Jobs["validate"].With
	if enterprise["runner"] == oldJob.With["runner"] || enterprise["clue-source"] == oldJob.With["clue-source"] || enterprise["clue-install-directory"] == oldJob.With["clue-install-directory"] {
		t.Fatalf("enterprise fixture does not exercise caller-owned choices: %v", enterprise)
	}

	raw, workflow := readCIWorkflow(t, filepath.Join("..", "..", ".github", "workflows", "clue-validation.yml"))
	if err := assertReusableWorkflow(t, raw, workflow); err != nil {
		t.Fatal(err)
	}
}

// AC-063 negative: changing a local policy input is not an upstream-only
// update and must remain visible in the caller diff.
func TestAC063_UnitNegative_LocalChoiceChangeIsDetected(t *testing.T) {
	oldRaw, oldWorkflow := readCIWorkflow(t, filepath.Join("testdata", "ci-boundary", "default-old.yml"))
	if err := assertThinCaller(t, oldRaw, oldWorkflow); err != nil {
		t.Fatal(err)
	}
	changed := strings.Replace(oldRaw, "clue-source: vendored", "clue-source: release", 1)
	var changedWorkflow ciWorkflow
	if err := yaml.Unmarshal([]byte(changed), &changedWorkflow); err != nil {
		t.Fatal(err)
	}
	if equalChoices(callerChoices(oldWorkflow.Jobs["validate"]), callerChoices(changedWorkflow.Jobs["validate"])) {
		t.Fatal("fixture failed to detect the changed local clue-source choice")
	}
}

func TestSanity_ReusableWorkflowHasOneImmutableActionReference(t *testing.T) {
	raw, workflow := readCIWorkflow(t, filepath.Join("..", "..", ".github", "workflows", "clue-validation.yml"))
	if err := assertReusableWorkflow(t, raw, workflow); err != nil {
		t.Fatal(err)
	}
}
