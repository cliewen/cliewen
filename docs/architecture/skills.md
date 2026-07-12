---
id: ARCH-002
type: architecture
status: verified
links: [ARCH-001]
title: The skills layer — process knowledge as versioned artifacts
---

# The skills layer

## Why skills exist

Skills are the **process-knowledge actor** in the four-actor model ([architecture.md](architecture.md)): they tell the agent what the next right step is, so the CLI only has to judge whether it was done right. Everything that can be convention in a skill is deliberately kept **out of the CLI** — the judge stays boring and finished; the guidance stays editable prose. Skills live in `.agents/skills/` (what agents read), not in `/docs` (what the system is): this file documents *why the set looks like it does*; each skill file is the operational instruction and its own single source of truth for *how*.

## The four skills and how they complement each other

The set is not arbitrary — it is the lifecycle (Foundation §10) cut at its phase boundaries, and each skill hands off to the next:

| Skill | Lifecycle phase | Hands off to |
|---|---|---|
| `clue-analysis` | Risks and unknowns first: spikes ending in findings docs (`/docs/analysis`) | `clue-plan` (findings feed plans) or `clue-delta` (findings feed a change) |
| `clue-plan` | Campaign layer: create or revise a plan with verifiable milestones | `clue-delta` (every plan mutation is itself a change) |
| `clue-delta` | The change loop: branch → propose → implement → digest → merge | `clue-verify` (before every PR) |
| `clue-verify` | The pre-merge checklist — the human-readable twin of `clue validate` | the PR (human verifies meaning; CI verifies form) |

Two invariants tie them together: **every path through the skills ends at the same gate** (verify → PR → merge), and **every skill writes artifacts the next one consumes** — findings feed plans, plans frame deltas, deltas produce the corpus that analysis reads next time. A skill whose output nothing consumes gets removed, the same §2 rule that governs corpus artifacts.

## Rules for future skills

- A new skill must claim a **lifecycle slot or a recurring hand-off** the existing four don't cover, and name the artifacts it reads and writes. "Might be useful" is not a slot.
- `clue-` prefix, namespacing against user skills; independent skills (e.g. `dot-*`) live beside them, never coupled.
- Skills mutate through the change loop like everything else: improving a skill is a change with a branch, a PR, and — when the improvement stems from a decision — an ADR.
