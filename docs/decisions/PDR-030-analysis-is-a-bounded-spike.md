---
id: PDR-030
type: decision
status: verified
links: [P-013, PDR-029, PDR-006, ADR-035, C-011, G-001, AN-018]
title: Analysis is a bounded spike that ends in a findings document with a named consumer
author: agent
accepted-by: Flemming N. Larsen (2026-08-08, conversation)
---

# PDR-030 — Analysis is a bounded spike that ends in a findings document with a named consumer

## Context and problem statement

Cliewen's `clue-analysis` skill has stated the shape of a spike since the corpus was seeded: run it before planning or implementing, name the risk in one sentence, treat the investigation as disposable and the findings as durable, end in a findings document under `/docs/analysis/`, record what was rejected, and feed the result to a plan or a change. Seventeen findings documents have been produced under those rules without exception.

None of them was ever decided. The rules entered the skill from the Foundation Document, which is a frozen analysis whose own banner states that the corpus wins where the two disagree, and no decision, constraint, goal, or criterion ever restated them. Under [PDR-029](PDR-029-simplification-tests-by-surface.md)'s carrier test they trace to nothing, which makes the skill that governs how this project writes its shared memory the least defensible carrier it ships — a reviewer asking *why do you have this statement?* would receive "the founding document said so", which PDR-029 rules out by name.

The rules are also not separable. Remove "a spike is disposable" and "its findings are not" says nothing; remove "end in a findings document" and the disposability rule licenses losing the work; remove "feed it to a consumer" and the corpus accumulates analysis nobody asked for. A record per sentence would state ten rules that each depend on the other nine.

## Decision outcome

**Analysis is a bounded investigation with a stated risk, a disposable body of work, a durable findings document, and a named consumer. All four are one rule, and this record is its trace.**

*It runs before planning or implementing, not beside them.* Analysis exists to retire the largest risk or unknown first, so it precedes the work whose shape that risk would otherwise determine. Analysis run afterwards documents a decision already taken and is a different activity.

*The risk is named in one sentence, and failing to name it is the first finding.* An investigation whose subject cannot be stated in a sentence has not been scoped, and reporting that is more useful than proceeding.

*The spike is disposable; the findings are not.* A prototype, a measurement, a literature scan — the artifact of the investigation is thrown away. What survives is the findings document, and that asymmetry is deliberate: a spike kept alive becomes a second implementation nobody maintains, and findings thrown away make the same investigation happen twice.

*Every spike ends in a findings document under `/docs/analysis/`, and that document records what was rejected.* Discarded options are half of why the system looks as it does. A findings document stating only what was chosen leaves the next reader free to re-propose everything that was already examined and refused, which is the cost this clause exists to prevent.

*Analysis has a named consumer or it is not written.* The consumer is a plan or a change — a milestone that will use the findings, or the change that will act on them. Analysis with no consumer is not cheap documentation; it is an artifact that enters the mandatory read path forever and repays nothing, and this is the clause the campaign that produced this record cares about most.

This record does not change any of those rules. It states them, so that the skill has something a reader can open.

## Rejected: one decision record per rule

Ten records, each stating a sentence that is meaningless without the other nine. The register that found this gap counted ten untraceable statements, and the tempting repair is ten traces. It would produce a corpus where reversing any one clause leaves the remaining nine internally contradictory, and where the reason the rules hold together — that a bounded investigation, a disposable body, a durable finding, and a real consumer are one design — is stated nowhere.

## Rejected: promote the Foundation Document to a valid trace

It would close this gap and eleven like it in one edit, and it is wrong for the reason the document itself states: it is frozen, it is never edited, and where it and the corpus disagree the corpus wins. An artifact that cannot be amended cannot carry a live rule, because the first time the rule needs to change there is nowhere to change it.

## Rejected: delete the untraceable statements

The mechanical reading of the carrier test, and [PDR-029](PDR-029-simplification-tests-by-surface.md) forbids it explicitly: a statement that traces to nothing may be a real rule nobody recorded, in which case the repair is to write the missing decision. Seventeen findings documents produced under these rules are the evidence that this is that case.

## Carrier

The `clue-analysis` skill's canonical source and its generated trees state the rules; this record states why they hold. [PDR-006](PDR-006-decision-records-are-typed.md) carries the neighbouring rule that a rejected alternative which is itself a decision gets its own rejected decision record — that clause is about decision records rather than about analysis, and lives there rather than here. [ADR-035](ADR-035-bounded-provenance-and-reality-edges.md) continues to govern the `reality: contradicted` marker unchanged.
