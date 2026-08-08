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

The brownfield migration assessment asked for a committed per-criterion registry that a reviewer could read against a pull-request diff. ADR-049 deliberately made the target side of migration parity derived instead of editable, and ADR-054 made a durable extraction report derive its counts and mapping table from that same source manifest. Those mechanisms make the retained report trustworthy, but they do not deliver the requested registry. AN-017 therefore records the request as declined pending an explicit accepted limit. What review surface does extraction provide, and what does an adopter give up by not receiving a second per-criterion artifact?

## Decision outcome

**Cliewen does not create or require a committed per-criterion registry for an extraction pull request.** The durable extraction report's derived region, its pinned source manifest, and a clean `clue parity` result are the review surface. The report gives a reviewer a readable summary and mapping table; the manifest remains the inspectable per-criterion source, and parity establishes that it agrees with the derived target corpus.

This refuses a second stored representation of the criterion mapping. An adopter gives up a single committed document that renders every criterion as a diff-oriented registry. To inspect an individual criterion, the reviewer follows the report's manifest reference and reads the pinned manifest alongside the target corpus or the parity output. That is more navigation than a dedicated registry, and it remains an explicit cost rather than being hidden behind a claim that the report is the requested artifact.

The report is not automatically a pull-request attachment or a substitute for the human merge boundary. An extraction's ordinary proposal, rehearsal report, source manifest, durable analysis report, and command results remain the review material under PDR-020 and PDR-021. A repository may choose to attach generated parity output in its own workflow, but that operational choice does not create a new Cliewen corpus artifact or change what the deterministic judge asserts.

**Carrier:** CAP-003's extraction contract and design, the canonical `clue-extract` skill, the public adoption guide's extraction section, AN-017's re-derived gate register, and P-012/M-060 state this accepted limit. ADR-049 and ADR-054 continue to define the report and manifest mechanisms that replace the declined registry.

### Rejected: commit a generated per-criterion registry beside the report

A generated registry would be derived from the manifest or target corpus and would add a second review artifact without adding a second source of truth. It would still need regeneration, validation, and carrier guidance, and a reviewer could mistake it for the authoritative mapping. The manifest already supplies the per-criterion detail while the derived report supplies the readable summary; duplicating both is not simplification.

### Rejected: let the report claim to be the per-criterion registry

The report is deliberately a summary and mapping table, not an exhaustive rendered registry. Calling it the requested artifact would close the assessment gate by changing its words rather than by delivering the thing it asks for.
