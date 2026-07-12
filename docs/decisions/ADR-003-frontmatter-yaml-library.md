---
id: ADR-003
type: decision
status: verified
links: [CAP-002]
title: Parse frontmatter with gopkg.in/yaml.v3
author: agent
accepted-by: Flemming N. Larsen (2026-07-12)
---

# ADR-003 — Parse frontmatter with gopkg.in/yaml.v3

## Context and problem statement

`clue validate` must parse YAML frontmatter. Hand-rolling a parser keeps the binary dependency-free; a library parses what agents and humans will actually write. (Agent decision, promoted to `verified` by human acceptance 2026-07-12.)

## Considered options

1. **gopkg.in/yaml.v3** — the de facto standard Go YAML library.
2. **Hand-rolled flat `key: value` subset** — zero dependencies.

## Decision outcome

**yaml.v3.** Frontmatter is written by agents and humans, not generated: a strict hand-rolled subset would reject legitimate YAML (quoted strings, flow lists, multi-line values) with confusing errors, and the judge must never be wrong about form. One small, stable, pure-Go dependency does not threaten single-binary distribution (ADR-001's driver). Revisit only if the dependency ever blocks a `clue` release.
