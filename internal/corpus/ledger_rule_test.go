package corpus

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeLedger(t *testing.T, root, yamlBody string) {
	t.Helper()
	dir := filepath.Join(root, ".clue")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "id-ledger.yaml"), []byte(yamlBody), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestAC103_UnitPositive_NonLiveLedgerIDRejectedWhenLiveAgain(t *testing.T) {
	root := writeCorpus(t, validFiles)
	writeLedger(t, root, "counters: {G: 1}\nentries:\n  - id: G-001\n    kind: numeric\n    state: retired\n    prefix: G\n    component: 1\n")
	c, scanIssues := Scan(root)
	if len(scanIssues) != 0 {
		t.Fatalf("scan issues: %v", scanIssues)
	}
	issues := Validate(c, Options{})
	found := false
	for _, is := range issues {
		if is.Msg == "id G-001 is marked retired in .clue/id-ledger.yaml and cannot be used by a live artifact" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected a retired-id issue, got: %v", issues)
	}
}

func TestAC103_UnitNegative_LiveOrReservedIDPasses(t *testing.T) {
	root := writeCorpus(t, validFiles)
	writeLedger(t, root, "counters: {G: 1}\nentries:\n  - id: G-001\n    kind: numeric\n    state: live\n    prefix: G\n    component: 1\n")
	c, scanIssues := Scan(root)
	if len(scanIssues) != 0 {
		t.Fatalf("scan issues: %v", scanIssues)
	}
	issues := Validate(c, Options{})
	for _, is := range issues {
		if strings.Contains(is.Msg, "id G-001 is marked") {
			t.Fatalf("unexpected state issue for a live entry: %v", issues)
		}
	}
}

func TestAC103_UnitPositive_ReservedIDRejectedWhenLiveAgain(t *testing.T) {
	root := writeCorpus(t, validFiles)
	writeLedger(t, root, "counters: {G: 1}\nentries:\n  - id: G-001\n    kind: numeric\n    state: reserved\n    prefix: G\n    component: 1\n")
	c, scanIssues := Scan(root)
	if len(scanIssues) != 0 {
		t.Fatalf("scan issues: %v", scanIssues)
	}
	issues := Validate(c, Options{})
	for _, is := range issues {
		if is.Msg == "id G-001 is marked reserved in .clue/id-ledger.yaml and cannot be used by a live artifact" {
			return
		}
	}
	t.Fatalf("expected a reserved-id issue, got: %v", issues)
}

func TestAC103_UnitPositive_LiveIDMissingFromLedgerRejected(t *testing.T) {
	root := writeCorpus(t, validFiles)
	writeLedger(t, root, "counters: {}\nentries: []\n")
	c, scanIssues := Scan(root)
	if len(scanIssues) != 0 {
		t.Fatalf("scan issues: %v", scanIssues)
	}
	issues := Validate(c, Options{})
	for _, is := range issues {
		if is.Msg == "id G-001 is missing from .clue/id-ledger.yaml" {
			return
		}
	}
	t.Fatalf("expected a missing-ledger-entry issue, got: %v", issues)
}

func TestAC103_UnitPositive_NonLiveLedgerCriterionRejectedWhenDeclaredLive(t *testing.T) {
	root := writeCorpus(t, capFiles("active"))
	writeLedger(t, root, "counters: {G: 1, CAP: 101, AC: 101}\nentries:\n  - id: G-001\n    kind: numeric\n    state: live\n    prefix: G\n    component: 1\n  - id: CAP-101\n    kind: numeric\n    state: live\n    prefix: CAP\n    component: 101\n  - id: CAP-101-criteria\n    kind: opaque\n    state: live\n  - id: AC-101\n    kind: numeric\n    state: retired\n    prefix: AC\n    component: 101\n")
	c, scanIssues := Scan(root)
	if len(scanIssues) != 0 {
		t.Fatalf("scan issues: %v", scanIssues)
	}
	issues := Validate(c, Options{})
	for _, is := range issues {
		if is.Msg == "criterion AC-101 is marked retired in .clue/id-ledger.yaml but its declaration is live" {
			return
		}
	}
	t.Fatalf("expected a retired-criterion issue, got: %v", issues)
}

func TestAC103_UnitPositive_RetiredCriterionInDraftArtifactMustRemainRetired(t *testing.T) {
	files := with(capFiles("active"), map[string]string{
		"docs/capabilities/CAP-101-x/retired-criteria.md": "---\nid: CAP-101-retired-criteria\ntype: criteria\nstatus: draft\nlinks: [CAP-101]\ntitle: Retired criterion\n---\n\n```gherkin\n\n  @AC-102 @retired\n  Scenario: retired\n```\n",
	})
	root := writeCorpus(t, files)
	writeLedger(t, root, "counters: {AC: 102}\nentries:\n  - id: AC-102\n    kind: numeric\n    state: live\n    prefix: AC\n    component: 102\n")
	c, scanIssues := Scan(root)
	if len(scanIssues) != 0 {
		t.Fatalf("scan issues: %v", scanIssues)
	}
	issues := Validate(c, Options{})
	for _, is := range issues {
		if is.Msg == "criterion AC-102 is marked live in .clue/id-ledger.yaml but its declaration is retired" {
			return
		}
	}
	t.Fatalf("expected a draft-tombstone ledger issue, got: %v", issues)
}

func TestAC104_UnitPositive_UnknownLedgerStateRejectedForLiveArtifact(t *testing.T) {
	root := writeCorpus(t, validFiles)
	writeLedger(t, root, "counters: {G: 1}\nentries:\n  - id: G-001\n    kind: numeric\n    state: pending\n    prefix: G\n    component: 1\n")
	c, scanIssues := Scan(root)
	if len(scanIssues) != 0 {
		t.Fatalf("scan issues: %v", scanIssues)
	}
	issues := Validate(c, Options{})
	want := map[string]bool{
		"entry G-001 has invalid state pending":                                                    false,
		"id G-001 is marked pending in .clue/id-ledger.yaml and cannot be used by a live artifact": false,
	}
	for _, is := range issues {
		if _, ok := want[is.Msg]; ok {
			want[is.Msg] = true
		}
	}
	for msg, found := range want {
		if !found {
			t.Fatalf("missing expected issue %q in: %v", msg, issues)
		}
	}
}

func TestAC103_UnitNegative_NoLedgerFileIsUnaffected(t *testing.T) {
	root := writeCorpus(t, validFiles)
	c, scanIssues := Scan(root)
	if len(scanIssues) != 0 {
		t.Fatalf("scan issues: %v", scanIssues)
	}
	issues := Validate(c, Options{})
	if len(issues) != 0 {
		t.Fatalf("expected no issues from a corpus without a ledger file, got: %v", issues)
	}
}

func TestAC104_UnitPositive_MalformedLedgerEntryShapeRejected(t *testing.T) {
	root := writeCorpus(t, validFiles)
	writeLedger(t, root, "counters: {}\nentries:\n  - id: G-900\n    kind: numeric\n    state: live\n  - id: P-001\n    kind: numeric\n    state: live\n    prefix: G\n    component: 99\n  - id: AC-imported-uuid\n    kind: opaque\n    state: live\n    prefix: AC\n    component: 5\n  - id: unknown-001\n    kind: other\n    state: live\n")
	c, scanIssues := Scan(root)
	if len(scanIssues) != 0 {
		t.Fatalf("scan issues: %v", scanIssues)
	}
	issues := Validate(c, Options{})
	wantMsgs := map[string]bool{
		"entry G-900 is numeric-kind but its ID, prefix, and component do not agree": false,
		"entry P-001 is numeric-kind but its ID, prefix, and component do not agree": false,
		"entry AC-imported-uuid is opaque-kind but carries numeric fields":           false,
		"entry unknown-001 has invalid kind other":                                   false,
	}
	for _, is := range issues {
		if _, ok := wantMsgs[is.Msg]; ok {
			wantMsgs[is.Msg] = true
		}
	}
	for msg, got := range wantMsgs {
		if !got {
			t.Fatalf("missing expected issue %q in: %v", msg, issues)
		}
	}
}

func TestAC104_UnitNegative_WellFormedEntriesPass(t *testing.T) {
	root := writeCorpus(t, validFiles)
	writeLedger(t, root, "counters: {G: 1, P: 1}\nentries:\n  - id: G-001\n    kind: numeric\n    state: live\n    prefix: G\n    component: 1\n  - id: P-001\n    kind: numeric\n    state: live\n    prefix: P\n    component: 1\n  - id: AC-imported-uuid\n    kind: opaque\n    state: live\n")
	c, scanIssues := Scan(root)
	if len(scanIssues) != 0 {
		t.Fatalf("scan issues: %v", scanIssues)
	}
	issues := Validate(c, Options{})
	if len(issues) != 0 {
		t.Fatalf("expected no shape issues for well-formed entries, got: %v", issues)
	}
}

func TestAC106_UnitPositive_RetiredOpaqueIDRejectedWhenLiveAgain(t *testing.T) {
	root := writeCorpus(t, with(validFiles, map[string]string{
		"docs/goals/G-002-second.md": "---\nid: 8f14e45f-ceea-467e-9a2b-a1c8b9d2f7a1\ntype: goal\nstatus: accepted\nlinks: []\ntitle: Imported goal\n---\n",
		"docs/goals/README.md":       "# Goals\n\n<!-- clue:index:start -->\n- [G-001](G-001-first.md)\n- [8f14e45f-ceea-467e-9a2b-a1c8b9d2f7a1](G-002-second.md)\n<!-- clue:index:end -->\n",
	}))
	writeLedger(t, root, "counters: {}\nentries:\n  - id: 8f14e45f-ceea-467e-9a2b-a1c8b9d2f7a1\n    kind: opaque\n    state: retired\n")
	c, scanIssues := Scan(root)
	if len(scanIssues) != 0 {
		t.Fatalf("scan issues: %v", scanIssues)
	}
	issues := Validate(c, Options{})
	found := false
	for _, is := range issues {
		if is.Msg == "id 8f14e45f-ceea-467e-9a2b-a1c8b9d2f7a1 is marked retired in .clue/id-ledger.yaml and cannot be used by a live artifact" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected a retired-opaque-id issue, got: %v", issues)
	}
}

func TestAC106_UnitNegative_FreshDistinctOpaqueIDPasses(t *testing.T) {
	root := writeCorpus(t, with(validFiles, map[string]string{
		"docs/goals/G-002-second.md": "---\nid: 8f14e45f-ceea-467e-9a2b-a1c8b9d2f7a1\ntype: goal\nstatus: accepted\nlinks: []\ntitle: Imported goal\n---\n",
		"docs/goals/README.md":       "# Goals\n\n<!-- clue:index:start -->\n- [G-001](G-001-first.md)\n- [8f14e45f-ceea-467e-9a2b-a1c8b9d2f7a1](G-002-second.md)\n<!-- clue:index:end -->\n",
	}))
	// A distinct retired opaque ID must never block an unrelated freshly
	// generated one from passing.
	writeLedger(t, root, "counters: {}\nentries:\n  - id: 11111111-1111-1111-1111-111111111111\n    kind: opaque\n    state: retired\n  - id: 8f14e45f-ceea-467e-9a2b-a1c8b9d2f7a1\n    kind: opaque\n    state: live\n")
	c, scanIssues := Scan(root)
	if len(scanIssues) != 0 {
		t.Fatalf("scan issues: %v", scanIssues)
	}
	issues := Validate(c, Options{})
	for _, is := range issues {
		if strings.Contains(is.Msg, "retired") {
			t.Fatalf("unexpected retired-id issue for an unrelated, live opaque id: %v", issues)
		}
	}
}
