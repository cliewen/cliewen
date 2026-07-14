---
id: CH-013
type: change
status: open
links: [P-002]
title: Plan revision — P-002 gains the foreign-soil milestone (M-007)
---

# CH-013 — Foreign-soil milestone

**Plan revision to [P-002](../../docs/plans/P-002-leaves-home.md)** (semantic mutation → own change + ADR, per P-002's mutation rules). Fourth and last methodology adjustment from the CH-009 retro: Cliewen has only been proven on repos that share a maintainer and a mindset. The genuine test is foreign ground.

## What

P-002 gains milestone **M-007 — Foreign soil**: the Cliewen skills are trialed on at least two external open-source repositories (selected by the human; not yet chosen); each trial produces an `AN-xxx` findings doc in `docs/analysis/`, and at least one methodology adjustment traces back to trial findings.

Scope guard: trials are findings-producing experiments, not adoptions — no PRs against the foreign repos, no new extraction mappings (P-002 explicitly keeps non-OpenSpec mappings out; a trial that demands one surfaces that as a finding, not as new scope).

## Files

- `docs/plans/P-002-leaves-home.md` — M-007 row (status `todo`)
- `docs/decisions/ADR-019-*.md` — the decision that validation requires foreign trials, born `inferred`

No changelog entry: nothing user-visible changes in `clue` or the skills — this is campaign bookkeeping for this repo.
