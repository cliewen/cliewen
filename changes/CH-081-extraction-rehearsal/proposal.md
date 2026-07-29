---
id: CH-081
type: change
status: open
links: [P-008]
title: Extraction rehearses before it mutates
---

# CH-081 — Extraction rehearses before it mutates

## What

Require every brownfield extraction to begin with a report-only rehearsal in its full-change workspace. The rehearsal inventories the source and proposed conversion, records uncertainty and reversal cost, identifies conflicts and plan doors, and leaves the target repository, its source corpus, routing, tests, and hosted state untouched. Mutation may begin only after explicit human direction; its durable extraction report then digests the rehearsal findings.

This change establishes that workflow boundary in a PDR, updates the canonical and generated extraction guidance and public guidance, adds an acceptance criterion and focused evidence, and records the completed M-034 bookkeeping and user-facing release note in its digest.

## Why

AN-003 found that extraction turns uncertain brownfield interpretation into durable truth and can introduce governance changes. A reversible rehearsal makes the interpretation, conflicts, and proposed destructive work inspectable before the source corpus or repository workflow is changed, satisfying P-008/M-034 without designing an adoption interface prematurely.

## Scope

This is a full methodology change serving P-008 milestone M-034. It does not perform an extraction against another repository, select an adopter, or implement the configuration and distributed-work interfaces reserved for M-035 and M-036.
