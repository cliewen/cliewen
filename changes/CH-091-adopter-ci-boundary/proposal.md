---
id: CH-091
type: change
status: open
links: [P-009, M-040, AN-012, AN-014, CAP-001, CAP-006, CAP-004]
title: Replace the forked CI wall with an upstream validation workflow
---

# CH-091 — Replace the forked CI wall with an upstream validation workflow

## What

This full change serves P-009 milestone M-040 by replacing the copied all-in-one CI wall emitted by `clue init` with a thin caller and an upstream-owned reusable validation workflow. A synthetic fixture will establish the demonstrated caller choices for runner labels and approved `clue` installation, while upstream action references and validation logic remain fixed and reviewable.

The reusable workflow will preserve the stable `validate` check, changed-surface detection, the unarmed warning, the acceptance-brief gate, and `clue validate --forbid-changes`. The generated caller will carry only the versioned upstream workflow reference and the demonstrated local choices; it will not introduce a general Cliewen configuration file.

The old vendored-binary wall and its upgrade instructions will receive a documented transition. The release contract, CAP-001/CAP-006 criteria and designs, scaffold output, public guide, tests, and `[Unreleased]` note will describe one consistent adoption path.

## Why

AN-012 measured that an adopter's edit to the copied wall permanently forked upstream validation, while AN-014 reports enterprise runner, action-pinning, and installation constraints that have not yet been reproduced locally. A reusable workflow can preserve upstream fixes without pretending those unverified observations justify a general configuration interface. M-040 requires the inputs to be priced by a reproducible fixture before they become public contract.

## Decision boundary

This change makes the CI adoption boundary and its release reference explicit. It does not add a general configuration file, alter the deterministic judge's meaning, weaken the human merge boundary, or make forge state the system of record.
