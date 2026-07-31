---
id: CH-090
type: change
status: open
links: [P-009, M-039, PDR-007, PDR-016, ARCH-001, ARCH-003]
title: Define the supported merge history for Cliewen changes
---

# CH-090 — Define the supported merge history for Cliewen changes

## What

This full change serves P-009 milestone M-039 by reproducing merge-commit, squash, and rebase outcomes for a complete Cliewen change, then recording which merge modes preserve the accepted proposal, implementation, digest, and durable corpus history in the repository itself. It aligns the PDR, operations and setup guidance, contributor workflow, architecture, generated agent guidance, and branch-protection probe around one explicit support boundary.

## Why

Cliewen treats a human merge as acceptance and Git history as the provenance archive, but the current guidance does not say which pull-request merge shapes preserve that archive. A squash merge can retain the net tree while discarding the branch commits that carried the proposal and implementation, so presenting merge, squash, and rebase as equivalent would leave the acceptance record dependent on forge state.

## Decision boundary

The change decides only the repository-history contract for a full Cliewen change and the configuration needed to enforce it. It does not change the human merge boundary, add a forge as a system of record, alter the one-change-in-flight rule, or introduce stacked changes or a general repository configuration file.
