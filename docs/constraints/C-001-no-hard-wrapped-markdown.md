---
id: C-001
type: constraint
status: active
links: []
title: Markdown prose is never hard-wrapped
source: AGENTS.md rule 5
enforcement: agent
---

# C-001 — Markdown prose is never hard-wrapped

One line per paragraph and per list item; wrapping is the reader's IDE concern. Line breaks are structural only: headings, lists, tables, code fences.

**Promotion trigger:** `clue validate` gains a prose-layout lint that flags mid-paragraph line breaks in `docs/**` markdown — then `enforcement: machine`.
