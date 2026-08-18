---
id: ADR-034
type: decision
status: verified
links: [ADR-025, PDR-003, C-008, CAP-002, P-007]
title: Retirement is deletion; supersedes carries the pointer forward
author: agent
accepted-by: Flemming N. Larsen (2026-08-02, conversation)
---

# ADR-034 — Retirement is deletion; supersedes carries the pointer forward

> **Amended by [PDR-046](PDR-046-decisions-route-by-subject.md):** decision records are retained and no longer demote into a log. Other retired artifacts still delete their files, with the named exceptions below.

## Context and problem statement

The default lifecycle called `retired` a terminal status even though retired artifacts were deleted in practice, leaving no durable pointer from a deleted ID to the reader's next stop. Criteria tombstones and completed plans are different because they deliberately remain on disk.

## Decision outcome

**Retiring an artifact deletes its file in the same change; `status: retired` is not a surviving default-lifecycle state.** A `supersedes:` frontmatter field is optional on any artifact and is carried by the direct successor or best live next stop, while Git history remains the archive of the deleted text. `clue validate` rejects a `supersedes:` target that still resolves to a live artifact.

Criteria retain their `@retired` tombstones so stale test tags fail loudly, and completed plans remain frozen under C-008 rather than being retired. Decision records are retained under PDR-046, which retires the former decision-log demotion instead of deleting those records. ADR-025's default lifecycle and goals are corrected accordingly.

**Carrier:** the `supersedes:` validator rule, lifecycle and status descriptions, and the successor artifact that carries each retirement pointer.
