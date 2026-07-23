---
id: C-013
type: constraint
status: active
links: [PDR-013, C-012]
title: Changes to a core carrier require an explicit decision record and human acceptance
source: PDR-013, AGENTS.md rule 8
enforcement: agent
---

# C-013 — Core changes need an explicit decision record

Cliewen's core is defined in [ARCH-003](../architecture/core.md): the verifiable thread, the human merge boundary, and `clue validate` as deterministic judge. A change that alters the meaning of a core element — what the thread connects, what a merge accepts, what a green validate asserts — is never plain and never light: it uses the full change loop, mints an explicit decision record ([C-011](C-011-decision-records-typed.md) routes the type), and takes effect only through human merge. Periphery never constrains the core: a peripheral rule that conflicts with a core element yields or is retired, and that retirement is ordinary change content. The rule reaches adopters as the red-line rule in the AGENTS.md that `clue init` ships.

The protected carrier set includes this rule's own carriers: this constraint, its defining decision [PDR-013](../decisions/PDR-013-explicit-core-red-line.md), and the shipped AGENTS.md red-line rule. Weakening or removing the red line is itself a core-meaning change and crosses the red line.

**Promotion trigger:** `clue validate` learns to detect a diff touching a core carrier (ARCH-003, the merge-boundary rule C-012, this constraint and PDR-013, the AGENTS.md red-line rule and its init template, or the validator's own check set) without an accompanying decision-record change in the same branch — then `enforcement: machine` for that detectable subset. Judging whether a change alters core *meaning* stays agent-held.
