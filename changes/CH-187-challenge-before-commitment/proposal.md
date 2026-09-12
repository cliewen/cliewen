---
id: CH-187
type: change
status: open
links: [P-023, M-094, G-014, AN-024, VIS-001, ADR-062, PDR-056, C-013]
title: This repository challenges a consequential commitment before making it, and says in the same rule when not to
---

# CH-187 — Challenge a consequential commitment before making it

[G-014](../../docs/goals/G-014-riskiest-assumption-is-challenged.md) names a gap the method has no answer for: a criterion, its implementation, and its test are usually written by one author in one sitting, so a wrong assumption travels through all three and comes out looking verified. Nothing currently asks, before the repository commits to a course of work, what would have to be true for the work to be worth doing and what the cheapest way to find out is. [PDR-056](../../docs/decisions/PDR-056-plan-health-checks-and-lightweight-replanning.md)'s plan-health check comes closest, but it asks whether a plan *still* holds once work starts; it never asks whether the plan should have been adopted.

This change delivers P-023/M-094. It states the challenge obligation for this repository and records the decision that shapes it.

**The obligation.** Before a consequential commitment, the proposer names four things where the reviewer will read them: the assumption most likely to undermine the work, a credible alternative to the chosen course, the cheapest useful test of that assumption, and the result that would stop or revise the work. It also asks what an implementation could look like that satisfied every criterion and still failed the person the work is for. A milestone whose deliverable is decision-changing evidence rather than a feature is a legitimate way to spend a milestone.

**The limit, in the same rule.** The challenge scales with consequence and uncertainty. Simple work, and full work that carries out a course already challenged and well understood, proceeds without it. No prototype is ever universally required. The limit is part of the rule rather than an exemption from it, because a challenge demanded of everything becomes a form filled in by habit, and a form filled in by habit challenges nothing.

**Where it lives.** One place: a section of `AGENTS.md`, the file every agent session in this repository reads before planning anything, linked to a new decision record, PDR-061, which fixes what counts as consequential and why proportionality belongs inside the rule. The plan binds this milestone to the repository only. Nothing is written onto `internal/skills/source/` or the scaffolded templates: an obligation reaches adopters only through M-100, after M-095 has tried it on real work in both directions. PDR-061 therefore carries no `binds: adopter`, and [ADR-062](../../docs/decisions/ADR-062-repository-role-is-declared-machine-state.md)'s carrier check correctly does not apply to it.

**No acceptance criterion.** The rule changes no behaviour of `clue`, no generated skill, and nothing an adopter receives, so there is no executable surface to prove, and a `Human`-class criterion stating "the rule is followed" would only restate M-095 before M-095 has run. M-095 is where the rule earns evidence; this change states it.

## This proposal's own challenge

The rule applies to the change introducing it, so it is stated here rather than waived.

- **Riskiest assumption:** that writing an obligation where agents read it makes them meet it. [AN-025](../../docs/analysis/AN-025-unenforced-integration-boundary.md) is direct evidence against that: this repository's strictest integration rule was stated in four shipped documents and unmet for the project's whole life.
- **Credible alternative:** put the four questions into the pull-request template or `proposal.md` skeleton, so the obligation is a visible blank rather than a sentence to remember. Rejected for now because a blank invites a filled-in form on work that needed none, which is the failure the proportionality limit exists to prevent — but it is the obvious next move if the trial shows the sentence being skipped.
- **Cheapest useful test:** M-095 itself, which requires one real proposal redirected by the challenge and one real proposal that proceeds because the rule said it did not need investigating. No separate experiment is needed.
- **Stop or revise if:** M-095 cannot find a redirected proposal among real campaign work, or finds the challenge written as boilerplate on proposals that did not need it. Either result means the rule should be revised or not shipped, and M-100 records which.
- **Satisfied but failing:** an agent that writes four plausible bullets under every proposal meets the letter of this change and defeats its purpose. That is why the limit sits in the rule and why M-095's evidence must include a case where no challenge was written.
