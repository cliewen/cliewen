---
id: CH-112
type: change
status: open
links: []
title: Preserve diagram links and assets during corpus extraction
---

# CH-112 — Preserve diagram links and assets during corpus extraction

This plan-less full change corrects C-007: its Mermaid-only machine ban caused an OpenSpec-to-Cliewen extraction to remove local SVG architecture views and their image directories. The corpus should prefer embedded Mermaid, use embedded ASCII art where that is clearest, and allow SVG where neither representation communicates the diagram adequately. Image links and assets, whether local or external, are valid corpus content.

The change removes the invalid validator prohibition, records the methodology decision and extraction finding, and makes the extraction contract preserve links and assets. A deterministic source-to-corpus mapping may retarget a link; an unmapped target prevents destructive source deletion rather than being removed or silently broken.
