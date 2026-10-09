---
id: PDR-067
type: decision
status: inferred
links: [G-016, CAP-012, ARCH-003, PDR-021, C-012]
title: Adopters may explicitly select human-controlled local acceptance
author: agent
accepted-by: []
binds: adopter
---

# PDR-067 — Local human acceptance

## Context

A PR provides a hosted review and admission boundary, but requiring a forge excludes repositories that can retain the same proposal history and explicitly accept a reviewed branch locally. Identity allocation and acceptance are independent concerns.

## Decision

PR acceptance remains the default. An adopter may select local acceptance on its accepted base and in its repository instructions. A human runs `clue accept` after verification and review; the command creates a merge whose tree is exactly the reviewed candidate, whose parents retain the base and candidate, and whose message retains the complete acceptance brief. Preparation and cancellation are not acceptance. The command never pushes or silently incorporates a newer base. Agents never run its accepting form for their own change.

This refines PDR-021 and C-012 for opted-in adopters only. Source repositories retain their PR requirement. Local acceptance does not provide a forge-enforced admission gate or authenticate human presence: the prompt, Git identity, verification results, and review declarations retain explicitly stated procedural trust. `clue validate` still reads repository state; `clue accept` owns history and integration checks. The ID ledger, allocation modes, and coordination requirements remain unchanged.

The shipped carriers are `internal/skills/source/shared/review-boundary.md.tmpl`, `internal/skills/source/skills/clue-delta.md.tmpl`, the local acceptance reference, and `internal/scaffold/templates/acceptance/brief.md`. The [design overview](../design/README.md) links the record and recovery contract.
