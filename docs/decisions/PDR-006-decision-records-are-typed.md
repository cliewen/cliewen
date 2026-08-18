---
id: PDR-006
type: decision
status: verified
links: [PDR-003, ADR-010, P-013]
title: Decision records are typed — ADRs for architecture, PDRs for project/process, log rows for the cheap
author: agent
accepted-by: Flemming N. Larsen (2026-07-15, PR #13 approval); Flemming N. Larsen (2026-08-08, conversation — the rejected-record clause)
---

# PDR-006 — Decision records are typed

> **Superseded in part by [PDR-046](PDR-046-decisions-route-by-subject.md):** its future-shaping subject test replaces the reversal-cost and decision-log clauses; provenance, retention, subject definitions, and the rejected-record principle survive.

## Context and problem statement

Process rules filed as ADRs make the architectural term misleading, while expensive-to-reverse process decisions need more context than a log row can carry. The corpus needed distinct record types and a routing test that kept both architecture and project decisions findable.

## Decision outcome

**The original taxonomy used three types: ADRs for architecture and corpus format, PDRs for expensive-to-reverse project or process decisions, and log rows for cheap local decisions.** The record type followed reversal cost and subject; rejected future-shaping alternatives received a decision record under the same test, while cheap rejections stayed in analysis. Reclassification renamed a wrongly typed record into the right series without silently changing its text or meaning.

PDR-046 now routes every future-shaping choice by subject to ADR, PDR, or IDR and retires the log. The surviving contract is that public record names retain their ordinary meaning, durable decisions and consequential rejections remain findable, and a carrier change does not discard their provenance.
