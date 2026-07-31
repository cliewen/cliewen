---
id: CH-088
type: change
status: open
links: []
title: Refresh the local clue only for Cliewen releases
---

# CH-088 — Refresh the local clue only for Cliewen releases

## What

This explicitly plan-less, full change adds a repository-local release convention to `AGENTS.md`: when releasing Cliewen itself, update the locally installed `clue` from this checkout with `go install ./cmd/clue` before release verification. Ordinary changes and work in adopter repositories do not trigger that local installation update.

The digest records the cheap, local decision in the decision log. The generic scaffold `AGENTS.md` template, distributed skills, release workflow, and `clue` behavior remain unchanged because this convention belongs only to Cliewen's own repository.

## Why

This checkout carries the current generated skills and source validator while a user's PATH may still contain an older released `clue`. Release verification must use the checkout's judge rather than accidentally invoking that stale binary, but ordinary repository work should not mutate a user's local installation as a side effect.

## Decision boundary

The rule is repository-local release housekeeping. It does not change the shipped CLI, the release artifact, the adopter contract, the validator's drift semantics, or the generic instructions emitted by `clue init`.
