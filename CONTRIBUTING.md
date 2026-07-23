# Contributing to Cliewen

Thank you for helping improve Cliewen. Participation is governed by the [Code of Conduct](CODE_OF_CONDUCT.md).

## Choose the Right Route

- Suspected security vulnerability: follow the private reporting route in [SECURITY.md](SECURITY.md). Never disclose it in a public issue or pull request.
- Private conduct concern: use the [private conduct-reporting address](mailto:flemming&#46;n&#46;larsen&#43;cliewen-conduct&#64;gmail&#46;com) with the subject `[Cliewen Conduct] <private report>`. Never open a public conduct issue.
- Reproducible defect: open the structured bug form.
- Desired outcome or unmet need: open the proposed-goal form. A goal issue records demand for consideration; it does not add the goal to Cliewen's accepted plan.
- Small editorial correction with no effect on behavior, intent, evidence, decisions, plans, policy, or methodology: use the plain-change route below.

## Before Starting a Change

Classify the work before loading the corpus; three rules set the tier, and you take the first that matches. A change is plain when nothing about meaning changes: it has no effect on product behavior, intent, executable evidence, decisions, plans, policy, or methodology. Protected surfaces are never plain: `/docs`, `/changes`, code, tests, configuration, build and release machinery, security and governance policy, `AGENTS.md`, skills, and lint rules. Changes to commands, contracts, user workflow, or normative instructions are not editorial.

Two guards hold above the rules, from this first classification onward. When the tier is unclear, take the higher one; and the moment a decision, an open question, a meaning change, or a methodology-carrier edit appears during work, move up a tier before continuing.

A plain change starts from the current tip of `main`, uses an ordinary branch, runs checks relevant to the changed surface, and opens a ready pull request. It needs no CH number, plan declaration, proposal, corpus read, Cliewen skill, full verification checklist, plan bookkeeping, or changelog entry.

For every other change, search existing issues, pull requests, and the system-of-record under [`docs/`](docs/README.md). Every Cliewen change serves an existing item under [`docs/plans/`](docs/plans/README.md) or explicitly declares itself plan-less. A contributor may have one Cliewen change in flight at a time; plain changes do not consume that slot. Every branch starts from accepted `main` and never from unmerged work.

Use the next free `CH-xxx` identifier visible in git history and any active `/changes/` workspace, then create a descriptive Cliewen branch such as `ch-031-short-slug`.

## Choose the Change Tier

A change is light when meaning is touched but not changed: it makes no decision, changes no acceptance-criterion or capability meaning, makes no semantic plan mutation, and touches no methodology carrier such as an agent skill, `AGENTS.md` rule, or lint rule. A light change has no `/changes/` workspace; its pull-request description is the proposal and states what, why, and the plan item or plan-less declaration.

Every other change is full. Before implementation, add `/changes/<CH-xxx-slug>/proposal.md`, `tasks.md`, and `open-questions.md`, then commit that proposal by itself. Record unresolved decisions in `open-questions.md` and stop until a human answer can be captured as a typed decision.

## Implement and Digest

Keep the change focused on its proposal and tick each task immediately when it is complete. Update permanent capability, acceptance-criteria, decision, constraint, architecture, and plan artifacts when their meaning changes.

Active acceptance criteria require positive and negative executable evidence. Never weaken a test, lint rule, or quality gate to make a build pass. If a Cliewen-owned skill changes, edit `internal/skills/source/` and run `go generate ./internal/skills`; do not edit `.agents/skills/` or `internal/scaffold/templates/skills/` directly.

Before review, digest a full change into the permanent corpus, update its plan bookkeeping and release-relevant [`CHANGELOG.md`](CHANGELOG.md) entry where applicable, and remove its `/changes/` workspace. Plain editorial changes add no release note. The final tree proposed for merge must not contain transient change files.

## Verify Locally

For a plain change, run only checks relevant to its changed surface. A guide-Markdown-only edit runs `git diff --check` and `npm run guide:build`.

For a Cliewen change, commit the complete candidate, then run the repository's full mechanical gates against that commit:

```text
go build ./...
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
go run ./cmd/clue validate --forbid-changes
git diff --check
```

Total Go statement coverage must remain at least 80%. `clue-verify` then automatically reviews that same commit before publication. A coding-agent host with context-isolated delegation starts a fresh read-only reviewer; other hosts disclose an in-context fallback. Actionable findings return to the implementing context, and every substantive fix is committed, checked against that commit, and reviewed again until the current commit receives a clean pass. The final verification evidence identifies the review mode and reviewed commit.

## Open the Pull Request

For a plain change, complete only the pull-request summary and relevant verification, then open the pull request after the applicable checks pass. For a Cliewen change, also complete the template's proposal, traceability, and Cliewen checklist, and open the pull request only after the applicable checks and automatic agentic review pass. Keep review fixes on the same branch and pull request; for a Cliewen change, each substantive fix invalidates the earlier clean pass.

The branch and pull request are a proposal; merge is acceptance. A human maintainer merges accepted changes. Agents must never merge their own pull requests, create local merge commits into `main`, or push directly to `main`.

Cliewen does not currently require a Contributor License Agreement or Developer Certificate of Origin sign-off.
