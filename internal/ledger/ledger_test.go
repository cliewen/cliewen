package ledger

import (
	"math/big"
	"os"
	"path/filepath"
	"testing"
)

func decimal(n int64) *big.Int { return big.NewInt(n) }

func nextNumeric(t *testing.T, l *Ledger, prefix string) string {
	t.Helper()
	id, err := l.NextNumeric(prefix)
	if err != nil {
		t.Fatalf("NextNumeric(%q): %v", prefix, err)
	}
	return id
}

func TestAC101_UnitPositive_NextNumericIncrementsStoredCounter(t *testing.T) {
	root := t.TempDir()
	l, err := Load(root)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	first := nextNumeric(t, l, "PDR")
	if first != "PDR-001" {
		t.Fatalf("first id = %q, want PDR-001", first)
	}
	if err := l.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}

	l2, err := Load(root)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	second := nextNumeric(t, l2, "PDR")
	if second != "PDR-002" {
		t.Fatalf("second id = %q, want PDR-002 (counter must persist across Load)", second)
	}
	if e, ok := l2.Lookup("PDR-001"); !ok || e.State != StateReserved {
		t.Fatalf("PDR-001 entry = %+v, ok=%v, want reserved", e, ok)
	}
}

func TestAC101_UnitNegative_PrefixAbsentFromLedgerStartsAtOne(t *testing.T) {
	root := t.TempDir()
	l, _ := Load(root)
	id := nextNumeric(t, l, "ZZZ")
	if id != "ZZZ-001" {
		t.Fatalf("id = %q, want ZZZ-001 for an unseen prefix", id)
	}
}

func TestAC101_UnitNegative_NextNumericRejectsNonCanonicalPrefix(t *testing.T) {
	root := t.TempDir()
	l, _ := Load(root)
	for _, prefix := range []string{"snap-sqs", "SNAP--SQS", "SNAP_SQS"} {
		if _, err := l.NextNumeric(prefix); err == nil {
			t.Fatalf("NextNumeric accepted non-canonical prefix %q", prefix)
		}
	}
}

func TestAC102_UnitPositive_OpaqueIDPreservedVerbatim(t *testing.T) {
	root := t.TempDir()
	l, _ := Load(root)
	const id = "AC-8f14e45f-ceea-467e-9a2b-a1c8b9d2f7a1"
	if err := l.ReserveOpaque(id, "abc123", "legacy/tests.feature"); err != nil {
		t.Fatalf("ReserveOpaque: %v", err)
	}
	e, ok := l.Lookup(id)
	if !ok || e.ID != id || e.Kind != KindOpaque {
		t.Fatalf("Lookup(%q) = %+v, ok=%v, want opaque entry with exact id", id, e, ok)
	}
}

func TestAC102_UnitNegative_OpaqueIDRejectedOnReuse(t *testing.T) {
	root := t.TempDir()
	l, _ := Load(root)
	const id = "AC-8f14e45f-ceea-467e-9a2b-a1c8b9d2f7a1"
	if err := l.ReserveOpaque(id, "rev1", "loc1"); err != nil {
		t.Fatalf("first ReserveOpaque: %v", err)
	}
	if err := l.ReserveOpaque(id, "rev2", "loc2"); err == nil {
		t.Fatal("second ReserveOpaque of the same id succeeded, want rejection")
	}
}

func TestAC105_UnitPositive_RetiredNumericIDNeverReissued(t *testing.T) {
	root := t.TempDir()
	l, _ := Load(root)
	// Simulate an archived id above the live maximum, as an
	// imported/backfilled entry would be: the counter sits just below it,
	// so the very next allocation would collide with it if the skip loop
	// did not fire.
	l.byID["PDR-050"] = &Entry{ID: "PDR-050", Kind: KindNumeric, State: StateRetired, Prefix: "PDR", Component: decimal(50)}
	l.counters["PDR"] = decimal(49)

	id := nextNumeric(t, l, "PDR")
	if id == "PDR-050" {
		t.Fatalf("NextNumeric reissued a retired id: %q", id)
	}
	if l.IsUsed("PDR-050") == false {
		t.Fatal("retired id must remain marked used")
	}
}

func TestAC105_UnitNegative_FreshlyIssuedNumberIsUnused(t *testing.T) {
	root := t.TempDir()
	l, _ := Load(root)
	l.byID["PDR-050"] = &Entry{ID: "PDR-050", Kind: KindNumeric, State: StateRetired, Prefix: "PDR", Component: decimal(50)}
	l.counters["PDR"] = decimal(49)
	id := nextNumeric(t, l, "PDR")
	if l.IsUsed(id) == false {
		t.Fatal("a freshly issued id must be recorded as used")
	}
	if e, ok := l.Lookup(id); !ok || e.State != StateReserved {
		t.Fatalf("freshly issued id state = %+v, ok=%v, want reserved", e, ok)
	}
}

func TestUnit_MarkLiveClassifiesShapeAndSeedsCounter(t *testing.T) {
	root := t.TempDir()
	l, _ := Load(root)
	l.MarkLive("CH-116")
	e, ok := l.Lookup("CH-116")
	if !ok || e.Kind != KindNumeric || e.Prefix != "CH" || e.Component.Cmp(decimal(116)) != 0 || e.State != StateLive {
		t.Fatalf("MarkLive(CH-116) entry = %+v, ok=%v", e, ok)
	}
	if l.counters["CH"].Cmp(decimal(116)) != 0 {
		t.Fatalf("counters[CH] = %d, want 116 seeded from the live component", l.counters["CH"])
	}

	l.MarkLive("SNAP-SQS-001")
	se, ok := l.Lookup("SNAP-SQS-001")
	if !ok || se.Kind != KindNumeric || se.Prefix != "SNAP-SQS" || se.Component.Cmp(decimal(1)) != 0 || se.State != StateLive {
		t.Fatalf("MarkLive(SNAP-SQS-001) entry = %+v, ok=%v", se, ok)
	}
	if l.counters["SNAP-SQS"].Cmp(decimal(1)) != 0 {
		t.Fatalf("counters[SNAP-SQS] = %d, want 1", l.counters["SNAP-SQS"])
	}
	if got := nextNumeric(t, l, "SNAP-SQS"); got != "SNAP-SQS-002" {
		t.Fatalf("next segmented ID = %q, want SNAP-SQS-002", got)
	}

	const opaque = "AC-8f14e45f-ceea-467e-9a2b-a1c8b9d2f7a1"
	l.MarkLive(opaque)
	oe, ok := l.Lookup(opaque)
	if !ok || oe.Kind != KindOpaque || oe.State != StateLive {
		t.Fatalf("MarkLive(%q) entry = %+v, ok=%v, want live opaque entry", opaque, oe, ok)
	}
}

func TestUnit_ValidNumericEntryRequiresCanonicalStructuredIdentity(t *testing.T) {
	if !ValidNumericEntry(Entry{ID: "ADP-045b", Kind: KindNumeric, Prefix: "ADP", Component: decimal(45)}) {
		t.Fatal("canonical suffixed numeric entry rejected")
	}
	if !ValidNumericEntry(Entry{ID: "SNAP-SQS-001", Kind: KindNumeric, Prefix: "SNAP-SQS", Component: decimal(1)}) {
		t.Fatal("canonical segmented numeric entry rejected")
	}
	if !ValidNumericEntry(Entry{ID: "AC-000", Kind: KindNumeric, Prefix: "AC", Component: decimal(0)}) {
		t.Fatal("canonical zero-valued numeric entry rejected")
	}
	for _, e := range []Entry{
		{ID: "G-001", Kind: KindNumeric, Prefix: "PDR", Component: decimal(1)},
		{ID: "G-001", Kind: KindNumeric, Prefix: "G", Component: decimal(99)},
		{ID: "snap-sqs-001", Kind: KindNumeric, Prefix: "snap-sqs", Component: decimal(1)},
	} {
		if ValidNumericEntry(e) {
			t.Fatalf("invalid numeric entry accepted: %+v", e)
		}
	}
}

func TestUnit_MarkLivePreservesArbitraryPrecisionComponent(t *testing.T) {
	root := t.TempDir()
	l, _ := Load(root)
	const component = "999999999999999999999999"
	l.MarkLive("AC-" + component)
	e, ok := l.Lookup("AC-" + component)
	if !ok || e.Kind != KindNumeric || e.Component.String() != component {
		t.Fatalf("large numeric entry = %+v, ok=%v", e, ok)
	}
	got := nextNumeric(t, l, "AC")
	if got != "AC-1000000000000000000000000" {
		t.Fatalf("next large numeric ID = %q", got)
	}
}

func TestUnit_PromoteReservedMarksAllocatedIDLive(t *testing.T) {
	root := t.TempDir()
	l, _ := Load(root)
	id := nextNumeric(t, l, "PDR")
	if err := l.PromoteReserved(id); err != nil {
		t.Fatalf("PromoteReserved: %v", err)
	}
	if e, ok := l.Lookup(id); !ok || e.State != StateLive {
		t.Fatalf("entry after promotion = %+v, ok=%v, want live", e, ok)
	}
	if err := l.PromoteReserved(id); err == nil {
		t.Fatal("second promotion succeeded, want non-reserved ID rejection")
	}
	if err := l.PromoteReserved("PDR-999"); err == nil {
		t.Fatal("promotion of an unreserved ID succeeded")
	}
}

func TestUnit_LoadMissingFileIsNotAnError(t *testing.T) {
	root := t.TempDir()
	l, err := Load(root)
	if err != nil {
		t.Fatalf("Load on missing ledger: %v", err)
	}
	if l.IsUsed("PDR-001") {
		t.Fatal("empty ledger must not report any id as used")
	}
	if Exists(root) {
		t.Fatal("Exists must be false before any Save")
	}
}

func TestAC104_UnitPositive_LoadRejectsDuplicateCanonicalID(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, ".clue")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	data := "counters: {}\nentries:\n  - id: PDR-001\n    kind: numeric\n    state: retired\n    prefix: PDR\n    component: 1\n  - id: PDR-001\n    kind: numeric\n    state: live\n    prefix: PDR\n    component: 1\n"
	if err := os.WriteFile(filepath.Join(dir, "id-ledger.yaml"), []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(root); err == nil {
		t.Fatal("Load accepted duplicate canonical ledger IDs")
	}
}

func TestUnit_SaveThenLoadRoundTrips(t *testing.T) {
	root := t.TempDir()
	l, _ := Load(root)
	nextNumeric(t, l, "PDR")
	if err := l.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if !Exists(root) {
		t.Fatal("Exists must be true after Save")
	}
	if _, err := os.Stat(filepath.Join(root, DefaultPath)); err != nil {
		t.Fatalf("ledger file not written: %v", err)
	}

	l2, err := Load(root)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if !l2.IsUsed("PDR-001") {
		t.Fatal("reloaded ledger lost its entry")
	}
}
