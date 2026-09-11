---
id: CH-180
type: change
status: open
links: [P-022]
title: clue validate reports a milestone ID collision
---

# CH-180 — clue validate reports a milestone ID collision

M-091 and M-092 seeded the `M` counter and made `clue id next M` allocate
safely, but nothing yet catches a milestone ID that collides despite that
safe allocation: the same `M-xxx` hand-typed into two plans' tables, or a
milestone whose ledger state disagrees with what its own declared status
implies. Every native prefix already gets this treatment from
`checkDuplicateIDs` and `checkLedger` because each one carries a
frontmatter `id:` field; a milestone carries neither a file nor
frontmatter, so it is invisible to both checks (P-022).

This change adds `corpus.checkMilestoneCollisions`, which reuses
`corpus.LedgerMilestoneIdentities` — the same read path M-091's backfill
and M-092's allocation already share — to report both collision shapes,
folding it into `clue validate`'s rule pipeline the same way every other
structural rule is registered.
