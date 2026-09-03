---
id: PDR-054
type: decision
status: inferred
links: [G-012, P-021, M-087, UC-001, CAP-009, ARCH-003, PDR-037]
title: Use cases stay optional, and Cliewen introduces no generic requirement artifact
binds: adopter
author: agent
accepted-by: []
---

# PDR-054 — Use cases stay optional, and there is no requirement artifact

## Context

Adding an actor-oriented artifact invites two adjacent additions that would each look like tidying up. The first is making it mandatory, so every capability has one. The second is generalizing it into a `requirement` type that a goal, a use case, and an acceptance criterion are all special cases of.

Both are familiar from heavier methods, and both are answers to a question Cliewen does not have. The corpus already carries what someone wants (goal), what the system can do (capability), and what proves it (criterion). What it lacks is only the actor's path *across* capabilities, and most behaviour has no such path worth writing down.

## Decision

**A use case is created when it changes what a reader understands, and omitted otherwise.** It is worth writing when an actor's journey crosses several capabilities, when ordering or failure recovery carries meaning, when the criteria are each locally correct and do not add up to the outcome, when several actors collaborate, or when a brownfield system holds behaviour no single capability explains. It is not worth writing when one capability and its criteria already describe the behaviour, when there is no actor interaction, when it would restate a goal, or when the subject is internal implementation behaviour that design, architecture, a constraint, or an IDR owns.

The agent recommends, with its reason; the human decides. Absence is not a gap, and no check, report, or metric treats it as one.

**No generic requirement artifact.** A type that subsumes goal, use case, and criterion would make the thread's steps indistinguishable and cost the corpus the thing that makes it verifiable: a criterion is provable and a goal is not, and a single type would have to either drop that distinction or carry a discriminator that recreates it under another name. [PDR-037](PDR-037-tooling-is-judged-by-what-it-holds.md)'s test settles it — removing the distinct types would move an obligation from the machine back to a human, which is the direction the core does not go.

**Plan, milestone, and change stay out of the intent hierarchy.** They describe how intent is delivered, not what the product means. Vision → goal → use case → capability → criterion → evidence is the intent thread; goal → plan → milestone → change → accepted commit is the delivery thread. They meet at the goal and nowhere else, which is why a change never edits a vision merely because it edited code.

## Rejected: a coverage report over use cases

Cliewen already reports proof coverage by capability, and the shape is tempting to reuse. It would be actively harmful here. A percentage over an optional artifact reads as a target, and the only way to move it is to write use cases nobody needs — which is precisely the outcome this record exists to prevent. `clue validate --intent` lists what exists and computes no ratio.

## Carrier

The optionality of every use-case rule in `internal/corpus`, the recommendation test in `internal/skills/source/shared/intent-model.md.tmpl`, the counterexample in `guide/intent.md`, and the intent report's deliberate absence of a figure.
