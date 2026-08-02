---
id: C-016
type: constraint
status: active
links: [G-001, ADR-041, CAP-002, CAP-005]
title: A generated index row states the record it links, never just its filename
source: ADR-041
enforcement: machine
---

# C-016 — A generated index row states its record

Given a taxonomy README's `clue:index` block, when index generation appends a row for an artifact, then the row states the artifact's own `id — title` from its frontmatter followed by its `status` — never the target's filename restated as a label.

A row referencing a subfolder README states a section rather than a record and carries no title or status. An artifact whose frontmatter carries no readable id and title falls back to the plain link, because index generation reports nothing and fails on nothing; naming a malformed artifact is the judge's job.

The rule governs a row's opening only. Curated text after the status — supersession notes, section descriptions — is part of the contract rather than an exception to it, and regeneration preserves every row whose target still exists.

**Enforcement:** `machine` — index generation emits the stated form (AC-073), and `clue validate` counts rows whose label restates only their own link, naming them behind `--index-rows` (AC-074). The count is reported and never failed on, for the reasons in [ADR-041](../decisions/ADR-041-index-rows-state-their-record.md); under [C-004](C-004-never-weaken-checks.md) it is never softened to make the number look better.
