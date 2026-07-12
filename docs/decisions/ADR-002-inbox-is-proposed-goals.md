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

Where does an idea or bug report live *before* it is a change? The Foundation Document flagged this as an open gap to decide early in the baseline.

## Decision drivers

- Provenance must be **repo-native, never forge-native** (§3).
- **Status lives in frontmatter, never in paths**; status views are generated (§4).
- Intake friction must stay near zero or the inbox gets bypassed.

## Considered options

1. **Goals with `status: proposed`** in `/docs/goals`
2. **GitHub Issues as a non-binding forecourt**

## Decision outcome

**Goals with `status: proposed`.** An idea enters the corpus as a goal file with `status: proposed`; promotion to `accepted` is a human decision made through a change/PR, which puts intake on the same provenance chain as everything else. The generated goals index doubles as the backlog view.

### Carrier (added 2026-07-12, per the carrier rule)

The `/docs/goals` README that `clue init` scaffolds: its prose declares the folder the inbox and the `proposed` status the entry state. This repo's own [goals/README.md](../goals/README.md) is the template source.

### Rejected: GitHub Issues forecourt

Zero-friction intake, but it places the pre-change record outside the corpus and leans on exactly the forge dependency §3 forbids. Issues may still exist socially, but the corpus-native record is the proposed goal, and tickets reference IDs, never the reverse.
