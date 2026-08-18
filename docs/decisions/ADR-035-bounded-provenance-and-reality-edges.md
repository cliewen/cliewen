---
id: ADR-035
type: decision
status: verified
links: [ADR-010, PDR-004, CAP-002, CAP-003, AN-007, P-007]
title: Cost bounds inferred provenance and incident analyses return an edge from reality
author: agent
accepted-by: Flemming N. Larsen (2026-08-02, conversation)
---

# ADR-035 — Bounded provenance and reality edges

## Context and problem statement

Inferred non-decision artifacts can be cheap or expensive to reverse, but a single undifferentiated inferred count cannot distinguish safe deferral from a costly ambiguity or connect an observed failure to the capability whose claim was wrong.

## Decision outcome

**An inferred non-decision declares `reversal-cost: low | high`; high-cost inferred meaning blocks activation only within the active capability's one-edge link slice.** The field is required while `provenance: inferred` and optional after verification; ADRs and PDRs do not carry it because their type already represents the high-cost route. Low-cost inferred artifacts remain valid. The CLI reports high-cost inferred non-decisions as activation blockers and inferred ADRs/PDRs as non-blocking decisions awaiting verification; the old count of every inferred non-decision is removed.

**An analysis records an observed contradiction with `reality: contradicted` and links every affected capability or live acceptance criterion.** `clue validate --reality-gaps` derives a sorted capability view from those edges, and no incident registry or production feedback loop is added.

**Carrier:** provenance and active-slice checks, acceptance-criterion link resolution, the derived reality-gap view, and the `clue-analysis` incident convention.
