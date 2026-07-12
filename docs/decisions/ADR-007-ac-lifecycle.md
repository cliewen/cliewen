---
id: ADR-007
type: decision
status: verified
links: [CAP-002, ADR-005]
title: AC lifecycle — meaning-immutable IDs, retirement by tombstone
author: human
accepted-by: Flemming N. Larsen (2026-07-12)
---

# ADR-007 — AC lifecycle: meaning-immutable IDs, retirement by tombstone

## Context and problem statement

Requirements change; acceptance criteria die or change meaning. But AC IDs leak into places the repo cannot fix — tickets, commit messages, review threads — so an ID whose meaning changes silently poisons every old reference (§4: external systems reference IDs; IDs are eternal identity).

## Decision outcome

- **An AC ID's meaning is immutable.** A change that alters what the criterion *means* retires the old ID and mints a new one. Cosmetic rewording that preserves meaning may keep the ID — whether meaning changed is exactly what PR review judges (machines enforce form, humans verify meaning).
- **Retire, don't delete.** A dead AC stays in `criteria.md` as a tombstone: its scenario keeps its tag line with `@retired` added (`@AC-012 @retired`). The tombstone is what makes ID-reuse physically visible without git archaeology, and it preserves what old references meant.
- **A retired AC's tests must die.** Retired ACs require no tests, and a test still referencing one fails the build — the linter forces the cleanup rather than requesting it. (Deleting the AC outright also breaks the build via the unknown-reference rule, but loses the tombstone; deletion is for mistakes, retirement is for requirement changes.)
- **Duplicate AC declarations fail.** The criteria files are the AC registry — `clue` builds it by scanning — and a registry's one hard guarantee is uniqueness.

**Carrier:** the retirement and uniqueness rules in `clue`'s `checkACTests` (machine); the retire-don't-delete convention in the `clue-delta` skill (agent).

### Rejected: a standalone AC registry

Unnecessary — the AC↔test contract means an unimplemented AC fails the build, so there is no phantom-AC backlog to track, and a registry file would be the hand-maintained index that rots first. The corpus is the registry.
