package main

import (
	"os"
	"path/filepath"
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

// Sanity: the release workflow builds versioned cross-platform binaries.
// A repo invariant guarding M-004's release pipeline against regression;
// the operational proof is the first tagged release itself.
func TestSanity_ReleaseWorkflowIsCrossPlatform(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", ".github", "workflows", "release.yml"))
	if err != nil {
		t.Fatalf("release workflow not found: %v", err)
	}
	wf := string(data)
	for _, want := range []string{"main.version", "linux", "darwin", "windows", "arm64", "amd64"} {
		if !strings.Contains(wf, want) {
			t.Errorf("release workflow does not mention %q — expected a stamped cross-platform build", want)
		}
	}
	// A manual dispatch must not publish a branch-named release: the ref
	// must be guarded to a tag before anything is built.
	if !strings.Contains(wf, "GITHUB_REF_TYPE") {
		t.Error("release workflow does not guard GITHUB_REF_TYPE — a branch dispatch could publish a branch-named release")
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
