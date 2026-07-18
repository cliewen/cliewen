---
id: CAP-005-design
type: design
status: active
links: [CAP-005, ADR-019]
title: Design for CAP-005 index generation
---

# Design — CAP-005 Index generation

One engine, two exposures ([ADR-019](../../decisions/ADR-019-init-regenerates-indexes.md)): the scaffold package's index regeneration runs inside `clue init` after materialization, and stands alone as `clue scaffold` via an exported entry point (`Regen`) that materializes nothing.

- **The regeneration contract mirrors the judge**: the set of READMEs regenerated (docs/README.md and each first-level folder README) is exactly the set `checkIndexes` validates, and the wanted-entry rule (sibling artifacts plus markdown-bearing subfolders) is the same — including subfolder coverage, where any live link into the subfolder counts, so a curated descendant entry survives regeneration exactly as it survives validation. A validator rule the regeneration violates fails the AC-026 test the moment it lands.
- **Prose safety**: only the block between the `clue:index` markers is rewritten; the file's own line endings are preserved; a marker-less taxonomy README gains an appended block; a lone or reversed marker is an error naming the file — never a guess.
- **Nothing is invented**: missing folder READMEs are reported (recursively, matching the validator's reading) and left missing; a root without a `docs/` tree exits 1 with the error on stderr — `clue init` is the tool that materializes, and a silent no-op would mask a mistyped path.
- **Report surface**: `indexed` lines per regenerated README, `missing` lines per absent folder README, one summary line — same vocabulary as `clue init`.
