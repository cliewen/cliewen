---
id: PDR-026
type: decision
status: verified
links: [P-011, P-012, P-013, PDR-025, C-013, G-001]
title: A campaign closes on re-derived gate status, not on its own evidence column
author: agent
accepted-by: Flemming N. Larsen (2026-08-08, conversation)
---

# PDR-026 — A campaign closes on re-derived gate status, not on its own evidence column

## Context and problem statement

A plan's milestone table carries an evidence column, and each implementing change fills its own row in its own merge digest. That column is therefore written by the same agent that did the work, in the change that did it, and never read back against the corpus afterwards. When a campaign exists to answer an external assessment, the question at closing time is not whether every row says `done` — it is whether the assessment's gaps are actually closed, which is a different question that nothing currently asks.

The two can diverge in a specific and predictable way. A gap phrased as a judgment ("this migration hides real coverage debt behind a bulk `@draft`") is naturally answered with a machine check on form ("a draft disposition must carry a justification"), because form is what a deterministic judge can hold. The row is then true and the gap is not closed: the check counts sentences, and the assessment was about whether the sentences are honest. A request the campaign deliberately declined can also be recorded as satisfied, because the milestone that declined it is `done`. Neither divergence is visible from the table.

This lands on [G-001](../goals/G-001-verifiable-thread.md). A thread whose last link is a self-reported status is not verifiable in the sense the goal means.

## Decision outcome

**A campaign closes on gate status re-derived from the corpus and the tool, not on its own evidence column.** Before a plan moves to `completed` against an external assessment, each gap that assessment raised is checked again from the artifacts and commands as they now stand, and each one lands in exactly one of three states: closed with a named mechanism and its failure-path evidence; closed as a *declined request*, recording explicitly what was not granted and what it costs the adopter; or open, as a milestone in the successor campaign. A row saying `done` is the beginning of that check, never its result.

**A declined request is recorded as declined.** Where a campaign deliberately refuses what an assessment asked for — because the refusal is the better design — the refusal and its cost are stated. Silently recording a declined request as a satisfied one is the failure this decision exists to prevent, and it is worse than the original gap, because it removes the reason to look again.

**P-012 finishes the brownfield migration gap on re-derived evidence, and simplification moves to P-013.** This supersedes only [PDR-025](PDR-025-brownfield-migration-precedes-simplification.md)'s clause naming P-012 as the simplification campaign; its reasoning that migration credibility precedes simplification is not weakened by this decision but extended by it, since the evidence that closes the migration boundary must itself be re-derived before the boundary can be handed off. Every other clause of PDR-025 stands.

Milestones that change what `clue validate` or `clue parity` asserts remain core-adjacent under [C-013](../constraints/C-013-core-changes-need-decision.md) and carry their own decision records in their own changes; this decision authorizes the campaign, not its mechanisms.

## Rejected: treat a completed milestone table as closure

It is the cheapest possible reading and it is the one that fails. The table records what each change believed it had done at the moment it did it, which is exactly the claim under review. Accepting it as the answer makes the campaign its own judge.

## Rejected: reopen P-011 rather than open a successor

[C-008](../constraints/C-008-completed-plans-immutable.md) freezes a completed plan, and the milestones themselves were genuinely delivered — the mechanisms exist and their evidence holds at the scope they claim. What remains is different work at a different scope, and a campaign that can be reopened whenever its results are re-read is not a boundary at all.

## Rejected: fold the remaining migration work into simplification

A campaign carrying both would let either half postpone the other, which is the failure mode PDR-025 was written to prevent. It also mixes work that adds machinery with work whose purpose is to remove it.

## Carrier

P-011's terminal state, P-012's scope and milestones, P-013's deferred scope, PDR-025's superseded clause, the plans index, and this record carry the decision.
