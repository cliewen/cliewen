---
id: ADR-049
type: decision
status: verified
links: [ADR-008, ADR-032, ADR-033, ADR-048, PDR-020, PDR-024, P-011, C-013]
title: A pinned source manifest and a derived target manifest give migration parity one comparable shape
author: agent
accepted-by: Flemming N. Larsen (2026-08-08, conversation)
---

# ADR-049 — Migration parity manifests

## Context and problem statement

Extraction needs to compare what a pinned source revision declared with what the current corpus carries, including evidence and deliberate exclusions or deferrals, without authoring two independent registries that can drift.

## Decision outcome

**The rehearsal writes a human/agent-authored source manifest pinned by `source-revision` and `source-location`; `clue parity` derives the target manifest from the same corpus/evidence harvest used by validation and compares the two.** A source entry records classified proof and every direction/location, or exactly one `excluded` or `draft | human | retired` disposition with justification; ADR-053 adds its source location and unique plan door. The derived target records ledger state, proof class, directions, evidence locations, and draft/Human/retired state.

Comparison reports missing criteria, orphaned target evidence, changed direction or location, stale source fingerprints, and unjustified dispositions. Invalid or ambiguous rows are usage errors; clean output is deterministic, sorted, and byte-identical for the same inputs. Parity writes neither manifest nor corpus.

**Carrier:** the source manifest schema, the shared evidence harvest, `clue parity`, ledger cross-checks, and the `clue-extract` rehearsal contract.
