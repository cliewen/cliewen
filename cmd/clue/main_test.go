package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/cliewen/cliewen/internal/corpus"
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

// AC-022 (wiring): runValidate threads the binary's stamp through
// corpus.Options.Version into the drift rule — a released clue fails
// against lagging skills, a matching release passes.
func TestAC022_RunValidateThreadsVersionIntoDriftRule(t *testing.T) {
	root := validCorpus(t)
	writeFile(t, root, ".agents/skills/clue-delta/skill.md", "---\nversion: 0.1.0\n---\n\n# clue-delta\n")
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
