---
id: IDR-008
type: decision
status: inferred
links: [C-014, CAP-008]
title: Go coverage is an aggregate tripwire
author: agent
accepted-by: []
---

# IDR-008 — Go coverage is an aggregate tripwire

## Context

Acceptance criteria hold behavior directly, but code between their focused paths still needs a backstop that does not force artificial tests into thin entry-point packages.

## Decision

Repository verification enforces at least 80% total Go statement coverage. The threshold is aggregate rather than per package and remains a tripwire behind the AC-to-evidence contract, not a substitute for focused criterion evidence.
