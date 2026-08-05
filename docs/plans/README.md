# Plans

The campaign layer: P-xxx plans with M-xxx milestones. This folder is **flat** — status lives in frontmatter, never in folder names. Plans live on `main`, mutate continuously (bookkeeping in merge digests; semantic changes only via their own change + ADR), and are **frozen, never deleted** when their goal is reached: a completed plan is immutable and this index doubles as the project's achievement overview.

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
- [P-012 — Cliewen closes the brownfield migration gap on re-derived evidence](P-012-migration-gap-closes-on-evidence.md) · `active` — Finishes what P-011's milestone table did not close: a deferral that must be honest rather than merely justified, a report that cannot contradict its own tree, the declined per-criterion artifact decided in the open, and the whole contract tested at a real corpus's size.
- [P-013 — Cliewen is simplified against a stated criterion](P-013-simplification.md) · `draft` — Deferred by four campaigns; holds PDR-013's "does the core need it?" applied to a measured surface, the accumulate-only corpus, and the distributed-work findings P-011 closed without touching.
<!-- clue:index:end -->
