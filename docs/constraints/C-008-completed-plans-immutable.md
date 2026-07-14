---
id: C-008
type: constraint
status: active
links: []
title: Completed plans are immutable
source: docs/plans/README.md, ARCH-001 lifetime classes
enforcement: agent
---

# C-008 — Completed plans are immutable

A plan at `status: completed` is frozen and never deleted: the plans index doubles as the project's achievement overview, and rewriting a finished campaign rewrites history.

**Promotion trigger:** `clue` gains git-diff context and fails a change that modifies a file whose `main`-side status is `completed` — then `enforcement: machine`. Already noted as a deliberate limit in CAP-002's design.
