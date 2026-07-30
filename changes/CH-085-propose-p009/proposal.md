---
id: CH-085
type: change
status: open
links: []
title: Open P-009 to route P-008's two closing candidate sets
---

# CH-085 — Open P-009 to route P-008's two closing candidate sets

## What

Create [P-009](../../docs/plans/P-009-adopter-and-distributed-boundaries.md) `active`, a sequential campaign whose four milestones are the load-bearing candidates [AN-012](../../docs/analysis/AN-012-adopter-configuration-cost.md) and [AN-013](../../docs/analysis/AN-013-distributed-work-and-evidence-boundaries.md) each named for a successor plan when P-008 closed. This change writes the plan, its index row, and one decision-log row recording the campaign's scope and order. It implements none of the four milestones and is plan-less because its product is the plan.

## Why

P-008 closed 2026-07-30 with both analyses ending in findings rather than interfaces: M-035 measured that a real adopter's actual needs were a wall it had to fork and an upgrade path it could not follow, not a configuration file ADR-013 had already rejected; M-036 reproduced that no artifact records a dependent change's base or authorization, that cross-repository evidence is discharged by hand, and that no foreign identifier can be expressed at all. Both analyses named P-009 as consumer and explicitly could not route their own candidates — M-035 and M-036 forbid implementing an interface. Leaving the candidates unrouted lets fresh, pinned evidence age past the point a plan can still speak to what was actually measured.

## Decision boundary

This change decides only the campaign scope and order. It changes no corpus format, validator behavior, generated skill, or public command, and drafts no ADR or PDR content. Each milestone is a separate full change that routes its own decision record: M-037 and M-038 are process and corpus-format decisions with no adopter-facing surface; M-039 and M-040 become public contracts the moment an adopted repository depends on them, so they are ordered after the format decisions they build on.
