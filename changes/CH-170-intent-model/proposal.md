---
id: CH-170
type: change
status: open
links: [P-021, M-087, M-088, M-089, M-090, G-012, CAP-009]
title: An intent layer — one vision per corpus and optional use cases
---

# CH-170 — An intent layer: one vision per corpus and optional use cases

## What

Add two durable artifact types and the agent workflows that produce them:

- **Vision** — one per corpus, at the fixed address `docs/vision.md`, identity `VIS-001`, stating what the product is, whom it serves, what is in and out of scope, what would count as succeeding, and what remains uncertain.
- **Use case** — optional, `docs/use-cases/UC-xxx-slug.md`, describing one actor's end-to-end path across capabilities. Zero is a valid number.

With them: link direction and bounded traversal, deterministic form checks that make no semantic claim, a `--intent` report, a migration that reports rather than invents, greenfield interview and brownfield inference in the shipped skills, and the guide material that explains when *not* to write a use case.

## Why

Serves P-021 (M-087 … M-090) and G-012. The thread starts at a goal and answers why someone wanted that goal; nothing answers why the product exists. An agent orienting on a change therefore has capability-local truth and no direction to weigh a judgement against, and the tempting substitute — inferring intent from implementation structure — answers a question code cannot answer.

## Shape, and what it deliberately is not

This is not a requirements layer. No generic `requirement` artifact is introduced, use cases never restate acceptance criteria, and the delivery artifacts (plan, milestone, change) stay out of the semantic hierarchy — [PDR-054](../../docs/decisions/PDR-054-use-cases-are-optional-and-no-requirement-artifact.md) records both refusals.

Three shapes were reconciled against 0.22.0 rather than adopted from the request:

- The vision is a **single file**, not a folder. `indexTargets` and `checkIndexes` already carry a sibling `.md` under `docs/`; a folder would add a README and an index block to hold one artifact ([ADR-065](../../docs/decisions/ADR-065-the-vision-is-a-singleton-at-a-fixed-address.md)).
- Lifecycle reuses what exists: `draft` → `active` for the artifact, `provenance: inferred|verified` with `reversal-cost` for agent-drafted meaning. No new status vocabulary ([ADR-010](../../docs/decisions/ADR-010-decision-provenance.md), [ADR-035](../../docs/decisions/ADR-035-inferred-meaning-routes-by-reversal-cost.md)).
- Links point **down**: a use case names its goal and the capabilities it crosses; a goal names the vision. A capability does not name its use case, so no edge has to be kept in step in two files. Reaching a use case from a capability is a bounded, names-only report from `clue context`, not a reverse traversal ([ADR-066](../../docs/decisions/ADR-066-intent-links-point-down-and-context-names-what-points-back.md)).

## Compatibility

A 0.22.0 corpus with no vision and no use cases stays valid. `clue migrate` writes no vision — nothing in a repository proves why a product exists — and reports the absence as a notice that blocks nothing. `clue init` writes a bootstrap for a *new* repository, which validate rejects until replaced, which is the existing architecture/design overview pattern and is what separates "not yet established" from "accidentally omitted" ([ADR-067](../../docs/decisions/ADR-067-a-corpus-without-a-vision-stays-valid.md)).

## Plan item

P-021, milestones M-087 through M-090. The plan is created in this same proposal commit: P-020 closed naming no successor, and a proposal cannot link a plan item that does not resolve.
