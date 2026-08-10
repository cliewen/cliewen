---
id: CH-143-proposal
type: change
status: open
links: [P-014, M-069, PDR-038, AN-022, ADR-034]
title: A surviving superseded decision stays prose; the residue is declined a second time
---

# CH-143 — Supersession residue declined a second time

## What

This change closes [P-014](../../docs/plans/P-014-supersession-edge.md)'s only milestone, M-069. It writes an ADR that decides the door [PDR-038](../../docs/decisions/PDR-038-supersession-residue-declined.md) named: whether `supersedes:` widens to cover a decision record that is superseded but survives, or whether the decline stands a second time with a fuller argument. The answer is the second decline. The new ADR states, at minimum, what obligation a superseding change would gain that it does not carry today, how `clue validate` would have to distinguish a live superseded record from a stale one, and how the reverse question — what was downstream of a given decision — stays unanswered without the reverse walk [PDR-034](../../docs/decisions/PDR-034-the-corpus-is-read-narrowly-by-default.md) argues against. The same change settles the nine decision index rows [AN-022](../../docs/analysis/AN-022-remaining-surface-scored.md) counted into one consistent prose shape, so the corpus's own record of a live supersession is a stated, named convention rather than nine independently worded sentences.

## Why

[PDR-038](../../docs/decisions/PDR-038-supersession-residue-declined.md) declined the widening once, inside a simplification campaign that [PDR-026](../../docs/decisions/PDR-026-campaigns-close-on-re-derived-evidence.md) forbade from adding machinery to argue for itself, and routed the question to this campaign so the argument could be made on its own terms rather than assumed. Revisiting it here, on a mechanism-owning milestone, the three costs a widening would add — a new obligation on every superseding change, a validator rule with nothing observable to check a live edge against, and a reverse-walk question the field alone does not answer — still exceed what closing a nine-record, already-readable gap would buy. The nine rows already carry the fact in prose; what they lack is a named, single shape, which this change gives them without adding a field, a check, or a reverse index.
