---
id: CAP-005
type: capability
status: active
links: [G-001]
title: Index generation — README index blocks regenerate from folder contents
goal: G-001
---

# CAP-005 — Index generation

## What

`clue scaffold [path]` regenerates the taxonomy README index blocks (`docs/README.md` and each `docs/<folder>/README.md`) from folder contents: hand-written entries whose targets survive keep their lines, missing entries are appended, entries whose targets are gone are dropped, and prose outside the `clue:index` markers is never touched. The command materializes nothing — missing folder READMEs are reported, never invented, and a path without a `docs/` tree is an error. `checkIndexes` (in `clue validate`) remains the judge of the result.

## Why

Delivers the maintenance half of [G-001](../../goals/G-001-verifiable-thread.md)'s navigable corpus: every Cliewen change's digest must leave the index blocks matching folder contents, and the validator only reports the mismatch — this command repairs it. The engine is the one `clue init` runs (shared per [ADR-019](../../decisions/ADR-019-init-regenerates-indexes.md)); this capability is its standalone, materialization-free exposure.

Acceptance criteria: [criteria.md](criteria.md) · engine and command surface: [design.md](design.md).
