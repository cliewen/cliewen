package scaffold

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/cliewen/cliewen/internal/role"
)

// A fresh repository is an adopter, and init says so. The marker is what
// every rule that differs by repository kind reads, so a new corpus must
// not start out silent about it.
func TestAC153_UnitPositive_InitMaterializesAnAdopterMarker(t *testing.T) {
	root, rep := runInto(t)
	got, declared, err := role.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if !declared || got != role.Adopter {
		t.Fatalf("init produced role (%q, declared=%v), want adopter", got, declared)
	}
	if !slices.Contains(rep.Created, role.DefaultPath) {
		t.Fatalf("init did not report creating %s: %v", role.DefaultPath, rep.Created)
	}
}

// init is a materializer, never an updater. Cliewen's own repository
// declares role: source by hand, and re-running init there must not quietly
// demote it to an adopter.
func TestAC153_UnitNegative_InitNeverOverwritesADeclaredRole(t *testing.T) {
	root := t.TempDir()
	if err := role.Write(root, role.Source); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(filepath.Join(root, role.DefaultPath))
	if err != nil {
		t.Fatal(err)
	}
	rep, err := Run(root)
	if err != nil {
		t.Fatal(err)
	}
	after, err := os.ReadFile(filepath.Join(root, role.DefaultPath))
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Fatalf("init rewrote a declared role: %q became %q", before, after)
	}
	if !slices.Contains(rep.Skipped, role.DefaultPath) {
		t.Fatalf("init did not report skipping %s: %v", role.DefaultPath, rep.Skipped)
	}
	if slices.Contains(rep.Created, role.DefaultPath) {
		t.Fatalf("init reported creating a marker that already existed: %v", rep.Created)
	}
}
