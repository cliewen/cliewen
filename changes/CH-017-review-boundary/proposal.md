---
id: CH-017
type: change
status: open
links: []
title: Review-boundary rules — branch from main, one change in flight per author, humans merge
---

# CH-017 — Review-boundary rules

## What

Tighten `clue-delta` and `clue-verify` so the PR review boundary cannot be reinterpreted away: every change branches from the current tip of `main`, each author has one change in flight at a time, agents never merge their own PRs or push to `main`, review fixes stay on the reviewed branch, and stacking on unmerged work is an explicit human decision. Record the decision as a new PDR, register it as a new constraint, and surface the rule in AGENTS.md. The records this change creates will be linked from this proposal's frontmatter once they exist.

## Why

The model2diagram 0.2.0 audit (its `AN-002-cliewen-workflow-audit.md`) found CH-006 through CH-010 there were stacked on each other's unmerged tips, merged with agent-created local merge commits, and pushed directly to `main` — zero PRs, zero review surface. The skills were satisfiable without a human ever seeing anything: `clue-delta` says "PR + merge" but never says who merges, never forbids starting a new change on unmerged work, and never says review fixes belong to the current branch. This change closes those textual gaps; enforcement in adopter CI is a separate change in the adopter repo.

## Plan item

Plan-less, explicitly: this is a methodology correction prompted by an adopter audit. P-002/M-007 (foreign soil) does not apply — model2diagram shares a maintainer with this repo, so its findings are adopter feedback, not a foreign-soil trial.
