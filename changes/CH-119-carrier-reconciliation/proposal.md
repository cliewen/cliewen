---
id: CH-119
type: change
status: open
links: [P-011]
title: Carrier reconciliation before migration deletion
---

# CH-119 — Carrier reconciliation before migration deletion

Serves [P-011](../../docs/plans/P-011-truthful-brownfield-migration.md) M-055.

## What

M-052 (identity ledger) and M-053 (migration parity) both close a gap for
criteria and evidence. [PDR-025](../../docs/decisions/PDR-025-brownfield-migration-precedes-simplification.md)
also names a separate finding: the current extraction contract can leave
operational carriers stale — CI workflows, instruction files, freshness
inputs, cross-reference registries, links, and diagram assets that a source
repository's own operations depend on, none of which is a criterion or a
piece of acceptance evidence.

Today these are handled ad hoc in the `clue-extract` mapping (a
hand-maintained-index row, a "preserve every link" paragraph, a diagram
choice rule). M-055 requires a systematic inventory and a deterministic gate:
every operational carrier maps to a retained target or blocks mutation, a
mapped carrier must be current at the pinned source revision, and CI proves
that a stale deleted-path reference, a lost fingerprint, or a missing asset
fails before a migration PR can merge.

## Why

Without this, a migration PR can look green — `clue validate` passes,
`clue parity` is clean — while the target repository silently loses a CI
gate, a freshness check, or a diagram nobody re-linked, because none of
those are acceptance criteria the parity contract already covers. This is
the same red-line concern PDR-025 raises for identity and evidence: a
migration that looks complete while quietly discarding operational meaning.

## Design (see ADR-050 for the full record)

- A **carrier inventory** (YAML), pinned to `source-revision`/`source-location`
  like ADR-048/ADR-049, source-side and human/agent-authored during the
  `clue-extract` rehearsal. Each entry names one carrier (instruction,
  workflow, freshness-input, registry, link, or diagram-asset), its source
  path, and either a mapped target path plus a content fingerprint, or an
  explicit `blocked` marker with a reason.
- A **reconciliation check**, `clue carriers <inventory> [root]`, mirroring
  `clue parity`'s shape: deterministic, derived from the current corpus,
  excluded from the ambient release notice. It fails on a stale
  deleted-path reference (something in the corpus still points at a path
  the inventory's `deleted-paths` names), a lost fingerprint (a mapped
  carrier's current content no longer matches what was pinned), or a
  missing asset (a mapped target path does not exist).
- No source-format parsing in `clue` — inventorying what counts as a
  carrier for a given source format stays a `clue-extract` mapping concern,
  same split as parity's source manifest.

## Scope

- `internal/carriers/` (new package): inventory model, load/validate,
  `Reconcile`.
- `cmd/clue/main.go`: `clue carriers` command.
- `docs/decisions/ADR-050-carrier-reconciliation-manifests.md` (new).
- `docs/capabilities/CAP-003-extract/criteria.md`: AC-115..AC-119.
- `docs/capabilities/CAP-003-extract/design.md`, `README.md`: carrier
  reconciliation section.
- `.claude/skills/clue-extract/mappings/openspec.md` and its canonical
  source under `internal/skills/source/`: systematic carrier inventory
  rows for an OpenSpec source, regenerated via `go generate ./internal/skills`.
- `docs/plans/P-011-truthful-brownfield-migration.md`: M-055 row to `done`.
- `CHANGELOG.md`: `[Unreleased]` entry.
