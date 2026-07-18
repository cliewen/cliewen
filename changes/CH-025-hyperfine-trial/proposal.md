---
id: CH-025
type: change
status: open
links: [P-002, M-007]
title: Trial Cliewen on hyperfine
---

# CH-025 — Trial Cliewen on hyperfine

## What

Run the first qualifying foreign-soil trial of Cliewen against [`sharkdp/hyperfine`](https://github.com/sharkdp/hyperfine) at commit `f12f3d9f86f3643b3b7deace5e160b1f0f44d2b7`.

Treat the target as read-only: use a disposable clean checkout, make no commits or pull requests in the external repository, and keep any trial-only artifacts outside its tracked tree. Record the evidence, interpretations, verification results, limits, and methodology findings as AN-004 in Cliewen.

## Why

P-002/M-007 requires trials on at least two human-selected external open-source repositories with no shared maintainer. `hyperfine` is a compact Rust command-line application whose statistical measurements, platform-specific behavior, documentation, tests, and CI exercise Cliewen's ability to distinguish deterministic acceptance criteria from quality evidence while remaining practical on an unfamiliar codebase.

The human selected `hyperfine` as the first qualifying repository and `toss/es-toolkit` as the second. Completing this change counts as one trial; it does not close M-007 by itself.

## Decision boundary

This change may make a narrowly evidenced methodology adjustment when the trial exposes a reusable defect or omission. It must trace that adjustment to AN-004 and must not change methodology merely to manufacture the M-007 adjustment requirement. Any larger or unresolved choice is recorded as an open question and stops the change.
