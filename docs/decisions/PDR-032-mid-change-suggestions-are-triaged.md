---
id: PDR-032
type: decision
status: verified
links: [ADR-002, PDR-017, PDR-007, C-003, C-005, G-001, P-013, AN-018]
title: A suggestion raised mid-change is triaged into the current work or into the inbox, never into memory
author: agent
accepted-by: Flemming N. Larsen (2026-08-08, conversation)
---

# PDR-032 — A suggestion raised mid-change is triaged into the current work or into the inbox, never into memory

## Context and problem statement

A suggestion can arrive after a change starts, while its workspace is transient and conversation is not durable; silently absorbing it changes the reviewed scope, while silently forgetting it loses the human's direction.

## Decision outcome

**Every suggestion is triaged immediately into one of two carriers, and the choice is stated to the human.** If the change is wrong or incomplete without it, it becomes a task in `tasks.md` and is handled before merge. Otherwise it becomes a `status: proposed` goal written in the digest, so it survives workspace deletion and can be accepted later through its own change. There is no third carrier called memory, and neither an agent's silent rejection nor an unrecorded out-of-scope judgment is allowed.

## Rejected: sort suggestions at digest

The digest deletes the workspace, so deferred sorting can destroy an overlooked suggestion at the same moment it is supposed to preserve it.

## Rejected: create a new backlog artifact or absorb adjacent work

ADR-002 already makes proposed goals the inbox, while absorbing work after authorization changes the scope no reviewer approved.

## Carrier

The goals README carries ADR-002's inbox; the shared `durable-work` fragment and `clue-delta` routing carry the triage; C-021 states the human-enforced residual that `clue` cannot observe suggestions in conversation.
