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

**Outcome, from the test actually run.** `ARCH-003`'s three core elements are the verifiable thread, the human acceptance boundary, and the deterministic judge; it does not mention vision or use-case at all — both are periphery, exactly as `intent.md` states. The delivery-thread description in `methodology.md` (Goal → Plan → Change → Capability → Acceptance criterion → Evidence) and the delivery thread drawn in `intent.md`'s diagram (Goal → Plan → Milestone → Change → Merge, converging on the same criteria and evidence) agree. Nothing in either guide contradicts `ARCH-003`. The mismatch is that `methodology.md` calls its subject "one red thread from motivation to acceptance evidence" without acknowledging the intent thread that `intent.md` later introduces — an omission, not a competing claim.

**What this settles for M-101.** No ADR or PDR is needed; the fix is an editorial correction to `methodology.md`'s framing (its opening sentence and diagram should not claim singularity that `intent.md` already qualifies). This routing conclusion is recorded here so M-101, when picked up, can proceed straight to that edit rather than re-deriving whether it is a red-line change. The edit itself remains M-101's own evidence and is not made by this change.

**What it caught.** A default assumption that any guide-level disagreement is a core-meaning conflict, when the actual disagreement was scope, not substance.

**What it cost.** Three file reads already needed to understand the milestone; no prototype, no separate investigation task.

## What neither trial exercised

Both commitments were caught before any implementation existed to discard, so the rule has not yet been tested against a course already partly built — the more expensive case PDR-061 is meant to prevent. Both are also corpus-and-methodology-shaped rather than user-facing, so PDR-061's fourth question — what implementation could satisfy every criterion and still fail the person the work is for — had no concrete feature to bite on in either trial; it was asked and answered vacuously ("no feature, no such implementation exists here") rather than doing real work. A future trial against an in-flight capability change would test that part of the rule for the first time.

## Disposition

M-095 is closed by this record: one commitment redirected (this milestone's own delivery approach), one proceeded without extra investigation (M-101's routing), both real and unsettled at the time the rule was applied, neither invented for the purpose.
