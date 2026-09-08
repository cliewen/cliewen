package ledger

import (
	"math/big"
	"os"
	"path/filepath"
	"strings"
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

func TestAC174_UnitPositive_NextNumericIncrementsStoredHighWater(t *testing.T) {
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

func TestAC174_UnitNegative_PrefixAbsentFromLedgerStartsAtOne(t *testing.T) {
	root := t.TempDir()
	l, _ := Load(root)
	id := nextNumeric(t, l, "ZZZ")
	if id != "ZZZ-001" {
		t.Fatalf("id = %q, want ZZZ-001 for an unseen prefix", id)
	}
}

func TestAC174_UnitNegative_NextNumericRejectsNonCanonicalPrefix(t *testing.T) {
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

func TestAC104_UnitPositive_LoadRejectsNullOrNegativeCounter(t *testing.T) {
	for _, counters := range []string{"PDR: null", "PDR: -1"} {
		root := t.TempDir()
		dir := filepath.Join(root, ".clue")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "id-ledger.yaml"), []byte("counters: {"+counters+"}\nentries: []\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, err := Load(root); err == nil {
			t.Fatalf("Load accepted malformed counter %q", counters)
		}
	}
}

func TestAC104_UnitPositive_LoadRejectsMissingOrLowNumericCounter(t *testing.T) {
	for _, counters := range []string{"{}", "{PDR: 1}"} {
		root := t.TempDir()
		dir := filepath.Join(root, ".clue")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		data := "counters: " + counters + "\nentries:\n  - id: PDR-002\n    kind: numeric\n    state: retired\n    prefix: PDR\n    component: 2\n"
		if err := os.WriteFile(filepath.Join(dir, "id-ledger.yaml"), []byte(data), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, err := Load(root); err == nil {
			t.Fatalf("Load accepted insufficient counters %s", counters)
		}
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

func TestAC169_UnitPositive_EventLedgerFoldsLifecycleMonotonically(t *testing.T) {
	root := t.TempDir()
	l, _ := Load(root)
	id := nextNumeric(t, l, "CH")
	if err := l.PromoteReserved(id); err != nil {
		t.Fatal(err)
	}
	l.Retire(id)
	if err := l.Save(); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(root, DefaultPath))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "version: 2") || strings.Count(string(data), "id: CH-001") != 3 {
		t.Fatalf("ledger did not retain one event per transition:\n%s", data)
	}
	reloaded, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if e, ok := reloaded.Lookup(id); !ok || e.State != StateRetired {
		t.Fatalf("effective entry = %+v, ok=%v", e, ok)
	}
}

func TestAC169_UnitNegative_EventLedgerRejectsConflictingIdentityMetadata(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".clue"), 0o755); err != nil {
		t.Fatal(err)
	}
	data := "version: 2\ncoordination:\n    mode: local\nevents:\n    - {id: CH-001, kind: numeric, state: reserved, prefix: CH, component: \"1\", source-revision: one}\n    - {id: CH-001, kind: numeric, state: live, prefix: CH, component: \"1\", source-revision: two}\n"
	if err := os.WriteFile(filepath.Join(root, DefaultPath), []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(root); err == nil || !strings.Contains(err.Error(), "conflicting identity metadata") {
		t.Fatalf("Load error = %v", err)
	}
}

func TestAC169_UnitNegative_EventLedgerRejectsUnknownSchemaVersion(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".clue"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, DefaultPath), []byte("version: 3\ncoordination: {mode: local}\nevents: []\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(root); err == nil || !strings.Contains(err.Error(), "unsupported ledger version 3") {
		t.Fatalf("Load error = %v", err)
	}
}

func TestAC170_UnitPositive_MergeClaimsIsIdempotentAndNeverDowngrades(t *testing.T) {
	root := t.TempDir()
	l, _ := Load(root)
	id := nextNumeric(t, l, "CH")
	if err := l.PromoteReserved(id); err != nil {
		t.Fatal(err)
	}
	claim := Entry{ID: id, Kind: KindNumeric, State: StateReserved, Prefix: "CH", Component: decimal(1)}
	if err := l.MergeClaims([]Entry{claim, claim}); err != nil {
		t.Fatal(err)
	}
	if e, _ := l.Lookup(id); e.State != StateLive {
		t.Fatalf("remote reservation downgraded local state to %s", e.State)
	}
	if err := l.Save(); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(filepath.Join(root, DefaultPath))
	if strings.Count(string(data), "id: CH-001") != 2 {
		t.Fatalf("idempotent claim merge appended duplicates:\n%s", data)
	}
}

func TestAC170_UnitNegative_MergeClaimsRejectsConflictingMetadata(t *testing.T) {
	root := t.TempDir()
	l, _ := Load(root)
	first := Entry{ID: "CH-001", Kind: KindNumeric, State: StateReserved, Prefix: "CH", Component: decimal(1), SourceRevision: "one"}
	second := cloneEntry(first)
	second.SourceRevision = "two"
	if err := l.MergeClaims([]Entry{first}); err != nil {
		t.Fatal(err)
	}
	if err := l.MergeClaims([]Entry{second}); err == nil || !strings.Contains(err.Error(), "conflicting identity metadata") {
		t.Fatalf("MergeClaims error = %v", err)
	}
}

// writeLedger puts raw ledger bytes at root, for the shapes a merge produces
// that no API of ours would ever write.
func writeLedger(t *testing.T, root, data string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(root, ".clue"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, DefaultPath), []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
}

const halfA = "version: 2\ncoordination:\n    mode: local\nevents:\n    - {id: CH-001, kind: numeric, state: reserved, prefix: CH, component: \"1\"}\n    - {id: CH-001, kind: numeric, state: live, prefix: CH, component: \"1\"}\n"

const halfB = "version: 2\ncoordination:\n    mode: local\nevents:\n    - {id: CH-002, kind: numeric, state: reserved, prefix: CH, component: \"2\"}\n"

func TestAC179_UnitPositive_UnionMergedLedgerLoadsAndFoldsBothHalves(t *testing.T) {
	root := t.TempDir()
	writeLedger(t, root, halfA+halfB)

	l, err := Load(root)
	if err != nil {
		t.Fatalf("Load of a union-merged ledger failed: %v", err)
	}
	if l.Damage() == "" {
		t.Fatal("Damage() is empty; the repeated header was not reported")
	}
	for _, want := range []string{"version", "coordination", "events"} {
		if !strings.Contains(l.Damage(), want) {
			t.Fatalf("Damage() = %q, want it to name %q", l.Damage(), want)
		}
	}
	// Both branches' identities survive, at their furthest-along state.
	if e, ok := l.Lookup("CH-001"); !ok || e.State != StateLive {
		t.Fatalf("CH-001 = %+v, ok=%v; want live", e, ok)
	}
	if e, ok := l.Lookup("CH-002"); !ok || e.State != StateReserved {
		t.Fatalf("CH-002 = %+v, ok=%v; want reserved", e, ok)
	}
	// The counter reflects the higher half, so nothing is reissued.
	if id := nextNumeric(t, l, "CH"); id != "CH-003" {
		t.Fatalf("next after recovery = %s, want CH-003", id)
	}
}

func TestAC179_UnitPositive_SavingARecoveredLedgerRewritesItClean(t *testing.T) {
	root := t.TempDir()
	writeLedger(t, root, halfA+halfB)

	l, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if err := l.Save(); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(root, DefaultPath))
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Count(string(data), "version: 2"); got != 1 {
		t.Fatalf("repaired ledger declares version %d times:\n%s", got, data)
	}
	if got := strings.Count(string(data), "coordination:"); got != 1 {
		t.Fatalf("repaired ledger declares coordination %d times:\n%s", got, data)
	}
	reloaded, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if reloaded.Damage() != "" {
		t.Fatalf("reloaded Damage() = %q, want empty", reloaded.Damage())
	}
	if len(reloaded.Entries()) != 2 {
		t.Fatalf("repaired ledger kept %d identities, want 2", len(reloaded.Entries()))
	}
}

// The one case recovery must refuse: two halves that chose different
// allocation modes. Picking either would silently move a team off the mode it
// agreed, so the ledger fails closed and says what the choice is.
func TestAC179_UnitNegative_DisagreeingCoordinationFailsClosed(t *testing.T) {
	root := t.TempDir()
	writeLedger(t, root, halfA+"version: 2\ncoordination:\n    mode: git\n    remote: origin\nevents: []\n")

	_, err := Load(root)
	if err == nil {
		t.Fatal("Load accepted two disagreeing coordination settings")
	}
	for _, want := range []string{"two different coordination settings", "local", "git through origin"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("Load error = %q, want it to mention %q", err, want)
		}
	}
}

func TestAC179_UnitNegative_UndamagedLedgerReportsNoDamage(t *testing.T) {
	root := t.TempDir()
	writeLedger(t, root, halfA)

	l, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if l.Damage() != "" {
		t.Fatalf("Damage() = %q for a well-formed ledger, want empty", l.Damage())
	}
}

// Backfilling an identity that was retired before the ledger existed must not
// invent a live event it never had.
func TestUnit_MarkRetiredRecordsOneEventForAnIdentityNeverLive(t *testing.T) {
	root := t.TempDir()
	l, _ := Load(root)
	l.MarkRetired("AC-042")
	l.MarkRetired("AC-042")
	if err := l.Save(); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(root, DefaultPath))
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Count(string(data), "id: AC-042"); got != 1 {
		t.Fatalf("retiring an unseen ID wrote %d events, want 1:\n%s", got, data)
	}
	if strings.Contains(string(data), "state: live") {
		t.Fatalf("retirement invented a live transition:\n%s", data)
	}
	if e, ok := l.Lookup("AC-042"); !ok || e.State != StateRetired {
		t.Fatalf("AC-042 = %+v, ok=%v; want retired", e, ok)
	}
}

func TestUnit_RetireIsIdempotentAndNeverDowngrades(t *testing.T) {
	root := t.TempDir()
	l, _ := Load(root)
	id := nextNumeric(t, l, "CH")
	if err := l.PromoteReserved(id); err != nil {
		t.Fatal(err)
	}
	l.Retire(id)
	l.Retire(id)
	l.MarkLive(id)
	if err := l.Save(); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(root, DefaultPath))
	if err != nil {
		t.Fatal(err)
	}
	// reserved, live, retired — and nothing after it.
	if got := strings.Count(string(data), "id: "+id); got != 3 {
		t.Fatalf("repeated retirement wrote %d events, want 3:\n%s", got, data)
	}
	if e, _ := l.Lookup(id); e.State != StateRetired {
		t.Fatalf("%s = %s, want retired", id, e.State)
	}
}

func TestUnit_UnionMergeRuleIsRecognizedHoweverItIsWritten(t *testing.T) {
	for name, line := range map[string]string{
		"canonical":         UnionAttribute,
		"no leading slash":  ".clue/id-ledger.yaml merge=union",
		"extra attributes":  "/.clue/id-ledger.yaml merge=union -text",
		"among other rules": "* text=auto\n# the ledger\n.clue/id-ledger.yaml merge=union\n",
	} {
		if !declaresUnionMerge([]byte(line)) {
			t.Errorf("%s: %q not recognized as declaring the union rule", name, line)
		}
	}
}

func TestUnit_UnrelatedAttributeLinesDoNotDeclareTheUnionRule(t *testing.T) {
	for name, line := range map[string]string{
		"another file":   "/.clue/role.yaml merge=union",
		"another driver": "/.clue/id-ledger.yaml merge=ours",
		"commented out":  "#/.clue/id-ledger.yaml merge=union",
		"pattern only":   "/.clue/id-ledger.yaml",
		"nothing at all": "* text=auto",
	} {
		if declaresUnionMerge([]byte(line)) {
			t.Errorf("%s: %q wrongly recognized as declaring the union rule", name, line)
		}
	}
}

// The damage a real union merge produces. Two branches each coordinating to
// their own remote leave `version` and `mode` untouched — identical lines merge
// as context — and duplicate only `remote`, inside the coordination block. A
// struct decode rejected that with a parser message about a line nobody wrote.
func TestAC179_UnitNegative_DuplicateRemoteInsideCoordinationFailsClosed(t *testing.T) {
	root := t.TempDir()
	writeLedger(t, root, "version: 2\ncoordination:\n    mode: git\n    remote: origin\n    remote: other\nevents:\n    - {id: CH-001, kind: numeric, state: live, prefix: CH, component: \"1\"}\n")

	_, err := Load(root)
	if err == nil {
		t.Fatal("Load accepted a ledger naming two allocator remotes")
	}
	for _, want := range []string{"two different allocator remotes", "origin", "other", "only a person can decide"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("Load error = %q, want it to mention %q", err, want)
		}
	}
	if strings.Contains(err.Error(), "already defined") {
		t.Fatalf("Load error = %q, want a decidable message rather than a parser one", err)
	}
}

// A key repeated with the same value is not a disagreement, so it reconciles
// rather than stopping the repository.
func TestAC179_UnitPositive_DuplicateCoordinationKeyWithOneValueReconciles(t *testing.T) {
	root := t.TempDir()
	writeLedger(t, root, "version: 2\ncoordination:\n    mode: git\n    remote: origin\n    remote: origin\nevents:\n    - {id: CH-001, kind: numeric, state: live, prefix: CH, component: \"1\"}\n")

	l, err := Load(root)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if c := l.Coordination(); c.Mode != "git" || c.Remote != "origin" {
		t.Fatalf("coordination = %+v, want git through origin", c)
	}
}
