---
id: ADR-056
type: decision
status: inferred
links: [CAP-007, PDR-034, ADR-007]
title: A context slice is bounded by default, and states what the bound held back
author: agent
accepted-by: []
---

# ADR-056 — Bounded context slice

## Context and problem statement

The prior context read followed the whole resolvable graph, making routine orientation expensive even when a reader needed only an identity and its immediate dependencies. A bounded read must remain honest about what it omitted.

## Decision outcome

**A context slice defaults to one link hop from the root and reports the frontier beyond that bound; `--depth` widens the bound and `--depth=all` preserves exhaustive traversal.** Included artifacts are printed whole. The frontier names only artifacts one hop out and counts the rest, and an unfollowed edge is reported only when it leaves an included artifact. Zero depth remains available for the root alone.

The default is a reading cost, not a semantic limit: it changes neither corpus meaning nor validation, and widening is the reader's judgment based on the named frontier.

**Carrier:** `clue context`, the focused-context capability, and the frontier output and guidance.
