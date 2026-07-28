---
id: ADR-035
type: decision
status: inferred
links: [ADR-010, PDR-004, CAP-002, CAP-003, AN-007, P-007]
title: Cost bounds inferred provenance and incident analyses return an edge from reality
author: agent
accepted-by: []
---

# ADR-035 — Cost bounds inferred provenance and incident analyses return an edge from reality

## Context and problem statement

`provenance: inferred` distinguishes agent-reconstructed meaning from human-verified meaning, but treating every inferred artifact as one backlog makes the signal monotonically increase and eventually become ignorable. Decisions express the same distinction in `status`, yet the existing count omits them. The graph also has no stated representation for the case where every corpus check passed and later evidence showed that a capability or criterion was wrong.

## Decision outcome

**An inferred non-decision artifact declares `reversal-cost: low | high`.** The field records the semantic routing judgment already used for decisions: cheap and local to reverse is `low`; expensive to reverse is `high`. It is required while `provenance: inferred` and optional after verification so promotion need not erase useful classification. A present value outside that vocabulary fails validation. Decisions do not carry the field: an ADR or PDR is already the high-cost route, while a cheap decision lives as a row in the active decision log.

**High-cost inferred meaning bounds capability activation.** For an active capability, its activation slice is the capability itself plus every live artifact joined to it by one `links:` edge in either direction. A high-cost non-decision artifact in that slice cannot remain inferred; validation names both the artifact and the active capability it blocks. Low-cost inferred artifacts remain valid indefinitely. The one-edge boundary is deliberate: it is deterministic, inspectable in frontmatter, and makes the owner state the dependency instead of asking the judge to infer transitive semantics.

**The CLI reports costly unverified meaning as two actionable populations.** High-cost inferred non-decision artifacts are counted as activation blockers; inferred ADRs and PDRs are counted as decisions awaiting verification. Decisions remain non-blocking under [PDR-004](PDR-004-merge-binds-approval-signs.md): merge makes them binding and only explicit approval signs them, so turning their unsigned state into a validation failure would contradict that boundary. The old count of every inferred non-decision artifact is removed; low-cost deferred findings no longer create permanent counter noise.

**An analysis records “green met wrong” with `reality: contradicted` and a `links:` edge to every capability or acceptance criterion whose claim failed.** Acceptance-criterion IDs therefore resolve as graph links alongside artifact and milestone IDs. The marker is valid only on an analysis and requires at least one capability or live criterion link. `clue validate --reality-gaps` derives a sorted capability listing from those edges, mapping a criterion through its owning criteria artifact to its capability; no incident registry is committed. The analysis may also link the decisions, constraints, or process carriers that failed to prevent the incident, preserving the shape AN-007 established.

This is an incident convention, not the production feedback loop. Cliewen records an observed contradiction after it reaches the repository; it does not ingest telemetry, operate deployments, or automatically create goals or constraints. The V3 production door remains closed.

**Carrier:** `checkProvenance`, the active-capability slice rule, acceptance-criterion link resolution, and the derived reality-gap view in `clue`; the incident shape in `clue-analysis`; the field and status descriptions in the corpus and public guide.

### Rejected: block every inferred artifact

Some findings legitimately remain unverified because the authoritative maintainer is external or unavailable, and reversing them is cheap. Blocking them recreates the monotonic backlog as a hard gate instead of making it useful.

### Rejected: make inferred decisions fail validation

That reverses “merge binds, approval signs” by making a separate approval ritual a prerequisite for a green change. Decisions stay visible in the costly count without pretending that signature status changes whether the repository is allowed to use them.

### Rejected: infer incidents from prose or production systems

Parsing headings or titles is nondeterministic, while ingesting production feedback opens the explicitly deferred operations door. One frontmatter marker and ordinary graph edges give the judge a stable local fact without adding orchestration.
