---
id: CH-106-open-questions
type: change
status: open
links: [CH-106, ADR-017, ADR-035, C-004, C-013]
title: Open questions for CH-106
---

# CH-106 — Open questions

One question remains open. Two earlier forks were closed by narrowing the scope and are recorded below as resolved, so the reasoning is not lost with the conversation that produced it.

## Q1 — Does the judge *count* a filler row, or *fail* on it? (blocking)

**Why this is still the human's call even after narrowing.** It changes what `clue validate` fails on. The judge is a core carrier under C-013, and two repositories — Tank Royale and model2diagram — are live on Cliewen with index blocks written before any of this was stated. Neither's current row inventory is knowable from here, so a failing check may turn a green adopter red on upgrade for filler the generator itself wrote into their README.

**Recommendation: count, do not fail.**

This repository has already run the experiment. The inferred-decision counter sat on the OK line as a visible, non-blocking number for months; it did not rot, it drove a campaign milestone, and CH-105 drove it to zero. [ADR-035](../../docs/decisions/ADR-035-bounded-provenance-and-reality-edges.md) settled the general form — costly unverified meaning is reported as an actionable population rather than turned into a build failure — and [ADR-017](../../docs/decisions/ADR-017-conventions-are-constraints.md) already applies it to the constraint backlog in the words "the backlog is visible, not archival". A counted population is the established Cliewen answer for something true but not yet repaired everywhere.

Failing would also punish adopters for the tool's own output. Every bare row in this corpus was written by `regenIndex`, not by a person; making the judge reject text the generator emitted last release, in a file the adopter owns, inverts where the fault lies. The generator fix stops the source, and the count shows the remainder.

**The cost, stated plainly.** Nothing clears this count automatically: `regenIndex` only appends missing entries and deliberately preserves rows that still cover a live target, so every repair is by hand. A count no command can drive down is weak, which is why the rows must be listable behind a flag (the `--reality-gaps` precedent) rather than existing only as a number. And under [C-004](../../docs/constraints/C-004-never-weaken-checks.md) the count can never be softened to make it look better. That is the same trade ADR-017 and ADR-035 each already accepted.

## Resolved — should the judge check that a label matches its record's title?

**No.** Recognizing a filename stem is mechanical and carries no opinion; requiring a label to equal the frontmatter title character for character is house style, and several of this corpus's titles exceed a hundred characters. That is a reasonable choice for this repository and not one to impose on every adopter's README through the judge. The consequence is accepted and stated in the proposal: this repository's own 58-row invariant is held by the generator and by review rather than by a machine.

## Resolved — should `regenIndex` normalize an existing divergent row?

**No, not in this change.** `internal/scaffold/scaffold.go` documents the opposite as a deliberate contract: regeneration "keeps existing entries whose targets still exist (hand-written descriptions survive)". Reversing it needs its own decision against [ADR-019](../../docs/decisions/ADR-019-init-regenerates-indexes.md), because it converts `clue scaffold` from a tool that appends to one that edits prose inside adopter-owned files. It is also ambiguous in practice: the tail boundary is defined as "after the status", so a row carrying no status gives the tool nothing to split on and it must either guess or refuse. Worth revisiting only with evidence that the counted backlog does not fall on its own.
