---
id: ADR-059
type: decision
status: inferred
links: [P-015, CAP-004, ARCH-002, ADR-021, PDR-029]
title: A generated skill is a standalone directory with a routing entry point
author: agent
accepted-by: []
---

# ADR-059 — Progressive standalone skill directories

## Context and problem statement

ADR-021 made each generated `skill.md` a complete standalone carrier, but loading every rule before a workflow branch is known imposes unnecessary reading cost. Deferring detail must reduce entry-point cost without making a copied skill depend on another directory or hidden source.

## Decision outcome

**A generated skill is a complete standalone directory: `skill.md` is a short routing entry point, and mandatory deferred instructions live in that directory's `references/`.** The entry point states purpose, the read-before-action boundary, and when each reference is required; the directory contains every named reference and no methodology rule disappears merely because it moved.

Copying one directory preserves its marker, version, workflow, and complete instructions without runtime dependency on another skill or `internal/skills/source/`. Canonical sources remain centralized, generation emits both distributed trees, drift checks cover every managed path, and tests inspect the whole directory except where entry-point placement is itself the contract.

**Carrier:** canonical skill sources, generated directories and scaffold templates, directory drift checks, and routing tests.
