---
id: C-010
type: constraint
status: active
links: []
title: Milestone status values in plan tables follow one vocabulary
source: docs/plans/README.md (milestone table convention)
enforcement: agent
---

# C-010 — Milestone status values follow one vocabulary

Milestone rows in plan tables carry a status column, but no vocabulary is enforced anywhere — plans have used values like `todo` and `done` by convention. Agents keep milestone statuses to a small consistent set within each plan.

**Promotion trigger:** milestones get a declared status vocabulary (a decision to make when promoting) and `clue validate` parses plan milestone tables against it — then `enforcement: machine`.
