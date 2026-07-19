# Get started

The shortest path is: install one binary, initialize the convention in your repository, and run the judge.

## Prerequisites

- **Required:** Git (`git`). The next step installs the `clue` binary.
- **Required for source installation:** the Go toolchain (`go`). You do not need Go when using a prebuilt `clue` release.
- **Recommended for GitHub:** an authenticated [GitHub CLI](https://cli.github.com/) (`gh`) for the pull-request loop. Cliewen itself works with plain Git and any forge.

Node.js (`node`) and npm (`npm`) are not required to use Cliewen. They are only needed to build this guide or contribute to Cliewen itself.

## 1. Install `clue`

Install from source:

```sh
go install github.com/cliewen/cliewen/cmd/clue@latest
clue version
```

The binary is installed into `$(go env GOPATH)/bin`; add that directory to `PATH` if your shell cannot find `clue`.

Alternatively, download the release asset for your operating system and architecture from [GitHub Releases](https://github.com/cliewen/cliewen/releases), verify it against `SHA256SUMS`, and place the executable on `PATH`.

## 2. Initialize your repository

Run these commands at the repository root:

```sh
clue init
clue validate
```

`clue init` creates the documentation taxonomy, agent routing file, lifecycle skills, and a CI workflow. It never replaces an existing file. Re-running it refreshes generated index blocks and otherwise leaves your prose alone.

A fresh repository validates immediately. In a repository with existing specifications, use the brownfield `clue-extract` skill to map the existing truth into the corpus as the first change.

## 3. Read the result

Start at `docs/README.md`. It explains the artifact graph and links to each taxonomy. `AGENTS.md` tells a coding agent which lifecycle skill to load for analysis, planning, implementation, extraction, or verification.

You do not need to learn the whole corpus before making a change. The [change loop](./change-loop) introduces each artifact when it becomes useful.

## 4. Make the wall real

The generated CI workflow runs immediately but stays visibly unarmed until you vendor the pinned `clue` binary described in its comments. Once armed, CI runs `clue validate --forbid-changes`; a pull request cannot pass with an undigested `/changes` workspace.

That is your first complete thread: a repository-local intent model, the same deterministic judge on a developer machine and in CI, and a human-controlled merge boundary.
