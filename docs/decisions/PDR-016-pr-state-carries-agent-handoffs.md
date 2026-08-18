---
id: PDR-016
type: decision
status: verified
links: [G-001, AN-009, PDR-007, PDR-012, PDR-040, PDR-042, C-012, CAP-006]
title: Hosted PR state carries review findings and updater handoffs across agents
author: agent
accepted-by: Flemming N. Larsen (2026-08-02, conversation)
---

# PDR-016 — Hosted PR state carries review findings and updater handoffs across agents

> **Scope amended by [PDR-042](PDR-042-routing-recommends-contract-aware-effort.md):** these handoff guarantees govern a chosen full loop; simple work follows user authority and repository policy.

> **Amended by [PDR-040](PDR-040-push-is-durability-ready-is-explicit.md):** every changed turn pushes, a repair to a ready PR returns it to draft, and a change is first published from its first commit; accepted `main` is merged into the branch rather than hidden by a rebase.

## Context and problem statement

Private conversation state cannot carry an unfinished finding or repair across agents, while the review boundary requires the exact hosted head to be known. The pull request must therefore hold the durable handoff without serializing independent authors.

## Decision outcome

**The pull request, its hosted head, and resolvable review conversations are the shared handoff state for a full change.** Reviews name the SHA they inspected; edits invalidate that result; actionable findings remain unresolved until a hosted commit contains the verified repair, where the forge supports that enforcement.

Any agent that edits an open PR becomes its updater for that turn: it fetches and records the current head, repairs and verifies the complete change, commits and pushes without force, obtains a clean review of the resulting SHA, confirms the hosted head matches it, and then resolves satisfied findings. A changed head or rejected non-fast-forward push requires reconciliation and renewed checks and review. The updater role follows mutation, not the identity of the first implementer.

Independent changes still branch from accepted `main`. After publication, accepted `main` advances by merging into the open PR branch and repeating verification and review; before first publication, rebasing onto current `main` remains allowed. A merged or closed PR stops updates and reports preserved local work rather than creating a follow-up without human scope.

The shared review-boundary source, `clue-verify`, pull-request template, required-conversation setting, CAP-006, and generated-contract tests carry the handoff. CI cannot detect uncommitted local work or a private finding.
