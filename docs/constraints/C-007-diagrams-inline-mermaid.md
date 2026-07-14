---
id: C-007
type: constraint
status: active
links: []
title: Diagrams are inline Mermaid
source: clue-verify checklist
enforcement: agent
---

# C-007 — Diagrams are inline Mermaid

Diagrams in the corpus are Mermaid blocks in the markdown they illustrate — versionable, diffable, rendered wherever the corpus is read. No binary images, no externally hosted diagrams.

**Promotion trigger:** `clue validate` flags image links and image files under `docs/**` — then `enforcement: machine`.
