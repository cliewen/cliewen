---
id: C-001
type: constraint
status: active
links: []
title: Markdown prose is never hard-wrapped
source: Cliewen methodology (scaffolded by clue init)
enforcement: agent
---

# C-001 — Markdown prose is never hard-wrapped

One line per paragraph and per list item; wrapping is the reader's IDE concern. Line breaks are structural only: headings, lists, tables, code fences.

The generated AGENTS.md declares this rule as its rule 5; a repo that kept its own AGENTS.md holds it through this artifact — the register entry is the rule's authoritative declaration either way.

**Promotion trigger:** `clue validate` gains a prose-layout lint that flags mid-paragraph line breaks in `docs/**` markdown — then `enforcement: machine`.
