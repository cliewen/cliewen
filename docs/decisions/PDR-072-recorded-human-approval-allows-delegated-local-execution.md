---
id: PDR-072
type: decision
status: inferred
links: [G-016, G-027, G-028, CAP-012, C-012, PDR-067]
title: Recorded human approval permits delegated local execution with safe committed snapshots
author: agent
accepted-by: []
binds: adopter
---

# Human decision and delegated execution

## Context

Requiring the owner to execute a terminal command adds manual work after the owner has already approved an exact reviewed candidate. The source repository's internal skill mirrors also make a blanket tracked-link refusal incompatible with its local acceptance policy. Delegation must retain the human decision, exact approved inputs and complete history, rather than become unattended agent authorization.

## Decision

A human may approve an exact candidate in a conversation or another explicitly identified venue. The agent records the actual statement and may mechanically execute local acceptance with a strict approval record containing candidate, base, complete brief SHA-256, decision, approver, source and recording time. The command validates these bindings and repeats the existing freshness, corpus, history and locked-reference checks. It retains the full brief and approval record in the original-candidate-tree merge. It never pushes; an agent pushes accepted main only with explicit user authorization. A changed candidate, base or brief requires a new approval. Default interactive confirmation remains available. No agent decides or invents its own approval.

The record is procedural provenance, not cryptographic authentication of a human or venue. The existing prompt and Git identity have the same authenticity limitation. Verification and observations remain recorded claims; neither approval path executes tests or certifies their meaning. This refines the local-execution parts of the [human acceptance boundary](../constraints/C-012-agents-never-merge-own-changes.md).

Snapshot construction resolves internal Git links against the committed revision and copies their regular file content into a private snapshot. It never follows operating-system links or live checkout targets. External targets, missing targets, metadata paths, cycles, excessive expansion and submodules are refused. The acceptance tree preserves the original Git objects and link modes. The [design overview](../design/README.md) and [local acceptance design](../capabilities/CAP-012-local-acceptance/design.md) describe the boundary.

The shipped carrier is `internal/skills/source/shared/review-boundary.md.tmpl`, together with the local command, approval template and scaffold hub. This replaces the earlier rule that the human must personally execute the local accepting command; the human-controlled acceptance decision remains mandatory.
