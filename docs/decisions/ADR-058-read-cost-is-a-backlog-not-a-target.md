---
id: ADR-058
type: decision
status: verified
links: [P-015, ADR-057, ADR-056, CAP-007, C-022, C-004, AN-023]
title: A read-cost report is a backlog judged artifact by artifact, and a link is never deleted to move its count
author: agent
accepted-by: Flemming N. Larsen (2026-08-11, conversation)
---

# ADR-058 — Read cost is a backlog, not a target

## Context and problem statement

ADR-057 deliberately reports structural read cost without a registry or threshold, so a campaign needs a truthful way to inspect the population without deleting valid relationships merely to make the number smaller.

## Decision outcome

**The reported over-budget population is worked artifact by artifact, and a `links` entry is never removed just to reduce the count.** Inspection either removes relationships genuinely redundant for that artifact's reader or records an accepted exception and its reason in the artifact or plan. Accepted rows remain in the derived report, so zero is not an exit criterion; the requirement is that every row has an inspected outcome.

No mechanical rule selects surviving links, and an accepted row is intentionally not distinguishable from an unexamined one by the validator. If the reader's actual output remains too expensive, CAP-007 owns the presentation boundary rather than the corpus relationship.

**Carrier:** P-015 M-072, ADR-057's backlog note, and C-022's residual judgment.
