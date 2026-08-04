---
id: C-007
type: constraint
status: active
links: []
title: Diagrams are inline Mermaid
source: clue-verify checklist
enforcement: machine
---

# C-007 — Diagrams are inline Mermaid

Diagrams in the corpus are Mermaid blocks in the markdown they illustrate — versionable, diffable, rendered wherever the corpus is read. No binary images, no externally hosted diagrams.

**Checked by:** `clue validate` ([AC-093](../capabilities/CAP-002-validate/criteria.md)) — an image link in `docs/**` prose, or an image file stored under `docs/`, fails. An image link inside a fence or a code span is an example of the form, never a diagram.
