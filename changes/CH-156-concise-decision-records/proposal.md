---
id: CH-156
type: change
status: active
links: [PDR-003, PDR-006]
title: Plan concise typed decision records
---

# CH-156 — Plan concise typed decision records

## Proposal

Create a dedicated campaign plan for replacing the decision log with subject-typed ADRs, PDRs, and IDRs, making new and modified records concise by default, migrating older adopter repositories safely, and compacting this repository's existing decision corpus in reviewable batches.

The accepted method currently routes cheap decisions to `docs/decisions/log.md` and uses reversal cost before subject. The requested model removes that artifact, routes durable decisions by subject, and records only future-shaping choices. Because this changes methodology meaning and several adopter-facing carriers, implementation needs staged full changes rather than one unreviewable rewrite.

This change will record the campaign decision, create the plan with independently verifiable milestones, and regenerate the plan and decision indexes. It is plan-less because it creates the plan that subsequent changes will serve.

## Scope boundary

This change plans and authorizes the campaign only. It does not yet change validation, scaffolding, migration, extraction, generated skills, or existing decision records beyond the new campaign decision.
