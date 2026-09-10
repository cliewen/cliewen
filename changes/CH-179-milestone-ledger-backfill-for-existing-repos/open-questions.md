---
id: CH-179-open-questions
type: open-questions
status: open
links: [CH-179]
title: Open questions for CH-179
---

# CH-179 — open questions

None blocking. P-022 already scoped M-092 as seeding an already-existing ledger's `M` counter, distinct from M-091's fresh-ledger path; this change implements exactly that, reusing `corpus.LedgerMilestoneIdentities` and mirroring the criteria backfill's live/retired classification.

This repository's own `.clue/id-ledger.yaml` still carries zero `M-xxx` entries after this change merges — `clue migrate --apply` is not run against it here, since doing so would also fold in an unrelated pending `MIG-003` clue-version bump this change did not scope. Applying the milestone backfill for real is a follow-up the maintainer can run at any time; the mechanism and its disposable-Git integration evidence are what M-092 asks for.
