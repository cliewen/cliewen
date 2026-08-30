---
id: PDR-051
type: decision
status: inferred
links: [P-018, PDR-019, PDR-043, PDR-046, CAP-001, CAP-002, CAP-006]
title: Every Cliewen corpus maintains concise system architecture and design overviews
author: agent
accepted-by: []
---

# PDR-051 — Every Cliewen corpus maintains concise system architecture and design overviews

## Context

Capability design and decisions can be correct while the system picture is absent, stale, or spread across unrelated documents. A second historical log would make the corpus larger without making the current system easier to understand.

## Decision

Every Cliewen corpus has canonical, concise, current-truth overviews at `docs/architecture/README.md` and `docs/design/README.md`. Architecture describes structure, boundaries, actors, and durable technology choices. Design describes cross-cutting flows, interactions, and patterns. Capability `design.md` remains the home of capability-local detail.

Agents inspect repository evidence, draft the big picture rather than exhaustive detail, and ask the user when a material boundary or intent remains unclear. Every change assesses documentation impact and updates only the overview, design, capability, or decision material needed to keep the corpus true; it records its disposition in the change or pull-request handoff, not as a permanent history document. A Mermaid or retained SVG diagram is used when it materially clarifies a relationship, boundary, or flow.

The judge holds the canonical paths and explicit bootstrap state. Lifecycle guidance holds truthful drafting, relevance, consent, and link-update judgment. On discovery of a suitable existing overview outside the convention, an agent presents one grouped relocation mapping and moves content only with explicit user consent; a declined move leaves the source in place and a concise canonical pointer. Relevant new or materially revised ADRs and IDRs link to the affected overview; PDRs do so only when they govern this methodology.

## Carrier

The corpus overview READMEs, CAP-001/CAP-002/CAP-006, their criteria and designs, the validator and migration planner, the scaffold and generated lifecycle skills, the routing hub, extraction mapping, guidance, and release notes carry this contract. [ARCH-003](../architecture/core.md) records why the validator portion crosses the red line.
