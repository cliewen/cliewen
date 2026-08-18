---
id: ADR-020
type: decision
status: verified
links: [ADR-013, ADR-017, CAP-001]
title: The scaffolded register seeds only conventions without a versioned carrier
author: agent
accepted-by: Flemming N. Larsen (2026-07-18, PR #20 review conversation)
---

# ADR-020 — Scope of the scaffolded constraint register

## Context and problem statement

`clue init` must seed a useful constraint register without duplicating rules already carried and versioned by generated skills.

## Decision outcome

**Seed exactly the methodology conventions no versioned skill carries.** Seeded constraints are self-sourced as scaffolded by `clue init`; the generated AGENTS.md mirrors them as readable prose. The hard-wrap prohibition is the first seed. From there, seeded constraints follow ADR-017's lifecycle: an agent-enforced entry names its promotion trigger and remains in the register when promoted.

Seeding nothing would leave prose-only rules untracked; seeding every rule would create an unversioned duplicate of each skill. The boundary governs init output, not what an adopter later registers or promotes.
