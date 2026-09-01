---
id: PDR-033
type: decision
status: verified
links: [PDR-017, PDR-007, PDR-019, PDR-040, C-003, C-012, G-001, P-013, AN-018]
supersedes: [AN-019]
title: Planning and implementation are separate steps, and the boundary between them is a human decision
author: agent
accepted-by: Flemming N. Larsen (2026-08-08, conversation)
---

# PDR-033 — Planning and implementation are separate steps

> **Amended by [PDR-040](PDR-040-push-is-durability-ready-is-explicit.md):** every changed turn pushes its commit and the proposal is already visible on the draft PR at the pause, so the pause asks only whether implementation should begin; the opt-in, task record, and ordinary-loop default remain.

## Context and problem statement

PDR-017 permits an opt-in spec-first pause but did not state its purpose or what the agent must ask, leaving the boundary between a proposal and its implementation easy to cross without a human decision.

## Decision outcome

**The boundary between proposing and implementing is a step in the loop, and the human decides whether to cross it.** The proposal is committed before the pause; the agent reports briefly and asks whether implementation should begin; the pause is recorded in `tasks.md`; and without an opt-in the ordinary loop proceeds directly to implementation. PDR-040 removes the former push question because publication is now durable on every changed turn.

## Rejected: make the pause an agent judgment

An agent cannot reliably decide which proposal needs a second reader, and making the pause mandatory would charge ordinary changes for a handoff they do not need.

## Carrier

PDR-017 and the canonical/generated `clue-delta` source carry the opt-in and task record; `guide/change-loop.md` carries the public explanation; PDR-040 carries the publication amendment.
