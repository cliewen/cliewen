---
id: ADR-008
type: decision
status: verified
links: [AN-002, CAP-003, G-001]
title: Brownfield extraction is one generic skill with per-source mappings
author: agent
accepted-by: Flemming N. Larsen (2026-07-12)
---

# ADR-008 — Extraction is a skill, generic with per-source mappings

## Context and problem statement

Brownfield adoption must transform varied source corpora into the `/docs` taxonomy without losing traceability or making the deterministic CLI interpret source-specific meaning.

## Decision outcome

**One generic `clue-extract` skill owns extraction.** Its source-independent target contract requires the target artifacts, born provenance (`provenance: inferred` for non-decisions and `status: inferred` for decisions), a validating corpus, same-change removal of source registries, and installed routing and skills. Per-source mapping files describe formats such as OpenSpec; a new format adds a mapping, not a skill.

Meaning-level decisions remain agent work under human review. `clue validate` judges the resulting corpus and never parses source formats.

**Carrier:** `clue-extract`, backed by `clue validate`. A deterministic `clue extract` command is rejected because source parsers cannot reliably decide meaning; one skill per format is rejected because it duplicates and drifts the target contract.
