---
id: CH-122
type: change
status: open
links: [P-012, M-057, PDR-026, P-011, AN-016]
title: Re-derive and record the brownfield migration assessment gates
---

# CH-122 — Re-derive and record the brownfield migration assessment gates

Serves [M-057](../../docs/plans/P-012-migration-gap-closes-on-evidence.md) of [P-012](../../docs/plans/P-012-migration-gap-closes-on-evidence.md).

## What

This change makes P-012's starting evidence inspectable and truthful. It records a sanitized version of the originating assessment as an analysis artifact; re-derives every assessment gap and blocking gate from the current corpus and commands; classifies each as closed, declined, or open under PDR-026; seeds this repository's identity ledger through the shipped migration; repairs CAP-003's stale adopter-CI distribution statement; and records the limit of branch protection as acceptance evidence.

## Why

P-011's completed milestone table describes mechanisms the campaign delivered, but PDR-026 requires its outcome to be re-derived rather than accepted as self-reported. P-012 exists because that pass found gaps that remain open or were deliberately declined. Without a durable, sanitized assessment record and a current gate-by-gate disposition, the successor campaign would repeat the same failure: prose claims about what its predecessor established.

## Design

The assessment record will distinguish supplied observations from Cliewen evidence, name no adopter, and link each assessment item to its present disposition. The gate register will make the three PDR-026 states explicit: named mechanism plus failure-path evidence for a closure, an explicit declined request and its adopter cost, or the P-012 milestone that remains open.

The identity ledger is seeded by the existing `clue migrate --apply` contract, not hand-authored. This exercises the repository's shipped migration and removes the now-inapplicable MIG-008 offer. The branch-protection boundary is a process decision: protection can enforce a PR's admission conditions, but cannot itself become the human acceptance evidence for a criterion.

## Scope

- A sanitized analysis artifact and its index row.
- A process decision defining the branch-protection evidence boundary and its live carriers.
- CAP-003 extraction guidance and the plan's M-057 evidence bookkeeping.
- The generated `.clue/id-ledger.yaml` migration result.
- A user-facing `[Unreleased]` entry because the shipped migration's own repository now exercises its identity-ledger path and the extraction guidance corrects a user workflow claim.

## Not in this change

No validation or parity rule changes. M-058 through M-061 remain separate changes: this change only makes their open/declined status explicit and establishes an honest evidence base for them.
