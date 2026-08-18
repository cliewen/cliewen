---
id: ADR-021
type: decision
status: verified
links: [G-002, CAP-004, P-002, ADR-011, ADR-018]
title: Skills are generated as standalone artifacts from shared canonical sources
author: agent
accepted-by: Flemming N. Larsen (2026-07-18, conversation)
---

# ADR-021 — Generated standalone skills

> **ADR-059 revises completeness from one generated file to one generated skill directory:** `skill.md` is now a short router to required local references, while copying the complete directory still preserves every instruction and creates no dependency outside it. Shared canonical authoring, deterministic generation, independent lifecycle entry points, matching output trees, and version stamps remain unchanged.

## Context and problem statement

Skills need independent entry points while sharing cross-cutting rules, and the embedded init tree must stay identical to the installed skill set.

## Decision outcome

**Author skills under `internal/skills/source/` and deterministically generate standalone skill directories into `.agents/skills/` and `internal/scaffold/templates/skills/`.** Shared instructions have one source; workflow-specific instructions remain separate. Generated directories have no runtime includes, inheritance, symlinks, or dependency on another skill. The six public entry points are `clue-analysis`, `clue-delta`, `clue-extract`, `clue-plan`, `clue-upgrade`, and `clue-verify`.

The generator and drift tests own both output trees, including per-skill version markers. The generated embedded tree replaces manual skill synchronization; ADR-018's embedding decision remains unchanged.
