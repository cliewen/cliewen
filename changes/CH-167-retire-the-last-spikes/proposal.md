---
id: CH-167
type: change
status: open
links: [P-020, M-086, PDR-052, PDR-053, ADR-034, ADR-063, G-011]
title: The last three spikes leave, and P-020 closes
---

# CH-167 — The last three spikes leave, and P-020 closes

## What

Retire `AN-001`, `AN-015`, and `AN-016` — the only analyses no standing rule cites — repair every reference a live artifact makes to them, close `P-020`'s `M-086`, complete the plan, and close `G-011`.

## Why

`PDR-053` settled which spikes may leave: those no live decision or constraint cites. Three qualify. Every other analysis in the corpus is retained by rule and needs no per-file argument, which is what turned `M-086` from a folder-wide review into a bounded job.

**`AN-001` is the Foundation Document, and the human directed its retirement.** It is the one file in the folder that was never a spike: it retired no risk and has no findings that landed anywhere, so the expiry mechanism does not reach it — it serves no plan, and the derivation deliberately leaves such an analysis for a human to retire on its own terms. Its substance is carried. The founding argument is in `guide/what-is-cliewen.md` and `guide/design.md`; the structural claims are `ARCH-001` (actors, lifetime classes, the frontmatter graph) and `ARCH-003` (the core red line); and the eleven skill statements that once rested on it were re-traced to `PDR-030` and `PDR-006`, which `AN-018` records as written for that purpose. Its own banner already concedes the point — where the document and the corpus disagree, the corpus wins — and a frozen v0.4 contract inside a corpus at v0.21 is the accumulation this campaign exists to end.

**`AN-015` and `AN-016` are ordinary spent spikes.** Each served a completed campaign, each proved a contract that now lives in `CAP-003` and the decisions it links, and neither is cited by any standing rule.

## No decision record is owed

`ADR-034` (retirement is deletion), `PDR-052` (a spent analysis is reported, not retained), and `PDR-053` (a cited spike stays) already decide everything this change applies. Retiring the Foundation Document is a judgement about this repository's own corpus, not a rule any adopter inherits — an adopter has no foundation document — so it is recorded in `M-086`'s evidence rather than as a record that would bind nobody.

## References to repair

Live artifacts are repaired; frozen ones are not, which is exactly what `ADR-063` bought.

- `AN-018` names `AN-001` twice in prose as the untraceable carrier its statements rested on. It is active, so both mentions are repointed to where those statements went.
- `AN-017` links `AN-016` in frontmatter and names it twice in prose as recording the composed failure paths. It is active, so the link is dropped as redundant — it already links `CAP-003` — and the prose names where that evidence now lives.
- `AN-023`'s row naming `AN-015` is left intact: a pinned inspection recording what it covered at a revision, breaking nothing, and rewriting it would make a pinned record a false account of itself. This follows `CH-165`.
- `P-009`, `P-011`, and `P-012` are completed and are not edited.
