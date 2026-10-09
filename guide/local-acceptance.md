---
type: guide
title: Local acceptance
---

# Local acceptance

An adopting repository can accept a tracked change without a pull request. The agent prepares a branch and its evidence. You review a brief and run `clue accept` to create a merge commit containing the exact candidate and the acceptance record.

PR acceptance remains the default. Cliewen's own source repository requires PRs. Local acceptance relies on the human and agent following their roles; it provides no hosted admission gate, and neither Git identity nor a confirmation prompt proves who operated the terminal.

## Choose the workflow

Under your repository's existing integration policy, accept this configuration on the integration branch before starting the first locally accepted change:

```yaml
# .clue/acceptance.yaml
mode: local
branch: main
```

Declare the same choice in your repository's `AGENTS.md`. Both the accepted base and candidate must carry the configuration. A candidate cannot enable local acceptance for itself. If your repository requires PRs, follow that policy; this setting does not authorize changing branch protections.

Keep the existing ID ledger and allocation settings. Multiple contributors still reserve IDs through the shared Git allocator. Local acceptance and ID coordination are separate: accepting locally does not permit an allocator to fall back when its remote is unavailable.

## Prepare and review

The agent follows the normal tracked workflow: reserve an ID, commit a proposal, implement, digest the change workspace, then verify and review the complete candidate commit. It commits each handoff and publishes only as your repository permits. A change to the reviewed candidate requires another verification and review pass.

The agent copies the brief template from `.clue/acceptance/brief.md` to the location reported by:

```sh
git rev-parse --git-path clue/acceptance/CH-xxx.md
```

Create the parent directory if needed. This works with linked worktrees as well as ordinary checkouts. The completed brief stays available there until you accept; the acceptance commit will retain its full text. Existing adopters can run `clue init` to materialize the missing template without overwriting their files.

The brief's frontmatter names `type: acceptance-brief`, a title, the change ID, and full lowercase `candidate`, `base`, and `reviewed` commit hashes. `reviewed` must equal `candidate`. Fill all five sections: Intent, Criteria and evidence, Binding decisions, Verification, and Review. Remove placeholders. Name actual commands and observations; an agent must not invent a human observation to complete the form. A structurally valid brief can still describe an outcome you should reject.

On a clean checkout of the configured integration branch at the specified base, the agent runs:

```sh
clue accept <candidate-sha> --base <base-sha> --brief <file> --check
```

The command reads the candidate from Git objects, validates its digested corpus and checks the proposal history. It does not execute tests, exporters, or a review. Their recorded results remain claims for you to assess. Preflight leaves refs, the index, and tracked files unchanged.

The first version requires regular tracked files. It refuses tracked symlinks and submodules rather than claiming to have validated content it cannot safely materialize in isolation.

## Accept or decline

Read the brief and evidence. If the change is wanted and the evidence supports it, run the same command yourself without `--check`:

```sh
clue accept <candidate-sha> --base <base-sha> --brief <file>
```

The command shows the brief and asks you to type `accept CH-xxx` using the change's actual ID. There is no unattended confirmation flag. To decline, cancel or enter another answer; the integration branch stays unchanged. Explain needed changes to the agent before it prepares a new candidate.

After confirmation, the command repeats preflight, locks and compares the integration base, and creates a merge with the base as its first parent and the candidate as its second. The merge tree is exactly the candidate tree. The command runs no Git hooks and does not push. Run your required verification before acceptance and publish the accepted branch through your repository's normal authorized process.

Deleting the change branch afterwards leaves its history reachable. Read the retained brief with `git show --no-patch --format=%B <acceptance-commit>`; the proposal remains in the history of its second parent.

## Handle changed or interrupted work

If the integration branch advances, incorporate it into the change branch, repeat verification and review, and prepare a brief for the new revisions. A dirty checkout, incomplete brief, missing proposal, undigested workspace, or invalid corpus prevents acceptance.

If acceptance is interrupted, inspect `git status` and `git log -1` before retrying. A clean acceptance merge means it completed. If HEAD is still the base but the candidate tree is staged, preserve unrelated edits first. Only when both the index and working tree match the recorded candidate with no intervening edits, restore the base using `git read-tree -m -u <candidate-sha> <base-sha>` and repeat preflight. If the state differs, stop and investigate. Do not use a force reset to hide an uncertain outcome.

## Next

[Set up the hosted gate for PR acceptance.](./ci-wall)
