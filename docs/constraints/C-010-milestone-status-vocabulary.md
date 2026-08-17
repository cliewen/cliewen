---
id: C-010
type: constraint
status: active
links: [G-001, P-013, AN-018]
title: Milestone status values in plan tables follow one vocabulary
source: docs/plans/README.md (milestone table convention); PDR-048
enforcement: partial
---

# C-010 — Milestone status values follow one vocabulary

A milestone row's status is one of `todo`, `doing`, `done`, or `dropped`. The vocabulary codifies what ten campaigns already did — every row in P-001 through P-010 read `todo` or `done` — and names the two states that were being written as prose instead.

**Every milestone carries a verifiable exit criterion.** A row states a condition a reader can check against the corpus and the tool, so that `done` is a finding rather than an assertion. A milestone whose exit criterion cannot be checked cannot be closed on evidence, and closing it anyway is how a campaign comes to report work it did not do — which is the failure [G-001](../goals/G-001-verifiable-thread.md)'s thread exists to prevent, applied to the plan layer.

**Checked by:** `clue validate` ([AC-095](../capabilities/CAP-002-validate/criteria.md)) — a plan table declaring a `Status` column has every cell in it read against `todo | doing | done | dropped` ([PDR-048](../decisions/PDR-048-plan-and-workspace-bookkeeping-stays-truthful.md)). A table with no status column is not a milestone table and is not read. The vocabulary half is held mechanically; the exit-criterion half is not.

**Residual:** whether an exit criterion is genuinely verifiable is a judgement about prose, and no machine will make it — a row can satisfy every lint and still promise something nobody can check. The cost lands at closure: a campaign closes on a row that reads like evidence and is not, and the error surfaces only when someone later tries to use what the milestone claimed to deliver.
