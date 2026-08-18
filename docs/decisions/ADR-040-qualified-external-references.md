---
id: ADR-040
type: decision
status: verified
links: [P-009, AN-013, ADR-005, ADR-009, ADR-039, C-013, CAP-002, CAP-004]
title: External references name their target, and resolving them stays outside the judge
author: agent
accepted-by: Flemming N. Larsen (2026-08-02, conversation)
---

# ADR-040 — Qualified external references

## Context and problem statement

Corpus references cross repository and forge boundaries, but bare issue numbers and unqualified IDs can silently name the wrong target. Network resolution would make the deterministic judge online, non-reproducible, and unavailable during outages.

## Decision outcome

**An external reference names its target, while the judge enforces only its form; a separate resolver reports target state and never gates validation or a merge.** Addresses are full URLs, with readable labels free to retain forge shorthand; bare `#N` is invalid. Foreign corpus identities use `clue:owner/repo/ID`, and foreign acceptance evidence adds its pinned revision. Local IDs and relative paths remain unchanged. Code spans, fences, link targets, and heading anchors are excluded from reference scanning.

The resolver is advisory and read-only. It reports `reachable`, `restricted`, `redirected`, `gone`, or `unreachable`; only `gone` is an error, and private or unavailable targets never fail the corpus. Credentials are sent only to their matching API host after an unauthenticated gone result, and pinned historical references are reported but never rewritten. Foreign evidence remains named but locally unproven; `Test-type: Human` remains its proof of record.

**Carrier:** the offline form rule in CAP-002, the resolver command in CAP-004, shared reference scanning, migration guidance, and the separate coverage listing for foreign evidence.
