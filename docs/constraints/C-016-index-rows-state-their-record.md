---
id: C-016
type: constraint
status: active
links: [G-001, ADR-041, ADR-046, CAP-002, CAP-005]
title: A generated index row states the record it links and says what it is about, never just its filename
source: ADR-041
enforcement: machine
---

# C-016 — A generated index row states its record

Given a taxonomy README's `clue:index` block, when index generation appends a row for an artifact, then the row states the artifact's own `id — title` from its frontmatter followed by its `status` — never the target's filename restated as a label.

A row referencing a subfolder README states a section rather than a record and carries no title or status. An artifact missing any of a readable id, title, and status falls back to the plain link — a row is one shape or the other, never a third carrying an empty status badge — because index generation reports nothing and fails on nothing; naming a malformed artifact is the judge's job.

The row also says what the artifact is about. Index generation seeds that sentence from the artifact's own body when it appends the row, and an artifact with no readable sentence keeps the shorter form rather than carrying an empty tail ([ADR-046](../decisions/ADR-046-index-rows-say-what-the-artifact-is-about.md)). The seed is a first draft and never an assertion: regeneration rewrites nothing that already exists, and no command backfills a description into rows a corpus already carries.

Curated text after the status — the description, supersession notes, section descriptions — is part of the contract rather than an exception to it, and regeneration preserves every row whose target still exists.

**Checked by:** index generation, which emits the stated form (AC-073) and seeds the description (AC-096, AC-097), and `clue validate`, which counts rows whose label restates only their own link and rows that state their record but say nothing about it, naming both behind `--index-rows` (AC-074, AC-098). Both counts are reported and never failed on, for the reasons in [ADR-041](../decisions/ADR-041-index-rows-state-their-record.md) and [ADR-046](../decisions/ADR-046-index-rows-say-what-the-artifact-is-about.md); under [C-004](C-004-never-weaken-checks.md) neither is softened to make the number look better.
