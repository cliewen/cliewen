---
id: CH-164
type: change
status: open
links: [P-020, M-083, PDR-052, ADR-034, C-008, CAP-002]
title: Retire the simplification campaign's spent registers
---

# CH-164 — Retire the simplification campaign's spent registers

## Proposal

Serve P-020's M-083 by retiring AN-019, AN-020, and AN-021 — the transient registers P-013's simplification campaign produced, which no live decision cites.

P-019 shipped the mechanism for a spike to leave and pruned nothing. This is the first application, on the repository that wrote the rule. These three are the honest place to start: each was written to answer one question at one revision, each served a campaign that is `completed`, and none is the evidence base of a decision that still stands. AN-021 pins a statement register to revision `9a632f9` and says so in its own evidence boundary; AN-019 audits one candidate's post-trim state; AN-020 answers a placement question in the vocabulary of the `plain` and `light` tiers, which the methodology has since retired — so it does not merely repeat what the corpus says, it describes a corpus that no longer exists.

Each retirement declares its `carried-by:` artifacts before the file is deleted, so the claim that its findings survived is written down and checkable rather than asserted in a commit message. Deletion is the retirement ([ADR-034](../../docs/decisions/ADR-034-retirement-is-deletion.md)); Git history is the archive, and no document is emptied into a stub.

The inbound references are P-013 and AN-023. P-013 is `completed` and therefore frozen under [C-008](../../docs/constraints/C-008-completed-plans-immutable.md); the only edit it accepts is a link-target repair, the same narrow allowance P-016's decision-log retirement used. AN-023 is an ordinary live artifact whose rows are corrected.

## Scope boundary

This change retires three analyses and no others. AN-018 and AN-022 are deliberately excluded: live decisions cite them as their evidence base, and whether a decision's cited evidence may be deleted while the decision still stands is a question this change does not answer. It belongs to M-084 and needs a decision record, not a judgement made in passing while retiring something else.

Nothing is hollowed out, no completed plan changes except a link target, and no decision loses a citation it rests on.
