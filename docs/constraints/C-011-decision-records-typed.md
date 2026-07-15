---
id: C-011
type: constraint
status: active
links: [PDR-006]
title: Decision records are routed by type — ADR for architecture, PDR for project/process, log row for the cheap
source: docs/decisions/PDR-006-decision-records-are-typed.md
enforcement: agent
---

# C-011 — Decision records are routed by type

Two questions route every decision: cheap and local to reverse → a decision-log row; otherwise architecture (software or corpus format) → ADR, how the project works → PDR. An ADR about process, or a PDR about architecture, is misfiled and gets renamed into the right series. Decisions adopting a well-established practice cite it by name and record only the local why and deviations.

**Promotion trigger:** the routing judgment is semantic (what a decision is *about*) and stays agent-enforced; no machine check is foreseeable.
