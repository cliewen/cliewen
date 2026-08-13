---
id: C-005
type: constraint
status: active
links: [PDR-042]
title: Every Cliewen proposal declares its plan item or declares itself plan-less
source: clue-delta skill step 2 (Propose)
enforcement: partial
---

# C-005 — Every Cliewen proposal declares its plan item or plan-less

A full change's `proposal.md` references the plan item it serves (P/M-IDs in `links`) or states plan-less explicitly. A simple change under [PDR-042](../decisions/PDR-042-routing-recommends-contract-aware-effort.md) has no Cliewen proposal and makes no plan declaration. No fake plan items.

**Checked by:** `clue validate` ([AC-092](../capabilities/CAP-002-validate/criteria.md)) — a `changes/*/proposal.md` whose links name no `P-xxx` or `M-xxx` and whose body never says `plan-less` fails.

**Residual:** whether the plan item named is the one the full change actually serves. A proposal can link a real milestone it has nothing to do with and pass. The cost is a plan that reads as accounted-for while the work drifted somewhere else; the merge review is where that is caught.
