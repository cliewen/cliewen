---
id: PDR-013
type: decision
status: inferred
links: [G-001, P-005, PDR-006, C-012, C-013]
title: Cliewen has an explicit core behind a red line; everything else is periphery adopters may extend
author: agent
accepted-by: []
---

# PDR-013 — Cliewen has an explicit core behind a red line

## Context and problem statement

Cliewen has grown a corpus of fourteen artifact types, five skills, three change tiers, and a public guide, and every simplification debate so far has been argued case by case without a criterion. The methodology names many protected surfaces but has never said which parts are load-bearing — the parts that, if changed casually, would silently change what a green `clue validate` or a merged PR *means*. Without that statement, simplification has no test ("does the core need it?"), protection has no boundary (everything feels equally sacred, so nothing is), and adopters cannot tell what they may freely extend from what they must not touch. What is Cliewen's core, and what rule protects it?

## Decision outcome

**Cliewen is an explicit core behind a red line. The core is three things; everything else is periphery.**

The core:

1. **The verifiable thread** — the chain goal → plan → change → capability → acceptance criterion → test, in which every durable claim traces to executable evidence.
2. **The human merge boundary** — an agent never merges its own change; the human merge is the act of acceptance ([C-012](../constraints/C-012-agents-never-merge-own-changes.md)).
3. **The deterministic judge** — `clue validate` as the machine check of corpus form, enforced as a wall by protected CI.

**The red line:** a change that alters the meaning of any core element — what the thread connects, what a merge accepts, what a green validate asserts — requires an explicit decision record and human acceptance; it is never plain, never light, and never rides silently inside another change. Periphery exists to serve the core and must never constrain it: when a peripheral rule and a core element conflict, the peripheral rule yields or is retired.

**Periphery** is everything else — analysis records, the public guide, foreign-soil trials, quality bars, tier prose, index generation, templates — changeable through the ordinary change loop at ordinary cost.

**Extension:** adopters extend Cliewen by adding their own artifacts and artifact types to their corpus under `/docs`; the core does not enumerate what a corpus may contain, only what the thread, the boundary, and the judge mean. The mechanical half of this (validator tolerance for adopter-defined types) is a corpus-format decision delivered separately.

**Carrier:** the register entry [C-013](../constraints/C-013-core-changes-need-decision.md) holds the red-line rule; [ARCH-003](../architecture/core.md) holds the durable core statement the rule protects.

### Rejected: keep protecting everything equally

The current protected-surface list treats a typo in guide prose and a change to merge semantics as the same class of event. Uniform protection is uniform friction: it neither warns harder at the actual load-bearing parts nor lets the periphery move cheaply, and it gives simplification debates no criterion.

### Rejected: define the core as a fixed artifact-type list

Enumerating blessed types would freeze the taxonomy — the opposite of the extension story — and would misplace the boundary: the core is not which folders exist but what the thread, the merge, and the validate *mean*. A type list also re-opens on every taxonomy change; a meaning statement does not.

### Rejected: core as guide prose only

Prose in the guide is periphery by this decision's own definition and carries no enforcement. A red line that lives only in explanatory text is a wish; the rule needs a register entry with an enforcement field and a named promotion trigger.
