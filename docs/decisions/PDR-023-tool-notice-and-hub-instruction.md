---
id: PDR-023
type: decision
status: verified
links: [P-010, CAP-001, CAP-004, ADR-011, ADR-018, ADR-042, ADR-043, PDR-019, PDR-022, C-013]
title: The tool carries the notice and the hub carries the instruction, and no vendor configuration is ever emitted
author: agent
accepted-by: Flemming N. Larsen (2026-08-08, conversation)
---

# PDR-023 — The tool carries the notice and the hub carries the instruction, and no vendor configuration is ever emitted

## Context and problem statement

`clue latest` and `clue-upgrade` cannot help an agent that does not know an upgrade exists, while executable session hooks are vendor-specific configuration that a cross-agent methodology cannot safely emit.

## Decision outcome

**The tool carries the notice and the hub carries the instruction.** Ordinary `clue` workflow commands — `init`, `scaffold`, `context`, `migrate`, and `refs` — report a known newer release on standard error; `validate`, `version`, and `latest` remain excluded. The notice runs only outside CI and when `CLUE_NO_UPDATE_NOTIFIER` is unset, honors the cached release answer, uses a short ambient budget, and degrades to silence at exit 0. It writes no repository files, changes no verdict, and leaves standard output and exit codes unchanged.

The routing hub's instruction remains the fallback for sessions that run no `clue` command, and an unknown-command error from an old binary is itself the signal that the binary is behind. No executable vendor configuration is emitted for any assistant, including this repository.

The hub pointer remains inert under [PDR-022](PDR-022-vendor-entry-points-only-point.md); the notice is the tool's cross-agent mechanism, not a vendor entry point's behavior.

## Rejected: emit a vendor session-start hook

A hook would take a position on an adopter's assistant, require vendor-specific schemas and migration behavior, and make the methodology less cross-agent. Keeping the notice in the tool and the fallback in the hub reaches ordinary work without shipping executable configuration.

## Carrier

CAP-004 carries the release check and notifier gate; CAP-001, ADR-018, the scaffolded hub, and migration behavior carry the adopter surfaces; `guide/operations.md` documents the notice and opt-out; this repository's hub carries the no-command fallback.
