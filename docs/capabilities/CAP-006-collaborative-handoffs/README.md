---
id: CAP-006
type: capability
status: active
links: [G-001, PDR-016, C-012]
title: Collaborative PR handoffs remain visible and exact across agents
goal: G-001
---

# CAP-006 — Collaborative PR handoffs

## What

Agents may independently implement, review, and repair changes without losing work between sessions: review results name the hosted commit they cover, actionable findings remain visible on the pull request, and every agent that edits completes an exact, non-destructive publication handoff before calling the pull request merge-ready.

## Why

The human merge boundary of [G-001](../../goals/G-001-verifiable-thread.md) can accept only hosted state. A private finding or local fix is invisible to the next agent and to the human, while a global change lock would needlessly serialize independent work.
