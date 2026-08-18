---
id: ADR-010
type: decision
status: verified
links: [AN-002, ADR-008, CAP-003]
title: Extracted artifacts carry a provenance field, born inferred
author: agent
accepted-by: Flemming N. Larsen (2026-07-12)
---

# ADR-010 — The provenance field

## Context and problem statement

Extraction creates non-decision artifacts whose lifecycle may already be active even though an agent reconstructed their meaning. Reusing `status` for both facts would either misclassify lifecycle or hide the lack of human review.

## Decision outcome

**Keep provenance and lifecycle separate.** Non-decision artifacts may carry `provenance: inferred | verified`; extraction writes `inferred`, and human review promotes it. Decisions continue to express provenance through `status: inferred | verified` and do not carry the separate field.

[ADR-035](ADR-035-bounded-provenance-and-reality-edges.md) bounds the inferred state: non-decisions declare cheap or expensive reversal, and expensive inferred meaning blocks activation in its immediate graph slice. The CLI reports those blockers separately from decisions awaiting verification.

**Carrier:** provenance vocabulary and reversal-cost lint in `clue`, plus the born-inferred rule in `clue-extract`. Adding `inferred` to every lifecycle vocabulary or relying only on transient PR review is rejected because each conflates axes or loses durable evidence.
