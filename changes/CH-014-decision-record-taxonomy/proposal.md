---
id: CH-014
type: change
status: open
links: [P-002]
title: Decision-record taxonomy — ADRs for architecture, PDRs for project/process
---

# CH-014 — Decision-record taxonomy

ADR means *Architectural* Decision Record, and five of ours are not about architecture — they are rules about how this project works (the light change tier, the decision log, how decisions get accepted, the foreign-soil trials). Using the well-known ADR name for them makes the term lie to anyone who knows it. The industry already tested the alternative: the MADR template project renamed itself from "Architectural" to "Any" decision records in version 3.0 and renamed **back** in 4.0 — the stretched meaning did not survive contact with its own community.

## What

Decision records split into three explicit tiers:

- **ADR** — architectural decisions only (the software and the corpus format), Nygard's strict sense. ADR-001…ADR-013 and ADR-017 qualify and keep their IDs.
- **PDR** (Project/Process Decision Record) — expensive-to-reverse decisions about how the project works. Same MADR template, same `inferred`/`verified` provenance, own ID series. ADR-014, 015, 016, 018, 019 are reclassified as PDR-001…005 (retroactive cleanup: the corpus must not lie about itself).
- **Decision log** — unchanged: one dated row for the cheap-and-local-to-reverse.

A new record, **PDR-006**, carries the taxonomy decision itself and one companion rule: a decision that adopts a well-established practice cites it by name instead of re-deriving it — the record captures only the local why and any deviation. PDR-003 (né ADR-016) is superseded by PDR-006. A new register entry, **C-011**, makes the routing rule agent-enforced.

**Plan revision to [P-002](../../docs/plans/P-002-leaves-home.md)**: the mutation rule said plan revisions are "backed by an ADR" — that wording is what forced the foreign-soil decision to be misfiled as an ADR in the first place; it becomes "backed by a decision record (PDR for direction and process, ADR if architectural)".

Pending clerical signing performed en route: the foreign-soil decision (now PDR-005) was explicitly approved by Flemming N. Larsen on 2026-07-15 (stated in session, with PR #12 approval) — first approval, so it turns `verified`.

## Files

- `docs/decisions/` — five renames (git mv, IDs/titles/cross-links updated, bodies otherwise untouched), new PDR-006 born `inferred`, README rewritten around the three tiers
- `docs/constraints/C-011-*.md` + constraints README index
- Cross-references: AGENTS.md, docs/README.md status table, P-002 (M-007 link + mutation rule), skills (clue-delta, clue-verify, clue-plan, clue-analysis — wording only)
- `CHANGELOG.md` — user-visible methodology change
- No Go code changes: the linter keys on `type: decision`, not on ID prefixes
