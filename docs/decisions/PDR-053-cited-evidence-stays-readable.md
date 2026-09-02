---
id: PDR-053
type: decision
status: inferred
links: [P-020, M-085, PDR-052, ADR-034, PDR-030, ADR-044, CAP-001]
title: A spike a standing record cites stays readable while the citation stands
binds: adopter
author: agent
accepted-by: []
---

# PDR-053 — A cited spike stays readable

## Context

[PDR-052](PDR-052-a-spent-analysis-is-reported-not-retained.md) reports a spike as spent when every plan it serves is complete and it names a durable artifact carrying its findings. Neither condition asks what still depends on it.

A decision record is compact on purpose. It states what was decided and why in a few hundred words, and the material that made the reasoning checkable — the population that was measured, the revision it was measured at, the cases that were tried and rejected — lives in the spike the decision came out of. A campaign completing says nothing about whether the rules it produced have stopped standing.

## Decision

**An analysis that a live decision or constraint names in `links` is not reported as spent, whatever its plans and carriers say.** Both are standing rules that outlive the campaign that produced them, and both are read by someone deciding whether the rule still holds. Deleting the evidence under a standing rule leaves it readable and no longer reviewable: a later reader can see what was concluded and cannot see what it was concluded from. That is the failure the verifiable thread exists to prevent, and a slightly larger corpus is the cheaper side of the trade.

A `links` entry is how a Cliewen record cites, so the condition needs no new field and no judgement per file. Where a decision's link to a spike turns out to be provenance rather than evidence — this came out of that, and nothing more — the repair is to remove the citation from the live decision in a reviewed change, after which the spike is reported on the next preview. A live decision is editable; that is exactly what distinguishes it from the frozen plan of [ADR-063](ADR-063-a-frozen-plan-links-are-historical.md). Pruning past this point stays possible and stays an explicit human act.

Only decisions and constraints gate. A plan does not: an active plan's link keeps the spike unspent through PDR-052's own completion test, and a completed plan's links are history rather than a live dependency. An analysis citing another analysis does not gate either, because a spike is not a standing rule.

This does not reach [`clue validate`](../capabilities/CAP-002-validate/README.md). A cited analysis is not an invalid corpus, and the judge reads state rather than judgment ([ADR-044](ADR-044-the-judge-reads-state.md)).

**Accepted cost:** in a corpus where decisions routinely cite the spikes they came from, this leaves most of the analysis population in place. That is the intended result and not a shortfall — the population that can leave without cost is the transient register, and this rule is what keeps the distinction from having to be re-argued at each retirement.

## Carrier

The citation condition in `internal/corpus`, its criterion, this record, and the shipped analysis guidance in `internal/skills/source/skills/clue-analysis.md.tmpl` and `internal/scaffold/templates/docs/analysis/README.md` stating what keeps a spike from being reported.
