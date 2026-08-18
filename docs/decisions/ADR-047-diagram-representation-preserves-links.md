---
id: ADR-047
type: decision
status: verified
links: [C-007, CAP-002, ADR-039, ADR-040, PDR-019]
title: Diagram representation preserves links and assets
author: agent
accepted-by: Flemming N. Larsen (2026-08-08, conversation)
---

# ADR-047 — Diagram representation preserves links

## Context and problem statement

Corpus diagrams need a readable representation and extraction must not destroy the links or assets that make a diagram meaningful, while validation must remain offline and deterministic.

## Decision outcome

**Embedded Mermaid is preferred when it communicates the diagram clearly, embedded ASCII is preferred for structures such as directory trees, and SVG is permitted when neither is adequate, including complex C4 views.** Inline, reference, collapsed-reference, and HTML image forms may target repository-local or absolute addresses; image assets may live under `docs/`, and `clue validate` neither resolves nor fetches them.

Extraction preserves link meaning by applying a deterministic source-to-target mapping. When no mapping exists it reports the target and leaves the source link or asset intact rather than deleting content merely to satisfy a validator.

**Carrier:** diagram guidance, extraction mapping and findings, and the offline link/asset validation boundary.
