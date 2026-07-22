# Get started

This path installs one binary and lets you see Cliewen reject a broken intent-to-test thread in a disposable Git repository. It does not touch an existing project, and you can remove the whole experiment by deleting one directory.

## Prerequisites

- **Required:** Git (`git`).
- **Required for the default installation:** permission to add one directory to your user `PATH`.
- **Optional:** the Go toolchain, only if you choose the source-installation route.
- **Recommended later for GitHub:** an authenticated [GitHub CLI](https://cli.github.com/) (`gh`) for the pull-request loop. Cliewen itself works with plain Git and any forge.

Node.js and npm are needed only to build this guide or contribute to Cliewen itself.

## 1. Install `clue`

Open the [latest Cliewen release](https://github.com/cliewen/cliewen/releases/latest) and download `SHA256SUMS` plus the binary for your machine into an otherwise empty download directory:

| Machine | Release asset |
|---|---|
| Windows x64 | `clue-<version>-windows-amd64.exe` |
| Windows ARM64 | `clue-<version>-windows-arm64.exe` |
| macOS on Intel | `clue-<version>-darwin-amd64` |
| macOS on Apple silicon | `clue-<version>-darwin-arm64` |
| Linux x86-64 | `clue-<version>-linux-amd64` |
| Linux ARM64 | `clue-<version>-linux-arm64` |

Run the matching block from that download directory. Each block selects the downloaded asset, verifies it against the release checksum, installs it under your user account, and makes it available in the current shell.

::: code-group

```powershell [Windows PowerShell]
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

```sh [macOS]
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

```sh [Linux]
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

:::

Persist `$HOME/.local/bin` in your shell's `PATH` on macOS or Linux if it is not already there. `clue version` prints the version encoded in the asset, for example `clue 0.5.1`; it should match the release you downloaded.

### Install from source instead

If you already have Go and prefer a source installation:

```sh
go install github.com/cliewen/cliewen/cmd/clue@latest
clue version
```

Go installs the executable under `$(go env GOPATH)/bin`. Add that directory to `PATH` if your shell cannot find `clue`. A tagged install reports the release version; an untagged checkout build reports `dev`.

## 2. Initialize a disposable repository

Create an empty directory instead of experimenting in an existing project:

```sh
mkdir cliewen-demo
cd cliewen-demo
git init
clue init
clue validate
```

The current release reports every created file, then ends like this:

```text
clue init: 25 created, 0 skipped, 0 index block(s) regenerated
next: run `clue validate` — green on a fresh scaffold; then read docs/README.md
clue validate: OK (2 artifacts, 1 agent-enforced constraint(s) awaiting machine checks)
```

The exact count can grow in a future release. The important result is the final `OK`. The top-level tree is:

```text
cliewen-demo/
├── .agents/skills/       versioned Cliewen skills
├── .claude/skills/       Claude Code mirror of those skills
├── .github/workflows/    the CI wall template
├── docs/                 the permanent corpus
└── AGENTS.md             routing instructions for coding agents
```

`clue init` copies defaults but does not take ownership of your repository. You and your agent own the corpus prose and repository-specific instructions. `clue scaffold` and repeated `clue init` regenerate only the marked README index blocks; existing files are otherwise skipped, never replaced. The copied skills and workflow are versioned repository files, not background-managed services.

## 3. See `clue` catch a broken thread

Add a tiny goal and capability while keeping its acceptance criterion in `draft`. Create these three files:

::: code-group

```powershell [Windows PowerShell]
New-Item -ItemType Directory -Force docs/capabilities/CAP-001-greeting | Out-Null
```

```sh [macOS and Linux]
mkdir -p docs/capabilities/CAP-001-greeting
```

:::

::: code-group

```markdown [docs/goals/G-001-demo.md]
---
id: G-001
type: goal
status: accepted
links: []
title: A greeting is available
---

# G-001 — A greeting is available

A user wants a greeting they can request by name.
```

```markdown [docs/capabilities/CAP-001-greeting/README.md]
---
id: CAP-001
type: capability
status: active
links: [G-001]
title: Return a greeting
goal: G-001
---

# CAP-001 — Return a greeting

The system returns a greeting for a supplied name.
```

````markdown [docs/capabilities/CAP-001-greeting/criteria.md]
---
id: CAP-001-criteria
type: criteria
status: draft
links: [CAP-001]
title: Acceptance criteria for greetings
---

```gherkin
Feature: Return a greeting

  @AC-001
  Scenario: Greet a supplied name
    Given the name "Ada"
    When a greeting is requested
    Then the result is "Hello, Ada"
```
````

:::

Regenerate the two taxonomy indexes and validate:

```sh
clue scaffold
clue validate
```

The draft criterion does not claim verified behavior, so the result is green:

```text
clue scaffold: 2 index block(s) regenerated
clue validate: OK (5 artifacts, 1 agent-enforced constraint(s) awaiting machine checks)
```

Now change only `status: draft` to `status: active` in `criteria.md` and run `clue validate` again. The command exits with status 1 and names the broken edge:

```text
docs/capabilities/CAP-001-greeting/criteria.md: AC-001 has no test (convention per ADR-005/ADR-009: a Go test named TestAC001_… or a framework tag "AC_001")
clue validate: 1 issue(s)
```

That is the product's job: an active promise cannot silently lose all executable evidence. To return this demo to green, set the criterion back to `draft`. In a real change, keep it draft until the implementation has focused positive and negative tests carrying the AC reference, then activate it.

`clue validate` currently detects whether at least one supported test reference exists; it does not distinguish or count the positive and negative pair. The Cliewen change loop and human review enforce that stronger requirement. `clue` also does not run the tests—your normal test runner remains responsible for whether they pass.

## 4. Remove the experiment or continue

The experiment changed only `cliewen-demo`. Leave that directory, then delete it with your file manager or normal directory-removal command to undo the entire trial. Removing the separately installed `clue` binary is not required.

If the failure made sense and you want to use Cliewen, start again in a new project or read the [greenfield and brownfield guide](./adoption) before initializing an existing repository. The [change loop](./change-loop) shows how an agent carries a real request from proposal to human merge.
