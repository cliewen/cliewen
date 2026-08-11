---
id: ADR-049
type: decision
status: verified
links: [ADR-008, ADR-032, ADR-033, ADR-048, PDR-020, PDR-024, P-011, C-013]
title: A pinned source manifest and a derived target manifest give migration parity one comparable shape
author: agent
accepted-by: Flemming N. Larsen (2026-08-08, conversation)
---

# ADR-049 — Migration parity compares a pinned source manifest against a derived target manifest

## Context and problem statement

[CAP-003](../capabilities/CAP-003-extract/README.md)'s extraction preserves canonical IDs, evidence classification, and provenance, and [PDR-024](PDR-024-extraction-never-silently-discards-evidence.md) already forbids silently discarding, selecting, or relocating source evidence during the rehearsal's human-reviewed resolution. Nothing today re-checks that resolution once mutation lands in a migration PR: a mapping can drop a criterion between rehearsal and merge, an evidence reference can move without a trace, or a `@draft`/`Human`/retired disposition can go unrecorded, and `clue validate` — which only reads the target corpus — has no source-side state to compare against. [P-011](../plans/P-011-truthful-brownfield-migration.md)'s M-053 commits to a "deterministic parity command" closing that gap. What must each side of the comparison hold, and what makes an entry pass or fail?

## Decision outcome

**A source manifest is a pinned, human/agent-authored file; a target manifest is deterministically derived by `clue` from the current corpus and ledger. `clue parity` compares them and reports every unmatched or altered entry.**

- **Source manifest** (YAML, one per source-mapping run, written during the `clue-extract` rehearsal and revised only when the source is re-read at a new revision): top-level `source-revision` and `source-location` pin what was read; each `entries[]` row names the criterion `id` and either:
  - one proof class (`proof-class`, the `Test-type` vocabulary), `direction` (`positive | negative | single-direction`), and `evidence-location` (a file:line or equivalent pointer into the *source*); a criterion has one such row for every classified source reference, so parity compares complete direction and location sets, or
  - `excluded: true` with a `reason` (the criterion is deliberately not carried forward), or
  - a `disposition` of `draft | human | retired` with a `justification`; [ADR-053](ADR-053-deferred-parity-dispositions-are-accountable.md) adds the separate `disposition-source-location` and `plan-door` fields, so both are inspectable and each door resolves uniquely in the target corpus. An entry may not combine `excluded` with a proof class, nor a proof class with a `disposition` — each row states exactly one outcome; an exclusion or disposition is the sole row for its ID.
- **Target manifest** is never authored: `clue parity` derives it from the same declaration-and-evidence harvest `clue validate` already runs (`internal/corpus`), reusing its classified-evidence walk rather than re-parsing the tree a second way. Per live or tombstoned criterion it holds the ledger state (`live | reserved | retired`, cross-checked against [ADR-048](ADR-048-corpus-wide-id-ledger.md)'s ledger when one exists), proof class, every classified direction, every evidence location recorded for it, and whether it is `@draft`, `Human`-classed, or retired.
- **`Compare`** matches source and target entries by ID and reports exactly these failure classes, matching M-053's stated exit criterion:
  1. **Missing criterion** — a non-excluded source entry has no target entry at all.
  2. **Orphaned tag** — the target carries classified evidence for an ID absent from the source manifest entirely (neither present nor excluded), so nothing pinned accounts for it.
  3. **Changed direction or location** — a source entry with a proof class has a matching target entry, but the target's direction set or evidence location(s) disagree with what the source recorded.
  4. **Stale source fingerprint** — the source manifest's `source-revision` disagrees with the revision already recorded for that ID's ledger entry (`ADR-048`'s `source-revision` field), meaning the manifest was not regenerated against what the corpus actually pins.
  5. **Unjustified disposition** — the target shows `@draft`, `Human`, or retired for an ID whose source entry has no matching `disposition` and `justification`. A clean comparison — every source entry matched, no orphaned target evidence, no disagreement — is the only passing result; invalid or ambiguous manifest rows are usage errors, and the report is deterministic (stable ordering, no wall-clock or environment-dependent content) so a CI artifact from the same inputs is byte-identical.
- **The report stays derived.** `clue parity` never writes back into the source manifest or the target corpus; a failing run is repaired by fixing the mapping, the corpus, or the manifest's own pinned fields, then re-running, never by hand-editing a stored coverage table.

## Rejected: one manifest format, hand-merged for both sides

Making an agent write the "target" half by hand as well would let a migration PR assert coverage that the corpus does not actually contain — exactly the unchecked-claim problem this ADR closes. Deriving the target deterministically from the same harvest `clue validate` trusts keeps the comparison honest against one authority.

## Rejected: compare against Git history instead of a pinned manifest

P-011's own scope excludes treating Git history as an implicit registry (see the plan's "Explicitly out of this campaign"). A pinned manifest names its source revision and location explicitly and survives the source repository's deletion; a history-based comparison would not.

## Rejected: fold parity into `clue validate`

`clue validate` judges one corpus against its own declared rules and needs no external input. Parity is a two-sided comparison that requires a source manifest nothing else in the corpus produces; conflating the two would make `validate`'s exit code depend on a file most repositories never have, which ADR-048's ledger gating already treats as the wrong default for a machine-checked corpus rule.

## Carrier

This record, `internal/parity/`, `docs/capabilities/CAP-003-extract/criteria.md` (AC-109..AC-114), and the canonical `clue-extract` skill's rehearsal guidance are the carriers of this decision.
