# Contributing to Cliewen

Thank you for helping improve Cliewen. Participation is governed by the [Code of Conduct](CODE_OF_CONDUCT.md).

## Choose the Right Route

- Suspected security vulnerability: follow the private reporting route in [SECURITY.md](SECURITY.md). Never disclose it in a public issue or pull request.
- Private conduct concern: use the [private conduct-reporting address](mailto:flemming&#46;n&#46;larsen&#43;cliewen-conduct&#64;gmail&#46;com) with the subject `[Cliewen Conduct] <private report>`. Never open a public conduct issue.
- Reproducible defect: open the structured bug form.
- Desired outcome or unmet need: open the proposed-goal form. A goal issue records demand for consideration; it does not add the goal to Cliewen's accepted plan.
- Small editorial correction with no effect on behavior, intent, evidence, decisions, plans, policy, or methodology: use the plain-change route below.

## Before Starting a Change

Classify the work before loading the corpus. Only this repository's own `clue` release cut is a dedicated short route: its ordinary PR changes only this repository's release files, has relevant release checks and a human merge, and needs no CH identity, workspace, digest, acceptance brief, or agentic review. This rule does not apply to an adopter's release process. Every other change uses the three rules — plain, then light, then full — and takes the first that matches; light and full are stated under "Choose the Change Tier" below. A change is plain when nothing about meaning changes: it has no effect on product behavior, intent, acceptance evidence, decisions, plans, policy, or methodology. Acceptance evidence includes executable evidence and genuine `Test-type: Human` proof carried by the pull-request acceptance brief. Protected surfaces are never plain: `/docs`, `/changes`, code, tests, configuration, build and release machinery, security and governance policy, `AGENTS.md`, skills, and lint rules. Changes to commands, contracts, user workflow, or normative instructions are not editorial.

Two guards hold above the rules, from this first classification onward. When the tier is unclear, take the higher one; and the moment a decision, an open question, a meaning change, or a methodology-carrier edit appears during work, move to the full loop before continuing.

A plain change starts from the current tip of `main`, uses an ordinary branch, runs checks relevant to the changed surface, and opens a ready pull request. It needs no CH number, plan declaration, proposal, corpus read, Cliewen skill, full verification checklist, plan bookkeeping, or changelog entry.

For every other change, search existing issues, pull requests, and the system-of-record under [`docs/`](docs/README.md). Every Cliewen change serves an existing item under [`docs/plans/`](docs/plans/README.md) or explicitly declares itself plan-less. A contributor may initiate one Cliewen change at a time; plain changes, reviewing, and helping update an existing pull request do not consume another initiated-change slot. Every branch starts from accepted `main` and never from unmerged work unless a human has explicitly authorized a genuine dependency: the committed answered workspace question then names the base, authorization, and meaning the dependent merge would bind, and the ready acceptance brief repeats it.

After classification, load the smallest durable slice that can govern the work: run `clue context <id>` when the request names or resolves to an artifact; otherwise orient at `docs/README.md`, select the closest plan, capability, criterion, or decision, and run `clue context` from there. Follow additional edges only when the task discovers them.

Allocate the next `CH-xxx` identifier with `clue id next CH`, which reserves it in the repository's identity ledger, and mark it live with `clue id live CH-xxx` once the proposal exists. Never derive the number by reading git history: an identifier that a deleted artifact once used would be re-minted, and `clue validate` rejects an artifact missing from the ledger. Then create a descriptive Cliewen branch such as `ch-031-short-slug`.

## Choose the Change Tier

A change is light when meaning is touched but not changed: it makes no decision, changes no acceptance-criterion or capability meaning, makes no semantic plan mutation, and touches no methodology carrier such as an agent skill, `AGENTS.md` rule, or lint rule. A light change has no `/changes/` workspace; its pull-request description is the proposal and states what, why, and the plan item or plan-less declaration.

Every other change is full. Before implementation, add `/changes/<CH-xxx-slug>/proposal.md`, `tasks.md`, and `open-questions.md`, commit that proposal by itself, push it, and open the pull request as a draft; a light change opens its draft at its first commit. Record unresolved decisions in `open-questions.md` and stop until a human answer can be captured as a typed decision.

A product behavior change remains full when an existing criterion already describes the intended behavior. The criterion avoids inventing new acceptance meaning; it does not make the implementation or its executable evidence semantically inert.

## Implement and Digest

Keep the change focused on its proposal, and give a task you mark `[-]` as infeasible its reason on the same line. Update permanent capability, acceptance-criteria, decision, constraint, architecture, and plan artifacts when their meaning changes.

Every working session that changed anything ends by committing and pushing the change branch: a push claims nothing about readiness — the draft pull request simply always shows the work as far as it got, and nothing waits in a local worktree where the next contributor, agent, or the maintainer cannot see it ([PDR-040](docs/decisions/PDR-040-push-is-durability-ready-is-explicit.md)). On a draft, CI validates without the digest gate and requires no acceptance brief; both bind when the pull request is marked ready.

New or revised machine-proven acceptance criteria declare `Test-type: Unit`, `Integration`, `E2E`, or `Performance` and require supported Go, JVM, or Cucumber evidence classified by that type and positive/negative direction, unless the criterion records `(single-direction)`. JVM evidence carries its AC identity, type, and direction on the same Java or Kotlin executable through literal JUnit method tags or the stable named-executable form; class tags and unrelated methods cannot complete the triple. A genuine `Test-type: Human` criterion is named in the acceptance brief as its proof and needs no code test; `@draft` exempts only one genuinely not-yet-proven criterion; an unannotated legacy criterion retains its one-supported-reference contract. `clue validate` validates declarations and references but does not execute tests. Never weaken a test, lint rule, or quality gate to make a build pass. If a Cliewen-owned skill changes, edit `internal/skills/source/` and run `go generate ./internal/skills`; do not edit `.agents/skills/` or `internal/scaffold/templates/skills/` directly.

Before review, digest a full change into the permanent corpus, update its plan bookkeeping, and remove its `/changes/` workspace. Add a [`CHANGELOG.md`](CHANGELOG.md) entry when the change alters what an adopter receives — `clue`'s behaviour, a generated skill's text, or an artifact `clue init` or `clue scaffold` materializes into an adopter repository. A change confined to this repository's own corpus, contributor guidance, and local conventions adds no release note even when it is a full Cliewen change, and neither does a plain editorial change; a change spanning both writes the entry for the adopter-visible part alone. Those categories are a shortcut and the test above is the rule: `.github/pull_request_template.md` is held byte-identical to what `clue init` writes and `.github/workflows/clue-validation.yml` is the reusable workflow an adopter's caller references, so a change to either owes an entry even though `.github/` reads as this repository's own CI ([C-002](docs/constraints/C-002-changelog-per-user-visible-change.md)). The final tree proposed for merge must not contain transient change files.

The generated `.github/workflows/clue.yml` is a thin caller for Cliewen's upstream reusable validation workflow. Keep runner labels, binary source, and writable install-directory choices in that caller; do not copy validation steps or action references into it. A reusable-workflow reference update is the reviewed path for importing upstream scope, warning, acceptance-brief, and digest-gate fixes.

## Verify Locally

For a plain change, run only checks relevant to its changed surface. A guide-Markdown-only edit runs the whitespace check below and `npm run guide:build`.

For a Cliewen change, commit the complete candidate, then run the repository's full mechanical gates against that commit:

```text
go build ./...
go test ./... -coverprofile coverage.out
go tool cover -func coverage.out
go run ./cmd/clue validate --forbid-changes
git diff --check $(git merge-base HEAD origin/main) HEAD
```

Total Go statement coverage must remain at least 80%. `clue-verify` then automatically reviews that same commit before the pull request is marked ready. A coding-agent host with context-isolated delegation starts a fresh read-only reviewer; other hosts disclose an in-context fallback. The loop owns its classification regardless of the reviewer brief: a blocking finding is actionable and enters the hosted repair lifecycle; an advisory is a non-actionable observation for the readiness gate and stays in the verification handoff. Counts and arithmetic disagreements are advisory, while a wrong, missing, or reused identity remains blocking, and the reviewer spends no pass re-deriving figures. A blocking finding returns to the implementing context, is committed, checked against that commit, and reviewed again — scoped to what changed and the carriers it declares — until the current commit receives a pass with no blocking findings. An advisory repair may ride before a pass already required by a blocking repair; an advisory first reported by a pass with no blocking findings stays in the handoff for a later change so the published candidate remains the exact reviewed commit without making the advisory a merge gate. The loop runs up to the maximum number of passes [C-017](docs/constraints/C-017-agentic-review-loop-is-bounded.md) states for this repository, and a further pass runs only after a pass with a blocking finding. Reaching the maximum with blocking findings outstanding stops the loop, reports them to the maintainer, and asks whether to run more; it never permits marking the pull request ready — the pushed branch and its draft simply show where the work stands. The final verification evidence identifies the review mode, reviewed commit, number of passes run, and advisory findings left open.

## Mark the Pull Request Ready

For a plain change, complete only the pull-request summary and relevant verification, then open the pull request after the applicable checks pass. For a Cliewen change, the draft pull request already exists; complete the template's brief, proposal, traceability, and Cliewen checklist, and mark it ready for review only after the applicable checks and automatic agentic review pass on the current head. Keep review fixes on the same branch and pull request; for a Cliewen change, a blocking repair invalidates the earlier clean pass, returns the pull request to draft, and starts a new pass, while an advisory first reported by a clean pass stays in the handoff for a later change instead of editing the reviewed commit.

The branch and pull request are a proposal; a human maintainer performs the human-controlled merge commit that accepts a full Cliewen change. Configure the protected default branch to allow merge commits and disable squash and rebase-and-merge, so the proposal, implementation, digest, and durable corpus history remain reachable from `main`. Agents must never merge their own pull requests, create local merge commits into `main`, or push directly to `main`.

Cliewen does not currently require a Contributor License Agreement or Developer Certificate of Origin sign-off.
