---
id: ADR-032
type: decision
status: verified
links: [ADR-005, ADR-006, CAP-002, CAP-003, P-007]
title: Acceptance criteria declare classified proof and paired directions
author: agent
accepted-by: Flemming N. Larsen (2026-08-02, conversation)
---

# ADR-032 — Acceptance criteria declare classified proof and paired directions

## Context and problem statement

An unclassified test reference proves only that an acceptance-criterion ID occurs in a test suite; it does not prove which boundary or direction was exercised. `checkACTests` therefore needs a scenario-level declaration that it can compare with named evidence.

## Decision outcome

> **Partially superseded by [ADR-036](ADR-036-jvm-evidence-per-executable.md) and extended by [ADR-037](ADR-037-brownfield-ac-id-grammar.md):** proof classes, paired directions, Go names, and Cucumber scenario tags remain current; JVM evidence is attributed to one executable, and segmented or letter-suffixed IDs use ADR-037's carrier normalization.

New or materially revised classified scenarios declare `Test-type: Unit`, `Integration`, `E2E`, or `Performance` on their first non-blank line. Each declared class requires one `positive` and one `negative` reference, except an explicit `(single-direction)` declaration. Go names, JVM executable tags, and Cucumber scenario tags carry the AC identity, class, and direction; profiles without native tags use the stable named-executable fallback, while proximity comments remain unsupported. Unannotated legacy scenarios retain ADR-006's one-reference rule.

The purpose taxonomy remains `AC`, `Unit`, `Sanity`, and `Arch`; test type and direction are evidence metadata consumed by `checkACTests`, not new test purposes. A missing class or direction is diagnosed rather than inferred.

**Carrier:** `checkACTests`, the supported Go/JVM/Cucumber evidence harvesters, CAP-002 and CAP-003 criteria, and the generated guidance that describes the declaration.
