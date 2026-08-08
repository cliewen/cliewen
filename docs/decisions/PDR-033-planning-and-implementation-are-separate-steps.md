---
id: PDR-033
type: decision
status: verified
links: [PDR-007, PDR-008, C-003, C-005, C-012, G-001, P-013, AN-018]
title: Planning and implementation are separate steps, and the boundary between them is a human decision
author: agent
accepted-by: Flemming N. Larsen (2026-08-08, conversation)
---

# PDR-033 — Planning and implementation are separate steps

## Context and problem statement

The change loop runs from proposal to implementation without stopping. That is right for most changes and wrong for the ones where the proposal is the work: a plan that turns out to need a second plan beside it, a proposal a human wants to read before anything is built, or work that should be implemented by a different agent — or a different person — than the one that designed it.

The skill has carried an opt-in pause at that point since early in the corpus, and no decision ever recorded it. Read as written it is a convenience — *a human may opt into a spec-first pause* — which understates what it is for. The proposal and the implementation are two different kinds of work, and the moment between them is the cheapest place in the loop to change direction, split the work, or hand it to someone else.

Whether the branch is pushed at that moment is a separate question with a real cost either way, and it is not one an agent should answer on the human's behalf: an unpushed branch cannot be picked up by anyone else, and a pushed branch is unfinished work visible to everyone and can no longer be rebased onto a newer `main` without rewriting hosted history.

## Decision outcome

**The boundary between proposing and implementing is a step in the loop, and the human decides whether to cross it.**

*The proposal is committed before the pause.* This is the loop's existing rule and it is what makes the pause useful rather than merely polite: a proposal that exists only in a conversation cannot be read later, revised by someone else, or implemented by a different agent. The commit is what turns a plan into something that can change hands.

*At the pause the agent reports briefly and asks.* A short status — what the proposal says and what implementing it would involve — and then two questions: **should implementation begin**, and **should the branch be pushed**. Not a summary of the proposal, which the human can read; a statement of where the work stands and what happens next.

*Pushing is asked, never assumed.* Pushing makes the work recoverable and hands it to whoever comes next; it also publishes unfinished work and ends the branch's freedom to be rebased onto a newer `main` under [PDR-007](PDR-007-review-boundary.md). Both are real, the trade differs by change, and the human owns it. An unpushed pause is preserved local work, exactly as [C-012](../constraints/C-012-agents-never-merge-own-changes.md) already describes a human-requested stopping point.

*The pause is recorded, not remembered.* It goes in `tasks.md` under [C-003](../constraints/C-003-tasks-tick-immediately.md), so a later reader — or a later agent — can see that the work stopped deliberately rather than being abandoned.

*Without an opt-in the loop proceeds directly to implementation.* Making the pause mandatory would charge every ordinary change for a handoff that almost never happens.

## Rejected: push on every commit

Considered and withdrawn in the same conversation that produced this record. It guarantees the work survives a lost machine and makes handoff always possible, and it costs the unpublished-branch rebase that [PDR-007](PDR-007-review-boundary.md) permits: a branch published from its first commit can only take a newer `main` by merging it in. Asking at the pause buys the same handoff where handoff is actually wanted, and leaves the ordinary change's history alone.

## Rejected: treat the pause as an agent's judgement call

An agent that decides for itself when a proposal is worth stopping on will stop either always or never, and neither is right. The signal that a proposal needs a second reader is one the human has and the agent does not.

## Carrier

No carrier outside this record states the reporting and asking rules today; the `clue-delta` loop's Propose step carries the pause itself in its earlier, narrower form. Updating it — and its scaffolded copies — is [P-013](../plans/P-013-simplification.md)'s M-063, which owns every edit to a skill in that campaign.
