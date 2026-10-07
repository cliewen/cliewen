---
id: CH-202
type: change
status: open
links: [CAP-010, AN-036, AN-035, G-025]
title: A sole contributor's first tracked change proceeds on local allocation
---

# CH-202 — A sole contributor's first tracked change proceeds on local allocation

## Proposal

This change is plan-less: P-024 is closed and no active plan carries this. It follows [AN-036](../../docs/analysis/AN-036-what-stops-the-first-tracked-change-in-a-fresh-repository.md), which ran ten first tracked changes in a repository with a seeded ledger and local allocation: nine stopped on the change loop's instruction to stop at the local-allocation warning and ask the maintainer, and one proceeded on its own judgement and said so in the acceptance brief. Both of the rule's branches, serialize on the integration branch or enable `clue id coordinate`, need a human to pick, so every first tracked change in a repository without coordination costs one maintainer answer before anything starts.

Change `clue-delta`'s change loop, step 1, so that on the local-allocation warning the agent looks for another contributor and, finding none, goes on. It checks the remote for any `ch-*` branch other than its own and any open pull request, and a `clue/id-allocator` branch carrying change claims. Where none exists and the host could be read, it serializes allocation on the integration branch: it continues with the reserved identity, commits and pushes the ledger with the proposal, and states under the acceptance brief's ledger note that allocation was local, what was checked, and that `clue id coordinate` comes before a second contributor. Where another contributor's branch exists, where the remote or host cannot be read, or where the check finds anything it cannot classify, the agent stops as today.

The maintainer chose this direction over the alternative of keeping the rule and having `clue init` say what to do, from AN-036's results. It records a decision (a PDR, inferred, agent-authored) because it changes what the method lets an agent do unasked, and it updates every live carrier of the rule: the canonical skill source and its generated and scaffolded copies, the guide's first-change page, CAP-010's criteria and design, and the changelog.

## Challenge

**The assumption most likely to undermine the work:** that an agent can tell from the remote that it is alone. A second clone with an unpushed branch, or a teammate who has allocated an identity locally and not yet pushed, is invisible to any check the agent can make. The rule would then turn a stop that asked a human into a silent proceed exactly where the stop mattered, and the two changes would receive the same CH number. The ledger's union merge and `clue validate` would catch the duplicate at merge time, but only after both authors had built on it.

**A credible alternative:** keep the stop, and make it cheaper: `clue init` or the first `clue id next` prints the two commands, and the agent offers to run `clue id coordinate` and carry on in one exchange. That keeps a human in the decision. Its cost is that the first change in every new repository still waits on an answer, which is the cost AN-036 measured; it was the other direction the maintainer weighed.

**The cheapest useful test:** run the changed rule against the same fixture AN-036 used, five runs each on Claude Code and Codex, expecting that most proceed and say so in the brief; and against a fixture variant in which a second contributor's `ch-*` branch is already pushed, three runs each, expecting that all stop. A variant with an unpushed second clone cannot be seen by the agent either, so it is not run; it is the stated limit of the rule, written into the skill text and the acceptance brief's ledger note, not something this test can show.

**What would revise or stop the work:** the agent proceeding in a run where another contributor's branch is on the remote, in which case the check is not reliable enough to replace a human answer and the rule is dropped in favour of the alternative above; or the agent still stopping in most runs without another contributor, in which case the wording does not change behaviour and the change is revised before anything ships.

**How an implementation could meet every criterion and still fail the maintainer:** the skill text says "look for another contributor" and the agent proceeds after a check that could not have found one (no remote read, no pull-request list), recording "no other contributor found" in the brief as though it had looked. The rule therefore names what must be read and says that failing to read it is a stop, and the ledger note states what was actually checked, not that nothing was found.

## Scope boundary

This changes the local-allocation branch of the change loop, one criterion's worth of generated guidance, and the carriers that repeat it. It does not change `clue id next`, `clue id coordinate`, the ledger or the coordination settings, and it does not change `clue init`: AN-036 did not test a missing ledger, so nothing here speaks to it. A repository with a second contributor is expected to enable coordination first, as the skill already says.
