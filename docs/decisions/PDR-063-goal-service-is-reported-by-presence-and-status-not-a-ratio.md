---
id: PDR-063
type: decision
status: inferred
links: [G-020, ADR-025, PDR-054, ARCH-003, CAP-009]
title: A goal's current service is reported by capability and plan presence, never a new status or a ratio
author: agent
accepted-by: []
---

# PDR-063 — A goal's service is reported, never scored

## Context

G-020 found that an adopter's agent, asked about goal status during brownfield extraction, could only say `proposed` or `accepted` — [ADR-025](ADR-025-one-status-lifecycle.md)'s only goal states — regardless of how well a goal was actually served, because nothing in the corpus or its tooling answers that separately. [CAP-007](../capabilities/CAP-007-focused-context/README.md)'s `clue context` deliberately never follows reverse links from a goal, so there was also no mechanized way to see this by hand.

Two familiar-looking fixes were available and both are rejected here.

## Decision

**`clue validate --intent` states, per goal, the capabilities whose required `goal:` field names it and the plans whose `links` name it, each shown with its own status.** A goal named by neither prints that plainly. This reuses two existing, already-bounded one-hop edges — capability's required `goal:` field, and the corpus's existing practice of a plan's `links` naming the goal it serves — so the report stays derived, is never stale, and adds no new required field or corpus obligation.

Both threads matter: a capability is the intent thread's ongoing capacity, while a plan reaching `completed` is the delivery thread's existing signal that a one-time achievement finished. A goal such as "Cliewen is public" is served by a completed plan and no standing capability at all; reporting capabilities alone would have reported it as unserved. This was found by running the extended report against this repository during CH-193's implementation, not decided in the abstract.

**Rejected: a new goal status** (e.g. `served` or `completed`). A goal is a durable want, modeled like the vision rather than like a plan's delivery-shaped `draft → active → completed` (ADR-025): the same want does not stop being wanted merely because it is well served today, and a regression could unserve it again without anyone remembering to revert a status that looked like "done." `proposed → accepted` is unchanged.

**Rejected: a coverage percentage or ratio.** [PDR-054](PDR-054-use-cases-are-optional-and-no-requirement-artifact.md) already rejected this shape for use cases: a percentage reads as a target, and the only way to move it is to write artifacts nobody needs. The same applies to goals, capabilities, and plans; the report lists identities and statuses and computes no figure.

## Carrier

`internal/corpus/intent.go` (`GoalState`, `GoalRef`, `Intent`), the `--intent` printer in `cmd/clue/main.go`, and [AC-202](../capabilities/CAP-009-product-intent/criteria.md) alongside AC-164's existing no-ratio rule.
