---
id: CH-103
type: change
status: open
links: []
title: An adopter gets the Claude Code entry point this repository gave itself
---

# CH-103 — An adopter gets the Claude Code entry point this repository gave itself

**Plan-less.** P-009 is closed and no successor is designated.

## What is wrong

`clue init` writes `AGENTS.md` as the routing hub and mirrors the skills into `.claude/skills/` (ADR-018). It does not write `CLAUDE.md`, and Claude Code reads `CLAUDE.md`, not `AGENTS.md` — the vendor states this outright in its own documentation and gives the `@AGENTS.md` import as the supported way to bridge the two.

The consequence is precise and one-sided. An adopter's Claude Code session lists the five `clue-*` skills, so it can be told to follow the method, but the routing hub never reaches it: it never learns to classify a change's tier before touching the corpus, that a Cliewen change branches and never merges itself, or that `/changes/` is transient. It has the manuals and not the instruction to open them.

This repository had exactly that gap until CH-101 closed it here, and the errors it produced are on the record — a plan-closure rule written into a file no adopter reads, and a release manifest never updated because the design that named it was read without opening it.

## What changes

**`clue init` emits a `CLAUDE.md` whose whole job is to point.** It imports `AGENTS.md` and says why it exists, so the pointer cannot be mistaken for a second place to write rules. It is adopter-owned from the moment it lands: init never overwrites it, and `clue migrate` never rewrites it, so Claude-specific instructions added below the import are safe.

**`clue migrate` reports an entry point that does not reach `AGENTS.md`** as MIG-005 — missing, or present but never importing the hub. It reports and never repairs, like MIG-004: a `CLAUDE.md` an adopter already wrote is their prose, and the remedy for the missing case is `clue init`, which is already non-destructive. Existing adopters do not re-run init, so without this the change would reach nobody who has already onboarded.

## What does not change

The methodology is still cross-agent and `AGENTS.md` is still the only hub. The emitted file carries no rules — a rule that lived in `CLAUDE.md` would be invisible to every other agent, which is the failure mode this repository already warns about in its own `CLAUDE.md`.

No vendor gains a claim on the corpus. `/docs`, the skills, and the routing hub are unchanged.

## Reversal cost

Moderate and public: an emitted file becomes part of what adopters have, and a migration ID is append-only once published. The vendor-neutrality question — whether a cross-agent methodology may ship a vendor's flagship file at all — is the kind of decision the project must be able to point at later, so it gets a PDR rather than a log row.
