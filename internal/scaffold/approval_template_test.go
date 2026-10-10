package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUnit_InitShipsApprovalTemplateWithoutOverwritingActualRecords(t *testing.T) {
	root := t.TempDir()
	if _, e := Run(root); e != nil {
		t.Fatal(e)
	}
	filename := filepath.Join(root, ".clue/acceptance/approval.yaml")
	data, e := os.ReadFile(filename)
	if e != nil {
		t.Fatal(e)
	}
	for _, field := range []string{"type: acceptance-approval", "decision: accept", "candidate:", "base:", "brief-sha256:", "approved-by:", "approval-source:", "statement:", "recorded-at:"} {
		if !strings.Contains(string(data), field) {
			t.Fatalf("missing %s", field)
		}
	}
	actual := []byte("repository-owned approval preparation\n")
	if e = os.WriteFile(filename, actual, 0644); e != nil {
		t.Fatal(e)
	}
	if _, e = Run(root); e != nil {
		t.Fatal(e)
	}
	after, e := os.ReadFile(filename)
	if e != nil || string(after) != string(actual) {
		t.Fatal("existing record overwritten")
	}
}
