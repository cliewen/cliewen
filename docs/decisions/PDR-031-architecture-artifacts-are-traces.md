---
id: PDR-031
type: decision
status: verified
links: [P-013, PDR-029, PDR-013, C-013, G-001, AN-018]
title: An architecture artifact is a valid trace when it states the rule
author: agent
accepted-by: Flemming N. Larsen (2026-08-08, conversation)
---

# PDR-031 — An architecture artifact is a valid trace when it states the rule

## Context and problem statement

[PDR-029](PDR-029-simplification-tests-by-surface.md) judges carrier prose by three conditions, of which tracing is the operative one: a rule-bearing statement traces when a live corpus artifact *states* the rule it carries. It names four artifact types as valid targets — decision, constraint, goal, acceptance criterion — and bounds the test hard, because derivability is not a trace and a looser rule would let every statement pass by argument from [G-001](../goals/G-001-verifiable-thread.md).

Architecture is not among the four, and the omission bites immediately at the most load-bearing place in the corpus. The routing hub's rule about the core — what the verifiable thread, the human merge boundary, and the deterministic judge are — cites [ARCH-003](../architecture/core.md) by ID, because ARCH-003 is the durable statement of what the core *is*: the thing [C-013](../constraints/C-013-core-changes-need-decision.md) protects and [PDR-013](PDR-013-explicit-core-red-line.md) drew a red line around. Read literally, PDR-029 reports the hub's core definition as untraceable, and the only repairs available are to delete it from the hub or to write a decision record restating what ARCH-003 already says in numbered form.

Both repairs are worse than the defect. The second is the one that matters: a corpus that forbids tracing to its own architecture forces every architectural fact to be duplicated into a decision before a carrier may state it, which manufactures exactly the second stored representation this project refuses everywhere else.

## Decision outcome

**An architecture artifact is a valid trace under PDR-029's carrier test, under the same restriction the other four types already carry: it traces when the architecture file states the rule, and not when the rule is merely derivable from the file's broader description.**

The restriction is what keeps this from being a loophole, and it is not a new one. Architecture prose is more narrative than a constraint or a criterion, so the temptation to trace a statement to "the architecture explains this area" is stronger — and it is the same failure PDR-029 already blocks when someone traces a rule to G-001's intent. ARCH-003 passes because it states the core's three elements in a numbered list and states the red line in one sentence; a paragraph describing how a subsystem is organised does not become a trace for every rule about that subsystem.

The accepted set is therefore five types: decision, constraint, goal, acceptance criterion, and architecture. Everything else PDR-029 says about tracing is unchanged, including that the register names the *narrowest* artifact stating the rule — where both a decision and an architecture file state it, the decision is narrower and is what the register names.

This decision amends PDR-029's trace-type clause and touches nothing else in it. The two tests, the surface split, the ordering rule, the overlap rules, and the refusal to score the campaign in words all stand as written.

## Rejected: leave the four types and restate ARCH-003 as a decision

The reading that requires no amendment, and it produces a decision record whose whole content is a copy of an architecture file — with the two then obliged to move together forever under [PDR-019](PDR-019-methodology-contract-carriers-move-together.md). The corpus would gain a record, lose nothing, and acquire a drift risk at the one place it can least afford one.

## Rejected: accept any corpus artifact as a trace

Symmetrical, simple, and it dissolves the test. Analysis records and plans describe rules constantly without deciding them; admitting them would let a carrier statement trace to a findings document that merely observed the practice, which is the failure that produced this campaign's first escalation rather than a solution to it.

## Carrier

This record and [PDR-029](PDR-029-simplification-tests-by-surface.md)'s trace-type clause, which is amended here rather than restated. [AN-018](../analysis/AN-018-skill-statement-register.md)'s register applies the widened set. P-013's milestones are unaffected: M-062 recommended and did not build, and M-063 trims against the register as amended.
