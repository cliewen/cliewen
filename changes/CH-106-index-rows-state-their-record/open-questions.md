---
id: CH-106-open-questions
type: change
status: open
links: [CH-106, ADR-035, ADR-039, C-013]
title: Open questions for CH-106
---

# CH-106 — Open questions

## Q1 — Does the judge *fail* on a divergent index row, or *count* it? (blocking)

**Why this blocks.** It decides the shape of everything else in the change: whether a corpus migration is required, whether adopters can upgrade without going red, whether the new constraint is born `machine` or `agent`, and whether M-049's own counter is affected. It also changes what `clue validate` fails on, which is a core carrier under C-013 and therefore explicitly not an agent's call.

**The stakes.** Two repositories are live on Cliewen — Tank Royale and model2diagram — and both carry index blocks written before any of this was stated. A hard failure turns them red the moment they upgrade the binary, for rows that were legal when written. Under the skill's own rule, a release that adds or narrows a corpus obligation has to ship a supported `clue migrate` migration; a failing check therefore drags a sixth migration, its preview/apply pair, and its evidence into this change.

**Recommendation: the judge counts, and does not fail.** `clue validate` reports divergent index rows as a counted population on its OK line, exactly as it already reports agent-enforced constraints and — until CH-105 emptied it — decisions awaiting verification.

The reasoning is that this repository has already run the experiment. The inferred-decision counter was a visible, non-blocking number for months; it did not rot, it drove a campaign milestone, and CH-105 drove it to zero. [ADR-035](../../docs/decisions/ADR-035-bounded-provenance-and-reality-edges.md) settled the general form of this — costly unverified meaning is reported as an actionable population rather than turned into a build failure — and [ADR-017](../../docs/decisions/ADR-017-conventions-are-constraints.md) already applies it to the constraint backlog with the words "the backlog is visible, not archival". A counted row is the established Cliewen answer for "true but not yet repaired everywhere", and it costs no adopter a red build for a label written before the rule existed.

It also keeps this change small enough to be reviewable: no migration, no adopter coordination, no release note about a breaking check. The generator stops producing violations, the counter shows the remainder, and if the count later proves it does not move on its own, promoting it to a failure is a second, cheap decision made with evidence — the direction that is easy to take later and hard to take back.

The cost of the recommendation, stated plainly: a rule that only counts is a rule a repository can ignore indefinitely, and C-004 forbids ever weakening it to make a number look better. That is the trade being accepted, and it is the same trade ADR-017 and ADR-035 already accepted twice.

## Q2 — Does M-049 own this, or is it plan-less? (non-blocking)

Recorded rather than assumed. This change does not move M-049's stated count (see the proposal's opening paragraph); it adds a machine-enforced constraint rather than converting one of the thirteen. It was scoped to M-049 by the human after CH-105 merged. Proceeding under that declaration with the qualification stated in the proposal, the plan evidence column, and the PR body. Correcting it is a one-line frontmatter change and needs no rework.
