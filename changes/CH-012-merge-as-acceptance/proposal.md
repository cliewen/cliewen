---
id: CH-012
type: change
status: open
links: []
title: Merge is acceptance — team-first ADR promotion
---

# CH-012 — Merge is acceptance

**Plan-less.** Methodology adjustment from the CH-009 retro: everything currently routes through one approver, and real teams do not actively approve ADRs even when a human wrote them. The acceptance event moves from the PR approval to the merge itself, so one rule serves solos and teams — and a solo project that becomes a team changes nothing.

## What

A new ADR supersedes [ADR-014](../../docs/decisions/ADR-014-pr-approval-promotes-adrs.md)'s promotion trigger:

- **The merge is the acceptance event.** Merging a PR accepts the `inferred` decisions it introduces; whoever merges accepts on behalf of the reviewing team. Silence through review plus a merge is consent.
- **Objection blocks, and is the disagreement mechanism.** An explicit review objection to a specific ADR keeps that ADR `inferred` and becomes an open question for a follow-up change. This is also the recorded answer to "what happens when reviewers disagree" — a defined door, nothing more built until a real multi-reviewer project needs it.
- **The agent still does the clerical flip** — `status: inferred → verified`, `accepted-by:` citing the merger, date, and PR — in the next digest (or a final pre-merge commit when the human has explicitly said "merge it").
- ADR-014 stays in the folder (decisions are never deleted) with a superseded-by note; its status remains `verified` — it was a real decision, now replaced.

## Files

- `docs/decisions/ADR-018-merge-is-acceptance.md` (new, born `inferred`)
- `docs/decisions/ADR-014-pr-approval-promotes-adrs.md` (superseded-by note)
- `docs/decisions/README.md` (promotion paragraph rewritten; index)
- `.agents/skills/clue-delta/skill.md` step 5 and `.agents/skills/clue-verify/skill.md` final item (promotion wording)
- CHANGELOG `[Unreleased]` entry (the skills ship to adopters — user-visible)

## Retroactive cleanup

Nothing to rewrite: past `accepted-by:` entries cite real approval events and stay historically accurate. Only forward mechanics change. The pending clerical flip of ADR-017 (accepted by the merge of PR #10) is performed in this change's digest — under the old rule and the new one alike.
