---
id: CH-133
type: change
status: open
links: [P-013, M-064, PDR-029, PDR-013, PDR-026, AN-008, AN-018]
title: Score the remaining surface and register the remaining carrier prose
---

# CH-133 — Score the remaining surface and register the remaining carrier prose

## What this change does

[P-013](../../docs/plans/P-013-simplification.md)'s **M-064** carries four separable obligations, and this change performs the three that produce evidence and questions rather than edits. It scores the surface [PDR-013](../../docs/decisions/PDR-013-explicit-core-red-line.md)'s test governs — every remaining rule, artifact type, required field, command, and check — against *does the core need it?*; it registers the carrier prose M-062 and M-063 did not reach, statement by statement, against [PDR-029](../../docs/decisions/PDR-029-simplification-tests-by-surface.md)'s three conditions and its ordering rule; and it assembles the determination [AN-008](../../docs/analysis/AN-008-methodology-critiques.md)'s pattern C needs, so a human can close it, decline it with a stated cost, or route it to the successor campaign on evidence rather than on the sentence three campaigns have repeated.

The carriers are not edited. Nothing is removed, reworded, or reordered, and no trace or marker is added to any carrier — the same boundary [AN-018](../../docs/analysis/AN-018-skill-statement-register.md) held for the skills.

## Why the milestone splits here

M-062 registered and escalated; M-063 trimmed against the answers. That order was load-bearing, and PDR-029 states why: a statement cannot be defended or removed on merit until someone has written down what it is for, and an agent that resolves an untraceable statement silently has made a methodology decision without a human. M-064 reaches a wider surface than M-062 did and reaches the shared-memory surface that M-063 never touched, so the same split applies with more force. Every candidate this change scores is either carried out with a guard or declined with its cost stated, and declining is a decision — so the scoring produces open questions, and the carrying out belongs to the change that follows this one.

Two items in M-064 are therefore explicitly deferred to that follow-on change and named here so the deferral is not silent: `guide/change-loop.md` gaining [PDR-033](../../docs/decisions/PDR-033-planning-and-implementation-are-separate-steps.md)'s report-and-ask beside the pause it already describes, and every candidate this change's scoring approves.

## Scope

**Registered as carrier prose:** the public guide (`guide/*.md`), `CONTRIBUTING.md`, and the CLI's user-facing text. AN-018's segmentation rule is reused verbatim rather than restated, so the two registers are comparable and a second reader has one definition to disagree with.

**Scored against *does the core need it?*:** the live rule population the corpus states (constraints and the obligations decisions carry), the artifact types, the required frontmatter fields per type, the `clue` commands, and the checks `clue validate` runs.

**Pattern C's determination** is assembled, not decided. Closing it requires three things AN-008 names — a retired artifact provably leaving the mandatory read path as `clue context` and `clue validate` compute it, a supersession edge that answers what was downstream of a reversed decision, and a bounded rather than monotonic born-`inferred` population — and each is priced here so the human's answer has something to weigh. Because that surface is shared memory, the price includes what would become harder to recover.

## Out of scope

Any carrier edit. Any removal. [AN-013](../../docs/analysis/AN-013-distributed-work-and-evidence-boundaries.md)'s three findings, which are M-065's. The mechanisms M-062's answers asked for, which are M-067's. The re-derived cost evidence that closes the campaign, which is M-066's.
