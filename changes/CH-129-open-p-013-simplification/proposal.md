---
id: CH-129
type: change
status: open
links: [P-013, PDR-013, PDR-025, PDR-026, AN-006, AN-008, AN-010, AN-012, AN-013, ADR-021]
title: Open P-013 with a milestone table and the two tests simplification is judged by
---

# CH-129 — Open P-013 with a milestone table and the two tests simplification is judged by

## Why

P-012 completed and designated [P-013](../../docs/plans/P-013-simplification.md) as its successor, but P-013 is `draft` and holds no milestones. Its own text says why that is deliberate: a `draft` plan does not hold milestone identities it has not decided, and deciding them is the work of the change that opens the campaign. Until that happens the campaign cannot be worked, and simplification stays where it has been for four campaigns.

P-013 also names its own first problem plainly: simplification has been deferred three times, so the shape it can be executed in has never existed. [PDR-013](../../docs/decisions/PDR-013-explicit-core-red-line.md) supplies a criterion — *does the core need it?* — that has never been applied systematically to a measured surface. The milestone table this change writes is that shape.

One thing became clear while sizing the surface, and it changes what the campaign must measure. The criterion PDR-013 supplies does not fit every surface simplification has to reach. Applied to Cliewen's own methodology carriers it returns the wrong answer, because nearly every rule stated in a skill *is* needed by the core; what makes the skill set expensive is not unnecessary rules but the same necessary rule stated in several places and, in places, compressed until it cannot be checked. A campaign that judged carrier prose by *does the core need it?* would conclude there was nothing to do, which is exactly how the previous three deferrals were argued.

## What changes

Three things, and no removal of anything.

**P-013 opens.** Status moves `draft` → `active` and the plan gains a milestone table, M-062 through M-066, continuing corpus-global numbering from P-012's M-061. The milestones are ordered so that measurement precedes removal: M-062 inventories and scores, M-063 makes the corpus able to shrink at all, M-064 executes or declines each candidate on the record, M-065 determines [AN-013](../../docs/analysis/AN-013-distributed-work-and-evidence-boundaries.md)'s three open findings, and M-066 closes the campaign on re-derived evidence as [PDR-026](../../docs/decisions/PDR-026-campaigns-close-on-re-derived-evidence.md) requires.

**PDR-029 records the second test.** Simplification is judged by two questions, not one: *does the core need it?* for rules, artifact types, required fields, commands, and checks; and *is it stated once, and can a reader check it?* for carrier prose. The second test is what makes the methodology carriers a measurable surface rather than a matter of taste, and it also states the trap: satisfying it sometimes makes a passage longer, so word count is not the campaign's measure of success.

**The record notes what is not safe to compact.** Carrier prose is instruction and can be rewritten as long as the rules survive; the corpus is the memory multiple humans and agents share, and its job is to still mean the same thing to a different reader later. Reducing what the corpus can remember is not simplification, and M-063 — the milestone that touches that surface — carries the heaviest guard in the campaign for that reason.

## Evidence this change already has

The skill set was measured while sizing the campaign, at this branch's base commit `ab946bc`. The six shipped skills total 11,422 words, of which 6,508 — 57% — is shared fragments rendered into more than one skill. [ADR-021](../../docs/decisions/ADR-021-generated-standalone-skills.md) chose that duplication deliberately so a copied skill folder stays complete, and this change does not reopen it; the file-level repetition is not the cost. The cost is that `clue-delta` step 5 directs the agent to `clue-verify`, both render the same five fragments, and `AGENTS.md` has already stated the tier rules and the ready contract a third time — so an ordinary full change reads about 6,329 words of instruction of which roughly 1,965 is already in context.

These figures are recorded here as the reason M-062 names the skill set as a surface. They are not M-062's evidence: this change measures one carrier family at one commit, and the milestone requires the full inventory.

## Acceptance

P-013 is `active` and carries M-062…M-066, each with a verifiable exit criterion. PDR-029 records the two tests, their split by surface, and the word-count trap. No artifact is removed, no rule changes, and no methodology carrier is edited by this change — the campaign that will edit them is merely opened.
