---
id: ADR-051
type: decision
status: verified
links: [ADR-008, ADR-048, ADR-049, PDR-019, PDR-025, P-011, C-013]
title: A pinned carrier inventory and a deterministic reconciliation check close the operational-carrier gap
author: agent
accepted-by: Flemming N. Larsen (2026-08-08, conversation)
---

# ADR-051 — A pinned carrier inventory and a deterministic reconciliation check close the operational-carrier gap

## Context and problem statement

[ADR-048](ADR-048-corpus-wide-id-ledger.md) closes the identity gap and [ADR-049](ADR-049-migration-parity-manifests.md) closes the evidence-parity gap for a brownfield migration, but neither reaches a third class of thing a source repository depends on: material that is not a criterion or a piece of evidence at all, yet still governs how that repository actually operates — its CI workflows, its own agent instructions, whatever it reads to decide something is current (a version pin, a last-checked date), its cross-reference registries and indexes, its local and external Markdown links, and its diagram assets.

[PDR-025](PDR-025-brownfield-migration-precedes-simplification.md) names this directly: "the current extraction contract can lose ... operational carriers." [PDR-019](PDR-019-methodology-contract-carriers-move-together.md) already defines "live carrier" for Cliewen's own corpus; a migrated source's operational surface is the same concept applied to a different repository.

Today the `clue-extract` OpenSpec mapping handles pieces of this ad hoc — a hand-maintained-index row, a "preserve every link" rule, a diagram-choice rule per C-007 — with no inventory step and no gate: a mapping can silently drop a CI workflow, break a freshness check, or leave a dangling link, and neither `clue validate` nor `clue parity` would notice, because none of that is a criterion.

[P-011](../plans/P-011-truthful-brownfield-migration.md)'s M-055 commits to a systematic inventory and a CI gate. What must the inventory hold, and what makes reconciliation pass or fail?

## Decision outcome

**A carrier inventory is a pinned, human/agent-authored manifest, written during the `clue-extract` rehearsal like ADR-049's source manifest; `clue carriers` reconciles it against the current corpus and reports every finding, mirroring ADR-049's split between an authored source side and a derived comparison.**

- **Operational carrier**, for a migrated source, is any of six kinds: `instruction` (an agent-facing entry file such as an OpenSpec-side `AGENTS.md` equivalent), `workflow` (a CI definition, e.g. under `.github/workflows/`), `freshness-input` (anything a CI job or script reads to decide something is current — a version pin, a last-checked timestamp, a lockfile), `registry` (a cross-reference index or changelog), `link` (a local or external Markdown link the source corpus carries), and `diagram-asset` (an SVG, embedded Mermaid, or ASCII diagram the source references). This list is the inventory's fixed vocabulary; a new migrated kind extends it through a plan decision, the same route ADR-049's proof-class vocabulary uses.
- **Carrier inventory** (YAML, one per source-mapping run, revised only when the source is re-read at a new revision): top-level `source-revision` and `source-location` pin what was read, reusing ADR-048/ADR-049's field names for the same purpose. `deleted-paths[]` names every source-repository path the migration will delete, so reconciliation can recognize a reference left pointing at one. `entries[]` rows each name one carrier: an `id`, its `kind`, its `source-path`, and either — a `target-path` naming its retained location in the reconciled corpus plus a `fingerprint` (a content hash) of that target's state at inventory time, or `blocked: true` with a `reason` (an explicit block-mutation marker: this carrier has no target yet, and its presence in the inventory is what stops the source deletion it would otherwise silently survive as nothing). An entry may not combine `target-path`/`fingerprint` with `blocked`; every carrier the rehearsal finds gets exactly one outcome, so an inventory that names a carrier with neither is rejected at load, matching M-055's "every carrier maps to a retained target or blocks mutation."
- **Reconciliation is derived, never authored:** `clue carriers <inventory> [root]` recomputes each mapped entry's target fingerprint from what is actually in `root` right now and compares it against the pinned value, the same "pin, then re-derive, then compare" shape ADR-049 uses for evidence rather than authoring both sides by hand.
- **`Reconcile`** reports exactly three failure classes, matching M-055's stated exit criterion:
  1. **Stale deleted-path reference** — a local or external Markdown link anywhere in the reconciled corpus still resolves to a path named in `deleted-paths`, meaning something the migration is about to delete is still depended on.
  2. **Lost fingerprint** — a mapped entry's `target-path` exists, but its current content fingerprint disagrees with the one the inventory pinned, meaning the retained carrier drifted after the rehearsal recorded it as current and was never re-pinned.
  3. **Missing asset** — a mapped entry's `target-path` does not exist in the reconciled corpus at all.

  A clean run — every mapped fingerprint matches, no deleted path is still referenced, no target is missing — is the only passing result. `blocked` entries are not reconciled against a target (there is none yet); their presence is itself the record that the carrier is a known gap, not a silent one, and closing one converts it to a mapped entry in a later inventory revision. The report is sorted and holds no wall-clock or environment-dependent content, matching ADR-049's reproducibility requirement, so a CI artifact from the same inputs is byte-identical.
- **`clue carriers` never writes back** into the inventory or the target corpus; a failing run is repaired by fixing the mapping, the corpus, or the inventory's own pinned fields, then re-running.
- **Positioning matches `clue parity`, not `clue validate`:** reconciliation needs a source-side inventory nothing else in an ordinary corpus produces, so it stays a separate CLI entry point excluded from the ambient release notice (a deterministic judge's output must not depend on another system's present state) rather than a `clue validate` rule. A migrating repository's own CI workflow is what makes the check required before a migration PR can merge — the same route `clue parity` already takes, per the `clue-extract` skill's target contract, not a rule this repository's own `ci.yml` runs for itself.
- **No source-format parsing in `clue`, ever:** what counts as a carrier for a given source format, and how it is discovered during the rehearsal without executing that source's own CI or requiring live network access, stays a `clue-extract` mapping concern — the same split ADR-049 already draws for the source manifest.

## Rejected: fold carrier reconciliation into `clue parity`

Parity's shape is criterion-keyed: one comparable ID per row, evidence classified by proof type and direction. A carrier is not a criterion and has no proof-class or direction; forcing it through parity's schema would either invent a fake criterion ID for a CI workflow file or bend parity's finding classes to mean something else for one column, weakening a contract ADR-049 already made precise for evidence.

## Rejected: infer carriers from source-format knowledge inside `clue`

Teaching the deterministic judge what an OpenSpec-side `AGENTS.md` equivalent or a GitHub Actions workflow looks like would put source-format parsing in `clue`, which every prior extraction decision (ADR-008, ADR-049) deliberately keeps out. The rehearsal, executed by an agent that already reads the source, is where that judgment belongs; `clue` only reconciles what the rehearsal recorded.

## Rejected: treat a missing fingerprint match as a warning, not a failure

An operational carrier that drifted silently is exactly the failure mode PDR-025 names. Making it advisory would let a migration PR merge on the same unchecked claim the parity contract already refuses for evidence.

## Carrier

This record, `internal/carriers/`, `docs/capabilities/CAP-003-extract/criteria.md` (AC-118..AC-122 — renumbered from the M-055 milestone row's original AC-115..AC-119 range because CH-118/M-054 merged into `main` first and already claims AC-115..AC-117), and the canonical `clue-extract` skill's rehearsal guidance and OpenSpec mapping are the carriers of this decision.
