---
id: CH-072
type: change
status: open
links: [P-007, M-025, ADR-006, ADR-005, CAP-002, CAP-003]
title: Classify acceptance-criterion evidence and count it
---

# CH-072 — Classify acceptance-criterion evidence and count it

M-025 makes the existing test-purpose taxonomy sufficient to express what proof each acceptance-criterion scenario requires, then makes `clue validate` check the declared proof. It replaces ADR-006's deferred per-AC test-type annotation with a coverage-pair convention that declares proof classes beside each scenario, requires positive and negative evidence unless a declared single-direction exemption applies, and keeps the declaration visible at specification review.

The change will record the vocabulary and its language-specific evidence carriers in an ADR, extend `checkACTests` to harvest classified Go, JVM, and Cucumber evidence, and resolve the existing Gatling/no-tag-mechanism and QS-lane doors explicitly. Its new validation behavior will be specified by acceptance criteria and exercised with positive and negative tests.

## Scope

- Extend the durable test-purpose and evidence-classification decisions.
- Define and parse per-scenario proof-class declarations in `criteria.md`.
- Enforce required coverage classes and the positive/negative pair convention in `checkACTests`.
- Harvest Cucumber `.feature` tags alongside existing Go and JVM evidence.
- Update the corpus, capability contracts, source skills, generated skills, and user-facing changelog where they describe the former limit.

## Out of scope

- The human proof class and per-criterion draft exemption belong to M-026.
- Retiring corpus artifacts, bounded provenance, and task-proportionate context remain later P-007 milestones.
- General test execution or runner-level selection is not changed.
