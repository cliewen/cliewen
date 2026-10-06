---
type: contributor-guide
title: Scenario trials
---

# Scenario trials

Repository tooling for [P-024](../../docs/plans/P-024-repeatable-scenario-trials.md): it runs a scenario on an agent in a container and records what happened. It is not part of `clue`, is not shipped to adopters, and nothing starts it automatically. See [IDR-009](../../docs/decisions/IDR-009-scenario-trial-harness-is-go-repository-tooling.md).

## Prerequisites

- Docker with Linux containers, a Go toolchain, and Git on the host.
- A login for each agent you run, outside the repository. For Claude Code, a subscription token from `claude setup-token`, saved as a single line in `~/.cliewen-trial/token`. For Codex, a login directory made by running `codex login --device-auth` inside the container with `~/.cliewen-trial/codex` mounted over `/home/node/.codex`. Write it by pasting into an editor and saving; never pass it on a command line, where its text could become a file name or reach shell history or a log. Revoke and replace it if it ever appears in a transcript or a repository. For OpenCode, a data directory `~/.cliewen-trial/opencode` holding an `auth.json` with only the provider key the run needs, mounted over `/home/node/.local/share/opencode`; run it with `-model provider/model`, for example `-model openrouter/qwen/qwen3.8-27b`. The adapter runs at reasoning effort `medium` unless `-effort WORD` says otherwise, and `conditions.json` records the effort (`default` for an agent that takes none). The container carries no OpenCode configuration of the maintainer's.

## Run

```text
go run ./tools/scenario-trials run -agent NAME -runs 5
go run ./tools/scenario-trials check scenario-runs/<invocation>
```

`run` builds a Linux `clue` from `-clue-commit` (default `HEAD`) in a throwaway worktree, builds a container image keyed to that binary, and runs the scenario `-runs` times with an empty agent home for a run that has one; a Codex run keeps the session files of its login directory, which persists between runs. Each invocation writes `scenario-runs/<time>-<scenario>-<agent>/` with `conditions.json`, a `run-N/` directory per run (event stream, final working-tree status, read-back facts, `result.json`) and `summary.json`. `scenario-runs/` is ignored by Git. `-login` names where the agent's login lives, outside the repository; each adapter has a default under `~/.cliewen-trial`.

`check` evaluates the stored transcripts again with the current checks and rewrites the results and summary. It runs no agent, which is how a revised check is tried on earlier runs.

## Agents

An agent is reached through an adapter ([IDR-010](../../docs/decisions/IDR-010-agents-are-reached-through-adapters-over-one-interface.md)): a file that registers how its login is supplied, the command that runs it headless, a probe for its version and model, and a parser from its event stream to the neutral transcript the checks read. The runner and the checks name no vendor. `run` without a known `-agent` fails and lists the registered adapters. Both installed agents run as a user who can read their own login inside the container, so the login is the maintainer's own subscription login in a throwaway container and is kept out of image layers, argument lists, result files and logs.

## What a result means

The harness reports, for each run, what the deterministic checks found and groups the runs by that signature. It gives no pass or fail. The checks read an agent's wording and can be wrong on ordinary phrasing; read the transcripts before relying on a signature ([AN-029](../../docs/analysis/AN-029-five-routing-runs-on-claude-code.md)). The container runs Linux and bash, so a result says nothing about PowerShell.

## Scenarios

A scenario is a directory under `scenarios/` holding `prompt.txt`, which is given to the agent on standard input, and `setup.sh`, which builds the fixture repository inside the container, and a registration in `scenarios.go` that states the obligation the scenario exercises, the failure that obligation prevents, and the checks. The prompt is an ordinary request and names none of the rules it tests. The checks live in Go outside the container, so the agent under test cannot read them; the setup script is visible in the container's process arguments.

| Scenario | Obligation exercised | Failure it prevents |
|---|---|---|
| `routing` | State a route before editing | Editing first and saying why afterwards |
| `routing-code` | State a route, and take a new capability through the tracked route | Building a new capability directly |
| `upgrade` | Report a newer release and ask whether to upgrade now or later | Upgrading unasked |
| `brownfield` | Rehearse an extraction report-only before the human authorises | Converting or deleting the source specifications first |
| `greenfield` | Draft a vision as inferred meaning, never accepted, and end with a green `clue validate` | Guessed intent presented as accepted, or a red repository described as set up |
| `tracked-brief` | Propose before implementing, digest, and hand over a brief with no template text | Implementing first, leaving the workspace, placeholders in the brief, a red repository called ready |

Each is traced to a goal or use case in [AN-031](../../docs/analysis/AN-031-four-scenarios-and-a-removal-trial.md), which also states what a machine checks, what a person reads, and where the run is judged. `run` with an unknown `-scenario` lists them. A run is judged at the point the agent stops; no script answers its questions. A long scenario carries its own turn and time limits, and a `post.sh` that runs in the fixture after the agent and writes `key=value` lines the checks read (`greenfield` runs `clue validate`; `tracked-brief` also reads the git history and what its `gh` stand-in recorded); `check` re-scores a stored run from those saved facts. See [AN-034](../../docs/analysis/AN-034-what-the-two-long-scenarios-needed.md).

## Method variants

`-variant NAME` applies a named removal to the fixture repository inside the container after setup, so a scenario can be run with an obligation absent. A variant is a script under `variants/`, its hash is recorded in `conditions.json`, and it runs only inside the throwaway fixture: it never touches this repository's skills or hub. A run that did not finish carries `unfinished` in its signature, and an OpenCode run may write under `/tmp`, which OpenCode would otherwise reject in a run nobody can answer ([AN-033](../../docs/analysis/AN-033-first-comparison-period.md)). A variant that finds nothing to remove fails the run instead of running as a no-op. `baseline` is no removal. `no-routing` removes the hub's routing paragraphs and `no-routing-skill` also removes the `clue-delta` skill.
