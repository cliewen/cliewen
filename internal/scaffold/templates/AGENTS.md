# Agent routing hub

This repo runs **Cliewen**. Before doing anything else, read [`docs/README.md`](docs/README.md): the `/docs` corpus is the system-of-record and your working memory.

## The rules that bind every change

1. **Everything that mutates `main` goes through branch + PR.** The branch is the proposal; transient files live in `/changes/<CH-xxx-slug>/` on the branch only and are deleted in the digest commit before merge — `main` never contains `/changes/`. A **light** change — no decision, no AC/capability meaning change, no semantic plan mutation, no methodology carrier touched — skips the workspace: the PR description is the proposal. Every change branches from the current tip of `main`, one change in flight per author, and **agents never merge their own PRs or push to `main`** — the merge is a human act.
2. **Every proposal declares which plan item it serves** (see [`docs/plans/`](docs/plans/README.md)) or explicitly declares itself plan-less. The merge digest updates plan bookkeeping in the same commit.
3. **Open questions are artifacts.** When blocked on a decision, write it to the change's `open-questions.md` and stop; human answers become recorded decisions (ADR, PDR, or decision-log row — see [`docs/decisions/`](docs/decisions/README.md)).
4. **Machines enforce form; humans verify meaning.** Never weaken a test or a lint rule to make a build pass — surface the conflict instead.
5. **Markdown prose is never hard-wrapped.** One line per paragraph and per list item; wrapping is the reader's IDE concern. Line breaks are structural only (headings, lists, tables, code fences).

## Skills

| Skill | When |
|---|---|
| [`clue-analysis`](.agents/skills/clue-analysis/skill.md) | Risks/unknowns first: spikes that end in findings docs |
| [`clue-plan`](.agents/skills/clue-plan/skill.md) | Creating or revising a plan |
| [`clue-delta`](.agents/skills/clue-delta/skill.md) | The change loop: branch → implement → digest → merge |
| [`clue-verify`](.agents/skills/clue-verify/skill.md) | Pre-merge checklist before any PR |
| [`clue-extract`](.agents/skills/clue-extract/skill.md) | Brownfield adoption: transform an existing corpus into `/docs` |

## Repo-local conventions

<!-- Add your project's own layer here: tech stack, build commands, code style, review conventions. Repo-local conventions extend the methodology, never override it — when a rule here would contradict a skill, that conflict is an open question for a human, not a silent choice. -->
