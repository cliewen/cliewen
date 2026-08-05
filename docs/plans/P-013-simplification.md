---
id: P-013
type: plan
status: draft
links: [P-012, G-001, PDR-013, PDR-025, PDR-026, AN-006, AN-008, AN-010, AN-012, AN-013]
title: Cliewen is simplified against a stated criterion
---

# P-013 — Cliewen is simplified against a stated criterion

This campaign is [P-012](P-012-migration-gap-closes-on-evidence.md)'s successor and stays `draft` until P-012 completes. [PDR-026](../decisions/PDR-026-campaigns-close-on-re-derived-evidence.md) moved simplification here from P-012; [PDR-025](../decisions/PDR-025-brownfield-migration-precedes-simplification.md)'s reason for deferring it in the first place still holds, and this file exists so that the deferral is a named destination rather than a phrase in four campaigns' out-of-scope sections.

Simplification has now been deferred by P-008, P-009, P-010, and P-011. That is worth stating plainly, because a thing deferred four times is either not actually wanted or has never been given a shape it could be executed in. This campaign's first problem is the second one.

## The material this campaign has to work with

[PDR-013](../decisions/PDR-013-explicit-core-red-line.md) already supplies the criterion — *does the core need it?* — and records why simplification kept stalling without one: every debate was argued case by case, protection had no boundary, and uniform protection turned out to be uniform friction. What the criterion has never been is systematically applied to a measured surface.

The measured cost is in [AN-006](../analysis/AN-006-plain-change-overhead.md) (overhead on changes that carry no meaning), [AN-010](../analysis/AN-010-adopter-change-overhead.md) (what a change costs an adopter rather than an authoring repository), and [AN-012](../analysis/AN-012-adopter-configuration-cost.md) (what adoption costs before any change happens). [AN-008](../analysis/AN-008-methodology-critiques.md)'s pattern C names the structural half: the graph only accumulates, and a corpus with no compaction path grows monotonically whatever its rules say.

[AN-013](../analysis/AN-013-distributed-work-and-evidence-boundaries.md)'s F1 and F2 are a different matter and are recorded here so they are not lost: a dependent change's base and authorization have no durable form, and cross-repository evidence has no way to be named without importing a foreign verdict. Both were routed to P-011, which closed without touching them. They add machinery rather than removing it, so this campaign either adopts them deliberately or re-routes them somewhere they fit — it may not drop them silently.

## Milestones

Unassigned. Milestone numbering continues corpus-global numbering from P-012 when this campaign opens; a `draft` plan does not hold milestone identities it has not yet decided, and deciding them is the work of the change that opens the campaign.

## Mutation rules

Status and evidence fields in the milestone table may mutate in an implementing change's merge digest. Everything else changes only through a declared plan revision backed by a decision record routed by reversal cost. Plan adjustments are decisions.
