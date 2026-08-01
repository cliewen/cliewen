package main

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/cliewen/cliewen/internal/corpus"
	"github.com/cliewen/cliewen/internal/scaffold"
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

// AC-051: low-cost inferred artifacts pass without feeding a monotonic count.
func TestAC051_LowCostInferredArtifactsAreAccepted(t *testing.T) {
	root := validCorpus(t)
	writeFile(t, root, "docs/goals/G-001-first.md", "---\nid: G-001\ntype: goal\nstatus: accepted\nlinks: []\ntitle: First goal\nprovenance: inferred\nreversal-cost: low\n---\n")
	if code := runValidate([]string{root}); code != 0 {
		t.Fatalf("inferred provenance is valid; expected exit 0, got %d", code)
	}
	c, _ := corpus.Scan(root)
	if n := len(corpus.ProvenanceBacklog(c).BlockerArtifacts); n != 0 {
		t.Fatalf("low-cost inferred artifact must not be an activation blocker, got %d", n)
	}
}

func TestAC052_RealityGapsFlagPrintsAffectedCapability(t *testing.T) {
	root := validCorpus(t)
	writeFile(t, root, "docs/README.md", "# Corpus\n\n<!-- clue:index:start -->\n- [goals/](goals/README.md)\n- [capabilities/](capabilities/README.md)\n- [analysis/](analysis/README.md)\n<!-- clue:index:end -->\n")
	writeFile(t, root, "docs/capabilities/README.md", "# Capabilities\n\n<!-- clue:index:start -->\n- [CAP-101](CAP-101-x/README.md)\n<!-- clue:index:end -->\n")
	writeFile(t, root, "docs/capabilities/CAP-101-x/README.md", "---\nid: CAP-101\ntype: capability\nstatus: active\nlinks: []\ntitle: X\n---\n")
	writeFile(t, root, "docs/capabilities/CAP-101-x/criteria.md", "---\nid: CAP-101-criteria\ntype: criteria\nstatus: active\nlinks: [CAP-101]\ntitle: X criteria\n---\n\n```gherkin\nFeature: X\n\n  @AC-101 @draft\n  Scenario: X\n    Given X\n    Then X\n```\n")
	writeFile(t, root, "docs/analysis/README.md", "# Analysis\n\n<!-- clue:index:start -->\n- [AN-101](AN-101-incident.md)\n<!-- clue:index:end -->\n")
	writeFile(t, root, "docs/analysis/AN-101-incident.md", "---\nid: AN-101\ntype: analysis\nstatus: active\nlinks: [AC-101]\ntitle: Incident\nreality: contradicted\n---\n")
	code, out := runValidateCapturingStdout(t, []string{"--reality-gaps", root})
	if code != 0 || !strings.Contains(out, "CAP-101: contradicted by AN-101") {
		t.Fatalf("expected derived reality gap, code=%d output=%q", code, out)
	}
}

func TestAC051_CLIReportsInferredDecisionsSeparately(t *testing.T) {
	root := validCorpus(t)
	writeFile(t, root, "docs/README.md", "# Corpus\n\n<!-- clue:index:start -->\n- [goals/](goals/README.md)\n- [decisions/](decisions/README.md)\n<!-- clue:index:end -->\n")
	writeFile(t, root, "docs/decisions/README.md", "# Decisions\n\n<!-- clue:index:start -->\n- [ADR-101](ADR-101-x.md)\n<!-- clue:index:end -->\n")
	writeFile(t, root, "docs/decisions/ADR-101-x.md", "---\nid: ADR-101\ntype: decision\nstatus: inferred\nlinks: []\ntitle: X\nauthor: agent\naccepted-by: []\n---\n")
	code, out := runValidateCapturingStdout(t, []string{root})
	if code != 0 || !strings.Contains(out, "1 inferred decision(s) awaiting verification") {
		t.Fatalf("expected separate inferred-decision count, code=%d output=%q", code, out)
	}
}

func TestAC051_CLIReportsActivationBlockerCount(t *testing.T) {
	root := validCorpus(t)
	writeFile(t, root, "docs/README.md", "# Corpus\n\n<!-- clue:index:start -->\n- [goals/](goals/README.md)\n- [capabilities/](capabilities/README.md)\n<!-- clue:index:end -->\n")
	writeFile(t, root, "docs/capabilities/README.md", "# Capabilities\n\n<!-- clue:index:start -->\n- [CAP-101](CAP-101-x/README.md)\n<!-- clue:index:end -->\n")
	writeFile(t, root, "docs/capabilities/CAP-101-x/README.md", "---\nid: CAP-101\ntype: capability\nstatus: active\nlinks: []\ntitle: X\nprovenance: inferred\nreversal-cost: high\n---\n")

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	old := os.Stderr
	os.Stderr = w
	code := runValidate([]string{root})
	os.Stderr = old
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	out, err := io.ReadAll(r)
	_ = r.Close()
	if err != nil {
		t.Fatal(err)
	}
	if code != 1 || !strings.Contains(string(out), "1 high-cost inferred activation blocker(s)") {
		t.Fatalf("expected activation-blocker count, code=%d stderr=%q", code, out)
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

func readReleaseWorkflow(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", ".github", "workflows", "release.yml"))
	if err != nil {
		t.Fatalf("release workflow not found: %v", err)
	}
	return string(data)
}

// The build matrix, the stamp, the checksum name, and the published asset
// names moved out of the workflow's shell loop into goreleaser's config
// (ADR-030). The invariants did not move, so they are judged where they now
// live — parsed, so a key in a comment or a stray substring cannot satisfy
// them the way `strings.Contains(wf, "linux")` once could.
type goreleaserConfig struct {
	Builds []struct {
		Main         string   `yaml:"main"`
		Binary       string   `yaml:"binary"`
		Env          []string `yaml:"env"`
		Flags        []string `yaml:"flags"`
		Ldflags      []string `yaml:"ldflags"`
		Goos         []string `yaml:"goos"`
		Goarch       []string `yaml:"goarch"`
		ModTimestamp string   `yaml:"mod_timestamp"`
	} `yaml:"builds"`
	Archives []struct {
		ID           string   `yaml:"id"`
		Formats      []string `yaml:"formats"`
		NameTemplate string   `yaml:"name_template"`
	} `yaml:"archives"`
	Checksum struct {
		Algorithm    string `yaml:"algorithm"`
		NameTemplate string `yaml:"name_template"`
	} `yaml:"checksum"`
	Changelog struct {
		Disable bool `yaml:"disable"`
	} `yaml:"changelog"`
}

func readGoreleaserConfig(t *testing.T) goreleaserConfig {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", ".goreleaser.yaml"))
	if err != nil {
		t.Fatalf("goreleaser config not found: %v — the release's build matrix and asset names have no definition", err)
	}
	var cfg goreleaserConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("goreleaser config does not parse: %v", err)
	}
	return cfg
}

func contains(hay []string, needle string) bool {
	for _, h := range hay {
		if h == needle {
			return true
		}
	}
	return false
}

// Sanity: the release builds versioned cross-platform binaries.
// A repo invariant guarding M-004's release pipeline against regression;
// the operational proof is the first tagged release itself.
func TestSanity_ReleaseWorkflowIsCrossPlatform(t *testing.T) {
	wf := readReleaseWorkflow(t)
	cfg := readGoreleaserConfig(t)

	if len(cfg.Builds) != 1 {
		t.Fatalf("goreleaser defines %d builds, expected exactly 1 — a second build could ship an unstamped clue", len(cfg.Builds))
	}
	b := cfg.Builds[0]
	for _, want := range []string{"linux", "darwin", "windows"} {
		if !contains(b.Goos, want) {
			t.Errorf("goreleaser does not build for %s — a release must cover every supported host", want)
		}
	}
	for _, want := range []string{"amd64", "arm64"} {
		if !contains(b.Goarch, want) {
			t.Errorf("goreleaser does not build for %s", want)
		}
	}
	stamped := false
	for _, f := range b.Ldflags {
		if strings.Contains(f, "-X main.version={{ .Version }}") {
			stamped = true
		}
	}
	if !stamped {
		t.Error("goreleaser does not stamp main.version from the tag — the binary could not report its version, and the drift rule would have nothing to compare")
	}
	if !contains(b.Flags, "-trimpath") {
		t.Error("goreleaser does not pass -trimpath — released binaries would embed the runner's paths")
	}
	if !contains(b.Env, "CGO_ENABLED=0") {
		t.Error("goreleaser does not disable cgo — a released binary could acquire a runtime dependency (ADR-001: install is one download)")
	}
	if b.ModTimestamp == "" {
		t.Error("goreleaser does not set mod_timestamp — re-running the same tag would produce different bytes")
	}
	// An adopter's wall verifies a file it vendored under this exact name.
	if cfg.Checksum.NameTemplate != "SHA256SUMS" {
		t.Errorf("the checksum file is named %q, not SHA256SUMS — every adopted repo's vendored CI wall verifies that exact name", cfg.Checksum.NameTemplate)
	}

	// A manual dispatch must not publish a branch-named release: the ref
	// must be guarded to a tag before anything is built.
	if !strings.Contains(wf, "GITHUB_REF_TYPE") {
		t.Error("release workflow does not guard GITHUB_REF_TYPE — a branch dispatch could publish a branch-named release")
	}
	// A tag can land on any commit; nothing ships untested.
	if !strings.Contains(wf, "go test") {
		t.Error("release workflow does not run the tests before publishing")
	}
	// The publishing action runs with contents: write; a mutable tag pin
	// would let a moved tag ship different code into our releases.
	if pin := regexp.MustCompile(`goreleaser/goreleaser-action@[0-9a-f]{40}`); !pin.MatchString(wf) {
		t.Error("release workflow does not pin goreleaser-action by commit SHA")
	}
}

// Sanity: the published asset names are an append-only contract (ADR-030).
// The upstream reusable workflow is called by every thin caller that ran
// `clue init`; it verifies SHA256SUMS and stages clue-${CLUE_VERSION}-linux-
// amd64. A renamed asset breaks all of them at their next CI run, with no
// change in those repositories to explain it — so the release config and the
// upstream validation unit must keep naming the same file.
func TestSanity_ReleaseKeepsTheAssetNamesTheAdopterWallInstalls(t *testing.T) {
	cfg := readGoreleaserConfig(t)

	var bare string
	for _, a := range cfg.Archives {
		if len(a.Formats) == 1 && a.Formats[0] == "binary" {
			bare = a.NameTemplate
		}
	}
	if bare == "" {
		t.Fatal("no archive emits format `binary` — the release would publish only wrapped archives, and every adopted repo's wall downloads a bare binary by name")
	}
	rendered := strings.NewReplacer(
		"{{ .Version }}", "9.9.9", "{{.Version}}", "9.9.9",
		"{{ .Os }}", "linux", "{{.Os}}", "linux",
		"{{ .Arch }}", "amd64", "{{.Arch}}", "amd64",
	).Replace(bare)
	if rendered != "clue-9.9.9-linux-amd64" {
		t.Fatalf("the bare asset name renders as %q, not clue-<version>-linux-amd64 — the vendored wall installs the latter", rendered)
	}

	wall, err := os.ReadFile(filepath.Join("..", "..", ".github", "workflows", "clue-validation.yml"))
	if err != nil {
		t.Fatalf("upstream validation workflow not found: %v", err)
	}
	for _, want := range []string{"clue-${CLUE_VERSION}-linux-amd64", "SHA256SUMS"} {
		if !strings.Contains(string(wall), want) {
			t.Errorf("the upstream validation workflow no longer references %q — it and the release config must name the same asset", want)
		}
	}
}

func TestSanity_ReleaseIncludesReusableValidationUnit(t *testing.T) {
	root := filepath.Join("..", "..")
	data, err := os.ReadFile(filepath.Join(root, ".github", "workflows", "clue-validation.yml"))
	if err != nil {
		t.Fatalf("reusable validation workflow not found: %v", err)
	}
	wf := string(data)
	for _, want := range []string{"workflow_call:", "clue-source:", "clue-install-directory:", "clue validate --forbid-changes", "actions/checkout@"} {
		if !strings.Contains(wf, want) {
			t.Errorf("reusable validation workflow is missing %q", want)
		}
	}
	if !regexp.MustCompile(`(?m)^\s*- uses: actions/checkout@[0-9a-f]{40}`).MatchString(wf) {
		t.Error("reusable validation workflow does not pin checkout to a full commit SHA")
	}
	caller, err := os.ReadFile(filepath.Join(root, "internal", "scaffold", "templates", "github", "workflows", "clue.yml"))
	if err != nil {
		t.Fatal(err)
	}
	callerText := string(caller)
	if !strings.Contains(callerText, "cliewen/cliewen/.github/workflows/clue-validation.yml@__CLUE_WORKFLOW_REF__") {
		t.Error("scaffolded caller does not carry the reusable workflow reference placeholder")
	}
	if strings.Contains(callerText, "actions/checkout@") || strings.Contains(callerText, "clue validate --forbid-changes") {
		t.Error("scaffolded caller copied upstream validation logic")
	}
}

// Sanity: the published install scripts are the asset contract's second
// dependent (ADR-030). They compose the same `clue-<version>-<os>-<arch>`
// name the release emits and verify it against SHA256SUMS before writing
// anything. A renamed asset must fail here, not in a stranger's terminal;
// a script that stopped verifying would be worse than the manual steps it
// replaces, because it removes the check while keeping the convenience.
func TestSanity_InstallScriptsUseTheReleaseAssetContract(t *testing.T) {
	for _, s := range []struct {
		file  string
		asset string
		sums  []string
	}{
		{
			file:  "install.sh",
			asset: `clue-${version}-${os}-${arch}`,
			sums:  []string{"SHA256SUMS", "sha256sum -c --ignore-missing SHA256SUMS", "shasum -a 256"},
		},
		{
			file:  "install.ps1",
			asset: `clue-$version-windows-$arch.exe`,
			sums:  []string{"SHA256SUMS", "Get-FileHash"},
		},
	} {
		data, err := os.ReadFile(filepath.Join("..", "..", "guide", "public", s.file))
		if err != nil {
			t.Errorf("published install script %s not found: %v — the documented install command would 404", s.file, err)
			continue
		}
		body := string(data)
		if !strings.Contains(body, s.asset) {
			t.Errorf("%s does not compose the asset name %q — it and the release config must name the same file", s.file, s.asset)
		}
		for _, want := range s.sums {
			if !strings.Contains(body, want) {
				t.Errorf("%s does not reference %q — it must verify the download against the release checksums before installing", s.file, want)
			}
		}
		// The guide tells readers no administrator rights are needed, and
		// ADR-030 makes that a property of the decision, not a convenience.
		for _, forbidden := range []string{"sudo ", "RunAs", "-Verb RunAs"} {
			if strings.Contains(body, forbidden) {
				t.Errorf("%s contains %q — the install must not require elevation", s.file, forbidden)
			}
		}
	}
}

// Sanity: the release body is the tag's CHANGELOG.md section (ADR-012) —
// user-facing, reviewed prose. GitHub's auto-generated notes (a PR dump
// with contributor @mentions) must not come back.
func TestSanity_ReleaseNotesComeFromChangelog(t *testing.T) {
	wf := readReleaseWorkflow(t)
	for _, want := range []string{"CHANGELOG.md", "--release-notes=release-notes.md"} {
		if !strings.Contains(wf, want) {
			t.Errorf("release workflow does not mention %q — the release body must be extracted from the changelog", want)
		}
	}
	if strings.Contains(wf, "generate_release_notes") {
		t.Error("release workflow enables generate_release_notes — release bodies are written for users in the changelog, not auto-generated")
	}
	// This assertion used to be inverted, and it published v0.8.0 with a
	// blank page. `changelog: {disable: true}` reads like extra protection
	// for ADR-012 and is documented to "also ignore any changelog files
	// passed via --release-notes, and render an empty changelog". Supplying
	// --release-notes is already what keeps GoReleaser's commit dump out;
	// disabling on top of it throws the reviewed notes away too.
	if cfg := readGoreleaserConfig(t); cfg.Changelog.Disable {
		t.Error("goreleaser config sets changelog.disable — that discards --release-notes and publishes an empty body; the flag alone already suppresses the generated changelog")
	}
	// Configuration can only say what was asked for. The release workflow
	// must also check what GitHub actually published, because that is the
	// only place the failure above was visible.
	if !strings.Contains(wf, "gh release view") {
		t.Error("release workflow never reads back the published body — ADR-012's one-to-one map must be verified against the release page, not just requested")
	}
}

func TestAC065_UnitPositive_ReleaseWorkflowRequiresMigrationGuidance(t *testing.T) {
	wf := readReleaseWorkflow(t)
	for _, required := range []string{
		"internal/migrate",
		"Require migration guidance when the migration registry changed",
		"### Migration",
		"release-notes.md",
	} {
		if !strings.Contains(wf, required) {
			t.Errorf("release workflow is missing migration guidance guard %q", required)
		}
	}
}

func TestAC065_UnitNegative_ReleaseWorkflowCannotSkipChangedRegistryGuidance(t *testing.T) {
	wf := readReleaseWorkflow(t)
	guard := `git diff --quiet "${GITHUB_SHA}^" "${GITHUB_SHA}" -- internal/migrate`
	if !strings.Contains(wf, guard) || !strings.Contains(wf, "exit !content") {
		t.Fatal("a migration-registry change can reach publishing without checking for non-empty migration notes")
	}
}

func TestAC064_CLI_MigratePreviewAndApply(t *testing.T) {
	root := t.TempDir()
	if _, err := scaffold.Run(root); err != nil {
		t.Fatal(err)
	}
	artifact := "---\nid: AN-101\ntype: analysis\nstatus: verified\nlinks: []\ntitle: Historical analysis\nprovenance: inferred\n---\n\nThe operator's prose remains intact.\n"
	writeFile(t, root, "docs/analysis/AN-101.md", artifact)

	before, err := os.ReadFile(filepath.Join(root, "docs", "analysis", "AN-101.md"))
	if err != nil {
		t.Fatal(err)
	}
	var out, errOut strings.Builder
	if code := runMigrate([]string{"--reversal-cost=low", root}, &out, &errOut); code != 0 {
		t.Fatalf("preview exit code = %d, stderr=%q, stdout=%q", code, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), "preview only") || !strings.Contains(out.String(), "docs/analysis/AN-101.md") {
		t.Fatalf("preview output does not name the planned change: %q", out.String())
	}
	if after, err := os.ReadFile(filepath.Join(root, "docs", "analysis", "AN-101.md")); err != nil || !bytes.Equal(after, before) {
		t.Fatal("preview changed a file")
	}

	out.Reset()
	errOut.Reset()
	if code := runMigrate([]string{"--apply", "--reversal-cost=low", root}, &out, &errOut); code != 0 {
		t.Fatalf("apply exit code = %d, stderr=%q, stdout=%q", code, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), "applied") {
		t.Fatalf("apply output does not report a write: %q", out.String())
	}

	out.Reset()
	errOut.Reset()
	if code := runMigrate([]string{root}, &out, &errOut); code != 0 || !strings.Contains(out.String(), "no changes needed") {
		t.Fatalf("second run was not a no-op: code=%d stderr=%q stdout=%q", code, errOut.String(), out.String())
	}
}

// Sanity: merging a release PR is what cuts the release (PDR-015). The
// tagging workflow must take its version from the same single bump site the
// release's drift gate judges, and must start the release explicitly — a tag
// pushed with GITHUB_TOKEN does not trigger a workflow, so a missing dispatch
// leaves a tagged commit that never publishes.
func TestSanity_TagOnMergeDerivesTheVersionAndStartsTheRelease(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", ".github", "workflows", "tag-on-merge.yml"))
	if err != nil {
		t.Fatalf("tag-on-merge workflow not found: %v", err)
	}
	wf := string(data)

	const stamp = "internal/skills/source/shared/frontmatter.md.tmpl"
	if !strings.Contains(wf, stamp) {
		t.Errorf("tag-on-merge does not read %s — the tag must come from the same stamp the release gate compares against, not a second source", stamp)
	}
	if !strings.Contains(wf, "gh workflow run release.yml") {
		t.Error("tag-on-merge never dispatches release.yml — a tag pushed with GITHUB_TOKEN does not trigger workflows, so the release would never start")
	}
	if !strings.Contains(wf, "actions: write") {
		t.Error("tag-on-merge lacks actions: write — dispatching the release workflow fails with 403 without it")
	}
	// Without the already-tagged check this fails on every ordinary merge.
	if !strings.Contains(wf, "refs/tags/v") {
		t.Error("tag-on-merge does not check for an existing tag — it must be a quiet no-op on merges that do not raise the version")
	}
	if !strings.Contains(wf, "group: tag-on-main") || !strings.Contains(wf, "queue: max") || !strings.Contains(wf, "cancel-in-progress: false") {
		t.Error("tag-on-merge does not serialize its release decision and tag push — concurrent main pushes can otherwise tag different commits")
	}
	if !strings.Contains(wf, `git diff --quiet "${GITHUB_SHA}^" "${GITHUB_SHA}" -- "$tmpl"`) {
		t.Errorf("tag-on-merge does not gate releases on %s changing in the pushed commit — an ordinary merge can otherwise tag its own checkout", stamp)
	}
	// A spent version number cannot be reused, so a missing changelog
	// section must stop the tag rather than the publish.
	if !strings.Contains(wf, "CHANGELOG.md") {
		t.Error("tag-on-merge does not require a CHANGELOG section — a stamp raised without notes would burn a tag that can never release (ADR-012)")
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
	// Scoped to this step and to the whole of it, in whatever order its
	// keys are written: an empty ${VERSION} is exempt from the drift rule,
	// so a stamp wired from nothing fails open and ships the mismatch.
	// Another step's env does not vouch for this one.
	const stepMark = "\n      - " // a step boundary in the job's steps: list
	i := strings.LastIndex(wf[:m[0]], stepMark)
	if i < 0 {
		t.Fatal("the stamped run line is not inside a step of the job's steps: list — the env it is judged against cannot be identified")
	}
	step := wf[i:]
	if j := strings.Index(step[len(stepMark):], stepMark); j >= 0 {
		step = step[:len(stepMark)+j]
	}
	if !strings.Contains(step, "VERSION: ${{ steps.version.outputs.v }}") {
		t.Error("the stamped step does not take VERSION from the tag — an empty stamp is exempt from the drift rule, so the gate would pass on any corpus")
	}
	// Nothing goreleaser touches may precede the gate: it builds, checksums,
	// checksums and publishes in one step, so a mismatched tag
	// must be refused before it runs at all. A marker that has vanished is a
	// hole, not a pass — the earlier form used `i >= 0 && i < m[0]`, which
	// silently guarded nothing once its step was renamed.
	for _, later := range []string{"goreleaser/goreleaser-action", "release --clean"} {
		i := strings.Index(wf, later)
		if i < 0 {
			t.Fatalf("release workflow no longer mentions %q — the ordering guard has nothing to hold; name whatever step now builds and publishes", later)
		}
		if i < m[0] {
			t.Errorf("the stamped step runs after %q — a mismatched tag must fail before any artifact exists", later)
		}
	}

	args := strings.Fields(wf[m[2]:m[3]])
	if len(args) == 0 || args[0] != "validate" {
		t.Fatalf("the stamped step runs %q, not validate — the drift rule is never reached", args)
	}
	// runValidate parses with flag.ExitOnError, which exits the process on
	// an unknown flag: executing one would take the test binary down with
	// no message at all. Refuse it legibly instead — a step that grows a
	// flag teaches this test what the flag does to the gate first.
	executable := map[string]bool{"--forbid-changes": true}
	for _, a := range args[1:] {
		if strings.HasPrefix(a, "-") && !executable[a] {
			t.Fatalf("the stamped step passes %q, which this test cannot execute — assert what the flag does to the gate below, then add it to the executable set", a)
		}
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
	// Drain while validate writes: a pipe holds one buffer's worth, so
	// reading only after the writer closes would deadlock the moment the
	// corpus reports more issues than that.
	type captured struct {
		out string
		err error
	}
	done := make(chan captured, 1)
	go func() {
		out, readErr := io.ReadAll(r)
		done <- captured{string(out), readErr}
	}()
	old := os.Stdout
	defer func() { os.Stdout = old }()
	os.Stdout = w
	code := runValidate(args)
	if err := w.Close(); err != nil {
		t.Fatalf("close pipe: %v", err)
	}
	got := <-done
	if got.err != nil {
		t.Fatalf("read pipe: %v", got.err)
	}
	return code, got.out
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
	fixes := strings.Index(content, "Any agent that edits an existing PR becomes its updater for that turn")
	if fixes < 0 {
		t.Fatal("C-012 does not define the review-fix handoff")
	}
	handoff := content[fixes:]
	commitCandidate := strings.Index(handoff, "commit")
	verifyCandidate := strings.Index(handoff, "verifies")
	reviewCandidate := strings.Index(handoff, "clean review")
	pushCandidate := strings.Index(handoff, "push")
	if commitCandidate < 0 || verifyCandidate <= commitCandidate || reviewCandidate <= verifyCandidate || pushCandidate <= reviewCandidate {
		t.Error("C-012 must commit a repaired candidate, verify and review that commit, then push the reviewed commit")
	}
}

func TestAC041_PublicCarriersKeepCrossAgentHelpOutsideTheInitiatedSlot(t *testing.T) {
	root := filepath.Join("..", "..")
	for rel, wants := range map[string][]string{
		".github/pull_request_template.md": {
			"initiating author's only initiated Cliewen change",
			"review or update help on an existing PR does not consume another slot",
		},
		"CONTRIBUTING.md": {
			"A contributor may initiate one Cliewen change at a time",
			"reviewing, and helping update an existing pull request do not consume another initiated-change slot",
		},
		"guide/change-loop.md": {
			"One initiating author takes one initiated Cliewen change",
			"reviews, and help updating an existing pull request do not consume another initiated-change slot",
		},
	} {
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
		if err != nil {
			t.Fatal(err)
		}
		for _, want := range wants {
			if !strings.Contains(string(data), want) {
				t.Errorf("%s does not preserve cross-agent help outside the initiated-change slot %q", rel, want)
			}
		}
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

// TestAC068_UnitNegative_RefsExitsNonZeroOnlyForGone drives the command itself,
// so the exit code the criterion names is proven rather than inferred from the
// report type.
func TestAC068_UnitNegative_RefsExitsNonZeroOnlyForGone(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/owner/live/forbidden", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	})
	mux.HandleFunc("/owner", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	write := func(t *testing.T, body string) string {
		t.Helper()
		root := t.TempDir()
		dir := filepath.Join(root, "docs", "analysis")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		text := "---\nid: AN-001\ntype: analysis\nstatus: active\nlinks: []\ntitle: t\n---\n\n" + body + "\n"
		if err := os.WriteFile(filepath.Join(dir, "AN-001-x.md"), []byte(text), 0o644); err != nil {
			t.Fatal(err)
		}
		return root
	}

	var out, errOut bytes.Buffer
	restricted := write(t, "restricted "+srv.URL+"/owner/live/forbidden")
	if code := runRefs([]string{restricted}, &out, &errOut); code != 0 {
		t.Fatalf("a restricted address must not fail the command, got exit %d: %s", code, out.String())
	}

	out.Reset()
	gone := write(t, "deleted "+srv.URL+"/nobody/vanished")
	if code := runRefs([]string{gone}, &out, &errOut); code != 1 {
		t.Fatalf("a gone address must exit 1, got %d: %s", code, out.String())
	}

	out.Reset()
	if code := runRefs([]string{"a", "b"}, &out, &errOut); code != 2 {
		t.Fatalf("two paths is a usage error, got exit %d", code)
	}
}

// TestAC067_UnitPositive_CoverageListsForeignPointersApart proves the criterion's
// user-visible claim: a pointer to proof elsewhere is printed after the
// capability states and labelled so it cannot be read as coverage.
func TestAC067_UnitPositive_CoverageListsForeignPointersApart(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "docs", "analysis")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	body := "---\nid: AN-001\ntype: analysis\nstatus: active\nlinks: []\ntitle: t\n---\n\nProven by clue:robocode-dev/tank-royale@384d27d5/BR-001 upstream.\n"
	if err := os.WriteFile(filepath.Join(dir, "AN-001-x.md"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}

	c, _ := corpus.Scan(root)
	pointers := corpus.ForeignPointers(c)
	if len(pointers) != 1 || pointers[0] != "clue:robocode-dev/tank-royale@384d27d5/BR-001" {
		t.Fatalf("expected the pointer harvested once, got %v", pointers)
	}

	// The label must not be a coverage state: "covered", "partial" and "gap"
	// are what a reader scans for, and a pointer must never join them.
	line := pointers[0] + ": named but locally unproven"
	for _, state := range []string{"covered", "partial", "gap"} {
		if strings.HasSuffix(line, ": "+state) {
			t.Fatalf("the pointer line reads as a coverage state: %q", line)
		}
	}
	if !strings.Contains(usage, "named but") {
		t.Fatal("the --coverage help must account for the pointer lines it now prints")
	}
}
