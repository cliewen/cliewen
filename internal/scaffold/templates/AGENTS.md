# Agent routing hub

This repo runs **Cliewen**. Before loading the corpus, classify the requested work. Three rules set the tier, by how deeply the change reaches into meaning; take the first rule that matches.

1. **Plain — nothing about meaning changes.** No product behavior, intent, evidence, decision, plan, policy, or methodology changes. Protected surfaces — `/docs`, `/changes`, product code, tests, configuration, build/release, governance/security policy, this file, skills, and lint rules — are never plain; neither are commands, contracts, user workflows, or normative instructions. Use an ordinary branch from `main`, relevant checks for the changed surface, a ready PR, and human merge, with no CH identity or Cliewen bookkeeping. Plain changes do not consume the one-Cliewen-change-in-flight slot and never build on unmerged work.
2. **Light — meaning is touched but not changed.** No decision, acceptance-criterion or capability meaning change, semantic plan mutation, or methodology carrier. Typical: protected-surface clarity, dependency bumps, pure refactors, and CI plumbing. Use a Cliewen branch and ready PR whose description is the proposal, but no `/changes` workspace.
3. **Full — everything else.** Product behavior changes are full even when an existing criterion already states the behavior. Use the whole loop with `/changes/CH-xxx-slug/`.

Two guards hold above the rules. **Uncertainty escalates:** when the tier is unclear, take the higher one. **Discovery escalates immediately:** the moment a decision, an open question, a meaning change, or a methodology-carrier edit appears during work, move to the full loop before continuing.

For a light or full change, read [`docs/README.md`](docs/README.md) only when the request does not name or resolve to an artifact; use it to identify the closest artifact, then run `clue context <id>`. When an identity is already known, run `clue context` directly and read its outgoing-link slice. Read beyond that slice only when the task or a discovered edge requires it. The `/docs` corpus remains the system-of-record and working memory.

## The rules that bind every change

1. **Everything that mutates `main` goes through branch + PR.** For a Cliewen change, the branch is the proposal; transient files live in `/changes/<CH-xxx-slug>/` on the branch only and are deleted in the digest commit before merge — `main` never contains `/changes/`. A **light** change skips the workspace: the PR description is the proposal. Every change branches from the current tip of `main`, one Cliewen change is in flight per initiating author, and **agents never merge their own PRs or push to `main`** — the merge is a human act. Reviewing or helping update an existing PR does not mint another change or create a global lock.
2. **Ready means the hosted PR contains the reviewed and verified state.** Every review of an existing PR names its hosted head; actionable findings become unresolved hosted review conversations where supported, and a clean result applies only to its named commit. Any agent that edits becomes the updater for that turn: fetch and record the hosted head, commit and verify the complete repair, obtain a clean review of that commit, push without force, confirm the PR head equals the reviewed commit, and only then resolve satisfied findings. A changed head or non-fast-forward rejection requires reconciliation and renewed verification; newer accepted `main` is merged into an already-published PR branch without rewriting its history; a merged or closed PR stops with local work reported as unpublished. Before publishing a Cliewen change, automatically run `clue-verify`, require a clean worktree, and complete this exact hosted-head handoff. A human-requested local stopping point is preserved work, but it is incomplete and not mergeable.
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
