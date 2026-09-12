---
id: PDR-058
type: decision
status: verified
links: [CAP-011, AC-190, AC-191, PDR-057, G-012]
title: Fresh-context orientation is derived and conditionally disclosed
author: agent
accepted-by: Flemming N. Larsen (2026-09-12, planning conversation)
---

# PDR-058 — Fresh-context orientation is derived and conditionally disclosed

## Context

PDR-057 made active plan milestones visible when a user asks what is next, but a fresh agent still begins without repository position and the command reaches a dead end between campaigns. Persisting a personal focus would add a stale source of truth and would conflict with Cliewen's boundary against storing product state outside the repository.

## Decision

The routing hub asks every fresh agent context to inspect `clue next --all` once after the mandatory release check. The agent discloses a brief status only when the opening request leaves direction open or existing work affects it; a concrete request receives no routine status preamble. Before recommending a recorded choice, the agent reads its bounded context. Inspection never authorizes starting proposed work.

`clue next` derives repository position in this order: an open change workspace to resume, unfinished active milestones with `doing` before `todo`, unfinished draft milestones, and proposed goals. The last two remain explicitly non-actionable. With no recorded candidate, the command points to capturing a proposed goal. The command stores no focus, user identity, or session state.

This contract is carried by CAP-011 and its criteria and design, `cmd/clue` and `internal/corpus`, the repository and scaffold routing hubs, the canonical durable-work skill source and its generated copies, the plans and cross-cutting design overviews, public prompting and operating guidance, and the release notes.
