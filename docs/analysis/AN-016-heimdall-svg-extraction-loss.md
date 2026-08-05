---
id: AN-016
type: analysis
status: active
provenance: inferred
reversal-cost: low
reality: contradicted
links: [CAP-002, C-007, ADR-039, ADR-040, PDR-020]
title: Heimdall extraction removed SVG architecture views to satisfy C-007
---

# AN-016 — Heimdall extraction removed SVG architecture views to satisfy C-007

## Finding

The OpenSpec-to-Cliewen extraction recorded at [Heimdall commit `b443613c`](https://github.com/elsevier-research/iip-heimdall/commit/b443613c) removed seven `![...](./images/....svg)` links from C4 and C5 architecture-view Markdown files and deleted their `images/` directories. The accompanying extraction analysis states that these removals were made to satisfy C-007, while the Structurizr DSL sources remained.

This was a green-corpus contradiction: C-007 and AC-093 made the validator reject a valid local documentation representation, and the extraction contract did not explicitly protect its links and assets. The loss is not repaired in this methodology change; a separate Heimdall change restores the deleted views.

## Evidence boundary

The finding is a read-only inspection of the named commit in the local Heimdall checkout on 2026-08-05. It establishes the committed extraction result and its recorded reason, not whether a hosted pull request has since restored the assets.

## Rejected response

Requiring all diagram authors to translate SVG to Mermaid was rejected. A complex C4 context view can need SVG, and a format conversion cannot establish that the result preserves the source diagram's clarity or meaning.

## Consumer

ADR-047 changes the diagram contract, and the extraction skill now preserves links and assets or stops destructive source deletion when it cannot map a target.
