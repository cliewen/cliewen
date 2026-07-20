# Agent routing hub

This repo runs **Cliewen** — its own methodology, dogfooded from commit one. Before loading the corpus, classify the requested work.

A change is **plain** only when it has no effect on product behavior, intent, executable evidence, decisions, plans, policy, or methodology. `/docs`, `/changes`, product code, tests, configuration, build and release machinery, security and governance policy, this file, skills, and lint rules are protected and never plain. Changes to commands, contracts, user workflow, or normative instructions are not editorial. When uncertain, the change is not plain.

A plain change uses an ordinary branch from the current tip of `main`, relevant checks for the changed surface, a ready pull request, and human merge. It has no CH identity, plan declaration, proposal metadata, corpus read, Cliewen skill, `clue validate`, Cliewen verification checklist, plan bookkeeping, or Cliewen-mandated changelog entry. Plain changes do not consume the one-Cliewen-change-in-flight slot and never build on unmerged work.

For every other change, read [`docs/README.md`](docs/README.md) before acting: the `/docs` corpus is the system-of-record and your working memory.

## The rules that bind every change

1. **Everything that mutates `main` goes through branch + PR.** For a Cliewen change, the branch is the proposal; transient files live in `/changes/<CH-xxx-slug>/` on the branch only and are deleted in the digest commit before merge. `main` never contains `/changes/`. A **light** change ([PDR-002](docs/decisions/PDR-002-light-change-tier.md)) — no decision, no AC/capability meaning change, no semantic plan mutation, no methodology carrier touched — skips the workspace: the PR description is the proposal. Every change branches from the current tip of `main`, one Cliewen change is in flight per author, and **agents never merge their own PRs or push to `main`** — the merge is a human act ([C-012](docs/constraints/C-012-agents-never-merge-own-changes.md)).
2. **Ready means the hosted PR contains the verified state.** Before reporting any change ready, commit every intended edit, require a clean worktree, push the branch, and confirm that a ready hosted PR's head branch and SHA equal the locally verified branch and `HEAD` ([C-012](docs/constraints/C-012-agents-never-merge-own-changes.md)). Review fixes repeat that handoff on the existing branch and PR. A human-requested local stopping point is preserved work, but it is incomplete and not mergeable.
3. **Every Cliewen proposal declares which plan item it serves** (see [`docs/plans/`](docs/plans/README.md)) or explicitly declares itself plan-less. The merge digest updates plan bookkeeping in the same commit.
4. **Open questions are artifacts.** When blocked on a decision, write it to `open-questions.md` and stop; human answers become recorded decisions (ADR, PDR, or log row — [C-011](docs/constraints/C-011-decision-records-typed.md)).
5. **Machines enforce form; humans verify meaning.** Never weaken a test or a lint rule to make a build pass — surface the conflict instead ([C-004](docs/constraints/C-004-never-weaken-checks.md)).
6. **Markdown prose is never hard-wrapped.** One line per paragraph and per list item; wrapping is the reader's IDE concern. Line breaks are structural only (headings, lists, tables, code fences) ([C-001](docs/constraints/C-001-no-hard-wrapped-markdown.md)).
7. **Release notes are written for users, in [`CHANGELOG.md`](CHANGELOG.md)** ([C-002](docs/constraints/C-002-changelog-per-user-visible-change.md)). A Cliewen change that affects shipped behavior, a capability, a contract, a command, or a user workflow adds its entry to the `[Unreleased]` section in the digest — what the change means to a user, never a PR title or commit subject. Plain editorial changes add no release note. Cutting a release renames that section to the version; the release workflow publishes the section verbatim as the GitHub release body and fails without it. Auto-generated changelogs, PR lists, and @mentions never appear on a release ([ADR-012](docs/decisions/ADR-012-release-notes-from-changelog.md)).

## Skills

| Skill | When |
|---|---|
| [`clue-analysis`](.agents/skills/clue-analysis/skill.md) | Risks/unknowns first: spikes that end in findings docs |
| [`clue-plan`](.agents/skills/clue-plan/skill.md) | Creating or revising a plan |
| [`clue-delta`](.agents/skills/clue-delta/skill.md) | The change loop: branch → implement → digest → merge |
| [`clue-verify`](.agents/skills/clue-verify/skill.md) | Pre-merge checklist before any Cliewen PR |

The skill files are generated artifacts ([ADR-021](docs/decisions/ADR-021-generated-standalone-skills.md)): to change a skill, edit `internal/skills/source/` and run `go generate ./internal/skills` — never edit `.agents/skills/` or `internal/scaffold/templates/skills/` directly; the repository tests reject hand-edited generated files.
