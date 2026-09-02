---
id: IDR-004
type: decision
status: verified
links: [ADR-011, ADR-012, ADR-013, CAP-004]
title: Release verification reuses the shipped judge and reviewed notes
author: agent
accepted-by: Flemming N. Larsen (2026-09-02)
---

# IDR-004 — Release verification reuses the shipped judge and reviewed notes

## Context

Release-only shell checks and generated changelogs would duplicate rules already held by the shipped judge and the reviewed changelog, while a stale locally installed development binary can invalidate repository release checks.

## Decision

This repository refreshes its local `clue` from the checkout before release verification. The release workflow runs the shipped judge stamped as the tag to enforce tag-to-skill consistency and passes the reviewed changelog section through `--release-notes`; `.goreleaser.yaml` carries no competing changelog block.
