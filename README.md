# Cliewen

> A verifiable thread from goal to test — docs as the system of record for agent-driven development.

**Cliewen** (Old English *cliewen*, "ball of thread" — the word that became *clue*) is a methodology and toolchain for agent-driven software development: the documentation corpus is the system-of-record, and the chain **goal → capability → acceptance criterion → test** is mechanically enforced. CLI binary: `clue`.

SDD frameworks document the *change*; Cliewen documents the *system*. Changes are transient deltas digested into the permanent corpus at merge — `git log docs/` is the provenance archive.

## Install

`clue` is a single binary with no runtime dependencies. Either route below ends with `clue` on your `PATH`:

```sh
# From source (needs the Go toolchain) — installs into "$(go env GOPATH)/bin"; make sure that directory is on your PATH:
go install github.com/cliewen/cliewen/cmd/clue@latest

# Or a stamped prebuilt binary from the latest release (linux/darwin/windows × amd64/arm64) — e.g. on linux-amd64:
gh release download --repo cliewen/cliewen --pattern 'clue-*-linux-amd64'
mkdir -p ~/.local/bin && install -m 0755 clue-*-linux-amd64 ~/.local/bin/clue
export PATH="$HOME/.local/bin:$PATH"   # skip if ~/.local/bin is already on your PATH; persist it in your shell profile
```

Either route ends with `clue version` working from any directory. On Windows, download the `clue-*-windows-*.exe` asset for your architecture and place it on `PATH` as `clue.exe`. The repository is currently private while the [public-release campaign](docs/plans/P-003-goes-public.md) completes: collaborators using `go install` additionally need `GOPRIVATE=github.com/cliewen` and git authentication for github.com, while `gh release download` authenticates through `gh auth login`. Anonymous installation is part of the campaign's final visibility milestone.

`clue version` reports the release it was built from — a checkout build (`go build ./cmd/clue`) or an install of an untagged commit reports `dev`. A tagged release (`vX.Y.Z`) builds the cross-platform binaries and stamps each with its version; the agent skills carry the same version, and `clue validate` flags drift between them ([CAP-004](docs/capabilities/CAP-004-ship/README.md), [ADR-011](docs/decisions/ADR-011-version-stamping.md)).

## Quickstart

From nothing to a validated corpus in a few minutes. Prerequisites: `git`, the `clue` binary (install above), and — for the pull-request loop as practiced — an authenticated [`gh`](https://cli.github.com/) CLI (`gh` is the convenient path; the loop works with plain git and any forge).

**1. Initialize.** In your repository (new or existing):

```sh
clue init
clue validate
```

`init` materializes the whole convention in one call: the `docs/` corpus (each folder README explains its record type — start reading at the generated `docs/README.md`), an `AGENTS.md` routing hub for coding agents, the agent skills (`.agents/skills/` plus a `.claude/skills/` mirror), and a CI workflow under `.github/workflows/`. It never replaces existing files — anything you already have (say, your own `AGENTS.md`) is reported and skipped, and a taxonomy README of your own just gains an index block. On a fresh repository `validate` is green immediately; a repo with an existing spec corpus under `docs/` is the brownfield path (see `clue-extract` below). `validate` stays the judge of every change from here on.

**2. Make your first change.** Work follows the change loop in the [`clue-delta`](.agents/skills/clue-delta/skill.md) skill: branch, propose in `/changes/CH-001-your-slug/`, implement against the corpus, digest into `docs/`, then open a PR — checking it first against the [`clue-verify`](.agents/skills/clue-verify/skill.md) checklist. Your coding agent picks both skills up through `AGENTS.md`; a human merges.

**3. Arm the wall.** The generated workflow runs on every PR from the start, but passes with a visible warning until you vendor the pinned release binary it expects — the two commands are in the workflow's comments. Once armed, it runs `clue validate --forbid-changes` on every PR: the same binary that judges your corpus locally judges every merge.

Adopting a repo with an existing spec corpus instead? That is the [`clue-extract`](.agents/skills/clue-extract/skill.md) skill — a one-time transform into `docs/`, run as the repo's first change loop.

## Developing the skills

The five standalone `clue-*` skills are generated from `internal/skills/source/`: shared methodology instructions live under `shared/`, while each lifecycle workflow has its own template under `skills/`. Edit those sources and run:

```sh
go generate ./internal/skills
go test ./...
```

The generator rewrites `.agents/skills/` and the embedded `clue init` copies under `internal/scaffold/templates/skills/`. Tests fail if either generated tree drifts from the canonical rendering.

## Status

Baseline and distribution are complete ([P-001](docs/plans/P-001-elaboration-baseline.md), [P-002](docs/plans/P-002-leaves-home.md)). The repository remains private while the public-release campaign completes readiness, community, guide, and release work before the visibility flip ([P-003](docs/plans/P-003-goes-public.md)). User-visible history lives in [CHANGELOG.md](CHANGELOG.md); each GitHub release body is its version's section there. This repo dogfoods its own conventions from commit one — start reading at [docs/README.md](docs/README.md). Agents: see [AGENTS.md](AGENTS.md).

## License

[Apache 2.0](LICENSE)
