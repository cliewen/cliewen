---
id: CAP-002
type: capability
status: active
links: [G-001]
title: clue validate — the deterministic judge for the corpus
goal: G-001
---

# CAP-002 — `clue validate`

## What

`clue validate [--forbid-changes] [--coverage] [--reality-gaps] [--index-rows] [--read-cost] [path]` scans `docs/` and `changes/` for frontmatter artifacts and fails (exit 1) on any breach of the corpus rules: missing frontmatter or core fields, UTF-8 byte-order marks, leftover second frontmatter blocks, duplicate IDs, unresolvable `links` (milestones and live acceptance criteria also resolve), status values outside the per-type vocabulary, unbounded high-cost inferred meaning behind active capabilities, malformed incident edges, `/docs` folders without README.md, index-block drift, the convention register's fields and declarations, hard-wrapped prose, unexplained skipped tasks, a proposal that declares no plan item, missing per-type frontmatter, milestone statuses outside the vocabulary, a live artifact whose ID is absent from the persisted identity ledger (`.clue/id-ledger.yaml`) or not marked live there, a malformed ledger entry shape or identity, and — with `--forbid-changes` — the presence of a transient workspace. Image links and assets are valid corpus content; `clue validate` does not fetch their local or external targets. The optional reports derive proof coverage, capabilities contradicted by incident analyses, index rows that say nothing a reader can use, and structural read-cost backlogs, without committing registries.

After `clue migrate --apply` has backfilled that ledger, `clue id next <prefix>` allocates the next numeric ID for a prefix through it — an O(1) counter increment, never a corpus scan — and `clue id live <id>` promotes the reservation after its artifact is created, so a retired ID is never silently reissued ([ADR-048](../../decisions/ADR-048-corpus-wide-id-ledger.md)).

## Why

The judge actor of [G-001](../../goals/G-001-verifiable-thread.md): machines enforce form so humans only have to verify meaning. The same binary runs locally and in CI — a green local run means a green PR.

That guarantee is what fixes the judge's edge. It reads the repository as it stands and never a diff, a base revision, or history: a verdict computed from a transition would depend on which branch the caller is standing on and how deep the clone is, and two people with identical bytes would get different answers ([ADR-044](../../decisions/ADR-044-judge-reads-state-not-transitions.md)). Rules about what a *change* did are held by machines that are allowed to have a base — CI, the release gates, branch protection — and the constraint register names which one holds each.

Acceptance criteria: [criteria.md](criteria.md) · implementation notes: [design.md](design.md).

## Status note

`active`: implemented and covered by Go tests whose names carry the AC IDs and evidence classifications. AC-009 and its positive and negative unit tests enforce the active-criterion reference rule; criteria declaring a machine proof type additionally require supported evidence classified by that type and direction, unannotated legacy criteria retain one reference, Human-class criteria route proof to the acceptance brief, and `@draft` exempts one criterion. AC-135 reports multi-document artifacts and over-budget default context slices as a non-blocking structural backlog. JVM evidence is credited only when the AC identity, type, and direction belong to the same supported Java or Kotlin executable. `clue validate` validates these declarations and references but does not execute tests.
