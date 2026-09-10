---
id: CH-179
type: change
status: open
links: [P-022, CAP-010, ADR-048]
title: Backfill milestone identities into a ledger that already exists
---

# CH-179 — Backfill milestone identities into a ledger that already exists

## What

Implements [P-022](../../docs/plans/P-022-milestone-ids-in-the-ledger.md)/M-092. `planLedgerBackfill` (`internal/migrate/migrate.go`) only runs the first time a repository gets a ledger at all — its own doc comment calls it idempotent "once the file exists." CH-178/M-091 taught it to read milestone tables, but a repository whose ledger already existed before that change (this one) never gets a second pass: its `M` counter stays at an implicit zero forever, exactly the false-zero collision ADR-048 exists to prevent for every other prefix.

A new migration, `planMilestoneLedgerBackfill` (MIG-016), closes that gap: it runs only when a ledger already exists, adds one live-or-retired event per milestone identity `corpus.LedgerMilestoneIdentities` reports that the ledger does not already carry under any state, and leaves every other identity untouched. It reports nothing once every declared milestone is covered. To avoid two migrations proposing conflicting rewrites of the same ledger file in one plan — which `Apply` cannot detect, since both would record the same on-disk bytes as `Before` and whichever writes last would silently discard the other's work — it skips a ledger this same plan is about to create or rewrite (a fresh ledger, or `planLedgerEvents` converting or repairing a version-one one); that narrower combination backfills its milestones on the run after the conversion.

## Why

Confirmed against this repository's own `.clue/id-ledger.yaml`: it carries zero `M-xxx` entries despite 90-odd milestones declared across `docs/plans/*.md`, because its ledger already existed before CH-178 gave the backfill a read path for milestones. `clue id next M` would allocate a colliding `M-001` the moment anyone runs it here, or in any adopter repository that adopted the ledger before this release.

## Compatibility

A repository with no ledger yet is unaffected — `planLedgerBackfill` already covers it in full. A repository already carrying every declared milestone plans no further change. A repository still on a version-one ledger and also missing milestones needs one extra `clue migrate --apply` run after its conversion, rather than getting both in the same pass; this is the one combination the conflict guard above defers.

## Plan item

P-022/M-092. Sets the milestone `todo` → `done` with this change's evidence in the digest.
