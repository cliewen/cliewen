---
id: AN-030
type: analysis
status: active
links: [P-024, G-025, G-017, AN-029, AN-028, IDR-010]
title: The routing scenario on two agents through one interface, and what differed
---

# AN-030 — The routing scenario on two agents through one interface, and what differed

## Purpose

[P-024](../plans/P-024-repeatable-scenario-trials.md)/M-104 asks for two agents to run the same scenario through one adapter interface from a clean configuration, with what each could still see of the host stated. This records the runs, whether Codex's event stream could carry the same checks as Claude Code's, and the one difference between the agents that the runs show.

## Evidence boundary

Run on 2026-10-05 on the maintainer's Windows 11 machine with Docker Desktop (server 29.7.2, Linux containers on WSL2), with a `clue` and skills (version 0.28.0) built from commit `5bd7247f54f1cd9b5959697419de43faf0c2488b`, which is `main` plus CH-200's proposal, and the harness code itself taken from the working tree of CH-200 at that time. Both invocations used image `cliewen-trial:17f1c40ab4c3`. Agents: Claude Code 2.1.289 (`claude-sonnet-5-5`) with a subscription token, and `codex-cli` 0.160.0 whose model name the harness read from Codex's own session record as `gpt-6.1-sol`, with the maintainer's subscription login; no API key for either. Both are in one container image with pinned versions. This is a **prepared environment**: Docker, a Go toolchain, and two logged-in subscriptions on the host. The scenario is the routing scenario of [AN-029](AN-029-five-routing-runs-on-claude-code.md), unchanged. Records are under `scenario-runs/`, ignored by Git.

The sample is Codex five runs (`20261005-152601-routing-codex`), Claude Code five runs (`20261005-152944-routing-claude`) and, for Claude, AN-029's five (`20261005-145310-routing-claude`), ten in all. Two single smoke runs on the new interface and one Codex smoke run before it are not part of the sample.

## What was observed

**Codex's stream carries what the checks need, and the order is preserved.** The first Codex run's native events were read beside the neutral transcript the adapter produced: text, shell commands and the file change appear in the same order, and the login-shell wrapper Codex records around each command is removed before the edit detector sees it. The stream has no turn count or cost, so those fields are zero for Codex and the harness records wall time for both agents instead: 41 to 46 seconds across Codex's five runs and 10 to 15 across Claude's five in this invocation. The probe read Codex's version and model from inside the container.

**The two agents differ in when they state a route, and in whether they read the routing skill.** After the checks were revised, both agents recommended the direct route for the typo and the tracked route for the CSV export in every run, did not start the export, and asked where the tool's source lives. Claude Code edited the README before stating any route in 8 of its 10 runs. Codex stated both routes before editing in all 5. The counts give a two-sided Fisher exact p of 0.007 for that difference; it describes these fifteen runs on one prompt and does not say how often either agent would do it. Scanning each run's commands for the routing skill and its scope reference, Codex read the routing skill in all 5 runs and Claude Code in none of its 10. Claude's runs read the README and the vision and answered from `AGENTS.md`. That the skill was read is an observation; that reading it is what made Codex state the route first is a hypothesis these runs do not test.

**The checks misjudged one run per agent at first.** One Codex run wrote "direct for the README typo; tracked for CSV export" as two clauses of one sentence, and one Claude run said "so I treated it as a direct change" two sentences after naming the typo. The loose route reader now splits on semicolons and carries the last named topic forward through a block of text. After this the checks agree with my reading of all fifteen transcripts, which is by construction, as in AN-029. Re-checking AN-029's five Claude runs with the revised reader gave the same signatures as before.

**What each agent could see of the host.** Neither transcript names CodeGraph, which only the maintainer's own instructions mention. Claude Code ran with `--setting-sources project,local` and its initialisation event lists no MCP servers; its plugin and skill counts were not traced, as in AN-029. Codex ran with a login directory created by a login inside a container, mounted over the container's `.codex`, so its configuration is that directory's and not the host's `~/.codex`; the directory's other files were not read. Both agents can read their own credential from inside the container, Claude Code's from an environment variable and Codex's from the mounted `auth.json`: the credential is the maintainer's own subscription login in a throwaway container, and neither reaches an image layer, an argument list, a result file, or a log. Codex runs with its own sandbox off because it cannot start in the container, and the container is the boundary.

**Where vendor names remain.** The runner, the checks, the conditions record, the summary and the command line name no vendor; each adapter's file does, as do its tests and the container image's Dockerfile, which installs both agents. Adding an agent means a file that registers an adapter, and does not touch the runner.

## What was rejected

**A terminal-driving orchestrator (Orca).** Recorded in [IDR-010](../decisions/IDR-010-agents-are-reached-through-adapters-over-one-interface.md). It was set aside on a two-page reading of its documentation and not run.

**A second copy of the Claude path for Codex.** Rejected for the reason in IDR-010, and because the checks would have named vendors.

## What this does not establish

Fifteen runs of one scenario on one prompt, two agents, and a hypothesis about why they differ that was not tested. It says nothing about how often either agent behaves a given way beyond these runs, about other prompts or scenarios (M-105), or about a model reached through another tool (M-106). The container runs Linux and bash, so nothing here concerns PowerShell. Codex's turn count and cost are not measured.

## What follows from it

M-104's evidence is this record and IDR-010. The difference between the agents is a first observation in G-017's categories of an obligation followed by one agent and, in its first step, not by the other. For M-107, two things follow: the removal trial can use the routing-state difference between the agents as a check that the harness sees a real effect before it is used on a removal, and the runs per arm have to be stated beforehand in terms of the claim, which the declared revision of M-107 in P-024 does.
