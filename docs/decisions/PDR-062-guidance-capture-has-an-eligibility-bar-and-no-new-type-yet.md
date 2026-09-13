---
id: PDR-062
type: decision
status: inferred
author: agent
links: [P-023, M-096, G-015, AN-024, ADR-026]
title: A reusable discovery earns capture only past an eligibility bar, and files into an existing home rather than a new corpus type
accepted-by: []
---

# PDR-062 — Guidance capture has an eligibility bar, an existing home, and no new corpus type yet

## Context

[G-015](../goals/G-015-what-the-work-teaches-is-retained.md) names a real gap: operational knowledge discovered mid-change — a command form that actually works, a recovery from a specific check failure — is neither product meaning, architecture, nor a future-shaping decision, so it has nowhere to go and is lost when the change ends. G-015 also names the opposite failure: a capture obligation with no eligibility bar produces a folder nobody reads in full, and writing down a workaround is how a workaround becomes permanent. [AN-024](../analysis/AN-024-methodology-review-acceptance-and-learning.md) proposed the mechanism and deferred the concrete choices rather than pre-deciding them. This decision settles those choices before [P-023](../plans/P-023-challenge-plans-and-retain-what-work-teaches.md)/M-097 needs a home to write its first discovery into.

## Decision

**Eligibility.** A discovery is capturable only if it cost something to find — a wrong first attempt, an undocumented flag, a nonobvious recovery — and is plausibly recurring: a fresh agent doing the same class of task would hit it again. Restating already-documented behaviour, or something observed once with no reason to expect recurrence, is not eligible. This is the low-value-accumulation guard.

**Home, and no new corpus type yet.** AN-024's own home table already lists five existing homes (analysis, guide, design, decision, script), and [ADR-026](ADR-026-adopter-types-default-lifecycle.md) already lets an adopter add an unrecognized type without touching the validator, so the mechanism to mint one is not the open question — whether one is needed *yet* is. A capability-specific discovery goes in that capability's own `design.md`. A discovery about running this repository's own tooling, not tied to one capability, goes in `CONTRIBUTING.md`, which already carries exactly this kind of procedural material (its "Verify Locally" block). Neither is a new document, and correcting either in preference to adding beside it is the duplication guard.

**Evidence sufficiency.** One observation in one session supports only a claim scoped to what was actually seen — the platform, the command form, the version — never a general claim from a single trial. Broadening the claim requires a later session confirming it under different conditions; until then, the guidance states its scope rather than a universal rule.

**Discovery.** No new index. Guidance placed in a capability's `design.md` is found the way all local design is found — `clue context <capability-id>` on the next relevant task. Guidance placed in `CONTRIBUTING.md` is found because [`AGENTS.md`](../../AGENTS.md)'s source-repository-conventions table already sends an agent there before verifying. Riding an existing read path meets G-015's "without reading every guide" bar without building anything new.

**Retirement.** Guidance is corrected or removed in the same change that discovers it is stale, the same way any other design-document content is kept current — no separate lifecycle and no tombstone, because it was never a standalone artifact with its own identity.

**The workaround guard.** Before documenting a workaround, ask whether the confusing step can instead be removed or automated. Only what cannot yet be fixed gets written down.

## Scope

This decision binds this repository only, ahead of [P-023](../plans/P-023-challenge-plans-and-retain-what-work-teaches.md)/M-100's trial-before-ship ordering; it carries no `binds: adopter`.

## Reopening

M-097's first real discovery and M-098's fresh-agent trial are the test of the home choice above. If the discovery fits neither a capability's `design.md` nor `CONTRIBUTING.md` without distorting it, or belongs to no single capability and is not about this repository's own tooling, that is the case a dedicated corpus type exists to solve, and this decision is revised rather than left beside a contradicting answer.

## Rejected: a new `runbook` (or similar) corpus type now

Minting a type under ADR-026's extension mechanism was the credible alternative: a dedicated folder and template would give every future discovery one obvious address regardless of scope. Rejected for now because no discovery has yet tested whether the two existing homes are actually too narrow — building a new home before that evidence exists is exactly the accumulation G-015 warns against, and the mechanism to add one later costs nothing to defer.
