---
id: ADR-019
type: decision
status: verified
links: [ADR-013, CAP-001, P-002]
title: Index regeneration runs in clue init
author: agent
accepted-by: Flemming N. Larsen (2026-07-18, PR #20 review conversation)
---

# ADR-019 — Index regeneration is part of init

## Context and problem statement

Existing repositories need their README indexes regenerated during `clue init` to become green immediately; postponing that work to `clue scaffold` would leave onboarding red.

## Decision outcome

**`clue init` regenerates taxonomy README index blocks on every run.** This supersedes only ADR-013's emits-empty clause. The scaffold engine keeps surviving hand-written single-line entries, appends missing entries, and leaves prose outside markers untouched. `clue scaffold` remains the standalone exposure of the same engine, and fresh templates are idempotent. `checkIndexes` judges both paths.
