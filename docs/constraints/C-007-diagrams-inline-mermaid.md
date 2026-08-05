---
id: C-007
type: constraint
status: active
links: []
title: Diagrams use the clearest renderable form
source: ADR-047
enforcement: human
---

# C-007 — Diagrams use the clearest renderable form

Use an embedded Mermaid diagram when it clearly communicates the diagram. Use embedded ASCII art for a structure such as a file tree when that is clearer. Use SVG when neither embedded form adequately communicates the diagram, including a complex C4 context view.

Local and absolute image links, HTML image tags, and image assets under `docs/` are valid corpus content. Keep an existing link and its target through extraction; a deterministic source-to-corpus mapping may retarget it, but an unmapped target prevents destructive source deletion rather than being removed or silently broken.

**Residual:** choosing the clearest form depends on the diagram's meaning and its reader. A machine can see the syntax, but cannot determine whether a file tree is clearer as ASCII art or whether a C4 view needs SVG. The cost of a poor choice is documentation that is less readable or loses necessary detail; human review holds that judgment.
