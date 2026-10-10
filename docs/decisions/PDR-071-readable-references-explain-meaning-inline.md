---
id: PDR-071
type: decision
status: inferred
links: [G-026, CAP-013]
title: Human-facing references explain relevant meaning inline
author: agent
accepted-by: []
binds: adopter
---

# References explain the relevant meaning

## Context

An unexplained identity makes a reader navigate before understanding a document, report or handoff. Adding a link or reference table alone does not state why its target matters. Canonical identities remain necessary for machine traceability and exact technical work.

## Decision

Explain the relevant meaning and required action in the surrounding human-facing prose. Use descriptive links for optional detail, and secondary IDs where precise identification helps. Changelogs describe user-visible behavior and upgrade actions without internal rule IDs. Keep literal identities in metadata, tags, command syntax and structured data. Preserve pinned history.

CLI presentation names known local references and explicitly states uncertainty for unknown, unnamed, ambiguous or foreign references. A source-rule ID is not resolved against an adopter's same-spelled identity. Stored titles, raw findings, proof rules and exit codes are unchanged. This is a writing and presentation contract; the deterministic judge does not claim to measure comprehension.

The canonical shipped carrier is `internal/skills/source/shared/readable-references.md.tmpl`, routed by all six standalone skills. The scaffold hub, local acceptance and PR templates, source contributor guidance, guide and human CLI presentation carry the same method. The [design overview](../design/README.md) describes its separation from raw corpus state.
