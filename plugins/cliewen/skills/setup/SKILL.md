---
name: setup
description: Install the clue binary and, with the user's agreement, scaffold Cliewen into the current repository. Use when the user asks to set up, install, or start using Cliewen or clue.
---

# Set up Cliewen

Get a working `clue` onto this machine, then stop and ask before changing the user's repository. This skill installs the binary and nothing else — read **What this plugin does not install** before offering to do more.

## 1. Check whether `clue` is already there

Run `clue version`.

- It prints a release version such as `clue 0.8.0` — skip to step 3.
- It prints `dev` — the binary was built from source rather than installed from a release. It works, but it is exempt from the skill-drift check, so say that and offer to replace it with a release build.
- The command is not found — continue to step 2.

## 2. Install `clue`

Pick the command for the host and run it.

```sh
curl -fsSL https://cliewen.dev/install.sh | sh          # macOS and Linux
irm https://cliewen.dev/install.ps1 | iex               # Windows (PowerShell)
go install github.com/cliewen/cliewen/cmd/clue@latest   # any host with Go
```

Do not pass a version. The scripts install the newest release, which is the one whose skills `clue init` writes; naming a version here would pin new users to a release nobody remembers to bump.

The scripts verify the download against the release's `SHA256SUMS` before writing anything, install into a directory the user owns (`~/.local/bin`, or `%LOCALAPPDATA%\Programs\clue` on Windows), and need no administrator rights. If the install directory is not on `PATH`, the script prints the exact line to add — pass it on rather than editing a shell profile on the user's behalf.

The Go route reports `dev` rather than a release version unless a tagged version is installed.

## 3. Confirm the install

Verify the binary with the command for the installation route you used:

- macOS or Linux script: `PATH="${CLUE_INSTALL:-$HOME/.local/bin}:$PATH" clue version`
- Windows script: `& "$env:LOCALAPPDATA\Programs\clue\clue.exe" version`
- Go: `"$(go env GOPATH)/bin/clue" version`

Each of these calls the just-installed binary by its exact path rather than the bare `clue` command, because the host application's inherited `PATH` is independent of the registry or profile change the script just made — a bare `clue version` here can silently resolve to an unrelated binary earlier on `PATH` instead of failing outright. Use the same explicit path for every later `clue` command in this skill. If the script printed an `export PATH=…` line (Unix) or reported adding to the user `PATH` (Windows), pass that along so the user can open a new terminal later; do not edit the profile or registry on their behalf beyond what the script already did. Confirm that a script installation prints a release version. The Go route can print `dev` unless it installed a tagged version. If the relevant command does not produce the expected result, stop and report what it printed. Everything below assumes a working binary, and a wrong answer here is much easier to explain now than after files have been written.

## 4. Ask before scaffolding

`clue init` creates a `/docs` taxonomy, an `AGENTS.md`, the Cliewen skills, and a CI workflow in the current repository. That is a decision about how the project will be run, not a setup step — so ask, and do not run it until the user agrees.

Say what it will do:

- It writes into **this** repository and never replaces an existing file; anything it would have shadowed is reported as skipped.
- It is safe to re-run: a second run refreshes the generated README index blocks and leaves prose alone.
- It writes the six managed Cliewen skills into `.agents/skills/`, stamped with this binary's version. They are committed files, and `clue validate` fails if the binary and those skills ever disagree.

If the user agrees, run `clue init`, then `clue validate`. A fresh scaffold is deliberately **not** green: `init` writes marked bootstraps at `docs/architecture/README.md`, `docs/design/README.md`, and `docs/vision.md`, and `validate` names each of them until it carries concise repository-specific truth. That is the expected first result, not a broken installation. Use the same explicit path from step 3 for both commands:

- macOS or Linux script: `PATH="${CLUE_INSTALL:-$HOME/.local/bin}:$PATH" clue init`, then `PATH="${CLUE_INSTALL:-$HOME/.local/bin}:$PATH" clue validate`
- Windows script: `& "$env:LOCALAPPDATA\Programs\clue\clue.exe" init`, then `& "$env:LOCALAPPDATA\Programs\clue\clue.exe" validate`

Report both outputs. Then draft the two overviews from the repository's own code and documentation — structure, boundaries, actors, and durable technology choices in the architecture overview; cross-cutting flows, interactions, and shared patterns in the design overview — asking the user wherever a material boundary or intent is unclear. Re-run `clue validate` afterwards; that run is the first green one.

If the repository already has a corpus of its own (existing decision records, specifications, requirements), say so and point at the `clue-extract` skill that `clue init` installs, rather than scaffolding over it.

## What this plugin does not install

This plugin ships this skill and nothing else. It deliberately does **not** carry `clue-analysis`, `clue-delta`, `clue-extract`, `clue-plan`, or `clue-verify`.

Those six are not portable helpers. `clue init` writes them into the repository as committed files, each stamped with the version of the binary that wrote it, and `clue validate` fails when the binary and its skills disagree. A plugin's components are copied into a per-user cache instead, where that check cannot see them — one copy spanning every repository the user opens, while the entire point of a stamped skill is that it is pinned per repository.

So if the user wants the lifecycle skills, the answer is always `clue init` in the repository that needs them. Never copy a Cliewen skill into a plugin, a personal skills directory, or a settings file: a second, unversioned set of Cliewen instructions beside the committed set is exactly the drift the version stamp exists to catch, and it would be invisible.

Further reading: [cliewen.dev/plugin](https://cliewen.dev/plugin) for this plugin's boundary, and [cliewen.dev/getting-started](https://cliewen.dev/getting-started) for the full first run.
