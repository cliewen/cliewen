---
id: ADR-041
type: decision
status: verified
links: [CAP-002, CAP-005, ADR-017, ADR-019, ADR-035, ADR-046, C-004, C-013, C-016]
title: Generated index rows state their record, and rows that state only their link are counted
author: agent
accepted-by: Flemming N. Larsen (2026-08-08, conversation)
---

# ADR-041 — Generated index rows state their record

## Context and problem statement

An index row that names only a filename opens the right file but does not tell a reader which record or status it represents. The generator can correct newly appended rows, while existing adopter-owned rows must not be rewritten or turned into merge failures.

## Decision outcome

**An appended artifact row is `- [<id> — <title>](<file>) · \`<status>\``, using parsed frontmatter values.** A malformed artifact missing `id`, `title`, or `status` degrades to a plain link; folder rows remain plain because they name sections rather than records; and existing rows whose targets still exist are preserved. **Amended by [ADR-064](ADR-064-regeneration-owns-the-index-badge.md):** a kept row's badge is refreshed from its artifact on regeneration, because it is a copy of a frontmatter field rather than anything an author decided; the rest of the row is still preserved.

`clue validate` counts rows whose label is only the target filename stem and `--index-rows` lists them, but the population is not an `Issue`. The count is a repair backlog rather than a gate because the generator emitted the old rows and the adopter may have curated their text.

**Carrier:** `regenIndex`, the `IndexRowBacklog` and `--index-rows` surfaces, C-016, and CAP-002/CAP-005 criteria.
