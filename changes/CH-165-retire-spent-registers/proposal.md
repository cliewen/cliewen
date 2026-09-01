---
id: CH-165
type: change
status: open
links: [P-020, M-084, PDR-052, ADR-034, ADR-063, CAP-002]
title: Retire the simplification campaign's spent registers
---

# CH-165 — Retire the simplification campaign's spent registers

## Proposal

Serve P-020's M-084 by retiring AN-019, AN-020, and AN-021 — the transient registers P-013's simplification campaign produced, which no live decision cites.

This is the first artifact P-019's forgetting mechanism actually removes. CH-164 attempted it, found that a completed plan's link made retirement impossible, and settled that in [ADR-063](../../docs/decisions/ADR-063-a-frozen-plan-links-are-historical.md) instead. P-013 links all three of these spikes, so none of them could be retired before that fix; now they can, and P-013 needs no edit at all — its links stay exactly as they are, recording what the campaign referenced while it ran.

These three are the honest place to start. Each answered one question at one revision for a campaign that is `completed`. AN-021 pins its register to revision `9a632f9` in its own evidence boundary. AN-020 answers a placement question in the vocabulary of the `plain` and `light` tiers the methodology has since retired, so it does not merely repeat what the corpus says — it describes a corpus that no longer exists. AN-019 audits one candidate's post-trim state, and the boundary question it left open was closed by P-013's own M-067.

Each spike declares its `carried-by:` artifacts **in its own commit**, before anything is deleted, with `clue migrate`'s `MIG-013` report captured against that commit. CH-164's reverted attempt marked its milestone `done` on exactly this evidence while adding the field and deleting the files in one uncommitted step, so no commit ever held a state a reader could re-run. True and unverifiable is, for an evidence column, the same as false.

## Scope boundary

Deletion is the retirement ([ADR-034](../../docs/decisions/ADR-034-retirement-is-deletion.md)); Git history is the archive and nothing is hollowed into a stub, which would keep a document's index row and read cost while removing the only part with value.

P-013 is not edited. AN-023's rows naming two of these spikes are left intact: they are plain text recording what one inspection covered at a pinned revision, they break nothing, and rewriting them would make a pinned record a false account of itself. AN-022 is an active artifact and repoints like anything else.

AN-018 and AN-022 stay. Live decisions cite them as their evidence base, and whether a decision's cited evidence may be deleted while the decision stands belongs to M-085 with a decision record of its own — not to a judgement made in passing here.
