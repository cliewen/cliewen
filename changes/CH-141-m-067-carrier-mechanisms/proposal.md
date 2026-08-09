---
id: CH-141-proposal
type: proposal
status: active
links: [P-013]
title: M-067 — the answers M-062 produced reach their carriers and their machines
---

# CH-141 — M-067: carrier and machine mechanisms

P-013's M-062 registered every rule-bearing statement in the shipped skills and the routing hub, escalated the ones that traced to nothing or needed a human call, and the human answered them on 2026-08-08 (PDR-030…PDR-034, amendments to PDR-006 and PDR-017, log rows recorded under CH-130). Those answers are decisions on paper; nothing yet builds the mechanisms they asked for. M-067 is that build. It is P-013's last `todo` milestone before M-066 closes the campaign.

## What this change does

1. **F-DW-03/F-DW-05 (no-computed-figures rule)** — mint constraint `C-018` (`enforcement: agent`, promotion trigger stated) over the rule already stated in the shared `durable-work` fragment. Tighten `checkConstraints` so a constraint's `source:` naming a corpus ID must resolve to a live artifact. Repair the five constraints [AN-022](../../docs/analysis/AN-022-remaining-surface-scored.md) found pointing at "AGENTS.md rule N", a numbering CH-132's hub restructuring removed (C-001, C-002, C-004, C-005, C-013). Repair `clue-verify`'s unbounded "assessed against every constraint" checklist item so it points at a countable population (the constraints README index) rather than an unbounded claim.
2. **Reference-names-what-it-is-about rule** — mint constraint `C-019` (`enforcement: human`, priced residual) applying [ADR-046](../../docs/decisions/ADR-046-index-rows-say-what-the-artifact-is-about.md)'s rule to agent-facing reports (review findings, verification handoffs, acceptance briefs), and state it in `clue-verify`'s agentic review loop where a finding's shape is defined.
3. **F-RB-09 (orient-after-merge)** — write the missing decision (a log row: cheap and local to reverse, restating CH-018's original practice in the form PDR-029's register can trace to), mint constraint `C-020` (`enforcement: human`) over it, and move the rule's statement from the `review-boundary` fragment to `durable-work`, whose subject — what an agent does with work and its own next steps once a change is behind it — it actually shares. `durable-work` is added to `clue-extract`'s render so the move does not silently drop the rule from that skill's reading path.
4. **PDR-032's mid-change triage rule** — state it in the shared `durable-work` fragment (the carrier PDR-032 itself names as owed) and in `clue-delta`'s Propose/Digest steps.
5. **PDR-006's rejected-record clause** — state it in the shared `decision-records` fragment, which PDR-006's own carrier note names as not yet stating it.
6. **Hub `clue-extract` row + generator test** — already landed early, in CH-130 (`log.md` row, 2026-08-08). Verified still green; no new work.

## Out of scope

M-066 (campaign close on re-derived evidence) — this change's own digest updates M-067's plan row only.

## Links

- [P-013](../../docs/plans/P-013-simplification.md) — M-067
- [AN-018](../../docs/analysis/AN-018-skill-statement-register.md) — F-DW-03, F-DW-05, F-RB-09
- [AN-021](../../docs/analysis/AN-021-remaining-carrier-register.md) — GL-77 (F-RB-09 restated to a human)
- [AN-022](../../docs/analysis/AN-022-remaining-surface-scored.md) — the five stale constraint sources
- [PDR-032](../../docs/decisions/PDR-032-mid-change-suggestions-are-triaged.md), [PDR-006](../../docs/decisions/PDR-006-decision-records-are-typed.md), [ADR-017](../../docs/decisions/ADR-017-conventions-are-constraints.md), [ADR-046](../../docs/decisions/ADR-046-index-rows-say-what-the-artifact-is-about.md)
