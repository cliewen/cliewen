---
id: PDR-034
type: decision
status: verified
links: [CAP-007, ADR-013, G-001, C-013, P-013, AN-018]
title: The corpus is read narrowly by default, and widened only when the work requires it
author: agent
accepted-by: Flemming N. Larsen (2026-08-08, conversation)
---

# PDR-034 — The corpus is read narrowly by default

## Context and problem statement

The `/docs` corpus is the system-of-record, and it only grows. Every change adds artifacts and none are removed cheaply, so the cost of reading it before doing anything is a cost every future change pays, forever. An agent that loads the corpus to answer a question spends most of its budget before the work starts, and a large fraction of what it loads has no bearing on the task.

[CAP-007](../capabilities/CAP-007-focused-context/README.md) exists to make that avoidable: `clue context <id>` emits one artifact and its transitive link slice, deterministically, so an agent can read what an identity actually depends on instead of a folder. The capability was built and the command ships.

Nothing obliges anyone to use it that way. The routing hub describes the narrow path, but no decision, constraint, goal, or criterion states that narrow reading is the rule — the capability supplies a mechanism and the obligation was never written down. A mechanism nobody is required to use is a mechanism that gets used when someone remembers.

## Decision outcome

**Reading the corpus starts at the narrowest point that answers the question, and widens only when the work or a discovered edge requires it.**

- *An identity that is already known is entered directly.* When the request names or resolves to an artifact, `clue context <id>` is the entry point and its outgoing-link slice is what gets read. The index is not consulted to find something already named.
- *The index is for finding the entry point, not for reading.* `docs/README.md` is consulted when the request does not resolve to an artifact, to identify the closest one — and then the previous rule applies.
- *Reading past the slice is a decision, not a default.* An edge discovered during the work, a rule that turns out to bind, a conflict that needs adjudicating — these widen the read, and the widening is driven by something encountered rather than by caution.

This is a rule about where reading *starts*, never a limit on what may be read. An agent that needs an artifact reads it; what is forbidden is loading the corpus first and deciding afterwards what mattered.

**This decision states an obligation and adds no mechanism.** CAP-007's command, its determinism, and its output are unchanged, and nothing here reaches [C-013](../constraints/C-013-core-changes-need-decision.md)'s core: what the verifiable thread connects, what a merge accepts, and what a green `clue validate` asserts are all untouched. What changes is that the narrow path is the rule rather than the clever option.

## Rejected: leave it to judgement, since the command already exists

The state this decision ends. The command shipped, the hub described the path, and no artifact said it was expected — so a reader asking *why does the hub tell me to read this way?* had no answer, and an agent under no obligation reads defensively. Defensive reading is the expensive default precisely because it never looks wrong.

## Rejected: cap what may be read

A limit — a file count, a depth — would be checkable and would be wrong. Some work genuinely needs a wide read, and a cap turns that into a rule to be worked around rather than a judgement to be made. The obligation is about where reading begins, which is where the cost is decided.

## Carrier

The routing hub states the three steps of this rule, and the scaffolded hub states them for adopters ([ADR-013](ADR-013-ships-generic-vs-repo-local.md) keeps that layer repository-local). This record is what those statements trace to; before it they traced to nothing, having been recorded against an acceptance criterion that describes what `clue context` emits rather than when to run it.
