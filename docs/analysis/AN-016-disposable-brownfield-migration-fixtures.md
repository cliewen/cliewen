---
id: AN-016
type: analysis
status: active
provenance: verified
reversal-cost: low
links: [P-011, CAP-003, PDR-020, ADR-048, ADR-049, ADR-050, ADR-051]
title: Disposable fixtures prove the composed brownfield migration contract
---

# AN-016 — Disposable fixtures prove the composed brownfield migration contract

## Scope and pins

CH-120 exercised P-011/M-056 with two disposable OpenSpec-shaped sources rather than a production adopter. `numeric-archive-fixture-v1` at `fixtures/numeric-archive` contributes live criterion `ARC-002`, classified positive and negative source proof, pending work, and retired ledger identity `ARC-099`. `opaque-identifier-fixture-v1` at `fixtures/opaque-identifier` contributes `OPA-001`, classified proof, pending work, and source-owned opaque identity `8f14e45f-ceea-467e-9a2b-a1c8b9d2f7a1`.

The initial report-only rehearsal recorded the source pins, manifests, proposed mappings, and carrier inventories before any fixture target existed. The maintainer then explicitly authorized mutation in conversation on 2026-08-05, satisfying PDR-020's checkpoint before the target fixtures were materialized.

## Mappings and durable work

Each source maps to an isolated target corpus with its own capability, criteria, ledger entries, and complete `imported-change` record. The records retain source revision and location, intent, design rationale, dependencies, and proof-link tables after the fixture source is removed from the target view.

The numeric target retains `ARC-099` as retired while `ARC-002` is live; allocation advances to `ARC-100`, so a removed numeric identity cannot return. The opaque target reserves its UUID-like identity verbatim and rejects a second reservation. Neither fixture treats an opaque source token as a criterion identifier.

## Parity, carriers, and failure boundaries

The fixture writes and loads pinned source manifests and finalized carrier inventories, then invokes `clue validate`, `clue parity`, and `clue carriers` against each target. Both clean runs pass. Each inventory maps the rehearsal's instruction, registry, and link carriers to a retained target with a current fingerprint.

Deliberate variants prove that `clue parity` rejects missing criteria, orphaned evidence, changed proof direction or location, stale source revisions, and unjustified dispositions. Other variants prove that `clue carriers` rejects stale deleted-path links, lost fingerprints, and missing assets. The test invokes the public commands, so manifest loading, command routing, reports, and non-zero exits are part of the evidence instead of only package internals.

The source fixture's own tests are never executed. Their locations enter the source-side manifest only; Cliewen acceptance evidence is the target validation and deterministic command results.

## Result

The composed contract holds for a disposable numeric archive and an opaque source-owned identifier, including the report-before-approval-before-mutation ordering. This is fixture evidence for the common migration contract, not a claim about any production adopter's corpus, tests, or CI.

## Rejected alternatives

Treating the focused component tests as complete migration evidence was rejected because individually valid ledger, parity, imported-work, and carrier components may still fail to compose.

Using a production adopter was rejected because the contract needs repeatable fixture evidence and must not make an adopter's operational state into this repository's test fixture.
