---
id: CH-212-tasks
type: tasks
status: open
links: [CH-212]
title: Cover every declared criterion identity in the identity ledger
---

- [ ] Write AC-252 and AC-253 in the team-safe identity allocation capability.
- [ ] Return criteria from draft criteria files from the shared reading as identities in use (AC-252, AC-253).
- [ ] Report a declared draft criterion missing from the ledger in `clue validate` (AC-253).
- [ ] Add the migration that backfills declared criterion identities into an already-existing ledger and raises the counters (AC-252).
- [ ] Prove both criteria in both directions, and update the capability design, the changelog and the generated evidence.
- [ ] Run the backfill against a copy of tank-royale and check that `clue validate` reports nothing missing and `clue id next BR` returns a number above `BR-049`.
