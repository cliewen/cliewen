# The change loop

The change loop applies when work belongs in Cliewen. Before loading the corpus, classify the request: if it changes no product behavior, intent, evidence, decision, plan, policy, command, contract, user workflow, or methodology carrier, use an ordinary branch, relevant checks, a ready pull request, and human merge. That plain route has no CH number, proposal metadata, corpus work, Cliewen verification, plan bookkeeping, or changelog entry.

When meaning may change or classification is uncertain, use this loop.

## 1. Branch

Create `ch-xxx-your-slug` from the current tip of `main`. One author takes one Cliewen change to its pull request before starting another, and a change never starts from unaccepted work. Plain changes do not consume the Cliewen slot.

## 2. Propose

A full change commits `/changes/CH-xxx-your-slug/proposal.md` before implementation. The proposal says what will change, why it matters, which plan item it serves or that it is plan-less, and where the decision boundary lies.

`tasks.md` is an ordered checklist with dependencies first. Tick a task the moment it completes. If a blocking decision appears, write it to `open-questions.md` and stop; the answer becomes a typed decision record rather than disappearing into chat.

## 3. Implement

Change the permanent corpus and implementation together. Behavior-changing work names the acceptance criteria it serves. Every active criterion gets focused positive and negative tests in the same change, and every test declares one purpose: an AC ID, unit, sanity, or architecture.

Never weaken a test or lint rule to make the build pass. A failing check is evidence about the change.

## 4. Digest

Once every implementation task is complete or explicitly infeasible, update durable documentation, decisions, indexes, plan bookkeeping, and release notes for shipped behavior or workflow changes. Then delete the `/changes` workspace.

Deletion is the digest: the proposal has been absorbed into the current system truth, and Git retains the delta. `main` never contains `/changes`.

## 5. Verify and review

Commit the complete candidate, run the repository tests and `clue validate --forbid-changes` against that commit, then run `clue-verify` on the same commit. The skill automatically challenges the committed candidate before publication: a host with context-isolated delegation starts a fresh read-only reviewer with the declared intent but without the implementation conversation; other hosts disclose an in-context fallback. The reviewer returns actionable correctness, regression, security, evidence, intent, or unjustified-complexity findings without editing. The implementing context fixes them, commits the repaired candidate, reruns checks against that commit, and starts a new review pass on the same commit; every substantive edit invalidates the earlier clean result. The current commit needs a clean pass before it is locally ready. Fetch the latest `main`; if another change is merged first, rebase and repeat review and verification.

## 6. Open the review gate

The pull request is an authorization and protected-integration gate, not a demand for duplicate human code review. A solo developer may already have accepted the local candidate; the PR still prevents the agent that prepared it from accepting its own work. The agent may publish the branch, but it never merges the pull request or pushes to `main`. The human-controlled merge is the acceptance act.

The PR also gives hosted CI an exact candidate, but a PR alone does not enforce anything. Enforcement requires the CI workflow to run on the PR, its result to be a required status check, and branch protection to block merge until that check passes. Local verification remains fast evidence; protected hosted CI is the safeguard that the agent cannot silently skip. Workflow and protection changes must never weaken the gate merely to make a change pass.

Open a ready pull request only after local review and verification pass. Report the review mode and reviewed commit, then confirm that the hosted head branch and SHA equal the clean, locally reviewed branch and `HEAD` before reporting it ready. A requested local branch or commit stopping point preserves work, but it is incomplete and not mergeable.

Review fixes are committed, locally verified, and agent-reviewed again on the same branch. Once the current commit has a clean pass and the worktree is clean, push it to the existing pull request and repeat the hosted-head check before reporting it ready again. After the human merges, orient on the next unfinished plan milestone and ask before beginning it.

The lifecycle instructions live in [the skills](./skills).
