---
id: C-011
type: constraint
status: active
links: [PDR-006]
title: Decision records are routed by type — ADR for architecture, PDR for project/process, log row for the cheap
source: docs/decisions/PDR-006-decision-records-are-typed.md
enforcement: human
---

# C-011 — Decision records are routed by type

Two questions route every decision: cheap and local to reverse → a decision-log row; otherwise architecture (software or corpus format) → ADR, how the project works → PDR. An ADR about process, or a PDR about architecture, is misfiled and gets renamed into the right series. Decisions adopting a well-established practice cite it by name and record only the local why and deviations.

**Residual:** the whole routing judgment, permanently. The two questions are what a decision is *about* and what undoing it would cost — neither is a property of the file. A misfiled record is well-formed, links correctly, and validates green.

The cost is a corpus whose decisions are hard to find by the question they answer, and a log that fills up with decisions expensive enough to deserve their own record. It is caught by review and by the periodic act of reading the decisions index, which is one reason that index is curated rather than merely generated.
