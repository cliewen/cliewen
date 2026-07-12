---
id: ADR-006
type: decision
status: inferred
links: [ADR-005, CAP-002]
title: Every test declares its purpose from a small taxonomy
author: agent
---

# ADR-006 — Every test declares its purpose

## Context and problem statement

Not every test verifies an acceptance criterion: coverage backstops, sanity checks and architectural checks are legitimate craft. The AC↔test contract ([ADR-005](ADR-005-test-reference-convention.md)) says nothing about them, so a test suite can silently accumulate tests whose intent nobody can name. The principle — tag every test so its purpose is clear — comes from the maintainer's Intent Engineering practice; this ADR is the agent's concrete taxonomy design awaiting human promotion.

## Decision outcome

**Every test declares exactly one purpose, machine-checked.** The starting vocabulary, extended only by ADR when a new class earns it:

| Purpose | Meaning |
|---|---|
| `AC<digits>` | Verifies that acceptance criterion — the red thread's edge (ADR-005) |
| `Unit` | Implementation-detail / coverage backstop for code paths between the ACs |
| `Sanity` | Invariants of the environment or the repo itself (e.g. "this repo's own corpus validates") |
| `Arch` | Structural/architectural checks (dependency direction, layering) |

Carried per ADR-005's mechanics: a framework tag where tags exist (`@Tag("AC-004")`, `@Tag("unit")`); in Go, the name prefix — `TestAC004_…`, `TestUnit_…`, `TestSanity_…`, `TestArch_…`. A test matching no purpose fails `clue validate` (AC-011). This also tightens the AC reference to a **prefix**: `TestUnit_HandlesAC004Edge` is a unit test that mentions an AC, not an AC test.

**Carrier:** the purpose check in `clue`'s per-language harvesters (machine), plus the taxonomy table in the `clue-delta` skill (agent).

### Rejected: leaving non-AC tests unclassified

The status quo. It makes the AC-010 check one-sided and lets intent-free tests accumulate — the test-suite version of doc-slop.

### Purpose tags and runner tags coexist — separate namespaces, same mechanism

An AC may legitimately need verification at several levels — unit, integration, E2E, possibly with different tools. Level answers *how a test runs*, and its consumer is the **runner**, not the methodology: level tags exist so pipelines can filter (fast tests on every PR, slow E2E nightly). The rule that avoids redundancy: a test carries **exactly one purpose tag, which `clue` reads, and any number of runner tags, which `clue` ignores**. `@Tag("AC-022")` and `@Tag("integration")` on the same test is correct and not redundant — each tag has exactly one consumer. In Go, purpose lives in the name prefix and level uses Go's own runner idioms (`//go:build integration`, `testing.Short()`).

An AC requiring several levels gets several tests referencing the same AC, each with its own runner tags. Whether that set covers the criterion *adequately* is meaning, not form — PR review, not lint. (Door, only if real use demands it: a per-AC required-levels annotation in `criteria.md` that `checkACTests` cross-checks; deliberately not built, as it doubles every AC's bookkeeping.)

### Rejected: level and quality classes inside the purpose taxonomy (Integration, Performance, …)

**Integration/E2E are execution levels, not purposes** — runner metadata per the rule above, with no consumer in `clue` (§7: a field nobody reads gets removed before it exists). And **Performance is not a purpose class but the QS lane**: a performance test must trace to a specific quality scenario, mirroring how an AC test traces to its AC.

**Door:** a `QS<digits>` purpose class joins the taxonomy when the first QS-verifying test exists, extending `checkACTests` the same way for the thread's second strand.
