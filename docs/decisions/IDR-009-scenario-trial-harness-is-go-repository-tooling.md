---
id: IDR-009
type: decision
status: verified
links: [P-024, G-025]
title: The scenario-trial harness is Go repository tooling under tools/
author: agent
accepted-by: Flemming N. Larsen (2026-10-05, conversation)
---

# IDR-009 — The scenario-trial harness is Go repository tooling under tools/

## Context

[P-024](../plans/P-024-repeatable-scenario-trials.md) needs a harness that runs a scenario on a chosen agent in a container and records what happened. The maintainer works on Windows, and the repository already carries Go tooling under `tools/` that is not part of the `clue` executable.

## Decision

The harness is a Go program in `tools/scenario-trials`, in the repository's own module. It is repository tooling: it is not shipped to adopters, `clue` does not call it, and no workflow starts a trial run. A run starts only when the maintainer runs it. Run records go to `scenario-runs/`, which Git ignores.

Go keeps one program for Windows and Linux hosts instead of a bash script and a PowerShell script that drift apart. CI's `./...` compiles and tests the package like any other, which is neither a run nor a dependency of one.
