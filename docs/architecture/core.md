---
id: ARCH-003
type: architecture
status: active
links: [G-001, PDR-013, C-013]
title: The Cliewen core — three load-bearing elements, a red line, and an extensible periphery
---

# The Cliewen core

Cliewen is organized like a kernel: a small core whose meaning is protected, surrounded by periphery that exists to serve it and may change cheaply. This file is the durable statement of that core; the red line protecting it is the register entry [C-013](../constraints/C-013-core-changes-need-decision.md), decided in [PDR-013](../decisions/PDR-013-explicit-core-red-line.md).

## The three core elements

1. **The verifiable thread.** Goal → plan → change → capability → acceptance criterion → test: every durable claim about the system traces to executable evidence. The thread is what makes the corpus a system-of-record rather than documentation.
2. **The human merge boundary.** An agent never merges its own change; the human merge is the act of acceptance ([C-012](../constraints/C-012-agents-never-merge-own-changes.md)). This boundary is what lets agents do the work without owning the truth.
3. **The deterministic judge.** `clue validate` is the machine check of corpus form — same binary locally and in CI, enforced as a wall by branch protection. The judge is what makes "the corpus is well-formed" a fact rather than an opinion.

Remove any one element and the other two stop meaning anything: evidence without a thread is trivia, a thread without the merge boundary is unaccepted, and both without the judge are unenforced.

## The red line

A change that alters the *meaning* of a core element — what the thread connects, what a merge accepts, what a green validate asserts — requires an explicit decision record and human acceptance. It is never a plain or light change, and it never rides silently inside another change. When a peripheral rule conflicts with a core element, the peripheral rule yields or is retired.

## Periphery

Everything else is periphery: analysis records, the public guide, foreign-soil trials, quality bars, change-tier prose, index generation, scaffold templates, skills wording. Periphery changes go through the ordinary change loop at ordinary cost, and "does the core need it?" is the standing test for whether a peripheral concept earns its place.

## Extension

Adopters extend Cliewen by putting their own artifacts — including their own artifact types — into their corpus under `/docs`. The core does not enumerate what a corpus may contain; it only fixes what the thread, the merge boundary, and the judge mean. Adopter-defined types are validated against the same form rules as everything else (core frontmatter, unique IDs, resolvable links, the default status lifecycle) without needing Cliewen's permission to exist.
