---
id: ADR-061
type: decision
status: verified
links: [ADR-040, PDR-039, G-001]
title: Cliewen has no general external coordination interface
author: agent
accepted-by: Flemming N. Larsen (2026-09-02)
---

# ADR-061 — Cliewen has no general external coordination interface

## Context

Cross-repository evidence and tracker metadata cannot be made deterministic from one repository's bytes, and treating foreign mutable state as acceptance data would weaken the local thread and judge.

## Decision

Cliewen provides no general interface for cross-repository evidence or tracker metadata. External references remain qualified and outside the deterministic verdict; cross-repository proof stays Human or draft until a future-shaping architecture decision supplies a boundary that preserves the local judge. Authorized dependency on an unmerged branch remains the separate, repository-local workflow PDR-039 defines.
