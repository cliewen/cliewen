# Agent routing hub

This repo runs **Cliewen** — its own methodology, dogfooded from commit one. Before doing anything else, read [`docs/README.md`](docs/README.md): the `/docs` corpus is the system-of-record and your working memory.

## The rules that bind every change

1. **Everything that mutates `main` goes through branch + PR.** The branch is the proposal; transient files live in `/changes/<CH-xxx-slug>/` on the branch only and are deleted in the digest commit before merge. `main` never contains `/changes/`. A **light** change ([PDR-002](docs/decisions/PDR-002-light-change-tier.md)) — no decision, no AC/capability meaning change, no semantic plan mutation, no methodology carrier touched — skips the workspace: the PR description is the proposal. Every change branches from the current tip of `main`, one change in flight per author, and **agents never merge their own PRs or push to `main`** — the merge is a human act ([C-012](docs/constraints/C-012-agents-never-merge-own-changes.md)).
2. **Every proposal declares which plan item it serves** (see [`docs/plans/`](docs/plans/README.md)) or explicitly declares itself plan-less. The merge digest updates plan bookkeeping in the same commit.
3. **Open questions are artifacts.** When blocked on a decision, write it to `open-questions.md` and stop; human answers become recorded decisions (ADR, PDR, or log row — [C-011](docs/constraints/C-011-decision-records-typed.md)).
4. **Machines enforce form; humans verify meaning.** Never weaken a test or a lint rule to make a build pass — surface the conflict instead ([C-004](docs/constraints/C-004-never-weaken-checks.md)).
5. **Markdown prose is never hard-wrapped.** One line per paragraph and per list item; wrapping is the reader's IDE concern. Line breaks are structural only (headings, lists, tables, code fences) ([C-001](docs/constraints/C-001-no-hard-wrapped-markdown.md)).
6. **Release notes are written for users, in [`CHANGELOG.md`](CHANGELOG.md)** ([C-002](docs/constraints/C-002-changelog-per-user-visible-change.md)). A change with user-visible impact adds its entry to the `[Unreleased]` section in the digest — what the change means to a user, never a PR title or commit subject. Cutting a release renames that section to the version; the release workflow publishes the section verbatim as the GitHub release body and fails without it. Auto-generated changelogs, PR lists, and @mentions never appear on a release ([ADR-012](docs/decisions/ADR-012-release-notes-from-changelog.md)).

## Skills

| Skill | When |
|---|---|
| [`clue-analysis`](.agents/skills/clue-analysis/skill.md) | Risks/unknowns first: spikes that end in findings docs |
| [`clue-plan`](.agents/skills/clue-plan/skill.md) | Creating or revising a plan |
| [`clue-delta`](.agents/skills/clue-delta/skill.md) | The change loop: branch → implement → digest → merge |
| [`clue-verify`](.agents/skills/clue-verify/skill.md) | Pre-merge checklist before any PR |
