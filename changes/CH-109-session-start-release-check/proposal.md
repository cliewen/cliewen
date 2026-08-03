---
id: CH-109
type: change
status: open
links: [P-010, CAP-001, CAP-004, PDR-022, ADR-042, ADR-043, C-013]
title: A coding agent learns a release is available when its session starts
---

# CH-109 — A coding agent learns a release is available when its session starts

Serves [P-010](../../docs/plans/P-010-adopters-keep-current.md) milestone **M-047**, and revises it.

## What

Cliewen has a command that answers "is there something newer?" ([ADR-042](../../docs/decisions/ADR-042-release-check-outside-the-judge.md), M-045) and a skill that carries out the upgrade once someone decides to ([ADR-043](../../docs/decisions/ADR-043-upgrade-skill-is-a-managed-carrier.md), M-046). Neither runs on its own. `clue latest` only helps someone who already knows the command exists, and `clue-upgrade` only runs when an agent is asked to upgrade — the request nobody makes, because nothing told them there was anything to upgrade to. The quiet mode built in M-045 was built for a caller that did not exist.

This change supplies that caller through the routing hub: the `AGENTS.md` that `clue init` materializes asks the agent to run `clue latest --quiet` once when it starts, and to route a non-empty answer to `clue-upgrade`. `clue migrate` reports a hub that never asks, and rewrites nothing.

## The revision this change carries

M-047's original exit criterion required the scaffold to emit a vendor's session-start configuration — for Claude Code, a `.claude/settings.json` declaring a hook. That was implemented and then rejected: Cliewen is agent-agnostic, and shipping executable configuration for one assistant takes a position on an adopter's tooling that a methodology may not take. The rejection is recorded in `open-questions.md`.

The replacement is not a workaround for the ban; it is the better mechanism, and the vendor framing was what obscured it. Every assistant already reads one file when a session starts, and it is `AGENTS.md` — the cross-agent standard, re-read when a context window is rebuilt, and already the file that says how to work here. Using it needs no schema, names no vendor, and reaches the same moment.

The one objection P-010 itself raised against the hub — that it is the adopter's file, so `clue init` never overwrites it and `clue migrate` never rewrites it, leaving new repositories the only ones a hub line reaches — is answered by the pattern already shipped for the entry point: emit it in the template, and report its absence for everyone else without touching their prose. That pattern is better suited here than where it was first used, because the hub is vendor-neutral, so nobody is being asked to accommodate a tool they never chose.

[PDR-023](../../docs/decisions/PDR-023-the-hub-is-the-session-start-channel.md) records the channel and the prohibition; the M-047 revision and P-010's new out-of-scope boundary are declared in the plan and backed by it.

## Scope

**In:**

- A decision record naming the hub as the session-start channel and putting executable vendor configuration out of bounds for every assistant, in adopted repositories and in this one.
- The M-047 revision in P-010, plus a standing out-of-scope boundary for the campaign.
- The instruction in the emitted `AGENTS.md` template and in this repository's own hub.
- A migration notice reporting a hub that never asks, and an absent hub; neither repaired.
- Acceptance criteria with positive and negative Unit evidence, plus a `Human` criterion for the claim no test can make: that a real session reads the hub and runs the check.
- Every live carrier that states what `clue init` emits or what `clue migrate` reports: CAP-001's README, criteria, and design; CAP-004's design where the quiet mode's caller is named; `guide/operations.md`; `guide/getting-started.md`; the migration inventory; `[Unreleased]` and its migration section.

**Out:**

- Any configuration file for any assistant, including in this repository. That is the decision, not an omission.
- A line in the change loop's branch step. The skills reach existing adopters mechanically, which is a real advantage, but the start of a change is a different moment from the start of a session — an agent answering a question or reviewing a pull request never reaches it.
- Any repair of an adopter-owned hub. Reporting is the whole mechanism, exactly as it is for the entry point.
- Changing what `clue latest --quiet` prints. Its contract was settled in M-045 and this change only calls it.
