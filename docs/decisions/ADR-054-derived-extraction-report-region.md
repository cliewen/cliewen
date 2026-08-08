---
id: ADR-054
type: decision
status: verified
links: [ADR-049, ADR-053, PDR-020, P-012, CAP-003, C-013]
title: An extraction report's figures are a rendered region, not prose
author: agent
accepted-by: Flemming N. Larsen (2026-08-08, conversation)
---

# ADR-054 — An extraction report's figures are a rendered region, not prose

## Context and problem statement

PDR-020 requires the mutate phase to digest the rehearsal into a durable extraction report under `/docs/analysis`, and ADR-049 makes `clue parity` compare a pinned source manifest against the corpus-derived target manifest. The report's own criterion counts and mapping tables are outside both mechanisms: they are prose an agent types, so a report can state a criterion population, a preserved-and-minted mapping, or a deferral set that the manifest beside it never contained, and every check still passes. P-012 M-059 requires that a report claiming a different corpus than the one committed not be expressible. What makes the report's figures checkable without turning the report into a second registry?

## Decision outcome

**An extraction report states its figures in one delimited region rendered from the pinned source manifest, and `clue validate` re-renders that region and fails on any difference.**

- The region is delimited by `<!-- clue:derived-from: <manifest-path> -->` and `<!-- clue:derived-end -->`. The opening marker names the source manifest, relative to the repository root, so the report declares its own origin and no second index maps reports to manifests.
- `clue report <report-path>` renders the region in place from the manifest the marker names. The figures are therefore generated; typing them is not a supported way to produce them.
- `clue validate` re-renders every region it finds in the corpus and reports a difference as an issue. A typed figure, a stale region left behind by a revised manifest, and a marker naming a manifest that is not there all fail the required check rather than the optional one.
- The region renders the source manifest alone and never the live corpus. A durable report is history: it must stay stable as the corpus moves on, and a region that re-derived itself from the present tree would go stale on every later change to an unrelated criterion.
- A marker written inside a code span or a fenced block is an example, not a region. A decision record, a skill, and a design note have to be able to state the contract, and a scan that could not tell an example from a claim would forbid describing it — the same exemption the outward-reference scan already makes.
- Prose outside the region is free. The region holds no fact of its own — it is a rendering of a file that is already committed beside it — which is what keeps the report a readable document instead of a registry that must itself be maintained.

Agreement between the report and the committed tree follows from the two checks together: the region is rendered from the same manifest `clue parity` compares against the derived target manifest, so a report cannot claim a corpus that parity did not check.

A corpus with no such region is unaffected. This contract reaches extraction reports, not greenfield analysis records.

## Rejected: render the region from the current corpus

It would let `clue validate` alone assert agreement with the tree, but every later change touching any criterion would invalidate every historical extraction report, and regenerating a merged migration's report to describe a corpus it never described would destroy the record it exists to keep.

## Rejected: check the figures with a lint rule over ordinary prose

Recognizing which sentence states a criterion count is natural-language parsing, and it would leave the counts hand-maintained — the thing M-059 identifies as the defect.

## Carrier

`internal/parity`, the `clue report` and `clue validate` command surface, AC-126 and AC-127 in CAP-003, and the canonical `clue-extract` skill carry this decision.
