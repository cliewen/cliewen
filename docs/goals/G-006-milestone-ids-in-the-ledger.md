---
id: G-006
type: goal
status: proposed
links: [G-001, ADR-048]
title: Milestone IDs are covered by the corpus-wide identity ledger
---

# G-006 — Milestone IDs are covered by the corpus-wide identity ledger

**Who wants it:** repository contributors and agents allocating a new milestone number (2026-08-10), found while closing [P-013](../plans/P-013-simplification.md) and opening [P-014](../plans/P-014-supersession-edge.md): `clue id next M` returns `M-001`, even though every `M-xxx` up to `M-069` that any plan declares already exists as a milestone row across `docs/plans/*.md` — 68 milestone IDs, not 69 contiguous ones, since M-043 was minted, used, and then withdrawn ([docs/decisions/log.md](../decisions/PDR-039-dependent-changes-carry-authorization.md), 2026-08-01) without a plan ever restoring it. That withdrawn-but-minted number is exactly why a counter, and not a max-scan of what currently exists, is needed. Confirmed empirically: running `clue id next M` returns `M-001` and writes a reserved `M-001` entry into `.clue/id-ledger.yaml` (reverted after the check) — the ledger does not merely fail to seed the `M` counter, it records a false fact on first use, which raises this from a wrong suggestion to a corrupted ledger.

**Why:** [ADR-048](../decisions/ADR-048-corpus-wide-id-ledger.md) states the ledger covers "every native prefix Cliewen mints for itself" and its backfill seeds each prefix's counter "at its current max." Every other prefix examined here (`P`, `CH`, `AN`, `PDR`, `ADR`, `G`) is seeded correctly. `M` is not: its backfill scan almost certainly walks corpus files' own `id:` frontmatter, and a milestone has no file and no frontmatter of its own — it is a row inside a plan's markdown table, identified only by the `ID` column [docs/plans/README.md](../plans/README.md) already defines as a milestone table's declared door. The scan that seeds every other counter has nothing to scan for `M`, so its counter was never created, and `clue id next M` allocates from an implicit zero — silently reissuing a number already used in a plan's own milestone table, which is exactly the collision the ledger exists to prevent for every other prefix.

**Success looks like:**

- `clue id next M` returns a milestone ID higher than every `M-xxx` currently declared in any plan's milestone table, corpus-wide.
- The ledger's backfill (or an equivalent mechanism) seeds the `M` counter from milestone table rows rather than from artifact frontmatter, since a milestone has none.
- `clue validate` reports a milestone ID reused across two plans, or a milestone ID that collides with the ledger's own records, the same way it already reports a native-prefix collision for every other ID kind.
- Nothing about how milestones are declared — a table with `ID` and `Status` headers inside a plan file — changes; a milestone still mints no separate file.
