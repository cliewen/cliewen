---
id: PDR-003
type: decision
status: verified
links: [ADR-010]
title: ADRs are for expensive-to-reverse decisions; the rest is a decision log
author: agent
accepted-by: Flemming N. Larsen (2026-07-14, PR #9)
---

# PDR-003 — ADRs for the expensive-to-reverse; a decision log for the rest

> **Superseded by [PDR-006](PDR-006-decision-records-are-typed.md) and [PDR-046](PDR-046-decisions-route-by-subject.md):** the reversal-cost and decision-log routing was first split by subject and is now replaced by future-shaping subject-typed records; the provenance and retention principles survive.

## Context and problem statement

Putting cheap defaults beside architecture in one ADR series hides which choices constrain future work and discourages recording smaller decisions. The corpus needed a cheap, repository-native record without weakening retention or provenance.

## Decision outcome

**The original routing used reversal cost as its litmus test: expensive-to-reverse decisions were ADRs, while cheap and local decisions were dated rows in `docs/decisions/log.md`.** The log was a linted corpus artifact with `Date | Decision | Why | Change/PR` columns; demotion preserved the decision as a row and promotion moved a constraining row into a full record.

PDR-006 retained the need to distinguish architecture from process, and PDR-046 superseded the reversal-cost/log clauses by routing future-shaping choices to ADR, PDR, or IDR records and retiring the log. The durable lesson is that decisions remain repository-native and are not silently discarded when their carrier changes.
