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

**Fix:** extend `clue validate --intent` (CAP-009, which already reports the vision and use cases as a state, never a scorecard) to also print each goal with its status and the capability identities whose required `goal:` field names it — a goal no capability yet names prints that plainly rather than nothing. This reuses an existing one-hop, required field (every capability already declares exactly one `goal:`), so the traversal stays bounded, the report stays derived and never goes stale, and — matching [AC-164](../../docs/capabilities/CAP-009-product-intent/criteria.md)'s and [PDR-054](../../docs/decisions/PDR-054-use-cases-are-optional-and-no-requirement-artifact.md)'s existing rule for use cases — it prints no percentage, ratio, or count.

**Explicitly rejected, and recorded in a new PDR:**
- A new goal status (e.g. `served`/`completed`). A goal is a durable want, not a delivery item: it does not stop being wanted merely because it is well served today, and a status a human must remember to revert when a regression unserves it again would drift silently. `proposed → accepted` stays exactly as ADR-025 left it.
- A coverage percentage or ratio. PDR-054 already rejected this shape for use cases for the same reason: a percentage reads as a target, and the only way to move it is to write artifacts nobody needs.

**Bundled by explicit human direction:** promote G-018, G-019, and this change's own motivating goal, G-020, from `proposed` to `accepted` (ADR-002 requires promotion through a change and PR; the human decision was already made in conversation).

## Challenge

The riskiest assumption: that naming which capabilities serve a goal, without also surfacing their criteria's evidence state, is enough signal for "is this goal actually served" to be useful — versus needing to drill one hop further into `@draft`/`@retired`/evidenced criteria.

Credible alternative: report per-criterion evidence counts beneath each capability instead of (or as well as) capability identities.

Cheapest useful test: this repository's own corpus already has both shapes of goal — G-001 has seven long-active capabilities serving it, G-018/G-019 (until this change) had none. Running the extended report against this repository during implementation and reading whether "served by CAP-001, CAP-002, ... (active)" versus "served by no capability yet" already answers the question that prompted G-020 is the test; if a reader still cannot tell whether a goal is served from that line alone, drill one hop further before this change is marked ready.

An implementation that passed AC-202 and still failed the adopter it is for: printing the capability list without each capability's own `status` — two capabilities both named but one `draft` and one `active` would look identically "served" when only one of them actually is.
