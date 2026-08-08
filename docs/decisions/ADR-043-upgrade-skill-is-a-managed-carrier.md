---
id: ADR-043
type: decision
status: verified
links: [P-010, CAP-001, CAP-004, ADR-021, ADR-022, ADR-031, ADR-039]
title: The managed skill set includes a human-authorized upgrade entry point
author: agent
accepted-by: Flemming N. Larsen (2026-08-08, conversation)
---

# ADR-043 — The managed skill set includes a human-authorized upgrade entry point

## Context and problem statement

The five lifecycle skills tell an agent how to work inside a Cliewen repository but give an already-onboarded repository no named route for responding to a newer release. `clue latest` can identify the release and the machine-specific installation route, while `clue migrate` can move the reviewed repository state, but neither names the human decision that must connect them. A sixth managed skill changes the generated topology and must not weaken the ownership, migration, or plugin boundaries that were written for the former set.

## Decision outcome

**The managed Cliewen set contains six generated, repository-scoped skills: `clue-analysis`, `clue-delta`, `clue-extract`, `clue-plan`, `clue-upgrade`, and `clue-verify`.** `clue-upgrade` is the named entry point for an existing adopter to check a release, read its migration guidance, ask the human whether to act now or later, and carry an explicit yes through a normal reviewed repository change. It carries no installation command of its own, because `clue latest` is the only component that can select the route for the machine it is running on.

The sixth skill has the same ownership marker, version stamp, canonical source, generated agent tree, and embedded scaffold template as the other five. Its directory is a reserved managed slot, so an unmarked manifest there fails toward reinstall rather than disappearing from the versioned set.

MIG-003 may add the new canonical skill directory only when every carrier a supported preceding release shipped exactly matches that release's manifest and that release did not ship the new path. The preview names both roles: the target release supplies the bytes, and the recognized preceding release proves the addition is safe. A missing, partial, modified, or otherwise unrecognized set remains a finding and blocks every write. A Claude mirror remains optional and is never materialized by migration.

The Claude Code plugin remains a bootstrap only. It ships none of the six managed skills; `clue init` and the reviewed migration are still their only supported writers into the repository.

## Rejected: let any matching sibling authorize a new directory

The new directory has no sibling of its own, but one unrelated matching file is not evidence that the rest of the managed set is original. Treating it as sufficient could add a new carrier beside a partial or locally repaired installation and make an unresolved upgrade look complete.

## Rejected: let the upgrade skill install or merge on its own

The installation route is platform-specific and belongs to `clue latest`; a repository upgrade is a reviewed change and belongs behind the human's explicit choice and merge boundary. Combining either authority with a generic distributed skill would bypass the two boundaries this entry point is meant to make visible.

## Carrier

`internal/skills`, `internal/scaffold`, `internal/corpus`, `internal/migrate`, the generated skill trees, CAP-001 and CAP-004 criteria, the operations and plugin guidance, and the bootstrap plugin carry this decision.
