---
id: ADR-018
type: decision
status: inferred
links: [ADR-014, ADR-010]
title: The merge is ADR acceptance — objection blocks, teams included
---

# ADR-018 — The merge is ADR acceptance

## Context and problem statement

[ADR-014](ADR-014-pr-approval-promotes-adrs.md) made PR *approval* the acceptance event for the ADRs a PR introduces. But approval is exactly the act teams do not reliably perform: in practice reviewers read, comment, and merge — an explicit per-decision "I approve this ADR" rarely happens even for human-authored decisions, so decisions rot `inferred` for reasons that have nothing to do with doubt. The methodology must also be team-first: a solo project can become a team without changing its rules, so the acceptance event cannot depend on a widget or a ritual only solos perform consistently.

## Decision outcome

**Merging a PR is accepting the `inferred` decisions it introduces. Supersedes ADR-014's trigger; the clerical mechanics survive unchanged.**

- **The merge is the acceptance event.** Whoever merges accepts on behalf of the reviewing team; review-then-merge with silence on a decision is consent to it. No separate approval act is required — the merge is the one act every workflow, solo or team, actually performs.
- **Objection blocks.** An explicit review objection to a specific decision keeps that decision `inferred` even through the merge, and becomes an open question for a follow-up change. This is also the multi-reviewer disagreement rule: an unresolved disagreement is an objection, the decision stays `inferred`, and nothing further is built until a real multi-reviewer project demands more — a defined door.
- **The agent performs the clerical flip**, exactly as under ADR-014: `status: inferred → verified` and `accepted-by:` citing the merger, date, and PR — in the next change's digest, or in a final pre-merge commit when the human has explicitly instructed the merge.
- **Nothing else self-promotes.** The flip always cites a specific merge event; this covers decision `status` and the `provenance` field ([ADR-010](ADR-010-provenance-field.md)) alike.

**Carrier:** the promotion paragraph in `docs/decisions/README.md` (ships as `clue init` template prose); the promotion wording in the `clue-delta` and `clue-verify` skills (agent). No machine carrier — whether a merge happened is git's record already.

### Rejected: keeping explicit approval as the trigger (ADR-014)

Approval and merge are the same judgment expressed twice, and only the merge reliably happens. A rule tied to the act people skip makes the `inferred` count mean "unbookkept" instead of "unverified" — the exact failure ADR-014 was written to fix, reproduced one act earlier.

### Rejected: decisions born `verified` because merge review is guaranteed anyway

Already rejected by ADR-014 and still wrong: birth provenance and acceptance are different facts ([ADR-010](ADR-010-provenance-field.md)), and an unreviewed draft must never carry the accepted marker, even between commit and merge.

### Rejected: quorum or multi-approver rules for teams

Designing consensus mechanics before any multi-reviewer project exists is speculation. The objection rule covers disagreement minimally; richer semantics wait for a real demand.
