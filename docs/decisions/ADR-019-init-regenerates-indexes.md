---
id: ADR-019
type: decision
status: inferred
links: [ADR-013, CAP-001, P-002]
title: Index regeneration runs in clue init; ADR-013's emits-empty clause is superseded
author: agent
accepted-by:
---

# ADR-019 — Index regeneration is part of init

## Context and problem statement

[ADR-013](ADR-013-ships-generic-vs-repo-local.md) gave the README index blocks a split ownership: `clue init` emits them empty, and `clue scaffold` (M-006) regenerates them from folder contents. CH-020 then built init's green-after-init contract for existing repositories: a pre-existing `docs/` tree reaches a green `clue validate` only if the taxonomy README index blocks reference its artifacts — which requires regeneration (and marker-append for READMEs that predate init) at init time, not in a later command. Keeping regeneration exclusive to a future `clue scaffold` would leave init red on exactly the repositories the onboarding promise covers.

## Decision outcome

**`clue init` regenerates the taxonomy README index blocks on every run. The emits-empty/regenerates-later split in ADR-013 is superseded in that one clause; the rest of ADR-013 stands.**

- The regeneration engine lives in the scaffold package and follows the contract `checkIndexes` judges: hand-written single-line entries whose targets survive are kept, plain entries are appended for anything missing, prose outside the markers is never touched (AC-024/AC-025 hold this in tests).
- M-006 (`clue scaffold`) remains open as the standalone exposure of the same engine — regeneration without materializing anything — and `checkIndexes` remains the judge either way. The milestone's exit criterion is unchanged.
- On a fresh scaffold the templates' index blocks already match what regeneration produces, so the first run reports nothing regenerated — idempotence is observable from the first command.
