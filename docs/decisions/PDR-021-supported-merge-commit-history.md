---
id: PDR-021
type: decision
status: verified
links: [G-001, P-009, M-039, AN-014, PDR-007, PDR-016, ARCH-001, ARCH-003, C-012]
title: Full Cliewen changes are accepted with a merge commit that preserves their branch history
author: agent
accepted-by: Flemming N. Larsen (2026-08-02, conversation)
---

# PDR-021 — Full Cliewen changes are accepted with a merge commit that preserves their branch history

## Context and problem statement

The proposal, implementation, digest, and durable corpus form a provenance chain, but integration modes can leave the same final tree while retaining different history. The repository needs a supported mode whose accepted `main` keeps that chain without making the forge the system of record.

## Decision outcome

**A full Cliewen change is accepted only through a human-controlled merge commit.** The protected default branch must allow merge commits and disable squash and rebase-and-merge, while retaining required validation, pull requests, resolved conversations, no-bypass, deletion, and force-push protections. The original proposal, implementation, digest, and corpus commits therefore remain reachable from `main`; squash loses those commits and rebase-and-merge rewrites their identities, so both are unsupported.

Rebasing an unpublished local branch onto accepted `main` remains allowed. After publication, current `main` is incorporated with the normal non-rewriting merge and checks are repeated. This decision refines PDR-007 and ARCH-001 without authorizing agents to merge or treating forge state as the repository's durable meaning. Plain changes are affected by branch-scoped enforcement but do not carry a proposal or digest chain.

The merge-boundary source, C-012, ARCH-001, ARCH-003, contributor and operations guidance, branch-protection probe, history fixture, and content guards carry this support boundary.
