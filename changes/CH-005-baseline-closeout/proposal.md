---
id: CH-005
type: change
status: open
links: [P-001, M-003]
title: Baseline closeout — M-003 done, versioning enters the inbox
---

# CH-005 — Baseline closeout

## Why

M-003's evidence exists: model2diagram merged its CH-001 extraction PR (merge commit `8a5a7af`, 2026-07-12) — `/docs` is its system-of-record, every AC ID survived, `clue validate` green. The plan table should say so. Separately, the maintainer proposed versioning `clue` and the skills; per ADR-002 that enters the corpus as a proposed goal, not as unplanned work.

## What

1. P-001 milestone table: M-003 → `done` with evidence (status mutation, explicitly permitted in any merge by P-001's own rules).
2. `G-002` (proposed): `clue` and the skills carry versions — intake only; design and acceptance belong to the next plan.
