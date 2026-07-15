---
id: P-002
type: plan
status: active
links: [G-001, G-002]
title: Cliewen leaves home — distribution, greenfield bootstrap, generated indexes
---

# P-002 — Cliewen leaves home

The baseline ([P-001](P-001-elaboration-baseline.md)) proved the methodology inside two repos that share a maintainer and a checkout. This campaign makes Cliewen adoptable without either: a versioned binary anyone can install, a one-command greenfield start, and generated index blocks. Serves [G-001](../goals/G-001-verifiable-thread.md) and [G-002](../goals/G-002-versioned-clue-and-skills.md). Milestone numbering is corpus-global and continues from P-001 (M-001…M-003).

## Milestones

| ID | Milestone (exit criterion) | Status | Evidence |
|---|---|---|---|
| M-004 | **clue ships**: a tagged release produces versioned cross-platform binaries; `clue --version` reports the release; the skills carry the same version stamp (G-002); an adopted repo's CI can install and run `clue validate` (the model2diagram door closes) | `done` | CH-007: `clue version`, per-skill stamps, drift lint (CAP-004, ADR-011), release.yml. Adopter-CI proof: model2diagram CH-003 (its PR #3, merged 2026-07-13) runs `clue validate` in CI against a vendored, checksum-verified release binary. |
| M-005 | **Greenfield in one command**: `clue init` emits the taxonomy, AGENTS.md routing, skills and CI template; CAP-001 (onboarding) goes `active` with its install→first-green-validate-in-30-min criteria tested | `todo` | — |
| M-006 | **Indexes are generated**: `clue scaffold` regenerates README index blocks from folder contents (prose above markers untouched); `checkIndexes` remains the judge | `todo` | — |
| M-007 | **Foreign soil**: the skills are trialed on ≥2 external open-source repos (selected by the human; no shared maintainer, not built for the methodology); each trial produces an `AN-xxx` findings doc; at least one methodology adjustment traces back to trial findings. Trials are findings, not adoptions: no PRs against the foreign repos, no new extraction mappings ([PDR-005](../decisions/PDR-005-foreign-soil-trials.md)) | `todo` | — |

## Explicitly out of this campaign

Multi-agent orchestration; semantic consistency checking; `clue locate`; production feedback loop; non-OpenSpec extraction mappings; implementing PlantUML in model2diagram (that repo's own P-001/M-002); public release / repo visibility decisions (distribution targets the private-repo install story: `go install` and `gh release download`).

## Mutation rules (lintable)

Status fields in the milestone table may mutate in any merge (bookkeeping). Everything else in this file changes only via a change that declares itself a plan revision, backed by a decision record — a PDR for direction and process, an ADR if architectural ([C-011](../constraints/C-011-decision-records-typed.md)). Plan adjustments ARE decisions.
