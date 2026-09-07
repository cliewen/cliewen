---
id: CH-171
type: change
status: open
links: [G-001, CAP-002, ADR-048, C-013]
title: Team-safe parallel identity allocation
---

# CH-171 — Team-safe parallel identity allocation

## What

Make `clue id next` safe when contributors and agents allocate from the same accepted commit in parallel clones or worktrees. A repository may opt into Git-coordinated allocation: a dedicated remote branch serializes numeric claims with ordinary fast-forward pushes, while the checked-in ledger remains the deterministic lifecycle truth that `clue validate` reads.

The same change makes the ledger append-only and union-mergeable, adds synchronization and batch reservation commands, and preserves a local mode for single-user or deliberately serialized repositories.

## Why

The current allocator is checkout-local load → increment → save. Two branches starting from the same ledger can both reserve the same next ID, and identical ledger edits may merge without a textual conflict. That makes the documented independent-author workflow unsafe for a team even though validation eventually detects duplicate artifacts.

Cliewen must support parallel team work without giving up readable sequential IDs, a network-free deterministic judge, repository-owned state, or host independence.

## Chosen shape

- A dedicated `clue/id-allocator` branch on a configured Git remote is the compare-and-swap boundary. No forge API or external service is introduced.
- The remote branch records permanent numeric claims only. It never imports speculative `live` or `retired` lifecycle state from a feature branch.
- The checked-in ledger becomes a versioned append-only event log with a built-in Git union merge rule. State folds monotonically as `reserved → live → retired`.
- Coordinated allocation fails closed. It never falls back to checkout-local allocation when the remote is unavailable, unauthorized, missing after initialization, or malformed.
- Batch reservation plus read-only synchronization supports fork contributors and planned offline work without making ranges the primary allocator.

## Compatibility

Existing IDs are never renumbered. Migration converts the current ledger to the event format while leaving repositories in local mode; enabling coordination is an explicit team decision committed before contributors allocate in parallel. The validator continues to judge only checkout bytes and never contacts a remote.

## Plan item

plan-less. This change introduces and closes one team-allocation capability in a single reviewed change.
