---
id: CH-053
type: change
status: open
links: [P-005, M-018, G-001, C-013]
title: Fold quality scenarios into the constraints register
---

# CH-053 — Quality scenarios are constraints

## What

Retire the `quality` artifact type and re-home its two scenarios as constraints. Mints ADR-027 (the decision, with a supersession-reading clause for immutable prior texts), C-014 (coverage floor ≥ 80%, `enforcement: machine`), and C-015 (onboarding under 30 minutes, `enforcement: human`). Deletes `docs/quality/` and its `clue init` template. Adds a decision-log row as the QS-001/QS-002 tombstone. Updates every surface that named quality scenarios as a distinct type: the corpus README and its scaffold template, the architecture and methodology diagrams, CAP-001's links, the change-tier and extraction skills, the PR template, and the CI coverage-gate step name.

## Why

A quality scenario was always a constraint in disguise — a cross-cutting rule checked against every proposal, each with an enforcement mechanism. QS-001 is literally the CI coverage gate. Since CH-052 made `quality` an ordinary default-lifecycle type with no special handling, the type now carries no meaning the constraints register does not already provide; keeping it is a second name for one idea. "Does the core need it?" — the core needs cross-cutting rules with enforcement, which is exactly what a constraint is.

## Decision boundary

Folding an artifact type into another is a corpus-format change, so ADR-027 records it; the change alters what the documented taxonomy and the `clue init` scaffold contain, and it touches methodology carriers (skills), so by the red line ([C-013](../../docs/constraints/C-013-core-changes-need-decision.md)) it is a full change accepted only at human merge. Immutable prior texts (ADR-006's "QS lane", PDR-012, completed plans P-002/P-004) are never edited; ADR-027 states they are read as pointers to the successor constraints. The QS-001/QS-002 IDs are retired tombstones, never reused. Editing skill sources requires `go generate ./internal/skills` and a version stamp aligned to the coming 0.6.0 release.
