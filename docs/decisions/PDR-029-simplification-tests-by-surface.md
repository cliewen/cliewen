---
id: PDR-029
type: decision
status: inferred
links: [P-013, PDR-013, PDR-025, PDR-026, ADR-021, ADR-034, AN-006, AN-008, AN-010, AN-012, G-001, C-013]
title: Simplification is judged by two tests, chosen by surface
author: agent
accepted-by: []
---

# PDR-029 — Simplification is judged by two tests, chosen by surface

## Context and problem statement

[PDR-013](PDR-013-explicit-core-red-line.md) supplies simplification's criterion — *does the core need it?* — and records why simplification kept stalling without one: every debate was argued case by case, and protection with no boundary became uniform friction. That criterion is correct for what it was written about. It decides whether a rule, an artifact type, a required field, a command, or a check earns its place, and it decides it without appeal to taste.

It returns the wrong answer on one surface, and that surface is the one an adopter reads most closely. Applied to Cliewen's own methodology carriers — the skills, the routing hub, the guide — nearly every rule stated there *is* needed by the core, so the test passes almost everything and concludes there is nothing to simplify. What makes those carriers expensive is not unnecessary rules. It is the same necessary rule stated in several places; a rule compressed so far that a reader cannot tell whether it has been satisfied; and statements that are individually true but that nobody can say the purpose of.

The last of those is the one the criterion is least able to see, and it is the one that matters most, because the skills are read adversarially. A reviewer examining a skill closely does not ask whether the file is short. They ask, of a particular sentence, *why do you have this statement?* — and a methodology whose own instructions cannot answer that question has lost the argument before any of its rules are evaluated. Every statement in a carrier is therefore accountable individually, and a statement that survives only because it was already written is not surviving on merit.

None of the three is a matter of taste, and the first two are outright countable. Shared fragments are rendered into more than one skill, and because the change loop directs an agent from one skill to another that renders the same fragments — while the routing hub has already stated several of them — an ordinary full change re-reads a substantial fraction of the instruction it has already read. Compression's cost shows in the same files: a verification item stating five independent conditions in one sentence is offered as a single checkbox, which no reader can honestly tick when four conditions hold and the fifth does not. The third is settled per statement rather than counted, by writing down what each one is for and seeing which answers survive contact with a reader who did not write them. The arithmetic and the register both belong to the inventory P-013's M-062 performs, not to this record; what is enduring here is that all three failure modes exist and that none is visible to *does the core need it?*

A campaign judged only by PDR-013 would not see any of that, and would defer simplification a fourth time on the grounds that everything passed.

## Decision outcome

**Simplification is judged by two tests, and the surface decides which one applies.**

- ***Does the core need it?*** governs rules, artifact types, required fields, commands, checks, and anything else whose existence changes what the method obliges. [PDR-013](PDR-013-explicit-core-red-line.md) is unchanged and remains the test of record for this surface.
- ***Does every statement earn its place, state its rule once, and admit checking?*** governs carrier prose: skills, the routing hub, contributor and public guidance, and CLI text. The unit is the individual statement, not the file, and it must pass all three conditions. It **earns its place** when a reader who knows what Cliewen is about can say why the statement is there without being told — the justification is the method itself, not a remembered incident, not a rule that no longer exists, and not the fact that the sentence was already written. It **states its rule once** when no live carrier already states the same rule in the same reading path. It **admits checking** when a reader can determine whether it has been satisfied.

  The first condition is the operative one and the hardest to pass. A statement can be unique in the corpus and perfectly checkable and still fail it, by being a true sentence nobody needs. The practical form of the test is the question a human reviewer actually asks when reading a skill closely — *why do you have this statement?* — and a statement passes only when the answer is evident from what Cliewen is for, rather than requiring its history to be recounted.

**Word count is not the measure of success, and the second test will sometimes lengthen a passage.** Splitting an unreadable five-condition sentence into five checkable ones adds words and is a pass, not a regression. The campaign's measure is the set of statements that survive: each earning its place, each stating its rule once in a given reading path, each checkable — never the size of the diff. A campaign that scored itself on words removed would repair the duplication and make the compression worse.

**Carrier prose and corpus artifacts are not equally safe to compact.** Carrier prose is instruction: it may be rewritten, merged, or moved freely so long as every rule it carries survives somewhere in the reading path. The corpus is the memory that multiple humans and coding agents share, and its obligation is that a different reader, later, recovers the same meaning. Reducing what the corpus can remember, or making a past decision harder to reconstruct, is not simplification, and passing *does the core need it?* does not by itself license it. Within this campaign, any such removal carries its own decision record and a guard proving what still resolves; this clause states how P-013 judges its own candidates and adds no obligation to changes outside it, which continue to retire artifacts under [ADR-034](ADR-034-retirement-is-deletion.md) and the `supersedes:` mechanism unchanged.

**This decision authorizes the tests, not any removal.** Every milestone that changes what a command asserts, or that touches the verifiable thread, the human merge boundary, or the judge, remains core-adjacent under [C-013](../constraints/C-013-core-changes-need-decision.md) and carries its own record in its own change.

## Rejected: apply *does the core need it?* to every surface

It is the single-criterion answer and it is the reason simplification has been deferred four times. On carrier prose the test passes nearly everything, because the rules genuinely are needed; the cost lives in how many times they are stated and how checkable each statement is, which the test cannot see. One criterion that returns "nothing to do" on the largest reader-facing surface is not a boundary, it is an exemption.

## Rejected: measure the campaign in words or artifacts removed

A count is tempting because it is deterministic, and it is precisely wrong here. It scores the duplication repair and the compression repair with opposite signs, so an agent optimizing it would delete the second copy of a rule and then compress the survivor until nobody can check it — arriving at a corpus that is smaller and less usable, which is the outcome [G-001](../goals/G-001-verifiable-thread.md) exists to prevent.

## Rejected: remove the duplication by making skills reference each other

[ADR-021](ADR-021-generated-standalone-skills.md) considered exactly this and rejected it, so that copying one skill folder yields complete instructions with no dependency on a repository layout outside itself. That reasoning is untouched here. The file-level repetition is deliberate and is not the cost this decision targets; what the second test measures is repetition inside a single reading path, which is a different quantity and does not require reopening ADR-021.

## Carrier

P-013's milestones and prose, and this record, carry the decision. Nothing else states it yet: the carriers this decision will be applied *to* are edited by P-013's own milestones, not by the change that records it.
