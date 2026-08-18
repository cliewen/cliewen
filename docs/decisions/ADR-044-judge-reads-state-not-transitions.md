---
id: ADR-044
type: decision
status: verified
links: [P-010, CAP-002, ADR-017, ADR-040, ADR-042, C-002, C-004, C-006, C-008, C-012, C-013]
title: The judge judges a repository state, never a transition
author: agent
accepted-by: Flemming N. Larsen (2026-08-08, conversation)
---

# ADR-044 — The judge reads state, not transitions

## Context and problem statement

Some rules concern what a change did, what it is based on, or what a forge will merge, but a deterministic corpus judge must be reproducible from the repository it is given rather than from history, a diff, or an external base revision.

## Decision outcome

**`clue validate` reads the repository as it is and never reads history, a diff, or a base revision; no transition rule becomes its verdict.** CI, branch protection, release workflows, and other named machines may compare transitions when that is their purpose, and the register identifies which machine owns each part.

Meaning that no machine can hold is declared with its human judgment and exposure rather than queued as a future check. This includes whether a deletion weakened a rule, whether a decision is enduring, and whether a change crosses the core meaning boundary.

**Carrier:** the judge's input boundary, the enforcement register, and the workflow/forge guidance that owns transition checks.
