---
id: CH-203
type: change
status: open
links: [CAP-010, PDR-066, AN-037]
title: Local allocation stops on any trace of another contributor
---

# CH-203 — Local allocation stops on any trace of another contributor

## Proposal

This change is plan-less: no active plan carries it. It follows two advisories from the isolated review of CH-202 ([PR 253](https://github.com/cliewen/cliewen/pull/253)), which merged with them open. Both concern what the agent's `git ls-remote --heads origin` can tell it.

The command lists branch names and tips, so "change claims on `clue/id-allocator`" cannot be read: the agent sees that the allocator branch exists and nothing of its claims. AN-037's spoiled first round showed the effect, Codex stopping in 3 of 5 runs on a leftover allocator branch. The rule now says what it already does in practice: the existence of `clue/id-allocator` on the remote is a stop, because a coordinated repository whose local settings are missing is not a sole contributor's repository.

The review's second advisory suggested ignoring `ch-*` branches already merged into the integration branch. This change does not take that remedy. A merged branch of someone else shows that someone else works here, and the remote does not say who wrote it, so ignoring it would turn a stop into a silent proceed exactly where the stop matters. The rule instead says that any `ch-*` branch other than the agent's own stops it, merged or not, and that the maintainer's answer, enabling `clue id coordinate`, ends the question for good. The cost is a repository that keeps merged branches stops once.

It amends PDR-066 in place (inferred, unaccepted), edits the canonical skill source and its generated and scaffolded copies, AC-218 and its tests, the changelog entry, and CAP-010's design sentence.

## Challenge

**The assumption most likely to undermine the work:** that stopping on a merged `ch-*` branch is a rare, one-time cost. A repository that never deletes merged branches, and where the only contributor's own earlier changes remain on the remote, would stop on every change until coordination is enabled; the agent then asks a question whose answer is always the same.

**A credible alternative:** the review's remedy, excluding merged branches, or comparing the tip author with `git config user.email`. Both let a colleague's merged or same-address branch through, and the second also fails where agents commit under the human's address.

**The cheapest useful test:** none beyond the wording tests. The change only narrows what the agent may do unasked, so the failure direction is a false stop, which AN-037 already measured and found recoverable with one answer; a run would measure the same thing again.

**What would revise the work:** evidence that adopters keep merged branches and hit the stop repeatedly, in which case the answer is a better stop message or an author comparison, decided then.

**How an implementation could meet every criterion and still fail the maintainer:** the skill text stops on branches but the stop message does not tell the human that enabling coordination ends the question, so a sole contributor is asked the same thing each change. The stop already names `clue id coordinate`; the test checks it stays.

## Scope boundary

Only the local-allocation branch of the change loop and the carriers that repeat it. No change to `clue id next`, `clue id coordinate`, the ledger, or `clue init`.
