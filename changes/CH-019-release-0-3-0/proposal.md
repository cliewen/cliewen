---
id: CH-019
type: change
status: open
links: []
title: Cut release 0.3.0 — the review boundary and post-merge orientation ship to adopters
---

# CH-019 — Cut release 0.3.0

Plan-less, explicitly: a routine release cut. The three methodology changes accumulated since 0.2.0 — the review-boundary prohibitions, the digest-is-never-a-task rule (both CH-017), and post-merge orientation (CH-018) — all correct agent misbehavior observed in the adopter repo (model2diagram), and that repo runs against a vendored, drift-checked 0.2.0 pair: the fixes reach it only through a release. No milestone covers routine release cuts (M-004, the release *mechanics*, closed with 0.2.0).

## What

- `CHANGELOG.md`: the `[Unreleased]` section becomes `[0.3.0] - 2026-07-17`, with an install section matching 0.2.0's; the release workflow publishes this section verbatim as the GitHub release body.
- The five skill version stamps bump `0.2.0` → `0.3.0` — skills and binary version as one set (drift is a failure for released binaries, per the versioning decision recorded in ADR-011). The binary code is unchanged since v0.2.0; re-tagging it at 0.3.0 is how a skills-only release keeps the pair matched.
- After merge (human acts, not tasks): tag `v0.3.0` on the merge commit and push the tag; the release workflow builds the stamped cross-platform binaries and publishes the release. Then model2diagram updates its vendored binary + skills in its own change.

This change is full-loop only because it touches skill files (the version stamps); it decides nothing and changes no rule text — the version choice (minor bump: behavior rules added, still 0.x) is plain semver, following the 0.2.0 precedent.
