---
id: PDR-068
type: decision
status: inferred
links: [G-016, CAP-001, CAP-012, PDR-067]
title: New adoption defaults to local acceptance and existing policy is preserved
author: agent
accepted-by: []
binds: adopter
---

# Local default with policy continuity

## Context

Requiring each new adopter to opt into local acceptance adds a setup decision before its first change. Changing an existing repository's implicit PR workflow on upgrade would instead change meaning without human acceptance.

## Decision

Fresh `clue init` materializes local acceptance on `main`; `--acceptance=pr` actively chooses PR and explicit local selection is supported. Init never overwrites policy and refuses conflicting explicit options before writes. Existing Cliewen machine state, canonical skills, corpus artifacts or index markers preserve the legacy PR convention when no explicit choice exists. Ambiguous convention state is treated conservatively as existing adoption.

Missing policy retains PR until init or reviewed migration records it explicitly. Migration adds PR policy for legacy repositories and preserves existing valid choices, reporting malformed policy as a finding. User-owned hubs are not rewritten. This source repository explicitly selects PR and remains PR-only. Base and candidate still require the same accepted local policy; a candidate cannot choose its own acceptance mechanism. Human acceptance, evidence, identity allocation and coordination are unchanged.

This refines PDR-067. Its shipped carriers are `internal/skills/source/shared/review-boundary.md.tmpl`, `internal/scaffold/templates/AGENTS.md`, init policy materialization and the acceptance-policy migration. The [design overview](../design/README.md) reaches the shared policy lifecycle.
