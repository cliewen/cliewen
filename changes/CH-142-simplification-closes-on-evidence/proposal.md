---
id: CH-142
type: change
status: open
links: [P-013, M-066, PDR-026, AN-006, AN-010, AN-012, AN-008, PDR-038]
title: P-013 closes on re-derived cost evidence
---

# CH-142 — P-013 closes on re-derived cost evidence

## Proposal

Serve P-013's M-066, the campaign's closing milestone. [PDR-026](../../docs/decisions/PDR-026-campaigns-close-on-re-derived-evidence.md) requires a campaign to close on re-derived evidence rather than a milestone table read; this change re-derives the three cost surfaces the campaign named — [AN-006](../../docs/analysis/AN-006-plain-change-overhead.md) (plain-change overhead), [AN-010](../../docs/analysis/AN-010-adopter-change-overhead.md) (adopter-change overhead), and [AN-012](../../docs/analysis/AN-012-adopter-configuration-cost.md) (adopter-configuration cost) — against the corpus and the tool at head, and states in each what changed since the analysis was last written, what did not, and what the campaign declined.

The digest closes P-013 (`status: completed`) once every milestone is evidenced, and designates a successor plan when one is decided. [AN-008](../../docs/analysis/AN-008-methodology-critiques.md)'s pattern C names a door P-013 left explicitly for its successor — widening `supersedes:` so a superseded-but-surviving decision record carries a machine-visible edge — and [PDR-038](../../docs/decisions/PDR-038-supersession-residue-declined.md) is the record that declined building it inside this campaign. Whether a successor plan exists, and whether it carries that door, is a human decision this change surfaces as an open question rather than assumes.

## Scope boundary

This change does not implement the widened `supersedes:` mechanism itself, does not reopen M-062 through M-065 or M-067/M-068 (already `done`), and does not re-run AN-013's determinations (M-065 closed those). It only re-derives AN-006, AN-010, and AN-012, and closes the campaign in its digest.
