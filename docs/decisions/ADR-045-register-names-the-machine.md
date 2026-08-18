---
id: ADR-045
type: decision
status: verified
links: [P-010, CAP-002, ADR-017, ADR-044, C-004, C-011, C-013]
title: Every constraint names the machine that holds it or the judgment that remains
author: agent
accepted-by: Flemming N. Larsen (2026-08-08, conversation)
---

# ADR-045 — The register names the machine

## Context and problem statement

The constraint register previously forced rules that were partly mechanical or permanently judgmental into inaccurate `machine` or `agent` categories, hiding what a machine actually checks and what a human still accepts.

## Decision outcome

**Constraint enforcement has four classes: `machine`, `partial`, `agent`, and `human`.** A `partial` constraint names the machine and exact subset it checks plus a `Residual` describing the judgment and exposure; a `human` constraint names the machine fragment when one exists and always declares its residual. The validator requires the fields owed by each class.

`agent` keeps its existing meaning: a rule awaiting a real machine check, with a promotion trigger and a visible backlog count. `human` now means any permanent property no machine can hold, whether the remaining judgment is exercised by a person or agent; it is a finished declaration, not a smaller `agent` backlog.

**Carrier:** the constraint register schema and validator, the enforcement descriptions, and the machine/human evidence named by each constraint.
