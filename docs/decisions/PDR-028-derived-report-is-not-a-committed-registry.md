---
id: PDR-028
type: decision
status: verified
links: [P-012, ADR-049, ADR-054, PDR-020, PDR-021, AN-017]
title: A derived extraction report is not a committed per-criterion registry
author: agent
accepted-by: Flemming N. Larsen (2026-08-08, conversation)
---

# PDR-028 — A derived extraction report is not a committed per-criterion registry

## Context and problem statement

The assessment requested a committed per-criterion registry, while ADR-049 and ADR-054 deliberately make the migration target and report mapping derived from a pinned source manifest; AN-017 records the request as declined pending an accepted limit.

## Decision outcome

**An extraction does not create or require a committed per-criterion registry.** The durable report's derived region, pinned source manifest, and clean `clue parity` result are the review surface: the report gives a readable summary, the manifest remains the per-criterion source, and parity confirms agreement with the target corpus.

The cost is explicit: inspecting an individual criterion requires following the manifest reference rather than opening one diff-oriented registry. The report is not a pull-request attachment or a replacement for the proposal, rehearsal, analysis, source manifest, command results, or human merge boundary. A repository may attach generated parity output without creating a new Cliewen artifact.

## Rejected: commit or rename a generated registry

A generated registry would duplicate the manifest, add regeneration and drift risk, and make a second representation look authoritative; calling the summary a registry would change its words without delivering the requested artifact.

## Carrier

CAP-003, the canonical `clue-extract` skill, the public adoption guide, AN-017, P-012/M-060, ADR-049, and ADR-054 carry this accepted limit and its replacement review surface.
