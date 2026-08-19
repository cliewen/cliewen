# Cliewen

> Evidence-backed Intent Engineering for coding agents.

**Cliewen** is a methodology for repositories where coding agents implement real product changes through pull requests. It keeps requirements, decisions, implementation, and acceptance evidence connected in Git, and catches missing evidence before merge. **`clue`** is its command-line judge; the **corpus** under `docs/` is the permanent system record that agents maintain with the code.

*Evidence-backed Intent Engineering* is Cliewen's own description of its approach, not an established industry label: human intent is recorded as durable goals, capabilities, decisions, constraints, and acceptance criteria; every active criterion declares the evidence that accepts it — a classified executable reference or explicitly identified human verification — and tooling checks mechanically that the chain is complete before a human merges. It makes the connection between intent and acceptance evidence explicit, reviewable, and checkable. It does not prove your software satisfies your intent: `clue` does not execute tests, judge whether a test asserts the right behavior, or know whether the intent was right in the first place. Semantic acceptance stays with review and the human at the merge gate.

The name comes from Old English *cliewen*, “ball of thread” — the word that became *clue*. The enforced thread is **goal → capability → acceptance criterion → acceptance evidence**: classified executable references for machine-proven criteria or the pull request acceptance brief for genuine Human proof.

SDD frameworks document the *change*; Cliewen documents the *system*. Changes are transient deltas digested into the permanent corpus at merge — full Cliewen changes use a human-controlled merge commit so their proposal, implementation, digest, and durable corpus history remain reachable from `main`; reachable Git history is the provenance archive.

## Install

`clue` is a single binary with no runtime dependencies. One command:

```sh
curl -fsSL https://cliewen.dev/install.sh | sh          # macOS and Linux
irm https://cliewen.dev/install.ps1 | iex               # Windows (PowerShell)
go install github.com/cliewen/cliewen/cmd/clue@latest   # any host with Go
```

Open a new terminal, then `clue version` should print the release you installed. The script detects your platform, downloads the matching binary, and verifies it against the release's `SHA256SUMS` before installing — a mismatch stops with nothing written. It needs no administrator rights; `CLUE_INSTALL` and `CLUE_VERSION` override the target directory and release. Upgrading is the same command again, which moves the binary only; in an adopted repository, preview and apply the coordinated corpus upgrade with `clue migrate` as described in the operations guide.

Inside Claude Code, `/plugin marketplace add cliewen/cliewen` then `/plugin install cliewen@cliewen` runs the same install for you and offers to scaffold. It ships no lifecycle skills — `clue init` writes those into the repository, stamped with the binary's version; the [plugin page](https://cliewen.dev/plugin) explains why that boundary exists.

### Download a binary instead

To install by hand, open the [latest release](https://github.com/cliewen/cliewen/releases/latest) and download `SHA256SUMS` plus the asset for your machine into an otherwise empty directory:

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

The macOS binaries are unsigned and not notarized, so a binary downloaded through a browser can be blocked by Gatekeeper. First confirm the checksum matches, try `clue version` once, then open **System Settings → Privacy & Security** and click **Open Anyway**. Apple documents this exception in [Open a Mac app from an unknown developer](https://support.apple.com/guide/mac-help/open-a-mac-app-from-an-unknown-developer-mh40616/mac). The install script is not affected — a download made outside the browser carries no quarantine attribute.

The [installation guide](https://cliewen.dev/install) has the same short path with a little more context, but it is not required for the quickstart.

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

`init` materializes the whole convention in one call: the `docs/` corpus, an `AGENTS.md` routing hub, agent skills, and a GitHub workflow. On a fresh repository `validate` is green immediately. Continue with the [safe demo](https://cliewen.dev/getting-started#_2-see-clue-catch-a-broken-thread) to activate an acceptance criterion without evidence and watch `clue` name the missing test; remove `cliewen-demo` afterwards.

**2. Make your first change.** The generated `AGENTS.md` recommends simple when the accepted contract stays unchanged and full when acceptance-criterion, capability, decision, policy, plan-promise, methodology, or uncovered-behavior meaning changes. Simple work uses relevant checks and the integration mechanism the user authorizes and the repository permits. A chosen full loop uses `clue context <id>` to load the relevant outgoing-link slice, branches from accepted `main`, commits `/changes/CH-001-your-slug/`, and opens a draft PR with that proposal so work is durable from first publication; it then implements against the corpus, digests into `docs/`, and runs the pre-ready checks and automatic agentic review in [`clue-verify`](.agents/skills/clue-verify/skill.md). Your coding agent loads broader corpus context only when the task discovers that it needs it.

**3. [Arm the wall](https://cliewen.dev/ci-wall).** The generated caller invokes Cliewen's upstream reusable validation workflow at an immutable reference and exposes only runner and binary-source choices. Changed surfaces select relevant corpus checks independently of route, while proposal history selects the full-loop acceptance-brief gate unless the authored head commit carries the complete user-override trailers. Until you vendor the pinned release binary or choose the verified release source, corpus-relevant changes pass with a visible warning. Once armed, the upstream workflow runs `clue validate --forbid-changes`; configure the protected branch for merge commits only if you use the supported full loop, then use the guide's probe to prove both the history boundary and the failing pull-request block.

Adopting a repo with an existing spec corpus instead? That is the [`clue-extract`](.agents/skills/clue-extract/skill.md) skill — a one-time transform into `docs/`, run as the repo's first change loop. Its proposed change starts with a report-only rehearsal; explicit human direction is required before it mutates the repository.

## Public guide

The handwritten [Cliewen guide](https://cliewen.dev/) explains the methodology, corpus taxonomy, change loop, and skills for newcomers who are not yet inside a Cliewen repository. Its [source](guide/index.md) builds with strict dead-link checking in CI and deploys from `main` through GitHub Pages.

If you are evaluating Cliewen, read it in this order: [What is Cliewen?](https://cliewen.dev/what-is-cliewen) for the approach and where it came from, [The design of Cliewen](https://cliewen.dev/design) for why acceptance and evidence are central and what the method does not solve, [The verifiable thread](https://cliewen.dev/methodology) for the concrete goal-to-evidence model, and [Get started](https://cliewen.dev/getting-started) for a disposable example that deliberately breaks the thread so you can watch `clue` catch it.

## Developing the skills

The six standalone `clue-*` skills are generated from `internal/skills/source/`: shared methodology instructions live under `shared/`, while each lifecycle workflow has its own template under `skills/`. Edit those sources and run:

```sh
go generate ./internal/skills
go test ./...
```

The generator rewrites `.agents/skills/` and the embedded `clue init` copies under `internal/scaffold/templates/skills/`. Tests fail if either generated tree drifts from the canonical rendering.

## Status

Baseline, distribution, the public launch, and the first-try campaign are complete ([P-001](docs/plans/P-001-elaboration-baseline.md), [P-002](docs/plans/P-002-leaves-home.md), [P-003](docs/plans/P-003-goes-public.md), [P-004](docs/plans/P-004-first-try.md)). User-visible history lives in [CHANGELOG.md](CHANGELOG.md); each GitHub release body is its version's section there. This repo dogfoods its own conventions from commit one — start reading at [docs/README.md](docs/README.md). Agents: see [AGENTS.md](AGENTS.md).

## License

[Apache 2.0](LICENSE)
