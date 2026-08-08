---
id: PDR-038
type: decision
status: verified
links: [P-013, PDR-026, PDR-029, ADR-034, AN-008, AN-022]
title: A surviving superseded record carries no machine edge, declined with its cost stated
author: agent
accepted-by: Flemming N. Larsen (2026-08-08, conversation)
---

# PDR-038 — Supersession of a surviving record stays prose

## Context and problem statement

[AN-008](../analysis/AN-008-methodology-critiques.md)'s pattern C said the corpus only accumulates, in four claims, and three campaigns carried the paragraph forward without examining it. [AN-022](../analysis/AN-022-remaining-surface-scored.md) re-derived each claim at head. Retirement now removes an artifact from the read path by deleting it ([ADR-034](ADR-034-retirement-is-deletion.md)), the edge back from reality ships as an incident analysis, and the claim that unsigned decisions can only accumulate is false in this repository — one of eighty-four carries `status: inferred`.

One residue survives and it is narrower than the paragraph. `supersedes:` and its check exist and **no artifact in the corpus uses them**, because ADR-034 binds the field to deletion while the supersession that actually happens is of records that survive and must survive: a superseded decision is history a later reader has to be able to read. Nine of eighty-four decision index rows claim a supersession or amendment, all of it in prose — an index row, a successor's body, an amendment blockquote. None is an edge, and nothing answers *what depended on this decision* except reading.

That leaves one field which is simultaneously the strongest removal candidate on the whole scored surface and the only partial answer pattern C has.

## Decision outcome

**The residue is declined, with its cost stated, and the field stays as it is.**

- **`supersedes:` keeps its current meaning: it carries the pointer forward when an artifact is deleted.** It is not widened to a superseded record that survives.
- **It is not removed, although it is unused.** [PDR-037](PDR-037-tooling-is-judged-by-what-it-holds.md)'s test asks whether removal hands work back to a human, and here the answer is subtler than usual: nothing moves today, because nothing uses the field. What removal would take away is the only place a future change could record the edge without new machinery — and under [PDR-029](PDR-029-simplification-tests-by-surface.md) reducing what the corpus *can* remember is not simplification. An unused field whose check is already written costs a reader one line of a lifecycle rule; removing it costs a future answer its cheap route.
- **The cost of declining is stated rather than implied.** Supersession recorded in prose cannot be queried. Nobody can ask which decisions were downstream of a reversed one, and the answer stays a reading task over nine records today and more later. `clue context` follows outgoing links only and deliberately, so even a machine-visible edge would need a reverse walk that the narrow-reading obligation argues against — the mechanism this decision declines is therefore larger than one field.
- **Widening it is routed, not dropped.** Pattern C's residue becomes a named door for [P-013](../plans/P-013-simplification.md)'s successor campaign, stated in P-013's own prose so the deferral has a destination rather than a sentence. It is not built inside a simplification campaign: it adds machinery, and [PDR-026](PDR-026-campaigns-close-on-re-derived-evidence.md) requires an addition inside such a campaign to argue for itself.
- **Pattern C is otherwise determined and stops being carried forward as a paragraph.** Three claims answered, one false at head and unproven in general, one declined here. A campaign that names pattern C again names this record and the door, not AN-008's four-part sentence.

**Carrier inventory:** this record; [P-013](../plans/P-013-simplification.md)'s prose, which states the determination and the successor's door; and [AN-022](../analysis/AN-022-remaining-surface-scored.md), which holds the re-derived evidence. ADR-034 is untouched — its retirement rule and the `supersedes:` semantics it decided are exactly what this decision leaves in place. No skill, hub, guide, or CLI carrier states a rule about superseding a surviving record, so none is affected.

### Rejected: widen `supersedes:` to a superseded record that survives

It closes the residue and it is the answer this decision would give if the campaign had a milestone that owned a mechanism. It does not: every superseding change would gain an obligation, the judge would need a rule distinguishing a live superseded record from a stale one, and the reverse question — what was downstream — still needs a walk the tooling deliberately does not do. Adding all of that to close a gap nine records wide, inside the campaign whose subject is simplification, inverts the campaign.

### Rejected: remove the field and its check as unused

The candidate is real and the reasoning against it is the shared-memory clause, not sentiment. The field is the cheap route to an answer somebody will want; deleting it makes that answer expensive without making anything simpler today, because the check is written and passes silently.

### Rejected: close pattern C because three of its five claims are answered

Three campaigns carried the paragraph forward unexamined. Closing it unexamined is the same failure with the opposite sign, and PDR-026 exists to prevent exactly that trade.
