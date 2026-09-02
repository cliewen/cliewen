---
id: CH-166
type: change
links: [P-020, M-085, PDR-052, ADR-034, PDR-030]
title: A spike a standing record cites is not reported as spent
---

# CH-166 — A spike a standing record cites is not reported as spent

## What

Add the third condition to the spent-analysis derivation: an analysis that a live decision or constraint names in its `links` is not reported, however complete its plans and however resolved its carriers. Record the rule as `PDR-053`, cover it with `AC-158`, carry it on the shipped analysis guidance, and close `P-020`'s `M-085` by recording the disposition of the two spikes the milestone names.

## Why

`P-020`'s `M-085` asks a question the campaign deliberately refused to answer in passing: may a decision's cited evidence be deleted while the decision still stands?

The answer this change records is no. A decision record is compact by design — it states what was decided and the reasoning in a few hundred words, and everything that made the reasoning checkable lives in the spike it came out of. Delete the spike and the decision is still readable but no longer reviewable: a later reader can see what was concluded and cannot see what it was concluded from. That is a worse corpus than a slightly larger one, and it is the failure mode the whole verifiable thread exists to prevent.

The rule is also already the one this repository acted on. `AN-019`, `AN-020`, and `AN-021` were retired in `CH-165` on the stated ground that no live decision cited any of them; every other analysis in the corpus has a decision or constraint pointing at it. This change makes the ground the tool applies rather than a sentence in a merge digest.

Constraints count alongside decisions for the same reason: both are standing rules that outlive the campaign that produced them, and both are read by someone deciding whether the rule still holds.

## The lever this leaves the human

A `links` entry is how a Cliewen record cites. Where a decision's link to a spike turns out to be provenance rather than evidence — this came out of that, and nothing more — the correct repair is to remove the citation from the live decision in a reviewed change, which then makes the spike eligible on the next preview. A live decision is editable; that is the difference between it and a frozen plan. Pruning further therefore stays possible and stays an explicit human act, rather than an agent judging per file which citations are load-bearing.

## What this does not do

It does not retire anything. It does not touch `clue validate`: an analysis a decision cites is not an invalid corpus, and the judge reads state rather than judgment (`ADR-044`). It writes no file during migration.
