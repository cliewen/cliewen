---
id: CH-178
type: change
status: open
links: [P-022, CAP-010, ADR-048]
title: Seed the milestone counter from every plan's milestone table
---

# CH-178 — Seed the milestone counter from every plan's milestone table

## What

Implements [P-022](../../docs/plans/P-022-milestone-ids-in-the-ledger.md)/M-091. `planLedgerBackfill` (`internal/migrate/migrate.go`) now marks every `M-xxx` identity `corpus.LedgerMilestoneIdentities` reports across `docs/plans/*.md` live or retired in the same pass it marks artifact IDs, using each milestone's own status — `done`/`dropped` as retired, `todo`/`doing` as live — the same way it already does for acceptance criteria. `corpus.LedgerMilestoneIdentities` is new: it walks every plan artifact regardless of the plan's own status (draft, active, or completed) and returns declared milestones in ID order, sharing `PlanMilestones`'s existing table extraction.

P-022 moves from `draft` to `active` in this change's digest, since it now has actionable work in flight.

## Why

The backfill has never had a read path for milestone identities: a milestone is a table row inside a plan, not a frontmatter-bearing artifact, so `planLedgerBackfill`'s scan of `corpus.Scan`'s `ByID` map never saw one. Running the backfill against this repository's own corpus today seeds the `M` counter at an implicit zero despite 90-odd `M-xxx` identities already declared non-contiguously, so `clue id next M` would allocate a colliding `M-001`. This closes that gap using the same milestone extraction `clue next` (CH-172/CH-174) already relies on, so the backfill and the allocator read milestone identities one way.

## Compatibility

An adopter with no ledger yet now gets its `M` counter seeded correctly on first backfill; nothing changes for a repository that already carries a version-two ledger, since `planLedgerBackfill` only runs once, before the ledger file exists (M-092 covers seeding an *existing* ledger's `M` counter — out of scope here).

## Plan item

P-022/M-091. Sets the milestone `todo` → `done` with this change's evidence in the digest.
