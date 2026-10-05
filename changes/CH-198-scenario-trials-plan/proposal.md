---
id: CH-198
type: change
status: active
links: [G-025, P-024, G-017, AN-028]
title: Plan the repeatable scenario trials the method needs to be evidenced
---

# CH-198 — Plan the repeatable scenario trials the method needs to be evidenced

## Proposal

[AN-028](../../docs/analysis/AN-028-scenario-trials-on-different-agents.md) found that one headless run of the same task works on two agents, in a container, with no leaked user configuration, and that a repeatable, manually started trial harness is feasible. Nothing yet commits the repository to building one. This change records the commitment: a proposed goal, G-025, for re-observing what the method does on chosen agents and models on demand, and a draft plan, P-024, that delivers it in five milestones.

The digest is the two files, `docs/goals/G-025-*.md` and `docs/plans/P-024-*.md`, plus regenerated indexes. No code, no criteria, and no shipped surface change. Whether a milestone needs a capability criterion is settled when that milestone starts, not here.

## Challenge

**The assumption most likely to undermine the plan:** that repeated agent runs of one scenario are stable enough to compare, and informative enough to be worth the cost of running them. One run of the routing scenario on each of two agents ended in different states, one agent editing before asking and the other not, and both were correct under the method. If the spread is that wide on a simple scenario, a harness built around pass/fail would produce noise that looks like a result.

**A credible alternative:** keep observing by hand, as [AN-027](../../docs/analysis/AN-027-baseline-observation-first-period.md) did, with a written checklist per scenario and no harness at all. It costs nothing to build and needs no container, but it cannot be repeated after a skill changes and does not compare agents.

**The cheapest useful test:** M-103 is built to run it. Run one scenario five times on one agent in the container before any other scenario or adapter exists, and record the spread. If the outcomes agree on what the method decides, the plan proceeds as written. If they do not, M-105 and M-106 are revised before they start so that the harness reports a distribution and what a human should read, not a verdict.

**What would stop or revise the work:** a spread too wide to compare, or a harness whose maintenance costs more attention than the hand-run observation it replaces. In that case the plan is closed early with a finding, and the checklist alternative is what remains.

**How a criterion-satisfying build could still fail the maintainer:** scenarios that pass every deterministic check while the transcript shows an agent confused or lucky; skills tuned until the scenarios pass, which measures the scenarios and not the method; and a harness that works on the author's Windows machine and a Linux container and says nothing about any other setup. The plan answers the first with a transcript report that a human reads, the second by keeping scenario texts outside what the agent under test can read, and the third by stating the conditions in every recorded result.

## Scope boundary

This change adds a goal and a plan. It does not build the harness, choose the agent adapters beyond the three AN-028 examined, add anything to `clue`, or change a shipped skill. The harness is repository tooling and is not shipped to adopters, in line with the vision's out-of-scope line on depending on any vendor's agent.
