---
id: P-021
type: plan
status: active
links: [G-012, VIS-001, UC-001, CAP-009, ADR-065, ADR-066, ADR-067, PDR-054, PDR-055, C-023, G-001, CAP-002, CAP-007]
title: Cliewen states what the product is for
---

# P-021 — Cliewen states what the product is for

Cliewen traces a goal to its evidence and says nothing about why the goals are the goals. A reader arriving at a mature corpus can follow any criterion back to the capability and the goal that asked for it, and still cannot answer what the product is, whom it serves, or what is deliberately outside it. The same gap costs an agent more than a human: orienting on a change, it has capability-local truth and no statement of direction to weigh a judgement against, so it either asks the human every time or infers direction from implementation structure, which is exactly the inference [ADR-067](../decisions/ADR-067-a-corpus-without-a-vision-stays-valid.md) says code cannot support.

The second gap is narrower and older. A capability is a unit of what the system can do, and an acceptance criterion is a unit of proof. Neither carries an actor's path *across* capabilities — the ordering, the alternative, the failure recovery — so a set of locally correct criteria can describe a journey nobody would want. Cliewen has lived without that artifact deliberately, and this campaign keeps it optional for the same reason: most behaviour genuinely is one capability wide, and an artifact created to satisfy a process is worse than no artifact.

This campaign adds both concepts at the smallest size that carries their meaning, and its hardest constraint is what it must *not* do. It must not make a 0.22.0 repository invalid for lacking either one; it must not let migration write a vision, because nothing in a repository proves why a product exists; and it must not let the judge claim that a vision is correct or a use case complete. Those are human judgements, and a validator that appeared to check them would be the most expensive thing this campaign could ship.

Milestone numbering continues corpus-global numbering from P-020's M-086. The identity ledger still does not cover milestones ([G-006](../goals/G-006-milestone-ids-in-the-ledger.md)), so these are assigned by reading the corpus maximum.

## Milestones

| ID | Milestone (exit criterion) | Status | Evidence |
|---|---|---|---|
| M-087 | **A corpus can carry one vision and any number of optional use cases, and the judge checks their form without claiming their meaning.** `clue validate` accepts a corpus with no vision and no use cases, accepts one active vision at its fixed address, and rejects a second vision, a vision elsewhere, an unreplaced bootstrap, a use case outside its folder or misnamed, a use case naming no goal or no capability, and a use case missing a required structural section. No rule requires a use case for any goal, capability, or criterion. Focused positive and negative evidence covers each. | `todo` | |
| M-088 | **Intent is reachable from the artifact a task starts at, without reverse-expanding the corpus.** A goal's outgoing link reaches the vision, a use case's outgoing links reach its goal and the capabilities it crosses, and `clue context` names — by identity and title only — the use cases whose links name the root. The bounded default slice is unchanged and no request loads the intent graph as a whole. Focused positive and negative evidence covers naming, ordering, and the absence of content expansion. | `todo` | |
| M-089 | **Agents elicit intent on greenfield, infer it on brownfield, and never present either as confirmed.** The shipped skills carry a proportionate interview, an evidence-first inference that cites repository sources and records conflicts rather than resolving them, the documented escape hatch for drafting from little, and the test for when a use case is worth creating and when it is not. Drafted content is `draft` with `provenance: inferred`; promotion is a human act. | `todo` | |
| M-090 | **A 0.22.0 repository migrates without a fabricated vision, and one without a vision stays valid while meaning-changing work discloses it.** `clue migrate` adds the optional-use-case folder structure and reports a missing vision as a notice that blocks nothing, writing no vision content. `clue validate --intent` states what exists without a coverage figure. The acceptance brief carries the vision the change proceeds under as a required line the shipped CI wall already enforces. | `todo` | |

## Mutation rules

Status and evidence fields in the milestone table may mutate in an implementing change's merge digest. Everything else changes only through a declared plan revision backed by a correctly typed decision record. Plan adjustments are decisions.
