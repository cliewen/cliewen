---
id: CH-109
type: change
status: open
links: [P-010, CAP-001, CAP-004, PDR-022, ADR-042, ADR-043, C-013]
title: A coding agent learns a release is available without being asked
---

# CH-109 — A coding agent learns a release is available without being asked

Serves [P-010](../../docs/plans/P-010-adopters-keep-current.md) milestone **M-047**, and revises it.

## What

Cliewen has a command that answers "is there something newer?" ([ADR-042](../../docs/decisions/ADR-042-release-check-outside-the-judge.md), M-045) and a skill that carries out the upgrade once someone decides to ([ADR-043](../../docs/decisions/ADR-043-upgrade-skill-is-a-managed-carrier.md), M-046). Neither runs on its own. `clue latest` only helps someone who already knows the command exists, and `clue-upgrade` only runs when an agent is asked to upgrade — the request nobody makes, because nothing told them there was anything to upgrade to. The quiet mode built in M-045 was built for a caller that did not exist.

This change supplies that caller twice over. The mechanism is in the tool: the ordinary `clue` workflow commands — `context`, `migrate`, `refs`, `init`, `scaffold` — print one advisory line to standard error when a newer release exists, so an agent learns it is behind through work it was already doing. The fallback is in the routing hub: the `AGENTS.md` that `clue init` materializes asks the agent to run `clue latest --quiet` before its first tool call, which covers a session that runs no `clue` command at all, and routes a non-empty answer to `clue-upgrade`. `clue migrate` reports a hub that never asks, and rewrites nothing.

## The revision this change carries

M-047's original exit criterion required the scaffold to emit a vendor's session-start configuration — for Claude Code, a `.claude/settings.json` declaring a hook. That was implemented and then rejected: Cliewen is agent-agnostic, and shipping executable configuration for one assistant takes a position on an adopter's tooling that a methodology may not take. The rejection is recorded in `open-questions.md`.

The replacement is not a workaround for the ban; it is the better mechanism, and the vendor framing was what obscured it. Every assistant already reads one file when a session starts, and it is `AGENTS.md` — the cross-agent standard, re-read when a context window is rebuilt, and already the file that says how to work here. Using it needs no schema, names no vendor, and reaches the same moment.

The one objection P-010 itself raised against the hub — that it is the adopter's file, so `clue init` never overwrites it and `clue migrate` never rewrites it, leaving new repositories the only ones a hub line reaches — is answered by the pattern already shipped for the entry point: emit it in the template, and report its absence for everyone else without touching their prose. That pattern is better suited here than where it was first used, because the hub is vendor-neutral, so nobody is being asked to accommodate a tool they never chose.

## The second revision, and why it is here

The hub-only answer above was implemented and then tested. **AC-085 was run against a real session on 2026-08-03 and failed.** The hub was loaded into the session's context and the release-check paragraph was its first content; the agent read past it into the requested work. A second finding came with it: the installed binary predated `clue latest`, so the instruction, had it been obeyed, would have produced an unknown-command error and a usage dump rather than one line or silence.

So the mechanism moved into the tool, and the prose kept only what nothing else can carry. An instruction asking for an unprompted action before work begins is weaker than the rules beside it, which constrain work already underway and are consulted at the moment they bind; no rewording changes that. A notice the tool volunteers does not depend on remembering. Both channels ship, because a session that runs no `clue` command is still a session.

This is the second revision M-047 has taken in one change, and the sequence is recorded rather than tidied away: agnosticism ruled out the reliable mechanism, the obvious replacement was assumed to be good enough, and a Human criterion actually being run is what found out that it was not.

[PDR-023](../../docs/decisions/PDR-023-tool-notice-and-hub-instruction.md) records both channels and the prohibition; it was revised in place rather than superseded, being unmerged, `inferred`, and unaccepted. The M-047 revisions and P-010's new out-of-scope boundaries are declared in the plan and backed by it.

## Scope

**In:**

- A decision record naming both channels and putting executable vendor configuration out of bounds for every assistant, in adopted repositories and in this one.
- The two M-047 revisions in P-010, plus standing out-of-scope boundaries for the campaign.
- The notifier: one advisory line from the workflow commands, sharing the quiet check's exact wording, on a shorter ambient budget, gated by command, `CI`, and `CLUE_NO_UPDATE_NOTIFIER`, with the judge and the version command excluded. No terminal condition — a coding agent captures output through a pipe, so requiring one would exclude the audience (Q3).
- The instruction in the emitted `AGENTS.md` template and in this repository's own hub, with a sharpened trigger and an unknown-command error read as the answer it is.
- A migration notice reporting a hub that never asks, and an absent hub; neither repaired.
- A bound on what a check that cannot answer costs: the release list's failure to answer is remembered for an hour under its own field, honoured only by the unrequested notice. Without it the day-long cache's "repeating this costs no request" was true only when the list answered, and the notice's own budget was the thing being spent repeatedly. Recorded as a decision-log row, carried by its own criterion, and leaving ADR-042's unreadable-cache rule untouched — that rule is about the cache mechanism failing, not the release list.
- Acceptance criteria with positive and negative Unit evidence, plus a restated `Human` criterion for the claim no test can make: that a real session learns the repository is behind, by either channel.
- Every live carrier that states what `clue init` emits, what `clue migrate` reports, or what the release check does: CAP-001's README, criteria, and design; CAP-004's criteria and design; `guide/operations.md`; `guide/getting-started.md`; the migration inventory; `[Unreleased]` and its migration section.

**Out:**

- Any configuration file for any assistant, including in this repository. That is the decision, not an omission.
- A line in the change loop's branch step. The skills reach existing adopters mechanically, which is a real advantage, but the start of a change is a different moment from the start of a session — an agent answering a question or reviewing a pull request never reaches it.
- Any repair of an adopter-owned hub. Reporting is the whole mechanism, exactly as it is for the entry point.
- Changing what `clue latest --quiet` prints. Its contract was settled in M-045; this change adds a second caller sharing its exact wording.
- Any notice from `clue validate`. The judge's verdict, exit code, and output stay a statement about the repository alone — that exclusion is the reason this stays clear of the core's red line.
