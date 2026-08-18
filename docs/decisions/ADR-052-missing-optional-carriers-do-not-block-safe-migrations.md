---
id: ADR-052
type: decision
status: verified
links: [ADR-038, ADR-039, CAP-001, P-012]
title: A missing optional generated carrier does not block an independent safe migration
author: agent
accepted-by: Flemming N. Larsen (2026-08-06, conversation)
---

# ADR-052 — Missing optional carriers do not block safe migrations

## Context and problem statement

An adopter may lack the thin generated CI caller while another migration is independently safe. Treating that absence like an ambiguous or modified carrier would either block safe work or give migration authority over adopter-owned runner and installation choices.

## Decision outcome

**A missing thin CI caller is a notice, not a blocking migration finding; migration names `clue init` as its materialization route and never creates or rewrites the caller.** An independent safe change may apply. A present caller whose content cannot be safely recognized remains blocking, as do ambiguous corpus meaning and locally modified managed carriers; preflight remains atomic.

**Carrier:** `planCaller`, the notice and safe-write behavior, AC-124, and CAP-001 guidance.

## Rejected alternatives

Creating the caller during migration would transfer adopter-owned choices from `init`; allowing every finding to permit partial writes would weaken the atomic boundary for unsafe semantic or modified-carrier states.
