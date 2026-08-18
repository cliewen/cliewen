---
id: PDR-002
type: decision
status: verified
links: [PDR-003, PDR-042]
title: A light change tier — the PR description is the proposal
author: agent
accepted-by: Flemming N. Larsen (2026-07-14, PR #9)
---

# PDR-002 — A light change tier: the PR description is the proposal

> **Superseded by [PDR-042](PDR-042-routing-recommends-contract-aware-effort.md):** light and plain are now the simple recommendation outside the loop, while accepted-contract change is recommended for the full loop.

## Context and problem statement

The full change workspace is disproportionate for changes that preserve decisions, acceptance criteria, plan promises, and methodology carriers. The smallest honest route still needs a reviewable branch and pull request without transient proposal files.

## Decision outcome

**The original light tier used a branch and pull request whose description carried the proposal, with no `/changes/` workspace; all other work used the full loop.** The qualifying test was semantic rather than a file or line threshold: no new decision, acceptance-criterion or capability meaning, semantic plan mutation, or methodology-carrier change. Discovery of any such meaning required immediate escalation to the full loop.

PDR-042 replaces this light/full vocabulary with the simple/full recommendation and defines the current simple exceptions. The surviving carrier principle is that contract-preserving work carries only the metadata, workspace, evidence, and checks relevant to its surfaces; CH numbering remains global for changes that enter Cliewen.
