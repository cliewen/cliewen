---
id: CH-011
type: change
status: open
links: []
title: Convention register — prose rules become constraint artifacts
---

# CH-011 — Convention register via constraints

**Plan-less.** Methodology adjustment from the CH-009 retro: rules that live only in prose (AGENTS.md, READMEs, skills) have no visible backlog and will drift. The register makes every prose-only convention a linted artifact with a visible enforcement class.

## What

The constraints folder finally gets used as designed: `docs/constraints/README.md` has declared `enforcement: machine|agent|human` since the baseline, with only `machine` implemented. Every prose-only convention becomes a constraint artifact `C-xxx-slug.md` carrying `source:` (the doc that states it), `enforcement: agent`, and a **promotion trigger** line — the condition under which it becomes a `machine` rule in clue. The constraints README index is the register table.

### Lint extension (AC-023, new rule in `internal/corpus/rules.go`)

Constraints must carry non-empty `source` and `enforcement`; enforcement vocabulary is `machine|agent|human`. `clue validate`'s OK line reports the count of `enforcement: agent` constraints the same way it reports the born-`inferred` count — the visible backlog. Positive + negative tests mirror the existing rule tests.

### ADR-017

One new ADR: prose conventions are registered as constraint artifacts with enforcement classes, deliberately opening the "`enforcement:` classes beyond `machine`" door that `docs/architecture/architecture.md` lists as out; that line is updated in the same change.

### Seed register (the backfill — this change IS the retroactive cleanup)

- C-001 markdown prose is never hard-wrapped (AGENTS.md rule 5)
- C-002 changelog entry per user-visible change (AGENTS.md rule 6, ADR-012)
- C-003 tasks tick immediately; a `[-]` carries its reason (clue-delta skill)
- C-004 never weaken a test or lint rule to pass (AGENTS.md rule 4)
- C-005 every proposal declares its plan item or plan-less (AGENTS.md rule 2, clue-delta skill)
- C-006 ADRs are timeless prose; method decisions name their carrier (decisions README)
- C-007 diagrams are inline Mermaid (clue-verify checklist)
- C-008 completed plans are immutable (plans README, architecture lifetime classes)
- C-009 type-specific frontmatter fields must be present (clue-verify checklist; the constraint fields themselves are promoted to `machine` by this very change)
- C-010 milestone status values in plan tables follow one vocabulary (currently unlinted anywhere)

Prose statements in AGENTS.md/READMEs stay as the human-readable carrier and gain pointers to their C-xxx entries where that helps; shipped skills stay generic (no repo doc-IDs) so they carry no pointers.

## Files

- `docs/constraints/C-001…C-010` (new) + constraints README rewrite (register intro + index)
- `docs/decisions/ADR-017-*.md` (new, born `inferred`); `docs/architecture/architecture.md` door line updated
- `internal/corpus/rules.go` (checkConstraints), `cmd/clue/main.go` (agent-count on the OK line), tests (AC-023 pair + unit), `docs/capabilities/CAP-002-validate/criteria.md` (AC-023) and design.md
- `AGENTS.md` rules 4/5/6 gain C-xxx pointers
