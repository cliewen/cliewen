---
id: ADR-002
type: decision
status: verified
links: [G-001]
title: "The inbox is goals with status: proposed"
author: human
accepted-by: Flemming N. Larsen (2026-07-12)
---

# ADR-002 — The inbox is goals with `status: proposed`

## Context and problem statement

Ideas and bug reports need a pre-change home whose provenance remains in the repository, whose status is explicit, and whose intake is low-friction.

## Decision outcome

**Use goal files with `status: proposed` in `/docs/goals`.** Promotion to `accepted` is a human decision made through a change and pull request, and the generated goals index is the backlog view.

The `/docs/goals` README that `clue init` scaffolds carries the inbox rule and entry status; this repository's [goals/README.md](../goals/README.md) is the template source. GitHub Issues may provide social intake, but they are not the corpus record: tickets reference goal IDs and never replace them.
