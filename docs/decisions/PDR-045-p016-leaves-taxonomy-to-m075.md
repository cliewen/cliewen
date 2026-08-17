---
id: PDR-045
type: decision
status: verified
links: [P-016, PDR-044]
title: P-016 leaves the decision-record contract to M-075
author: human
accepted-by: Flemming N. Larsen (2026-08-17, commit 35e734c and conversation)
---

# PDR-045 — P-016 leaves the decision-record contract to M-075

## Context and problem statement

P-016 named the taxonomy and routing outcome that M-075 was supposed to decide, while its closing batch covered new IDRs but not other decision records created during the campaign.

## Decision outcome

P-016 states the questions M-075 must settle and the evidence it must produce, not their answers. Its final compaction batch covers every decision record created by earlier milestones, regardless of type.

## Consequences

- M-075 remains responsible for an accepted decision defining the record types, routing test, and compact shape.
- The campaign cannot close while any decision record it created remains outside the compaction audit.
