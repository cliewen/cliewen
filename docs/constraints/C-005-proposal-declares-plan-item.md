---
id: C-005
type: constraint
status: active
links: []
title: Every Cliewen proposal declares its plan item or declares itself plan-less
source: AGENTS.md rule 3, clue-delta skill step 2
enforcement: partial
---

# C-005 — Every Cliewen proposal declares its plan item or plan-less

A full change's `proposal.md` references the plan item it serves (P/M-IDs in `links`) or states plan-less explicitly; a light change makes the same declaration in its PR description. A plain change under [PDR-011](../decisions/PDR-011-plain-changes-bypass-cliewen.md) has no Cliewen proposal and makes no plan declaration. No fake plan items.

**Checked by:** `clue validate` ([AC-092](../capabilities/CAP-002-validate/criteria.md)) — a `changes/*/proposal.md` whose links name no `P-xxx` or `M-xxx` and whose body never says `plan-less` fails.

**Residual:** the light tier, whose declaration lives in a pull-request description outside the tree, and whether the plan item named is the one the change actually serves. A proposal can link a real milestone it has nothing to do with and pass. The cost is a plan that reads as accounted-for while the work drifted somewhere else; the merge review is where that is caught.
