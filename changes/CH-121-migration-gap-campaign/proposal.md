---
id: CH-121
type: change
status: open
links: [P-011, P-012, PDR-025, PDR-026]
title: Close P-011; open P-012 for the migration gap and P-013 for simplification
---

# CH-121 — Close P-011; open P-012 for the migration gap and P-013 for simplification

Serves [P-011](../../docs/plans/P-011-truthful-brownfield-migration.md) as its closing change, and opens its successor.

## What

P-011's five milestones are all `done` (CH-116…CH-120), but the campaign was never closed: its frontmatter and its index row still read `active`, and no successor is designated. CH-120's digest updated M-056's row only.

This change closes P-011, records the sequencing decision that follows from re-deriving the campaign's own results, opens **P-012** as the campaign that finishes the brownfield migration gap on real evidence, and opens **P-013** as the deferred simplification campaign in `draft`.

## Why

[PDR-025](../../docs/decisions/PDR-025-brownfield-migration-precedes-simplification.md) named P-012 as the simplification campaign. Re-deriving the originating assessment's gaps from the corpus and the tool — rather than from P-011's own evidence column — does not support handing off yet:

- **A bulk-drafted migration still reaches green.** AC-114 fails a `@draft`, `Human`, or retirement disposition carrying no matching justification. A justification is prose, and nothing distinguishes one honest deferral from a formulaic repetition of it across most of a corpus, which is the shape the assessment actually reported.
- **The durable extraction report can still disagree with its own tree.** The assessment's second blocker was a narrative report whose scenario and draft counts contradicted the committed result. `clue parity` reconciles manifests; the `clue-extract` target contract asks the `/docs/analysis` report to record what was found, with nothing binding its figures to the derived manifest.
- **Per-criterion inspectability was declined, not delivered.** The assessment asked for a committed per-criterion artifact. M-053 deliberately refuses one — "the report stays derived rather than becoming an editable coverage registry". That may be the better design, but a declined request is currently recorded as a satisfied one.
- **Cliewen does not run its own identity ledger.** `clue id next` fails in this repository because the ledger is missing, and `clue migrate` still offers MIG-008. The mechanism M-052 shipped has never been exercised in the repository that shipped it.

Scale is also unproven: every proof in M-052…M-056 is fixture-sized, and [AN-016](../../docs/analysis/AN-016-disposable-brownfield-migration-fixtures.md) says so itself, disclaiming any claim about a production corpus.

Two smaller findings from the same sweep ride into the new campaign's first milestone rather than this change: `docs/capabilities/CAP-003-extract/design.md` still calls adopter-CI binary distribution unsolved and parked, which [ADR-038](../../docs/decisions/ADR-038-upstream-validation-workflow.md) and the verified install scripts closed; and the originating assessment is cited in prose by both P-011 and PDR-025 with no analysis record holding it, unlike [AN-014](../../docs/analysis/AN-014-confidential-migration-assessment.md).

## Design

**PDR-026** supersedes PDR-025's naming clause: P-012 finishes the migration gap on re-derived evidence, P-013 is simplification. Its enduring rule is the one the findings above expose — a campaign closes on gate status re-derived from the corpus and the tool, not on its own evidence column. PDR because it is process and expensive to unwind once the plan register carries it.

**P-012** opens `active` with M-057…M-061, continuing corpus-global milestone numbering. M-057 makes the campaign's own evidence base honest before extending it; M-058 and M-059 repair the two mechanisms that check form where the assessment asked about meaning; M-060 decides the declined request explicitly; M-061 proves the contract at the assessment's order of magnitude and closes the campaign.

**P-013** opens `draft`, held until P-012 completes — the shape PDR-025 used for P-011.

Each P-012 milestone that changes what `clue validate` or `clue parity` asserts is core-adjacent under [C-013](../../docs/constraints/C-013-core-changes-need-decision.md) and carries its own decision record in its own change. None of that mechanism work happens here.

## Scope

- `docs/decisions/PDR-026-<slug>.md` (new); `docs/decisions/PDR-025-…` gains its forward pointer; `docs/decisions/README.md` index row.
- `docs/plans/P-011-truthful-brownfield-migration.md`: `status: completed` (frontmatter only — [C-008](../../docs/constraints/C-008-completed-plans-immutable.md) freezes it after this change).
- `docs/plans/P-012-<slug>.md` and `docs/plans/P-013-<slug>.md` (new); `docs/plans/README.md` index rows.
- No `CHANGELOG.md` entry: closing one campaign and opening two moves no shipped behaviour, capability, contract, command, or user workflow.

## Not in this change

No corpus obligation, validator rule, parity rule, or skill changes. The identity ledger stays unseeded here, and IDs are assigned by continuing the existing sequence rather than by `clue id next`; seeding it is M-057's work, and leaving the broken allocator visible is the campaign's own first finding rather than something to quietly repair on the way past.
