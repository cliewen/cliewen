---
id: PDR-034
type: decision
status: verified
links: [CAP-007, ADR-013, G-001, C-013, P-013, AN-018]
title: The corpus is read narrowly by default, and widened only when the work requires it
author: agent
accepted-by: Flemming N. Larsen (2026-08-08, conversation)
---

# PDR-034 — The corpus is read narrowly by default

## Context and problem statement

The corpus grows while `clue context` can already emit a deterministic artifact slice, but no decision required starting with that slice instead of loading unrelated history.

## Decision outcome

**Reading starts at the narrowest point that answers the question and widens only when the work or a discovered edge requires it.** A known identity enters through `clue context <id>` and its outgoing-link slice; the index is used only to find an entry point when no identity resolves; a binding edge, conflict, or necessary artifact justifies widening. This sets the starting point, not a cap on what an agent may read, and adds no mechanism to CAP-007.

## Rejected: leave narrow reading to judgment or cap the corpus

The first leaves an existing capability optional and makes defensive reading the default; the second would forbid work that genuinely needs a wider slice.

## Carrier

The repository and scaffolded routing hubs state the three-step reading rule, CAP-007 carries the command, and this record is the live trace for those statements.
