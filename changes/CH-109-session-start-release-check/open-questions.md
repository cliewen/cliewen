---
id: CH-109-open-questions
type: open-questions
status: open
links: [CH-109, P-010, PDR-022, ADR-018]
title: Open questions — CH-109
---

# Open questions — CH-109

## Q1 — M-047's exit criterion requires emitting vendor configuration, which violates agent-agnosticism

**Raised:** 2026-08-03, during implementation, by the human.

M-047 states that "the scaffold emits that configuration for a new session and for a restarted context window alike". The only mechanism that satisfies that literally is a vendor's own settings file — for Claude Code, `.claude/settings.json` — because no cross-agent standard exists for running anything at session start. The implementation was written and then rejected: Cliewen must be agent-agnostic, and emitting executable configuration for one assistant crosses that line.

The distinction matters and is not the one [PDR-022](../../docs/decisions/PDR-022-vendor-entry-points-only-point.md) already settled. That decision permits a vendor entry point that only *points* at the hub — prose, inert, and the minimum needed to make a cross-agent hub reachable at all. Configuration that *executes* is a different thing: it takes a position on what an adopter's tool does, for one vendor, in a project whose whole premise is that the method belongs to no assistant.

The milestone cannot be delivered as written. The question is what replaces it.

**Recommendation: make the carrier a skill, and document the vendor wiring instead of shipping it.**

Two parts, and the first is the one that matters.

Cliewen already has an agent-agnostic carrier that reaches existing adopters: the generated skills, refreshed by `clue migrate`. `clue-delta` opens with **Branch**, and that is the honest moment to ask whether there is anything to upgrade to — better than session start, in fact, because a session start catches an agent that may be about to do something entirely unrelated, while the start of a change is when the answer changes what happens next. One line in the skill — run `clue latest --quiet` before branching, and route a non-empty answer to `clue-upgrade` — reaches every adopter on every assistant, through a carrier this project already owns and already upgrades. It is less deterministic than a hook, because an agent can skip a line of a skill in a way it cannot skip a hook; that is the price of agnosticism, and it is the same price every other Cliewen rule already pays.

The second part covers the literal wish without shipping it: `guide/operations.md` gains a short, vendor-neutral section saying that `clue latest --quiet` is safe to run at session start — it is silent when current, silent when offline, bounded, and exits 0 — and showing how an adopter wires it into whichever assistant they use, as their configuration in their file. Cliewen describes; the adopter decides. That also settles the migration question by dissolving it: with nothing emitted, there is nothing whose absence to report, so no MIG-006 and no notice about a vendor an adopter may not use.

Both parts change M-047's exit criterion, which P-010's mutation rules put behind a declared plan revision backed by a decision record.

**Answered:** 2026-08-03, by the human, who accepted the prohibition and asked whether it could be handled better than by documenting a recipe. It could. The recommendation above was superseded before implementation: the skills reach existing adopters mechanically, but the start of a change is not the start of a session, and the second part left the actual wish unimplemented. The routing hub is the session-start channel — the cross-agent file every assistant reads when a session begins — and the objection that it is adopter-owned is answered by the entry point's existing pattern of emitting into the template and reporting for everyone else. Recorded as [PDR-023](../../docs/decisions/PDR-023-the-hub-is-the-session-start-channel.md), which carries the M-047 revision.

**Status:** resolved. No open questions remain.
