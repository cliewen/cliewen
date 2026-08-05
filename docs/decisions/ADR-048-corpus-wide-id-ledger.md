---
id: ADR-048
type: decision
status: inferred
links: [ADR-009, ADR-034, ADR-007, P-011, C-013]
title: A persisted ledger replaces scan-and-max allocation for every native ID prefix
author: agent
accepted-by: []
---

# ADR-048 — A persisted ledger replaces scan-and-max allocation for every native ID prefix

## Context and problem statement

[ADR-009](ADR-009-ac-id-namespaces.md) states the current allocator for AC IDs — "the corpus is the registry: next-free-ID per prefix is the next numeric slot after the maximum numeric component already declared" — and the same scan-and-max practice is followed by convention, uncodified, for every other native prefix (`CH`, `C`, `P`, `M`, `PDR`, `ADR`, `CAP`, `G`, `AN`, ...). [ADR-034](ADR-034-retirement-is-deletion.md) makes retirement of any non-criterion artifact a file deletion, with no committed `status: retired` state surviving on `main`. Combined, these two decisions leave a gap: once a file is deleted, its ID is invisible to any future scan, so nothing stops that number being re-minted later. A re-minted ID silently collides with historic evidence, commit history, or external references that still carry the deleted artifact's original meaning — a correctness gap in identity, not a cosmetic one.

[P-011](../plans/P-011-truthful-brownfield-migration.md)'s M-052 already commits to a "permanent validated identity ledger" with live/reserved/retired state, but scopes it to criteria imported during brownfield migration. The same gap exists for every native prefix Cliewen mints for itself, so confining the fix to imports would leave the defect live everywhere else.

## Decision outcome

**A persisted ledger becomes the registry for every native ID prefix, superseding the "corpus is the registry" clause of ADR-009.** The corpus scan (`clue validate`) becomes a cross-check against the ledger rather than the source of truth for allocation.

- **File and format:** `.clue/id-ledger.yaml` at the repository root — not under `docs/`, which holds authored corpus prose with per-file frontmatter, while the ledger is one generated/operational registry file. One entry per ID: `id`, `prefix`, `component`, `state` (`reserved | live | retired`), and, only for entries imported from a brownfield source, `source-revision` and `source-location` (M-052's original fields). A per-prefix `counters` map holds the last-issued numeric component.
- **State machine:** `reserved` (an ID has been claimed but no artifact with it exists yet) → `live` (an artifact with that ID is present in the corpus) → `retired` (the artifact was deleted). A retired entry is never removed and never reissued.
- **O(1) lookups:** both queries an agent needs — "is this ID used" and "what is the next ID for this prefix" — are in-memory map lookups after a single `Load()`, never a corpus scan. The `counters` map makes allocation an increment; a `byID` map makes existence a lookup. Scanning the corpus stays necessary only once, for backfill, and afterward only as a cross-check inside `clue validate`.
- **Validator cross-check:** `clue validate` rejects a live artifact whose ID is `retired` or `reserved` under a different reservation in the ledger, and rejects a live artifact with no ledger entry once the ledger is the corpus's source of truth (gated so a corpus without a ledger yet is unaffected). This composes with, and does not replace, ADR-034's `supersedes:` mechanism: `supersedes:` gives a reader a pointer to what replaced a deleted ID; the ledger gives the allocator a reason never to reissue that ID's number in the first place.
- **Relationship to AC tombstones:** [ADR-007](ADR-007-ac-lifecycle.md)'s `@retired` tag is unchanged and remains the record for a criterion whose *file* survives; the ledger's `retired` state is the equivalent memory for every ID type whose file does not survive deletion (ADR-034's case), plus the criterion IDs M-052 already scoped in.
- **Backfill:** a one-time, idempotent migration step seeds the ledger from the current corpus scan — one `live` entry per currently-live ID, `counters` seeded at each prefix's current max — so existing history is not renumbered and the transition is invisible to already-issued IDs.

This decision revises M-052's stated scope accordingly (tracked in the same change): the ledger it commits to building covers all native prefixes, with the imported-criterion fields (`source-revision`, `source-location`) as one populated case among others, not the whole of its purpose.

**This decision does not itself implement the ledger.** `internal/ledger/`, the `clue id next` allocation subcommand, the `checkLedger` validator rule, and the `clue migrate` backfill step are follow-up work under M-052, each its own full change with its own acceptance.

## Rejected: keep scan-and-max, add only a retired-ID blocklist

A blocklist recording just "these numbers are dead" would stop reuse without giving allocation an O(1) path — every `next-free-ID` computation would still need the live max from a scan plus a lookup against the blocklist. The counters map costs no more to maintain and removes the scan from the hot path entirely.

## Rejected: keep the ledger scoped to imported criteria only, per M-052 as written

The defect being closed — an ID surviving in nobody's memory once its file is gone — is identical for a native `CH-xxx` or `PDR-xxx` ID and an imported `AC-xxx`. Two allocators for the same failure mode would mean fixing it twice, once now and once whenever native-prefix reuse is eventually noticed as its own problem.

## Carrier

`docs/plans/P-011-truthful-brownfield-migration.md` M-052 (scope), this record, and — once implemented — `internal/corpus/rules.go`'s validator rules and the `clue-delta` skill's ID-allocation guidance are the carriers of this decision.
