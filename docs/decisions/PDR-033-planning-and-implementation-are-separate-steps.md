---
id: PDR-033
type: decision
status: verified
links: [PDR-017, PDR-007, PDR-019, C-003, C-012, G-001, P-013, AN-018]
title: Planning and implementation are separate steps, and the boundary between them is a human decision
author: agent
accepted-by: Flemming N. Larsen (2026-08-08, conversation)
---

# PDR-033 — Planning and implementation are separate steps

## Context and problem statement

The change loop runs from proposal to implementation without stopping. That is right for most changes and wrong for the ones where the proposal is the work: a plan that turns out to need a second plan beside it, a proposal a human wants to read before anything is built, or work that should be implemented by a different agent — or a different person — than the one that designed it.

[PDR-017](PDR-017-merge-gate-has-content.md) decided that pause and states it in one line: a full change may opt into a spec-first pause after Propose, recording it in tasks and waiting for human direction, with the ordinary loop as the default. What it settles is that the pause exists and is opt-in. What it does not settle is what the pause is *for*, or what happens while the work is stopped — and read without that, the clause is a convenience, which is how it has been used.

The purpose is not convenience. The proposal and the implementation are two different kinds of work, and the moment between them is the cheapest place in the loop to change direction, split the work into more than one plan, or hand it to a different agent or a different person. A pause that nobody reports at, and that ends without anyone being asked, delivers none of that.

Whether the branch is pushed at that moment is a separate question with a real cost either way, and it is not one an agent should answer on the human's behalf: an unpushed branch cannot be picked up by anyone else, and a pushed branch is unfinished work visible to everyone and can no longer be rebased onto a newer `main` without rewriting hosted history.

## Decision outcome

**The boundary between proposing and implementing is a step in the loop, and the human decides whether to cross it.** This amends [PDR-017](PDR-017-merge-gate-has-content.md) by stating what the pause is for and what the agent does while it holds; PDR-017's rules that the pause is opt-in, is recorded in tasks, and leaves the ordinary loop as the default are unchanged.

*The proposal is committed before the pause.* This is the loop's existing rule and it is what makes the pause useful rather than merely polite: a proposal that exists only in a conversation cannot be read later, revised by someone else, or implemented by a different agent. The commit is what turns a plan into something that can change hands.

*At the pause the agent reports briefly and asks.* A short status — what the proposal says and what implementing it would involve — and then two questions: **should implementation begin**, and **should the branch be pushed**. Not a summary of the proposal, which the human can read; a statement of where the work stands and what happens next.

*Pushing is asked, never assumed.* Pushing makes the work recoverable and hands it to whoever comes next; it also publishes unfinished work and ends the branch's freedom to be rebased onto a newer `main` under [PDR-007](PDR-007-review-boundary.md). Both are real, the trade differs by change, and the human owns it. An unpushed pause is preserved local work, exactly as [C-012](../constraints/C-012-agents-never-merge-own-changes.md) already describes a human-requested stopping point.

*The pause is recorded, not remembered.* It goes in `tasks.md` under [C-003](../constraints/C-003-tasks-tick-immediately.md), so a later reader — or a later agent — can see that the work stopped deliberately rather than being abandoned.

*Without an opt-in the loop proceeds directly to implementation.* This is PDR-017's rule and it stands: making the pause mandatory would charge every ordinary change for a handoff that almost never happens.

## Rejected: push on every commit

Considered and withdrawn in the same conversation that produced this record. It guarantees the work survives a lost machine and makes handoff always possible, and it costs the unpublished-branch rebase that [PDR-007](PDR-007-review-boundary.md) permits: a branch published from its first commit can only take a newer `main` by merging it in. Asking at the pause buys the same handoff where handoff is actually wanted, and leaves the ordinary change's history alone.

## Rejected: treat the pause as an agent's judgement call

An agent that decides for itself when a proposal is worth stopping on will stop either always or never, and neither is right. The signal that a proposal needs a second reader is one the human has and the agent does not. PDR-017 rejected making the pause the default for the neighbouring reason, and that rejection is untouched here.

## Carrier

No carrier outside this record states the reporting and asking rules today. The pause itself, in PDR-017's narrower form, is stated by the `clue-delta` source template and its generated and scaffolded copies, and by `guide/change-loop.md`. All of them are owed this amendment, and [P-013](../plans/P-013-simplification.md) splits them by surface: the skill carriers belong to M-063, which owns every edit to a skill in that campaign, and the guide belongs to M-064, which carries the public guide and the rest of the carrier prose M-063 does not reach.

PDR-017's remaining carriers — `clue-verify`, the shared `durable-work` fragment, the pull-request template, the reusable workflow and its thin caller, and CAP-006's criteria — carry its acceptance-brief and scenario-resolution rules and state nothing about the pause, so this amendment does not reach them. That was checked rather than assumed, because [PDR-019](PDR-019-methodology-contract-carriers-move-together.md)'s inventory obligation runs in both directions.
