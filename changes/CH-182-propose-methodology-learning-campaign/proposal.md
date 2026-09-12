---
id: CH-182
type: change
status: open
links: [G-014, G-015, G-016, G-017, P-023, AN-024, VIS-001]
title: The methodology review's four priorities enter the corpus with an owner and an order
---

# CH-182 — The methodology review's four priorities enter the corpus with an owner and an order

[AN-024](../../docs/analysis/AN-024-methodology-review-acceptance-and-learning.md) reviewed the methodology and named four priorities: challenge ideas and plans before commitment, preserve and reuse practical learning, make acceptance an informed decision about usable behaviour, and evaluate the method's value so it can be simplified against evidence. It is a reviewer's recommendation, explicitly not accepted policy, and it allocated no identities. `clue next --all` reports no unfinished milestone — P-022 closed on 2026-09-11 naming no successor — so the four priorities currently live in one analysis file and nothing owns them.

This change puts them in the corpus. Each priority becomes a goal in the inbox, `status: proposed`, because a goal is where an idea enters and human promotion to `accepted` is the decision that admits it ([ADR-002](../../docs/decisions/README.md)). The first two, which AN-024 recommends as the initial delivery, also get an active campaign: [P-023](../../docs/plans/P-023-challenge-plans-and-retain-what-work-teaches.md). The remaining two wait in the inbox until that campaign reports, which is AN-024's own sequencing advice — later scope should follow the trials' findings rather than being committed before them.

P-023's hardest constraint is what it must not do. Every rule it proposes is a methodology rule, and a methodology rule written onto a carrier under `internal/skills/source/` reaches every adopter on the next release. So the campaign trials its rules as repository-binding conventions first and moves only the proven text onto the shipped carriers, at its last milestone. This is AN-024's "demonstrate before expanding the mechanism" read as a delivery order, and it keeps a retraction a local edit rather than a second full change against adopters. The campaign is also its own first test subject: the guides proposal that motivates priority 2 is exactly the kind of consequential commitment priority 1 says to challenge before making.

This change introduces no behaviour and adds no acceptance criterion. Criteria arrive with the milestones that change what the method obliges, under [CAP-006](../../docs/capabilities/CAP-006-collaborative-handoffs/README.md) or a new capability if [M-096] concludes one is needed.

## Identity note

`clue id next CH` first returned `CH-181`, which open draft PR #207 (`ch-181-fresh-context-orientation`) already held: that branch predates AN-024 on `main`, so `main`'s ledger had never seen its reservation. Identity allocation was local-only. Per the change loop, allocation moved to Git coordination (`clue id coordinate`, `.clue/id-coordination.yaml`) before this change took an identity, and this change is the later artifact, so it took a fresh one. PR #207 keeps CH-181.
