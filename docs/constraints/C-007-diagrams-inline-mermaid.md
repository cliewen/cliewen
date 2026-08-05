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

**Checked by:** `clue validate` ([AC-093](../capabilities/CAP-002-validate/criteria.md)) — an image in `docs/**` prose, or an image file stored under `docs/`, fails. Every form an image takes is one form: inline, reference and collapsed-reference links, and an `<img>` tag, which is how an externally hosted diagram usually arrives. An image inside a fence or a code span is an example of the form, never a diagram.
