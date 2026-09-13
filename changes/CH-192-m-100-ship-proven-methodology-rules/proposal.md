---
links: [P-023, M-100]
---

# CH-192 — M-100: ship the two rules P-023 proved

## What

P-023's M-094–M-099 trialled two repository-only rules before letting either reach an adopter: the challenge-a-consequential-commitment rule ([PDR-061](../../docs/decisions/PDR-061-challenge-consequential-commitments-proportionally.md), trialled by M-095/[AN-026](../../docs/analysis/AN-026-trialling-the-challenge-rule.md)) and the guidance-capture rule ([PDR-062](../../docs/decisions/PDR-062-guidance-capture-has-an-eligibility-bar-and-no-new-type-yet.md), trialled by M-097/M-098). Both trials succeeded: M-095 shows the challenge rule redirecting one proposal and correctly waving another through, and M-098 shows a fresh agent finding, using, and (where warranted) correcting captured guidance from the entry point the rule names. M-100 is the single adopter-binding milestone in the campaign: it moves this proven text onto the shipped carriers under `internal/skills/source/`, regenerated into `.agents/skills/` and `internal/scaffold/templates/skills/`, or records why a proven piece does not move.

Nothing from the trials was rejected outright — the one declined alternative (a new `runbook` corpus type, PDR-062's "Rejected" section) was already declined before shipping and needs no further action here — so both rules ship, generalized from this repository's own specifics (`AGENTS.md`, `CONTRIBUTING.md`) to adopter-neutral language.

## Why

Every methodology rule this campaign produces binds only this repository until a change explicitly promotes it, per [ADR-062](../../docs/decisions/ADR-062-repository-role-is-declared-machine-state.md)'s carrier-boundary check. Shipping before trialling would have cost adopters a rule this repository had not yet itself sustained; not shipping after a successful trial leaves the campaign's own stated purpose (G-014, G-015) unmet for everyone but this repository.

## Approach

- Add a new shared skill template, `challenge-commitments.md.tmpl`, carrying the generalized challenge rule; include it in `clue-plan.md.tmpl` (before adopting or revising a plan's promise) and `clue-delta.md.tmpl` (before writing a consequential change's proposal).
- Extend the existing `durable-work.md.tmpl` shared template with the generalized guidance-capture rule (eligibility bar, home selection, evidence scoping, discovery, retirement, workaround guard) — no new heading, since it fits the section's existing "durable state" topic and every skill that already includes it (`clue-delta`, `clue-extract`, `clue-upgrade`, `clue-verify`) gains it without a routing change.
- Regenerate via `go generate ./internal/skills`.
- Promote [PDR-061](../../docs/decisions/PDR-061-challenge-consequential-commitments-proportionally.md) and [PDR-062](../../docs/decisions/PDR-062-guidance-capture-has-an-eligibility-bar-and-no-new-type-yet.md) to `binds: adopter`, each naming its shipped carriers per ADR-062's convention.
- Replace `AGENTS.md`'s repository-only challenge paragraph with a pointer to the now-shipped rule, since restating it would duplicate the single carrier and risk drift; keep the campaign-specific proportionality trial reference out, since the rule is no longer trial-only.
- Add acceptance criteria for both rules to CAP-006 (the capability already carrying AC-184's plan-health guidance, the nearest existing home for cross-cutting change/plan-lifecycle rules).
- CHANGELOG `[Unreleased]` entry naming the adopter-facing surface: `clue-plan`, `clue-delta`, `clue-extract`, `clue-upgrade`, and `clue-verify`'s generated skill text.
- Digest: mark M-100 done in P-023; since it is the plan's last unfinished milestone, set the plan `status: completed` in the same digest (no successor plan is named or decided).
