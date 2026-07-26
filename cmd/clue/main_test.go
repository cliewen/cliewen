package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/cliewen/cliewen/internal/corpus"
	"gopkg.in/yaml.v3"
)

func writeFile(t *testing.T, root, rel, content string) {
	t.Helper()
	full := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func validCorpus(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	writeFile(t, root, "docs/README.md", "# Corpus\n\n<!-- clue:index:start -->\n- [goals/](goals/README.md)\n<!-- clue:index:end -->\n")
	writeFile(t, root, "docs/goals/README.md", "# Goals\n\n<!-- clue:index:start -->\n- [G-001](G-001-first.md)\n<!-- clue:index:end -->\n")
	writeFile(t, root, "docs/goals/G-001-first.md", "---\nid: G-001\ntype: goal\nstatus: accepted\nlinks: []\ntitle: First goal\n---\n")
	return root
}

// AC-004: exit 0 on a valid corpus.
func TestAC004_ExitCodeZeroOnValidCorpus(t *testing.T) {
	if code := runValidate([]string{validCorpus(t)}); code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
}

// AC-005: exit 1 on a broken corpus.
func TestAC005_ExitCodeOneOnBrokenCorpus(t *testing.T) {
	root := validCorpus(t)
	writeFile(t, root, "docs/goals/G-001-first.md", "---\nid: G-001\ntype: goal\nlinks: []\ntitle: First goal\n---\n")
	if code := runValidate([]string{root}); code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
}

// AC-018: a valid corpus with inferred artifacts passes and their count
// feeds the OK line.
func TestAC018_InferredArtifactsCountedAndAccepted(t *testing.T) {
	root := validCorpus(t)
	writeFile(t, root, "docs/goals/G-001-first.md", "---\nid: G-001\ntype: goal\nstatus: accepted\nlinks: []\ntitle: First goal\nprovenance: inferred\n---\n")
	if code := runValidate([]string{root}); code != 0 {
		t.Fatalf("inferred provenance is valid; expected exit 0, got %d", code)
	}
	c, _ := corpus.Scan(root)
	if n := inferredCount(c); n != 1 {
		t.Fatalf("expected 1 inferred artifact, got %d", n)
	}
}

// AC-023: a valid corpus passes with agent-enforced constraints and their
// count feeds the OK line as the promotion backlog.
func TestAC023_AgentConstraintCountReported(t *testing.T) {
	root := validCorpus(t)
	writeFile(t, root, "docs/README.md", "# Corpus\n\n<!-- clue:index:start -->\n- [goals/](goals/README.md)\n- [constraints/](constraints/README.md)\n<!-- clue:index:end -->\n")
	writeFile(t, root, "docs/constraints/README.md", "# Constraints\n\n<!-- clue:index:start -->\n- [C-001](C-001-rule.md)\n<!-- clue:index:end -->\n")
	writeFile(t, root, "docs/constraints/C-001-rule.md", "---\nid: C-001\ntype: constraint\nstatus: active\nlinks: []\ntitle: A rule\nsource: AGENTS.md\nenforcement: agent\n---\n")
	if code := runValidate([]string{root}); code != 0 {
		t.Fatalf("an agent-enforced constraint is valid; expected exit 0, got %d", code)
	}
	c, _ := corpus.Scan(root)
	if n := agentConstraintCount(c); n != 1 {
		t.Fatalf("expected 1 agent-enforced constraint, got %d", n)
	}
}

func TestAC023_MachineConstraintNotInBacklogCount(t *testing.T) {
	root := validCorpus(t)
	writeFile(t, root, "docs/README.md", "# Corpus\n\n<!-- clue:index:start -->\n- [goals/](goals/README.md)\n- [constraints/](constraints/README.md)\n<!-- clue:index:end -->\n")
	writeFile(t, root, "docs/constraints/README.md", "# Constraints\n\n<!-- clue:index:start -->\n- [C-001](C-001-rule.md)\n<!-- clue:index:end -->\n")
	writeFile(t, root, "docs/constraints/C-001-rule.md", "---\nid: C-001\ntype: constraint\nstatus: active\nlinks: []\ntitle: A rule\nsource: AGENTS.md\nenforcement: machine\n---\n")
	c, _ := corpus.Scan(root)
	if n := agentConstraintCount(c); n != 0 {
		t.Fatalf("machine-enforced constraints are not backlog; expected 0, got %d", n)
	}
}

// AC-019: version reports the stamp injected at build time.
func TestAC019_VersionCommandReportsStamp(t *testing.T) {
	old := version
	version = "9.9.9"
	defer func() { version = old }()
	var b strings.Builder
	if code := runVersion(&b); code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if !strings.Contains(b.String(), "9.9.9") {
		t.Fatalf("version output %q does not report the stamp", b.String())
	}
}

// AC-019 (negative): an unstamped source build reports "dev", not a
// release number.
func TestAC019_UnstampedBuildReportsDev(t *testing.T) {
	old := version
	version = "dev"
	defer func() { version = old }()
	var b strings.Builder
	runVersion(&b)
	if !strings.Contains(b.String(), "dev") {
		t.Fatalf("unstamped build should report dev, got %q", b.String())
	}
}

// Unit: the build-info fallback stamps `go install module@vX.Y.Z` builds
// and nothing else — checkout builds and commit installs stay unstamped
// (ADR-011: a pseudo-version is a commit, not a release).
func TestUnit_ReleaseFromModuleVersion(t *testing.T) {
	cases := map[string]string{
		"v0.2.0":                               "0.2.0",
		"v1.2.3-rc.1":                          "1.2.3-rc.1",
		"":                                     "",
		"(devel)":                              "",
		"v0.0.0-20260101120000-abcdef123456":   "",
		"v0.1.1-0.20260101120000-abcdef1234ab": "",
		"v0.0.0-20260101120000-abcdef123456+dirty": "",
		"v0.2.0+dirty": "",
	}
	for in, want := range cases {
		if got := releaseFromModuleVersion(in); got != want {
			t.Errorf("releaseFromModuleVersion(%q) = %q, want %q", in, got, want)
		}
	}
}

// AC-033 (wiring): runValidate threads the binary's stamp through
// corpus.Options.Version into the drift rule — a released clue fails
// against lagging skills, a matching release passes.
func TestAC033_RunValidateThreadsVersionIntoDriftRule(t *testing.T) {
	root := validCorpus(t)
	writeFile(t, root, ".agents/skills/clue-delta/skill.md", "---\ncliewen-skill: true\nversion: 0.1.0\n---\n\n# clue-delta\n")
	old := version
	defer func() { version = old }()
	version = "0.2.0"
	if code := runValidate([]string{root}); code != 1 {
		t.Fatalf("clue 0.2.0 against skills at 0.1.0: expected exit 1 (drift), got %d", code)
	}
	version = "0.1.0"
	if code := runValidate([]string{root}); code != 0 {
		t.Fatalf("clue 0.1.0 against skills at 0.1.0: expected exit 0, got %d", code)
	}
}

// Sanity: the release workflow builds versioned cross-platform binaries.
// A repo invariant guarding M-004's release pipeline against regression;
// the operational proof is the first tagged release itself.
func TestSanity_ReleaseWorkflowIsCrossPlatform(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", ".github", "workflows", "release.yml"))
	if err != nil {
		t.Fatalf("release workflow not found: %v", err)
	}
	wf := string(data)
	for _, want := range []string{"main.version", "linux", "darwin", "windows", "arm64", "amd64", "go test", "SHA256SUMS"} {
		if !strings.Contains(wf, want) {
			t.Errorf("release workflow does not mention %q — expected a tested, stamped, checksummed cross-platform build", want)
		}
	}
	// A manual dispatch must not publish a branch-named release: the ref
	// must be guarded to a tag before anything is built.
	if !strings.Contains(wf, "GITHUB_REF_TYPE") {
		t.Error("release workflow does not guard GITHUB_REF_TYPE — a branch dispatch could publish a branch-named release")
	}
	// The publishing action runs with contents: write; a mutable tag pin
	// would let a moved tag ship different code into our releases.
	if pin := regexp.MustCompile(`action-gh-release@[0-9a-f]{40}`); !pin.MatchString(wf) {
		t.Error("release workflow does not pin action-gh-release by commit SHA")
	}
}

// Sanity: the release body is the tag's CHANGELOG.md section (ADR-012) —
// user-facing, reviewed prose. GitHub's auto-generated notes (a PR dump
// with contributor @mentions) must not come back.
func TestSanity_ReleaseNotesComeFromChangelog(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", ".github", "workflows", "release.yml"))
	if err != nil {
		t.Fatalf("release workflow not found: %v", err)
	}
	wf := string(data)
	for _, want := range []string{"CHANGELOG.md", "body_path"} {
		if !strings.Contains(wf, want) {
			t.Errorf("release workflow does not mention %q — the release body must be extracted from the changelog", want)
		}
	}
	if strings.Contains(wf, "generate_release_notes") {
		t.Error("release workflow enables generate_release_notes — release bodies are written for users in the changelog, not auto-generated")
	}
}

// Sanity: the release runs the shipped drift rule stamped as its own tag
// (ADR-011), before anything is built or published. A tag that disagrees
// with the skills' frontmatter stamp must fail there; the alternative is an
// adopter meeting the mismatch as unexplained drift in their repository.
// The workflow's own arguments are executed, not merely read, so the gate
// is judged by what it does: refuse a mismatch naming both versions, admit
// agreement, and keep an undigested workspace out of a release.
func TestSanity_ReleaseRunsTheJudgeStampedAsTheTag(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", ".github", "workflows", "release.yml"))
	if err != nil {
		t.Fatalf("release workflow not found: %v", err)
	}
	wf := string(data)
	gate := regexp.MustCompile(`run: go run -ldflags "-X main\.version=\$\{VERSION\}" \./cmd/clue ([^\n]+)`)
	m := gate.FindStringSubmatchIndex(wf)
	if m == nil {
		t.Fatal("release workflow never runs clue stamped with the tag's version — nothing compares the tag to the shipped skill stamp")
	}
	// Scoped to this step: an empty ${VERSION} is exempt from the drift
	// rule, so a stamp wired from nothing fails open and ships the
	// mismatch. Another step's env does not vouch for this one.
	step := wf[:m[0]]
	if i := strings.LastIndex(step, "- name:"); i >= 0 {
		step = step[i:]
	}
	if !strings.Contains(step, "VERSION: ${{ steps.version.outputs.v }}") {
		t.Error("the stamped step does not take VERSION from the tag — an empty stamp is exempt from the drift rule, so the gate would pass on any corpus")
	}
	for _, later := range []string{"Build cross-platform binaries", "action-gh-release"} {
		if i := strings.Index(wf, later); i >= 0 && i < m[0] {
			t.Errorf("the stamped step runs after %q — a mismatched tag must fail before any artifact exists", later)
		}
	}

	args := strings.Fields(wf[m[2]:m[3]])
	if len(args) == 0 || args[0] != "validate" {
		t.Fatalf("the stamped step runs %q, not validate — the drift rule is never reached", args)
	}

	// The arguments the release actually passes, against a corpus whose
	// skills lag the stamp: exit 1, naming both versions.
	root := validCorpus(t)
	writeFile(t, root, ".agents/skills/clue-delta/skill.md", "---\ncliewen-skill: true\nversion: 0.1.0\n---\n\n# clue-delta\n")
	// runValidate receives the arguments after the subcommand, as main does.
	runArgs := append(append([]string{}, args[1:]...), root)
	old := version
	defer func() { version = old }()

	version = "0.2.0"
	code, out := runValidateCapturingStdout(t, runArgs)
	if code != 1 {
		t.Errorf("tag 0.2.0 against skills at 0.1.0: expected exit 1, got %d — the release would ship a self-inconsistent pair", code)
	}
	for _, want := range []string{"0.1.0", "0.2.0"} {
		if !strings.Contains(out, want) {
			t.Errorf("the refusal does not name %s: %q — the maintainer cannot see which half to bump", want, out)
		}
	}

	version = "0.1.0"
	if code, _ = runValidateCapturingStdout(t, runArgs); code != 0 {
		t.Errorf("tag 0.1.0 against skills at 0.1.0: expected exit 0, got %d — releases that agree must not be blocked", code)
	}

	// --forbid-changes is load-bearing in this step: an undigested
	// workspace must not reach a release either.
	writeFile(t, root, "changes/CH-009-x/proposal.md", "---\nid: CH-009\ntype: change\nstatus: open\nlinks: []\ntitle: X\n---\n")
	if code, _ = runValidateCapturingStdout(t, runArgs); code != 1 {
		t.Errorf("undigested workspace at a matching tag: expected exit 1, got %d — the release step dropped --forbid-changes", code)
	}
}

// runValidateCapturingStdout runs the validate command and returns its exit
// code with the issue lines it printed, which is where the drift rule names
// the two versions.
func runValidateCapturingStdout(t *testing.T, args []string) (int, string) {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	defer func() { _ = r.Close() }()
	old := os.Stdout
	defer func() { os.Stdout = old }()
	os.Stdout = w
	code := runValidate(args)
	if err := w.Close(); err != nil {
		t.Fatalf("close pipe: %v", err)
	}
	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read pipe: %v", err)
	}
	return code, string(out)
}

type communityIssueForm struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Body        []struct {
		Type        string          `yaml:"type"`
		ID          string          `yaml:"id"`
		Validations map[string]bool `yaml:"validations"`
	} `yaml:"body"`
}

// Sanity: the public community front door remains present, its GitHub
// configuration parses, and private reports cannot silently lose their
// routes while the visible templates continue to look complete.
func TestSanity_CommunityFrontDoorIsWellFormed(t *testing.T) {
	root := filepath.Join("..", "..")
	read := func(rel string) string {
		t.Helper()
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
		if err != nil {
			t.Fatalf("%s not found: %v", rel, err)
		}
		return string(data)
	}
	const (
		conductMailto  = "mailto:flemming&#46;n&#46;larsen&#43;cliewen-conduct&#64;gmail&#46;com"
		securityMailto = "mailto:flemming&#46;n&#46;larsen&#43;cliewen-security&#64;gmail&#46;com"
	)

	for rel, wants := range map[string][]string{
		"CONTRIBUTING.md": {"CODE_OF_CONDUCT.md", "SECURITY.md", conductMailto, "human maintainer", "plan-less", "plain change", "For a plain change", "For a Cliewen change", "automatic agentic review pass"},
		"CODE_OF_CONDUCT.md": {
			"Contributor Covenant 3.0 Code of Conduct",
			"## Encouraged Behaviors",
			"## Restricted Behaviors",
			"## Addressing and Repairing Harm",
			conductMailto,
			"[Cliewen Conduct]",
		},
		"SECURITY.md": {
			securityMailto,
			"[Cliewen Security]",
			"7 calendar days",
			"14 calendar days",
			"Do not open a public issue",
		},
		".github/pull_request_template.md": {
			"plain change",
			"Change ID",
			"Change tier",
			"Plan item served",
			"clue validate --forbid-changes",
			"human review and merge",
		},
	} {
		content := read(rel)
		for _, want := range wants {
			if !strings.Contains(content, want) {
				t.Errorf("%s does not contain required community-front-door text %q", rel, want)
			}
		}
	}
	prTemplate := read(".github/pull_request_template.md")
	cliewenSection := strings.Index(prTemplate, "## Cliewen proposal")
	reviewEvidence := strings.Index(prTemplate, "Agentic review mode and reviewed commit")
	if cliewenSection < 0 || reviewEvidence < cliewenSection {
		t.Error(".github/pull_request_template.md must keep agentic-review evidence inside the removable Cliewen-only section")
	}
	plainEmail := regexp.MustCompile(`[[:alnum:]._%+-]+@[[:alnum:].-]+\.[[:alpha:]]{2,}`)
	for _, rel := range []string{"CONTRIBUTING.md", "CODE_OF_CONDUCT.md", "SECURITY.md", "docs/decisions/PDR-010-community-participation.md"} {
		if match := plainEmail.FindString(read(rel)); match != "" {
			t.Errorf("%s exposes plain email address %q instead of an encoded reporting link", rel, match)
		}
	}

	for rel, requiredIDs := range map[string][]string{
		".github/ISSUE_TEMPLATE/bug.yml":  {"affected_version", "install_route", "environment", "description", "steps", "expected", "actual", "checks"},
		".github/ISSUE_TEMPLATE/goal.yml": {"audience", "need", "success", "checks"},
	} {
		var form communityIssueForm
		if err := yaml.Unmarshal([]byte(read(rel)), &form); err != nil {
			t.Fatalf("%s is not valid YAML: %v", rel, err)
		}
		if form.Name == "" || form.Description == "" || len(form.Body) == 0 {
			t.Errorf("%s must have a name, description, and body", rel)
		}
		ids := make(map[string]bool)
		for _, item := range form.Body {
			if item.Type == "" {
				t.Errorf("%s contains a body item without a type", rel)
			}
			if item.Type == "markdown" {
				continue
			}
			if item.ID == "" {
				t.Errorf("%s contains a %s item without an id", rel, item.Type)
				continue
			}
			if _, exists := ids[item.ID]; exists {
				t.Errorf("%s contains duplicate id %q", rel, item.ID)
			}
			ids[item.ID] = item.Validations["required"]
		}
		for _, id := range requiredIDs {
			if !ids[id] {
				t.Errorf("%s field %q must exist and be required", rel, id)
			}
		}
	}

	var config struct {
		BlankIssuesEnabled *bool `yaml:"blank_issues_enabled"`
		ContactLinks       []struct {
			URL string `yaml:"url"`
		} `yaml:"contact_links"`
	}
	if err := yaml.Unmarshal([]byte(read(".github/ISSUE_TEMPLATE/config.yml")), &config); err != nil {
		t.Fatalf("issue-template config is not valid YAML: %v", err)
	}
	if config.BlankIssuesEnabled == nil || *config.BlankIssuesEnabled {
		t.Error("issue-template config must explicitly disable blank issues")
	}
	hasSecurityPolicy := false
	for _, link := range config.ContactLinks {
		if strings.HasSuffix(link.URL, "/security/policy") {
			hasSecurityPolicy = true
		}
	}
	if !hasSecurityPolicy {
		t.Error("issue-template config must route vulnerability reports to the private security policy")
	}
}

func TestSanity_ReviewFixConstraintOrdersFinalCandidateBeforeReview(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", "docs", "constraints", "C-012-agents-never-merge-own-changes.md"))
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	fixes := strings.Index(content, "Review fixes to an unaccepted change")
	if fixes < 0 {
		t.Fatal("C-012 does not define the review-fix handoff")
	}
	handoff := content[fixes:]
	commitCandidate := strings.Index(handoff, "commit")
	verifyCandidate := strings.Index(handoff, "local verification")
	reviewCandidate := strings.Index(handoff, "agentic review")
	pushCandidate := strings.Index(handoff, "push")
	if commitCandidate < 0 || verifyCandidate <= commitCandidate || reviewCandidate <= verifyCandidate || pushCandidate <= reviewCandidate {
		t.Error("C-012 must commit a repaired candidate, verify and review that commit, then push the reviewed commit")
	}
}

func TestSanity_PRBoundaryExplainsAuthorizationAndCIEnforcement(t *testing.T) {
	root := filepath.Join("..", "..")
	for rel, wants := range map[string][]string{
		"docs/decisions/PDR-007-review-boundary.md": {
			"authorization and protected-integration boundary",
			"not a requirement that a solo developer repeat a code review",
			"The PR alone is insufficient",
			"required check",
			"branch protection",
		},
		"guide/change-loop.md": {
			"authorization and protected-integration gate",
			"not a demand for duplicate human code review",
			"a PR alone does not enforce anything",
			"required status check",
			"branch protection",
		},
		"guide/what-is-cliewen.md": {
			"pull request is the authorization boundary",
			"does not require repeating a code review",
			"required check and branch protection",
		},
		"docs/decisions/log.md": {
			"PDR-007 supersedes CH-023's “human review gate” interpretation",
			"mandatory authorization and protected-integration boundary",
		},
		"docs/constraints/README.md": {
			"C-012 — Changes are reviewed locally, root at main, and remain human-merged",
		},
	} {
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
		if err != nil {
			t.Fatalf("%s not found: %v", rel, err)
		}
		content := string(data)
		for _, want := range wants {
			if !strings.Contains(content, want) {
				t.Errorf("%s does not explain PR safeguard %q", rel, want)
			}
		}
	}
}

func TestSanity_AgenticFindingsRequireOperativeViolations(t *testing.T) {
	root := filepath.Join("..", "..")
	for rel, wants := range map[string][]string{
		"docs/decisions/PDR-012-agentic-review-before-publication.md": {
			"A finding is grounded in an operative requirement",
			"human-controlled merge is mandatory but duplicate human code review is not",
			"a release cut renames `[Unreleased]` to the versioned section",
			"lifecycle-correct state are not actionable defects by themselves",
		},
		"guide/change-loop.md": {
			"Every finding identifies the operative requirement or declared intent that is violated and its concrete consequence",
			"Human-controlled merge does not imply duplicate human code review",
			"a release cut uses its versioned changelog section instead of `[Unreleased]`",
		},
	} {
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
		if err != nil {
			t.Fatalf("%s not found: %v", rel, err)
		}
		content := string(data)
		for _, want := range wants {
			if !strings.Contains(content, want) {
				t.Errorf("%s does not contain grounded-finding rule %q", rel, want)
			}
		}
	}
}

// Sanity: guide-only editorial changes take the focused CI path, while
// the classifier and workflow retain a fail-closed full path.
func TestSanity_CIHasFailClosedFocusedGuideScope(t *testing.T) {
	root := filepath.Join("..", "..")
	for rel, wants := range map[string][]string{
		".github/scripts/ci-scope.mjs": {
			"files.length > 0",
			`file.startsWith("guide/")`,
			`file.endsWith(".md")`,
			"full: !focusedGuide",
		},
		".github/workflows/ci.yml": {
			"Classify changed surface",
			"full=true",
			"guide=true",
			"steps.scope.outputs.full == 'true'",
			"npm run guide:build",
			"git diff --check",
		},
	} {
		content, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
		if err != nil {
			t.Fatalf("%s not found: %v", rel, err)
		}
		for _, want := range wants {
			if !strings.Contains(string(content), want) {
				t.Errorf("%s does not contain fail-closed CI text %q", rel, want)
			}
		}
	}
}

// docID matches a Cliewen corpus doc-ID reference: ADR-011, G-002, CAP-004,
// AC-020, P-002, M-004, CH-007, AN-002, QS-001, and so on. Digits are what
// make it a reference — placeholder forms (CH-xxx, @AC-xxx) don't match.
var docID = regexp.MustCompile(`\b(?:ADR|CAP|AC|G|P|M|CH|AN|QS)-\d+\b`)

// Sanity: cmd/clue is the shipped CLI — the one package under this module
// actually exported to a user, unlike internal/corpus which Go itself
// keeps unimportable outside the module. A corpus doc-ID reference leaking
// into a string literal here means a user sees "(ADR-011)" in --help or
// command output with no way to know what that is (caught in PR #6 review:
// the usage string named ADR-011 in a line explaining `clue validate`).
// AST-based so this only inspects actual string literals, not source
// comments — comments citing ADR/CAP/AC/G/P/M/CH IDs for future readers of
// the code remain fine.
func TestSanity_NoDocIDInUserFacingStrings(t *testing.T) {
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	fset := token.NewFileSet()
	for _, f := range files {
		if strings.HasSuffix(f, "_test.go") {
			continue
		}
		node, perr := parser.ParseFile(fset, f, nil, 0)
		if perr != nil {
			t.Fatalf("parsing %s: %v", f, perr)
		}
		ast.Inspect(node, func(n ast.Node) bool {
			lit, ok := n.(*ast.BasicLit)
			if !ok || lit.Kind != token.STRING {
				return true
			}
			if m := docID.FindString(lit.Value); m != "" {
				t.Errorf("%s: string literal mentions %q — cmd/clue is user-facing, a corpus doc-ID means nothing to a CLI user", fset.Position(lit.Pos()), m)
			}
			return true
		})
	}
}

// Sanity: the skills under .agents/skills ship verbatim to adopting repos,
// where this repo's corpus doc-IDs resolve to nothing — or to that repo's
// own unrelated documents. A skill states each rule's content in its own
// text; the deciding document stays in this repo's corpus. Placeholder
// forms (CH-xxx, @AC-xxx) stay fine — digits are what make it a reference.
func TestSanity_SkillsCarryNoDocIDs(t *testing.T) {
	root := filepath.Join("..", "..", ".agents", "skills")
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}
		data, rerr := os.ReadFile(path)
		if rerr != nil {
			return rerr
		}
		for i, line := range strings.Split(string(data), "\n") {
			if ids := docID.FindAllString(line, -1); ids != nil {
				t.Errorf("%s:%d references corpus doc-ID(s) %v — skills are exported verbatim; state the rule, don't cite the document", path, i+1, ids)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

// AC-008: the --forbid-changes gate flips the exit code, nothing else.
func TestAC008_ForbidChangesFlagExitCodes(t *testing.T) {
	root := validCorpus(t)
	writeFile(t, root, "changes/CH-009-x/proposal.md", "---\nid: CH-009\ntype: change\nstatus: open\nlinks: []\ntitle: X\n---\n")
	if code := runValidate([]string{root}); code != 0 {
		t.Fatalf("without the gate: expected exit 0, got %d", code)
	}
	if code := runValidate([]string{"--forbid-changes", root}); code != 1 {
		t.Fatalf("with the gate: expected exit 1, got %d", code)
	}
}
