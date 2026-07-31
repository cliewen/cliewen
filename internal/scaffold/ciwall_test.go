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
	// GitHub's `on:` key survives as a string in YAML 1.2, which is what
	// yaml.v3 implements; it is not folded into the boolean `true`.
	On   map[string]any   `yaml:"on"`
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

// Both reference forms ADR-038 admits are immutable: the emitting source
// commit, and the protected release tag a build without VCS metadata falls
// back to. A branch or floating ref is neither, and is what this rejects.
var reusableWorkflowRef = regexp.MustCompile(`^cliewen/cliewen/\.github/workflows/clue-validation\.yml@(?:[0-9a-f]{40}|v[0-9]+\.[0-9]+\.[0-9]+)$`)

// A thin caller delegates every step, so a run step is copied upstream logic
// by definition, and the delegation itself is the only `uses:` it may carry —
// an action pin in the caller is the fork ADR-038 removes, one level down.
var (
	callerRunStep = regexp.MustCompile(`(?m)^\s*-?\s*run:\s*\S`)
	callerUses    = regexp.MustCompile(`(?m)^\s*-?\s*uses:\s*(\S+)`)
)

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
	// The caller owns the event boundary: without both triggers the stable
	// check exists but stops gating pull requests or stops covering main,
	// and nothing about the delegated job would reveal it.
	for _, trigger := range []string{"pull_request", "push"} {
		if _, ok := workflow.On[trigger]; !ok {
			return fmt.Errorf("caller does not run on %s", trigger)
		}
	}
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
	if callerRunStep.MatchString(raw) {
		return fmt.Errorf("caller carries a run step instead of delegating the validation unit")
	}
	uses := callerUses.FindAllStringSubmatch(raw, -1)
	if len(uses) != 1 || uses[0][1] != job.Uses {
		return fmt.Errorf("caller carries %d uses references, want only the reusable workflow", len(uses))
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

// callerFixture renders a thin caller, so each negative case below differs
// from a passing caller in exactly one way and names which rule caught it.
// A single fixture violating several rules at once would still pass while
// any one of those rules silently rotted.
func callerFixture(ref, triggers, extra string) string {
	return fmt.Sprintf(`name: clue
%s
jobs:
  validate:
    name: validate
    uses: cliewen/cliewen/.github/workflows/clue-validation.yml@%s
    with:
      runner: '["ubuntu-latest"]'
      clue-version: 0.10.0
      clue-source: vendored
      clue-install-directory: ''
%s`, triggers, ref, extra)
}

const defaultTriggers = `on:
  pull_request:
  push:
    branches: [main]`

const immutableRef = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

// AC-062 negative: each way of forking the wall back into the caller is
// rejected, and rejected for its own stated reason.
func TestAC062_UnitNegative_EmbeddedWallIsRejected(t *testing.T) {
	// The generator itself must produce a caller that passes, or every case
	// below would be vacuous — a fixture typo would "reject" for free.
	control := callerFixture(immutableRef, defaultTriggers, "")
	var controlWorkflow ciWorkflow
	if err := yaml.Unmarshal([]byte(control), &controlWorkflow); err != nil {
		t.Fatal(err)
	}
	if err := assertThinCaller(t, control, controlWorkflow); err != nil {
		t.Fatalf("the negative fixtures' own baseline is not a thin caller: %v", err)
	}

	for _, tc := range []struct {
		name     string
		workflow string
		wantErr  string
	}{
		{
			name:     "mutable branch reference",
			workflow: callerFixture("main", defaultTriggers, ""),
			wantErr:  "want the reusable workflow at a full commit SHA",
		},
		{
			name:     "copied validation step",
			workflow: callerFixture(immutableRef, defaultTriggers, "    steps:\n      - name: copied validation\n        run: clue validate --forbid-changes\n"),
			wantErr:  "caller embeds 1 steps",
		},
		{
			name:     "a second job carrying its own action pin",
			workflow: callerFixture(immutableRef, defaultTriggers, "  extra:\n    runs-on: ubuntu-latest\n    steps:\n      - uses: actions/checkout@v4\n"),
			wantErr:  "exactly one validate job",
		},
		{
			name:     "pull requests no longer gated",
			workflow: callerFixture(immutableRef, "on:\n  push:\n    branches: [main]", ""),
			wantErr:  "does not run on pull_request",
		},
		{
			name:     "an upstream-owned knob became a caller input",
			workflow: callerFixture(immutableRef, defaultTriggers, "      checkout-action: actions/checkout@v4\n"),
			wantErr:  "caller inputs are",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var workflow ciWorkflow
			if err := yaml.Unmarshal([]byte(tc.workflow), &workflow); err != nil {
				t.Fatal(err)
			}
			err := assertThinCaller(t, tc.workflow, workflow)
			if err == nil {
				t.Fatal("forked caller unexpectedly satisfied the thin-caller contract")
			}
			if !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("rejected for the wrong reason: got %q, want it to mention %q", err, tc.wantErr)
			}
		})
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

// Sanity: a `workflow_call`-only file is never run by the repository that
// owns it, so without an in-repository caller the validation unit every
// adopter depends on could ship having executed nowhere (ADR-038). This
// holds the caller's existence, its delegation, and the name boundary that
// keeps it out of this repository's own required check.
func TestSanity_ValidationUnitIsExercisedInThisRepository(t *testing.T) {
	raw, workflow := readCIWorkflow(t, filepath.Join("..", "..", ".github", "workflows", "clue-validation-smoke.yml"))
	if _, named := workflow.Jobs["validate"]; named {
		t.Error("the smoke caller declares a `validate` job, which collides with this repository's required check")
	}
	var delegating ciJob
	for _, job := range workflow.Jobs {
		if strings.Contains(job.Uses, "clue-validation.yml") {
			delegating = job
		}
	}
	if delegating.Uses == "" {
		t.Fatal("no job in the smoke workflow calls clue-validation.yml")
	}
	for _, input := range []string{"clue-version", "clue-source", "runner", "clue-install-directory"} {
		if _, ok := delegating.With[input]; !ok {
			t.Errorf("the smoke caller does not exercise the %q input", input)
		}
	}
	// The dispatch path is what reaches acquisition and validation; the
	// automatic path can only ever exercise the unarmed branch.
	if !strings.Contains(raw, "workflow_dispatch:") || !strings.Contains(raw, "release") {
		t.Error("the smoke caller cannot be dispatched against a released binary, so acquisition and validation stay unexercised")
	}
}

// Unit: the emitted reference falls back to the release tag unless the tree
// that built the binary is positively this project. A module cache entry or
// an extracted archive can sit inside an unrelated git repository, and that
// repository's HEAD is a reference no adopter could resolve.
func TestUnit_WorkflowReferenceRejectsAForeignCheckout(t *testing.T) {
	if isCliewenCheckout(t.TempDir()) {
		t.Fatal("an empty directory was accepted as this project's checkout")
	}

	// A tree with the workflow but another module's go.mod is the shape a
	// vendored or re-published copy would take.
	foreign := t.TempDir()
	if err := os.MkdirAll(filepath.Join(foreign, ".github", "workflows"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(foreign, ".github", "workflows", "clue-validation.yml"), []byte("name: copy\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(foreign, "go.mod"), []byte("module example.com/other\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if isCliewenCheckout(foreign) {
		t.Fatal("a foreign module carrying a copy of the workflow was accepted as this project's checkout")
	}

	if !isCliewenCheckout(filepath.Join("..", "..")) {
		t.Fatal("this project's own checkout was not recognized")
	}
}

// Unit: whichever branch workflowReference takes, what it emits is one of the
// two immutable forms ADR-038 admits — never a branch, and never empty.
func TestUnit_WorkflowReferenceIsAlwaysImmutable(t *testing.T) {
	version, err := PairVersion()
	if err != nil {
		t.Fatal(err)
	}
	ref := workflowReference(version)
	if !workflowRefRe.MatchString(ref) && ref != "v"+version {
		t.Fatalf("emitted reference %q is neither a full commit nor the release tag", ref)
	}
}
