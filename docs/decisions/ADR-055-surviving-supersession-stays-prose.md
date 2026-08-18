---
id: ADR-055
type: decision
status: inferred
links: [P-014, PDR-038, ADR-034, AN-022, PDR-034, PDR-029, PDR-026, ADR-019, PDR-007, PDR-016, PDR-031, ADR-046]
title: A live superseded decision stays prose a second time; the annotation gets one settled shape
author: agent
accepted-by: []
---

# ADR-055 — Surviving supersession stays prose

## Context and problem statement

ADR-034 defines `supersedes:` for deleted artifacts, but widening it to live decisions would add a per-change schema obligation without giving `checkSupersedes` an observable way to tell a stale live relationship from a current one. The corpus instead needs one readable convention for the nine live supersession or amendment rows.

## Decision outcome

**Live supersession and amendment remain prose; `supersedes:` keeps ADR-034's deletion meaning.** An older index row says `superseded by <ID>` or `amended by <ID>`, and the newer row says `supersedes <ID>` or `amends <ID>`; a full ADR-046 description may state the same relationship. The three nonconforming rows in ADR-019, PDR-007, and PDR-016 are corrected, while all stated facts remain.

No reverse index is added: a live edge's truth is judgment rather than an observable validator fact, and a forward field would not answer the reverse question without more machinery. This decision names and repairs the corpus convention without widening the schema.

**Carrier:** this record, the corrected index rows, and the decisions-index preamble; no skill, template, or CLI rule is added.
