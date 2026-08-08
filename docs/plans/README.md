# Plans

The campaign layer: P-xxx plans with M-xxx milestones. This folder is **flat** — status lives in frontmatter, never in folder names. Plans live on `main`, mutate continuously (bookkeeping in merge digests; semantic changes only via their own change + ADR), and are **frozen, never deleted** when their goal is reached: a completed plan is immutable and this index doubles as the project's achievement overview.

**A milestone table declares an `ID` column and a `Status` column.** Both headers are what make a table a milestone table: the status cells are read against `todo | doing | done | dropped` ([C-010](../constraints/C-010-milestone-status-vocabulary.md)), and the `M-xxx` identities in the ID column are the declared resolution doors a deferred migration criterion may name ([ADR-053](../decisions/ADR-053-deferred-parity-dispositions-are-accountable.md)). Any other columns are yours; the outer Markdown pipes are optional; a table missing either header is not a milestone table and is not read, so a plan that states its milestones some other way declares no doors. Verbatim examples never declare a door: a fenced block and an HTML block are both skipped when the doors are harvested, so a plan may show a milestone table without minting one.

**Closing is bookkeeping too.** The change that completes the last milestone sets the plan `completed` in its digest, not a separate PR — a campaign is over the moment its last milestone is evidenced, and leaving it `active` makes this index claim work is in flight that is not. A successor is designated in that same digest when one is decided; P-008 named P-009 as it closed, while P-005 and P-006 closed naming none. Because the closed plan is immutable afterwards, every milestone's evidence belongs in the table before the digest lands.

<!-- clue:index:start -->
- [P-002 — Cliewen leaves home](P-002-leaves-home.md) · `completed` (2026-07-18)
- [P-001 — Elaboration baseline](P-001-elaboration-baseline.md) · `completed` (2026-07-13)
- [P-003 — Cliewen goes public](P-003-goes-public.md) · `completed` (2026-07-21)
- [P-004 — Cliewen earns the first try](P-004-first-try.md) · `completed` (2026-07-23)
- [P-005 — Cliewen draws its core](P-005-explicit-core.md) · `completed` (2026-07-24)
- [P-006 — Cliewen digests its first adoption](P-006-first-adoption.md) · `completed` (2026-07-25)
- [P-007 — Cliewen hardens its core](P-007-core-hardening.md) · `completed`
- [P-008 — Cliewen agrees with itself before widening its promises](P-008-self-consistency.md) · `completed`
- [P-009 — Cliewen closes migration blockers and distributed boundaries](P-009-migration-readiness.md) · `completed`
- [P-010 — Cliewen makes staying current something an adopter notices](P-010-adopters-keep-current.md) · `completed` (2026-08-05)
- [P-011 — Cliewen makes brownfield migration truthful](P-011-truthful-brownfield-migration.md) · `completed` (2026-08-06) — Makes source deletion safe: criterion identity, proof parity, in-flight work, and operational carriers all survive an extraction before the source corpus is removed.
- [P-012 — Cliewen closes the brownfield migration gap on re-derived evidence](P-012-migration-gap-closes-on-evidence.md) · `completed` (2026-08-07) — Closed the re-derived migration gaps: an honest evidence base, accountable deferrals, a report bound to its own tree, the per-criterion registry declined in the open, and the ordered migration path proven at assessment scale under a pinned release. Successor: P-013.
- [P-013 — Cliewen is simplified against a stated criterion](P-013-simplification.md) · `active` — Deferred by every campaign that named it, now open with M-062…M-066 and the shipped skills as its primary object: trace every statement in the skills and routing hub to a corpus artifact, trim and reorder them so the surviving rules hold together and bind before they instruct, score and settle the remaining surface and pattern C, determine AN-013's three open findings, and close on re-derived cost evidence.
<!-- clue:index:end -->
