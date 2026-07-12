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

Adopting Cliewen in an existing repository means transforming whatever corpus is already there — OpenSpec specs, plain READMEs, wiki exports — into the `/docs` taxonomy without losing traceability that already exists. Where does that transform live: in the `clue` binary, in one skill per source format, or in one generic skill?

## Decision outcome

**One generic skill, `clue-extract`, owns extraction.** It has two parts: a **target contract** that is source-independent (which artifacts must exist, everything born `provenance: inferred` per [ADR-010](ADR-010-provenance-field.md), `clue validate` green before the PR, the source corpus and its parallel registries deleted in the same PR, routing and skills installed) and **per-source mapping sections**, of which OpenSpec ([AN-002](../analysis/AN-002-model2diagram-extraction.md)) is the first. A new source format is a new mapping section, not a new skill.

The transform is meaning-level work — deciding what a SHALL statement asserts, merging design prose into ADRs, splitting a spec into capabilities. That is exactly the work the methodology assigns to agents with human review, so `clue` stays out of it: the binary judges the *result* (the extracted corpus must validate like any other) and never parses source formats.

**Carrier:** the `clue-extract` skill itself (agent), backed by `clue validate` as the machine judge of the output.

### Rejected: a deterministic `clue extract` command

A Go converter would need a parser per source format and would still get the meaning wrong — AN-002's target alone has three notations for one ID scheme and prose that must become Gherkin. A converter that needs human post-editing on every run is a skill wearing a binary's clothes.

### Rejected: one skill per source format

`clue-extract-openspec`, `clue-extract-readme`, … would duplicate the target contract — the part that must stay identical for `clue validate` to be the single judge — across skills that drift independently.
