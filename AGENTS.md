# Agent routing hub

This repo runs **Cliewen** — its own methodology, dogfooded from commit one. Before loading the corpus, classify the requested work. Three rules set the tier, by how deeply the change reaches into meaning; take the first rule that matches.

1. **Plain — nothing about meaning changes.** No effect on product behavior, intent, executable evidence, decisions, plans, policy, or methodology. A plain change uses an ordinary branch from the current tip of `main`, relevant checks for the changed surface, a ready pull request, and human merge. It has no CH identity, plan declaration, proposal metadata, corpus read, Cliewen skill, `clue validate`, Cliewen verification checklist, plan bookkeeping, or Cliewen-mandated changelog entry. Plain changes do not consume the one-Cliewen-change-in-flight slot and never build on unmerged work. Protected surfaces are never plain: `/docs`, `/changes`, product code, tests, configuration, build and release machinery, security and governance policy, this file, skills, and lint rules. Changes to commands, contracts, user workflow, or normative instructions are not editorial.
2. **Light — meaning is touched but not changed** ([PDR-002](docs/decisions/PDR-002-light-change-tier.md)). No decision, no acceptance-criterion or capability meaning change, no semantic plan mutation, no methodology carrier touched.
3. **Full — everything else.** The whole loop, with a `/changes/CH-xxx-slug/` workspace.

Two guards hold above the rules. **Uncertainty escalates:** when the tier is unclear, take the higher one. **Discovery escalates immediately:** the moment a decision, an open question, a meaning change, or a methodology-carrier edit appears during work, move up a tier before continuing.

For a light or full change, read [`docs/README.md`](docs/README.md) before acting: the `/docs` corpus is the system-of-record and your working memory.

## The rules that bind every change

1. **Everything that mutates `main` goes through branch + PR.** For a Cliewen change, the branch is the proposal; transient files live in `/changes/<CH-xxx-slug>/` on the branch only and are deleted in the digest commit before merge. `main` never contains `/changes/`. A **light** change skips the workspace: the PR description is the proposal. Every change branches from the current tip of `main`, one Cliewen change is in flight per author, and **agents never merge their own PRs or push to `main`** — the merge is a human act ([C-012](docs/constraints/C-012-agents-never-merge-own-changes.md)).
2. **Ready means the hosted PR contains the reviewed and verified state.** Before publishing a Cliewen change, commit every intended edit and run the applicable local verification against that commit, then automatically run the `clue-verify` agentic review loop on that exact candidate: prefer a context-isolated read-only reviewer where the host supports it, otherwise disclose the in-context fallback; resolve actionable findings, commit and verify the fixes, and repeat until the current commit receives a clean pass. Then require a clean worktree, push the reviewed commit, and confirm that a ready hosted PR's head branch and SHA equal the locally reviewed branch and `HEAD` ([C-012](docs/constraints/C-012-agents-never-merge-own-changes.md)). Review fixes repeat that local-before-publish handoff on the existing branch and PR. A human-requested local stopping point is preserved work, but it is incomplete and not mergeable.
3. **Every Cliewen proposal declares which plan item it serves** (see [`docs/plans/`](docs/plans/README.md)) or explicitly declares itself plan-less. The merge digest updates plan bookkeeping in the same commit.
4. **Open questions are artifacts.** When blocked on a decision, write it to `open-questions.md` and stop; human answers become recorded decisions (ADR, PDR, or log row — [C-011](docs/constraints/C-011-decision-records-typed.md)).
5. **Machines enforce form; humans verify meaning.** Never weaken a test or a lint rule to make a build pass — surface the conflict instead ([C-004](docs/constraints/C-004-never-weaken-checks.md)).
6. **Markdown prose is never hard-wrapped.** One line per paragraph and per list item; wrapping is the reader's IDE concern. Line breaks are structural only (headings, lists, tables, code fences) ([C-001](docs/constraints/C-001-no-hard-wrapped-markdown.md)).
7. **Release notes are written for users, in [`CHANGELOG.md`](CHANGELOG.md)** ([C-002](docs/constraints/C-002-changelog-per-user-visible-change.md)). A Cliewen change that affects shipped behavior, a capability, a contract, a command, or a user workflow adds its entry to the `[Unreleased]` section in the digest — what the change means to a user, never a PR title or commit subject. Plain editorial changes add no release note. Cutting a release renames that section to the version; the release workflow publishes the section verbatim as the GitHub release body and fails without it. Auto-generated changelogs, PR lists, and @mentions never appear on a release ([ADR-012](docs/decisions/ADR-012-release-notes-from-changelog.md)).
8. **The core is behind a red line.** Cliewen's core is the verifiable thread (goal → plan → change → capability → criterion → test), the human merge boundary, and `clue validate` as deterministic judge ([ARCH-003](docs/architecture/core.md)). A change that alters what any of these means is never plain and never light: it requires an explicit decision record and human acceptance ([C-013](docs/constraints/C-013-core-changes-need-decision.md)). Periphery never constrains the core.

## Skills

| Skill | When |
|---|---|
| [`clue-analysis`](.agents/skills/clue-analysis/skill.md) | Risks/unknowns first: spikes that end in findings docs |
| [`clue-plan`](.agents/skills/clue-plan/skill.md) | Creating or revising a plan |
| [`clue-delta`](.agents/skills/clue-delta/skill.md) | The change loop: branch → implement → digest → merge |
| [`clue-verify`](.agents/skills/clue-verify/skill.md) | Pre-merge verification and automatic agentic review before any Cliewen PR |

The skill files are generated artifacts ([ADR-021](docs/decisions/ADR-021-generated-standalone-skills.md)): to change a skill, edit `internal/skills/source/` and run `go generate ./internal/skills` — never edit `.agents/skills/` or `internal/scaffold/templates/skills/` directly; the repository tests reject hand-edited generated files.
