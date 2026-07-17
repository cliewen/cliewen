---
id: CH-020
type: change
status: open
links: [P-002, M-005, CAP-001, G-001]
title: clue init + quickstart — greenfield onboarding in one command
---

# CH-020 — clue init + quickstart

Serves [P-002](../../docs/plans/P-002-leaves-home.md) milestone M-005: **greenfield in one command**. A new user installs `clue`, runs `clue init`, and gets the whole convention — `/docs` taxonomy, AGENTS.md routing, skills, CI workflow — then reaches a green `clue validate` via a quickstart guide, in under 30 minutes ([CAP-001](../../docs/capabilities/CAP-001-onboarding/README.md)).

## Why now

The manual equivalent of `init` was just performed on model2diagram and works: docs taxonomy since CH-001, vendored skills (`.agents/skills/` canonical + `.claude/skills/` mirror), clue-in-CI on a checksum-verified release binary. That hand-built setup is a fresh, verified fixture — `init` must reproduce it, and the longer we wait the more the fixture drifts from what `init` should generate.

## What

1. **`clue init`** (per the decisions already in [CAP-001 design](../../docs/capabilities/CAP-001-onboarding/design.md)):
   - Emits the `/docs` taxonomy (folder READMEs with index markers, starter artifacts with valid frontmatter), an AGENTS.md routing hub, the five skills (`.agents/skills/` + `.claude/skills/` mirror, stamped with the binary's version), and a CI workflow template that runs `clue validate`.
   - Template source is `/.cliewen/templates` in this repo, embedded into the binary (`go:embed`) so the installed binary is self-contained — no network, no checkout of this repo needed.
   - Idempotent: re-running regenerates index blocks between `clue:index` markers and touches nothing hand-written; existing files (e.g. an adopter's own AGENTS.md) are never overwritten — `init` reports what it skipped.
2. **Quickstart in the repo README**: install → `clue init` → first change loop → green `clue validate`. Prerequisites explicit (git, `clue`, and `gh` for the PR loop as practiced — noted as the convenient path, not a dependency of the method). Skills are linked at the moment the reader's first loop needs them, never explained upfront (guide requirements accumulated in CAP-001 design.md).
3. **Document-type introduction**: the generated `docs/README.md` and folder READMEs explain, in plain language, what each record type is and when a change updates it — ADR (architectural decision), PDR (how the project works), decision-log row (cheap to reverse), capability folder (README / criteria / design per capability), architecture (the whole, the expensive-to-change), constraints (external rules), quality scenarios. The corpus explains itself; the quickstart links into it instead of duplicating it.
4. **CAP-001 goes `active`** with its criteria tested. AC-002 and AC-003 get Go test pairs driving the built binary. AC-001's 30-minute clock spans a human journey (reading, installing) no focused test pair can verify — per the granularity rule it is split: the mechanical path (init → validate green, and its timing envelope) becomes testable criteria; the 30-minute end-to-end promise moves to a quality scenario owned by the quickstart. AC-001 is retired with a tombstone per the AC lifecycle (ADR-007).

## Out of scope

- `clue scaffold` as a standalone index regenerator (M-006) — `init` only regenerates indexes as part of its idempotent run.
- Foreign-soil trials (M-007).
- Any change to `clue validate` rules beyond what the generated corpus needs to pass them unchanged (AC-002).
