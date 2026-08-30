---
id: P-018
type: plan
status: active
links: [PDR-051, CAP-001, CAP-002, CAP-006]
title: Every Cliewen corpus keeps a living system overview
---

# P-018 — Every Cliewen corpus keeps a living system overview

Cliewen's durable corpus must let a reader understand the system as a whole without making them reconstruct it from capability folders or a stream of decisions. This campaign establishes concise architecture and design overviews, makes their absence visible to the judge, and gives every lifecycle a feedback loop that keeps the relevant documentation current without preserving a second change history.

Milestone numbering continues corpus-global numbering from P-016's M-079.

## Milestones

| ID | Milestone (exit criterion) | Status | Evidence |
|---|---|---|---|
| M-080 | **Every Cliewen corpus has concise, maintained architecture and design overviews.** `docs/architecture/README.md` covers system structure and `docs/design/README.md` covers cross-cutting behaviour; `clue init` creates explicit bootstrap files, and `clue validate` rejects a missing or unactivated overview. Migration and upgrade safely surface the requirement while lifecycle skills discover, draft, maintain, and—with consent—relocate existing overview content without duplication. Capability-local design remains local, relevant ADRs and IDRs link to the appropriate overview, and focused tests prove the validator, migration, and generated guidance. | `doing` | CH-162 proposal |

## Mutation rules

Status and evidence fields in the milestone table may mutate in an implementing change's merge digest. Everything else changes only through a declared plan revision backed by a correctly typed decision record.
