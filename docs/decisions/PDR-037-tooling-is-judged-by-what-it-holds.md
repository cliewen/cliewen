---
id: PDR-037
type: decision
status: verified
links: [P-013, PDR-013, PDR-029, ARCH-003, AN-022, C-013]
title: Tooling is judged by whether removing it hands work back to a human
author: agent
accepted-by: Flemming N. Larsen (2026-08-08, conversation)
---

# PDR-037 — Tooling is judged by what it holds, not by what the core names

## Context and problem statement

[PDR-013](PDR-013-explicit-core-red-line.md) supplies *does the core need it?* and [PDR-029](PDR-029-simplification-tests-by-surface.md) assigns it the surface it governs: rules, artifact types, required fields, commands, and checks. [AN-022](../analysis/AN-022-remaining-surface-scored.md) applied it to those populations for the first time and the criterion under-determined on four of the five.

The reason is structural rather than a defect in the populations. [ARCH-003](../architecture/core.md)'s core is three *guarantees* — the verifiable thread, the human merge boundary, the deterministic judge. A command, a check, a field, and a type are *means*. Asked of `clue validate` the criterion answers yes, and asked of every other command it answers no, because the core needs the judge and needs nothing else by name. `clue id next` is a plain example: the core does not need it, and it exists so a deleted artifact's identity is never re-minted, which is the criterion-identity guarantee the thread depends on. One step removed, and invisible.

PDR-029 already recorded the mirror failure on carrier prose, where the criterion passes almost everything because nearly every rule a skill states genuinely is needed. The same repair applies here, and the campaign that discovered both failures should not report either as a result.

## Decision outcome

**On the tooling surface — commands, checks, required fields, artifact types, and rules that hold a mechanism rather than a meaning — the operative question is: *does removing it move an obligation from a machine back to a human?***

- **A yes is a keep.** The obligation still exists; only its holder changed, and the holder it moved to forgets. That is the failure Cliewen was built from: a ritual that depends on a person remembering will eventually be skipped, by everyone, eventually.
- **A no is a removal candidate, and the standing preference is to remove it.** Where nothing moves back to a human and the item does not contribute meaningfully, it goes. A field nothing populates, a check nothing exercises, a rule nobody has needed — each is read cost with no holder, and the burden falls on retention rather than removal.
- **The two tests are not ranked; the surface decides.** *Does the core need it?* remains the test of record for anything whose existence changes what the method **obliges** — a rule reaching into meaning, a protected boundary, the red line itself. This test governs anything whose existence changes what a machine **holds**. Where an item is both, the core question is decisive and answering it *yes* ends the enquiry.
- **Neither test authorizes a removal on its own.** A removal that changes what the verifiable thread connects, what a merge accepts, or what a green `clue validate` asserts is a core meaning change under [C-013](../constraints/C-013-core-changes-need-decision.md) and carries its own record and human acceptance — and removing a check is exactly that, because green afterwards asserts less. [PDR-029](PDR-029-simplification-tests-by-surface.md)'s shared-memory clause also holds unchanged: reducing what the corpus can remember is not simplification, whatever either test returns.

The second test is not new reasoning. It reconstructs every decision the project has actually made: each command in `clue` exists because a finishing step was being skipped, each generated index exists because a hand-maintained one went stale, and each check exists because an obligation that lived in prose was not kept. Writing it down makes that reasoning available to a campaign instead of leaving it in the history of nine separate decisions.

**What this changes about P-013's result.** The campaign measured its surface and found that almost nothing should be removed. That is the finding, stated plainly, and it is not a disappointment: the criterion it began with would have reported either that most of Cliewen's tooling is unnecessary — false, and quotable — or that everything passed, which is the fourth deferral reached through a new test rather than the old one. A campaign may honestly conclude that a surface is already close to minimal, provided it says which test it applied and what it declined.

**Carrier inventory:** the amendment blockquote at the head of [PDR-029](PDR-029-simplification-tests-by-surface.md); [P-013](../plans/P-013-simplification.md)'s M-064 prose and the campaign's statement of its own measure; and this record. No skill, hub, guide, or CLI carrier states the campaign's tests, so none is affected — the tests govern how a simplification campaign judges, not how an ordinary change is made.

### Rejected: apply *does the core need it?* to the tooling surface as written

It is what PDR-029 assigned and AN-022 showed it returns *no* for eleven of thirteen commands, fourteen of twenty-three checks, and eleven of seventeen constraints, in nearly every case for something that should stay. A test whose answer must be overridden almost every time it is applied is not a boundary.

### Rejected: fold the second test into PDR-013 as a clarification of the criterion

It is not a clarification. *Does the core need it?* and *does removing it hand work back to a human?* can disagree, and they disagree in a specific direction — the second keeps things the first would discard. Presenting that as a reading of PDR-013 would make the red line look elastic, and the red line's value is that it is not.

### Rejected: measure the tooling surface by removals, artifacts, or lines

Refused for the same reason PDR-029 refused it for prose, and AN-022 is the demonstration: the honest result of the scoring is that almost nothing should be removed, so a campaign scored on removals would have to manufacture some.
