---
id: ARCH-002
type: architecture
status: active
links: [ARCH-001, ADR-021, ADR-043, PDR-042]
title: The skills layer — process knowledge as versioned artifacts
---

# The skills layer

## Why skills exist

Skills are the **process-knowledge actor** in the four-actor model ([architecture.md](architecture.md)): they tell the agent what the next right step is, so the CLI only has to judge whether it was done right. Everything that can be convention in a skill is deliberately kept **out of the CLI** — the judge stays boring and finished; the guidance stays editable prose. Skills live in `.agents/skills/` (what agents read), not in `/docs` (what the system is): this file documents *why the set looks like it does*. Each generated skill directory is the complete operational instruction; its short `skill.md` entry point routes the agent to required local references only when their condition is reached ([ADR-059](../decisions/ADR-059-progressive-standalone-skill-directories.md)).

## The six skills and how they complement each other

The set is not arbitrary — it is the lifecycle (Foundation §10) cut at its phase boundaries, and each skill hands off to the next:

| Skill | Lifecycle phase | Hands off to |
|---|---|---|
| `clue-analysis` | Risks and unknowns first: spikes ending in findings docs (`/docs/analysis`) | `clue-plan` (findings feed plans) or `clue-delta` (findings feed a change) |
| `clue-plan` | Campaign layer: create or revise a plan with verifiable milestones | `clue-delta` (every plan mutation is itself a change) |
| `clue-upgrade` | Existing-adopter release check and human-authorized coordinated upgrade | the reviewed repository change the affirmative choice becomes — simple work by default, and `clue-delta` only for a decision of the repository's own that the release forces ([PDR-043](../decisions/PDR-043-upgrade-routes-as-simple-work.md)) |
| `clue-delta` | The change loop: branch → propose → implement → digest → merge | `clue-verify` (before every Cliewen PR) |
| `clue-verify` | Pre-merge verification followed by automatic adversarial agent review | the locally verified and reviewed candidate, then the PR (human controls merge; CI verifies form) |
| `clue-extract` | Brownfield adoption: one-time transform of an existing corpus into `/docs`, everything born `inferred` and non-decisions classified by reversal cost (ADR-008, ADR-035) | `clue-delta` (the extraction runs as the adopted repo's first change loop) |

Two invariants tie them together: **every full-loop path through the skills ends at the same gate** (implement → commit → verify → context-isolated review where supported → PR → merge), and **every skill produces a hand-off the next one consumes** — findings feed plans or changes, review findings return to implementation, and deltas produce the corpus that analysis reads next time. Simple work is recommended before this layer and invokes no lifecycle skill ([PDR-042](../decisions/PDR-042-routing-recommends-contract-aware-effort.md)); an observational analysis may still invoke `clue-analysis` for the investigation without turning its resulting document into a full delta. A skill whose output nothing consumes gets removed.

## Shared rules and standalone outputs

The six skill directories repeat the cross-cutting instructions each one can need because every directory must remain independently installable. They are not independently authored: `internal/skills/source/skills/` holds the workflow templates and `internal/skills/source/shared/` holds shared rules such as decision routing, simple/full change routing, repository-local conventions, and the review boundary ([ADR-021](../decisions/ADR-021-generated-standalone-skills.md)). The generator splits those rendered sections into local references and writes the condition for reading each one into the directory's `skill.md` router.

`go generate ./internal/skills` composes complete skill directories into `.agents/skills/` and the embedded `clue init` template tree. Repository tests compare every entry point and reference in both trees with the canonical rendering and fail on changed, missing, or unexpected generator-owned files. Contributors edit the source tree and regenerate; generated outputs are never edited directly.

## Rules for future skills

- A new skill must claim a **lifecycle slot or a recurring hand-off** the existing set does not cover, and name the artifacts or hand-off it reads and produces. "Might be useful" is not a slot.
- `clue-` prefix, namespacing against user skills; independent skills (e.g. `dot-*`) live beside them, never coupled.
- A change to methodology meaning is recommended for the full loop, edits canonical source fragments or templates, runs the generator, and records a decision when the improvement stems from one.
