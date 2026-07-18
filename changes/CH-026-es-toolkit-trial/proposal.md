---
id: CH-026
type: change
status: open
links: [P-002, M-007]
title: Trial Cliewen on es-toolkit
---

# CH-026 — Trial Cliewen on es-toolkit

## What

Run the second qualifying foreign-soil trial of Cliewen against [`toss/es-toolkit`](https://github.com/toss/es-toolkit) at commit `fd4877295443a92655530735b0058c1d9ba1db4c`.

Treat the target as read-only: use a disposable clean checkout, make no commits or pull requests in the external repository, and keep trial-only artifacts outside its tracked tree. Record the evidence, interpretations, verification results, conflicts, limits, and cross-trial findings as AN-005 in Cliewen.

Include the human-approved bookkeeping correction from the post-merge CH-025 review by changing its decision-log reference from `CH-025` to `CH-025 / PR #24`.

## Why

P-002/M-007 requires trials on at least two human-selected external open-source repositories with no shared maintainer. `es-toolkit` is a TypeScript library monorepo with its own `AGENTS.md`, compatibility surface, benchmarks, bundle-size constraints, multilingual documentation, and broad CI. It tests Cliewen against a materially different repository and exposes how methodology rules interact with pre-existing agent instructions.

The human selected `es-toolkit` as the second repository after `sharkdp/hyperfine`. Completing this trial, publishing AN-005, and reconciling both trials can satisfy M-007 because CH-025 already delivered the required methodology adjustment.

## Decision boundary

This change may make another narrowly evidenced methodology adjustment when the trial exposes a reusable defect or omission. It must trace that adjustment to AN-005 and must not change methodology merely to manufacture additional findings. Any larger or unresolved choice is recorded as an open question and stops the change.
