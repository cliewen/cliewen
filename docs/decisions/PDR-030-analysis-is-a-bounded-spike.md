---
id: PDR-030
type: decision
status: verified
links: [P-013, PDR-029, PDR-006, ADR-035, C-011, G-001, AN-018]
title: Analysis is a bounded spike that ends in a findings document with a named consumer
author: agent
accepted-by: Flemming N. Larsen (2026-08-08, conversation)
---

# PDR-030 — Analysis is a bounded spike that ends in a findings document with a named consumer

## Context and problem statement

The analysis skill already required a scoped investigation, disposable work, durable findings, rejected options, and a consumer, but those rules had only been inherited from a frozen foundation document and had no live trace.

## Decision outcome

**Analysis is a bounded investigation with a stated risk, disposable work, a durable findings document, and a named consumer; the four parts form one rule.** It runs before the planning or implementation whose shape the risk could change, states the risk in one sentence, throws away the prototype or measurement, and keeps a findings document under `/docs/analysis/` that records what was rejected. The consumer is a plan or change; without one, the analysis is not written.

This record supplies the live trace for the `clue-analysis` skill and does not change the workflow it records.

## Rejected: make each clause its own decision

Each clause depends on the others, so separate records would preserve neither the coherent boundary nor a meaningful reversible choice.

## Rejected: promote the frozen foundation document or delete the statements

The foundation cannot carry live rules, while deleting untraced statements would discard a practice already used by the analysis corpus; the repair is this decision.

## Carrier

The canonical and generated `clue-analysis` skills carry the workflow; this record carries its rationale, PDR-006 carries the separate rule for rejected decisions, and ADR-035 remains the `reality: contradicted` rule.
