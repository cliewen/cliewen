---
id: C-008
type: constraint
status: active
links: []
title: Completed plans are immutable
source: docs/plans/README.md, ARCH-001 lifetime classes
enforcement: partial
---

# C-008 — Completed plans are immutable

A plan at `status: completed` is frozen and never deleted: the plans index doubles as the project's achievement overview, and rewriting a finished campaign rewrites history.

**Checked by:** the completed-plan guard in this repository's CI (`.github/scripts/completed-plans.sh`): a pull request that modifies a plan whose status on the merge base is `completed` fails. It compares against the base, which is a workflow's business and not the judge's — `clue validate` reads a state and never a transition ([ADR-044](../decisions/ADR-044-judge-reads-state-not-transitions.md)). The change that *closes* a plan is unaffected, because the file it edits is still `active` on the base.

**Residual:** everything outside an integration into `main`. A branch may rewrite a finished campaign freely until it is proposed, and an adopter's repository carries no such step — the guard is this repository's own, not part of the shipped wall. Both are bounded by the same thing that makes the guard sufficient here: nothing reaches `main` except through a pull request ([C-012](C-012-agents-never-merge-own-changes.md)).
