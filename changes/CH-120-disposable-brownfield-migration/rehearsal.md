---
id: CH-120-rehearsal
type: analysis
status: draft
links: [CH-120, P-011, CAP-003, PDR-020, ADR-048, ADR-049, ADR-050, ADR-051]
title: Report-only rehearsal for disposable brownfield migration fixtures
provenance: inferred
reversal-cost: low
---

# CH-120 — report-only rehearsal for disposable brownfield migration fixtures

## Scope and boundary

This rehearsal exercises P-011/M-056 with two disposable OpenSpec-shaped source fixtures. It writes only the transient CH-120 workspace. It does not create a target corpus, alter Cliewen documentation, tests, routing, or hosted state, execute either source fixture's own test suite, or delete source material.

The fixtures are intentionally invented test data rather than a production adopter. Their value is the repeatable migration shape: a source revision, a source location, an identity inventory, classified proof locations, in-flight work, operational carriers, and the failure paths the final target must reject.

## Pins and source entry points

`numeric-archive` is pinned as revision `numeric-archive-fixture-v1` at `fixtures/numeric-archive`. It contains an OpenSpec capability with live criterion ID `ARC-100`, source tests with positive and negative proof, a pending source change, and an index and instruction carrier; archived identity `ARC-099` remains reserved in the target ledger.

`opaque-identifier` is pinned as revision `opaque-identifier-fixture-v1` at `fixtures/opaque-identifier`. It contains an OpenSpec capability with criterion `OPA-001`, classified positive and negative source proof, a pending source change, and an opaque imported source identity `8f14e45f-ceea-467e-9a2b-a1c8b9d2f7a1`. The opaque ID is not an invented criterion grammar; it is a preserved source-owned artifact identity that the target ledger must reserve verbatim and refuse to reuse.

The pinned source manifests are [numeric-source-manifest.yaml](numeric-source-manifest.yaml) and [opaque-source-manifest.yaml](opaque-source-manifest.yaml). They record every source criterion proof at its source location. The manifests are source-side evidence, not editable target coverage reports.

## Proposed target mappings

The numeric fixture will map its capability and criteria into a disposable target corpus, preserve `ARC-099` as a retired ledger identity, and demonstrate that the numeric allocator does not reissue it. Its pending source change will become an `imported-change` record whose proposal, rationale, dependency, and proof links remain inspectable after the source tree is removed.

The opaque fixture will map its capability and criteria into a separate disposable target corpus, preserve `OPA-001` through the normal criterion grammar, and reserve the opaque source identity exactly as supplied by its source-owned generator. The fixture will attempt a duplicate reservation to prove that retired, reserved, and live opaque identities are equally unavailable for reuse.

Both targets will derive their parity side from the corpus and ledger, never from authored coverage data. The mutation phase will demonstrate a clean parity run and each required parity failure class by deliberate fixture variants. It will also demonstrate target validation and a clean carrier reconciliation run, then deliberate stale deleted-path, lost-fingerprint, and missing-asset variants.

## Operational-carrier inventory and deletions

The source trees will eventually delete their OpenSpec specifications, pending-change folders, source test registries, index files, and OpenSpec-specific instruction files. The rehearsal inventories the source carriers in [numeric-carriers.yaml](numeric-carriers.yaml) and [opaque-carriers.yaml](opaque-carriers.yaml).

Every listed carrier is currently `blocked` because report-only work has no approved target file to fingerprint. This is intentional: inventing a target path or hash before mutation would make the inventory claim a retained carrier that does not exist. On authorized mutation, each entry must be revised into a mapped target path with a current fingerprint before source deletion can proceed.

## Evidence and test-purpose work

The source manifests declare `Integration` proof and both directions for `ARC-099` and `OPA-001`. The future target fixtures must attach one criterion, proof type, and direction to each executable evidence carrier, then preserve the exact target locations in the derived parity manifest.

The source fixtures' own tests remain source evidence only. Their successful or failed execution is not evidence for Cliewen's deterministic commands. The CH-120 Go tests will separately prove target validation, ledger handling, parity, and carrier reconciliation.

## Confidence, risks, and named doors

Confidence is high that the existing M-052 through M-055 components can be composed in isolated fixtures because each component already has focused acceptance evidence. The remaining risk is integration drift: a source manifest, ledger entry, imported-change record, and carrier inventory may use individually valid representations that do not compose into a clean end-to-end target.

No new source-format mapping is required because both fixtures use the canonical OpenSpec shape. No semantic conflict was found in the source evidence. Mutation nevertheless requires the explicit authorization recorded in [open-questions.md](open-questions.md), as PDR-020 requires.

## Planned failure paths

- A source manifest entry omitted from the target must report a missing criterion.
- Added target evidence must report an orphaned tag, and changed direction or evidence location must report changed evidence.
- A source revision that differs from its ledger record must report a stale fingerprint.
- A target draft, Human, or retired state without its matching source disposition must report an unjustified disposition.
- A retained carrier that drifts, is absent, or remains linked from a deleted source path must report the applicable carrier finding.
- An imported change whose proof link is not complete must block source-work deletion, and a numeric or opaque identity that already exists in the ledger must reject reuse.

## Rejected alternatives

Treating the focused M-052 through M-055 tests as the M-056 evidence was rejected because they do not prove their manifests, ledger state, imported work, and carrier inventory compose in one migration.

Using a production adopter was rejected because M-056 needs a disposable, repeatable contract fixture and must not turn a real repository's operational state into this repository's test data.

Running source fixture tests during this rehearsal was rejected because the report-only phase must not let source-suite execution masquerade as Cliewen validation or mutate source build state.
