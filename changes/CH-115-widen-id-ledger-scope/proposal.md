---
links: [P-011]
---

# CH-115 — Widen M-052 to a corpus-wide identity ledger

## What

M-052 in P-011 currently commits to "a permanent validated identity ledger holding the pinned source revision, source location, exact canonical identifier, and live/reserved/retired state for every imported criterion" — scoped to criteria imported during brownfield migration. This change widens that scope: the same ledger mechanism also becomes the allocator and retirement record for Cliewen's own native ID prefixes (`CH`, `C`, `P`, `M`, `PDR`, `ADR`, `CAP`, `G`, `AN`, ...), not only imported criteria.

This change revises the plan text and adds a decision record. It does not implement the ledger — that is a separate full change once this scope decision is accepted.

## Why

Today, allocation for every native prefix is scan-and-max: `docs/decisions/ADR-009-ac-id-namespaces.md` states "the corpus is the registry" for AC IDs, and the same practice is followed by convention for every other prefix. Nothing persists once an artifact's file is deleted — Cliewen's retirement model is deletion ([ADR-034](../../docs/decisions/ADR-034-retirement-is-deletion.md)) — so a retired ID's number can silently be re-minted later, colliding with historic evidence, commit history, or external references that still carry its original meaning. This gap exists for every native ID type, not only for imported criteria, so scoping the fix to imports only would leave the same defect live everywhere else.

## Decision

A new ADR (ADR-046) is drafted alongside this proposal, extending ADR-009's "corpus is the registry" clause to a persisted ledger for all native prefixes. `docs/decisions/ADR-046-corpus-wide-id-ledger.md` records the design: a `.clue/id-ledger.yaml` ledger, in-memory maps for O(1) "is this ID used" and "next ID for this prefix" lookups (no scan on the hot path), a `reserved → live → retired` state machine, and how it composes with ADR-034 (retirement is deletion) and ADR-007 (AC tombstoning — the closest existing prior art for a persisted retired state). It starts `status: inferred`, `author: agent`, per the decision-records convention; human acceptance promotes it to `verified`.

## Out of scope

Implementing `internal/ledger/`, the new `clue id next` subcommand, the validator rule, and the `clue migrate` backfill step — those land in a follow-up full change once this scope decision and its ADR are accepted.
