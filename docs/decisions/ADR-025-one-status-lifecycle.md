---
id: ADR-025
type: decision
status: inferred
links: [P-005, ADR-010, ADR-002, C-008, PDR-003]
title: One default status lifecycle — draft → active → retired — plus named exceptions
author: agent
accepted-by: []
---

# ADR-025 — One default status lifecycle

## Context and problem statement

The validator carried fourteen per-type status vocabularies. Eleven of them were the same three-state lifecycle `draft → active → retired` or a subset of it: capability, criteria, design, constraint (which merely omitted `draft`), quality, and — spelled with a different terminal word — architecture (`draft → verified`) and analysis (`verified` only). Every reader and every adopter had to learn a table of near-duplicates to find the four rows that actually differ. A field's vocabulary should be as small as its meaning requires; this one was far larger. Which types genuinely need their own lifecycle, and what is the default for the rest?

## Decision outcome

**One default lifecycle, `draft → active → retired`, applies to every artifact type. A type gets its own vocabulary only for a stated semantic reason.**

The exceptions, and why each earns one:

- **goal** — `proposed → accepted → retired`: proposed goals *are* the inbox ([ADR-002](ADR-002-inbox-is-proposed-goals.md)); `proposed` is a distinct meaning, not a draft.
- **plan** — `draft → active → completed`: `completed` is immutable ([C-008](../constraints/C-008-completed-plans-immutable.md)), a terminal state `retired` does not capture.
- **decision** — `inferred → verified`: a decision expresses its provenance in `status` ([ADR-010](ADR-010-provenance-field.md)), where `verified` means a human accepted it.
- **log** — `active`: the decision log is one standing register; its rows, not the file, carry lifecycle ([PDR-003](PDR-003-decision-log.md)).
- **change, tasks** — `open`; **open-questions** — `open → resolved`: transient workspace artifacts that never reach `main`.

Every other type — capability, criteria, design, constraint, quality, architecture, analysis, and any adopter-defined type ([ADR-026](ADR-026-adopter-types-default-lifecycle.md)) — uses the default.

**Architecture and analysis move off `verified`-in-status to the default.** Their `verified` was never the decision sense of the word: it marked the single state those records were allowed to hold. Provenance — whether an agent-extracted record has been human-checked — already has its own field for every non-decision type (`provenance: inferred | verified`, [ADR-010](ADR-010-provenance-field.md)), so a second verified-in-status channel was redundant and easy to misread. An architecture or analysis record is `active` while it stands and `retired` when superseded; its immutability as a historical record is a documentation convention, not a status value.

**Constraint gains `draft`.** Folding constraint into the default widens its vocabulary from `{active, retired}` to `{draft, active, retired}`. This is a harmless widening — no constraint is required to pass through `draft` — that buys one fewer special case.

**Carrier:** the `defaultLifecycle` slice and `statusVocabExceptions` map in `internal/corpus/rules.go`, mirrored by the status tables in `docs/README.md` and the `clue init` template.

### Rejected: keep the per-type table

The table made the four genuine exceptions invisible among ten look-alikes and forced every adopter to read all fourteen rows to learn what amounts to "draft, active, retired, with four footnotes."

### Rejected: move decisions onto the default too

`inferred → verified` is load-bearing: it is how [ADR-010](ADR-010-provenance-field.md) and "merge binds, approval signs" ([PDR-004](PDR-004-merge-binds-approval-signs.md)) express whether a human has accepted a decision. Collapsing it into `draft → active` would erase the acceptance signal the whole decision workflow turns on. Decisions keep their vocabulary because their meaning requires it — which is exactly the test this decision applies.
