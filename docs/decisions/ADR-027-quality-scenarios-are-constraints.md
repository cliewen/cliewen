---
id: ADR-027
type: decision
status: verified
links: [P-005, ADR-025, ADR-026, ADR-017, C-014, C-015]
title: Quality scenarios are constraints — the quality type folds into the register
author: agent
accepted-by: Flemming N. Larsen (2026-08-02, conversation)
---

# ADR-027 — Quality scenarios are constraints

## Context and problem statement

Quality scenarios and constraints both express cross-cutting rules checked against proposals; keeping a separate `quality` type and folder adds no distinct processing.

## Decision outcome

**Quality scenarios are constraints.** QS-001 becomes [C-014](../constraints/C-014-coverage-floor.md) for total Go statement coverage of at least 80%, and QS-002 becomes [C-015](../constraints/C-015-onboarding-under-30-minutes.md) for the onboarding bound. The `quality` folder and type are removed; the old IDs are not reused, and Git history preserves their provenance.

Immutable prior texts that mention QS-001 or QS-002 are read as pointers to C-014 or C-015; this includes ADR-006, PDR-012, and completed plans. The constraints register is the carrier, along with removal from init templates and taxonomy prose. Keeping a quality subtype is rejected because it duplicates the register; demoting the rules to log rows is rejected because standing constraints need `source` and `enforcement`.
