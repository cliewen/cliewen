---
id: IDR-005
type: decision
status: verified
links: [ADR-009, ADR-029, PDR-046, CAP-003]
title: MADR conversion preserves meaning and resolves identities deterministically
author: agent
accepted-by: Flemming N. Larsen (2026-09-02)
---

# IDR-005 — MADR conversion preserves meaning and resolves identities deterministically

## Context

A source decision corpus can mix proposed, accepted, rejected, deprecated, and superseded records, reuse numeric prefixes across folders, and include proposals the repository never adopted.

## Decision

The MADR mapping preserves each source status as meaning without manufacturing Cliewen approval: an acted-on proposal may convert, an untouched proposal becomes an open question, and rejection, deprecation, and supersession remain explicit. Subject selects ADR, PDR, or IDR. Preserved numbers are placed within their routed namespace first; collisions keep the path-sorted first record and mint later identities only after preserved numbers are known.
