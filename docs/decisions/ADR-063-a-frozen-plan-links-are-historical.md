---
id: ADR-063
type: decision
status: verified
links: [P-020, M-083, ADR-034, C-008, ADR-044, CAP-002]
title: A completed plan's links are a historical record, not live navigation
binds: adopter
author: agent
accepted-by: Flemming N. Larsen (2026-09-02)
---

# ADR-063 — A completed plan's links are historical

## Context

Three accepted rules cannot all hold when a completed plan links an artifact being retired. [ADR-034](ADR-034-retirement-is-deletion.md) makes retirement a deletion. `checkLinks` rejects the surviving link and demands it be repointed to the successor. [C-008](../constraints/C-008-completed-plans-immutable.md) freezes the plan, and its CI guard fails on any modification to it whatsoever. Repairing the link fails the guard; leaving it fails the judge.

There is no third move, and this is not an edge case: a campaign's plan naturally links the spikes that fed it, so every one of those spikes becomes permanently unretirable the moment the campaign closes. The forgetting mechanism the corpus already accepts is unusable in its most common case, here and in any adopter.

## Decision

**A completed plan's `links` may name a retired artifact, and `clue validate` accepts it.** A finished campaign is a record of what it referenced while it ran, in the same way a pinned analysis records what it observed at a revision. Reading its outgoing links as live navigation is the error the rule made; a frozen document cannot be asked to keep a pointer current.

The allowance is narrow in three ways. It applies only to a plan whose status is `completed`, so an active plan still repoints. It applies only to an artifact that was *properly* retired — one some live artifact names in `supersedes:` — so a link to an ID that never existed remains a failure everywhere, because that is a typo rather than history. And it changes nothing about C-008: the freeze stands exactly as it was, which is what makes this the cheaper of the two available answers.

The alternative was to carve a link-repair exemption into C-008 and its guard. That weakens a core freeze *and* still requires editing a finished campaign's file — the precise act C-008 exists to prevent — to achieve the same end.

**Accepted cost:** a completed plan's *prose* may still contain a markdown link to a deleted file, which no check covers and which C-008 forbids repairing. Such a link will not resolve for a reader. The plan's sentence remains a true statement about what the campaign referenced, and Git history holds the target, but the convenience of following it is gone. Repairing it would require editing frozen history, which is a worse trade.

**Accepted cost:** `clue context` reports each such link as an edge it could not follow. That output is expected for a frozen plan rather than a defect, and it does not affect the command's exit code.

The corpus structure this affects is described in [`docs/architecture/README.md`](../architecture/README.md).

## Carrier

The completed-plan allowance in `internal/corpus`, its criterion, this record, the `checkLinks` description in `docs/capabilities/CAP-002-validate/design.md`, and the scaffolded corpus guidance in `internal/scaffold/templates/docs/README.md` describing what a frozen plan's links mean.
