---
id: AN-032
type: analysis
status: active
links: [P-024, G-025, G-017, AN-031, AN-030]
title: One OpenCode run on a model that is neither agent's default
---

# AN-032 — One OpenCode run on a model that is neither agent's default

## Purpose

[P-024](../plans/P-024-repeatable-scenario-trials.md)/M-106 asks whether a model reached through OpenCode runs the same scenario through the existing interface with no code specific to that model, whether OpenCode honours `AGENTS.md` and the repository's skills, and what leaks from its configuration. This records one run of five on one model.

## Evidence boundary

Run on 2026-10-06 on the maintainer's Windows 11 machine with Docker Desktop and Linux containers, with `clue` and skills 0.28.0 built from `main` at `fc9bfa9`, and the OpenCode adapter taken from the working tree before it was committed as `f5ad502`. Agent: OpenCode 1.18.34 reaching `openrouter/qwen/qwen3.8-27b` through OpenRouter on an API key, with the reasoning effort fixed at `medium` by the adapter. This is a **prepared environment**: Docker, a Go toolchain, and a key in a data directory outside the repository holding only the OpenRouter entry. The container carries no OpenCode configuration of the maintainer's; its state, cache and config directories point at `/tmp`. Records are under `scenario-runs/20261006-171327-routing-opencode/` (Git-ignored).

| Scenario | Agent | Model | Variant | Runs |
|---|---|---|---|---|
| `routing` | OpenCode | `qwen3.8-27b`, effort `medium` | baseline | 5 |

An earlier invocation of the same five runs ended in zero turns, because Docker created the parents of the credentials mount as root and OpenCode could not create its state directory. The adapter now sets `XDG_STATE_HOME`, `XDG_CACHE_HOME` and `XDG_CONFIG_HOME` to `/tmp`; that invocation is not counted.

## What was observed

All five runs finished, in 7 to 13 turns, 56 to 164 seconds and $0.02 to $0.05 each. The runner, the checks and the summary needed no change for the agent: the adapter is one file, one line in the Dockerfile and its tests.

| Check | Result over five runs |
|---|---|
| `clue latest --quiet` run first | Seen in runs 1 and 5, which were read through; not tallied for runs 2 to 4 |
| Route stated for the typo | `direct` in 5 of 5 |
| Route stated for the CSV export | `direct` in 2, `tracked` in 3 |
| Edit made before a route was stated | Yes in 4, no in 1 |
| Configuration leak, as the harness defines it | None |
| Skill files read from `.agents/skills` or `.claude/skills` | None in 5 of 5 |

Reasoning tokens in the stream were 1,429 to 3,586 per run, so the model thought at `medium` as set; they are in the stream's token counts, not in the events the adapter turns into a transcript.

OpenCode honoured `AGENTS.md`: in the two runs read through, the agent opened with the update check the hub requires and later wrote `Recommended route: direct`. No run touched a skill file, so this run does not show whether OpenCode lists the repository's skills to the model; an agent that never needed one gives no evidence either way. The CSV tool the prompt names does not exist in the fixture, and several runs said so and asked where it is; that is a property of the scenario, not of this agent.

## What this does not establish

One model, one scenario, five runs: it shows the interface carries a third agent, not how that agent behaves in general. The 4 of 5 edits before a route is a different distribution from the two agents in [AN-031](AN-031-four-scenarios-and-a-removal-trial.md), on a different scenario and a much smaller model, and no comparison is claimed. Whether skills are discovered needs a scenario that cannot be done without one, such as `routing-code` or `upgrade`. The configuration-leak check looks at the transcript the adapter produces, and I did not inspect OpenCode's own session records for anything the checks could not see. Runs 2 to 4 were read only through the checks' summary.

## What follows from it

M-106's exit criterion is met on the interface: a model reached through OpenCode ran the scenario with no code specific to that model. The open question this leaves for M-107 is skill discovery, answered by running a scenario that needs a skill on this agent. M-107 would also need to decide whether the effort setting is a per-run condition recorded in `conditions.json`, which today it is not.
