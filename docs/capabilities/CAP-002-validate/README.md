---
id: CAP-002
type: capability
status: active
links: [G-001]
title: clue validate — the deterministic judge for the corpus
goal: G-001
---

# CAP-002 — `clue validate`

## What

`clue validate [--forbid-changes] [--coverage] [--reality-gaps] [path]` scans `docs/` and `changes/` for frontmatter artifacts and fails (exit 1) on any breach of the corpus rules: missing frontmatter or core fields, UTF-8 byte-order marks, leftover second frontmatter blocks, duplicate IDs, unresolvable `links` (milestones and live acceptance criteria also resolve), status values outside the per-type vocabulary, unbounded high-cost inferred meaning behind active capabilities, malformed incident edges, `/docs` folders without README.md, index-block drift, and — with `--forbid-changes` — the presence of a transient workspace. The optional reports derive proof coverage and capabilities contradicted by incident analyses without committing registries.

## Why

The judge actor of [G-001](../../goals/G-001-verifiable-thread.md): machines enforce form so humans only have to verify meaning. The same binary runs locally and in CI — a green local run means a green PR.

Acceptance criteria: [criteria.md](criteria.md) · implementation notes: [design.md](design.md).

## Status note

`active`: implemented and covered by Go tests whose names carry the AC IDs and evidence classifications. AC-009 and its positive and negative unit tests enforce the active-criterion reference rule; criteria declaring a machine proof type additionally require supported evidence classified by that type and direction, unannotated legacy criteria retain one reference, Human-class criteria route proof to the acceptance brief, and `@draft` exempts one criterion. `clue validate` validates these declarations and references but does not execute tests.
