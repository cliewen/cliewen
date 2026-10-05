---
type: contributor-guide
title: Scenario trials
---

# Scenario trials

Repository tooling for [P-024](../../docs/plans/P-024-repeatable-scenario-trials.md): it runs a scenario on an agent in a container and records what happened. It is not part of `clue`, is not shipped to adopters, and nothing starts it automatically. See [IDR-009](../../docs/decisions/IDR-009-scenario-trial-harness-is-go-repository-tooling.md).

## Prerequisites

- Docker with Linux containers, a Go toolchain, and Git on the host.
- A Claude Code subscription token from `claude setup-token`, saved as a single line in `~/.cliewen-trial/token` (outside the repository). Write it with an editor or `Set-Content -NoNewline -Path <file> -Value <token>`; never pass it where its text could become a file name or reach a log. Revoke and replace it if it ever appears in a transcript or a repository.

## Run

```text
go run ./tools/scenario-trials run -runs 5
go run ./tools/scenario-trials check scenario-runs/<invocation>
```

`run` builds a Linux `clue` from `-clue-commit` (default `HEAD`) in a throwaway worktree, builds a container image keyed to that binary, and runs the scenario `-runs` times with an empty agent home. Each invocation writes `scenario-runs/<time>-<scenario>-<agent>/` with `conditions.json`, a `run-N/` directory per run (event stream, final working-tree status, read-back facts, `result.json`) and `summary.json`. `scenario-runs/` is ignored by Git.

`check` evaluates the stored transcripts again with the current checks and rewrites the results and summary. It runs no agent, which is how a revised check is tried on earlier runs.

## What a result means

The harness reports, for each run, what the deterministic checks found and groups the runs by that signature. It gives no pass or fail. The checks read an agent's wording and can be wrong on ordinary phrasing; read the transcripts before relying on a signature ([AN-029](../../docs/analysis/AN-029-five-routing-runs-on-claude-code.md)). The container runs Linux and bash, so a result says nothing about PowerShell.

## Scenarios

A scenario is a directory under `scenarios/` holding `prompt.txt`, which is given to the agent, and `setup.sh`, which builds the fixture repository inside the container. Only the `routing` scenario and the `claude` adapter exist so far.
