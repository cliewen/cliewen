# Make CI enforce Cliewen

`clue init` gives you a thin GitHub Actions caller for Cliewen's upstream reusable validation workflow. The caller starts unarmed when it uses the default vendored source, so the job warns and skips corpus validation until the pinned Linux release binary and its checksum file are committed under `.github/tools/`.

There are three separate jobs here:

1. Choose the caller inputs that your runner and binary-delivery policy require.
2. Arm the upstream workflow so hosted CI runs the same judge you use locally.
3. Protect `main` so a failed `validate` check actually blocks integration.

CI without branch protection is a dashboard. Branch protection without the validator cannot see a broken Cliewen thread. You need both.

## 1. Choose the caller inputs

Open `.github/workflows/clue.yml`. It should contain one `uses:` reference to `cliewen/cliewen/.github/workflows/clue-validation.yml` at an immutable reference, and four inputs: the generated `clue-version` plus the three local policy choices `runner`, `clue-source`, and `clue-install-directory`. Use the exact generated version; do not substitute `latest`.

The default caller uses `runner: '["ubuntu-latest"]'`, `clue-source: vendored`, and an empty `clue-install-directory`, which stages the verified executable in the runner's temporary directory without requiring root. A repository that needs a self-hosted/no-root runner changes only the caller's runner-label JSON and writable install directory. A repository that downloads the release instead of committing `.github/tools/` changes only `clue-source: release`; the upstream workflow downloads the matching binary and `SHA256SUMS` over HTTPS and verifies the checksum before execution.

Prefer an install directory outside the checkout. An empty value uses the runner's temporary directory, which is cleaned between runs; a path inside the workspace such as `.cliewen/bin` leaves an untracked executable in the working tree, and on a persistent self-hosted runner it survives into the next run. If your policy requires a path inside the workspace, add it to `.gitignore`.

The checksum and the binary come from the same release, so verification catches a truncated or corrupted download rather than establishing independent trust in the release itself. That is the same guarantee the published install scripts give.

The reusable workflow owns its action references and all validation steps. Do not copy its checkout, scope, warning, acceptance-brief, or `clue validate` steps into the caller. Updating the one upstream reference is the reviewed upgrade that imports those fixes while retaining the caller's local choices.

Release binaries emit the reusable workflow's source commit as the reference. A binary built without usable VCS metadata — installed with `go install`, or built from a modified tree — emits the corresponding `vX.Y.Z` release tag instead, because a commit it cannot vouch for would leave your CI unable to resolve the workflow at all. Both forms are immutable; protect release tags from force updates, and replace a tag with the exact source SHA when your hosting policy requires SHA-only references.

## 2. Arm the pinned judge

The default vendored path keeps the release asset and checksum file in the repository. The examples below use `0.7.0`; replace it with the `clue-version` in your caller.

Create the tools directory, then download the Linux amd64 binary and the release checksum file:

::: code-group

```powershell [Windows PowerShell]
New-Item -ItemType Directory -Force .github/tools | Out-Null
gh release download v0.7.0 --repo cliewen/cliewen --pattern 'clue-0.7.0-linux-amd64' --pattern 'SHA256SUMS' --dir .github/tools
```

```sh [macOS and Linux]
mkdir -p .github/tools
gh release download v0.7.0 --repo cliewen/cliewen --pattern 'clue-0.7.0-linux-amd64' --pattern 'SHA256SUMS' --dir .github/tools
```

:::

The runner is Linux amd64 even when you develop on Windows or macOS. Verify the vendored file before committing it:

| System | Check |
|---|---|
| Windows PowerShell | Run `Get-FileHash .github/tools/clue-0.7.0-linux-amd64 -Algorithm SHA256`, then compare it with the matching line in `.github/tools/SHA256SUMS` |
| macOS | Run `shasum -a 256 .github/tools/clue-0.7.0-linux-amd64`, then compare it with the matching line in `.github/tools/SHA256SUMS` |
| Linux | Run `(cd .github/tools && sha256sum -c --ignore-missing SHA256SUMS)` |

Commit both files with the generated caller:

```sh
git add .github/workflows/clue.yml .github/tools/SHA256SUMS .github/tools/clue-0.7.0-linux-amd64
git commit -m "Arm the Cliewen CI wall"
```

Do not edit `clue-version` without replacing both vendored files from the matching release. The upstream workflow verifies the checksum again on every Cliewen run before it stages or executes the binary. The staging directory is caller-selected and defaults to `RUNNER_TEMP`; it is never assumed to be `/usr/local/bin`.

## 3. Know what armed means

| State | What the `validate` job does | Can it protect the corpus? |
|---|---|---|
| Unarmed | Reports a warning and skips `clue validate` | No |
| Armed | Verifies the vendored checksum, stages `clue`, and validates corpus-relevant changes; only full proposal history enables the acceptance-brief gate | Yes, once the check is required |
| Release source | Downloads and verifies the matching binary, then applies the same corpus and full-loop distinction | Yes, once the check is required |
| Simple non-corpus Markdown change | Keeps the stable `validate` job green without running the corpus validator | Yes; the same required check still exists |

The `--forbid-changes` flag is the digest boundary. A pull request with a transient `/changes/CH-xxx-*` workspace is unfinished, even if ordinary validation passes. The hosted check turns red until the change is digested into the permanent corpus and the workspace is removed.

## 4. Require the check on GitHub

Push the caller and let its pull request run once. GitHub needs a recent check before it can offer it as a ruleset requirement. The caller and reusable job are both named `validate`; select the exact check GitHub displays, which may be qualified as `validate / validate`.

Then open **Settings → Rules → Rulesets → New ruleset → New branch ruleset**. GitHub's [ruleset instructions](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-rulesets/creating-rulesets-for-a-repository) and [branch-rule reference](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-rulesets/available-rules-for-rulesets) describe the current controls.

Two separate surfaces control which merge methods are available, and Cliewen needs both pointed the same way. The ruleset's pull-request rule constrains the branches it targets, and that is what enforces the boundary on `main`. **Settings → General → Pull Requests** constrains the whole repository: allow **Create a merge commit**, and disable **Squash and merge** and **Rebase and merge**. A merge commit keeps the exact proposal, implementation, and digest commits reachable from `main`; the other methods can leave the same final tree while losing the reviewed branch chain.

Configure one active ruleset:

| Setting | Value |
|---|---|
| Name | `protect-main` |
| Enforcement status | `Active` |
| Bypass list | Empty |
| Target branches | Include the default branch |
| Restrict deletions | Enabled |
| Require a pull request before merging | Enabled; zero required approvals is enough for Cliewen's human-controlled merge boundary |
| ↳ Allowed merge methods | `Merge` only; clear `Squash` and `Rebase` so the reviewed branch chain survives acceptance |
| Require conversation resolution before merging | Enabled; known agent-review findings remain blocking until their hosted fixes are reviewed |
| Require status checks to pass | Add the exact `validate` check displayed by the caller; require the branch to be up to date before merging |
| Expected source | Select GitHub Actions when GitHub offers a source for the recent `validate` check |
| Block force pushes | Enabled |

An empty bypass list matters. A rule that the normal maintainer or automation can silently bypass is not the merge boundary Cliewen assumes. If another ruleset or an older branch-protection rule also targets `main`, GitHub combines them and applies the most restrictive result.

Rulesets are available for public repositories on GitHub Free, while private repositories need a plan that includes them. If the Rulesets menu is unavailable but classic branch protection is offered, configure the same default-branch requirements there: pull requests only, strict `validate` status check, conversation resolution required, administrators included, and force pushes and deletions disabled. Classic branch protection carries no merge-method control of its own, so the repository-wide **Settings → General → Pull Requests** options are the only lever for the merge shape — enable merge commits and disable squash and rebase-and-merge there. If the hosting plan offers neither enforcement surface, the workflow can report failures and agents can warn about unresolved findings, but neither can block integration.

After saving, inspect the effective default-branch rules:

```sh
gh ruleset check --default --repo OWNER/REPOSITORY
```

You should see the pull-request requirement with merge commits as its only allowed merge method, the conversation-resolution requirement, the exact required `validate` check, deletion restriction, and force-push block.

The ruleset's allowed merge method is what enforces the boundary on the default branch. The repository-wide merge-method settings are a second, broader surface: they decide which buttons GitHub offers anywhere in the repository, including branches no ruleset targets. Align them too, so nobody is offered a button that the default branch will reject:

```sh
gh api repos/OWNER/REPOSITORY --jq '{merge: .allow_merge_commit, squash: .allow_squash_merge, rebase: .allow_rebase_merge}'
```

Expect `merge: true`, `squash: false`, and `rebase: false`. Set them under **Settings → General → Pull Requests**, or with `gh api -X PATCH repos/OWNER/REPOSITORY -F allow_squash_merge=false -F allow_rebase_merge=false`.

Do not remove an existing stronger requirement merely to match this minimum. If the default branch still permits squash or rebase-and-merge, the repository is not ready for a full Cliewen change.

## 5. Prove failure blocks merge

Do this once in a disposable branch. The probe creates a valid change workspace, so normal validation stays green while the merge-time command fails only because the workspace has not been digested. Run it only after the merge-method checks above pass; the final pull request must be accepted with a merge commit.

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

The reusable caller and upstream validation unit are GitHub Actions support. Other forges are outside this workflow contract unless they can provide an equivalent reusable-workflow boundary and the same immutable reference, checksum, stable-check, and protected-merge guarantees.

Forge menus differ, so copy the contract rather than GitHub's labels:

- Run the generated workflow's equivalent on every proposed change to the protected branch.
- Keep one stable check name, `validate`.
- Verify and execute the pinned `clue` release for corpus-relevant changes, while requiring an acceptance brief only when branch history contains a full proposal without a complete simple override.
- Require a merge request or pull request and the successful `validate` check before integration.
- Require resolvable review conversations to be closed before integration; an agent finding remains open until its reviewed fix is hosted.
- Give normal users and automation no bypass for direct pushes, failed checks, force pushes, or branch deletion.
- Run the same undigested-workspace probe and confirm the forge blocks integration.

If your forge or hosting plan cannot enforce those conditions, local Cliewen validation still catches mistakes, but CI is evidence rather than a wall.

## Next

[Choose the smallest Cliewen practice your repository needs.](./adoption)
