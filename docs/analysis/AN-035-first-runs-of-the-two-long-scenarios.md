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
