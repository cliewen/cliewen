---
id: ADR-038
type: decision
status: verified
links: [P-009, M-040, AN-012, AN-014, CAP-001, CAP-006, CAP-004, ADR-030, PDR-021]
title: The CI wall is an upstream reusable workflow with a thin caller
author: agent
accepted-by: Flemming N. Larsen (2026-08-02, conversation)
---

# ADR-038 — Upstream validation workflow

## Context and problem statement

Adopters need the same validation boundary without copying Cliewen's workflow, while the reference must be immutable and the caller must retain local runner and acquisition choices.

## Decision outcome

**Cliewen owns one reusable `.github/workflows/clue-validation.yml`, and an adopter owns only a thin `.github/workflows/clue.yml` caller.** The caller references a full source commit when the emitting binary has trustworthy repository VCS metadata, otherwise the corresponding protected release tag; neither reference is mutable. The caller carries triggers, the stable `validate` job name, the upstream reference, `clue-version`, runner labels, `clue-source`, and a writable install directory, while the reusable workflow owns checkout, changed-surface handling, acquisition, checksum verification, executable staging, validation, and the acceptance-brief gate.

The source revision is emitted only when the tree identifies this project, is the answering repository root, tracked files match, and Go reports no dirty build; publication is guaranteed only by the release path. Vendored and release sources both verify `SHA256SUMS`, action references use full commit SHAs, and the workflow never assumes a root-only install path. The required status check and human merge boundary remain protected workflow/forge concerns.

The pre-1.0 repository replaces copied validation walls rather than maintaining compatibility with an installed shape that does not exist.

**Carrier:** the reusable workflow, generated caller, generator and release metadata, synthetic fixtures, and the validation guidance.
