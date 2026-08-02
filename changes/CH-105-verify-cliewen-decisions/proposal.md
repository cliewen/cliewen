---
id: CH-105
type: change
status: open
links: [P-010, M-048, ADR-029, PDR-004]
title: Cliewen's own decisions become human-verified
---

# CH-105 — Cliewen's own decisions become human-verified

**Serves P-010's M-048.** Taken ahead of the campaign's first arc at the human's direction; the arcs express priority, not a prohibition on order, and nothing in M-045…M-047 depends on this milestone or is blocked by it.

## What is wrong

Twenty-seven decision records — seventeen ADRs and ten PDRs — carry `status: inferred`, `author: agent`, and `accepted-by: []`. They are binding, because merging makes an agent-authored decision binding. They are not verified, because [PDR-004](../../docs/decisions/PDR-004-merge-binds-approval-signs.md) and [ADR-029](../../docs/decisions/ADR-029-accepted-by-is-cliewen-approval-only.md) deliberately separate the two: a merge accepts a change, and only an explicit human approval accepts the *decision* inside it.

That separation is right, and it is why the backlog exists rather than being a bug. But nearly half of Cliewen's decision corpus was written by an agent and never read back by the human it binds — including the decision that draws the core's red line, the decision defining what `accepted-by` may record, and the decisions governing releases, migrations, and the merge boundary itself. A methodology whose core promise is that humans verify meaning cannot leave its own meaning unverified and still mean it.

## What changes

**Every one of the twenty-seven is promoted to `status: verified` with the approver, date, and venue recorded in `accepted-by:`.** The approval was given explicitly, in one conversation on 2026-08-02, against the complete list of the twenty-seven presented by identity and title. `author: agent` is unchanged: who wrote a decision is a fact about its origin, and promotion records who accepted it, not who authored it.

The approval was informed by two consequences stated before it was given: that ADR-029 governs this very operation and would be promoted by it, and that ADR-039 and PDR-022 constrain P-010's own remaining milestones. Neither turned out to require holding a record back — a command that reaches the network is not the judge going online, and a hook file is not the entry point file whose contents PDR-022 bounds — so M-045 and M-047 extend those decisions rather than superseding them.

No record's text changes. Promotion is a status and a signature; a decision whose content needed editing would be a different change with a different argument.

## What does not change

**The bar for future decisions is unmoved.** New agent-authored decisions are still born `inferred` with `accepted-by: []`, and merging still does not sign them. This change discharges a backlog; it does not make promotion automatic, and it creates no precedent that a batch may be signed without being read.

**No decision is retired here.** M-048 permits supersession as the alternative outcome, and none of the twenty-seven was found obsolete.

## Reversal cost

Cheap and local per record: promotion is two frontmatter lines, and an objection later reverts one file and opens a question. What is not cheap is a signature that was not meant, which is why the list was presented in full before the approval rather than after it.
