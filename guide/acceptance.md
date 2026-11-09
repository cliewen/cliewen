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

The agent publishes the proposal as a draft PR, implements and digests the change, verifies and reviews the exact candidate, and marks the PR ready. You decide by merging it with a merge commit. Hosted CI and branch protection can enforce admission where configured. See [the PR change loop](./change-loop) for the complete handoff and [the CI wall](./ci-wall) for enforcement setup.

## Local acceptance

Local acceptance relies on the human and agent following their roles. It provides no hosted admission gate, and Git identity and a confirmation prompt do not prove who operated the terminal.

### Prepare and review

The agent follows the normal tracked workflow: reserve an ID, commit a proposal, implement, digest the change workspace, then verify and review the complete candidate commit. It commits each handoff and publishes only as your repository permits. A change to the reviewed candidate requires another verification and review pass.

Some criteria can only be proved by a person looking at the result. These declare `Test-type: Human`. When the change has one, the agent asks you what you observed before it digests the change, and writes your answer into a `## Human evidence` section of the change's proposal: who observed it, under which conditions, what they saw, the result, and anything that was not observed. The agent writes only what you reported. The proposal stays in Git history after the change workspace is deleted.

On a clean checkout of the configured integration branch at the specified base, the agent runs:

```sh
clue accept <candidate-sha> --base <base-sha> --check
```

The command finds the change from the one proposal that the candidate's history adds after the base, and refuses a candidate that adds none or several. It reads the candidate from Git objects, validates its digested corpus and checks the proposal history. It does not execute tests, exporters, or a review, so their recorded results remain claims for you to assess. Preflight leaves refs, the index, and tracked files unchanged.

Cliewen materializes regular files and safe internal Git links from the committed revision. It refuses escaping or missing targets, metadata paths, cycles, excessive expansion and submodules.

The agent then tells you in a few plain sentences what changed, what was checked, and what was not observed or is still uncertain, and gives you the exact command.

### Accept or decline

If the change is wanted and the evidence supports it, either approve the exact candidate and delegate execution as described below, or run the same command yourself without `--check`:

```sh
clue accept <candidate-sha> --base <base-sha>
```

The command shows the change, its title, the candidate and base, and any human evidence, and asks you to type `accept CH-xxx` using the change's actual ID. Without a recorded approval, execution requires the interactive confirmation. To decline, cancel or enter another answer; the integration branch stays unchanged. Explain needed changes to the agent before it prepares a new candidate.

After confirmation, the command repeats preflight, locks and compares the integration base, and creates a merge with the base as its first parent and the candidate as its second. The merge tree is exactly the candidate tree, and the merge message names the change, its title, the candidate and base, and the proposal's human evidence. The command runs no Git hooks and does not push. Run your required verification before acceptance and publish the accepted branch through your repository's normal authorized process.

Deleting the change branch afterwards leaves its history reachable. Read the acceptance summary with `git show --no-patch --format=%B <acceptance-commit>`; the proposal remains in the history of its second parent.

### Handle changed or interrupted work

If the integration branch advances, incorporate it into the change branch and repeat verification and review. A dirty checkout, missing proposal, undigested workspace, or invalid corpus prevents acceptance.

If acceptance is interrupted, inspect `git status` and `git log -1` before retrying. A clean acceptance merge means it completed. If HEAD is still the base but the candidate tree is staged, preserve unrelated edits first. Only when both the index and working tree match the recorded candidate with no intervening edits, restore the base using `git read-tree -m -u <candidate-sha> <base-sha>` and repeat preflight. If the state differs, stop and investigate. Do not use a force reset to hide an uncertain outcome.

### Delegate execution after your approval

You can approve a reviewed change in your conversation with the coding agent and let the agent execute integration. Approve the exact candidate and base; a changed candidate or base needs a new decision. The agent records your actual approval using the template at `.clue/acceptance/approval.yaml`, with who approved it, the venue, your actual statement and the recording time.

```sh
clue accept <candidate-sha> --base <base-sha> --approval <record> --check
clue accept <candidate-sha> --base <base-sha> --approval <record>
```

The first command checks without integrating. The second executes the bound decision without a terminal prompt and keeps the approval record in the merge message. The tool verifies matching inputs and provenance fields; it cannot authenticate who authored the statement. An agent must never invent your approval. The command never pushes; explicitly authorize the agent to publish the accepted branch if you want that step delegated too.

Safe internal links are validated from the committed Git revision. Cliewen copies their target content into the private snapshot without following live filesystem links, while preserving the original candidate tree and link modes in acceptance. Links outside the repo, missing targets, cycles, metadata paths, excessive expansion and submodules are refused.

## Next

[Set up the hosted gate for PR acceptance.](./ci-wall)
