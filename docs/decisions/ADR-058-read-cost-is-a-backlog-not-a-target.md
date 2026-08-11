---
id: ADR-058
type: decision
status: verified
links: [P-015, ADR-057, ADR-056, CAP-007, C-022, C-004]
title: A read-cost report is a backlog judged artifact by artifact, and a link is never deleted to move its count
author: agent
accepted-by: Flemming N. Larsen (2026-08-11, conversation)
---

# ADR-058 — A read-cost report is a backlog, not a target

## Context and problem statement

[ADR-057](ADR-057-read-cost-measurements.md) makes two structural populations visible and states what they are for: the report "is a report, not a verdict", later work decides whether an artifact genuinely needs what it holds, and a failing check was rejected so that "reporting makes the backlog inspectable first". The same record states that the measure does not claim artifact count measures quality.

A reported population still arrives as a number, and a number invites an exit criterion phrased as reaching it. Once that criterion exists, the cheapest way to satisfy it is to delete `links` entries until the count falls, and nothing in the report distinguishes a citation that repeats what the prose already says from the one edge that carries a decision's evidence, a plan anchor on the goal-to-criterion thread, or the constraint a rule lives in.

That deletion is not deferral. [ADR-056](ADR-056-bounded-context-slice.md) bounds what a slice prints and names what the bound held back, so an artifact beyond the bound is still reachable by widening; `corpus.Context` walks `links` outward only, so an entry removed from the list is reachable at no depth from that artifact and is recoverable only by a reader who already knows what to look for. The bound moves reading cost. Deleting the entry moves meaning out of the corpus.

The corpus already holds the general form of this rule for machine checks. [C-004](../constraints/C-004-never-weaken-checks.md) states that a failing check is fixed at its cause or surfaced as a conflict, never loosened to go green, because the check is worth what it still catches. A read-cost report is not a failing check, but a corpus edited until the report is quiet reproduces exactly the failure C-004 names: the signal turns green while the thing it was measuring gets worse, and the edit is indistinguishable in a diff from legitimate pruning.

## Decision outcome

**The reported over-budget population is a backlog worked artifact by artifact, and a `links` entry is never deleted to move the reported count.**

- Repairing a reported identity means inspecting it. Its outcome is either the removal of entries that are genuinely redundant for a reader of that artifact, or an accepted exception with its reason stated where the artifact's plan or the artifact itself records it.
- An accepted exception keeps its row in the report. The count is therefore not expected to reach zero, and a plan may not state reaching zero as an exit criterion. What a plan may require is that every reported identity has been inspected and its outcome recorded.
- A `links` entry records a relationship that holds. Whether it stays is a question about that relationship, never about the population's size.
- No mechanical rule selects which entries survive. Position in the list, presence of the target's ID in prose, and the existence of a reverse edge are all proxies for relevance, and a proxy applied across a population produces exactly the deletions this record refuses.
- Where the cost of an entry path is genuinely too high after inspection, the repair belongs to what the reader loads rather than to what the corpus records: [CAP-007](../capabilities/CAP-007-focused-context/README.md) owns how much of a slice prints and how the remainder is named, and bounding it further is a change to that capability.

**Carrier:** P-015's M-072, which states the exit criterion this record governs; ADR-057, which gains a note that its report is worked as a backlog; and C-022, whose residual already holds the judgment half of the same measure.

## Rejected: drive the reported count to zero by trimming links

This is the reading that makes the report a check the corpus is edited to satisfy. It cannot distinguish a redundant citation from a load-bearing one, because the count carries no information about either, so the edits it produces are selected by whatever is cheapest to identify rather than by what the reader needs. The result is a quiet report over a corpus whose thread has holes, which is the outcome C-004 exists to prevent and the opposite of what ADR-057 set out to make visible.

## Rejected: make the budget a failing check so the corpus cannot drift

ADR-057 already rejected this, and this record does not reopen it. A hard failure would make repair order a condition of unrelated work, and the judgment each reported identity needs is exactly what a failing check cannot supply.

## Rejected: rank links automatically and keep the highest-scoring budget's worth

A ranking rule is attractive because it is auditable and cheap to apply to a whole population. Every available signal is a proxy: list order encodes when an edge was added, a mention in prose means a reader can recover the citation but says nothing about whether the relationship matters, and a reverse edge depends on the other artifact's unrelated editing history. A proxy that is wrong for one artifact in ten is applied silently to all of them, and the artifacts where it is wrong are the ones whose links were load-bearing enough to be worth keeping.
