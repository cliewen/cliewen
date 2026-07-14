---
id: CH-010
type: change
status: open
links: []
title: Lightweight change tier and decision log
---

# CH-010 — Lightweight change tier + decision log

**Plan-less.** This change serves no plan item; it is a methodology adjustment confirmed in the CH-009 retro: the full change loop is too heavy for small changes, and ADRs are being written for decisions that are not architectural.

## What

Two mechanisms, one ADR each, plus retroactive cleanup so the existing corpus conforms to the new rules in the same PR.

### Light change tier (ADR-015)

A change qualifies as **light** when ALL hold: no new decision needed (no ADR or log entry), no AC or capability meaning added/changed/retired, no semantic plan mutation, no methodology carrier touched (skills, AGENTS.md rules, lint rules). Examples: typos, doc clarity, dependency bumps, pure refactors, CI plumbing. A light change is branch `ch-xxx-slug` → commits → PR, where the **PR description is the proposal** — no `/changes/` folder. **Escalation rule:** the moment an open question, a decision, or an AC change appears mid-work, create `/changes/CH-xxx-slug/` and continue the full loop. CH numbering stays global.

### Decision log (ADR-016)

ADRs are reserved for decisions that are expensive to reverse (architecture, methodology mechanics, public contracts, AC semantics). Everything else becomes a row in `docs/decisions/log.md` — columns `Date | Decision | Why | Change/PR`. Litmus test: if reversing it later is cheap and local, it's a log row; if it constrains future changes, it's an ADR. The retention convention is amended: retention applies to *decisions*, and demoted entries remain decisions — as log rows; git history is the provenance archive (same philosophy as digests).

## Retroactive cleanup — proposed demotions (confirm in review)

Audited ADR-001…ADR-014 against the litmus test. Proposed for demotion to log rows (file deleted, content becomes a dated row, inbound references repointed):

- **ADR-003 (yaml.v3 library pick)** — cheap to swap; the ADR itself says "revisit only if the dependency ever blocks a release". Referenced by CAP-002 design.md (link + prose) — repointed to the log.
- **ADR-004 (80% coverage number)** — a tunable policy default, not a constraint on future changes. Referenced by QS-001 (link + prose) and cited in ADR-013's prose — repointed to the log.

Staying as ADRs (architectural or methodology mechanics): ADR-001, ADR-002, ADR-005…ADR-014.

## Files

- `docs/decisions/ADR-015-light-change-tier.md`, `ADR-016-decision-log.md` (new, born `inferred`)
- `docs/decisions/log.md` (new) + `log` type in `statusVocab` (`internal/corpus/rules.go`) + `docs/README.md` status table + unit test
- `.agents/skills/clue-delta/skill.md` (light tier, escalation, ADR-vs-log routing in digest)
- `.agents/skills/clue-verify/skill.md` (tier-correctness checklist item)
- `AGENTS.md` rule 1, `docs/decisions/README.md` (ADR vs log split, amended retention wording)
- Demotions: delete ADR-003/ADR-004, add log rows, fix QS-001 / CAP-002 design / ADR-013 / decisions README index
