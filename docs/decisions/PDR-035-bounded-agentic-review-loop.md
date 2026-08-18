---
id: PDR-035
type: decision
status: inferred
links: [P-013, PDR-012, PDR-016, PDR-019, PDR-029, C-004, C-012, C-017, ADR-021]
title: The agentic review loop owns its severity model and has a bounded ordinary budget
author: agent
accepted-by: []
---

# PDR-035 — The agentic review loop owns severity and budget

> **The budget is amended by [PDR-036](PDR-036-review-loop-budget-and-human-checkpoint.md):** the former three-pass ordinary budget is a maximum stated through C-017 and the generated skill, not a floor; reaching it with blocking findings stops the loop, reports the remainder, and asks the human whether to continue. The severity, exact-commit, and advisory-handoff rules below remain.

## Context and problem statement

PDR-012 required adversarial review but left its severity and pass cost vulnerable to each reviewer brief, risking both a weakened blocking gate and an unbounded repair cycle.

## Decision outcome

**The review loop owns its blocking/advisory classification.** A caller may state risks and intent but cannot redefine severity; computed figures and arithmetic disagreements are advisory, while wrong, missing, or reused identities and false normative claims remain blocking. The loop exits only when the current commit has no blocking findings, keeps the exact-reviewed-commit boundary, and reports its mode, commit, pass count, and advisory findings left open. A blocking repair earns another pass; an advisory does not become a gate or silently alter a clean candidate.

## Rejected: let each brief choose severity or make figures blocking

That would make publication depend on prompt wording and spend the repair loop re-deriving arithmetic instead of judging the operative requirement.

## Carrier

PDR-012, PDR-016, C-012, C-017, CAP-006, the review-boundary and `clue-verify` sources, generated skill trees, contributor and guide text, and pull-request templates carry the classification and handoff; PDR-036 carries the maximum-pass amendment.
