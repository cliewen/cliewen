---
id: CH-118
type: change
status: open
links: [P-011]
title: In-flight source work survives extraction as a durable imported-change record
---

# CH-118 — In-flight source work survives extraction as a durable imported-change record

## Serves

[P-011](../../docs/plans/P-011-truthful-brownfield-migration.md) M-054, currently `todo`.

## What

The current OpenSpec mapping row for a source repository's pending change (`changes/<name>/`) converts it into a plan milestone plus a `status: draft` capability, and states plainly that its `tasks.md` "dies" — `clue-delta` regenerates tasks once implementation starts. That discards the dependency graph between a source change's tasks and the design rationale and proof links that connected them, exactly the "reduces in-flight work to an insufficient plan row" gap M-011 names.

This change adds a new native artifact type, `imported-change`, one record per source change extraction preserves, holding:

- the pinned origin (`source-revision`, `source-location`, reusing ADR-048's ledger field names),
- intent and design rationale prose,
- dependency links to other `imported-change` records (via the ordinary `links:` field, resolved like any other artifact link),
- a task-to-criterion proof-links table, and
- a completion state (`in-progress` | `complete`).

`clue validate` gains a new rule requiring every proof-linked criterion on a `complete` record to exist, be non-`@draft`, and non-retired — an unproven proof link makes `complete` an unjustified claim, the same failure shape ADR-049's unjustified-disposition class already established for parity. The OpenSpec mapping row is rewritten to target this record instead of the lossy plan-row treatment, and the `clue-extract` skill's rehearsal guidance gains the agent-side rule that extraction must not delete an incomplete source change's in-flight work until its `imported-change` record is `complete` (a source-repository judgment `clue` cannot make, since it never reads the source repo).

## Why not defer further

P-011's M-054 is the next `todo` milestone in strict order; M-052 and M-053 already closed the identity and parity gaps this milestone's fixture composes with.
