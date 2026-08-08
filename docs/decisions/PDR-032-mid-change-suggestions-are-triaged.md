---
id: PDR-032
type: decision
status: verified
links: [ADR-002, PDR-017, PDR-007, C-003, C-005, G-001, AN-018]
title: A suggestion raised mid-change is triaged into the current work or into the inbox, never into memory
author: agent
accepted-by: Flemming N. Larsen (2026-08-08, conversation)
---

# PDR-032 — A suggestion raised mid-change is triaged into the current work or into the inbox

## Context and problem statement

Humans do not raise ideas at artifact boundaries. A suggestion arrives in the middle of a change — while a proposal is being read, while a review finding is being repaired, while a digest is being written — and it is usually adjacent to the work rather than part of it.

The corpus has a rule for where work lives ([PDR-017](PDR-017-merge-gate-has-content.md): an agent's private memory is never an authoritative carrier) and a rule for where an idea enters ([ADR-002](ADR-002-inbox-is-proposed-goals.md): the inbox is a goal with `status: proposed`). Neither reaches this moment. The change workspace is transient and is deleted in the digest, so a suggestion noted in `tasks.md` that does not belong to the current work is destroyed at merge; a suggestion held only in the conversation does not survive the session. The result is that the most frequent way a human contributes direction to this project is also the only one with no carrier.

Silently widening the change is not the alternative. A suggestion absorbed into work already under review changes what the reviewer approved after they approved it, which is the failure the review boundary exists to prevent.

## Decision outcome

**A suggestion raised during a change is triaged immediately into one of two carriers, and the triage is stated rather than assumed.**

*It belongs to the current work → it becomes a task, and it is handled before merge.* The test is whether the change is wrong or incomplete without it. A defect in what is being built, a missing criterion, a carrier the change should have moved — these are the change's own scope arriving late, and deferring them ships something known to be wrong. The task enters `tasks.md` under [C-003](../constraints/C-003-tasks-tick-immediately.md) like any other, and if the change is already published the review boundary's updater handoff governs the repair.

*It does not belong to the current work → it becomes a goal with `status: proposed`, written in the digest.* This is the ordinary case and it reuses ADR-002's inbox unchanged: the idea enters the corpus as a proposed goal, the generated goals index is the backlog view, and promotion to `accepted` is a later human decision through its own change. Writing it in the digest is what makes it survive — the workspace is deleted there, so a note that lives only in `tasks.md` or `open-questions.md` dies with it.

*Neither carrier is optional, and "I will remember" is not a third.* A suggestion that is neither actioned nor recorded has been declined without anyone deciding to decline it. If the agent judges a suggestion out of scope and not worth an inbox entry, that judgement is the human's to make, not the agent's: the entry is written and the human closes it.

*The triage is stated to the human when it is made.* Which carrier a suggestion went to, and why, is a sentence — not a silent filing. A human who offers an idea and hears nothing cannot tell whether it was captured, actioned, or lost.

## Rejected: hold suggestions in the change workspace and sort them at digest

The obvious economy: note everything in one place, decide later. The workspace is deleted by the digest, so "sort them at digest" means the sorting and the destruction happen in the same commit, and anything overlooked is gone with no recovery path except reading a deleted file out of Git history. Triage at the moment of arrival costs one sentence and removes the deadline entirely.

## Rejected: a new backlog artifact type

A dedicated record would carry more structure than a proposed goal — priority, origin, the change it was raised during. ADR-002 already decided this question for the general case and put the inbox in `/docs/goals`, and a second intake surface would mean an idea's home depends on when someone happened to think of it. Where the goal type genuinely cannot hold something, that is a defect in the goal type worth its own record.

## Rejected: absorb adjacent suggestions into the change under way

Fastest, and it is how scope creep enters a methodology that otherwise forbids it. A change is authorised by its proposal; work added after review is work no one approved, and a reviewer who returns to a changed diff has to re-derive what moved. Keeping the change to what it declared is what makes the declaration mean anything.

## Carrier

The `clue-delta` loop's Propose and Digest steps and the shared `durable-work` fragment state the triage; the goals folder README already states the inbox rule this decision reuses. `clue` holds none of it: a suggestion arriving in conversation is invisible to a judge that reads files, so this is agent-enforced and its residual is stated where the rule is registered.
