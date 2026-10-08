---
id: AN-037
type: analysis
status: active
links: [CH-202, AN-036, AN-035, CAP-010, G-025]
title: Does the local-allocation rule tell a sole contributor from a second one
---

# AN-037 — Does the local-allocation rule tell a sole contributor from a second one

## Purpose

[CH-202](../../changes/CH-202-sole-contributor-local-allocation/proposal.md) changes the change loop so that, on the local-allocation warning, an agent reads the remote with plain Git and goes on when it finds no other contributor. Its challenge named the cheapest test and what would stop the work: the agent proceeding where another contributor's branch is on the remote, or still stopping in most runs without one. This records that test.

## Method

The `tracked-brief` scenario of [AN-036](AN-036-what-stops-the-first-tracked-change-in-a-fresh-repository.md), run from the CH-202 branch at `286c2b5` with the changed rule, on Claude Code and Codex at each agent's default effort. Variant `local-allocation`: five runs per agent. Variant `other-contributor`: three runs per agent, with another clone's `ch-150-colleague-change` branch pushed to the bare remote. Prediction stated before the runs: most `local-allocation` runs proceed and say so in the acceptance brief; all `other-contributor` runs stop.

## Fixture fault, found and corrected

The first round (`20261008-135215-tracked-brief-claude`, `20261008-140543-tracked-brief-codex`, `20261008-141819-tracked-brief-claude`, `20261008-142030-tracked-brief-codex`) was spoiled. Coordination, set up by the fixture and then removed by the variant, leaves a `clue/id-allocator` branch on the remote, and the rule tells the agent to treat claims there as another contributor's. Codex stopped on that branch in 3 of 5 `local-allocation` runs, citing "allocator claims". A repository that never coordinated has no such branch, so the cell did not test a sole contributor. Both variants now delete the remote allocator branch and fail if it is absent. All four cells were re-run; the first round is not counted below, and its `other-contributor` stops are not evidence either, since the allocator branch alone could have caused them.

## Observations

Re-run directories: `20261008-142409-tracked-brief-claude`, `20261008-143841-tracked-brief-codex`, `20261008-150357-tracked-brief-claude`, `20261008-150545-tracked-brief-codex` (Git-ignored).

| Cell | Runs | Proceeded to a pull request | Stopped |
|---|---|---|---|
| Claude Code, `local-allocation` | 5 | 4 (3 ready, 1 created) | 1 |
| Codex, `local-allocation` | 5 | 5 (3 ready, 2 created) | 0 |
| Claude Code, `other-contributor` | 3 | 0 | 3 |
| Codex, `other-contributor` | 3 | 0 | 3 |

**The decisive condition held.** In all six `other-contributor` runs the agent stopped before any proposal, reverted or left unpushed its local `CH-001` reservation, and named the colleague's branch as the reason (read from the last messages of all six). Codex asked whether to run `clue id coordinate --remote origin` or wait; Claude Code offered the same two answers. No run proceeded past the other branch, so the condition that would stop the work did not occur.

**The sole-contributor condition mostly held.** 9 of 10 `local-allocation` runs proceeded. The one that stopped (Claude Code run 3, `20261008-142409-tracked-brief-claude`) stopped at the route recommendation with three clarifying questions about the CSV design, not on the allocation rule; it never reached the identity step. AN-036 saw the same kind of stop once in five Claude Code runs on this scenario.

**Not read.** I did not read the nine completed acceptance briefs for the ledger note the rule requires (that allocation was local, what was read, and that `clue id coordinate` comes before a second contributor); a search showed `git ls-remote` in the transcript of every completed run, which shows the check was made, not that the brief says so. The reading is open.

## The rule's stated limit

The agent can see only pushed branches and claims. A second clone whose branch is unpushed, or a teammate who allocated an identity locally and has not pushed, is invisible to it, and this test did not and cannot show otherwise. The fixture's second contributor is visible by construction. The limit stays written into the skill text and the acceptance brief's ledger note.

## Decision

By the stated rules the work continues: no run proceeded past another contributor's branch, and the wording changed behaviour (AN-036: nine of ten stopped; here nine of ten proceeded, once the fixture was corrected). The cost of the fault was one re-run of four cells. Reading the briefs' ledger notes remains before this change is digested.
