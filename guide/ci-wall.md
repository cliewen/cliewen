# Make CI enforce Cliewen

`clue init` gives you a GitHub Actions workflow, but it starts unarmed. The job runs and warns, yet it skips corpus validation until the pinned Linux release binary and its checksum file are committed under `.github/tools/`.

There are two separate jobs here:

1. Arm the workflow so hosted CI runs the same judge you use locally.
2. Protect `main` so a failed `validate` check actually blocks integration.

CI without branch protection is a dashboard. Branch protection without the validator cannot see a broken Cliewen thread. You need both.

## 1. Vendor the pinned judge

Open `.github/workflows/clue.yml` and read `CLUE_VERSION`. Use that exact release; do not substitute `latest`. The examples below use `0.6.0`.

Create the tools directory, then download the Linux amd64 binary and the release checksum file:

::: code-group

```powershell [Windows PowerShell]
New-Item -ItemType Directory -Force .github/tools | Out-Null
gh release download v0.6.0 --repo cliewen/cliewen --pattern 'clue-0.6.0-linux-amd64' --pattern 'SHA256SUMS' --dir .github/tools
```

```sh [macOS and Linux]
mkdir -p .github/tools
gh release download v0.6.0 --repo cliewen/cliewen --pattern 'clue-0.6.0-linux-amd64' --pattern 'SHA256SUMS' --dir .github/tools
```

:::

The runner is Linux amd64 even when you develop on Windows or macOS. Verify the vendored file before committing it:

| System | Check |
|---|---|
| Windows PowerShell | Run `Get-FileHash .github/tools/clue-0.6.0-linux-amd64 -Algorithm SHA256`, then compare it with the matching line in `.github/tools/SHA256SUMS` |
| macOS | Run `shasum -a 256 .github/tools/clue-0.6.0-linux-amd64`, then compare it with the matching line in `.github/tools/SHA256SUMS` |
| Linux | Run `(cd .github/tools && sha256sum -c --ignore-missing SHA256SUMS)` |

Commit both files with the generated workflow:

```sh
git add .github/workflows/clue.yml .github/tools/SHA256SUMS .github/tools/clue-0.6.0-linux-amd64
git commit -m "Arm the Cliewen CI wall"
```

Do not edit `CLUE_VERSION` without replacing both vendored files from the matching release. The workflow verifies the checksum again on every Cliewen run before it installs or executes the binary.

## 2. Know what armed means

| State | What the `validate` job does | Can it protect the corpus? |
|---|---|---|
| Unarmed | Reports a warning and skips `clue validate` | No |
| Armed | Verifies the vendored checksum, installs `clue`, prints its version, and runs `clue validate --forbid-changes` for Cliewen changes | Yes, once the check is required |
| Plain Markdown change | Keeps the stable `validate` job green without running the corpus validator | Yes; the same required check still exists |

The `--forbid-changes` flag is the digest boundary. A pull request with a transient `/changes/CH-xxx-*` workspace is unfinished, even if ordinary validation passes. The hosted check turns red until the change is digested into the permanent corpus and the workspace is removed.

## 3. Require the check on GitHub

Push the armed workflow on a branch and let its pull request run once. GitHub needs a recent `validate` check before it can offer that check as a ruleset requirement.

Then open **Settings → Rules → Rulesets → New ruleset → New branch ruleset**. GitHub's [ruleset instructions](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-rulesets/creating-rulesets-for-a-repository) and [branch-rule reference](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-rulesets/available-rules-for-rulesets) describe the current controls.

Configure one active ruleset:

| Setting | Value |
|---|---|
| Name | `protect-main` |
| Enforcement status | `Active` |
| Bypass list | Empty |
| Target branches | Include the default branch |
| Restrict deletions | Enabled |
| Require a pull request before merging | Enabled; zero required approvals is enough for Cliewen's human-controlled merge boundary |
| Require status checks to pass | Add `validate`; require the branch to be up to date before merging |
| Expected source | Select GitHub Actions when GitHub offers a source for the recent `validate` check |
| Block force pushes | Enabled |

An empty bypass list matters. A rule that the normal maintainer or automation can silently bypass is not the merge boundary Cliewen assumes. If another ruleset or an older branch-protection rule also targets `main`, GitHub combines them and applies the most restrictive result.

Rulesets are available for public repositories on GitHub Free, while private repositories need a plan that includes them. If the Rulesets menu is unavailable but classic branch protection is offered, configure the same default-branch requirements there: pull requests only, strict `validate` status check, administrators included, and force pushes and deletions disabled. If the hosting plan offers neither enforcement surface, the workflow can report failures but cannot block integration.

After saving, inspect the effective default-branch rules:

```sh
gh ruleset check --default --repo OWNER/REPOSITORY
```

You should see the pull-request requirement, required `validate` check, deletion restriction, and force-push block. Do not remove an existing stronger requirement merely to match this minimum.

## 4. Prove failure blocks merge

Do this once in a disposable branch. The probe creates a valid change workspace, so normal validation stays green while the merge-time command fails only because the workspace has not been digested.

```sh
git switch main
git pull --ff-only
git switch -c probe/cliewen-wall
```

Create `changes/CH-999-wall-probe/proposal.md`:

```markdown
---
id: CH-999
type: change
status: open
links: []
title: Verify the protected CI wall
---

# CH-999 — Verify the protected CI wall

This is a disposable, plan-less merge-blocking probe.
```

Check the distinction locally:

```sh
clue validate
clue validate --forbid-changes
```

The first command should pass. The second should fail with `changes: transient workspace present — digest before merge`. Commit the probe and open a pull request:

```sh
git add changes/CH-999-wall-probe/proposal.md
git commit -m "Probe the Cliewen merge wall"
git push -u origin probe/cliewen-wall
gh pr create --title "Probe: Cliewen wall blocks an undigested change" --body "Disposable merge-blocking probe; do not merge."
gh pr checks --watch
```

The `validate` check must fail and GitHub must show the pull request as blocked. If the check is red but the merge button is available, the workflow works but the ruleset does not enforce it yet.

Close the probe without merging it:

```sh
gh pr close --delete-branch
git switch main
git pull --ff-only
git branch -D probe/cliewen-wall
```

## Other forges

Forge menus differ, so copy the contract rather than GitHub's labels:

- Run the generated workflow's equivalent on every proposed change to the protected branch.
- Keep one stable check name, `validate`.
- Verify and execute the pinned `clue` release, then run `clue validate --forbid-changes` for every Cliewen change.
- Require a merge request or pull request and the successful `validate` check before integration.
- Give normal users and automation no bypass for direct pushes, failed checks, force pushes, or branch deletion.
- Run the same undigested-workspace probe and confirm the forge blocks integration.

If your forge or hosting plan cannot enforce those conditions, local Cliewen validation still catches mistakes, but CI is evidence rather than a wall.

## Next

[Choose the smallest Cliewen practice your repository needs.](./adoption)
