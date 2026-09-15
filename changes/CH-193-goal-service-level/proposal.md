---
id: CH-193
type: change
status: open
links: []
title: The intent report states which capabilities serve each goal
---

# CH-193 — The intent report states which capabilities serve each goal

## Proposal

This change is plan-less: no active plan carries an unfinished milestone for it, and it is small enough to serve [G-020](../../docs/goals/G-020-goal-service-level-is-observable.md) directly.

Serve G-020: an adopter running brownfield extraction asked an agent about the status of its goals and got only `proposed`/`accepted` back, because that is all [ADR-025](../../docs/decisions/ADR-025-one-status-lifecycle.md) gives a goal — modeled on the vision's durable-want shape, not a plan's delivery-shaped `→ completed`. There was, additionally, no mechanized way to see which capabilities currently serve a goal at all: [CAP-007](../../docs/capabilities/CAP-007-focused-context/README.md)'s `clue context` deliberately never follows reverse links from a goal, to avoid pulling most of a mature corpus into one result.

**Fix:** extend `clue validate --intent` (CAP-009, which already reports the vision and use cases as a state, never a scorecard) to also print each goal with its status, the capabilities whose required `goal:` field names it, and the plans whose `links` name it — each capability and plan shown with its own status, and "none" printed plainly rather than nothing when there is neither. Capabilities are the intent thread's ongoing capacity; plans are the delivery thread's one-time achievement, and Cliewen already has a real `completed` status there (ADR-025). Both are one-hop, bounded reads over an existing field or convention (every capability's required `goal:`; every plan's own `links` naming the goal it serves, already the corpus's practice), so the report stays derived and never goes stale, and — matching [AC-164](../../docs/capabilities/CAP-009-product-intent/criteria.md)'s and [PDR-054](../../docs/decisions/PDR-054-use-cases-are-optional-and-no-requirement-artifact.md)'s existing rule for use cases — it prints no percentage, ratio, or count.

**Explicitly rejected, and recorded in a new PDR:**
- A new goal status (e.g. `served`/`completed`). A goal is a durable want, not a delivery item: it does not stop being wanted merely because it is well served today, and a status a human must remember to revert when a regression unserves it again would drift silently. `proposed → accepted` stays exactly as ADR-025 left it.
- A coverage percentage or ratio. PDR-054 already rejected this shape for use cases for the same reason: a percentage reads as a target, and the only way to move it is to write artifacts nobody needs.

**Bundled by explicit human direction:** promote G-018, G-019, and this change's own motivating goal, G-020, from `proposed` to `accepted` (ADR-002 requires promotion through a change and PR; the human decision was already made in conversation).

## Challenge

The riskiest assumption: that naming which *capabilities* serve a goal is enough signal for "is this goal actually served" to be useful.

Credible alternative: also read the delivery thread's plans, not just the intent thread's capabilities.

Cheapest useful test, run during implementation: this repository's own corpus already has the exact motivating case. G-003 ("Cliewen is public") is unambiguously served — yet no capability's `goal:` field names it. The first implementation, reading capabilities alone, printed `G-003 ... capabilities: none` — a false negative for the adopter's own original example. Reading this repository's plans found `P-003 — Cliewen goes public` (`completed`, `links: [G-003]`): the delivery thread had already recorded the achievement Cliewen already tracks a real `completed` status for. The fix now reports both threads, and G-003 correctly shows `plans: P-003 (completed)`. This stands as the recorded result of the test named above, not a hypothetical: the alternative course was needed, and the change was revised before implementation continued.

An implementation that passed AC-202 and still failed the adopter it is for: printing capabilities alone (the first attempt, corrected above), or printing either list without each entry's own `status` — two capabilities or plans both named but one `draft`/`active` and one `completed` would look identically "served" when only one of them actually is.
