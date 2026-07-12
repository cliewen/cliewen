---
id: P-001
type: plan
status: completed
links: [G-001]
title: Elaboration baseline
---

# P-001 — Elaboration baseline

> **Completed 2026-07-13** — all milestones done; frozen immutable. The baseline's lessons live as [ADR-001…ADR-010](../decisions/README.md), distilled continuously as each change landed. Successor: [P-002 — Cliewen leaves home](P-002-leaves-home.md).

A 2–3 week **elaboration baseline** in RUP's sense: a running skeleton that proves the architecture and is built upon, not thrown away. Risk-driven — the milestones retire the biggest risks first, as running code. Serves [G-001](../goals/G-001-verifiable-thread.md).

## Milestones

| ID | Milestone (exit criterion) | Status | Evidence |
|---|---|---|---|
| M-001 | The loop closes once, end-to-end: goal in → proposal branch → implement → permanent `/docs` updated → transient files gone → clean merge | `done` | CH-001 (merge commit of branch `ch-001-bootstrap`) |
| M-002 | The thread is machine-checkable: linter fails the build when an AC lacks a test or a test lacks an AC | `done` | CH-002 (`clue validate` + CI wall) and CH-003 (AC↔test contract, ADR-005) |
| M-003 | Brownfield extraction works on model2diagram | `done` | CH-004 (clue namespaces/JVM harvest/provenance + clue-extract skill) and model2diagram CH-001 (extraction PR, merge commit `8a5a7af`, 2026-07-12) |

## Explicitly out of the baseline

Multi-agent orchestration; semantic consistency checking of shared docs (parked as known edge); plan milestone mechanics beyond one markdown file + one linter rule; production feedback loop (V3 door only); dot-principles / external catalog integration (door defined via constraint `source:`, nothing more); `enforcement:` classes beyond `machine`; `clue locate`.

## Mutation rules (lintable)

Status fields in the milestone table may mutate in any merge (bookkeeping). Everything else in this file changes only via a change that declares itself a plan revision, backed by an ADR — plan adjustments ARE decisions.
