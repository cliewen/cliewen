# Plans

The campaign layer: P-xxx plans with M-xxx milestones. This folder is **flat** — status lives in frontmatter, never in folder names. Plans live on `main`, mutate continuously (bookkeeping in merge digests; semantic changes only via their own change + ADR), and are **frozen, never deleted** when their goal is reached: a completed plan is immutable and this index doubles as the project's achievement overview.

**A plan closes in the digest of the change that completes its last milestone.** Closure is bookkeeping, not a separate decision: the campaign is over the moment its last milestone is evidenced, and a change that leaves the plan `active` publishes a state where the index claims work is in flight that is not. P-007 and P-008 set this; P-005 and P-006 used a separate closing change and it is not the pattern to copy.

**A successor is designated in that same digest when one is decided, and its absence never delays the closure.** P-008 named P-009 as it closed; P-005 and P-006 closed naming none. Whether the next campaign is known is a scheduling question with nothing to do with whether this one is finished, so waiting for it would hold a plan open on an unrelated decision.

The status flip itself is cheap and local to reverse, which is what keeps it bookkeeping. What it is *not* is reversible after the fact: C-008 makes a completed plan immutable, so every milestone's evidence belongs in the table before the digest lands, not after.

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
<!-- clue:index:end -->
