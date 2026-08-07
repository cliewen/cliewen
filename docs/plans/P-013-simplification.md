---
id: P-013
type: plan
status: draft
links: [P-012, G-001, ADR-040, PDR-013, PDR-025, PDR-026, AN-006, AN-008, AN-010, AN-012, AN-013]
title: Cliewen is simplified against a stated criterion
---

# P-013 — Cliewen is simplified against a stated criterion

This campaign is [P-012](P-012-migration-gap-closes-on-evidence.md)'s successor. P-012 has completed; this campaign stays `draft` until it is opened, which is what the milestone section below records. [PDR-026](../decisions/PDR-026-campaigns-close-on-re-derived-evidence.md) moved simplification here from P-012; [PDR-025](../decisions/PDR-025-brownfield-migration-precedes-simplification.md)'s reason for deferring it in the first place still holds, and this file exists so that the deferral is a named destination rather than a phrase in successive campaigns' out-of-scope sections.

Simplification has been deferred by every campaign that named it, each one to its own successor, three times over. That is worth stating plainly, because a thing deferred three times is either not actually wanted or has never been given a shape it could be executed in. This campaign's first problem is the second one.

## The material this campaign has to work with

[PDR-013](../decisions/PDR-013-explicit-core-red-line.md) already supplies the criterion — *does the core need it?* — and records why simplification kept stalling without one: every debate was argued case by case, protection had no boundary, and uniform protection turned out to be uniform friction. What the criterion has never been is systematically applied to a measured surface.

The measured cost is in [AN-006](../analysis/AN-006-plain-change-overhead.md) (overhead on changes that carry no meaning), [AN-010](../analysis/AN-010-adopter-change-overhead.md) (what a change costs an adopter rather than an authoring repository), and [AN-012](../analysis/AN-012-adopter-configuration-cost.md) (what adoption costs before any change happens). [AN-008](../analysis/AN-008-methodology-critiques.md)'s pattern C names the structural half: the graph only accumulates, and a corpus with no compaction path grows monotonically whatever its rules say.

[AN-013](../analysis/AN-013-distributed-work-and-evidence-boundaries.md)'s findings are a different matter and are recorded here so they are not lost. P-010's M-050 evidence names F1, F2's first half, and F3 as open findings feeding P-011, and P-011 closed without addressing any of them — but their histories differ, and this campaign owes each one a determination rather than a repetition of that sentence. P-009 declined F1's dependent-change candidate outright rather than deferring it, so what is open there is the decline's stated cost, not an unbuilt interface. [ADR-040](../decisions/ADR-040-qualified-external-references.md) consumed two of AN-013's candidates at once under P-009's M-044, so both remaining findings are narrower than they first read: it answered F3's qualified-reference candidate, leaving as open only whatever that decision did not reach, and it supplied the named-but-locally-unproven form F2 asked for, which leaves as F2's residue what AN-013 actually scopes to its first half — one repository's emitted CI wall kept in step by hand. That residue has never had a route at all. This campaign establishes for each whether it is answered, declined with a recorded cost, or genuinely open, and adopts or re-routes what remains; each of them adds machinery rather than removing it, which is exactly why the determination belongs here rather than in a milestone that assumes the answer. None of the three is dropped silently.

## Milestones

Unassigned. Milestone numbering continues corpus-global numbering from P-012 when this campaign opens; a `draft` plan does not hold milestone identities it has not yet decided, and deciding them is the work of the change that opens the campaign.

## Mutation rules

Once this campaign opens and its milestone table exists, status and evidence fields in that table may mutate in an implementing change's merge digest. Everything else changes only through a declared plan revision backed by a decision record routed by reversal cost. Plan adjustments are decisions.
