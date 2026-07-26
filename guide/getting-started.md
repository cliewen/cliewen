# Get started

This path installs one binary and lets you see Cliewen reject a broken intent-to-test thread in a disposable Git repository. It does not touch an existing project, and you can remove the whole experiment by deleting one directory.

## Prerequisites

- **Required:** Git (`git`).
- **Required for the install script:** `curl` or `wget` and `sha256sum` or `shasum` on macOS and Linux; PowerShell on Windows. All are normally present.
- **Required:** permission to add one directory to your user `PATH`. No administrator rights are needed.
- **Optional:** the Go toolchain, if you prefer installing from source.
- **Recommended later for GitHub:** an authenticated [GitHub CLI](https://cli.github.com/) (`gh`) for the pull-request loop. Cliewen itself works with plain Git and any forge.

Node.js and npm are needed only to build this guide or contribute to Cliewen itself.

## 1. Install `clue`

One command, on any supported machine:

::: code-group

```sh [macOS and Linux]
curl -fsSL https://cliewen.dev/install.sh | sh
```

```powershell [Windows]
irm https://cliewen.dev/install.ps1 | iex
```

```sh [Any host with Go]
go install github.com/cliewen/cliewen/cmd/clue@latest
```

:::

Then open a new terminal and run `clue version` — a `PATH` change does not reach an already-running shell. It should print the version you installed, for example `clue 0.7.0`.

The script detects your operating system and architecture, downloads the matching release binary, and **verifies it against the release's `SHA256SUMS` before installing anything**; a mismatch stops with nothing written. It installs to `~/.local/bin` (macOS and Linux) or `%LOCALAPPDATA%\Programs\clue` (Windows) and needs no administrator rights. Set `CLUE_INSTALL` to choose a different directory, or `CLUE_VERSION` to pin a release. You are piping a script into a shell, so read it first if you prefer: [install.sh](https://cliewen.dev/install.sh) and [install.ps1](https://cliewen.dev/install.ps1) are short, and the manual steps below do the same work by hand.

The Go route installs under `$(go env GOPATH)/bin`, which you may need to add to `PATH` yourself; it reports `dev` rather than a release version unless you install a tagged version.

Upgrading later means re-running the same command. That moves the binary only — for a repository already using Cliewen it is half an upgrade, and [Operate safely](./operations) explains why the resulting drift report is the check working.

### Download a binary instead

To install by hand — or on a machine where neither script can run — open the [latest Cliewen release](https://github.com/cliewen/cliewen/releases/latest) and download `SHA256SUMS` plus the binary for your machine into an otherwise empty download directory:

| Machine | Release asset |
|---|---|
| Windows x64 | `clue-<version>-windows-amd64.exe` |
| Windows ARM64 | `clue-<version>-windows-arm64.exe` |
| macOS on Intel | `clue-<version>-darwin-amd64` |
| macOS on Apple silicon | `clue-<version>-darwin-arm64` |
| Linux x86-64 | `clue-<version>-linux-amd64` |
| Linux ARM64 | `clue-<version>-linux-arm64` |

Then:

1. Verify the downloaded binary's SHA-256 matches its line in `SHA256SUMS`.

| System | Built-in checksum command |
|---|---|
| Windows PowerShell | `Get-FileHash <asset> -Algorithm SHA256` |
| macOS | `shasum -a 256 <asset>` |
| Linux | `sha256sum <asset>` |

2. Rename the binary to `clue.exe` on Windows or `clue` on macOS and Linux. On macOS and Linux, also make it executable with `chmod +x clue`.
3. Move it into a directory on your user `PATH`. On Windows, a folder such as `%LOCALAPPDATA%\Programs\clue` works once added through "Edit environment variables for your account." On macOS and Linux, `~/.local/bin` is a common choice; add it to your shell's `PATH` if needed.
4. Open a new terminal and run `clue version`. It should print the version you downloaded, for example `clue 0.7.0`.

The macOS binaries are unsigned and not notarized, so a binary you download through a browser can be blocked by Gatekeeper. First confirm the checksum matches, try `clue version` once, then open **System Settings → Privacy & Security** and click **Open Anyway**. Apple documents this exception in [Open a Mac app from an unknown developer](https://support.apple.com/guide/mac-help/open-a-mac-app-from-an-unknown-developer-mh40616/mac). The install script avoids this: a download made outside the browser carries no quarantine attribute.

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
clue init: 25 created, 0 skipped, 0 linked, 0 index block(s) regenerated
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

If the failure made sense and you want to use Cliewen, start again in a new project or read the greenfield and brownfield guide before initializing an existing repository. Before real work lands, arm the CI wall so a broken thread cannot merge.

## Next

[Choose a greenfield or brownfield adoption path.](./adoption)
