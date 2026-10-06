---
id: AN-034
type: analysis
status: active
links: [P-024, G-025, G-017, AN-031, AN-033]
title: What the two long scenarios needed, and what one run each on two agents showed
---

# AN-034 — What the two long scenarios needed, and what one run each on two agents showed

## Purpose

[P-024](../plans/P-024-repeatable-scenario-trials.md)/M-108 asks for greenfield adoption to a first green `clue validate`, and a tracked change carried to an acceptance brief, as scenarios on both agents from the same harness. This is the spike that comes before the runs: can each scenario reach its end state inside the harness, what did it need that the harness lacked, and does one run on each agent show anything.

## Evidence boundary

Run on 2026-10-06 on the maintainer's Windows 11 machine with Docker Desktop and Linux containers, with `clue` and skills 0.28.0 built from `main` at `e4356dd`. Claude Code 2.1.289 (`claude-sonnet-5-5`) and Codex 0.160.0, both on subscription logins; one run per agent and scenario, four runs in all. Records are under `scenario-runs/spike/` (Git-ignored). One run each is a feasibility check and says nothing about how often anything happens.

## What the harness lacked, and what was added

- A turn limit of 25 on Claude Code and a 15-minute run limit; a scenario now carries its own (`greenfield` 60 turns and 25 minutes, `tracked-brief` 120 turns and 45 minutes), and the short scenarios keep the defaults.
- Checks could read only the transcript and `git status`. A scenario may now have a `post.sh`, run in the fixture after the agent, whose `key=value` lines the checks read; `scenario-trials check` re-scores stored runs from the saved facts.
- A fixture could not offer a remote or a hosting CLI. `tracked-brief` builds a bare `origin` and a `gh` that records what it is given and never reaches a network.

## What was observed

| Scenario | Agent | End state | Turns, time, cost |
|---|---|---|---|
| `greenfield` | Claude Code | `clue validate` green; vision `draft`, `inferred`; the agent ran `clue validate` itself | 6, 23 s, $0.11 |
| `greenfield` | Codex | the same | n/a, 101 s |
| `tracked-brief` | Claude Code | proposal committed before source; workspace deleted in the digest; `clue validate --forbid-changes` green; PR created as a draft, edited and marked ready; no template text left in the brief | 46, 151 s, $0.84 |
| `tracked-brief` | Codex | the same, except the PR was created and edited but not marked ready | n/a, 198 s |

Neither scenario showed its failure on either agent in this one run each. Claude left a vision as `draft` and `inferred`, said so, and ran the judge; the brief it wrote names the criteria it added, says it binds an opt-in `--csv` flag, and records two review passes.

**The Codex run that stopped at a draft PR is about the fixture, not the method.** Codex's last message says hosted readiness was blocked because `gh` was a logging stub: the first shim printed `ok` for every call, so `gh pr view` returned no state. Claude proceeded anyway and marked the PR ready. The shim now answers `pr view` with draft state and the head commit, `pr ready` with a confirmation, and `repo view` and `api` with plausible JSON; it was tested in a container and the two `tracked-brief` runs have not been repeated with it. A `pr=created` result from the first shim must not be read as an agent failure to complete.

**Two gaps in the adoption path showed up in the fixtures, not as agent failures.** `clue id next CH` refused in a repository made by `clue init` because there was no identity ledger, and the Claude agent ran `clue migrate --apply` first to seed one; I did not check how Codex handled it. That is how a first tracked change goes in a freshly initialised repository, and it is recorded here as observed; whether `clue init` should leave a ledger is a question for the method, not for this harness. And Claude's brief calls the vision "still the bootstrap placeholder text" in a fixture whose vision was real, so the brief repeats something that was not true; the harness checks no claim inside the brief.

## What this does not establish

That either obligation is carried by anything: one run each, both clean, so neither scenario has been shown to provoke its failure. That the fixtures are hard enough: `greenfield` starts after `clue init`, whose red `clue validate` names the three placeholders to replace, so reaching green is close to following the judge's own message; a scenario from an empty directory would test more, and the cost of the run would rise. Anything about OpenCode, or about the effect of removing a skill, which no arm here tried. The checks read the end state and not the quality of the vision, the criteria or the brief; a person reads those.

## What follows from it

Both scenarios can reach a stated end state, so the runs can proceed: five per scenario on each of Claude Code and Codex, the count and reading protocol fixed before the runs in a successor analysis, as [AN-033](AN-033-first-comparison-period.md) did. The expected result is stated now so it cannot be fitted afterwards: both scenarios show no failure on either agent, and are recorded as scenarios that can be passed and checked, not as evidence that the obligations they exercise are needed. A scenario from an empty directory, and a removal variant for each obligation, are the next step if M-108's runs show nothing, and would be their own change.
