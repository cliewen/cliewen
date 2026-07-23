---
id: ADR-027
type: decision
status: inferred
links: [P-005, ADR-025, ADR-026, ADR-017, C-014, C-015]
title: Quality scenarios are constraints — the quality type folds into the register
author: agent
accepted-by: []
---

# ADR-027 — Quality scenarios are constraints

## Context and problem statement

The corpus carried a `quality` artifact type (QS-xxx) for verifiable non-functional requirements, in its own folder, alongside the `constraint` register (C-xxx). But the two were the same shape: a quality scenario is a cross-cutting rule checked against every proposal, and every constraint is a cross-cutting rule checked against every proposal ([ADR-017](ADR-017-conventions-are-constraints.md)). The corpus held exactly two quality scenarios — QS-001, which *is* the CI coverage gate, and QS-002, the onboarding-time bar — and neither was anything a constraint could not carry once constraints gained a `source` and `enforcement` field. [ADR-025](ADR-025-one-status-lifecycle.md) removed the last mechanical distinction by making `quality` an ordinary default-lifecycle type with no special validator handling. What does the separate type still buy, and if nothing, where do its scenarios live?

## Decision outcome

**A quality scenario is a constraint. The `quality` type folds into the constraints register, and its two scenarios are re-minted as constraints.**

- **QS-001 → [C-014](../constraints/C-014-coverage-floor.md)** — total Go statement coverage ≥ 80%, `enforcement: machine` (the CI coverage gate is the check).
- **QS-002 → [C-015](../constraints/C-015-onboarding-under-30-minutes.md)** — a new user reaches a first green validate in under 30 minutes, `enforcement: human` (a journey no test can time).
- The `quality` type leaves the documented taxonomy and the `clue init` scaffold. No validator code changes: since [ADR-025](ADR-025-one-status-lifecycle.md) the type was already just a default-lifecycle name, so an adopter who keeps a `quality/` folder still validates green ([ADR-026](ADR-026-adopter-types-default-lifecycle.md)) — the fold removes the type from *Cliewen's* corpus, not from the space of legal types.
- **The QS-001 and QS-002 IDs are retired tombstones, never reused.** Because the fold deletes `docs/quality/`, the tombstone cannot be a `@retired` marker in the file the way a retired acceptance criterion leaves one; the decision-log row dated with this change is the tombstone of record, and `git log docs/quality/` keeps the provenance.

**Supersession reading (for immutable prior texts).** Records that cannot be edited — [ADR-006](ADR-006-test-purpose-taxonomy.md)'s "QS lane" for performance tests, [PDR-012](PDR-012-agentic-review-before-publication.md)'s "constraints and quality scenarios", and completed plans [P-002](../plans/P-002-leaves-home.md) and [P-004](../plans/P-004-first-try.md) that name QS-002 — are read as pointers to the successor constraints: wherever they say "quality scenario QS-00x", read "constraint C-01x". [ADR-025](ADR-025-one-status-lifecycle.md)'s list of types on the default lifecycle names `quality` among them; that enumeration is illustrative, not the canonical taxonomy, and stays true either way — a corpus that keeps a `quality` type still validates against the default. Those files are never rewritten; this clause carries the mapping.

**Carrier:** the constraints register holds C-014 and C-015; the removal of `docs/quality/` from the corpus and the `clue init` template, and of the QS row from the taxonomy prose and diagrams, ships the fold to adopters.

### Rejected: keep the quality type as a specialization of constraint

A "quality scenario is a constraint that carries a measurement" subtype would keep two names for one register and a folder holding two files. The measurement is already expressible: C-014 names its threshold and its machine check, C-015 names its bound and its human check. A subtype earns its keep only when it changes how the artifact is processed, and after ADR-025 nothing processes `quality` differently.

### Rejected: demote the scenarios to decision-log rows

The coverage gate already has a log row (2026-07-12) recording *why* 80%; that rationale is cheap-to-reverse tuning. But the standing rule "coverage must not fall below the floor, checked every build" is a constraint the register exists to inventory — a log row cannot carry `source`/`enforcement` or appear in the agent-enforced count. The log keeps the origin rationale; the register carries the standing rule.
