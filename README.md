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

Run the matching block from that download directory. It verifies the release checksum, installs `clue` under your user account, adds it to the current shell, and prints the installed version.

### Windows PowerShell

```powershell
$Architecture = switch ([System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture) {
    ([System.Runtime.InteropServices.Architecture]::X64) { "amd64" }
    ([System.Runtime.InteropServices.Architecture]::Arm64) { "arm64" }
    default { throw "Cliewen has no Windows asset for this architecture" }
}
$Asset = Get-ChildItem -File "clue-*-windows-$Architecture.exe" | Select-Object -First 1
if ($null -eq $Asset) { throw "The Cliewen release asset is not in this directory" }
$ChecksumLine = Get-Content .\SHA256SUMS | Where-Object { $_ -match ("\s" + [regex]::Escape($Asset.Name) + "$") }
if ($null -eq $ChecksumLine) { throw "SHA256SUMS does not name $($Asset.Name)" }
$Expected = (($ChecksumLine -split "\s+")[0]).ToUpperInvariant()
$Actual = (Get-FileHash -Algorithm SHA256 -LiteralPath $Asset.FullName).Hash
if ($Actual -ne $Expected) { throw "Checksum mismatch: do not run this binary" }

$InstallDir = Join-Path $env:LOCALAPPDATA "Programs\clue"
New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
Copy-Item -LiteralPath $Asset.FullName -Destination (Join-Path $InstallDir "clue.exe")
$UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
if (@($UserPath -split ";") -notcontains $InstallDir) {
    $NewUserPath = if ([string]::IsNullOrWhiteSpace($UserPath)) { $InstallDir } else { "$UserPath;$InstallDir" }
    [Environment]::SetEnvironmentVariable("Path", $NewUserPath, "User")
}
$env:Path = "$InstallDir;$env:Path"
clue version
```

### macOS

```sh
case "$(uname -m)" in
  x86_64) arch=amd64 ;;
  arm64) arch=arm64 ;;
  *) echo "Cliewen has no macOS asset for this architecture" >&2; exit 1 ;;
esac
asset=$(find . -maxdepth 1 -type f -name "clue-*-darwin-$arch" -print -quit)
test -n "$asset" || { echo "The Cliewen release asset is not in this directory" >&2; exit 1; }
name=$(basename "$asset")
expected=$(awk -v file="$name" '$2 == file { print $1 }' SHA256SUMS)
actual=$(shasum -a 256 "$asset" | awk '{ print $1 }')
test -n "$expected" && test "$actual" = "$expected" || { echo "Checksum mismatch: do not run this binary" >&2; exit 1; }
mkdir -p "$HOME/.local/bin"
install -m 0755 "$asset" "$HOME/.local/bin/clue"
export PATH="$HOME/.local/bin:$PATH"
clue version
```

### Linux

```sh
case "$(uname -m)" in
  x86_64) arch=amd64 ;;
  aarch64 | arm64) arch=arm64 ;;
  *) echo "Cliewen has no Linux asset for this architecture" >&2; exit 1 ;;
esac
asset=$(find . -maxdepth 1 -type f -name "clue-*-linux-$arch" -print -quit)
test -n "$asset" || { echo "The Cliewen release asset is not in this directory" >&2; exit 1; }
name=$(basename "$asset")
expected=$(awk -v file="$name" '$2 == file { print $1 }' SHA256SUMS)
test -n "$expected" || { echo "SHA256SUMS does not name $name" >&2; exit 1; }
printf '%s  %s\n' "$expected" "$asset" | sha256sum -c -
mkdir -p "$HOME/.local/bin"
install -m 0755 "$asset" "$HOME/.local/bin/clue"
export PATH="$HOME/.local/bin:$PATH"
clue version
```

Persist `$HOME/.local/bin` in your shell's `PATH` on macOS or Linux if it is not already there. `clue version` should match the release you downloaded. The [installation guide](https://cliewen.github.io/cliewen/getting-started#_1-install-clue) explains the same path in more detail but is not required for the quickstart.

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

`init` materializes the whole convention in one call: the `docs/` corpus, an `AGENTS.md` routing hub, agent skills, and a GitHub workflow. On a fresh repository `validate` is green immediately. Continue with the [safe demo](https://cliewen.github.io/cliewen/getting-started#_3-see-clue-catch-a-broken-thread) to activate an acceptance criterion without evidence and watch `clue` name the missing test; remove `cliewen-demo` afterwards.

**2. Make your first Cliewen change.** The generated `AGENTS.md` first keeps unrelated editorial work out of Cliewen: a plain change uses an ordinary branch, relevant checks, a PR, and human merge. Work whose meaning belongs in Cliewen follows the change loop in [`clue-delta`](.agents/skills/clue-delta/skill.md): branch, propose in `/changes/CH-001-your-slug/`, implement against the corpus, digest into `docs/`, then run the pre-merge checks and automatic agentic review in [`clue-verify`](.agents/skills/clue-verify/skill.md) before opening a PR. Your coding agent loads the corpus and skills only when they are relevant.

**3. Arm the wall.** The generated workflow runs as the stable required check. Only Markdown outside protected corpus, policy, configuration, and methodology paths is eligible to pass through without Cliewen; every non-Markdown, protected, mixed, empty, or unreadable diff fails closed to the corpus wall. Until you vendor the pinned release binary it expects, those Cliewen changes pass with a visible warning — the two arming commands are in the workflow's comments. Once armed, it runs `clue validate --forbid-changes` for Cliewen changes.

Adopting a repo with an existing spec corpus instead? That is the [`clue-extract`](.agents/skills/clue-extract/skill.md) skill — a one-time transform into `docs/`, run as the repo's first change loop.

## Public guide

The handwritten [Cliewen guide](https://cliewen.github.io/cliewen/) explains the methodology, corpus taxonomy, change loop, and skills for newcomers who are not yet inside a Cliewen repository. Its [source](guide/index.md) builds with strict dead-link checking in CI and deploys from `main` through GitHub Pages.

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
