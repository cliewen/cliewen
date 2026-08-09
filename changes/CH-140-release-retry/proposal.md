---
id: CH-140
type: change
status: open
links: [P-013, PDR-015, CAP-004]
title: Release cuts are short and an unpublished version can be retried
---

# CH-140 — Release cuts are short and an unpublished version can be retried

## What and why

Release preparation is administrative work, not a Cliewen change: it needs a small reviewed release PR, not a CH workspace, acceptance brief, or digest. A tag that has no GitHub Release must remain recoverable so a repaired main commit can publish the same version rather than forcing a patch-number bump.

## Decision already accepted

The human chose a small release PR, automatic retagging of an unpublished version, and GitHub Release existence as the permanent-publication boundary in the 2026-08-09 conversation.
