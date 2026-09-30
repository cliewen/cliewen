package migrate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cliewen/cliewen/internal/evidence"
)

func migrationEvidenceFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	p := filepath.Join(root, "docs/capabilities/x")
	if err := os.MkdirAll(p, 0755); err != nil {
		t.Fatal(err)
	}
	text := "---\nid: CAP-900-criteria\ntype: criteria\nstatus: active\nlinks: []\ntitle: Evidence fixture\n---\n\n@AC-900\nScenario: proof\n  Test-type: Unit\n  Then a result\n"
	if err := os.WriteFile(filepath.Join(p, "criteria.md"), []byte(text), 0644); err != nil {
		t.Fatal(err)
	}
	return root
}

func TestAC208_UnitPositive_MigrationNamesMissingExportWithoutWritingTests(t *testing.T) {
	root := migrationEvidenceFixture(t)
	plan := MigrationPlan{}
	planEvidenceExport(root, &plan)
	if len(plan.Notices) != 1 || plan.Notices[0].Migration != MigrationEvidenceExport || len(plan.Changes) != 0 {
		t.Fatalf("missing evidence notice: %+v", plan)
	}
	if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(evidence.DefaultPath))); !os.IsNotExist(err) {
		t.Fatal("migration fabricated evidence")
	}
}

func TestAC208_UnitNegative_CurrentExportNeedsNoEstablishmentNotice(t *testing.T) {
	root := migrationEvidenceFixture(t)
	if err := evidence.Write(root, evidence.Manifest{Version: 1, Producers: []evidence.Producer{{ID: "suite", Include: []string{"tests/**/*.unknown"}}}}); err != nil {
		t.Fatal(err)
	}
	plan := MigrationPlan{}
	planEvidenceExport(root, &plan)
	if len(plan.Notices) > 0 || len(plan.Changes) > 0 {
		t.Fatalf("already established: %+v", plan)
	}
}
