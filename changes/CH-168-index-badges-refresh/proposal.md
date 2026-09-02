---
id: CH-168
type: change
status: open
links: [G-010, ADR-041, ADR-046, CAP-005, CAP-002]
title: A regenerated index row's badge follows the artifact it names
---

# CH-168 — A regenerated index row's badge follows the artifact it names

## What

`clue scaffold` refreshes an existing index row's status badge from the artifact it links, while leaving the row's title and description exactly as the author wrote them. Record the rule as `ADR-064`, cover it with `AC-159`, and close `G-010`.

## Plan

**This change is plan-less.** It serves `G-010`, a goal no campaign adopted: `P-020` closed with the corpus-forgetting work and never took the index defect on. The fix is one behaviour with one criterion, so a campaign around it would be bookkeeping rather than sequencing.

## Why

`regenIndex` builds a row only for an artifact the block does not already reference. A row that already covers its target is kept byte for byte — which is deliberate and right for the description, because `ADR-046` makes that sentence author-owned and seeded rather than asserted. The badge on the same row is generated content that no rule ever handed to the author, and it never refreshes. Change an artifact's status and the row keeps the old value; run the regenerator and nothing happens, which teaches a reader that the row is authoritative exactly when it is stale.

This is not theoretical. In the last three changes the defect produced two blocking review findings — a plans index calling a completed campaign `active` and a goals index calling an accepted goal `proposed`, each contradicting a file the same change had just edited — and a third, where promoting 28 decision records left all 28 rows saying `inferred`. One drift had stood since August. `clue validate` cannot see any of it: it counts rows that say nothing about their artifact, not rows that say something untrue about it.

## Decision to record

**Regeneration owns the badge; the author owns the description.** That split is the whole of it, and `ADR-064` states it so an adopter's row behaves the same way.

Two boundaries keep the fix from trading one silent staleness for a silent overwrite:

- **Only a row that already carries a badge is refreshed.** A row whose author deleted the badge is left without one. Adding it back would overwrite an authored shape, and a row that claims nothing is never wrong.
- **Only a row covering exactly one artifact is refreshed.** A curated line covering several entries has no single artifact its badge belongs to, so the generator leaves it alone rather than guessing.

The badge is whatever `RowIdentity` reports, so a constraint keeps showing its enforcement rather than its status — `IDR-001`, and the reason a naive status sweep would have rewritten twenty-two correct rows into nonsense.

## Rejected: make `clue validate` fail on a stale badge

It would catch a corpus where nobody runs the regenerator, and `G-010` accepts either answer — refresh it, or report the disagreement. It is rejected because a new required check turns every adopter's next upgrade red for rows they did not write and cannot have known about, to report something the regenerator can simply fix. The generator is where the badge came from, so the generator is where it should stay current.
