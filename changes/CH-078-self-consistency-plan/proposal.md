---
id: CH-078
type: change
status: open
links: []
title: Open P-008 to make Cliewen's methodology carriers agree
---

# CH-078 — Open P-008 to make Cliewen's methodology carriers agree

## What

Create [P-008](../../docs/plans/P-008-self-consistency.md) `active`, a sequential campaign that first repairs the evidence-model contradictions found after v0.9 and then consumes the evidence-backed analysis gaps that remain outside completed campaigns. Record the audit as [AN-011](../../docs/analysis/AN-011-methodology-self-consistency.md), including the edge from a green corpus back to the public-guide and extraction claims it contradicted. The plan carries milestones M-032…M-036, continuing corpus-global numbering from P-007. This change writes the analysis, plan, their index rows, and one decision-log row recording the campaign's scope and order. It implements none of the repairs and is plan-less because its product is the plan.

## Why

The v0.9 validator, tests, generated-skill drift checks, and guide build are green, but several live carriers still state the pre-v0.9 evidence contract: three guide pages say the judge cannot classify or count positive and negative evidence, the overview omits Cucumber and the Human and per-criterion draft cases, and `clue-extract` says activation is only per criteria file. ADR-032, ADR-033, CAP-002, and `clue-delta` state and implement the newer contract. This is the exact semantic-consistency edge the foundation acknowledged as unmechanized: every individual artifact passes its form checks while the set disagrees about what Cliewen means.

The audit also shows that AN-003's clean-environment, population-evidence, and adoption-governance findings remain only partly consumed, while P-007 explicitly leaves extraction rehearsal and several adopter-integration boundaries for later evidence. P-008 orders the demonstrated repair first, implements the analysis protocol that already has cross-trial support second, adds the non-destructive extraction boundary third, and requires real adopter evidence before committing interfaces for infrastructure and distributed work.

## Decision boundary

This change decides only the campaign scope and ordering. It changes no validator behavior, generated skill, guide statement, constraint, acceptance criterion, or public command. M-032 is a separate full change that records the lasting carrier-consistency rule, updates every stale evidence carrier, adds regression coverage, and supplies the user-facing changelog entry. M-033 and M-034 alter agent workflow and therefore each receive their own full change and decision routing. M-035 and M-036 are analysis milestones whose findings must precede any interface decision. Production operations, external constraint catalogs, and validation of foreign documentation kinds remain named future doors because the current evidence does not justify promising their implementation.
