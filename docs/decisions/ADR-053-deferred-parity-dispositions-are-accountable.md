---
id: ADR-053
type: decision
status: verified
links: [ADR-049, PDR-024, P-012, CAP-003, C-013]
title: Deferred parity dispositions name their source and plan door
author: agent
accepted-by: Flemming N. Larsen (2026-08-08, conversation)
---

# ADR-053 — Deferred parity dispositions are accountable

## Context and problem statement

Parity can deliberately defer a source criterion as `draft`, `human`, or `retired`, but a readable justification alone does not identify the source material or the milestone responsible for resolving the deferral.

## Decision outcome

**Every deferred source-manifest disposition carries `disposition-source-location` and a unique `plan-door` in addition to its justification.** The source location identifies the material that warranted the disposition; the plan door identifies a milestone whose declared identity must exist exactly once in the target corpus. `clue parity` derives the target milestone set and rejects an absent or reused door.

Parity reports the count of unique criterion IDs with dispositions on both clean and failing runs. The count is visible backlog, not a threshold, and the contract applies only to migration source manifests; `clue validate` remains a one-corpus judge.

**Carrier:** the source-manifest schema, parity validation, and the extraction plan/milestone corpus.
