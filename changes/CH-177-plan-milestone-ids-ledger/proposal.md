---
id: CH-177
type: change
status: open
links: [G-006, G-001, CAP-010, ADR-048]
title: Plan the ledger's missing milestone identities
---

# CH-177 — Plan the ledger's missing milestone identities

## What

Add P-022, a new plan, to `docs/plans/`. The plan commits to closing [G-006](../../docs/goals/G-006-milestone-ids-in-the-ledger.md): the identity ledger's `M` counter is never seeded, because a milestone is a table row inside a plan file, not a frontmatter-bearing artifact the existing backfill (`planLedgerBackfill`, `internal/migrate/migrate.go`) can scan. `clue id next M` currently returns `M-001` and writes that false reservation, even though 68 `M-xxx` identities already exist non-contiguously across `docs/plans/*.md` (current corpus maximum: `M-090`).

This change adds no implementation. It proposes the campaign — three milestones covering the backfill fix, the allocator behavior it enables, and a `clue validate` collision check for milestone IDs matching what every other native prefix already gets — so the shape of the work is human-reviewable before anyone starts it.

## Why

[CH-171](../CH-171-team-safe-id-allocation/proposal.md) (CAP-010) made allocation safe under concurrent access by serializing claims through a coordinated Git-remote ledger branch. That work assumed every prefix's starting counter was already trustworthy; it never touched how a counter is seeded, so the one prefix whose seed was always wrong — `M` — is untouched by it. G-006 records this as a corrupted-ledger bug, not a missing feature: the backfill does not merely fail to seed `M`, running `clue id next M` today writes a false live entry.

`internal/corpus/plans.go`'s `PlanMilestones` (added for `clue next`, CH-172/CH-174) already extracts milestone rows with their identity and status from every plan's table. This plan proposes reusing that same extraction as the read path for both the ledger backfill and the new collision check, rather than inventing a second way to parse a milestone table.

## Chosen shape

- P-022 links G-006, G-001, CAP-010, and ADR-048.
- Three milestones: seed the `M` counter from milestone tables (backfill), confirm `clue id next M` allocates correctly above the corpus maximum under CH-171's coordination, and make `clue validate` report a milestone-ID collision like it does for every other native prefix.
- Milestone declaration itself is unchanged — a table with `ID` and `Status` headers inside a plan file remains the only door; no milestone gains a file or frontmatter.

## Compatibility

No corpus obligation narrows or widens for an adopter yet — this change adds a plan artifact only. The implementing change(s) this plan authorizes will need to migrate an adopter's existing ledger without renumbering any already-declared `M-xxx`, per ADR-048.

## Plan item

Introduces P-022 (`draft`). This change's own digest sets P-022's status; no milestone in this change's own tasks is implementation work against it.
