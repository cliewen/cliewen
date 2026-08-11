---
id: C-004
type: constraint
status: active
links: [ADR-039]
title: Never weaken a test or a lint rule to make a build pass
source: the shared review-boundary fragment ("Never weaken the workflow or required-check policy to make a change pass."), clue-verify introduction ("Never fix a failure by weakening the check.")
enforcement: human
---

# C-004 — Never weaken a test or a lint rule to make a build pass

Machines enforce form so agents cannot cheat; weakening the check inverts that. A failing check is fixed at its cause or surfaced as a conflict — never deleted, skipped, or loosened to go green.

**Residual:** all of it. A deleted assertion, a loosened bound, and a narrowed lint rule are byte-for-byte what a legitimate refactor also produces; what separates them is whether the check still catches what it was written to catch, which is a question about meaning. [ADR-044](../decisions/ADR-044-judge-reads-state-not-transitions.md) settles that no diff-reading check in `clue` will be built for it, and no diff would answer it anyway.

The cost is the largest of any rule in this register, and it is stated plainly: this is the constraint that keeps every other machine honest, and nothing but review holds it. Weakening a check to go green produces a green build, and the register's other fifteen rules are worth exactly as much as this one is obeyed. What makes it survivable is that the act leaves evidence — a weakened check is visible in the diff a human merges, and [C-012](C-012-agents-never-merge-own-changes.md) guarantees a human is there.
