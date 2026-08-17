---
id: CH-159-log-classification
type: analysis
status: active
links: [CH-159, PDR-046]
title: Reviewed disposition of every legacy decision-log row
---

# Legacy decision-log classification

This transient audit accounts for every row in `docs/decisions/log.md`, newest first. “Record” names the subject-routed durable destination; “narrative” means the row records chronology, a completed plan transition, a carrier edit, or a withdrawn rule rather than a choice that constrains future work. Git history and the durable artifact named in the disposition retain the relevant state without converting that narrative into a decision record.

| Row | Date | Subject | Disposition |
|---|---|---|---|
| 01 | 2026-08-12 | Constraint index badges use enforcement | Record: IDR-001 |
| 02 | 2026-08-10 | P-013 closure and P-014 opening | Narrative: completed plan transition; P-013 and P-014 retain the state |
| 03 | 2026-08-09 | Orient after a reported merge | Record: PDR-047 |
| 04 | 2026-08-08 | Withdraw immediate task ticking | Narrative: withdrawn rule; C-003 states the surviving skipped-reason obligation |
| 05 | 2026-08-08 | Correct guide review-loop and small-delta prose | Narrative: carrier correction under PDR-029 |
| 06 | 2026-08-08 | Move poor-fit conditions into architecture | Record already present: PDR-031 governs architecture as a trace; ARCH-003 states the conditions |
| 07 | 2026-08-08 | Revise M-063 around named deferrals | Narrative: completed-plan history retained by P-013 |
| 08 | 2026-08-08 | Keep detailed handoff in generated boundary | Record already present: PDR-029's single-statement carrier rule and PDR-034's narrow-reading rule |
| 09 | 2026-08-08 | Withdraw end-of-plan decision distillation | Narrative: withdrawn rule |
| 10 | 2026-08-08 | Withdraw uncheckable small-delta advice | Narrative: withdrawn rule |
| 11 | 2026-08-08 | Milestones need verifiable exit criteria | Record: PDR-048 |
| 12 | 2026-08-08 | Allocate CH identities through the ledger | Record already present: ADR-048 |
| 13 | 2026-08-08 | Do not restate computed figures | Record: PDR-049 |
| 14 | 2026-08-08 | References say what their records hold | Record: PDR-049 |
| 15 | 2026-08-08 | Plain changes cannot depend on unmerged work | Record already present: PDR-011 |
| 16 | 2026-08-08 | Guard and report the managed-skill roster | Record: IDR-002 |
| 17 | 2026-08-04 | Milestone status vocabulary | Record: PDR-048 |
| 18 | 2026-08-04 | Cache release-list nonanswers | Record: IDR-003 |
| 19 | 2026-08-02 | Bulk decision provenance promotion | Narrative: approval history remains in each record's `accepted-by` field and Git |
| 20 | 2026-08-02 | Open P-010 | Narrative: completed plan history retained by P-010 |
| 21 | 2026-08-02 | Close a plan in the final milestone digest | Record: PDR-048 |
| 22 | 2026-08-01 | Close P-009 | Narrative: completed plan state retained by P-009 |
| 23 | 2026-08-01 | Widen P-009 M-044 | Narrative: completed-plan revision history retained by P-009 |
| 24 | 2026-08-01 | Remove dependent-change authorization record | Narrative: reversed by the later durable PDR-039 authorization contract |
| 25 | 2026-08-01 | One-off migration rehearsal direction | Narrative: one-off scope and outcome retained by AN-015 |
| 26 | 2026-07-31 | Refresh local clue before this repository's release | Record: IDR-004 |
| 27 | 2026-07-30 | Extend P-009 with M-043 and M-044 | Narrative: completed-plan revision history retained by P-009 |
| 28 | 2026-07-30 | Open P-009 | Narrative: completed plan history retained by P-009 |
| 29 | 2026-07-30 | Close P-008 and designate P-009 | Narrative: completed plan transition retained by P-008 and P-009 |
| 30 | 2026-07-30 | Decline distributed-work interfaces | Records: ADR-061 for external coordination; PDR-039 for authorized local dependencies |
| 31 | 2026-07-29 | Decline adopter configuration | Record already present: ADR-013 |
| 32 | 2026-07-29 | State analysis environment and population boundaries | Record: PDR-050 |
| 33 | 2026-07-28 | Open P-008 | Narrative: completed plan history retained by P-008 |
| 34 | 2026-07-27 | Keep goreleaser changelog generation disabled by omission | Record: IDR-004 |
| 35 | 2026-07-26 | Verify tag-to-skill consistency with the shipped judge | Record: IDR-004 |
| 36 | 2026-07-25 | Open P-007 | Narrative: completed plan history retained by P-007 |
| 37 | 2026-07-25 | Close P-006 | Narrative: completed plan state retained by P-006 |
| 38 | 2026-07-25 | Convert proposed and rejected MADR records without inventing acceptance | Record: IDR-005 |
| 39 | 2026-07-25 | Resolve duplicate MADR numbers deterministically | Record: IDR-005 |
| 40 | 2026-07-25 | Preserve deprecated and superseded MADR meaning | Record: IDR-005 |
| 41 | 2026-07-25 | Extraction entry points, phasing, and minted source IDs | Record: IDR-005 for the mapping choices; current extraction criteria retain the broader contract |
| 42 | 2026-07-25 | Init skips linked directories | Record: IDR-006 |
| 43 | 2026-07-24 | Open P-006 | Narrative: completed plan history retained by P-006 |
| 44 | 2026-07-24 | Close P-005 | Narrative: completed plan state retained by P-005 |
| 45 | 2026-07-24 | Present change routing as three rules and escalation guards | Record already present: PDR-002; current carrier shape is implementation narrative |
| 46 | 2026-07-23 | Retire quality IDs and fold quality into constraints | Record already present: ADR-027 |
| 47 | 2026-07-23 | Close P-004 | Narrative: completed plan state retained by P-004 |
| 48 | 2026-07-22 | Scope the first-try campaign and defer distribution claims | Record already present: PDR-014; P-004 retains campaign history |
| 49 | 2026-07-21 | Human-controlled merge does not require duplicate review | Record already present: PDR-007 |
| 50 | 2026-07-20 | Ready state binds a clean published commit | Record already present: PDR-040 |
| 51 | 2026-07-19 | Reject hidden or duplicate frontmatter | Record: IDR-007 |
| 52 | 2026-07-18 | Analysis evidence states revision, conditions, and confidence | Record: PDR-050 |
| 53 | 2026-07-18 | Open PRs ready rather than draft | Narrative: superseded by PDR-040's draft-from-publication and explicit-ready boundary |
| 54 | 2026-07-17 | Seed only non-skill constraints | Record already present: ADR-020 |
| 55 | 2026-07-17 | Split retired onboarding criterion | Record already present: ADR-027; C-015 holds the Human bar |
| 56 | 2026-07-17 | Place embedded templates under internal scaffold | Record already present: ADR-018 |
| 57 | 2026-07-17 | Orient after a reported merge | Record: PDR-047; duplicate of row 03's recovered trace |
| 58 | 2026-07-16 | Clarify reversal-cost routing | Narrative: superseded by PDR-046's subject routing |
| 59 | 2026-07-16 | Digest is not a task | Record: PDR-048 |
| 60 | 2026-07-12 | Aggregate Go coverage floor | Record: IDR-008 |
| 61 | 2026-07-12 | Parse frontmatter with yaml.v3 | Record: IDR-007 |
