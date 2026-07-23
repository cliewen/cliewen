---
id: P-005
type: plan
status: active
links: [G-001]
title: Cliewen draws its core — explicit core, one lifecycle, fewer types, memorable tiers
---

# P-005 — Cliewen draws its core

P-004 made the public path earn a newcomer's first try; this campaign simplifies what that newcomer finds. It states Cliewen's core explicitly and protects it behind a red line ([PDR-013](../decisions/PDR-013-explicit-core-red-line.md)), then spends the new criterion — "does the core need it?" — three times: one default status lifecycle instead of fourteen per-type vocabularies, quality scenarios folded into the constraints register, and change-tier routing stated as a few memorable rules instead of dense enumerations. It serves [G-001](../goals/G-001-verifiable-thread.md): a smaller, explicit methodology is easier to verify end to end.

The campaign is sequential. Each milestone is a separate full Cliewen change rooted at accepted `main`, published as a ready pull request, and accepted by human merge before the next begins. Milestone order is a dependency order: the core statement (M-016) supplies the criterion the rest spend, and adopter-type tolerance (M-017) is what makes removing the quality type (M-018) non-breaking for adopters.

## Milestones

| ID | Milestone (exit criterion) | Status | Evidence |
|---|---|---|---|
| M-016 | **The core is stated and protected**: PDR-013 defines the core (verifiable thread, human merge boundary, deterministic judge) and its red line; ARCH-003 `architecture/core.md` carries the durable statement including the periphery list and adopter extension story; C-013 is active with source, enforcement, and promotion trigger; P-005 is active; `clue validate` green | `done` | CH-051: PDR-013, C-013, and ARCH-003 written with red line, periphery list, extension story, and promotion trigger; P-005 active with M-016..M-019; indexes regenerated; `clue validate` OK (81 artifacts) |
| M-017 | **One status lifecycle**: the validator holds one default lifecycle `draft → active → retired` plus only semantically necessary exceptions; unknown artifact types validate against the default lifecycle instead of erroring; the status tables in `docs/README.md` and the scaffold template mirror the code; architecture and analysis records carry the default vocabulary; `go test ./...` and `clue validate` green | `done` | CH-052: ADR-025 (default lifecycle + named exceptions) and ADR-026 (adopter types validate against the default) written; `rules.go` restructured to `defaultLifecycle` + `statusVocabExceptions` with unknown-type fallback; unit tests for adopter types added; 7 analysis + 2 architecture records and `core.md` moved to `active`; both status tables rewritten as default + exceptions; no skill drift; `clue validate` OK (83 artifacts) |
| M-018 | **Quality scenarios are constraints**: `docs/quality/` and its scaffold template are gone; the coverage floor and onboarding bar live as constraints with `source` and `enforcement`; `quality` is absent from the validator's known types; every referencing surface (corpus, guide, skills, PR template, CI) is updated; the retirement is a decision-log row; skills regenerated without drift; strict guide build green | `done` | CH-053: ADR-027 folds the type; QS-001→C-014 (machine) and QS-002→C-015 (human); `docs/quality/` and its scaffold template deleted; tombstone log row added; corpus/guide/skills/PR-template/CI surfaces updated; skills regenerated at v0.6.0 with no drift; `clue validate` OK and strict guide build green. No validator change was needed for the known-types clause: after ADR-026 there is no known-type list, only the default lifecycle plus named exceptions, and `quality` was never an exception |
| M-019 | **Tier routing is memorable**: the canonical change-tier text states plain/light/full as at most three rules plus escalation guards with unchanged semantics; AGENTS.md, the scaffold AGENTS.md, and the guide's change-loop and methodology pages align with it; skills regenerated with a coordinated version bump; strict guide build green | `done` | CH-054: `change-tiers.md.tmpl` rewritten as three named rules (plain, light, full) plus two escalation guards with unchanged tier boundaries; AGENTS.md, the scaffold AGENTS.md, `guide/change-loop.md`, and `guide/methodology.md` restated in the same shape; routing sanity test repointed at the new uncertainty guard; skills regenerated at v0.6.0 with no drift; decision-log row names the template as carrier; `clue validate` OK and strict guide build green |

## Explicitly out of this campaign

Changes to the meaning of the verifiable thread, the merge boundary, or what `clue validate` asserts beyond the listed lifecycle and type-tolerance changes; merging ADR and PDR into one series (rejected in [PDR-006](../decisions/PDR-006-decision-records-are-typed.md)); removing any change tier; new capabilities or acceptance criteria; package-manager distribution channels; migration work in adopter repositories.

## Mutation rules (lintable)

Status and evidence fields in the milestone table may mutate in an implementing change's merge digest. Everything else in this file changes only through a declared plan revision backed by a decision record routed by reversal cost ([C-011](../constraints/C-011-decision-records-typed.md)). Plan adjustments are decisions.
