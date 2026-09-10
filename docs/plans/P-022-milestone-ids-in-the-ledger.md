---
id: P-022
type: plan
status: draft
links: [G-006, G-001, CAP-010, ADR-048]
title: Milestone identities join the corpus-wide ledger
---

# P-022 — Milestone identities join the corpus-wide ledger

Every native prefix the identity ledger covers — `P`, `CH`, `AN`, `PDR`, `ADR`, `G`, and the criterion prefixes — is seeded correctly because each one names an artifact with its own `id:` frontmatter, and the migration backfill (`planLedgerBackfill`, `internal/migrate/migrate.go`) walks `corpus.Scan`'s `ByID` map to find them. A milestone has no file and no frontmatter: it is a row in a plan's milestone table, identified only by the `ID` column the table declares ([docs/plans/README.md](README.md)). The backfill has nothing to scan for `M`, so its counter is never created, and `clue id next M` allocates from an implicit zero — confirmed empirically to return `M-001` and write that false reservation into `.clue/id-ledger.yaml`, even though 68 `M-xxx` identities already exist, non-contiguously, across `docs/plans/*.md` (highest currently declared: `M-090`, from [P-021](P-021-cliewen-states-what-the-product-is-for.md)).

CH-171 (CAP-010) made allocation safe under concurrent access by serializing claims through a Git-remote ledger branch. That work assumed every prefix's starting counter was already trustworthy; it did not touch how any counter is seeded. This plan closes the one prefix where that assumption is false, using the milestone extraction `corpus.PlanMilestones` (`internal/corpus/plans.go`) already provides for `clue next` (CH-172/CH-174) as the read path both the backfill and the new collision check share.

This campaign changes nothing about how a milestone is declared: a table with `ID` and `Status` headers inside a plan file remains the only door, and no milestone gains a file or frontmatter of its own.

## Milestones

| ID | Milestone (exit criterion) | Status | Evidence |
|---|---|---|---|
| M-091 | **The ledger backfill seeds the `M` counter from every milestone table in the corpus, not from frontmatter.** `planLedgerBackfill` marks every `M-xxx` identity `corpus.PlanMilestones` reports across `docs/plans/*.md` live in the same pass it marks artifact IDs live, using the plan's own status (`done`/`dropped` vs. `todo`/`doing`) to decide live vs. retired the same way it already does for criteria. Running the backfill against this repository's own corpus seeds the `M` counter above `M-090`, not at `M-001`. Focused unit evidence covers a corpus with non-contiguous milestone numbering (the withdrawn `M-043`) and a corpus with zero plans. | `todo` | |
| M-092 | **`clue id next M` allocates above every milestone ID declared anywhere in the corpus, safely under CH-171's coordination.** A repository migrated under M-091's backfill allocates its next milestone ID higher than any currently declared, and a repository that already has a ledger (this one) re-backfills or migrates forward without renumbering or double-counting an already-live `M-xxx`. Disposable-Git integration evidence exercises this repository's own migration path end to end. | `todo` | |
| M-093 | **`clue validate` reports a milestone ID collision the same way it already reports a native-prefix collision.** A milestone ID reused across two plans' tables, and a milestone ID that collides with the ledger's own recorded claims, are both reported findings — not a corpus that silently accepts the second table row as if it were a different identity. Focused positive and negative evidence covers both collision shapes and confirms a corpus with no collision stays silent. | `todo` | |

## Mutation rules

Status and evidence fields in the milestone table may mutate in an implementing change's merge digest. Everything else changes only through a declared plan revision backed by a correctly typed decision record. Plan adjustments are decisions.
