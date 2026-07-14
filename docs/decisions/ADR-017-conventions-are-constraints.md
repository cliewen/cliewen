---
id: ADR-017
type: decision
status: verified
links: [ADR-013]
title: Prose conventions register as constraint artifacts with enforcement classes
author: agent
accepted-by: Flemming N. Larsen (2026-07-14, PR #10)
---

# ADR-017 — Prose conventions register as constraint artifacts

## Context and problem statement

Methodology rules that live only in prose — AGENTS.md rules, README conventions, skill instructions — have no inventory, no visible backlog, and drift silently: nothing counts them, nothing tracks which ones a machine could enforce but doesn't. The constraints folder has declared `enforcement: machine|agent|human` since the baseline while only `machine` was implemented. Where is the single register of the rules that bind every change, and how does an unenforced rule stay visible until it is enforced?

## Decision outcome

**Every prose-only convention becomes a constraint artifact in `docs/constraints/`, and the `agent` enforcement class is the promotion backlog.**

- Each constraint carries `source:` (the doc that states the rule — the prose stays as the human-readable carrier and gains a pointer back where that helps) and `enforcement: machine|agent|human`. `clue validate` requires both fields and the vocabulary.
- An `agent`-enforced constraint states its **promotion trigger** in its body: the condition under which the rule becomes a `machine` check in `clue`. Promotion is an ordinary change that implements the check and flips the field.
- The backlog is visible, not archival: `clue validate` reports the count of `agent`-enforced constraints on its OK line, the same way it reports born-`inferred` artifacts awaiting verification.
- The constraints README index is the register table; there is no second inventory.
- This deliberately opens the "`enforcement:` classes beyond `machine`" door the architecture had listed as out. Shipped skills stay generic (no repo doc-IDs), so constraint pointers live in repo-local prose only — AGENTS.md and folder READMEs.

**Carrier:** the constraint-field lint and the agent-count report in `clue validate` (machine); the register convention in the constraints folder README, which `clue init` will scaffold (default); the pre-PR checklist item that assesses every change against the constraints folder (agent).

### Rejected: a standalone conventions table in a README or AGENTS.md

A table row is not an artifact: no frontmatter, no ID to link, no lintable fields — the register would itself be prose, reproducing the drift problem it exists to fix. The constraints folder already models exactly this with fields nothing yet consumed.

### Rejected: promoting everything to machine checks now

Several rules need git-diff context `clue` does not have, and some halves are meaning (timeless prose, weakening-by-refactor) that machines cannot judge. Registering first makes the gap countable; promotion proceeds trigger by trigger.
