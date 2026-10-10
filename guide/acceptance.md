---
type: guide
title: Acceptance
---

# Acceptance

Cliewen keeps acceptance with the human. The agent prepares and verifies a tracked change, and you decide whether its outcome and evidence are acceptable. The repository chooses where that decision is recorded.

Local acceptance on `main` is the default, including when `.clue/acceptance.yaml` is absent. `clue init` writes this default explicitly. PR acceptance is an active choice with `clue init --acceptance=pr`. Repository role does not change the policy. Cliewen's own source repository also uses local acceptance on `main`.

| Acceptance | Initial setup | Your approval |
|---|---|---|
| Local — default | `clue init` | Run `clue accept` and confirm in the terminal |
| PR — selected explicitly | `clue init --acceptance=pr` | Accept the ready PR with a merge commit |

## Choose the workflow

`clue init` writes this configuration when no policy exists. The same default also applies without this file:

```yaml
# .clue/acceptance.yaml
mode: local
branch: main
```

The generated `AGENTS.md` tells the agent to follow the repository's policy. Without this file, the policy is local acceptance on `main`. Init preserves an existing explicit choice and refuses conflicting options. Base and candidate must resolve to the same local policy; adding an explicit local/main file leaves the default unchanged. A candidate cannot override an accepted PR policy to accept itself locally. `clue migrate` does not create acceptance policy, and malformed policy blocks migration writes. If your repository requires PRs, select `mode: pr` explicitly and follow its integration rules. Existing user-authored hubs are preserved; update any wording that contradicts the selected policy.

Keep the existing ID ledger and allocation settings. Multiple contributors still reserve IDs through the shared Git allocator. Local acceptance and ID coordination are separate: accepting locally does not permit an allocator to fall back when its remote is unavailable.

## Choose PR acceptance

To select PR acceptance when initializing a repository, run:

```sh
clue init --acceptance=pr
```

This records `mode: pr` and `branch: main` in `.clue/acceptance.yaml`. Use the appropriate integration branch if your repository uses another name. Init preserves existing policy and refuses a conflicting option; change an existing decision under its current acceptance workflow.

The agent publishes the proposal as a draft PR, implements and digests the change, verifies and reviews the exact candidate, and marks the PR ready with its acceptance brief. You judge that brief and accept through a merge commit. Hosted CI and branch protection can enforce admission where configured. See [the PR change loop](./change-loop) for the complete handoff and [the CI wall](./ci-wall) for enforcement setup.

## Local acceptance

Local acceptance relies on the human and agent following their roles. It provides no hosted admission gate, and Git identity and a confirmation prompt do not prove who operated the terminal.

### Prepare and review

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

Cliewen materializes regular files and safe internal Git links from the committed revision. It refuses escaping or missing targets, metadata paths, cycles, excessive expansion and submodules.

### Accept or decline

Read the brief and evidence. If the change is wanted and the evidence supports it, either approve the exact candidate and delegate execution as described below, or run the same command yourself without `--check`:

```sh
clue accept <candidate-sha> --base <base-sha> --brief <file>
```

The command shows the brief and asks you to type `accept CH-xxx` using the change's actual ID. Without a recorded approval, execution requires the interactive confirmation. To decline, cancel or enter another answer; the integration branch stays unchanged. Explain needed changes to the agent before it prepares a new candidate.

After confirmation, the command repeats preflight, locks and compares the integration base, and creates a merge with the base as its first parent and the candidate as its second. The merge tree is exactly the candidate tree. The command runs no Git hooks and does not push. Run your required verification before acceptance and publish the accepted branch through your repository's normal authorized process.

Deleting the change branch afterwards leaves its history reachable. Read the retained brief with `git show --no-patch --format=%B <acceptance-commit>`; the proposal remains in the history of its second parent.

### Handle changed or interrupted work

If the integration branch advances, incorporate it into the change branch, repeat verification and review, and prepare a brief for the new revisions. A dirty checkout, incomplete brief, missing proposal, undigested workspace, or invalid corpus prevents acceptance.

If acceptance is interrupted, inspect `git status` and `git log -1` before retrying. A clean acceptance merge means it completed. If HEAD is still the base but the candidate tree is staged, preserve unrelated edits first. Only when both the index and working tree match the recorded candidate with no intervening edits, restore the base using `git read-tree -m -u <candidate-sha> <base-sha>` and repeat preflight. If the state differs, stop and investigate. Do not use a force reset to hide an uncertain outcome.

### Delegate execution after your approval

You can approve a reviewed change in your conversation with the coding agent and let the agent execute integration. Approve the exact candidate, base and complete acceptance brief; a changed candidate or brief needs a new decision. The agent records your actual approval using the template at `.clue/acceptance/approval.yaml`, including the complete brief's SHA-256 hash, who approved it, the venue, actual statement and recording time.

```sh
clue accept <candidate-sha> --base <base-sha> --brief <file> --approval <record> --check
clue accept <candidate-sha> --base <base-sha> --brief <file> --approval <record>
```

The first command checks without integrating. The second executes the bound decision without a terminal prompt and preserves the approval record with the brief in the merge. The tool verifies matching inputs and provenance fields; it cannot authenticate who authored the statement. An agent must never invent your approval. The command never pushes; explicitly authorize the agent to publish the accepted branch if you want that step delegated too.

Safe internal links are validated from the committed Git revision. Cliewen copies their target content into the private snapshot without following live filesystem links, while preserving the original candidate tree and link modes in acceptance. Links outside the repo, missing targets, cycles, metadata paths, excessive expansion and submodules are refused.

## Next

[Set up the hosted gate for PR acceptance.](./ci-wall)
