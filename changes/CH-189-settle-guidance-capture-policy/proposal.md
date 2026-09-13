---
id: CH-189
type: change
status: open
links: [P-023, M-096, G-015, AN-024, ADR-026]
title: This repository settles what a reusable operational discovery is, where it lives, and how it is found and retired
---

# CH-189 — Settle guidance capture before any capture obligation exists

[G-015](../../docs/goals/G-015-what-the-work-teaches-is-retained.md) names a real gap: operational knowledge discovered mid-change — which command form actually works here, what the recovery is when a check rejects a candidate — is neither product meaning, architecture, nor a future-shaping decision, so the corpus has nowhere to put it and it is lost when the change ends. But G-015 also names the failure mode of answering that gap carelessly: a capture obligation with no eligibility bar produces a folder nobody reads in full, and writing down a workaround is how a workaround becomes permanent. [AN-024](../../docs/analysis/AN-024-methodology-review-acceptance-and-learning.md) deferred the concrete choices rather than pre-deciding them, and P-023/M-096 exists to settle them before M-097 needs a home to write into.

This change delivers M-096. It answers, by decision, the five choices AN-024 deferred:

**Eligibility.** A discovery is capturable only if it cost something to find — a wrong first attempt, an undocumented flag, a nonobvious recovery — and is plausibly recurring: a fresh agent doing the same class of task would hit it again. Restating already-documented behaviour, or something observed once with no reason to expect recurrence, is not eligible. This is the low-value-accumulation guard.

**Home, and no new corpus type yet.** AN-024's own home table already lists five existing homes (analysis, guide, design, decision, script), and [ADR-026](../../docs/decisions/ADR-026-adopter-types-default-lifecycle.md) already lets an adopter add an unrecognized type without touching the validator, so the mechanism to mint one is not the open question — whether one is *needed yet* is. This decision routes a capability-specific discovery to that capability's own `design.md`, and a repository-wide discovery about running this repository's own tooling to `CONTRIBUTING.md`, which already carries exactly this kind of procedural material (the "Verify Locally" block). Neither is a new document. A new type is deferred until M-097 or M-098 shows a discovery that fits neither home without distorting it.

**Evidence sufficiency.** One observation in one session supports only a claim scoped to what was actually seen — the platform, the command form, the version — never a general claim from a single trial. Broadening the claim requires a later session confirming it under different conditions; until then the guidance states its scope, not a universal rule.

**Discovery.** No new index. Guidance placed in a capability's `design.md` is found the way all local design is found — `clue context <capability-id>` on the next relevant task. Guidance placed in `CONTRIBUTING.md` is found because [`AGENTS.md`](../../AGENTS.md)'s source-repository-conventions table already sends an agent there before verifying. Riding an existing read path satisfies G-015's "without reading every guide" bar without building anything new.

**Retirement.** Guidance is corrected or removed in the same change that discovers it is stale, the same way any other design-document content is kept current — no separate lifecycle, no tombstone, because it was never a standalone artifact with its own identity.

**The guards, restated as guards.** Low-value accumulation is blocked by the eligibility bar above. Duplicating an existing home is blocked by requiring a check of the target home (the capability's `design.md`, or `CONTRIBUTING.md`) before writing, and correcting what is already there in preference to adding beside it. A workaround hardening into policy is blocked by asking, before documenting any workaround, whether the confusing step can instead be removed or automated — only what cannot yet be fixed gets written down.

## This proposal's own challenge

Starting M-096 with an approach not yet settled is exactly the case [PDR-061](../../docs/decisions/PDR-061-challenge-consequential-commitments-proportionally.md) calls consequential, so the challenge is stated here rather than skipped.

- **Riskiest assumption:** that no new corpus type is needed — that a capability's `design.md` and `CONTRIBUTING.md` are wide enough homes for what M-097 and M-098 will actually turn up.
- **Credible alternative:** mint a `runbook` type now, under ADR-026's existing extension mechanism, with its own folder and template, so every future discovery has one obvious address regardless of scope.
- **Cheapest useful test:** M-097 itself. Draft its first discovery against both homes before choosing: if forcing it into a capability's `design.md` or `CONTRIBUTING.md` loses the trigger/prerequisite/procedure/recovery shape G-015 requires, or if the discovery does not belong to any one capability and is not about this repository's own tooling, that is the case a new type exists to solve.
- **Stop or revise if:** M-097's real discovery does not fit either home without distortion, or M-098's fresh-agent trial cannot find the guidance from the entry points this decision names. Either result reopens this PDR rather than accumulating a second, contradictory answer beside it.
- **Satisfied but failing:** a discovery correctly filed in the right home but never actually linked from a path an agent reads before the task — the home question answered while the discovery question (findability) is quietly skipped. M-097's own exit criterion, not this decision, is what has to catch that.
