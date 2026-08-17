---
id: PDR-048
type: decision
status: inferred
links: [PDR-026, PDR-033, C-003, C-010]
title: Plan and workspace bookkeeping stays truthful
author: agent
accepted-by: []
---

# PDR-048 — Plan and workspace bookkeeping stays truthful

## Context

Tasks and milestones become misleading when an infeasible item has no reason, a digest is listed as work that can finish before its own deletion, or a milestone status cannot be checked against a stated outcome.

## Decision

An infeasible task carries its reason on the same line. The digest is never a task. A milestone uses `todo`, `doing`, `done`, or `dropped`, states a verifiable exit criterion, and gives a dropped reason in its evidence cell. The digest completing a plan's last milestone closes that plan and names any successor rather than leaving bookkeeping for a later change.
