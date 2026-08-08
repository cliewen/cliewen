---
id: ADR-053
type: decision
status: verified
links: [ADR-049, PDR-024, P-012, CAP-003, C-013]
title: Deferred parity dispositions name their source and plan door
author: agent
accepted-by: Flemming N. Larsen (2026-08-08, conversation)
---

# ADR-053 — Deferred parity dispositions name their source and plan door

## Context and problem statement

ADR-049 required a migration parity disposition to carry a free-text justification. A repeated sentence can satisfy that shape for every criterion while leaving neither the original source location nor the work that will resolve the disposition inspectable. P-012 M-058 requires a disposition to be honest or fail.

## Decision outcome

**Every `draft`, `human`, or `retired` source-manifest disposition carries `disposition-source-location` and `plan-door` in addition to its readable `justification`.** The source location points to the particular source material that warranted the disposition. The plan door is a declared milestone identity such as `M-060`, unique to that deferred criterion; `clue parity` derives the target corpus's declared milestone set and fails an otherwise matching disposition whose door is absent or reused.

`clue parity` reports the count of unique criterion IDs carrying a disposition on both clean and failing runs. The count is derived from the authored source manifest and is a visible backlog, not a threshold: a migration may defer several criteria when each has its own accountable record, and a single genuine disposition continues to pass.

This contract applies only to source manifests consumed by migration parity. `clue validate` remains a one-corpus judge, and a greenfield criterion has no source migration location or plan door to supply.

## Rejected: infer accountability from the justification

Natural-language parsing cannot establish whether a location or plan door exists, and it would preserve the repeated-prose loophole that M-058 identifies.

## Carrier

`internal/parity`, AC-125 in CAP-003, the `clue parity` report, and the canonical `clue-extract` mapping guidance carry this decision.
