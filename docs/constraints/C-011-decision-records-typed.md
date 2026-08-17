---
id: C-011
type: constraint
status: active
links: [PDR-046]
title: Future-shaping decisions route by subject into ADR, PDR, or IDR
source: docs/decisions/PDR-046-decisions-route-by-subject.md
enforcement: human
---

# C-011 — Decision records are routed by type

Only a future-shaping choice earns a decision record. Subject routes it to exactly one type: software or corpus architecture → ADR, project/process or methodology → PDR, implementation → IDR. Reversal cost does not route a record; routine facts, chronology, and implementation history stay in their natural carriers. A misfiled record is renamed into the right series. Decisions adopting a well-established practice cite it by name and record only the local why and deviations.

**Residual:** whether a choice is future-shaping and what subject it has are judgments about meaning, not file shape. The validator can reject unsupported identities and filenames after the author chooses a type; it cannot tell whether the choice should have become a record or whether ADR, PDR, or IDR is the honest subject.

The cost is either durable intent lost as narrative or a corpus whose decisions are hard to find by the question they answer. Review and the curated decisions index hold that residual.
