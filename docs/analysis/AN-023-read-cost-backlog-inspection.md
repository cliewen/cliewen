---
id: AN-023
type: analysis
status: active
provenance: inferred
reversal-cost: low
links: [P-015, M-072, ADR-057, ADR-058, C-022]
title: The read-cost backlog is inspected without treating its count as a target
---

# AN-023 — Read-cost backlog inspection

## Scope and method

This record is CH-147's inspection of the populations named by `clue validate --read-cost` over the corpus at `ef93cfa`, the merge base this change branched from — one artifact holding multiple rendered documents, and thirty-three identities whose default slice exceeds the focused-read budget. The same command re-run on the change branch reports the same two populations, because an accepted row keeps its place in the report ([ADR-058](../decisions/ADR-058-read-cost-is-a-backlog-not-a-target.md)). It records the judgment that command deliberately cannot make: whether a rendered-document split has distinct consumers, and whether an over-budget identity has a genuinely redundant outgoing relationship. It does not turn either population into a pass/fail target or a registry; later corpus changes rerun the command and receive a new population.

## Multi-document artifact

| Artifact | Outcome | Reason |
| --- | --- | --- |
| AN-022 | Accepted | Anyone reading its Pattern C determination needs the non-prose scoring that precedes it to interpret the verdict; they are one investigation with shared evidence, so splitting them would fragment that reader's conclusion. |

## Over-budget slices

Every inspected path remains accepted. The listed relationships identify the decision's scope, its evidence, the contract it implements, or the durable carrier that a reader needs to follow; no entry was redundant for that artifact's reader.

| Identity | Outcome | Reason |
| --- | --- | --- |
| ADR-038 | Accepted | Its campaign and milestone, input analyses, capability owners, installation baseline, and merge boundary together define the reusable-workflow decision. |
| ADR-039 | Accepted | Its campaign and milestone, earlier migration and workflow decisions, capability owners, and governing constraints define the safe-upgrade contract. |
| ADR-040 | Accepted | Its campaign, originating analysis, evidence and migration decisions, judge boundary, and capability owners define qualified foreign references. |
| ADR-041 | Accepted | Its capability owners, prior index work, relevant constraints, and follow-on index decision define the generated-index rule. |
| ADR-042 | Accepted | Its plan, capability owners, release and install precedents, network-reference boundary, and governing constraints define the release-only check. |
| ADR-044 | Accepted | Its plan, judge owner, related state and reference decisions, and constraints identify the complete state-only validation boundary. |
| ADR-046 | Accepted | Its index and validator owners, prior index decisions, constraints, and plan anchor distinguish the curated index-row contract. |
| ADR-049 | Accepted | Its migration, extraction, ledger, rehearsal, and core-constraint links define the two-manifest parity architecture. |
| ADR-055 | Accepted | Its successor campaign, rejected predecessor, lifecycle rule, evidence analysis, review and simplification decisions, and corrected record set define the settled prose shape. |
| AN-008 | Accepted | Its plan, core architecture and goal, critique decisions, and governing constraints are the evidence boundary for the four-pattern finding. |
| AN-011 | Accepted | Its goal, campaign, criterion, capability owners, architecture, decisions, and constraints locate the contradiction it reports and the work that resolved it. |
| AN-012 | Accepted | Its campaign and milestone, adopter evidence, contract and migration decisions, capability owners, constraint, and successor-plan links preserve the measured upgrade finding. |
| AN-013 | Accepted | Its two campaigns and milestone, goal, evidence and criterion-identity decisions, versioning and installation contracts, CI-wall and external-reference decisions, authorization and dependent-change rules, routing, review, and core constraints, and three related analyses preserve all three distributed-work gaps. |
| AN-014 | Accepted | Its plan, two source analyses, evidence and identity decisions, rehearsal rule, and capability owners retain the confidential-assessment finding's boundary. |
| AN-015 | Accepted | Its plan, rehearsal and merge rules, criterion decision, source analyses, and extraction owner retain the retrospective rehearsal's proof boundary. |
| AN-017 | Accepted | Its campaigns, closure rule, migration decisions, delivery rules, fixture evidence, and extraction owner define the re-derived migration gate register. |
| AN-018 | Accepted | Its campaign, simplification and core rules, skill decision, identity rule, constraints, and goal are the statement register's scoring basis. |
| AN-020 | Accepted | Its campaign and milestone, source register, review and handoff decisions, and simplification rule explain the routing hub's carrier placement. |
| AN-021 | Accepted | Its campaign and milestone, scoring and carrier rules, source register, constraints, and quality bar are the public-surface register's evidence base. |
| AN-022 | Accepted | Its campaign and milestone, core and simplification rules, lifecycle and ledger decisions, source analyses, and core constraint are the scored surface's complete basis. |
| C-012 | Accepted | Its review, merge, handoff, acceptance, authorization, durability, and decision-log links are the sources of the protected human-merge constraint. |
| CAP-003 | Accepted | Its goal, extraction, provenance, and criterion-identity decisions, the manifest-parity, imported-change, carrier-inventory, and deferred-disposition records, and the derived-report rule define the brownfield-extraction capability. |
| CAP-004-design | Accepted | Its capability, release and skill contracts, installation and CI workflow decisions, external-reference boundary, and release check define the shipping design. |
| PDR-017 | Accepted | Its goal, originating critique, campaign, review, core, authorization, and constraints define the acceptance-brief boundary. |
| PDR-019 | Accepted | Its goal, campaign, contradiction analysis, core boundary, method decisions, and constraints define the live-carrier synchronization rule. |
| PDR-020 | Accepted | Its plan, extraction owner, source analyses, extraction decision, and carrier/core rules define report-only rehearsal before mutation. |
| PDR-021 | Accepted | Its goal, campaign and milestone, adopter evidence, review rules, architecture, and merge constraint define the supported merge-commit path. |
| PDR-023 | Accepted | Its plan, capability owners, versioning and notice decisions, carrier rule, and core constraint define the tool-notice and hub-instruction split. |
| PDR-029 | Accepted | Its campaign, core and carrier rules, lifecycle and analysis evidence, goal, constraints, and index rule define both simplification tests. |
| PDR-032 | Accepted | Its inbox, acceptance, review, task and goal rules, campaign, and source register define the mandatory suggestion-triage path. |
| PDR-033 | Accepted | Its acceptance, review, carrier, durability, task, merge, goal, campaign, and source-register links define the human planning/implementation boundary. |
| PDR-035 | Accepted | Its campaign, review and carrier rules, simplification rule, constraints, quality budget, and generated-skill decision define the review-loop contract. |
| PDR-040 | Accepted | Its goal, review and handoff rules, planning boundary, merge constraint, and collaboration capability define the durability and ready-state contract. |

## What was tried and rejected

**Trimming links until the report went quiet.** CH-146 did this first: it edited the reported artifacts, brought the over-budget population down, and then reverted every one of those edits at `343227c` before merging. The trim was cheap to produce and impossible to audit — the entries it removed were selected by what was easy to identify as duplicative, not by what a reader of each artifact needed, and a removed entry is unreachable at any `--depth` because `corpus.Context` walks links outward only. That reversal is what produced [ADR-058](../decisions/ADR-058-read-cost-is-a-backlog-not-a-target.md), and this record is the inspection the reverted trim skipped. The two further options ADR-058 weighed and refused — making the budget a failing check, and ranking links automatically to keep a budget's worth — are recorded there rather than repeated here.

**Recording the outcomes as a registry joined to the report.** This would remove the residual ADR-058 names, where an accepted row is indistinguishable from an uninspected one. It is refused for the reason [ADR-057](../decisions/ADR-057-read-cost-measurements.md) refused it: the report derives every row from current corpus state and writes nothing, so a registry would be a second population needing its own maintenance and going stale silently. The reason for each acceptance lives in this record, reachable from the plan whose milestone worked the backlog and from the decision that governs it.
