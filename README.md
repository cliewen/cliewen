# Cliewen

> Ship agent-written changes without losing the intent.

**Cliewen** is a methodology for repositories where coding agents implement real product changes through pull requests. It keeps requirements, decisions, implementation, and tests connected in Git, and catches missing evidence before merge. **`clue`** is its command-line judge; the **corpus** under `docs/` is the permanent system record that agents maintain with the code.

The name comes from Old English *cliewen*, “ball of thread” — the word that became *clue*. The enforced thread is **goal → capability → acceptance criterion → test**.

SDD frameworks document the *change*; Cliewen documents the *system*. Changes are transient deltas digested into the permanent corpus at merge — `git log docs/` is the provenance archive.

## Install

`clue` is a single binary with no runtime dependencies. Open the [latest release](https://github.com/cliewen/cliewen/releases/latest) and download `SHA256SUMS` plus the asset for your machine into an otherwise empty directory:

| Machine | Release asset |
|---|---|
| Windows x64 / ARM64 | `clue-<version>-windows-amd64.exe` / `clue-<version>-windows-arm64.exe` |
| macOS Intel / Apple silicon | `clue-<version>-darwin-amd64` / `clue-<version>-darwin-arm64` |
| Linux x86-64 / ARM64 | `clue-<version>-linux-amd64` / `clue-<version>-linux-arm64` |

Then:

1. Verify the downloaded binary's SHA-256 matches its line in `SHA256SUMS`. Use `Get-FileHash <asset> -Algorithm SHA256` on Windows, `shasum -a 256 <asset>` on macOS, or `sha256sum <asset>` on Linux.
2. Rename the binary to `clue.exe` on Windows or `clue` on macOS and Linux. On macOS and Linux, also make it executable with `chmod +x clue`.
3. Move it into a directory on your user `PATH`. On Windows, a folder such as `%LOCALAPPDATA%\Programs\clue` works once added through "Edit environment variables for your account." On macOS and Linux, `~/.local/bin` is a common choice.
4. Open a new terminal and run `clue version`. It should match the release you downloaded.

The macOS binaries are unsigned and not notarized. If macOS blocks `clue`, first confirm the checksum matches, try `clue version` once, then open **System Settings → Privacy & Security** and click **Open Anyway**. Apple documents this exception in [Open a Mac app from an unknown developer](https://support.apple.com/guide/mac-help/open-a-mac-app-from-an-unknown-developer-mh40616/mac).

The [installation guide](https://cliewen.dev/getting-started#_1-install-clue) has the same short path with a little more context, but it is not required for the quickstart.

To install from source instead, use the Go toolchain:

```sh
go install github.com/cliewen/cliewen/cmd/clue@latest
clue version
```

`clue version` reports the release it was built from — a checkout build (`go build ./cmd/clue`) or an install of an untagged commit reports `dev`. A tagged release (`vX.Y.Z`) builds the cross-platform binaries and stamps each with its version; the agent skills carry the same version, and `clue validate` flags drift between them ([CAP-004](docs/capabilities/CAP-004-ship/README.md), [ADR-011](docs/decisions/ADR-011-version-stamping.md)).

## Quickstart

From nothing to a validated corpus in a few minutes, without touching an existing repository. Prerequisites: `git` and the `clue` binary (install above). An authenticated [`gh`](https://cli.github.com/) CLI is convenient later for the pull-request loop; Cliewen itself works with plain Git and any forge.

**1. Initialize a disposable repository.**

```sh
mkdir cliewen-demo
cd cliewen-demo
git init
clue init
clue validate
```

`init` materializes the whole convention in one call: the `docs/` corpus, an `AGENTS.md` routing hub, agent skills, and a GitHub workflow. On a fresh repository `validate` is green immediately. Continue with the [safe demo](https://cliewen.dev/getting-started#_3-see-clue-catch-a-broken-thread) to activate an acceptance criterion without evidence and watch `clue` name the missing test; remove `cliewen-demo` afterwards.

**2. Make your first Cliewen change.** The generated `AGENTS.md` first keeps unrelated editorial work out of Cliewen: a plain change uses an ordinary branch, relevant checks, a PR, and human merge. Work whose meaning belongs in Cliewen follows the change loop in [`clue-delta`](.agents/skills/clue-delta/skill.md): branch, propose in `/changes/CH-001-your-slug/`, implement against the corpus, digest into `docs/`, then run the pre-merge checks and automatic agentic review in [`clue-verify`](.agents/skills/clue-verify/skill.md) before opening a PR. Your coding agent loads the corpus and skills only when they are relevant.

**3. [Arm the wall](https://cliewen.dev/ci-wall).** The generated workflow runs as the stable required check. Only Markdown outside protected corpus, policy, configuration, and methodology paths is eligible to pass through without Cliewen; every non-Markdown, protected, mixed, empty, or unreadable diff fails closed to the corpus wall. Until you vendor the pinned release binary it expects, those Cliewen changes pass with a visible warning. Once armed, the workflow verifies the pinned binary and runs `clue validate --forbid-changes`; the guide shows how to require that check and prove a failing pull request is blocked.

Adopting a repo with an existing spec corpus instead? That is the [`clue-extract`](.agents/skills/clue-extract/skill.md) skill — a one-time transform into `docs/`, run as the repo's first change loop.

## Public guide

The handwritten [Cliewen guide](https://cliewen.dev/) explains the methodology, corpus taxonomy, change loop, and skills for newcomers who are not yet inside a Cliewen repository. Its [source](guide/index.md) builds with strict dead-link checking in CI and deploys from `main` through GitHub Pages.

## Developing the skills

The five standalone `clue-*` skills are generated from `internal/skills/source/`: shared methodology instructions live under `shared/`, while each lifecycle workflow has its own template under `skills/`. Edit those sources and run:

```sh
go generate ./internal/skills
go test ./...
```

The generator rewrites `.agents/skills/` and the embedded `clue init` copies under `internal/scaffold/templates/skills/`. Tests fail if either generated tree drifts from the canonical rendering.

## Status

Baseline, distribution, and the public launch are complete ([P-001](docs/plans/P-001-elaboration-baseline.md), [P-002](docs/plans/P-002-leaves-home.md), [P-003](docs/plans/P-003-goes-public.md)). [P-004](docs/plans/P-004-first-try.md) is improving the first safe trial, enforced adoption path, and operating guidance. User-visible history lives in [CHANGELOG.md](CHANGELOG.md); each GitHub release body is its version's section there. This repo dogfoods its own conventions from commit one — start reading at [docs/README.md](docs/README.md). Agents: see [AGENTS.md](AGENTS.md).

## License

[Apache 2.0](LICENSE)
