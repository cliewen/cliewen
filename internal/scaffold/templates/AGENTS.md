# Agent routing hub

This repo runs **Cliewen**. Before loading the corpus, classify the requested work. Three rules set the tier, by how deeply the change reaches into meaning; take the first rule that matches.

1. **Plain — nothing about meaning changes.** No effect on product behavior, intent, executable evidence, decisions, plans, policy, or methodology. A plain change uses an ordinary branch from the current tip of `main`, relevant checks for the changed surface, a ready pull request, and human merge. It has no CH identity, plan declaration, proposal metadata, corpus read, Cliewen skill, `clue validate`, Cliewen verification checklist, plan bookkeeping, or Cliewen-mandated changelog entry. Plain changes do not consume the one-Cliewen-change-in-flight slot and never build on unmerged work. Protected surfaces are never plain: `/docs`, `/changes`, product code, tests, configuration, build and release machinery, security and governance policy, this file, skills, and lint rules. Changes to commands, contracts, user workflow, or normative instructions are not editorial.
2. **Light — meaning is touched but not changed.** No decision, no acceptance-criterion or capability meaning change, no semantic plan mutation, no methodology carrier touched.
3. **Full — everything else.** The whole loop, with a `/changes/CH-xxx-slug/` workspace.

Two guards hold above the rules. **Uncertainty escalates:** when the tier is unclear, take the higher one. **Discovery escalates immediately:** the moment a decision, an open question, a meaning change, or a methodology-carrier edit appears during work, move up a tier before continuing.

For a light or full change, read [`docs/README.md`](docs/README.md) before acting: the `/docs` corpus is the system-of-record and your working memory.

## The rules that bind every change

1. **Everything that mutates `main` goes through branch + PR.** For a Cliewen change, the branch is the proposal; transient files live in `/changes/<CH-xxx-slug>/` on the branch only and are deleted in the digest commit before merge — `main` never contains `/changes/`. A **light** change skips the workspace: the PR description is the proposal. Every change branches from the current tip of `main`, one Cliewen change is in flight per author, and **agents never merge their own PRs or push to `main`** — the merge is a human act.
2. **Ready means the hosted PR contains the reviewed and verified state.** Before publishing a Cliewen change, commit every intended edit and run the applicable local verification against that commit, then automatically run the `clue-verify` agentic review loop on that exact candidate: prefer a context-isolated read-only reviewer where the host supports it, otherwise disclose the in-context fallback; resolve actionable findings, commit and verify the fixes, and repeat until the current commit receives a clean pass. Then require a clean worktree, push the reviewed commit, and confirm that a ready hosted PR's head branch and SHA equal the locally reviewed branch and `HEAD`. Review fixes repeat that local-before-publish handoff on the existing branch and PR. A human-requested local stopping point is preserved work, but it is incomplete and not mergeable.
3. **Every Cliewen proposal declares which plan item it serves** (see [`docs/plans/`](docs/plans/README.md)) or explicitly declares itself plan-less. The merge digest updates plan bookkeeping in the same commit.
4. **Open questions are artifacts.** When blocked on a decision, write it to the change's `open-questions.md` and stop; human answers become recorded decisions (ADR, PDR, or decision-log row — see [`docs/decisions/`](docs/decisions/README.md)).
5. **Machines enforce form; humans verify meaning.** Never weaken a test or a lint rule to make a build pass — surface the conflict instead.
6. **Markdown prose is never hard-wrapped.** One line per paragraph and per list item; wrapping is the reader's IDE concern. Line breaks are structural only (headings, lists, tables, code fences).
7. **The core is behind a red line.** Cliewen's core is the verifiable thread (goal → plan → change → capability → criterion → test), the human merge boundary (agents never merge their own changes), and `clue validate` as deterministic judge. A change that alters what any of these means is never plain and never light: it requires an explicit decision record and human acceptance. Everything else is periphery you may freely extend — including your own artifact types under `docs/` — and periphery never constrains the core.

## Skills

| Skill | When |
|---|---|
| [`clue-analysis`](.agents/skills/clue-analysis/skill.md) | Risks/unknowns first: spikes that end in findings docs |
| [`clue-plan`](.agents/skills/clue-plan/skill.md) | Creating or revising a plan |
| [`clue-delta`](.agents/skills/clue-delta/skill.md) | The change loop: branch → implement → digest → merge |
| [`clue-verify`](.agents/skills/clue-verify/skill.md) | Pre-merge verification and automatic agentic review before any Cliewen PR |
| [`clue-extract`](.agents/skills/clue-extract/skill.md) | Brownfield adoption: transform an existing corpus into `/docs` |

## Repo-local conventions

<!-- Add your project's own layer here: tech stack, build commands, code style, review conventions. A convention that binds every change also registers as a constraint artifact in docs/constraints/ (enforcement: agent until a machine check holds it) — prose here is the readable carrier, the register is the inventory. Repo-local conventions extend the methodology, never override it — when a rule here would contradict a skill, that conflict is an open question for a human, not a silent choice. -->
