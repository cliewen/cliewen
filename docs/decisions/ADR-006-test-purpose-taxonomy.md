---
id: ADR-006
type: decision
status: verified
links: [ADR-005, CAP-002]
title: Every test declares its purpose from a small taxonomy
author: agent
accepted-by: Flemming N. Larsen (2026-07-12)
---

# ADR-006 — Every test declares its purpose

## Context and problem statement

The AC↔test contract does not classify legitimate non-AC tests such as coverage backstops, sanity checks, and architectural checks. Without a purpose, test intent becomes unreviewable.

## Decision outcome

**Every test declares exactly one machine-checked purpose.** The vocabulary is extended only by decision when a new class earns it:

| Purpose | Meaning |
|---|---|
| `AC<digits>` | Verifies that acceptance criterion — the red thread's edge (ADR-005) |
| `Unit` | Implementation-detail / coverage backstop for code paths between the ACs |
| `Sanity` | Invariants of the environment or the repo itself (e.g. "this repo's own corpus validates") |
| `Arch` | Structural/architectural checks (dependency direction, layering) |

Tests use ADR-005's tag or name mechanics. A test matching no purpose fails `clue validate`; an AC-looking token inside `TestUnit_HandlesAC004Edge` does not change its `Unit` purpose. Segmented prefixes, lowercase suffixes, and carrier aliases follow the canonical criterion grammar.

**Carrier:** the purpose check in `clue`'s per-language harvesters and the taxonomy table in the `clue-delta` skill.

### Purpose tags and runner tags coexist — separate namespaces, same mechanism

Terminology per the [Intent Engineering test strategy](https://intent-engineering-for-coding-agents.github.io/book/quality/test-strategy.html): **test type** answers *what a test verifies at which boundary* (unit, integration, E2E, …) and **level** answers *when it runs* (pre-commit, pre-merge, post-deploy). An AC may legitimately need verification at several test types, possibly with different tools. Both type and level have the **runner/pipeline** as their consumer, not the methodology: their tags exist so pipelines can filter (fast tests on every PR, slow E2E nightly). The rule that avoids redundancy: a test carries **exactly one purpose tag, which `clue` reads, and any number of runner tags (type, level, …), which `clue` ignores**. `@Tag("AC-022")` and `@Tag("integration")` on the same test is correct and not redundant — each tag has exactly one consumer. In Go, purpose lives in the name prefix and type/level use Go's own runner idioms (`//go:build integration`, `testing.Short()`).

An AC may have several tests with different runner metadata; whether that set adequately covers the criterion remains review judgment. Integration and E2E are runner test types, not purposes. Performance tests trace to a specific quality constraint rather than becoming a purpose class. A `QS<digits>` purpose is a future door only if a QS-verifying test exists.

Leaving non-AC tests unclassified is rejected because it makes the purpose check one-sided and lets intent-free tests accumulate.
