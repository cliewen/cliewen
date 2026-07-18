# Contributing to Cliewen

Thank you for helping improve Cliewen. Participation is governed by the [Code of Conduct](CODE_OF_CONDUCT.md).

## Choose the Right Route

- Suspected security vulnerability: follow the private reporting route in [SECURITY.md](SECURITY.md). Never disclose it in a public issue or pull request.
- Private conduct concern: email [flemming.n.larsen+cliewen-conduct@gmail.com](mailto:flemming.n.larsen+cliewen-conduct@gmail.com) with the subject `[Cliewen Conduct] <private report>`. Never open a public conduct issue.
- Reproducible defect: open the structured bug form.
- Desired outcome or unmet need: open the proposed-goal form. A goal issue records demand for consideration; it does not add the goal to Cliewen's accepted plan.
- Small, obvious correction: a pull request may be enough, but it must still use a branch, declare itself plan-less or name the plan item it serves, and follow the light-change rules below.

## Before Starting a Change

Search existing issues, pull requests, and the system-of-record under [`docs/`](docs/README.md). Discuss uncertain scope in an issue before investing in a large implementation.

Every change serves an existing item under [`docs/plans/`](docs/plans/README.md) or explicitly declares itself plan-less. A contributor may have one change in flight at a time, and every branch starts from the current tip of `main`.

Use the next free `CH-xxx` identifier visible in git history and any active `/changes/` workspace, then create a descriptive branch such as `ch-031-short-slug`.

## Choose the Change Tier

A change is light only when it makes no decision, changes no acceptance-criterion or capability meaning, makes no semantic plan mutation, and touches no methodology carrier such as an agent skill, `AGENTS.md` rule, or lint rule. A light change has no `/changes/` workspace; its pull-request description is the proposal and states what, why, and the plan item or plan-less declaration.

Every other change is full. Before implementation, add `/changes/<CH-xxx-slug>/proposal.md`, `tasks.md`, and `open-questions.md`, then commit that proposal by itself. Record unresolved decisions in `open-questions.md` and stop until a human answer can be captured as a typed decision.

## Implement and Digest

Keep the change focused on its proposal and tick each task immediately when it is complete. Update permanent capability, acceptance-criteria, decision, constraint, quality, architecture, and plan artifacts when their meaning changes.

Active acceptance criteria require positive and negative executable evidence. Never weaken a test, lint rule, or quality gate to make a build pass. If a Cliewen-owned skill changes, edit `internal/skills/source/` and run `go generate ./internal/skills`; do not edit `.agents/skills/` or `internal/scaffold/templates/skills/` directly.

Before review, digest a full change into the permanent corpus, update its plan bookkeeping and user-facing [`CHANGELOG.md`](CHANGELOG.md) entry where applicable, and remove its `/changes/` workspace. The final tree proposed for merge must not contain transient change files.

## Verify Locally

Run the repository's mechanical gates:

```text
go build ./...
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
go run ./cmd/clue validate --forbid-changes
git diff --check
```

Total Go statement coverage must remain at least 80%. Also review the meaning the machines cannot judge: links and traceability are truthful, decisions are typed correctly, active constraints and quality scenarios were assessed, and release notes describe user impact.

## Open the Pull Request

Complete the pull-request template and open the pull request as ready for review only after verification passes. Keep review fixes on the same branch and pull request.

The branch and pull request are a proposal; merge is acceptance. A human maintainer merges accepted changes. Agents must never merge their own pull requests, create local merge commits into `main`, or push directly to `main`.

Cliewen does not currently require a Contributor License Agreement or Developer Certificate of Origin sign-off.
