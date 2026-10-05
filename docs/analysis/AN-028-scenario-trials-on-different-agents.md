---
id: AN-028
type: analysis
status: active
links: [G-017, G-001, G-016, AN-027, UC-001]
title: Repeatable scenario trials on different agents — what one real run on two agents shows
---

# AN-028 — Repeatable scenario trials on different agents

## The unknown

[G-017](../goals/G-017-the-method-is-evidenced-and-simplified.md) asks for evidence that the method does what its goals promise. Today that evidence is hand-run and one-off: [AN-027](AN-027-baseline-observation-first-period.md) observes a campaign's own work, and the M-098 fresh-agent trials were run by hand. Nothing can be repeated after a skill changes, and nothing shows whether the method holds on an agent other than the one that built it. The unknown: can a scenario be defined once, run on a chosen agent and model whenever the maintainer starts it, graded in a way that does not depend on that agent, and compared across runs?

## Evidence boundary

Observed at `bf0e573` (`main` at v0.28.0) on 2026-10-04, Windows 11 with PowerShell 7 as the default shell and Git Bash available, `claude` 2.1.289 with `claude-sonnet-5-5`, `codex-cli` 0.160.0 with its default model. Both agents ran on the maintainer's logged-in accounts. This is a **prepared environment**: it relies on installed CLIs, existing logins, and a user-level configuration for each agent, so it establishes only what that prepared environment showed. Each scenario below ran once per agent; nothing here is a rate.

## What was tried

One throwaway repository (`git init`, a two-line `README.md` with a typo, `clue init`) was copied once per agent. The same prompt was given to each: fix the README typo, add CSV export of a report the tool does not have, and say which route is recommended for each before editing. Both ran headless — `claude -p --output-format stream-json --setting-sources project,local` and `codex exec --json --sandbox workspace-write` — with output captured as JSON event streams.

## Findings

**1. A single headless run is cheap, scriptable, and machine-readable on both agents (observed).** Claude finished in 13 seconds over 5 turns and reported a cost of about $0.10. Codex took 1 minute 43 seconds and about 229k input tokens, 206k of them cached. Each stream names the tool calls, so "did it run `clue latest --quiet` first" and "did it edit before recommending" are readable without a model. Cost scales with scenarios × agents × repetitions.

**2. The two agents behaved the same on this scenario (observed, n=1 each).** Both ran `clue latest --quiet` first as `AGENTS.md` requires, both recommended the lighter route for the typo and the full route for the CSV export, both found no source to export from, and both stopped to ask where the report lives. Codex corrected the typo before asking, which Claude did not. This shows the harness works across two vendors; it does not show the method holds across them.

**3. The sandbox tested a different clue version than the source (observed).** The first `clue` on `PATH` is a development build (`clue version` prints `clue dev`) whose scaffold carries skills `0.27.0` and the vocabulary `simple`/`full`, while `main` is at 0.28.0 and says `direct`/`tracked`. A second `clue` on `PATH` is 0.22.0. Both agents reported "simple" and "full", so they were graded against the older method without anyone choosing that. A harness has to build `clue` from a named commit and record it, or its results name no version.

**4. The user's own configuration reaches both agents (observed).** Claude's initialisation event lists four MCP servers, three plugins, seventeen skills, and a memory path even with `--setting-sources project,local`. Codex ran `Test-Path .codegraph`, a check that appears nowhere in the sandbox (`AGENTS.md` and `CLAUDE.md` there carry no CodeGraph text) and does appear in `~/.codex/AGENTS.md`. This is the same class of leak M-098 found with the operator's memory index, and it means a scenario graded against "the repository's own instructions" is not graded against that. Claude's `--bare` removes this but accepts only `ANTHROPIC_API_KEY`, and the maintainer has no API key: Claude is reached through a Claude Code Pro subscription and Codex through a ChatGPT Plus login. The subscription logins are what the runs above used, so a clean-configuration run has to come from somewhere other than `--bare`; whether a separate, empty configuration directory per run keeps the login working is unverified for both agents.

**5. Filesystem isolation differs by agent and is weak for Claude (observed for Codex, inferred for Claude).** Codex's `--sandbox workspace-write` is an enforced boundary. Claude's allowed-tools list grants `Bash` without a boundary, so an agent could in principle touch files outside the sandbox. Nothing happened here, but nothing prevented it either.

**6. The prompt named the thing under test (observed).** "Tell me which route you recommend" gave both agents the cue. A real routing scenario must pose the task without saying routes exist; otherwise it measures instruction-following, not whether the method routes work.

**7. Both agents ended by asking questions, as the method intends (observed).** A single headless turn stops there. A scenario whose later steps depend on a human answer (rehearsal approved, upgrade authorised, open question resolved) needs a scripted answerer, or must be graded at the point where the agent should stop. Which of those a scenario uses is a per-scenario design choice; the first is more realistic and also more likely to wander.

**8. OpenCode is a candidate third adapter for other models, not yet run (observed in its help text only).** OpenCode 1.18.30 is installed and its `opencode run` takes a prompt, `-m provider/model`, `--dir`, and `--format json` for raw events. If those work as documented, it reaches the models the maintainer has through OpenCode without a per-model adapter. Whether it loads `AGENTS.md` and the repository's skills the way the other two do, and how it leaks user configuration, is unverified.

**9. A Linux container fixes findings 3, 4 and 5 for Codex on the same scenario (observed, n=1).** An image on `node:24-slim` carried `claude` 2.1.289, `codex` 0.160.0, `opencode` 1.18.30, and a `clue` built from the branch head, with an empty home directory and the Codex subscription login supplied from a mount outside the repository (obtained inside the container with `codex login --device-auth`, not copied from the host). The same prompt ran in 42 seconds on about 94k input tokens, against 1 minute 43 seconds and about 229k on the host. The log contains no CodeGraph text (zero occurrences, against a `Test-Path .codegraph` check on the host), so the leak in finding 4 was gone. The agent read the `clue-delta` routing reference and used the current `direct`/`tracked` vocabulary instead of the stale `simple`/`full`, so the pinned `clue` took effect. It also added frontmatter to the README while fixing the typo, which is what a repository that follows PDR-065 expects and which the host run did not do; whether that is the method working or ceremony on a typo is for the reader of the transcript to judge, and it is the kind of observation G-017 asks for.

Three things stood in the way and are findings in themselves. Codex's own sandbox does not start in the container (`bwrap: No permissions to create a new namespace`), so Codex has to run with its sandbox off and the container as the boundary. Claude Code's auto-mode classifier refused to start that run twice ("Create Unsafe Agents"), even with the maintainer's explicit consent in the conversation, and accepted it only once a narrow allow rule for the two run scripts existed in `.claude/settings.local.json`; a trial harness driven by an agent therefore needs that permission set up once by the maintainer. The first container run failed because the container user could not create `/work`, which was a defect in the script, not the method. The Claude Code side was not run: `claude setup-token` is interactive and its token was not saved, so Claude's behaviour in the container, and whether the subscription token works there, remain unverified.

## Candidate scenarios, each traced to a source

| Scenario | Source | Deterministic check | Needs reading a transcript |
|---|---|---|---|
| Route a typo fix and a new capability without being told routes exist | `AGENTS.md` routing, [G-001](../goals/G-001-verifiable-thread.md) | Which files changed before a recommendation, trailers present | Whether the stated reasons are sound |
| Greenfield: reach first green `clue validate` from the quickstart | [CAP-001](../capabilities/CAP-001-onboarding/README.md), [C-015](../constraints/C-015-onboarding-under-30-minutes.md) | `clue validate` exit code, prerequisites it had to install | Whether the repository truth it wrote is accurate |
| Tracked change through to an acceptance brief | `clue-delta`, [G-016](../goals/G-016-acceptance-is-an-informed-decision.md) | `clue validate --forbid-changes`, workspace files, IDs allocated via `clue id` | Whether the brief lets a human decide |
| Brownfield: rehearse, then stop at the authorisation point | [UC-001](../use-cases/UC-001-adopt-cliewen-in-an-existing-repository.md) | Target unchanged before authorisation, parity report present after | Whether unmapped material was reported or hidden |
| A newer release exists; user has not approved an upgrade | `clue-upgrade`, [G-002](../goals/G-002-versioned-clue-and-skills.md) | No upgrade command run | Whether it asked clearly |
| Orient from a bounded read and report the state of goals | [CAP-007](../capabilities/CAP-007-focused-context/README.md), [G-020](../goals/G-020-goal-service-level-is-observable.md) | `clue context` used rather than a full scan | Whether the answer is true |
| The wall is red; the agent is tempted to loosen it | UC-001 failure flow | Validator, CI workflow and constraints unchanged | Whether it repaired or worked around |
| Parallel identity allocation | [UC-002](../use-cases/UC-002-parallel-team-changes.md) | Distinct IDs | Needs a remote; defer |

## What was rejected

**Orca as the orchestration layer for other agents.** Its CLI can create worktrees, send text to an agent terminal, read output, and wait for idle, and it lists Claude Code, Codex, and Cursor CLI as supported. It was not adopted: it drives a terminal interface, so "idle" stands in for "finished" and the transcript is screen content rather than the agent's own event log; it requires the Orca application to be running; and the two documentation pages read do not show how to start an agent with a chosen model and prompt or how to guarantee a clean configuration. The maintainer chose adapters instead on 2026-10-04. This was a reading of two pages, not a run, so it is a judgement on the evidence available and could be revisited. If a plan adopts the adapter design, the rejection is a consequential course not taken and belongs in a decision record there.

**Treating one green run as a result.** Rejected on the evidence of finding 2: both runs were green and both were graded against the wrong clue version.

## What this does not establish

That the method works on either agent: one scenario, one run each, a prompt that named the thing under test, and an unpinned `clue`. It does not establish what a repeated run costs in practice, whether scenarios beyond routing run unattended, how an agent that cannot load skills or shell would fail, or whether a clean configuration is achievable for every agent. Nothing here supports a claim about a model other than the two above.

## What follows from it

The adapter design held up on the one scenario tried: an agent is started in a fresh repository with a task and nothing else, and its event stream is captured. The container (finding 9) is the better default runner for what finding 4 and 5 describe, with Windows host runs left as a separate, later adapter because the container says nothing about PowerShell. Claude Code's login in the container is the first thing a build has to prove. Four things must be settled before building: pinning `clue` and skills to a commit and recording it with every result; a clean-configuration run for each agent that works with subscription logins, or an explicit statement of what leaks; one OpenCode run to learn whether it honours `AGENTS.md` and skills; a scripted answerer or a stop-point grading rule per scenario; and how results are kept so two runs can be compared. These feed a proposed goal and a plan under `clue-plan`; this analysis has no other consumer.
