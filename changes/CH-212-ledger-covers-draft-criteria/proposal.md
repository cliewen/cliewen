---
id: CH-212
type: change
status: open
links: [CAP-010, ADR-048]
title: The identity ledger covers every declared criterion identity, also those in a criteria file still at draft, so allocation never hands out a number already in use
---

# CH-212 — The identity ledger covers every declared criterion identity, also those in a criteria file still at draft

This change is plan-less: no active plan has an unfinished milestone, and it repairs an allocation defect found by an adopter. It proceeds under draft, inferred VIS-001 without confirming it.

## What and why

The [decision that a persisted ledger replaces scan-and-max allocation](../../docs/decisions/ADR-048-corpus-wide-id-ledger.md) makes `clue id next` increment the ledger's counter and never look at the corpus, so the ledger must already know every identity in use. Criterion identities reach the ledger through the shared reading of criteria declarations, and that reading keeps only criteria in an `active` criteria file plus tombstones. A criteria file still at `status: draft` is left out of the ledger backfill, out of the check in `clue validate` that every declared criterion is in the ledger, and so out of the counters.

Observed in an adopter, tank-royale, on 2026-11-20: its extracted criteria for the battle runner (`BR-001` to `BR-049`), the booter's fallback discovery (`BFD-…`) and the GUI's boot progress (`GBP-…`) sit in criteria files at `draft`, as the adopter's extraction left them. The ledger had no entry and no counter for those prefixes, `clue validate` reported nothing, and `clue id next BR` handed out `BR-001` and `BR-002`, `clue id next BFD` `BFD-001` and `clue id next GBP` `GBP-001`, all already declared, and reserved them on the coordinated allocator. `clue id repair` found nothing to repair. Nothing else stopped the collision; only reading the corpus by hand did.

The change makes the ledger cover every declared criterion identity whatever its file's status:

- **The shared reading** returns criteria in a criteria file at `draft` as identities in use, beside active criteria and tombstones, so the ledger, its backfill and `clue validate` read one population.
- **`clue validate`** reports a declared criterion in a `draft` criteria file that the ledger does not carry, as it already does for an active one; its ledger state may be `reserved` or `live`, since a draft criterion may not be proven yet, but never missing or `retired` unless it is a tombstone.
- **`clue migrate --apply`** gains a backfill, like the one that seeds milestone identities into a ledger made before milestones were covered: it adds each declared criterion identity the ledger does not carry, `live` or `retired` by its declaration, raising each prefix's counter above every declared number, and changes no identity the ledger already carries. A reservation that collided with a declared identity, as in tank-royale, is promoted to `live`, which is what the corpus says it is.

Allocation itself stays an increment of the stored counter; nothing here adds a corpus scan to `clue id next`.

## Challenge the commitment

The assumption most likely to be wrong is that `draft` criteria files are the only way a declared identity escapes the ledger. Others could exist, for example a criteria file whose frontmatter cannot be read or a tag outside the Gherkin block. The cheapest test is the adopter itself: after the backfill, `clue validate` on tank-royale must report no declared criterion missing from the ledger, and `clue id next BR` must return a number above `BR-049`. If either fails, the population is still incomplete and the change returns for a decision.

The credible alternative is a guard in `clue id next` that refuses a number the corpus declares. It catches every escape route, but it is the corpus scan the ledger decision rejected for its cost and for not knowing deleted identities, and it hides the incomplete ledger instead of repairing it.

An implementation could pass every criterion and still fail an adopter whose earlier collisions were already used for new artifacts. The backfill cannot tell a reused number from the declared one; it promotes the reservation to `live`, and `clue validate`'s existing duplicate-identity check reports two artifacts with one identity. That case is left to the human, as today.

## Carriers and proof

New criteria in the team-safe identity allocation capability: AC-252, that an already-existing ledger backfills every declared criterion identity it never covered, also from draft criteria files, and that the next allocation returns a number above every declared one (Unit, positive and negative); and AC-253, that `clue validate` reports a criterion declared in a draft criteria file and missing from the ledger (Unit, positive and negative). Carriers: the shared criteria reading, the migration list and its plan, the validation rule, the capability's criteria and design, the ledger decision's consequences if they say otherwise, and the changelog.

## Human evidence

The maintainer runs the released `clue migrate --apply` in tank-royale and sees `clue validate` report no missing criterion and `clue id next BR` return a number above every declared `BR` criterion.
