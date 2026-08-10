---
id: P-014
type: plan
status: completed
links: [P-013, PDR-038, PDR-026, PDR-029, AN-008, AN-022, ADR-034, ADR-055]
title: Cliewen gives a surviving superseded record a machine-visible edge
---

# P-014 — Cliewen gives a surviving superseded record a machine-visible edge

This campaign is [P-013](P-013-simplification.md)'s successor, opened in that campaign's own closing digest. [PDR-038](../decisions/PDR-038-supersession-residue-declined.md) named the door rather than building through it: `supersedes:` carries a pointer forward only when the artifact it names is deleted, so a decision record that is superseded but survives — the shape nine of the eighty-four decision index rows [AN-022](../analysis/AN-022-remaining-surface-scored.md) scored at its pin already describe in prose — carries no machine-visible edge, and nothing answers *what depended on this decision* except reading. Widening the field to that case was rejected inside [P-013](P-013-simplification.md) because [PDR-026](../decisions/PDR-026-campaigns-close-on-re-derived-evidence.md) forbids a simplification campaign from adding machinery to argue for itself; this campaign exists to make that argument on its own terms, not to assume the answer is yes.

## Milestones

Milestone numbering continues corpus-global numbering from [P-013](P-013-simplification.md)'s M-068.

| ID | Milestone (exit criterion) | Status | Evidence |
|---|---|---|---|
| M-069 | **The supersession residue PDR-038 named is decided, not merely revisited.** A decision record — an ADR, since this is a corpus-architecture question — either widens `supersedes:` to cover a decision record that is superseded but survives, or declines to do so a second time with a fuller argument than the first. A widening answer states, at minimum: what obligation a superseding change gains that it does not carry today; how `clue validate` distinguishes a live superseded record from a stale one, so the check cannot be satisfied by a pointer nobody maintains; and whether and how the reverse question — what was downstream of a given decision — is answered without the reverse walk the tooling deliberately does not do today. A second decline restates why each of those three costs still exceeds the nine-record problem it would close, rather than repeating PDR-038's reasoning unexamined. Either answer updates the nine affected decision index rows to the settled shape and closes this campaign's only milestone in the same digest. | `done` | [ADR-055](../decisions/ADR-055-surviving-supersession-stays-prose.md) declines the widening a second time, stating the three required costs, and names the settled prose shape for the nine affected rows in `docs/decisions/README.md` — seven of which already used it, with ADR-019's and PDR-007's corrected to match (CH-143). |

## Mutation rules

Status and evidence fields in the milestone table may mutate in an implementing change's merge digest. Everything else changes only through a declared plan revision backed by a decision record routed by reversal cost. Plan adjustments are decisions.
