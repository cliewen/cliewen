---
id: ADR-047
type: decision
status: inferred
links: [C-007, CAP-002, ADR-039, ADR-040, PDR-019]
title: Diagram representation preserves links and assets
author: agent
accepted-by: []
---

# ADR-047 — Diagram representation preserves links and assets

## Context and problem statement

Diagram syntax alone does not establish whether a diagram is readable. A Mermaid-only prohibition also treats a repository-local SVG and an external image link as invalid even when they are the only adequate representation, making a migration remove useful documentation to obtain a green corpus.

## Decision outcome

**Diagram representation is human judgment, not a `clue validate` verdict.** Embedded Mermaid is preferred when it clearly communicates the diagram. Embedded ASCII art is preferred for structures it represents clearly, such as a directory tree. SVG is permitted when neither embedded form adequately communicates the diagram, including a complex C4 context view.

**Image links and assets are valid corpus content.** Inline, reference, collapsed-reference, and HTML image forms may target repository-local or absolute addresses, and image assets may be stored under `docs/`. `clue validate` remains offline and deterministic: it neither resolves nor fetches an image target.

**Extraction preserves link meaning.** A deterministic conversion mapping rewrites a source target to its converted target. When no mapping exists, extraction reports the target and does not delete the source that would leave it broken. A validator rule must never make an extractor remove a link or asset merely to pass.

## Rejected: retain a SVG-only machine exception

SVG is not the only valid image representation, and a type-based exception would recreate the invalid conclusion that syntax establishes whether documentation is clear. The required preference is reviewable prose, not a MIME-type policy.

## Carrier

[C-007](../constraints/C-007-diagrams-inline-mermaid.md) states the review rule. CAP-002 removes image rejection from the deterministic judge. The canonical extraction skill and OpenSpec mapping preserve links and assets through conversion.
