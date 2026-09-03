package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cliewen/cliewen/internal/corpus"
)

// AC-166: init materializes the marked vision and the empty use-case folder,
// and the corpus index references both.
func TestAC166_UnitPositive_InitMaterializesTheVisionAndTheUseCaseFolder(t *testing.T) {
	root, rep := runInto(t)
	created := strings.Join(rep.Created, "\n")
	for _, rel := range []string{corpus.VisionPath, "docs/use-cases/README.md"} {
		if !strings.Contains(created, rel) {
			t.Fatalf("init did not create %s:\n%s", rel, created)
		}
	}
	vision, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(corpus.VisionPath)))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(vision), corpus.VisionBootstrapMarker) {
		t.Fatal("the materialized vision carries no bootstrap marker, so nothing tells an agent to replace it")
	}
	useCases, err := os.ReadFile(filepath.Join(root, filepath.FromSlash("docs/use-cases/README.md")))
	if err != nil {
		t.Fatal(err)
	}
	// Empty and staying empty is the normal outcome for an optional folder.
	between := string(useCases)
	body := between[strings.Index(between, IndexStart)+len(IndexStart) : strings.Index(between, IndexEnd)]
	if strings.TrimSpace(body) != "" {
		t.Fatalf("the use-case index is not empty: %q", body)
	}
	index, err := os.ReadFile(filepath.Join(root, filepath.FromSlash("docs/README.md")))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"(vision.md)", "use-cases/"} {
		if !strings.Contains(string(index), want) {
			t.Fatalf("the corpus index does not reference %s:\n%s", want, index)
		}
	}
	if issues := activatedValidateAt(t, root); len(issues) > 0 {
		t.Fatalf("expected green after activating the bootstraps, got %v", issues)
	}
}

// AC-166: a repository that already carries either file keeps it. Init is a
// materializer, never an updater.
func TestAC166_UnitNegative_InitNeverOverwritesAnExistingVisionOrFolder(t *testing.T) {
	root := t.TempDir()
	const mine = "---\nid: VIS-001\ntype: vision\nstatus: active\nlinks: []\ntitle: Mine\n---\n\n# VIS-001\n\nMy own direction.\n"
	full := filepath.Join(root, filepath.FromSlash(corpus.VisionPath))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(mine), 0o644); err != nil {
		t.Fatal(err)
	}
	rep, err := Run(root)
	if err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(full)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != mine {
		t.Fatalf("init overwrote an existing vision:\n%s", got)
	}
	if !strings.Contains(strings.Join(rep.Skipped, "\n"), corpus.VisionPath) {
		t.Fatalf("init did not report the existing vision as skipped: %v", rep.Skipped)
	}
}
