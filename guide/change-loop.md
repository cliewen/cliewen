# The change loop

Before editing, the agent recommends **simple** or **full**, says why, and names what discovery would change the recommendation. Simple work leaves the accepted contract unchanged: observational analysis with a named consumer, a defect correction restoring an unchanged criterion, regression evidence, in-contract configuration, refactoring, maintenance, and editorial work. It uses no CH identity, workspace, digest, acceptance brief, or mandatory agentic review and runs checks relevant to its surfaces. Full work changes acceptance-criterion, capability, decision, policy, plan-promise, methodology, or uncovered-behavior meaning; uncertainty makes full the honest recommendation.

Paths and diff size may warn but do not decide meaning. The agent reassesses when semantic scope grows and against the complete diff before integration. If it recommends full and the user explicitly chooses simple, the work proceeds without making the repository untruthful and the final authored commit records `Cliewen-Route: simple`, `Cliewen-Recommendation: full`, and a concise `Cliewen-Override` risk. The user and repository retain integration authority: a route never authorizes an agent push, and an agent pushes directly only with explicit user authorization and repository permission. Release is not a Cliewen route; each adopter defines or omits its own release process.

After classification, start from the smallest durable context that governs the task. `clue context <id>` prints the named artifact and the artifacts it links, out to a stated number of hops; an acceptance-criterion or milestone ID resolves to the artifact that declares it. The slice defaults to one hop and names what the bound held back, so `--depth` widens it on evidence rather than on caution. If the request gives no usable ID, orient at `docs/README.md`, choose the closest artifact, and run the command from there. Shared goals have many reverse dependents, so `context` deliberately follows declared outgoing dependencies only instead of recreating a full-corpus read.

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

A correction that restores behavior already promised by an unchanged criterion may remain simple, with focused regression evidence. If the criterion must change or the implementation introduces behavior it does not cover, the accepted contract changes and the full loop is recommended.

## 1. Branch

For a chosen full loop, create `ch-xxx-your-slug` from accepted `main`. One initiating author takes one full change to its ready pull request before starting another. From the first commit onward, every changed session commits and pushes the branch: a push is durability, not readiness. Simple work and help on another pull request consume no full-change slot. Full work depending on an unmerged change stops at a blocking question until a human authorizes and records that dependency.

## 2. Propose

A full change commits `/changes/CH-xxx-your-slug/proposal.md` before implementation, pushes it, and opens the pull request as a draft. The proposal says what will change, why it matters, which plan item it serves or that it is plan-less, and where the decision boundary lies. The draft is where unfinished work lives and stays visible: it claims nothing and cannot be merged.

When the human wants to review that shape before code exists, they can opt into a spec-first pause. It is not the default loop, and it is not merely a convenience: proposing and implementing are different work, and the boundary between them is where a change can still be split, redirected, or handed to a different agent — which is why the proposal is committed first. At the pause the agent records it in tasks, reports briefly what the proposal says and what implementation involves, and asks whether implementation begins — the proposal is already pushed and readable on the draft pull request. Then it stops until they answer.

`tasks.md` is an ordered checklist with dependencies first, and a task you mark `[-]` as infeasible says why on the same line. If a blocking decision appears, write it to `open-questions.md` and stop; the answer becomes a typed decision record rather than disappearing into chat.

## 3. Implement

Change the permanent corpus and implementation together. Behavior-changing work names the acceptance criteria it serves. Canonical criterion IDs use `<PREFIX>-<digits>[lowercase-suffix]`, including segmented brownfield prefixes such as `SNAP-SQS-001`; carrier aliases are limited to documented underscore tags and hyphen-free Go/JVM named prefixes. A new or revised machine-proven criterion declares `Test-type: Unit`, `Integration`, `E2E`, or `Performance` and gets supported Go, JVM, or Cucumber evidence classified by that type and positive/negative direction, unless it records `(single-direction)`; JVM evidence attaches all three parts to the same supported Java or Kotlin executable. A genuine `Test-type: Human` criterion is named in the acceptance brief as its proof; a not-yet-proven criterion carries `@draft` on its own tag line; an unannotated legacy criterion retains one supported reference. Every test declares one purpose: an AC ID, unit, sanity, or architecture.

Never weaken a test or lint rule to make the build pass. A failing check is evidence about the change.

## 4. Digest

Once every implementation task is complete or explicitly infeasible, update durable documentation, decisions, indexes, plan bookkeeping, and release notes for shipped behavior or workflow changes. Then delete the `/changes` workspace.

Plan bookkeeping includes closing the plan. When a change completes a campaign's last milestone, the same digest sets that plan `completed` — a campaign is over once its last milestone is evidenced, and leaving it open makes the plan index claim work is in flight that is not. A successor plan is designated in that digest when one is decided; not having decided one is no reason to keep the finished plan open.

Deletion is the digest: the proposal has been absorbed into the current system truth, and Git retains the delta. `main` never contains `/changes`.

## 5. Verify and review

Commit the complete candidate, run the repository tests and `clue validate --forbid-changes` against that commit, then run `clue-verify` on the same commit. The skill automatically challenges the full candidate before the pull request is marked ready: a host with context-isolated delegation starts a fresh read-only reviewer with the declared intent but without the implementation conversation; other hosts disclose an in-context fallback. Blocking findings are repaired, rechecked, and reviewed again; advisories stay in the handoff rather than rewriting a clean commit. The loop stops on a pass with no blocking findings or at the repository's maximum for a human decision. Newer accepted `main` is merged into the published branch without rewriting hosted history before checks repeat.

Several agents may collaborate without serializing everybody else. Separate authors keep separate branches from accepted `main`; collaboration is scoped to one pull request. A review of an existing PR names the hosted head it inspected, and actionable findings live as unresolved hosted review conversations where the forge supports them. Any agent asked to fix one becomes the updater for that turn: it fetches and records that head, commits and pushes the repair with the turn that made it — which returns a ready pull request to draft — reviews the repaired commit, confirms the hosted head is the reviewed commit, marks the pull request ready again, and only then resolves the finding. If another updater moved the head, normal Git non-fast-forward protection forces reconciliation and a new verification pass. If findings cannot be published as resolvable conversations, the agent says the PR is not merge-ready and names the enforcement gap; no forge can detect an edit or intention that remained solely in a private worktree.

## 6. Open the review gate

For a full loop, the pull request is an authorization and protected-integration gate, not a demand for duplicate human code review. The agent may publish the branch but does not accept its own full change; the human-controlled merge commit is the acceptance act. Configure that branch for merge commits so proposal, implementation, and digest remain reachable. Simple integration instead follows explicit user authorization and repository policy.

For a full change, the PR starts with an acceptance brief. It asks whether the plan item is still wanted, puts the added or changed criteria and their scenarios in front of the human, and names what merge binds. An authorized dependent change repeats its unmerged base, authorization, and binding meaning there; disclosure does not make the base accepted. The review loop adds an advisory verdict for each changed criterion — whether its referenced tests verify the scenario, something adjacent, or leave it undetermined. That is evidence for human judgment, not a semantic claim by `clue validate`: a green build and a fluent agent do not establish that the outcome is right.

The PR also gives hosted CI an exact candidate, but a PR alone does not enforce anything. Enforcement requires the CI workflow to run on the PR, its result to be a required status check, and branch protection to block merge until that check passes. Local verification remains fast evidence; protected hosted CI is the safeguard that the agent cannot silently skip. The [CI wall guide](./ci-wall) gives the setup and failing-PR probe. Workflow and protection changes must never weaken the gate merely to make a change pass.

Mark the pull request ready for review only after local review and verification pass on the current head. Report the review mode, reviewed commit, number of passes run, and advisory findings left open, then confirm that the hosted head branch and SHA equal the clean, locally reviewed branch and `HEAD` before and immediately after marking it ready. Stopping anywhere else is ordinary rather than an exception: the branch is pushed, the pull request is a draft, and no claim of readiness exists.

Review fixes are committed and pushed with the turn that made them, then locally verified and agent-reviewed again on the same branch; a repair returns a ready pull request to draft until the repaired head has its own clean pass. Once it does and the worktree is clean, repeat the hosted-head check and mark the pull request ready again. After the human merges, orient on the next unfinished plan milestone and ask before beginning it.

## Next

[Learn which Cliewen skill applies to your next change.](./skills)
