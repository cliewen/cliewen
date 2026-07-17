---
id: CH-022
type: change
status: open
links: [P-002, M-006]
title: clue scaffold — standalone index regeneration
---

# CH-022 — `clue scaffold`

## What

A new command, `clue scaffold [path]`, exposing the index-regeneration engine that already lives in the scaffold package (shared with `clue init`, per ADR-019): regenerate the taxonomy README index blocks from folder contents — hand-written entries whose targets survive are kept, missing entries are appended, prose outside the `clue:index` markers is never touched — and materialize nothing. Missing folder READMEs are reported, never invented; a path without a `docs/` tree is a loud error. `checkIndexes` remains the judge of the result.

New capability CAP-005 (index generation) anchors the criteria: AC-026 (regeneration keeps prose and reaches green) and AC-027 (scaffold touches only index blocks, creates nothing), each with a positive and negative test.

## Why

Serves P-002 / M-006 (exit criterion: "`clue scaffold` regenerates README index blocks from folder contents (prose above markers untouched); `checkIndexes` remains the judge"). The digest step of every change updates README index blocks; today that is hand work or a full `clue init` run in a repo that needs no materialization. The engine exists and is tested — this change is the thin command surface over it.

## Also riding (bookkeeping)

The capabilities README index still annotates CAP-001 as `draft`; it went `active` in CH-020. The annotation is corrected as index bookkeeping.
