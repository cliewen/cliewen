---
id: C-001
type: constraint
status: active
links: []
title: Markdown prose is never hard-wrapped
source: AGENTS.md rule 6
enforcement: machine
---

# C-001 — Markdown prose is never hard-wrapped

One line per paragraph and per list item; wrapping is the reader's IDE concern. Line breaks are structural only: headings, lists, tables, code fences.

**Checked by:** `clue validate`'s prose-layout lint ([AC-090](../capabilities/CAP-002-validate/criteria.md)) — two running-text lines in a row are one paragraph someone broke. Fenced and indented code, tables written with or without outer pipes, frontmatter, blockquotes, HTML blocks, and comments are structure and are read as such.
