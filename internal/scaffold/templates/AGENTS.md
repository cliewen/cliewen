# Agent routing hub

This repo runs **Cliewen**.

**Before your first tool call — whatever the request is, a question, a review, or a change — run `clue latest --quiet`.**

That is the whole instruction; the rest is why. It prints one line if a newer `clue` release exists and nothing at all when you are current, offline, or unable to tell, so it costs a line only when there is something to say; the answer is cached for a day, so asking again is free. If it says you are behind, route to [`clue-upgrade`](.agents/skills/clue-upgrade/skill.md) — the human decides whether to upgrade now or later, and nothing upgrades without that answer. If `clue` reports `latest` as an unknown command, that *is* the answer: your binary predates the check, so it is behind. The ordinary workflow commands — `clue init`, `scaffold`, `context`, `migrate`, `id`, `refs`, and `report` — also volunteer the same one-line notice on their own, so a session that runs any of them learns this without being asked; the check above is what covers a session that runs none. This is the only reason to run it unprompted; it reaches the network, so it never belongs in a validation verdict or a required check.

Before loading the corpus, classify the requested work by its effect on meaning.

1. **Plain — nothing about meaning changes.** No product behavior, intent, evidence, decision, plan, policy, or methodology changes. Protected surfaces — `/docs`, `/changes`, product code, tests, configuration, build/release, governance/security policy, this file, skills, and lint rules — are never plain; neither are commands, contracts, user workflows, or normative instructions. Use an ordinary branch from `main`, relevant checks for the changed surface, a ready PR, and human merge, with no CH identity or Cliewen bookkeeping. Plain changes do not consume the one-Cliewen-change-in-flight slot and never build on unmerged work.
2. **Light — meaning is touched but not changed.** No decision, acceptance-criterion or capability meaning change, semantic plan mutation, or methodology carrier. Use [`clue-delta`](.agents/skills/clue-delta/skill.md)'s change loop without a workspace.
3. **Full — everything else.** Use [`clue-delta`](.agents/skills/clue-delta/skill.md)'s full change loop with `/changes/CH-xxx-slug/`.

When the tier is unclear, take the higher one. When a decision, open question, meaning change, or methodology-carrier edit appears during work, move to the full loop before continuing.

For a light or full change, read [`docs/README.md`](docs/README.md) only when the request does not name or resolve to an artifact; use it to identify the closest artifact, then run `clue context <id>`. When an identity is already known, run `clue context` directly and read its outgoing-link slice. Read beyond that slice only when the task or a discovered edge requires it. The `/docs` corpus remains the system-of-record and working memory.

## Repository conventions

**Markdown prose is never hard-wrapped.** One line per paragraph and per list item; wrapping is the reader's IDE concern. Line breaks are structural only (headings, lists, tables, code fences).

**The core is behind a red line.** Cliewen's core is the verifiable thread (goal → plan → change → capability → criterion → acceptance evidence, including classified executable references and genuine Human proof in the acceptance brief), the human merge boundary, and `clue validate` as deterministic judge. A change that alters what any of these means is never plain and never light: it requires an explicit decision record and human acceptance. Everything else is periphery you may freely extend — including your own artifact types under `docs/` — and periphery never constrains the core.

## Skills

| Skill | When |
|---|---|
| [`clue-analysis`](.agents/skills/clue-analysis/skill.md) | Risks/unknowns first: spikes that end in findings docs |
| [`clue-plan`](.agents/skills/clue-plan/skill.md) | Creating or revising a plan |
| [`clue-upgrade`](.agents/skills/clue-upgrade/skill.md) | Checking for and, with human approval, carrying out a coordinated repository upgrade |
| [`clue-delta`](.agents/skills/clue-delta/skill.md) | The change loop: branch → implement → digest → merge |
| [`clue-verify`](.agents/skills/clue-verify/skill.md) | Pre-merge verification and automatic agentic review before any Cliewen PR |
| [`clue-extract`](.agents/skills/clue-extract/skill.md) | Brownfield adoption: transform an existing corpus into `/docs` |

## Repo-local conventions

<!-- Add your project's own layer here: tech stack, build commands, code style, review conventions. A convention that binds every change also registers as a constraint artifact in docs/constraints/ (enforcement: agent until a machine holds it, partial once one holds part of it, human when none ever can) — prose here is the readable carrier, the register is the inventory. Repo-local conventions extend the methodology, never override it — when a rule here would contradict a skill, that conflict is an open question for a human, not a silent choice. -->
