---
id: AN-036
type: analysis
status: active
links: [P-024, G-025, G-017, AN-034, AN-035, CAP-001, CAP-010]
title: What stops the first tracked change in a fresh repository, fixed before observing
---

# AN-036 — What stops the first tracked change in a fresh repository, fixed before observing

## Purpose

[AN-035](AN-035-first-runs-of-the-two-long-scenarios.md) found that in a repository made by `clue init`, agents stopped before a first tracked change could start: Codex in 3 of 5 runs, with reasons about missing coordination or a missing hosted forge, and Claude Code in 1 of 5 for clarifying questions. I recorded it as "`clue init` leaves no identity ledger". Reading the code and the transcripts again, that is only part of it. `clue id next` already tells an agent to run `clue migrate --apply` when the ledger is missing, and Claude Code did. The Codex stops cite the shipped change loop's rule to stop and have the maintainer serialize allocation or run `clue id coordinate` when allocation is local. The fix that followed set the ledger and the coordination together, so the AN-035 runs do not say which one mattered.

This analysis separates them, so that a change to `clue init`, to the skill's rule, or to neither is chosen on evidence. It is committed with the protocol alone; results, readings and the decision are added later under their own headings, in a later commit.

## Methodology, fixed before observing

**Question.** In a repository with a seeded identity ledger and local allocation, a remote and no coordination, does an agent asked to take a change through the change process start it, or stop? And if it stops, on what reason?

**Population.** The `tracked-brief` scenario with the new variant `local-allocation`, which removes `.clue/id-coordination.yaml` from the fixture after setup and leaves the seeded ledger, so allocation is local. Claude Code (2.1.289, `claude-sonnet-5-5`) and Codex (`codex-cli` 0.160.0, `gpt-6.1-sol`), five runs each, ten in all, built from one `clue` commit, the `main` commit at the merge of the pull request that carries this protocol, which the results name. OpenCode is not included: the question is what the method's own rule does to the two agents the criterion names, and OpenCode's `tracked-brief` cell already has an uncorrected fault.

**Conditions.** As AN-035: one maintainer's Windows 11 machine, Docker Desktop with Linux containers, subscription logins, an empty agent home in every run, reasoning effort at each agent's default, recorded in `conditions.json`. The `gh` stand-in now answers `gh ruleset check --default`, the branch-protection endpoint and the repository's merge settings as a protected repository would (a required `validate` check, pull requests required, force-push and deletion blocked, merge commits allowed, squash and rebase disabled). That is a stipulation about the fixture and not an observation about any real host; its purpose is that a pull request can reach `ready`, which AN-035 could not read. The comparison runs are therefore not comparable on the `pr` field with the AN-035 cells, whose stand-in answered `{}`.

**What each run is read for.** The checks score every run, with the fields of AN-035. A person reads, in full, every run whose signature differs from the most common one in its cell, and two further runs per cell, the two whose run numbers sort first by SHA-256 of `109-<invocation directory name>-<run number>`. For any run that stops before the proposal is pushed, the stop reason is quoted from the agent's last message and coded as one of: the identity ledger, local allocation, the hosted forge, branch protection, a clarifying question about the request, or another reason. This time I read the acceptance brief of each run that reaches a ready pull request, which AN-035 did not.

**Expected result, stated before the runs.** Claude Code starts and completes the change in most runs, as it did with both a ledger and coordination. Codex stops on the local-allocation rule in most runs, quoting `clue id coordinate` or serialized allocation, and so does not complete in them. That would say the ledger alone is cheap and the coordination rule is the friction. If Codex completes in most runs, the AN-035 stops were about something else, and the fixture's coordination was not what mattered.

**Decision rules.**

- If a majority of runs in a cell complete without stopping on local allocation, the rule does not stop that agent here, and no change to it is proposed for that agent.
- If a majority of Codex runs stop on local allocation and the stop reason is the rule's wording, the finding is recorded as a question for the method, with the two directions named before the run: let the agent proceed on local allocation and state it in the acceptance brief, or keep the rule and have `clue init` say what to do. Neither is chosen here, and either goes to the tracked route as its own change, whose `proposal.md` states the riskiest assumption.
- If any run stops on the missing ledger, which this variant seeds, that is a fixture fault, counted, corrected and the cell re-run.
- A change to `clue init` is proposed only if a stop is attributed to a missing ledger in a repository where `clue init` produced it. Nothing in this design runs `clue init` without a ledger, so this analysis cannot propose it; the AN-035 first-round runs are the only evidence for it and they do not separate the causes.

**What would make the runs uninformative.** More than two runs in a cell stopping on a clarifying question about the request rather than on a rule (AN-035 saw this once in five on Claude Code); a harness fault that spoils a cell, which is re-run and recorded; checks that disagree with the transcripts read in more than one run in ten of those read.

## Observations

Run on 2026-10-07 on `main` at `af11cd4`, with `clue` and skills 0.28.0, Claude Code 2.1.289 and Codex 0.160.0 as pinned, in `20261007-180052-tracked-brief-claude` and `20261007-180725-tracked-brief-codex` under `scenario-runs/` (Git-ignored), variant `local-allocation`, efforts at each agent's default. Cost: Claude Code about $0.1 for each run that stopped and about $0.7 for the one that completed; Codex reports none. All 10 runs were scored by the checks.

| Cell | Runs | Result |
|---|---|---|
| `tracked-brief` with local allocation, Claude Code | 5 | 1 completed: proposal before source, pushed, no workspace left, `validate` green, no template text, pull request marked ready. 4 stopped before the proposal, each after `clue id next CH` printed the local-allocation warning |
| `tracked-brief` with local allocation, Codex | 5 | 0 completed. 5 stopped on the same warning |

**Nine of ten runs stopped on the local-allocation rule, and the stop reason was the rule's wording.** Every last message of the nine quotes the change loop's instruction to stop and either serialize allocation on the integration branch or have a maintainer enable Git coordination, and asks the maintainer which to do. Codex asked to run `clue id coordinate --remote origin` in three runs and asked "coordinate or serialize" in two; Claude Code offered both options where it asked, and named serialize as its suggestion, for a sole contributor, in one. None stopped on the ledger, which the variant seeds, on the hosted forge, on branch protection or on a clarifying question about the request.

**The one run that completed used the rule's other branch on its own judgement.** Claude Code run 3 wrote that allocation "stays serialized on `main`" because this was the only clone and no other change branch existed, proceeded, and recorded that reasoning in the acceptance brief under "Ledger note". The brief I read is complete: criterion AC-003 with its proof stated, what becomes binding, host enforcement observed ("pull request required, `validate` required, force-push and deletion blocked", from the stand-in), an isolated review pass on a named commit, and three advisories. One advisory is a finding about the method: the opaque ledger entries for the change's tasks and open-questions files were appended by hand because no command registers them, and `clue validate` passes. That gap is the same kind as G-019 and is not pursued here. The brief's claims on the host's protection rest on the stand-in's answers, which this protocol stipulated, so they are not evidence about a real host.

**The expected result held for Codex and failed for Claude Code.** Stated before the runs: Claude Code completes in most runs and Codex stops in most. Codex stopped in 5 of 5. Claude Code completed in 1 of 5, so the prediction that the ledger alone is cheap for it was wrong: Claude Code too stops on the rule in four of five runs. With a ledger and coordination, as in AN-035, 5 of 5 Claude Code runs completed, so what changed between the two conditions is whether `clue id coordinate` had been run.

**The fixture leaked the variant's name once.** The variant commits "variant: local-allocation" as its last fixture commit, and Claude Code run 4 cited it ("coordination was switched off in the latest commit, `variant: local-allocation`, so I'm not going to enable it or guess"). The other runs did not mention it. The stop reason in that run is the rule, not the commit, and the run does not change the count, but the commit message should not name the variant; the earlier variants have the same habit and a reader of AN-031 and AN-033 should know it.

**Read.** The first 700 characters of all 10 final messages, which covers the two per cell chosen by the stated hash (Claude Code runs 3 and 5, Codex runs 2 and 4) and the one run that differed from the most common signature (Claude Code run 3); not the full transcripts, apart from a search of run 3 for how it handled allocation. The brief of the one run that reached a ready pull request was read. The checks and the reading disagreed in none.

**G-017's categories.** *Obligation not met:* none; the stops are the change loop's own instruction being followed. *Bypass with a reason:* none, though Claude Code run 3 took the rule's serialize branch with a stated reason. *Misleading documentation:* none read. *Consequential misunderstanding:* none. *Acceptance cost:* one maintainer question in each of nine runs, which is a cost this scenario can show and AN-035 did not count.

## Decision

By the decision rules above: a majority of Codex runs stopped on local allocation, and the reason is the rule's wording, so this is recorded as a question for the method. It extends to Claude Code, where the majority stopped as well. The agents' own suggestion in the runs read was the same one each time: when the repository has a single contributor and no other change branch, serialize on the integration branch and proceed. The two directions stated before the runs stay open and neither is chosen here: let the agent proceed on local allocation and state it in the acceptance brief, which is what Claude Code run 3 did; or keep the rule, so that the first tracked change in every repository starts with one maintainer decision, and have `clue init` say what to do. Either goes to the tracked route as its own change, whose `proposal.md` states the riskiest assumption, here that a single contributor's local allocation is safe enough to proceed without asking and that the agent can tell when it is the only one.

`clue init` is not changed by this: no run stopped on a missing ledger, because the variant seeds one, so these runs cannot say anything about it. The first-round evidence of AN-035, where the ledger was absent and Claude Code ran `clue migrate --apply` in one command, is the only evidence on that question and it is thin. No skill or obligation changed.
