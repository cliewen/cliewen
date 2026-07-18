---
id: ADR-022
type: decision
status: verified
links: [G-002, CAP-004, ADR-011, ADR-018]
title: Cliewen skills declare ownership in frontmatter
author: agent
accepted-by: Flemming N. Larsen (2026-07-18, planning conversation — explicitly approved the marker and non-deletion boundary)
---

# ADR-022 — Cliewen skills declare ownership in frontmatter

## Context and problem statement

[ADR-011](ADR-011-version-stamping.md) makes the skills a versioned set under `.agents/skills/`, but a repository may install unrelated skills in the same directory. A validator that enrolls every `skill.md` by location alone forces third-party skills to carry Cliewen's version and gives a future installer no machine-readable ownership boundary. How does Cliewen identify the skill files it owns without deleting or absorbing neighboring skills?

## Decision outcome

**Every distributed Cliewen skill declares `cliewen-skill: true` in YAML frontmatter; only marked skills participate in Cliewen's version set, while initialization remains non-destructive.**

- **Ownership is explicit.** A `skill.md` with the boolean marker `cliewen-skill: true` is a Cliewen-managed skill. The marker travels with a standalone skill just as its `version:` does. An absent marker means the skill is unmanaged and Cliewen version checks ignore it; a present marker whose value is not boolean `true` fails as malformed rather than silently escaping validation.
- **Legacy Cliewen names fail toward migration.** The five canonical directories (`clue-analysis`, `clue-delta`, `clue-extract`, `clue-plan`, and `clue-verify`) are reserved as legacy Cliewen slots. If one contains an unmarked `skill.md`, validation reports that the pre-marker skill must be reinstalled. A new binary therefore cannot mistake an old Cliewen installation for a repository with no managed skills.
- **Version rules apply after ownership.** Every marked skill must carry a string `version:`, all marked skills must agree, and a released binary must match them. Unmarked third-party skills do not affect those checks.
- **Initialization does not become an updater.** `clue init` continues to create missing files and skip existing files. It neither overwrites nor deletes skills. A future explicit upgrade operation may replace files inside marked Cliewen skill directories, but removal is needed only when a later decision renames, moves, replaces, or retires a managed skill.

**Carrier:** generated skill frontmatter and `corpus.checkSkillVersions` (machine); `clue init`'s existing never-overwrite behavior (machine).

### Rejected: directory location alone establishes ownership

`.agents/skills/` is a shared integration point, not a Cliewen-exclusive directory. Treating every child as Cliewen makes unrelated skills fail missing-version and drift checks.

### Rejected: a separate Cliewen skill manifest

A manifest can enumerate ownership centrally, but a copied standalone skill loses that evidence. Per-skill frontmatter keeps ownership and version on the artifact they describe.

### Rejected: delete and reinstall the skill tree during initialization

Deleting a shared tree risks destroying unrelated skills, while replacing existing files would break `clue init`'s idempotent, never-overwrite contract. No deletion is required while the managed skill topology remains stable.
