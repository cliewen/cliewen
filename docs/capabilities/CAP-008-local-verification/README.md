---
id: CAP-008
type: capability
status: active
links: [G-004]
title: Local verification — documented commands run on supported contributor environments
goal: G-004
---

# CAP-008 — Local verification

## What

Contributors can run the repository-local verification block in CONTRIBUTING.md verbatim on supported environments, including both the command that writes its coverage profile and the one that reports it.

## Why

The block is the repository’s verification contract. A command whose syntax fails on a supported contributor environment makes a valid change appear unverified and invites unsound local substitutions.

Acceptance criterion: [criteria.md](criteria.md) · implementation boundary: [design.md](design.md).
