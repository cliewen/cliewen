---
id: PDR-031
type: decision
status: verified
links: [P-013, PDR-029, PDR-013, C-013, G-001, AN-018]
title: An architecture artifact is a valid trace when it states the rule
author: agent
accepted-by: Flemming N. Larsen (2026-08-08, conversation)
---

# PDR-031 — An architecture artifact is a valid trace when it states the rule

## Context and problem statement

PDR-029's carrier test named decisions, constraints, goals, and acceptance criteria as trace targets, but the routing hub's core definition is stated by ARCH-003 and should not be duplicated into a decision merely to satisfy the register.

## Decision outcome

**An architecture artifact is a valid trace when it states the rule the carrier carries, not when the rule is merely derivable from its description.** The accepted trace set is therefore decision, constraint, goal, acceptance criterion, and architecture; the register still names the narrowest artifact that states the rule. PDR-029's two tests, surface split, ordering, overlap, and anti-derivability boundaries are otherwise unchanged.

## Rejected: restate ARCH-003 as a decision

That would duplicate architecture and create a second carrier required to move with it under PDR-019 without adding meaning.

## Rejected: accept every corpus artifact as a trace

Plans and analysis describe practices without deciding them, so admitting them would dissolve the bounded trace test.

## Carrier

This record amends PDR-029's trace-type clause; AN-018 and P-013's M-062 use the widened set, and no milestone work changes merely because the accepted trace set widened.
