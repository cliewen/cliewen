---
id: CAP-009
type: capability
status: active
links: [G-012, VIS-001, UC-001, ADR-065, ADR-066, ADR-067, PDR-054, PDR-055, C-023]
title: Product intent — one vision per corpus and optional use cases
goal: G-012
---

# CAP-009 — Product intent

## What

A corpus can state what the product is for, and can record an actor's path across capabilities when that path carries meaning the capabilities do not.

One `vision` artifact per corpus, at `docs/vision.md`, identity `VIS-001`. Zero or more `use-case` artifacts under `docs/use-cases/`. `clue validate` checks that what exists is well-formed and connected — a single vision at its address, a use case in its folder with a matching filename, naming at least one goal and at least one capability, and carrying its four structural sections. It requires neither artifact to exist and reports no coverage figure for either.

`clue validate --intent` states what the corpus has: the vision, its status and provenance, and the use cases with the capabilities each crosses. `clue context <id>` names — identity, title, and path only — the use cases whose links reach the artifact being read, which is how a capability's governing journey is reachable without a duplicate edge and without reverse traversal.

`clue init` writes a marked vision bootstrap into a new repository and the empty use-case folder. `clue migrate` reports a missing vision as a notice and writes no vision content.

## Why

The thread starts at a goal, and a goal answers why someone wanted *that goal*. Nothing answered why the product exists, whom it serves, or what is deliberately excluded — so an agent orienting on a change had capability-local truth and no direction to weigh a judgement against, and the available substitute was to infer direction from implementation structure, which answers a question code cannot answer.

Below the goal there was a matching gap: a capability is a unit of what the system can do and a criterion is a unit of proof, and neither carries an actor's path *across* capabilities. Most behaviour has no such path worth writing down, which is why that artifact is optional and why nothing measures its absence ([PDR-054](../../decisions/PDR-054-use-cases-are-optional-and-no-requirement-artifact.md)).

The judge checks form and claims nothing about meaning. Whether a vision is the right direction, whether a use case represents real users, and whether every actor has been found are human judgements, and a check that appeared to make them would be worse than no check.

`active`: AC-162 through AC-167 carry focused positive and negative unit evidence. AC-168 is `Human` — whether an interview is proportionate and whether inferred content is honestly labelled is a judgement, and the acceptance brief is its proof.
