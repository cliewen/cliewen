---
id: C-004
type: constraint
status: active
links: []
title: Never weaken a test or a lint rule to make a build pass
source: AGENTS.md rule 5, clue-verify preamble
enforcement: agent
---

# C-004 — Never weaken a test or a lint rule to make a build pass

Machines enforce form so agents cannot cheat; weakening the check inverts that. A failing check is fixed at its cause or surfaced as a conflict — never deleted, skipped, or loosened to go green.

**Promotion trigger:** `clue` gains git-diff context and can flag deletions or loosenings of test functions and lint rules within a change — then `enforcement: machine` for the detectable subset. Meaning-preserving refactors will always need human review; this rule never fully leaves the agent/human classes.
