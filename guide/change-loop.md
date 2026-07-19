# The change loop

Every mutation follows the same path from proposal to accepted system truth.

## 1. Branch

Create `ch-xxx-your-slug` from the current tip of `main`. One author takes one change to its pull request before starting another, and a change never starts from unaccepted work.

## 2. Propose

A full change commits `/changes/CH-xxx-your-slug/proposal.md` before implementation. The proposal says what will change, why it matters, which plan item it serves or that it is plan-less, and where the decision boundary lies.

`tasks.md` is an ordered checklist with dependencies first. Tick a task the moment it completes. If a blocking decision appears, write it to `open-questions.md` and stop; the answer becomes a typed decision record rather than disappearing into chat.

## 3. Implement

Change the permanent corpus and implementation together. Behavior-changing work names the acceptance criteria it serves. Every active criterion gets focused positive and negative tests in the same change, and every test declares one purpose: an AC ID, unit, sanity, or architecture.

Never weaken a test or lint rule to make the build pass. A failing check is evidence about the change.

## 4. Digest

Once every implementation task is complete or explicitly infeasible, update durable documentation, decisions, indexes, plan bookkeeping, and user-facing release notes. Then delete the `/changes` workspace.

Deletion is the digest: the proposal has been absorbed into the current system truth, and Git retains the delta. `main` never contains `/changes`.

## 5. Verify

Run the repository tests, `clue validate --forbid-changes`, and the human-readable `clue-verify` checklist. Fetch the latest `main`; if another change is merged first, rebase and repeat verification.

## 6. Open the review gate

Open a ready pull request. CI verifies form and the human reviews meaning. Agents never merge their own pull requests or push directly to `main`; merging is the human act that accepts the change.

Review fixes stay on the same branch and pull request. After the human merges, orient on the next unfinished plan milestone and ask before beginning it.

The lifecycle instructions live in [the skills](./skills).

