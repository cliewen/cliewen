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
