---
id: IC-001
type: imported-change
status: complete
links: [CAP-003]
title: "Fixture: add-retry-budget — core allocator"
source-revision: 4f1c9e2
source-location: github.com/example/adopter-repo/openspec/changes/add-retry-budget
---

# IC-001 — Fixture: add-retry-budget — core allocator

This record is M-054's worked fixture: it proves a proposal, a design rationale, a dependency, and a proof-linked task all remain inspectable from an `imported-change` record alone, once the source repository that originated them is gone.

## Origin

Extracted from the fictitious adopter repository `example/adopter-repo`, OpenSpec pending change `changes/add-retry-budget/`, at the pinned revision and location named in this record's frontmatter. The source change proposed a retry-budget allocator shared by every outbound call site.

## Intent

Give the adopter one place that decides whether a failing call may retry, instead of each call site inventing its own backoff policy — the source proposal's stated motivation.

## Design rationale

The source's `design.md` chose a token-bucket allocator over a fixed retry count because a shared budget degrades gracefully under a correlated failure (many call sites failing at once) where a per-site fixed count does not: each site would exhaust its own retries independently and the aggregate retry volume would spike exactly when the dependency is least able to absorb it.

## Dependencies

- [IC-002](IC-002-fixture-in-progress-dependent.md) depends on this record's allocator existing before its policy hooks can be wired in.

## Proof links

| Task | Criterion |
|---|---|
| implement the token-bucket allocator | AC-054 |
| record the allocator's verification environment and population claims | AC-055 |
