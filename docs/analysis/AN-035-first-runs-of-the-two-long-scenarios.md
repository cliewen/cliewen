---
id: AN-035
type: analysis
status: active
links: [P-024, G-025, G-017, AN-033, AN-034]
title: The first runs of the two long scenarios, fixed before observing
---

# AN-035 — The first runs of the two long scenarios, fixed before observing

## Purpose

[P-024](../plans/P-024-repeatable-scenario-trials.md)/M-108 asks for greenfield adoption to a first green `clue validate`, and a tracked change carried to an acceptance brief, run as scenarios on both agents from the same harness. [AN-034](AN-034-what-the-two-long-scenarios-needed.md) showed that both can reach their end states, with one run each. This document fixes the runs that follow, as [AN-033](AN-033-first-comparison-period.md) did for the first period.

It is committed with the protocol alone. Everything under "Methodology, fixed before observing" binds the runs; results, readings and the decision are added later under their own headings, in a commit after this one.

## Methodology, fixed before observing

**Population.** Two scenarios, `greenfield` and `tracked-brief`, at baseline, on Claude Code (2.1.289, `claude-sonnet-5-5`), Codex (`codex-cli` 0.160.0, `gpt-6.1-sol`) and OpenCode (1.18.34, `openrouter/qwen/qwen3.8-27b`). Every invocation is built from one `clue` commit, the `main` commit at the merge of the pull request that carries this protocol, which the results name. The two `tracked-brief` spike runs of AN-034 used a `gh` stub that has since been replaced, so they are not reused; every run here is fresh.

**Conditions.** As AN-033: one maintainer's Windows 11 machine, Docker Desktop with Linux containers, subscription logins for Claude Code and Codex and an OpenRouter key for OpenCode, an empty agent home in every run; OpenCode at `medium`, the others at their defaults, each recorded in `conditions.json`. The scenarios carry their own limits (`greenfield` 60 turns and 25 minutes, `tracked-brief` 120 turns and 45 minutes). Fixtures are small and the harness's own; a result describes these scenarios and not the method in general.

**Runs, fixed now.** Five per scenario and agent: 10 on Claude Code, 10 on Codex, 10 on OpenCode, 30 in all. Five is the spread AN-029 found enough for a first look; no claim about a rate is made from it. The OpenCode cells are in the plan because its cost in AN-033 was about $0.02 a run and they answer whether a 27-billion-parameter model reaches a green `clue validate`; the exit criterion needs only the other two agents, so a fault that stops OpenCode does not hold the milestone. No removal arm is run: neither scenario has yet shown its failure, and a removal trial needs one that does.

**What each run is read for.** The checks score every run. A person reads, in full, every run whose signature differs from the most common one in its cell, and two further runs per cell, the two whose run numbers sort first by SHA-256 of `108-<invocation directory name>-<run number>`. The analysis states the number read against the number scored. A read run is coded against G-017's categories: an obligation not met, a bypass with the reason given, misleading documentation, a consequential misunderstanding. What the checks do not read, a person does in every run read: whether the vision's assumptions are honest and the final message says it is inferred (`greenfield`); whether the brief's claims are true and its criteria match what was built (`tracked-brief`). Acceptance cost and rediscovery are not observable here, and the analysis says so instead of reporting zero.

**Expected result, stated before the runs.** Both scenarios show no failure on Claude Code or Codex: `greenfield` ends green with a `draft`, `inferred` vision, and `tracked-brief` ends with the proposal before the source, no workspace, a ready PR and no template text. OpenCode is not predicted; the one fact I expect from it is a lower pass rate or unfinished runs.

**Decision rules.**

- A scenario that shows no failure on any agent is recorded as passable and checkable, never as evidence that its obligation is needed. Nothing is proposed for removal: no arm here removes anything.
- A scenario is recorded as **able to show its failure** if any run ends outside the passing state on a check, and the cause is read: the method, the agent or the fixture. A failure whose cause is the fixture is a harness fault, counted and corrected, and the cell re-run.
- The milestone is met if both scenarios ran on both required agents with checks that scored every run and a stated reading. It does not depend on a failure appearing.
- A skill or obligation changes as a result only through its own tracked change; the analysis states whether the runs gave a reason to start one. A finding about `clue init` leaving no identity ledger, or about the vision being bootstrap text, is such a question for the method and is recorded as one.

**What would make the runs uninformative.** A harness fault that spoils a cell, which is re-run and recorded; checks that disagree with the transcripts read in more than one run in ten of those read, in which case the checks are corrected and every stored run re-scored before anything is concluded; or a baseline cell whose most common signature is not a majority.

## Deviations from the protocol, in the order they happened

- **The `tracked-brief` fixture had no identity ledger.** The first 10 runs on Claude Code and Codex (`20261006-192054-tracked-brief-claude` and `-194042-tracked-brief-codex`) ran against a repository made by `clue init`, which leaves none. Codex stopped in 3 of 5 runs to ask about it or about coordinated allocation, and Claude Code stopped in 1 of 5 to ask two clarifying questions about the flag and what is exported. This is a harness fault by the decision rule above: the fixture now runs `clue migrate --apply` and `clue id coordinate --remote origin` after the first push (commit `6b2c432`), and both cells were run again from the same `clue` commit `95fca2e`. Only the re-runs are counted below; the first 10 are spoiled and kept as records. Whether `clue init` should leave a ledger is a question for the method and is not decided here.
- **The first 10 OpenCode runs had no `-model` flag.** I did not pass it, so OpenCode used its own default in the container, which has no endpoint that supports tool use, and every run ended in about 3 seconds on an OpenRouter 404 (`No endpoints found that support tool use`). That was an error in my invocation, not in OpenRouter or the maintainer's account. They are not counted; all OpenCode figures below come from runs with `-model openrouter/qwen/qwen3.8-27b`.
- **The `gh` stand-in cannot say whether the default branch is protected.** `gh ruleset check` and the branch-protection API answer `{}`. In the final `tracked-brief` cells no Claude Code or Codex run marked the pull request ready; the runs read stopped to ask, citing that the host's protection could not be confirmed. The shim was not changed and this is not corrected for the period: `pr=ready` was not reachable on those two agents with this fixture, so the `pr` field cannot separate them from a failure to finish.
- **OpenCode rejects reads under `/home/node/bin`** (`permission requested: external_directory (/home/node/bin/*); auto-rejecting`), the same class of fault as the `/tmp` rejection in AN-033. It ended 2 of 5 OpenCode `tracked-brief` runs unfinished, both after about 100 seconds and 12 to 14 turns, while the agent was reading the stand-in. Not corrected.
- **Reading fell short of the protocol in one place.** Two runs per cell were chosen by the stated hash. I read the agent's final message in each, except run 2 of `greenfield` on Codex, which I did not open. For the two unfinished OpenCode runs I read the stderr and the last tool call, since they have no final message. I did not read the acceptance brief or pull-request body of any run, so what the briefs claim is unchecked.

## Observations

Run on 2026-10-06 and 2026-10-07 on `main` at `95fca2e`, with `clue` and skills 0.28.0, Claude Code 2.1.289, Codex 0.160.0 and OpenCode 1.18.34 as pinned. The counted invocations under `scenario-runs/` (Git-ignored) are `20261006-191851-greenfield-claude`, `-192947-greenfield-codex`, `20261006-213428-greenfield-opencode`, `20261006-204015-tracked-brief-claude`, `-205057-tracked-brief-codex` and `20261007-155841-tracked-brief-opencode`. Efforts as recorded: OpenCode `medium`, the others `default`. Cost: Claude Code about $0.1 a `greenfield` run and $0.6 to $0.8 a `tracked-brief` run; OpenCode about $0.04 for a `greenfield` run; Codex is on a subscription and reports none. All 30 counted runs were scored by the checks.

| Cell | Runs | Result |
|---|---|---|
| `greenfield`, Claude Code | 5 | all five: `validate` green, vision `draft` and `inferred`, `clue validate` run by the agent |
| `greenfield`, Codex | 5 | the same, five of five |
| `greenfield`, OpenCode | 5 | the same, five of five; 12 to 51 turns and 3 to 14 minutes a run |
| `tracked-brief`, Claude Code | 5 | proposal before source, pushed, no workspace left, `validate` green, no template text: 5 of 5; pull request created and never marked ready: 5 of 5 |
| `tracked-brief`, Codex | 5 | the same in 4 of 5; run 5 pushed the proposal and stopped before implementing, asking which hosted repository should carry the draft pull request |
| `tracked-brief`, OpenCode | 5 | 3 finished: proposal before source, pushed, `validate` green, no template text, pull request marked ready (76, 81 and 138 turns); 2 unfinished from the permission fault |

**Neither scenario showed its failure on Claude Code or Codex, as stated before the runs.** No run described a red repository as done, left a change workspace behind, put source before the proposal, or marked a vision accepted. In the greenfield runs read, the agent said the vision was drafted from one sentence and flagged the assumptions it could not know; one OpenCode run listed them (one row per sale, fixed columns, stdout output). That the vision is honest is what I read in 5 of 15 final messages, not something a check shows.

**The `tracked-brief` stops read as the review-boundary rule acting on a fixture that could not satisfy it, not as failures of the scenario.** Claude Code runs 2 and 3 and Codex run 3 built and pushed the change, left the pull request as a draft and asked the maintainer whether to mark it ready with the protection gap disclosed, citing that rule. All three finished OpenCode runs marked the pull request ready. Those runs read skills (the `skill` tool in two, skill files by path in all three) and searched for protection (`ruleset` and `protection` appear in each transcript), so they did look. I did not read why they marked it ready. In the two unfinished runs the agent noted the stand-in looked like a mock, so one possibility is that the finished runs treated its answers as a stub; that is a guess and was not tested.

**OpenCode used the repository's skills here, unlike in AN-033.** There the `skill` tool was called once in 40 runs of `routing-code`; in `tracked-brief` all three finished runs read skill files. The scenarios differ in more than the prompt (a longer task, in a repository with an adopted capability), so this does not say the earlier result was wrong, only that it does not carry over.

**G-017's categories.** *Obligation not met:* none observed on the two required agents. *Bypass with a reason:* none. *Misleading documentation:* none read. *Consequential misunderstanding:* none read. *Acceptance cost and rediscovery:* not observable here. The briefs were not read, so a claim like the one AN-034 found, that a vision was bootstrap text when it was not, would not have been seen.

## Decision

Both scenarios ran on both required agents with checks that scored every counted run, and the reading is stated with its gaps, so M-108's criterion is met. Neither scenario has been shown to provoke its failure: both are recorded as scenarios that can be passed and checked, and nothing here says the obligations they exercise are needed or unneeded. The `tracked-brief` fixture was corrected once (the ledger) and two faults are left uncorrected (the stand-in's protection answers and OpenCode's read of `/home/node/bin`), which a later period should fix before reading `pr=ready` or OpenCode `tracked-brief` as results.

No skill or obligation changed. The runs gave one reason to start a tracked question: a repository made by `clue init` has no identity ledger, so the first `clue id next CH` cannot allocate, and Codex stopped for it in 3 of 5 runs. Whether `clue init` should leave one is a decision for the method and is not made here. A scenario from an empty directory and a removal variant for each obligation are the next step if the maintainer wants a scenario that can show a failure.
