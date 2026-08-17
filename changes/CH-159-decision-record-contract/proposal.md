---
id: CH-159
type: change
status: open
links: [P-016, M-075, CAP-001, CAP-002, CAP-003, CAP-004, C-011, C-013]
title: Ship the concise typed decision-record contract
---

# CH-159 — Ship the concise typed decision-record contract

Cliewen currently routes cheap decisions into one growing log while ADRs and PDRs carry every expensive decision in a long-form shape. P-016 leaves the replacement taxonomy, routing test, and concise record shape for this milestone to decide, and requires that the contract ship atomically with a safe adopter migration and this repository's own conversion.

Write the accepted decision that settles which subject-typed records exist, how one subject test routes each future-shaping choice to exactly one type, and what compact content a new or modified record must retain. Carry that contract through validation, initialization and scaffolding, extraction guidance, generated lifecycle skills, the public guide, and release notes. Validation must reject the retired decision log and decision filenames outside the settled taxonomy.

Add a versioned migration that inventories every legacy log row without guessing its destination. The migration must block application with guided full-change classification until each durable row has a reviewed destination and the legacy log is removed; current repositories that already satisfy the new contract remain an idempotent no-op.

Apply the same classification review to this repository: preserve each durable decision in its routed record, repair live references, retire LOG-001, remove `docs/decisions/log.md`, and limit edits to immutable completed plans to link-target repairs. Focused positive and negative evidence will cover the validator, migration, emitted scaffold, extraction contract, and generated-skill parity before M-075 is digested.
