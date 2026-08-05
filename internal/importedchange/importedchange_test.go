package importedchange

import (
	"reflect"
	"testing"
)

const fixtureBody = `# IC-001 — a fixture record

## Origin

Extracted from example-repo at pinned revision abc123.

## Intent

Carry the source change's design forward.

## Design rationale

Chose the smallest shape that keeps proof links inspectable.

## Dependencies

- IC-000

## Proof links

| Task | Criterion |
|---|---|
| write the allocator | AC-201 |
| write the validator rule | AC-202 |
`

func TestUnit_ParseProofLinksExtractsTaskAndCriterionColumns(t *testing.T) {
	got := ParseProofLinks(fixtureBody)
	want := []ProofLink{
		{Task: "write the allocator", Criterion: "AC-201"},
		{Task: "write the validator rule", Criterion: "AC-202"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestUnit_ParseProofLinksStopsAtTheNextHeading(t *testing.T) {
	body := fixtureBody + "\n## Something else\n\n| Task | Criterion |\n|---|---|\n| unrelated | AC-999 |\n"
	got := ParseProofLinks(body)
	for _, l := range got {
		if l.Criterion == "AC-999" {
			t.Fatalf("a later section's table leaked into the proof-links parse: %#v", got)
		}
	}
}

func TestUnit_ParseProofLinksNoHeadingReturnsNil(t *testing.T) {
	if got := ParseProofLinks("# just a title\n\nno proof-links section here\n"); got != nil {
		t.Fatalf("expected nil, got %#v", got)
	}
}

// AC-115: extraction preserves in-flight source work as an inspectable
// imported-change record — its proposal, design rationale, dependency, and
// proof-linked task all remain readable from the record alone.
func TestAC115_UnitPositive_RecordHoldsProposalDesignDependencyAndProofTask(t *testing.T) {
	links := ParseProofLinks(fixtureBody)
	if len(links) != 2 {
		t.Fatalf("expected 2 proof links, got %d: %#v", len(links), links)
	}
	for _, want := range []string{"Extracted from example-repo", "Carry the source change", "smallest shape", "IC-000"} {
		if !contains(fixtureBody, want) {
			t.Fatalf("fixture body missing expected inspectable content %q", want)
		}
	}
}

func TestAC115_UnitNegative_RecordWithNoProofLinksSectionHasNoInspectableProof(t *testing.T) {
	body := "# IC-002\n\n## Origin\n\nsomething\n"
	if links := ParseProofLinks(body); len(links) != 0 {
		t.Fatalf("expected no proof links, got %#v", links)
	}
}

// AC-116: a complete record's proof links must all be proven.
func TestAC116_UnitNegative_CompleteRecordWithUnprovenLinkIsRejected(t *testing.T) {
	links := ParseProofLinks(fixtureBody)
	provable := map[string]bool{"AC-201": true} // AC-202 missing: not provable
	bad := UnprovenLinks(links, provable)
	if len(bad) != 1 || bad[0] != "AC-202" {
		t.Fatalf("expected [AC-202], got %#v", bad)
	}
}

func TestAC116_UnitPositive_CompleteRecordWithAllLinksProvenPasses(t *testing.T) {
	links := ParseProofLinks(fixtureBody)
	provable := map[string]bool{"AC-201": true, "AC-202": true}
	if bad := UnprovenLinks(links, provable); len(bad) != 0 {
		t.Fatalf("expected no unproven links, got %#v", bad)
	}
}

// AC-117: an in-progress record is not held to the same standard — the
// caller (checkImportedChanges) must not call UnprovenLinks for that status,
// so this test proves the pure check itself is status-agnostic and the
// exemption is the caller's responsibility, matching ADR-050's design.
func TestAC117_UnitPositive_UnprovenLinksIsAPureCheckStatusIsTheCallersGate(t *testing.T) {
	links := ParseProofLinks(fixtureBody)
	provable := map[string]bool{} // nothing provable yet — an in-progress record's normal state
	bad := UnprovenLinks(links, provable)
	if len(bad) != 2 {
		t.Fatalf("expected both links to be reported unproven by this pure check, got %#v", bad)
	}
}

// The in-progress exemption is a boundary, not a blanket pass: the identical
// unproven link that AC-117's positive case allows is exactly what AC-116
// rejects once a record claims complete. Proving that boundary from the
// complete side is this criterion's negative case.
func TestAC117_UnitNegative_TheSameUnprovenLinkFailsOnceStatusClaimsComplete(t *testing.T) {
	links := ParseProofLinks(fixtureBody)
	provable := map[string]bool{"AC-201": true} // AC-202 still unproven
	bad := UnprovenLinks(links, provable)
	if len(bad) == 0 {
		t.Fatalf("expected the unproven link to be reported once the caller applies it to a complete record, got none")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (func() bool {
		for i := 0; i+len(substr) <= len(s); i++ {
			if s[i:i+len(substr)] == substr {
				return true
			}
		}
		return false
	})()
}
