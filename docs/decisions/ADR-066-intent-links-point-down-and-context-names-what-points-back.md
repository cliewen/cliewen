---
id: ADR-066
type: decision
status: inferred
links: [G-012, P-021, M-088, UC-001, CAP-007, CAP-009, PDR-034, ADR-056]
title: Intent links point down the composition, and context names what points back without following it
binds: adopter
author: agent
accepted-by: []
---

# ADR-066 — Intent links point down, and context names what points back

## Context

A use case is only worth having if an agent can reach it from where a task starts, which is usually a capability or a criterion. The corpus reads narrowly and outward ([PDR-034](PDR-034-the-corpus-is-read-narrowly-by-default.md)), so reachability means an outgoing link, and every candidate direction has a cost.

If a capability names its use cases and a use case names its capabilities, one edge is written in two files and drifts the moment either is edited. If only the capability names the use case, the journey cannot enumerate its own parts. If only the use case names the capabilities, the capability cannot reach the journey that governs it.

## Decision

**Intent links follow composition, downward.** A use case names the goal it serves and every capability it crosses; a goal names the vision. A capability's links are unchanged and never name a use case. Each edge is written once, in the artifact whose meaning depends on the target, which is the direction the rest of the corpus already uses.

**`clue context` names the use cases whose links reach the root, by identity, title, and path only.** No content is emitted for them and no edge is traversed, so the slice's size and bound are unchanged and widening remains the reader's explicit act ([ADR-056](ADR-056-bounded-context-slice.md)). This is what makes a capability's governing journey reachable without a second copy of the edge and without a reverse traversal that would pull a shared vision's whole dependent graph into every read.

The naming is bounded by kind rather than by depth: only `use-case` artifacts are reported, and only those naming the root directly. A vision is deliberately not reported this way. Every goal names it, so the answer would be the corpus, which is the read this bound exists to prevent.

## Rejected: a general reverse-link mode on `clue context`

An `--incoming` flag reads as the natural generalization and is the wrong shape. For a leaf it returns almost nothing, and for a goal, constraint, or vision it returns most of the corpus — so the flag is useless exactly where a reader would reach for it, and dangerous exactly where it looks most useful. Restricting the answer to one artifact type is what keeps it bounded by construction rather than by the caller's care.

## Carrier

The use-case link rules in `internal/corpus`, the naming step in `clue context`, their criteria under `docs/capabilities/CAP-009-product-intent/`, the traversal instruction in `internal/skills/source/shared/intent-model.md.tmpl`, and `guide/intent.md`.
