---
id: ARCH-003
type: architecture
status: active
links: [G-001, PDR-013, PDR-019, PDR-037, ADR-044, C-013]
title: The Cliewen core — three load-bearing elements, a red line, and an extensible periphery
---

# The Cliewen core

Cliewen is organized like a kernel: a small core whose meaning is protected, surrounded by periphery that exists to serve it and may change cheaply. This file is the durable statement of that core; the red line protecting it is the register entry [C-013](../constraints/C-013-core-changes-need-decision.md), decided in [PDR-013](../decisions/PDR-013-explicit-core-red-line.md) and refined in [PDR-019](../decisions/PDR-019-methodology-contract-carriers-move-together.md).

## The three core elements

1. **The verifiable thread.** Goal → plan → change → capability → acceptance criterion → acceptance evidence: every durable claim about the system traces to its declared proof. Machine-proven criteria use supported, classified executable evidence; a genuine Human-class criterion uses the pull request acceptance brief as its evidence carrier. `clue validate` judges declarations and references but does not execute tests or replace the human judgment. The thread is what makes the corpus a system-of-record rather than documentation.
2. **The human merge boundary.** An agent never merges its own change; for a full Cliewen change, the human-controlled merge commit is the act of acceptance under [PDR-021](../decisions/PDR-021-supported-merge-commit-history.md) and [C-012](../constraints/C-012-agents-never-merge-own-changes.md). This boundary is what lets agents do the work without owning the truth.
3. **The deterministic judge.** `clue validate` is the machine check of corpus form — same binary locally and in CI, enforced as a wall by branch protection. The judge is what makes "the corpus is well-formed" a fact rather than an opinion.

Remove any one element and the other two stop meaning anything: evidence without a thread is trivia, a thread without the merge boundary is unaccepted, and both without the judge are unenforced.

## The red line

A change that alters the *meaning* of a core element — what the thread connects, what a merge accepts, what a green validate asserts — requires an explicit decision record and human acceptance. It is never a plain or light change, and it never rides silently inside another change. When a peripheral rule conflicts with a core element, the peripheral rule yields or is retired.

## Periphery

Everything else is periphery: analysis records, the public guide, foreign-soil trials, quality bars, change-tier prose, index generation, scaffold templates, skills wording. Periphery changes go through the ordinary change loop at ordinary cost, and "does the core need it?" is the standing test for whether a peripheral concept earns its place — with one qualification: for a command, check, required field, or artifact type, [PDR-037](../decisions/PDR-037-tooling-is-judged-by-what-it-holds.md) asks instead whether removing it would move an obligation from a machine back to a human, because the core needs guarantees and those are means.

## Where the method does not apply

The core fixes what the thread, the merge boundary, and the judge mean; it does not claim that every repository can carry them. Cliewen is the wrong choice for a repository that cannot own both the intent and its acceptance evidence, and specifically when any of these holds:

- Work does not go through Git branches and a human-controlled merge boundary — there is nowhere for the merge boundary to be.
- The project cannot run reliable tests or enforce a stable CI check before integration — the judge has nothing to be a wall in front of.
- The code is a disposable prototype, generated output, or vendored source whose behavior is accepted somewhere else — the thread would trace to a decision made in another repository.
- One corpus would need to claim test evidence spread across several repositories — validation is repository-local ([ADR-044](../decisions/ADR-044-judge-reads-state-not-transitions.md)).
- The team will not let agents maintain the corpus with the implementation, or will not review meaning before merge — both halves of the boundary are then unstaffed.

These are conditions on the adopting repository rather than rules a change can violate, which is why they live here and not in the constraint register. The honest answer in those cases is to keep the project's existing lightweight notes and tests: a corpus nobody maintains is worse than no corpus, because it makes a claim about the system that is not being kept true. Stating this is part of the design — a methodology that overclaims trains its readers to stop reading its claims.

## Extension

Adopters extend Cliewen by putting their own artifacts — including their own artifact types — into their corpus under `/docs`. The core does not enumerate what a corpus may contain; it only fixes what the thread, the merge boundary, and the judge mean. Adopter-defined types are validated against the same form rules as everything else (core frontmatter, unique IDs, resolvable links, the default status lifecycle) without needing Cliewen's permission to exist.
