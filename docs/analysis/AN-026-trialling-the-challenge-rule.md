---
id: AN-026
type: analysis
status: active
links: [G-014, P-023, PDR-061]
title: Trialling the challenge rule on real work, in both directions
---

# AN-026 — Trialling the challenge rule on real work, in both directions

## Purpose

[P-023](../plans/P-023-challenge-plans-and-retain-what-work-teaches.md)/M-095 requires [PDR-061](../decisions/PDR-061-challenge-consequential-commitments-proportionally.md)'s challenge to run against two real campaign decisions, not exercises, and to land in both directions: one redirected, one proceeding without extra investigation. This records both trials, run on 2026-09-13 while starting M-095 itself.

## Trial 1 — redirected: how M-095 should be satisfied

**The commitment.** M-095 is a milestone whose approach was not yet settled: how to produce the required evidence.

**The unexamined path.** Write two illustrative examples showing the rule saying yes and no.

**The challenge.**

- *Assumption most likely to undermine the work:* any two invented decisions would satisfy the exit criterion, since it only asks for one redirection and one non-redirection.
- *Credible alternative:* the milestone's own text excludes exercises — "Both are real campaign work" — so an invented example fails the criterion regardless of how well it illustrates the rule.
- *Cheapest useful test:* check whether real, currently unsettled commitments already exist in this campaign to run the rule against, before writing anything.
- *Result that would stop or revise:* if none existed, the milestone would have to wait for real work rather than manufacture a case.

**Outcome.** Two real candidates were already `todo` with unsettled approaches: M-096 (which choices to settle about guides, and how) and M-101 (whether reconciling `guide/methodology.md` and `guide/intent.md` changes accepted meaning). The challenge redirected M-095's delivery from an invented demonstration to running the rule on these two, and using whatever each trial actually produced as the evidence — which is what Trial 2 below is.

**What it caught.** A plausible-looking shortcut that would have produced a document about the rule instead of a document produced by the rule.

**What it cost.** One re-reading of the milestone's exit criterion and a look at `clue next --all`; no wasted work, since the redirection happened before any invented example was written.

## Trial 2 — proceeds without extra investigation: M-101's routing

**The commitment.** Whether M-101 — reconciling `guide/methodology.md`'s "one red thread from motivation to acceptance evidence" against `guide/intent.md`'s "two threads, meeting at the goal" — needs a decision record before the correction can be made, because a change to what a core element means is always full with an explicit decision under [ARCH-003](../architecture/core.md)'s red line.

**The unexamined path.** Treat the two guides' differing thread counts as a meaning conflict and open a PDR or ADR before touching either file.

**The challenge.**

- *Assumption most likely to undermine the work:* that the two guides assert incompatible claims about what Cliewen's core actually is, which would make the correction a red-line change.
- *Credible alternative:* the guides describe different scopes — `intent.md` separates an optional intent thread (vision → goal → use case → capability) from the delivery thread, while `methodology.md` only ever describes the delivery thread and predates the intent/vision split. The "conflict" is a framing gap in one guide's opening sentence, not a disagreement about what either thread does.
- *Cheapest useful test:* read `ARCH-003`, `guide/methodology.md`, and `guide/intent.md` together and compare what each actually asserts.
- *Result that would stop or revise:* if `ARCH-003` named vision or use-case among its three core elements, or if the guides described the delivery thread's steps differently from each other, the correction would touch accepted meaning and need a decision record.

**Outcome, from the test actually run — corrected after review.** The first pass of this test read only the guides' prose and missed a real edge-level discrepancy: `intent.md`'s diagram does not draw a Change reaching Capability, Acceptance criterion, or Evidence at all — Change's only edge there goes to `Merge`, and Capability is fed exclusively from Goal and (optional) Use case, with the diagram's own prose stating the two paths "touch at the goal and nowhere else." `methodology.md`'s diagram, and `ARCH-003`'s own prose ("Goal → plan → change → capability → acceptance criterion → acceptance evidence"), both draw Change flowing directly into Capability. Read literally, that is not merely a missing acknowledgment — it is two different claims about whether Change reaches Capability at all, and a change to what a core element connects is exactly what `ARCH-003`'s red line asks a decision record for.

Resolving that requires one more check the first pass skipped: whether the corpus's actual durable structure supports either drawing. Every capability in this repository carries a `goal:` frontmatter field naming the Goal it serves, and none carries any field naming a Change — because a Change is transient, its workspace is deleted at digest, and nothing durable could hold a link to an identity that no longer exists once the work lands. So `intent.md`'s diagram is the one that matches what the corpus actually, durably links: Capability is reachable from Goal (and optionally Use case), never from Change. `methodology.md`'s diagram and its "one red thread" framing draw a persistent Change → Capability edge that no artifact in this repository ever encodes; it describes Change's causal role in *changing* a capability, not a link that survives the change. `ARCH-003`'s prose sentence is the same shorthand, but its own text never claims the thread is a set of literal, persistent field links — only that a durable claim traces to its declared proof, which it does, through Goal.

**What this settles for M-101.** No ADR or PDR is needed, but not for the reason the first pass gave. `ARCH-003`'s core meaning is unchanged and uncontradicted by either guide once "the thread" is read as the durable claim-to-proof relationship it defines, not as a literal edge diagram. The correction routes to `methodology.md`, not `intent.md`: its diagram and opening sentence should stop drawing Change as flowing into Capability, since no durable link does, and should instead show Change acting on Capability without a persistent edge — matching what `intent.md` already draws correctly. This routing conclusion, and the specific edit target, are recorded here so M-101 can proceed straight to it rather than re-deriving which guide is wrong.

**What it caught, on the second pass.** The first pass accepted a comfortable answer (this is a wording gap, not a meaning conflict) without checking the diagrams' edges or the corpus's actual field-level structure — exactly the fluent-and-confident failure mode PDR-061 exists to interrupt. An independent review of this record caught the first pass's error before it reached M-101; the second pass, checked against every capability's actual frontmatter, is the one recorded here.

**What it cost.** Three file reads already needed to understand the milestone, plus a second pass after review — checking every capability's frontmatter for a `goal:`/Change field — that the first pass should have run before writing a conclusion. Still no prototype and no separate investigation milestone, but the true cost of "proceeds without extra investigation" here included one round of getting it wrong.

## What neither trial exercised

Both commitments were caught before any implementation existed to discard, so the rule has not yet been tested against a course already partly built — the more expensive case PDR-061 is meant to prevent. Both are also corpus-and-methodology-shaped rather than user-facing, so PDR-061's fourth question — what implementation could satisfy every criterion and still fail the person the work is for — had no concrete feature to bite on in either trial; it was asked and answered vacuously ("no feature, no such implementation exists here") rather than doing real work. A future trial against an in-flight capability change would test that part of the rule for the first time.

## Disposition

M-095 is closed by this record: one commitment redirected (this milestone's own delivery approach), one proceeded without extra investigation (M-101's routing), both real and unsettled at the time the rule was applied, neither invented for the purpose.
