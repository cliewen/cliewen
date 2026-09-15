---
id: G-020
type: goal
status: proposed
links: []
title: A goal's current service level is observable, not just its acceptance
---

# G-020 — A goal's current service level is observable, not just its acceptance

**Who wants it:** an adopter running brownfield extraction (2026-09-15), found while an agent answered "what's the status of our goals" in an adopted repository and could only say that G-001/G-002/G-003 all stay `accepted` regardless of how well their capabilities currently serve them, because [ADR-025](../decisions/ADR-025-one-status-lifecycle.md) gives `goal` only `proposed → accepted` — modeled on the vision's durable-want shape, not a plan's `→ completed` delivery shape.

`accepted` conflates two different facts: "the corpus endorses this want as legitimate" and "work satisfying this want is currently done." A goal proposed yesterday with no capability yet, and a goal whose capabilities have been active and evidenced for months, render identically. Making it worse, [CAP-007](../capabilities/CAP-007-focused-context/README.md)'s `clue context <goal-id>` deliberately never follows reverse links from a goal to the capabilities that serve it — by design, to avoid pulling most of a mature corpus into one result — so there is today no mechanized way, even read-only, to answer "how well is this goal currently served," let alone a place to record it.

Any fix should preserve the reason a goal is not a plan: a want does not stop being wanted merely because it is well served right now, and a status field that looks like "done" risks being read as "no longer needs attention" when a regression could unserve it again later. A derived, non-authoritative service-level view (e.g., a `clue context`-adjacent report over the goal's inbound capability links) may fit better than a new status value that a human must remember to revise.
