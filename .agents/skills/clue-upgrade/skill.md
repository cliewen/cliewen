---
cliewen-skill: true
version: 0.13.0
---

<!-- Generated from Cliewen's canonical skill sources; edit those sources, not this file. -->

# clue-upgrade

## Change scope and tiers

Classify scope before using the Cliewen loop. Three rules decide the tier, by how deeply the change reaches into meaning; take the first rule that matches.

1. **Plain — nothing about meaning changes.** No product behavior, intent, evidence, decision, plan, policy, or methodology changes. Protected product, corpus, test, configuration, build/release, governance/security, agent-rule, skill, and lint surfaces are never plain; neither are commands, contracts, workflows, or normative instructions. Plain work stays outside this skill: branch from `main`, run relevant checks, open a ready PR for human merge, and use no CH identity or Cliewen bookkeeping.
2. **Light — meaning is touched but not changed.** No decision, acceptance-criterion or capability meaning change, semantic plan mutation, or methodology carrier. Typical: protected-surface clarity, dependency bumps, pure refactors, and CI plumbing. Use a Cliewen branch and ready PR whose description names the plan item or plan-less scope, but no transient workspace.
3. **Full — everything else.** Product behavior changes are full even when an existing criterion already states the behavior. Use the whole loop with `/changes/CH-xxx-slug/`.

Two guards hold above the rules. **Uncertainty escalates:** when the tier is unclear, take the higher one. **Discovery escalates immediately:** the moment a decision, an open question, a meaning change, or a methodology-carrier edit appears during work, move to the full loop before continuing.

## Review boundary

Every change branches from the current tip of `main`, never from unaccepted work. Each initiating author takes one Cliewen change to its PR before starting another; independent authors may work in parallel from `main`, and plain changes do not consume this slot. Reviewing or helping update an existing PR does not mint another change or create a global lock. If work must build on an unmerged change, record a blocking open question and stop unless the human explicitly authorizes it. Keep its answered question committed until digest, naming the unmerged base, the authorization, and the meaning the dependent merge would bind; repeat the same disclosure in the ready acceptance brief. That record does not make the base accepted or permit an agent to merge either change. If another change merges before first publication, rebase onto the new `main` tip and repeat verification. After a PR is published, incorporate a newer accepted `main` by merging it into the PR branch with a normal push, never by rewriting hosted history, then repeat verification and review.

For a full Cliewen change, the human accepts the ready pull request with a merge commit. Configure the protected default branch to allow merge commits and disable squash and rebase-and-merge: the merge commit keeps the proposal, implementation, digest, and durable corpus commits reachable from `main`, while the other modes can discard or rewrite that reviewed chain. Rebasing an unpublished local branch before its first publication remains allowed; it is preparation, not the acceptance mode. A forge that cannot enforce the merge-commit boundary is outside the supported full-change adoption path.

Open the PR ready for review only after local verification and the automatic agentic review loop pass on the current commit, never as a draft. The PR is the completed proposal's authorization and protected-integration boundary, not a demand for duplicate human code review: the agent may publish the candidate, but only a human-controlled PR merge accepts it. Unfinished work stays on the branch. An agent never merges its own PR, creates a local merge commit into `main`, or pushes to `main`.

A PR alone displays hosted CI but does not enforce it. Where hosting supports enforcement, the PR triggers CI, branch protection makes its required status check a merge precondition, and the agent cannot silently skip the gate. Never weaken the workflow or required-check policy to make a change pass.

Every review of an existing hosted PR is bound to its observed head SHA. A clean result applies only to that commit; every substantive edit invalidates it. An advisory from a pass with no blocking findings stays in the verification handoff for a later change rather than editing the clean commit: the bounded loop is over, and an edit would create a new candidate that this exact-commit rule requires reviewing. A blocking finding is actionable durable PR state, not private agent memory: where the host supports resolvable review conversations, publish the finding there and leave it unresolved until a hosted commit contains the reviewed repair. Advisories do not become repair-required conversations. If the reviewer cannot publish a resolvable finding, report the PR as not merge-ready and disclose that the host cannot enforce the finding; this fail-safe applies to blocking findings, while advisories stay in the handoff; never claim a chat-only finding has equivalent protection.

Any agent that edits an existing PR becomes the updater for that turn. Before editing, fetch the PR and record its hosted head; before publishing, recheck that head and push only a normal fast-forward update, never force. If the head changed or the push is rejected as non-fast-forward, fetch and reconcile without overwriting remote work, then rerun verification and review on the resulting commit. If the PR merged or closed, stop and report local work as unpublished; never create a follow-up change without explicit human scope.

A review of an existing branch or PR ends in a commit and a push exactly when it changed something. Repairs are committed, verified, reviewed under the loop, and pushed to that PR in the same turn: publication is owed by the repair itself, not deferred to a later report of readiness, because a repair that exists only in a local worktree is private agent memory and no handoff survives it. An agent never elects a local stopping point for itself; only an explicit human request makes committed-but-unpushed work a legitimate stopping point, and the agent then says the repairs are unpublished and the PR is not merge-ready. A review that produced no repair commits and pushes nothing of its own, leaving the reviewed commit exactly as reviewed; incorporating a newer accepted `main` is the separate obligation stated above and is unaffected. This constrains when publication happens, never what is published: the pushed commit is still the one that passed its own local checks and its own review pass.

Ready means the hosted PR contains the exact locally reviewed and verified state. Before reporting any change ready, commit every intended edit, run the applicable local verification and a clean agentic review pass against that commit, require `git status --porcelain` to be empty, push the reviewed commit, and confirm that the ready hosted PR's head branch and SHA equal the current local branch and `HEAD`. Perform the hosted check immediately after opening or updating the PR. Resolve satisfied review conversations only after the hosted head contains their reviewed repair. If either side differs, apply the updater rule above, rerun verification and review on the resulting commit, require a clean worktree, push the reviewed commit, and check again. A local stopping point such as "commit only" exists only because a human asked for one, and is never the agent's own choice; it is preserved work, not a completed or mergeable change, and the agent says that no ready PR exists.

After opening and confirming its initiated PR, an agent stops before initiating another light or full Cliewen change; independent plain changes may still proceed from accepted `main`, and the agent may review or help update an existing PR under the handoff above. Review fixes stay on the same branch and PR and repeat the complete updater handoff before reporting ready again. A follow-up Cliewen change exists only when a human has accepted this one and explicitly scoped the follow-up.

After a human reports the merge, orient before starting anything else: describe the plan's next unfinished step in plain language and ask whether to start it, or say that the plan has nothing left and ask what comes next.


Use when a repository already uses Cliewen and the human wants to find out whether, or bring it up to, a newer release. This is a route into a reviewed repository change, never a background update or authority to merge.

1. Run `clue latest`. It determines whether a newer release is available and, when one is, prints the installation route for the machine it is running on. Do not reproduce an installation command here: one distributed skill cannot know the user's platform. If the release list cannot be reached, explain that Cliewen cannot tell and stop; do not call the repository current.
2. If a newer release is available, read that release's notes, including its `### Migration` section. Identify the coordinated corpus, generated-skill, and CI-caller changes before proposing any repository write.
3. Ask the human whether to upgrade now or later. Do nothing to the repository until they explicitly choose now. A later answer is complete: report the available release and stop without creating a branch, changing a file, or opening a pull request.
4. On an explicit yes, make the repository green and create a branch through its normal review process. Move the binary and repository together: preview `clue migrate`, resolve every finding and notice — including those no command may repair — and apply only the complete, preflighted plan with the required reversal-cost choice. Keep the managed skills, the thin caller, and any repository corpus obligations on the chosen release together.
5. Verify the upgraded repository, commit the complete candidate, run its required checks, and open a ready pull request. Never merge it: the repository's human merge boundary accepts the upgrade.

## Decision records

Route every decision by reversal cost. A cheap-and-local-to-reverse decision is a dated row in `docs/decisions/log.md` (columns `Date | Decision | Why | Change/PR`); otherwise write an ADR for software or corpus architecture, or a PDR for how the project works. A decision adopting a well-established practice cites it by name and records only the local why.

Agent-authored decisions start `status: inferred` and `author: agent`. Merging makes them binding without changing that status. Only explicit human approval promotes a decision to `verified`; record every approver in `accepted-by:`, use the first approval date, and cite the venue. An explicit objection keeps the decision `inferred` and becomes an open question.

`accepted-by:` records only approval given under Cliewen's merge boundary, never acceptance a source record already carried. A record converted from a format with its own acceptance history — names, roles, dates predating the corpus — preserves that history as body prose and keeps `accepted-by: []`, the same empty list any unsigned record carries.

Every decision record is timeless: state what is decided and only the enduring context and rationale needed to understand it. Keep triggering incidents, chronology, conversations, implementation details, and review history in findings, the change workspace, the PR, and Git history.

A decision that changes a methodology contract inventories every live carrier that states the affected contract and updates that complete inventory in the same change. Live carriers include current corpus truth, canonical and generated skills, templates, public or contributor guidance, implementation explanations, CLI text, and distribution metadata. Historical analyses, completed plans, and changelog entries remain pinned history. Add focused guards for stable repaired claims, but do not present those anchors as proof that an arbitrary future carrier inventory is complete; that general obligation remains agent-enforced until a mechanism can derive it.

## Repository-local conventions

For a Cliewen change, apply the repository-local conventions declared in AGENTS.md, including digest requirements such as a user-facing changelog entry. When a release adds or narrows a corpus obligation, preview and apply the supported `clue migrate` migrations before validating the adopted repository; `clue init` remains a non-destructive materializer, not an updater. Plain changes follow only the repository conventions that apply to their changed surface. Local conventions extend the methodology and never override it. If AGENTS.md conflicts with a skill, record the conflict in `open-questions.md` and stop for a human decision; never choose silently.


## Durable work state

An agent's private memory is never where work lives. Anything needed to implement, continue, review, or hand off work belongs in a corpus artifact, the change workspace, or the pull request; private conversation does not survive a change of agent.

A durable record never states a figure a command computes — an artifact count, a coverage percentage, a reported population size. Name the command instead. A number written into prose becomes a hand-maintained obligation that goes stale on the next change and that every later reviewer re-derives, and repairing one writes new prose carrying new numbers, so the finding regenerates instead of converging. Measurements that are the point of a record — an analysis's own results, a milestone's observed evidence — are stated with what produced them and when.
