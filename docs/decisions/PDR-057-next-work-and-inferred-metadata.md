---
id: PDR-057
type: decision
status: inferred
links: [AC-186, CAP-002, ADR-035, PDR-056]
title: Next work is derived from active plans and reversal cost ends at verification
author: agent
accepted-by: []
---

# PDR-057 — Next work is derived from active plans and reversal cost ends at verification

## Context

An agent asked what to do next can stop at the corpus index even when active milestones remain, while a draft plan's unfinished rows can be mistaken for authorized work. Inferred artifacts also retain a reversal-cost label after their provenance is verified, where it no longer changes any judgement.

## Decision

`clue next` is the deterministic read-only orientation path: unfinished `doing` milestones precede `todo` milestones in active plans, and draft-plan rows are reported as proposed rather than actionable. The command exposes alternatives without claiming or reassigning work. `reversal-cost` is required only for inferred non-decision meaning; `high` remains the activation blocker and `low` explicitly permits deferral, but verified non-decisions must not carry the field. Migration removes that obsolete field without rewriting surrounding content.
