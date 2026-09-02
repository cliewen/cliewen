---
id: PDR-040
type: decision
status: verified
links: [G-001, PDR-007, PDR-012, PDR-016, PDR-033, PDR-042, C-004, C-012, CAP-006]
title: Push is durability and the ready mark is the only readiness claim
author: agent
accepted-by: Flemming N. Larsen (2026-09-02, conversation)
---

# PDR-040 — Push is durability, ready is the explicit act

> **Scope amended by [PDR-042](PDR-042-routing-recommends-contract-aware-effort.md):** automatic branch publication and the ready mark govern a full loop the user chose; a simple route supplies no push authority by itself.

## Context and problem statement

Deferring publication until review left repaired work in a private worktree, while treating a push as readiness made unfinished work look accepted and forced exception machinery around non-converging review.

## Decision outcome

**A push claims nothing; the ready mark claims readiness.** Every changed turn commits and pushes its branch, and the full-change pull request is a draft from first publication. The PR becomes ready only when local verification and the review loop pass on the current commit and the hosted head equals it; any substantive edit returns it to draft. Hosted history is not rewritten, and accepted `main` is incorporated by a normal merge.

Draft CI may omit digest and acceptance-brief gates because a draft cannot merge; those gates bind on `ready_for_review` and `main`. The caller declares the draft-aware behavior fail-closed, so updating the reusable workflow without the trigger cannot leave a mergeable PR with only a lenient result. Nothing here changes the human merge boundary, deterministic judge, or verifiable thread.

## Rejected: make push a readiness signal or keep work local

Push-as-readiness conflicts with unresolved review and needs stopping exceptions; local work cannot be handed off or recovered. Durability and readiness are separate claims.

## Carrier

The review-boundary source and `clue-delta`/`clue-verify` templates, CI callers and reusable workflow, C-012, pull-request templates, `CONTRIBUTING.md`, the change-loop guide, and CAP-006 carry the rule; PDR-007, PDR-012, PDR-016, and PDR-033 receive the stated scope amendments.
