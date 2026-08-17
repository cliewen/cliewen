---
id: IDR-003
type: decision
status: inferred
links: [ADR-042, CAP-004]
title: Release-list nonanswers are cached briefly
author: agent
accepted-by: []
---

# IDR-003 — Release-list nonanswers are cached briefly

## Context

Repeated workflow commands should not repeatedly pay for the same offline, rate-limited, timed-out, or unrecognized release-list result, while an explicit check must be able to retry immediately.

## Decision

The release cache records a nonanswer separately from a version answer and gives it a shorter lifetime. Automatic notices honor that cached nonanswer; an explicit `clue latest` request always asks again. A corrupt, unreadable, or self-contradictory cache is absence and never suppresses a request.
