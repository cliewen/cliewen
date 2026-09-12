---
id: CAP-011
type: capability
status: active
links: [G-012]
title: Repository orientation
goal: G-012
---

# CAP-011 — Repository orientation

## What

Cliewen gives a newly started agent a small, read-only account of where repository work stands and which recorded choices can follow. `clue next [--all] [path]` derives that account from open change workspaces, unfinished plan milestones, and proposed goals without claiming work or storing private session state.

The routing hub asks an agent to inspect this account once after the mandatory release check in a fresh context. The agent reports it only when the opening request leaves direction open or the existing work affects that request, then reads bounded context before recommending a choice.

## Why

Serves [G-012](../../goals/G-012-corpus-states-the-product-direction.md): a fresh agent should not need the user to reconstruct repository position from memory, especially between campaigns when no active milestone exists.

Acceptance criteria: [criteria.md](criteria.md) · implementation notes: [design.md](design.md).

## Status note

`active`: the orientation ladder and conditional fresh-context guidance are implemented with focused unit evidence.
