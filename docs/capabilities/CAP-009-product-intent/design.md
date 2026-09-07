---
id: CAP-009-design
type: design
status: active
links: [CAP-009, CAP-002, CAP-007, ADR-065, ADR-066, ADR-067]
title: Product intent design
---

# Product intent design

## What the judge checks, and what it refuses to

`corpus.checkIntent` runs beside the other graph rules and is deliberately shallow. It reads the scanner's existing artifact list — `type`, `Path`, `Links`, `Body` — and adds no parsing of its own beyond looking for four literal section headings in a use-case body.

The vision half asserts four things: at most one artifact of type `vision` exists; it sits at `docs/vision.md`; the file at that address, if present, is a vision; and it no longer carries the scaffold bootstrap marker. Nothing about its length, its headings, or its sentences is checked. A vision that says the wrong thing is a well-formed vision, and pretending otherwise would be the most expensive claim this capability could make.

The use-case half asserts placement (`docs/use-cases/UC-<digits>-<slug>.md`, filename prefixed by the identity), that the links name at least one `G-` and at least one `CAP-` identity, and that the body carries `## Actors`, `## Trigger`, `## Main flow`, and `## Outcome`. Link *resolution* is not re-checked here: `checkLinks` already does it for every artifact, and a second rule would report the same defect twice. Nothing anywhere requires a use case to exist, and no rule reads a goal or capability looking for one.

The scaffold bootstrap rule fires only when the file exists, which is what lets a repository with no vision stay green while a repository that ran `clue init` and stopped is told to finish.

## The intent report

`clue validate --intent` derives `corpus.IntentState` and prints it. It is a state report, not a scorecard: the vision line names the identity, status, and whether `provenance` marks the meaning inferred; each use-case line names the identity, status, and the capability identities it crosses, in declaration order. A corpus with no vision prints one line saying so and exits with whatever the validation verdict was — the absence is state, never an issue.

There is no figure. The report computes no ratio of goals or capabilities with use cases, because a percentage over an optional artifact reads as a target and the only way to move it is to write artifacts nobody needs.

## Reaching a use case from a capability

Links point down: a use case names its goal and its capabilities, a goal names the vision, and a capability's links are unchanged. That leaves one direction unserved, and `clue context` serves it without a second edge.

After printing the slice, the command scans the corpus for artifacts of type `use-case` whose `links` contain the root's identity, and prints their identity, title, and repository-relative path, ordered by path. It follows nothing and emits no content, so the slice, its byte count, and its frontier are exactly what they would have been. The scan is by artifact type rather than by link direction in general: a generalized reverse-link mode returns almost nothing for a leaf and most of the corpus for a goal or a vision, which is precisely backwards.

A criterion or milestone identity given to `context` resolves to its owning artifact first, so naming is computed against the artifact the reader actually received.

## Materializing and migrating

`clue init` writes `docs/vision.md` from a template carrying `<!-- clue:vision:bootstrap -->` and an empty `docs/use-cases/README.md` index; both go through `writeIfAbsent`, so a repository that already has either keeps it. The corpus index picks both up through the ordinary regeneration — a sibling `.md` and an artifact-bearing subfolder are shapes `indexTargets` already knows.

`MIG-014` splits along what can be proven. The use-case folder README is structure with no meaning in it, so migration creates it and adds its index row the way `MIG-011` does for the overview folders. The vision is not: nothing in a repository proves why a product exists, so migration emits a notice naming the absence and the workflows that fill it, and never writes the file. The notice blocks nothing — a plan carrying it still applies — which is what keeps an established corpus from going red for an artifact its authors could not have written.
