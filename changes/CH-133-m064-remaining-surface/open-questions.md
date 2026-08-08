---
id: CH-133-questions
type: open-questions
status: open
links: [CH-133, P-013, M-064]
title: Open questions for CH-133
---

# CH-133 open questions

Six, all blocking. Q-01 through Q-03 come from [AN-021](../../docs/analysis/AN-021-remaining-carrier-register.md) and are about carriers; Q-04 through Q-06 come from [AN-022](../../docs/analysis/AN-022-remaining-surface-scored.md) and are about the method's own shape. Each names the statement or candidate, what it traces to or fails to, what removing or retaining it costs, and the judgement required. None is answerable by an agent under [PDR-029](../../docs/decisions/PDR-029-simplification-tests-by-surface.md).

The full statement of each is in the analysis that raised it; what follows is the question a human has to answer.

## Q-01 — `guide/design.md` states two rules the method has changed

GD-71 says a fix invalidates the review pass and the loop ends only on a clean review; C-017 and PDR-035 decided otherwise one change ago. GD-74 states "small deltas" as a named principle; that rule was withdrawn as uncheckable and removed from `clue-delta` by CH-132. Both are carrier-inventory misses under PDR-019, whose own carrier list names public guidance.

GD-71's repair is mechanical. GD-74's is not: the rule is withdrawn but its reason — two large parallel changes can be textually mergeable and semantically incompatible, and no tool notices — is why the one-change-in-flight and reconcile-and-recheck rules exist, and those survive.

**Question.** May `design.md` keep that argument as an explanation of the rules that remain, or does a withdrawn rule take its rationale with it?

## Q-02 — the poor-fit conditions trace to nothing

`guide/adoption.md` states five conditions under which Cliewen should not be adopted, in the imperative, and `guide/operations.md` links to them as its fallback. Nothing in the corpus states them. Four of the five are about the adopter's situation rather than Cliewen's design.

Removing them costs the method its only published statement of its own limits, on the page an evaluator reads before adopting. Retaining them untraced leaves five load-bearing conditions no change is obliged to keep true.

**Question.** Are these a real rule the corpus should state — most naturally a constraint, or prose in ARCH-003's periphery section — or are they editorial advice that should say so and stop being written as instructions?

## Q-03 — `CONTRIBUTING.md` contradicts C-017 on the page that cites it

CTB-44…CTB-49 carry the bounded loop and the intrinsic severity model. CTB-52, four paragraphs later, says every substantive fix invalidates the earlier clean pass. Under C-017 an advisory repair is substantive and does not invalidate. The generated `clue-verify` carrier states the same clause and is rescued by its following sentence; `CONTRIBUTING.md` carries the clause without it.

The direction looks settled and the repair looks free. It is asked anyway because a cheap repair is not evidence that it is the right one, and resolving a conflict silently is a methodology decision without a human.

**Question.** Confirm that C-017 governs and CTB-52 is redrafted to its language — or state why the contributor's rule should be stricter than the agent's.

## Q-04 — the tooling surface needs a second test

*Does the core need it?* returns *no* for eleven of thirteen commands, fourteen of twenty-three checks, six of fourteen types, and eleven of seventeen constraints, and in nearly every case the item should stay. PDR-029 recorded the mirror failure on carrier prose and repaired it by giving that surface its own test.

The proposed second test for tooling: **does removing it move an obligation from a machine back to a human?** It is checkable, and it reconstructs every decision the project actually made — `clue` exists because rituals that depend on remembering get skipped.

Cost of adopting it: a third test, and a campaign that must say which test applies where. Cost of not adopting it: P-013 either reports that most of Cliewen's tooling is unnecessary, which is false and would be quoted, or reports that everything passed, which is the fourth deferral by a new route.

**Question.** Add the second test for the tooling surface as an amendment to PDR-029? And does the answer change what P-013 claims to have measured?

## Q-05 — C-003 governs a transient artifact

C-003 obliges an agent to tick a task the moment it completes, on a file the digest deletes and that never reaches `main`. Its enforcement is `partial`: the check verifies a skipped-task line's shape, not the timing. It is the one constraint the scoring would put forward on the criterion alone.

Removing it costs handoff legibility — a batch-ticked checklist does not say which task was in flight when work stopped, and the workspace is what a second agent reads.

**Question.** Does that handoff value justify a constraint, or does the rule belong in `clue-delta` as procedure with no register entry?

## Q-06 — `supersedes:` is both the strongest removal candidate and pattern C's only partial answer

The field and its check ship. No artifact in the corpus carries it, because ADR-034 binds it to deletion while the supersession that actually happens here — nine of eighty-four decisions — is of records that survive and must survive. Nothing answers "what was downstream of this decision" except reading.

Three answers, not equivalent:

- **Widen it** — let a surviving superseded record carry the edge. Closes pattern C's residue; adds an obligation to every superseding change; needs a milestone that owns a mechanism, and P-013 has none.
- **Leave it** — keep it for retirement only, and decline pattern C's residue with the cost stated: prose supersession stays unqueryable.
- **Remove it** — accept that a future answer needs new machinery, against PDR-029's rule that reducing what the corpus can remember is not simplification.

**Question.** Which, and if *widen*, which milestone builds it? Pattern C cannot be reported closed or declined without this answer, so M-066 cannot close the campaign without it either.
