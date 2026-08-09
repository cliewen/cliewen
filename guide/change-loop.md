# The change loop

The change loop applies when work belongs in Cliewen. Before loading the corpus, classify the request: if nothing about meaning changes — no product behavior, intent, evidence, decision, plan, policy, command, contract, user workflow, or methodology carrier — it is plain, and it uses an ordinary branch, relevant checks, a ready pull request, and human merge. That plain route has no CH number, proposal metadata, corpus work, Cliewen verification, plan bookkeeping, or changelog entry.

When meaning is touched, use this loop. The light tier fits when meaning is only touched and no decision, acceptance meaning, plan semantics, or methodology carrier is affected; everything else is full. When the tier is unclear, take the higher one, and move to the full loop the moment a decision, an open question, a meaning change, or a methodology-carrier edit appears during work.

After classification, start from the smallest durable context that governs the task. `clue context <id>` prints the named artifact and the transitive closure of its outgoing links; an acceptance-criterion or milestone ID resolves to the artifact that declares it. If the request gives no usable ID, orient at `docs/README.md`, choose the closest artifact, and run the command from there. Shared goals have many reverse dependents, so `context` deliberately follows declared outgoing dependencies only instead of recreating a full-corpus read.

## One real change, end to end

[Cliewen pull request #2](https://github.com/cliewen/cliewen/pull/2) made the last edge of the thread machine-checkable. Before that change, an active acceptance criterion could lose its test reference without `clue` noticing; after it, the validator reports the criterion and exits non-zero.

| Stage | What the change carried |
|---|---|
| Need | Make acceptance-criterion-to-test traceability enforceable rather than conventional |
| Proposal | CH-003 declared the intended AC↔test contract and served the baseline plan's traceability milestone |
| Durable capability | [`clue validate`](https://github.com/cliewen/cliewen/tree/main/docs/capabilities/CAP-002-validate) owns criteria such as AC-009, which requires a supported reference for every active legacy machine-proven criterion |
| Classified evidence | Focused tests show both the missing-reference failure and referenced-evidence success path; new or revised criteria also declare proof type and evidence direction |
| Implementation | The validator harvests declared ACs and supported test references, then reports an active AC with no evidence |
| Digest | The transient CH-003 workspace disappeared; the capability, criteria, decisions, implementation, and tests remained |
| Acceptance | The branch became PR #2, CI ran the candidate, and a human merge accepted it into `main` |

That same shape applies to an ordinary product request: state the desired behavior, connect it to a criterion and evidence, implement until the thread and tests are green, digest the temporary proposal, and hand the exact verified commit to the protected pull request.

An existing criterion does not make a behavior change light. The implementation changes executable evidence and may reveal that criterion, test boundary, and product reality disagree, so behavior remains a full reviewed delta. The first adopter-history measurement found real workspace cost—144 transient lines for 76 durable corpus additions across its two full semantic changes—but no behavior-under-existing-criterion example from which to justify removing that boundary. Light routing already left three of five accepted units without a workspace; focused context reduces reading cost without weakening the behavior-to-evidence handoff.

## 1. Branch

Create `ch-xxx-your-slug` from the current tip of `main`. One initiating author takes one initiated Cliewen change to its pull request before starting another, and a change never starts from unaccepted work. Plain changes, reviews, and help updating an existing pull request do not consume another initiated-change slot. Work that genuinely depends on an unmerged change stops at a blocking question until a human authorizes it; the answered workspace record names the base, authorization, and meaning that a dependent merge would bind, and remains committed until digest.

## 2. Propose

A full change commits `/changes/CH-xxx-your-slug/proposal.md` before implementation. The proposal says what will change, why it matters, which plan item it serves or that it is plan-less, and where the decision boundary lies.

When the human wants to review that shape before code exists, they can opt into a spec-first pause. It is not the default loop, and it is not merely a convenience: proposing and implementing are different work, and the boundary between them is where a change can still be split, redirected, or handed to a different agent — which is why the proposal is committed first. At the pause the agent records it in tasks, reports briefly what the proposal says and what implementation involves, and asks two questions: whether implementation begins, and whether the branch is pushed. Pushing is what makes handoff possible and what ends the branch's freedom to be rebased, so it is the human's call rather than a default. Then it stops until they answer.

`tasks.md` is an ordered checklist with dependencies first, and a task you mark `[-]` as infeasible says why on the same line. If a blocking decision appears, write it to `open-questions.md` and stop; the answer becomes a typed decision record rather than disappearing into chat.

## 3. Implement

Change the permanent corpus and implementation together. Behavior-changing work names the acceptance criteria it serves. Canonical criterion IDs use `<PREFIX>-<digits>[lowercase-suffix]`, including segmented brownfield prefixes such as `SNAP-SQS-001`; carrier aliases are limited to documented underscore tags and hyphen-free Go/JVM named prefixes. A new or revised machine-proven criterion declares `Test-type: Unit`, `Integration`, `E2E`, or `Performance` and gets supported Go, JVM, or Cucumber evidence classified by that type and positive/negative direction, unless it records `(single-direction)`; JVM evidence attaches all three parts to the same supported Java or Kotlin executable. A genuine `Test-type: Human` criterion is named in the acceptance brief as its proof; a not-yet-proven criterion carries `@draft` on its own tag line; an unannotated legacy criterion retains one supported reference. Every test declares one purpose: an AC ID, unit, sanity, or architecture.

Never weaken a test or lint rule to make the build pass. A failing check is evidence about the change.

## 4. Digest

Once every implementation task is complete or explicitly infeasible, update durable documentation, decisions, indexes, plan bookkeeping, and release notes for shipped behavior or workflow changes. Then delete the `/changes` workspace.

Plan bookkeeping includes closing the plan. When a change completes a campaign's last milestone, the same digest sets that plan `completed` — a campaign is over once its last milestone is evidenced, and leaving it open makes the plan index claim work is in flight that is not. A successor plan is designated in that digest when one is decided; not having decided one is no reason to keep the finished plan open.

Deletion is the digest: the proposal has been absorbed into the current system truth, and Git retains the delta. `main` never contains `/changes`.

## 5. Verify and review

Commit the complete candidate, run the repository tests and `clue validate --forbid-changes` against that commit, then run `clue-verify` on the same commit. The skill automatically challenges the committed candidate before publication: a host with context-isolated delegation starts a fresh read-only reviewer with the declared intent but without the implementation conversation; other hosts disclose an in-context fallback. The reviewer returns correctness, regression, security, evidence, intent, or unjustified-complexity findings without editing. Every finding identifies the operative requirement or declared intent that is violated and its concrete consequence; authoritative decisions and explicit lifecycle rules govern before alternative readings become findings. Human-controlled merge does not imply duplicate human code review, and a release cut uses its versioned changelog section instead of `[Unreleased]`. The loop owns its classification regardless of the reviewer brief: a blocking finding is actionable and enters the hosted repair lifecycle; an advisory is a non-actionable observation for the publication gate and stays in the verification handoff. Counts and arithmetic disagreements are advisory, while a wrong, missing, or reused identity remains blocking, and the reviewer spends no pass re-deriving figures. The implementing context fixes blocking findings, commits the repaired candidate, reruns checks against that commit, and starts a new review pass on the same commit. An advisory repair may ride before a pass already required by a blocking repair; an advisory first reported by a pass with no blocking findings stays in the handoff for a later change so the published candidate remains the exact reviewed commit without making the advisory a merge gate. The loop runs up to a maximum number of passes — five by default, and a number your repository can set for itself in `AGENTS.md`. It is a maximum and not a quota: one pass that finds nothing blocking ends it, and a further pass runs only after a pass with a blocking finding. When the maximum is reached with blocking findings outstanding, the loop stops, reports them to you, and asks whether to run more, because a change that cannot converge in that many passes has stopped being a review problem. The current commit needs a pass with no blocking findings before it is locally ready, whether or not the maximum was reached. Fetch the latest `main`; if another change merges before the branch is first published, rebase and repeat review and verification. Once a PR exists, merge newer accepted `main` into its branch with a normal push instead of rewriting hosted history, then repeat both checks.

Several agents may collaborate without serializing everybody else. Separate authors keep separate branches from accepted `main`; collaboration is scoped to one pull request. A review of an existing PR names the hosted head it inspected, and actionable findings live as unresolved hosted review conversations where the forge supports them. Any agent asked to fix one becomes the updater for that turn: it fetches and records that head, commits and reviews the repair, pushes without force, confirms the hosted head is the reviewed commit, and only then resolves the finding. If another updater moved the head, normal Git non-fast-forward protection forces reconciliation and a new verification pass. If findings cannot be published as resolvable conversations, the agent says the PR is not merge-ready and names the enforcement gap; no forge can detect an edit or intention that remained solely in a private worktree.

## 6. Open the review gate

The pull request is an authorization and protected-integration gate, not a demand for duplicate human code review. A solo developer may already have accepted the local candidate; the PR still prevents the agent that prepared it from accepting its own work. The agent may publish the branch, but it never merges the pull request or pushes to `main`. For a full change, the human-controlled merge commit is the acceptance act; configure the forge to disable squash and rebase-and-merge so the original proposal, implementation, and digest commits remain reachable from `main`. A local rebase before first publication is allowed, but it is not the acceptance mode.

For a full change, the PR starts with an acceptance brief. It asks whether the plan item is still wanted, puts the added or changed criteria and their scenarios in front of the human, and names what merge binds. An authorized dependent change repeats its unmerged base, authorization, and binding meaning there; disclosure does not make the base accepted. The review loop adds an advisory verdict for each changed criterion — whether its referenced tests verify the scenario, something adjacent, or leave it undetermined. That is evidence for human judgment, not a semantic claim by `clue validate`: a green build and a fluent agent do not establish that the outcome is right.

The PR also gives hosted CI an exact candidate, but a PR alone does not enforce anything. Enforcement requires the CI workflow to run on the PR, its result to be a required status check, and branch protection to block merge until that check passes. Local verification remains fast evidence; protected hosted CI is the safeguard that the agent cannot silently skip. The [CI wall guide](./ci-wall) gives the setup and failing-PR probe. Workflow and protection changes must never weaken the gate merely to make a change pass.

Open a ready pull request only after local review and verification pass. Report the review mode, reviewed commit, number of passes run, and advisory findings left open, then confirm that the hosted head branch and SHA equal the clean, locally reviewed branch and `HEAD` before reporting it ready. A requested local branch or commit stopping point preserves work, but it is incomplete and not mergeable.

Review fixes are committed, locally verified, and agent-reviewed again on the same branch. Once the current commit has a clean pass and the worktree is clean, push it to the existing pull request and repeat the hosted-head check before reporting it ready again. A review of existing hosted work ends in a commit and a push exactly when it changed something: publication is owed by the repair itself, in the same turn, because a fix living only in somebody's worktree is the private state no handoff survives — and a review that found nothing to repair pushes nothing, leaving the reviewed commit exactly as reviewed. Holding a repair back is your call to make, not the agent's; when you ask for it, the agent says the work is unpublished and the pull request is not merge-ready. After the human merges, orient on the next unfinished plan milestone and ask before beginning it.

## Next

[Learn which Cliewen skill applies to your next change.](./skills)
