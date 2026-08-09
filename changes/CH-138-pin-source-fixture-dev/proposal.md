---
id: CH-138
type: change
status: open
links: [CAP-003, AC-128, AC-129, ADR-011]
title: Preserve dev fixtures on tagged checkouts
---

# CH-138 — Preserve dev fixtures on tagged checkouts

## What

Make an explicitly injected `dev` version remain `dev` when `clue` starts. An unstamped build continues to infer a release from Go module metadata when available, then falls back to `dev`.

## Why

v0.14.0's tagged release checkout made AC-128's ordinary source fixture inherit `0.14.0` through the module-version fallback. Its intentionally unrelated `1.4.0` fixture skills then appeared to drift, so the release's test step failed despite the fixture being designed to exercise the non-release path.

## How

- Distinguish an absent build stamp from an explicit `dev` stamp before reading module metadata.
- Add focused unit coverage for the version-resolution cases and retain AC-128's assessment-scale fixture as the tagged-checkout regression path.
- Keep the pinned-release fixture unchanged; it remains the only fixture that exercises release skill-drift behavior.

## Plan

This change is explicitly plan-less. It repairs release packaging evidence after v0.14.0 rather than implementing a P-013 milestone; that campaign's remaining milestones do not own version-resolution behavior.
