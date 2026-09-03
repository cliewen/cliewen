---
id: ADR-065
type: decision
status: inferred
links: [G-012, P-021, M-087, VIS-001, CAP-009, CAP-002, ADR-010, ADR-035, ADR-025, ADR-048]
title: A corpus has one vision, at a fixed address, and its lifecycle is the one the corpus already has
binds: adopter
author: agent
accepted-by: []
---

# ADR-065 — One vision, at a fixed address

## Context

A corpus needs one statement of what the product is and whom it serves. Three things about it have to be settled before it can exist: how many there are, how it is addressed, and how a draft is told apart from meaning a human has confirmed.

The corpus rule is normally that identity is the ID and the path is only the current address. A vision is the case where that rule buys nothing. There is exactly one, so nothing has to be searched for; an agent orienting wants to read it without first scanning frontmatter to find out whether it exists.

## Decision

**A corpus has at most one artifact of type `vision`, it lives at `docs/vision.md`, and its identity is `VIS-001`.** `clue validate` rejects a second vision, a vision anywhere else, and a `docs/vision.md` that is not one. The identity is still what links name — a goal names `VIS-001`, never a path — so the graph is unchanged; the fixed address is an additional guarantee, not a replacement for the ID.

A file rather than a folder. `indexTargets` and `checkIndexes` already carry a sibling `.md` directly under `docs/`, so one file needs no new index machinery, while a folder would need a README and an index block to hold a single artifact. That is ceremony a singleton cannot repay.

**The lifecycle is the corpus's existing one, with no new vocabulary.** The artifact runs `draft` → `active` like every other default-lifecycle type ([ADR-025](ADR-025-one-status-lifecycle.md)). Agent-drafted content additionally carries `provenance: inferred` with a `reversal-cost` ([ADR-010](ADR-010-decision-provenance.md), [ADR-035](ADR-035-inferred-meaning-routes-by-reversal-cost.md)), which is the corpus's existing way of saying "no human has confirmed this yet". Promotion to `verified` is a human act, exactly as it already is elsewhere.

**One vision per corpus, and no hierarchy.** A repository governing several independent products runs several corpora. Introducing a per-product vision namespace would add a scoping concept to every link, index, and context slice to serve a case no adopter has yet reported; a repository that reaches it can split, which is the cheaper move and the one that keeps each product's acceptance boundary whole.

## Rejected: give the vision a numbered namespace with no singleton rule

Allowing `VIS-002` costs nothing structurally and immediately makes every other question ambiguous — which vision a goal serves, which one an agent reads while orienting, which one a report means. The value of a vision is that there is one; a namespace that permits several removes it while looking harmless.

## Carrier

The `vision` form checks in `internal/corpus`, their criteria under `docs/capabilities/CAP-009-product-intent/`, the bootstrap in `internal/scaffold/templates/docs/vision.md`, the intent-model instruction in `internal/skills/source/shared/intent-model.md.tmpl`, and `guide/intent.md`.
