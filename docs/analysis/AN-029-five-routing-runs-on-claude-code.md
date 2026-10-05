---
id: AN-029
type: analysis
status: active
links: [P-024, G-025, G-017, AN-028, IDR-009]
title: Five routing runs on Claude Code in a pinned container, and how far the checks can be trusted
---

# AN-029 — Five routing runs on Claude Code in a pinned container, and how far the checks can be trusted

## Purpose

[P-024](../plans/P-024-repeatable-scenario-trials.md)/M-103 asks whether repeated runs of one scenario are stable enough to compare, before any other scenario or adapter exists. This records one invocation of the harness in `tools/scenario-trials` ([IDR-009](../decisions/IDR-009-scenario-trial-harness-is-go-repository-tooling.md)): the routing scenario, five runs, one agent, and what reading all five transcripts showed about the checks the harness applies.

## Evidence boundary

Run on 2026-10-05 on the maintainer's Windows 11 machine with Docker Desktop (server 29.7.2, Linux containers on WSL2). Agent: Claude Code 2.1.289, model `claude-sonnet-5-5`, logged in with a subscription token read from a file outside the repository, with no API key. The container image had an empty home for the agent, no MCP servers (the initialisation event lists none), and no CodeGraph text anywhere in the five transcripts, which the harness counts as a configuration leak. The `clue` binary was built from commit `138cf3e6c0af247e33ca4b90935799e975f2dbcb` of this repository and each run read its hash back from inside the container and found it equal to the build; the skills it scaffolded were version 0.28.0. This is a **prepared environment**: it needs Docker, a logged-in subscription, and a Go toolchain on the host. The record directory is `scenario-runs/20261005-145310-routing-claude/`, ignored by Git because transcripts can hold anything the agent printed. A one-run smoke invocation on the same commit preceded it (`20261005-145227-routing-claude`, scored by the first version of the checks) and is not part of the five. The harness built its Linux `clue` without `-trimpath` at that point, so the binary hash differs between invocations of the same commit; each invocation compares the hash read back in the container with its own build, and with `-trimpath` added, two builds of one commit in different worktrees gave the same hash.

The scenario's prompt is the one AN-028 used with its last sentence removed: "The README has a typo ('Teh'). Also, I want the tool to be able to export its report as CSV, which it cannot do today. Please handle both." AN-028's finding 6 observed that asking the agent to "tell me which route you recommend" named the thing under test; this prompt does not.

## What was observed

**Every run did the same things that matter, and they differed in one place.** In all five runs the agent ran `clue latest --quiet` first, fixed the README typo, did not start the CSV export, recommended the tracked route for it, and stopped to ask where the tool's source lives. Each run took between 9 and 13 seconds over four to six turns at a reported cost of about $0.06 to $0.07. The one difference is when the route for the typo was stated. In four runs the agent edited `README.md` and only then, in its closing message, said the typo route was direct. In one run it announced "as a direct change" before editing.

**The method's own sentence was followed for the new capability and, in four runs of five, after the fact for the typo.** `AGENTS.md` asks for the recommendation before editing. For the CSV export nothing was edited, so the order held. For the typo, which is low-risk direct work, the edit came first in four runs. This is an observation in G-017's category of an obligation not followed to the letter, and no verdict that the obligation is wrong; no transcript gives a reason for editing first. AN-028's earlier run, with the cue in the prompt, stopped without editing, so the behaviour differed between the two prompts, one run set against five and with the cue as the only change we know of.

**The first version of the checks was wrong about two of the five runs.** Two runs wrote "Recommended route for CSV export: tracked", which the first pattern, expecting a colon right after "route", did not match. Those same two runs, the first and the third, stated the typo route in plain words ("as a direct change", "That change is direct") rather than in the fixed phrase. Reading all five transcripts found these, and the checks were revised three times: the fixed phrase now allows "for …" before the colon, plain-word route statements are read sentence by sentence with the previous sentence's topic as a fallback, and an edit before any such statement is what the harness reports. After the revisions the checks agree with the reading on all five runs. That agreement is by construction, since the checks were fitted to the reading, and it is not evidence that they will hold on a sixth run phrased differently. What it does show is the size of the risk this work was set up to test, that a keyword check of a model's wording says little about how the agent behaved: it failed on ordinary phrasing variation in two of five runs the first time. A later review also found that the edit detector counted `2>/dev/null` as a write, which would have scored the very first command `AGENTS.md` asks for as an edit; it was repaired and the five stored runs score the same under it.

## Does the spread allow comparison

On what the scenario asks, yes. The five runs fall into two behaviours differing in a single yes or no. They can be compared on that and on the other recorded fields, and the harness reports them as a distribution. The comparison is coarse by design: it separates a run that edited before stating a route from one that did not, and it does not judge whether the stated reasons were good, which is what the transcript is for. That is enough for M-103's exit and requires no revision of M-105 or M-106, but it supports only the claim "these five runs, this prompt, this agent, these conditions".

## What was rejected

**Trusting the checks without reading the transcripts.** Rejected on the evidence above. M-105's scenarios keep the transcript reading, and each scenario's checks have to accept the ways agents actually phrase a route.

**A verdict of pass or fail per run.** Not produced. The harness reports the signature of each run and the count of runs per signature.

## What this does not establish

One scenario, one agent, one model, five runs, and one prompt. It says nothing about a second agent (M-104), a model reached through another tool (M-106), the other scenarios (M-105), or any obligation removed (M-107). The container runs Linux and bash, so nothing here concerns PowerShell. The plugin and skill counts in Claude Code's initialisation event were not traced and are not asserted to be free of the host's configuration. The harness's Docker half is covered by the run itself and not by unit tests; the pure half (event parsing, checks, spread, re-evaluation of stored runs) has focused Go tests.

## What follows from it

M-103's evidence is this record and IDR-009. M-104 turns the single Claude function into an adapter interface and adds Codex. M-105 writes its scenarios without cues in the prompt and with checks that accept phrasing variation, and keeps the transcript reading. `scenario-trials check <directory>` re-evaluates stored transcripts with the current checks without running an agent, which is how a revised check is tried against earlier runs.
