---
id: CH-157
type: change
status: open
links: [P-016, PDR-044]
title: Keep P-016 outcome-neutral until M-075
---

# CH-157 — Keep P-016 outcome-neutral until M-075

## Proposal

Revise P-016 so its milestones require M-075 to decide the decision-record taxonomy and compact shape instead of pre-committing the plan to those outcomes, and make M-079 cover every decision record created earlier in the campaign.

The correction already exists in the preserved human-authored commit `35e734c`, made after [cliewen/cliewen PR #168](https://github.com/cliewen/cliewen/pull/168)'s reviewed head and therefore absent from its merge. This change recovers that commit without rewriting it, records the plan revision in PDR-045, and integrates it through a new review boundary rooted at merged `main`.

## Scope boundary

This change revises P-016 and its index description, and clears the stray ledger reservations identified by the preserved commit. It does not choose the taxonomy, compact record structure, or migration behavior; those remain M-075 outcomes.
