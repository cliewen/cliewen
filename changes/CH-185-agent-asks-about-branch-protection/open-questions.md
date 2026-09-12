---
id: CH-185-open-questions
type: open-questions
status: open
links: [CH-185]
title: Open questions — the agent asks whether the merge boundary is enforced
---

# CH-185 — Open questions

None blocking.

One observation recorded rather than asked, because it does not block this change. `clue id coordinate` seeds the allocator branch from the local ledger, which knows nothing about identities claimed on unmerged branches. Allocating for this change returned `CH-181`, `PDR-058`, and `AC-190` through `AC-192`, all already held by open PR #207, and each had to be skipped by allocating again. Coordination prevents two contributors colliding from now on; it does not reconcile claims that were already in flight when coordination was enabled. That is a defect in `clue id coordinate` worth its own goal, not a question this change needs answered.
