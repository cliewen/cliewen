---
id: AN-027
type: analysis
status: active
links: [G-017, P-023, AN-024, AN-025, AN-026]
title: Baseline observation, first period — this campaign's own work
---

# AN-027 — Baseline observation, first period — this campaign's own work

## Purpose

[P-023](../plans/P-023-challenge-plans-and-retain-what-work-teaches.md)/M-099 requires what is observed to be fixed before the observing, and requires the first period's record to come from this campaign's own work rather than from a wider trial that has not happened yet. This document does both: it states the population, conditions, and sampling method before naming a single observation, then records everything in that population as of 2026-09-13. It is the first period, not the only one — [G-017](../goals/G-017-the-method-is-evidenced-and-simplified.md) continues past this campaign, and later periods will draw on wider populations, including willing adopters if any take part.

## Methodology, fixed before observing

**Population.** Every milestone in [P-023](../plans/P-023-challenge-plans-and-retain-what-work-teaches.md) that had reached `done` at the time this document was written: M-094, M-095, M-096, M-097, M-098, M-101, M-102. M-099 (this milestone) and M-100 (not yet started) are excluded because a milestone cannot be evidence for the period that produces it. No milestone is excluded for producing an unfavorable result; M-098's disclosed miss (below) is in the population precisely because exclusion would defeat the point.

**Conditions.** One repository (`cliewen`), one maintainer (Flemming N. Larsen), a mix of Claude Code sessions run directly by the maintainer and `general-purpose` subagents run in isolated git worktrees, on Windows with PowerShell as the default shell and Bash available. No adopter repository and no second organization is represented in this period. Every observation below names which of these conditions applied; none should be read as holding under different tooling, a different maintainer, or a different repository's history.

**Sampling.** This is not a sample. The population is small enough (seven milestones from one campaign) that every member is examined; nothing here supports a claim about milestones or campaigns not in the list. A future period drawing on more campaigns, or on adopters, will need to say explicitly whether it examines all of them or a defined subset, and why.

**What counts as evidence, and what does not.** An observation must be traceable to a specific milestone's evidence field, decision record, or linked pull request — not to a general impression of how the campaign went. Artifact counts (how many decision records, how many milestones) and green `clue validate` runs are explicitly not evidence of anything G-017 asks about; none appear as findings below, only as context for what a finding refers to.

## Observations

### Bypassed obligation, with the reason

[AN-025](AN-025-unenforced-integration-boundary.md): this repository's default branch had no branch protection of any kind for its entire history, while [PDR-021](../decisions/PDR-021-supported-merge-commit-history.md) — accepted and `status: verified` since 2026-08-02 — describes a protected branch that never existed, and four shipped documents tell every adopter that a host unable to enforce the same boundary is outside the supported adoption path. The reason was not disagreement or cost: nobody knew, because branch protection lives at the Git host and appears in no file a validator or a `git log` can see. This is the single clearest instance in the period of G-017's most informative category — a bypassed rule that reveals what enforcement the method actually has, as opposed to what it claims.

### Consequential misunderstanding, disclosed rather than hidden

M-098's first fresh-agent trial run: the subagent completed a real task ([PR #220](https://github.com/cliewen/cliewen/pull/220)) correctly, but sourced its safe commit- and PR-composition choice from the operator's own cross-session memory index rather than from anything in the repository, and said so candidly under direct questioning. The run is recorded as a miss rather than a pass because a memory index that already carries a lesson can pre-empt the very discovery the milestone exists to test — an interaction between two of Cliewen's own persistence mechanisms (the corpus and an operator's private memory) that no milestone had previously named as a risk. This is now visible for future guidance trials to control for; it was not visible before M-098 ran.

### Rediscovery, both the kind that worked and the kind that did not need to

M-097/M-098 together: the PowerShell backtick-escape corruption that produced [PR #49](https://github.com/cliewen/cliewen/pull/49)'s original description was rediscovered once, filed into `CONTRIBUTING.md` by CH-190, and then found again by a third fresh agent exactly where it was placed, with no separate pointer given, applying the documented procedure and producing an uncorrupted commit and PR body ([PR #221](https://github.com/cliewen/cliewen/pull/221)). A second trial run in the same milestone shows the other side of rediscovery: given the same task class after M-101 had already merged, the agent correctly recognized the work was already done and stopped, rather than rediscovering a solution to a problem that no longer existed. Both outcomes came from the same guidance and the same entry point; only the state of the repository differed.

### Misleading documentation, corrected

[M-101](../plans/P-023-challenge-plans-and-retain-what-work-teaches.md): `guide/methodology.md` drew a persistent Change → Capability edge that no artifact in this repository's corpus ever encodes — every capability's frontmatter names a Goal, none names a Change — while `guide/intent.md` already drew Capability fed only from Goal and optional Use case. [AN-026](AN-026-trialling-the-challenge-rule.md)'s Trial 2 records that the first analysis pass read this as a possible core-meaning conflict requiring a decision record, and only a second pass, checked against the corpus's actual field-level structure, found it to be an editorial gap in `methodology.md`'s diagram. The correction ([PR #220](https://github.com/cliewen/cliewen/pull/220)) merged as simple work. The misleading document and the initial misreading of how serious the correction was are both part of this observation: the guide misled a reader, and the same gap nearly misled the analysis meant to fix it.

### Acceptance effort

M-098's three trial runs each ended in a real pull request, and Flemming N. Larsen personally reviewed and merged each one — the human-acceptance act the review boundary reserves — including questioning each agent directly about where its choices came from rather than accepting a green run at face value. The effort this took is not quantified here (a raw time figure would misstate what G-017 asks for, since nothing establishes what effort would have been "too much"), but the record establishes that acceptance in this period involved direct questioning of the agent, not only reading the diff.

### Challenge rule, exercised in both directions

[AN-026](AN-026-trialling-the-challenge-rule.md) records M-095's two required trials on real, previously unsettled campaign decisions: one redirected (M-095's own delivery approach, away from an invented illustrative example and onto real unsettled decisions already `todo`), and one that proceeded to a conclusion without a decision record after the cheapest useful test was actually run (M-101's routing) — though Trial 2's own history shows the first pass of that test reached a wrong, too-comfortable conclusion, caught only by an independent review before it reached M-101. The rule caught something real in one direction and, on the other, needed a second pass to avoid a wrong answer of its own.

## What this period does not establish

Seven milestones from one campaign, one repository, and one maintainer cannot support a claim about whether any specific obligation should be kept, revised, or removed — that judgment needs more periods, and ideally an adopter's population, before G-017's evaluation is asked to bear weight. This period also does not exercise willing adopters, a different Git host, or a different maintainer, so nothing here should be read as evidence about conditions this repository does not have. What it does establish is that the categories G-017 named — consequential misunderstandings, acceptance effort, rediscovery, misleading documentation, and bypassed obligations with reasons — each had at least one real instance in this campaign's first period, none of them manufactured for the purpose, which is the bar M-099 sets.

## Next period

The next period should include M-100's adopter-facing carrier work once it exists, and should look for the same five categories again rather than assume the list is complete. If a willing adopter becomes available before M-100 closes, extending the population to include them is worth more than closing this analysis early.
