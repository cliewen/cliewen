---
id: C-010
type: constraint
status: active
links: []
title: Milestone status values in plan tables follow one vocabulary
source: docs/plans/README.md (milestone table convention)
enforcement: machine
---

# C-010 — Milestone status values follow one vocabulary

Milestone rows in plan tables carry a status column, but no vocabulary is enforced anywhere — plans have used values like `todo` and `done` by convention. Agents keep milestone statuses to a small consistent set within each plan.

**Checked by:** `clue validate` ([AC-095](../capabilities/CAP-002-validate/criteria.md)) — a plan table declaring a `Status` column has every cell in it read against `todo | doing | done | dropped` ([decision log](../decisions/log.md), 2026-08-04). A table with no status column is not a milestone table and is not read.
