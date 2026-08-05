package ledger

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAC101_UnitPositive_NextNumericIncrementsStoredCounter(t *testing.T) {
	root := t.TempDir()
	l, err := Load(root)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	first := l.NextNumeric("PDR")
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
	second := l2.NextNumeric("PDR")
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
	id := l.NextNumeric("ZZZ")
	if id != "ZZZ-001" {
		t.Fatalf("id = %q, want ZZZ-001 for an unseen prefix", id)
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
	l.byID["PDR-050"] = &Entry{ID: "PDR-050", Kind: KindNumeric, State: StateRetired, Prefix: "PDR", Component: 50}
	l.counters["PDR"] = 49

	id := l.NextNumeric("PDR")
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
	l.byID["PDR-050"] = &Entry{ID: "PDR-050", Kind: KindNumeric, State: StateRetired, Prefix: "PDR", Component: 50}
	l.counters["PDR"] = 49
	id := l.NextNumeric("PDR")
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
	if !ok || e.Kind != KindNumeric || e.Prefix != "CH" || e.Component != 116 || e.State != StateLive {
		t.Fatalf("MarkLive(CH-116) entry = %+v, ok=%v", e, ok)
	}
	if l.counters["CH"] != 116 {
		t.Fatalf("counters[CH] = %d, want 116 seeded from the live component", l.counters["CH"])
	}

	const opaque = "AC-8f14e45f-ceea-467e-9a2b-a1c8b9d2f7a1"
	l.MarkLive(opaque)
	oe, ok := l.Lookup(opaque)
	if !ok || oe.Kind != KindOpaque || oe.State != StateLive {
		t.Fatalf("MarkLive(%q) entry = %+v, ok=%v, want live opaque entry", opaque, oe, ok)
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

func TestUnit_SaveThenLoadRoundTrips(t *testing.T) {
	root := t.TempDir()
	l, _ := Load(root)
	l.NextNumeric("PDR")
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
