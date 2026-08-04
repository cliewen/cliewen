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

A milestone row's status is one of `todo`, `doing`, `done`, or `dropped`. The vocabulary codifies what ten campaigns already did — every row in P-001 through P-010 read `todo` or `done` — and names the two states that were being written as prose instead. A `dropped` milestone states its reason in the evidence column, the way a `[-]` task states it on its line.

**Checked by:** `clue validate` ([AC-095](../capabilities/CAP-002-validate/criteria.md)) — a plan table declaring a `Status` column has every cell in it read against `todo | doing | done | dropped` ([decision log](../decisions/log.md), 2026-08-04). A table with no status column is not a milestone table and is not read.
