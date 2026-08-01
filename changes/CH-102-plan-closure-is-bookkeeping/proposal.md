---
id: CH-102
type: change
status: open
links: []
title: Closing a plan is bookkeeping, and every carrier says so
---

# CH-102 — Closing a plan is bookkeeping, and every carrier says so

**Plan-less.** P-009 is closed and no successor is designated.

## What is wrong

Six live carriers state that plan bookkeeping rides in the merge digest. Not one of them says whether *closing the plan* is part of that bookkeeping, and the practice since P-007 is that it is: the change completing the last milestone also sets the plan `completed`.

The gap is not theoretical. CH-095 ticked M-044, left P-009 `active`, and argued in its own commit message that closure was a scheduling decision an implementing change may not take. P-008 had already rejected exactly that by closing itself and naming its successor in one digest. CH-099 then existed only to repair the omission, and its first attempt put the correction in `docs/plans/README.md` — a file no adopter ever reads, since `clue init` scaffolds a different one.

Every carrier is honest and every carrier is incomplete in the same way, which is why the rule survived being wrong.

## What changes

`clue-plan` is the normative carrier and gains the rule in full: closing the plan belongs in the same digest as the last milestone, a campaign is over the moment its last milestone is evidenced, and a successor is designated there when one is decided — its absence never holding the closure open.

The five other carriers gain the clause that removes the ambiguity, sized to what each is for. `AGENTS.md` and its scaffolded twin route; the two plan READMEs describe the folder; the public guide narrates the loop. None of them restates the method.

## What does not change

The rule itself. This is not a new obligation — P-007 and P-008 already followed it, and the decision log already records it as their precedent. Nothing here decides that plans close this way; it writes down that they do.

No plan is reopened. Completed plans are pinned history and are not edited to match the sharper wording.

## Reversal cost

Cheap and local: reverting one clause in each carrier restores the previous text exactly. A log row, not an ADR.
