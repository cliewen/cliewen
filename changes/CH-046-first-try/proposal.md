---
id: CH-046
type: change
status: open
links: [P-004, M-013, CAP-001, QS-002]
title: Make Cliewen's first proof concrete
---

# CH-046 — Make Cliewen's first proof concrete

## What

Complete P-004/M-013 by making Cliewen's value concrete before its terminology, making checksum-verified release binaries the primary installation path, and giving newcomers a disposable demonstration that first validates green and then exposes a real missing-test thread when a draft criterion is activated without evidence. Show the created tree, expected command output, file ownership boundaries, recovery, and cleanup. Add a compact request-to-merge narrative to the change-loop guide and link the real pull request that introduced AC↔test validation.

Update CAP-001's accumulating guide requirements and the user-facing changelog, then close M-013 in the digest when the documented path has been rehearsed from a clean temporary repository and the strict guide build passes.

## Why

The existing guide explains the methodology but makes a newcomer learn its abstractions and consider changing an existing repository before seeing the judge catch anything. Its source-first installation order also makes a cross-platform single binary look Go-specific. M-013 earns the first try by moving the concrete outcome, safe experiment, and observable failure ahead of the deeper method.

## Decision boundary

This change implements the already accepted M-013 scope. It does not alter `clue`, validation semantics, skills, scaffold templates, QS-002, CI enforcement, compatibility claims, package-manager distribution, or operations guidance. The demo must state that `clue validate` enforces at least one AC reference, while positive and negative test pairing remains an agent/human requirement. Any need to cross that boundary becomes an open question and stops the change.
