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

Acceptance-criterion IDs escape the repository into tickets, commits, and reviews, so changing an ID's meaning silently corrupts old references.

## Decision outcome

- **An AC ID's meaning is immutable.** A meaning change retires the old ID and mints a new one; cosmetic wording may keep it, with the distinction judged in review.
- **Retire, don't delete.** A retired scenario remains as a tombstone tagged `@retired`, preserving what old references meant. The canonical grammar and aliases follow ADR-037.
- **Retired ACs have no tests.** A test referencing one fails validation, forcing cleanup; deletion is reserved for mistakes.
- **Duplicate declarations fail.** The criteria files are the registry and `clue` enforces corpus-wide uniqueness.

**Carrier:** `checkACTests` and the retire-don't-delete rule in `clue-delta`. A standalone registry is rejected because it would duplicate the corpus and become a hand-maintained index.
