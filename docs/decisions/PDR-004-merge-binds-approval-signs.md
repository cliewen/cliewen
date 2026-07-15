---
id: PDR-004
type: decision
status: verified
links: [PDR-001, ADR-010]
title: Merge makes a decision binding; approval verifies it — approvers sign, first signature dates it
author: agent
accepted-by: Flemming N. Larsen (2026-07-14, PR #11 review)
---

# PDR-004 — Merge binds, approval signs

## Context and problem statement

[PDR-001](PDR-001-pr-approval-promotes-adrs.md) made PR approval the acceptance event for the decision records a PR introduces, and the merge was read as implying that approval. But merging and approving are different judgments: a merge says "ship it", a signature says "I stand behind this decision" — and in real teams the explicit per-decision approval is exactly the act that rarely happens, even for human-authored records. Conflating the two either blocks merges on a ritual nobody performs, or silently inflates `verified` until it means "was in a merged PR". The methodology must also be team-first: a solo project can become a team without changing its rules.

## Decision outcome

**Merging a PR makes its `inferred` decisions binding. Only an explicit human approval makes them `verified` — each approver signs, and the first signature dates the acceptance.**

- **Merge = binding.** Once merged, a decision is in force: the code follows it, future changes respect it, nothing waits for a signature. Merge does not touch `status:` — nobody approved anything by merging.
- **Approval = verification.** `verified` means at least one human has explicitly endorsed the decision — a PR review approval, a review comment, or a stated "approved" in conversation all count; what matters is the recorded human judgment. Each approver signs `accepted-by:` (approvals accumulate; later approvers append), and the acceptance date is the **first** approval. The agent performs the clerical signing, citing approver, date, and venue (PR or conversation).
- **The `inferred` count is the honest gauge**: decisions in force that no human has explicitly endorsed — a visible backlog, exactly like the register's count of constraints awaiting machine checks. It shrinks by real approvals, never by bookkeeping.
- **Objection blocks verification, not the merge.** An explicit objection to a specific decision keeps it `inferred` regardless of other approvals, and becomes an open question for a follow-up change. This is also the multi-reviewer disagreement rule — a defined door; richer consensus mechanics wait for a real multi-reviewer demand.
- **Nothing self-promotes.** Every signature cites a specific human approval event ([ADR-010](ADR-010-provenance-field.md)'s two-tier provenance, unchanged).

**Carrier:** the promotion paragraph in `docs/decisions/README.md` (ships as `clue init` template prose); the promotion wording in the `clue-delta` and `clue-verify` skills (agent). No machine carrier — merges and signatures are git's and the frontmatter's record already.

### Rejected: merge auto-flips decisions to verified

Destroys what `verified` exists to mean: a human actually read and endorsed this decision. Under auto-flip the corpus lies about how much human judgment it contains — the exact silent decay the two-tier status ([ADR-010](ADR-010-provenance-field.md)) was built to make visible.

### Rejected: requiring explicit approval before merge (PDR-001's trigger, enforced)

Blocks shipping on the act teams demonstrably skip; decisions would either rot unmerged or the ritual would degrade into rubber-stamping. Binding-on-merge keeps work flowing while the unapproved backlog stays visible instead of fictional.

### Rejected: quorum or multi-approver thresholds for teams

Designing consensus mechanics before any multi-reviewer project exists is speculation. Accumulating signatures plus the objection rule cover disagreement minimally; more waits for real demand.
