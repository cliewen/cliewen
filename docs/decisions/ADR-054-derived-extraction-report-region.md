---
id: ADR-054
type: decision
status: verified
links: [ADR-049, ADR-053, PDR-020, P-012, CAP-003, C-013]
title: An extraction report's figures are a rendered region, not prose
author: agent
accepted-by: Flemming N. Larsen (2026-08-08, conversation)
---

# ADR-054 — Derived extraction report region

## Context and problem statement

An extraction report can state criterion counts and mappings that its parity manifest never contained, while both the report and manifest still pass independent checks. The report needs a durable, readable region whose figures have one pinned source.

## Decision outcome

**A report's derived figures occupy a delimited region rendered from the pinned source manifest, and `clue validate` re-renders it and fails on any difference.** The markers name the manifest relative to the repository root; `clue report` renders in place; typed or stale figures, a missing manifest, or a marker outside an example code span fail validation. The region renders the manifest alone, never the current live corpus, so historical reports remain stable as later corpus changes occur.

Prose outside the region remains free, a corpus without such a region is unaffected, and a marker in code or a fence is an example rather than a claim. Parity supplies the corresponding target comparison, so the report cannot claim a corpus that the pinned manifest did not check.

**Carrier:** report rendering, marker scanning and validation, the source manifest, and extraction-report guidance.
