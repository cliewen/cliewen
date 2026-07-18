---
id: CH-030
type: change
status: open
links: [P-003, M-009, G-003]
title: Cliewen has a community front door
---

# CH-030 — Cliewen has a community front door

## What

Add the repository-local files that make Cliewen ready for public participation: contribution guidance, Contributor Covenant 3.0, a coordinated vulnerability disclosure policy, structured bug and proposed-goal issue forms, and a pull-request template that carries the change loop into GitHub.

The contribution path preserves Cliewen's review boundary. Issues express demand or defects but do not become accepted plans by themselves; pull requests declare a real plan item or plan-less status, follow the full or light tier, provide verification evidence, and remain for a human maintainer to merge.

## Why

Serves P-003/M-009. A public repository without an explicit intake path invites contributors to bypass the corpus, disclose vulnerabilities in public, or treat issue activity as accepted direction. The front door must make the safe and reviewable path the obvious path before visibility flips.

## Decision boundary

PDR-010 records the durable repository governance: Contributor Covenant 3.0, separate confidential conduct and security aliases, coordinated disclosure expectations, bug-versus-goal issue intake, and the unchanged human merge gate. These files govern Cliewen only and are not emitted by `clue init`.

GitHub private vulnerability reporting is activated with M-012 after the repository becomes public; email remains the stable reporting channel. This change adds no support forum, discussions strategy, CLA, DCO, or community files for adopting repositories.
