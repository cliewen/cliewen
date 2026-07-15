---
id: CH-015
type: change
status: open
links: [P-002, M-004]
title: Cut release 0.2.0 — the accumulated methodology changes ship
---

# CH-015 — Cut release 0.2.0

Serves [P-002](../../docs/plans/P-002-leaves-home.md)/M-004 (clue ships): five methodology changes have accumulated unreleased — the light change tier, the decision log, the convention register, merge-binds/approval-signs, and the typed decision records (ADR/PDR). Releasing them is also the first real run of the release-cut mechanics: 0.1.0's changelog section was written retroactively, so this change is the precedent.

## What

- `CHANGELOG.md`: the `[Unreleased]` section becomes `[0.2.0] - 2026-07-15`, with an install line matching 0.1.0's; the release workflow publishes this section verbatim as the GitHub release body.
- The five skill version stamps bump `0.1.0` → `0.2.0` — the skills and the binary version as one set (drift is a failure for released binaries).
- After merge: tag `v0.2.0` on the merge commit; the release workflow builds the stamped cross-platform binaries and publishes the release.

This change is full-loop only because it touches skill files (the version stamps); it decides nothing and changes no rule text — releasing is executing M-004, and the version choice (minor bump, features added, still 0.x) is plain semver.
