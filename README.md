# Cliewen

> A verifiable thread from goal to test — docs as the system of record for agent-driven development.

**Cliewen** (Old English *cliewen*, "ball of thread" — the word that became *clue*) is a methodology and toolchain for agent-driven software development: the documentation corpus is the system-of-record, and the chain **goal → capability → acceptance criterion → test** is mechanically enforced. CLI binary: `clue`.

SDD frameworks document the *change*; Cliewen documents the *system*. Changes are transient deltas digested into the permanent corpus at merge — `git log docs/` is the provenance archive.

## Install

`clue` is a single Go binary. Install from source, or download a stamped prebuilt binary from a release:

```sh
# From source — a release install stamps itself from Go's build info:
go install github.com/cliewen/cliewen/cmd/clue@latest

# A stamped prebuilt binary for your platform (linux/darwin/windows × amd64/arm64):
gh release download vX.Y.Z --repo cliewen/cliewen --pattern 'clue-*-<os>-<arch>*'
```

While the repo is private, `go install` additionally needs `GOPRIVATE=github.com/cliewen` and git authentication for github.com; `gh release download` authenticates through `gh auth login`.

`clue version` reports the release it was built from — a checkout build (`go build ./cmd/clue`) or an install of an untagged commit reports `dev`. A tagged release (`vX.Y.Z`) builds the cross-platform binaries and stamps each with its version; the agent skills carry the same version, and `clue validate` flags drift between them ([CAP-004](docs/capabilities/CAP-004-ship/README.md), [ADR-011](docs/decisions/ADR-011-version-stamping.md)).

## Quickstart

From nothing to a validated corpus in a few minutes. Prerequisites: `git`, the `clue` binary (install above), and — for the pull-request loop as practiced — an authenticated [`gh`](https://cli.github.com/) CLI (`gh` is the convenient path; the loop works with plain git and any forge).

**1. Initialize.** In your repository (new or existing):

```sh
clue init
clue validate
```

`init` materializes the whole convention in one call: the `docs/` corpus (each folder README explains its record type — start reading at the generated `docs/README.md`), an `AGENTS.md` routing hub for coding agents, the agent skills (`.agents/skills/` plus a `.claude/skills/` mirror), and a CI workflow template under `.github/workflows/`. It never overwrites existing files — anything you already have (say, your own `AGENTS.md`) is reported and skipped. `validate` is green immediately; it stays the judge of every change from here on.

**2. Make your first change.** Work follows the change loop in the [`clue-delta`](.agents/skills/clue-delta/skill.md) skill: branch, propose in `/changes/CH-001-your-slug/`, implement against the corpus, digest into `docs/`, then open a PR — checking it first against the [`clue-verify`](.agents/skills/clue-verify/skill.md) checklist. Your coding agent picks both skills up through `AGENTS.md`; a human merges.

**3. Put up the wall.** The generated workflow runs `clue validate --forbid-changes` on every PR — its comments explain how to vendor the pinned release binary it expects. From then on, the same binary that judges your corpus locally judges every merge.

Adopting a repo with an existing spec corpus instead? That is the [`clue-extract`](.agents/skills/clue-extract/skill.md) skill — a one-time transform into `docs/`, run as the repo's first change loop.

## Status

Baseline complete ([P-001](docs/plans/P-001-elaboration-baseline.md)); distribution and greenfield bootstrap under way ([P-002](docs/plans/P-002-leaves-home.md)). User-visible history lives in [CHANGELOG.md](CHANGELOG.md); each GitHub release body is its version's section there. This repo dogfoods its own conventions from commit one — start reading at [docs/README.md](docs/README.md). Agents: see [AGENTS.md](AGENTS.md).

## License

[Apache 2.0](LICENSE)
