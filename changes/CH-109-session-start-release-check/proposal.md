---
id: CH-109
type: change
status: active
links: [P-010, CAP-001, CAP-004, PDR-022, ADR-042, ADR-018, C-013]
title: A coding agent learns a release is available when its session starts
---

# CH-109 — A coding agent learns a release is available when its session starts

Serves [P-010](../../docs/plans/P-010-adopters-keep-current.md) milestone **M-047**.

## What

Cliewen now has a command that answers "is there something newer?" ([ADR-042](../../docs/decisions/ADR-042-release-check-outside-the-judge.md), M-045) and a skill that carries out the upgrade once someone decides to ([ADR-043](../../docs/decisions/ADR-043-upgrade-skill-is-a-managed-carrier.md), M-046). Neither runs on its own. `clue latest` is a command an adopter has to already know exists, and `clue-upgrade` is a skill an agent invokes when asked to upgrade — which is the request nobody makes, because nothing told them there was anything to upgrade to. The quiet mode built in M-045 was built for a caller that does not yet exist.

This change supplies that caller: `clue init` materializes `.claude/settings.json` with a `SessionStart` hook that runs `clue latest --quiet`, so a session that opens in a behind repository is told so, in one line, once, and a session that opens in a current one is told nothing at all.

That crosses a line [PDR-022](../../docs/decisions/PDR-022-vendor-entry-points-only-point.md) drew but did not consider. PDR-022 permits a scaffolded vendor entry point and confines it to *pointing* — the file it governs is prose a session reads. This one is configuration a session *executes*, which is a different category with a different failure mode: prose that is wrong misroutes an agent, while configuration that is wrong can hang a session or greet every start with an error. A new decision record extends the boundary rather than stretching PDR-022's wording over a case it never examined.

## Why now

M-047 is the last milestone of P-010's first arc, and the arc is incomplete without it. The first two milestones built the machinery for staying current and left the ignition unwired: an adopter who never runs `clue latest` and never asks to upgrade is exactly as uninformed as they were before P-010 opened. Discovery at session start is the one moment the system reaches an adopter who is not looking for it.

## Scope

**In:**

- A decision record extending the vendor entry-point boundary to executable vendor configuration, binding it to the same limits — adopter-owned on arrival, never overwritten by `clue init` and never rewritten by `clue migrate` — and to the limits executability adds: it runs the quiet release check and nothing else, it is bounded by a timeout, and it cannot fail a session.
- `clue init` emits `.claude/settings.json` carrying that hook, for every start reason a session has, and skips an existing file byte-for-byte.
- A migration notice reporting a settings file that never runs the check, and an absent one; neither is repaired.
- Acceptance criteria with positive and negative Unit evidence, plus a `Human` criterion for the one claim no test can make: that the line arrives in a real session.
- This repository's own `.claude/settings.json`, so the change is dogfooded and the human observation has something to observe.
- Every live carrier that states what `clue init` emits or what `clue migrate` reports: CAP-001's README, criteria, and design; CAP-004's design where the quiet mode's caller is named; `guide/operations.md`; `guide/getting-started.md`; the migration inventory; `[Unreleased]`.

**Out:**

- Any hook that is not the release check. The decision fixes the emitted configuration at one command; a repository that wants formatting, linting, or notification hooks writes them itself, below or beside what Cliewen emitted, in a file Cliewen never rewrites.
- A second vendor's configuration. PDR-022's evidence bar is unchanged: an entry point is emitted for an assistant whose published behaviour makes the check unreachable without one.
- Any repair of an adopter-owned settings file. Reporting is the whole mechanism, exactly as it is for the entry point.
- Changing what `clue latest --quiet` prints. Its contract was settled in M-045 and this change only calls it.
