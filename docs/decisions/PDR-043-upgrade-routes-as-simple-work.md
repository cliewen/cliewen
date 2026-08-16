---
id: PDR-043
type: decision
status: inferred
links: [CAP-004, ADR-039, ADR-043, ADR-060, PDR-042]
title: An upgrade routes as simple work unless the release makes the adopter decide
author: agent
accepted-by: []
---

# PDR-043 — An upgrade routes as simple work unless the release makes the adopter decide

## Context and problem statement

`clue-upgrade` sends the agent to the shared change-routing reference before it recommends a route, and nothing in the skill answers the question that reference raises. The agent therefore answers from the tier text alone, which routes a change in "capability, policy, plan promise, decision, or methodology meaning" to the full loop and resolves uncertainty the same way. An upgrade rewrites the six managed skills, the thin caller, and sometimes corpus shape, so full is the honest reading of a carrier set that never states otherwise — and it is the wrong answer. What route does moving a repository onto a published release take?

## Decision outcome

**An upgrade is simple work.** It moves a repository onto a release whose contract changes were made, argued, and accepted upstream before publication. The adopting repository's own accepted contract — its capabilities, criteria, decisions, plans, and constraints — is unchanged by the act of adopting the release, so the upgrade carries no CH identity, workspace, plan declaration, digest, acceptance brief, or mandatory agentic review. It runs the checks relevant to the surfaces it changed and follows the repository's integration rules. In PDR-042's own terms it is maintenance, and the routing rule already classified it; only the carriers were silent.

**The escalation is semantic, not mechanical.** Applying a release sometimes requires the repository to decide something of its own: an obligation it must choose how to meet, a reconciliation that changes what its own criteria or its own CI wall promise, or a migration finding whose resolution is a judgment about this repository rather than a transcription of the release. That decision is full work, recommended and taken on its own terms. It does not make the surrounding upgrade full; the mechanical move stays simple and the decision is routed as what it is.

The number of files a migration rewrites, the size of its diff, and the presence of `/docs` or skill paths in it never make the route full — PDR-042 already settled that paths and counts warn and do not decide.

Nothing else about the upgrade changes. The human is still asked whether to upgrade now or later before any repository write, and the pull request is still accepted by a human under the repository's merge boundary; ADR-043's two visible boundaries are exactly what a simple route leaves intact.

**Carrier inventory:** the canonical `clue-upgrade` skill source and both generated skill trees; CAP-004's criteria (AC-081 retired, AC-142) and design; the skills architecture hand-off table; `guide/operations.md`; and the tests pinning the generated upgrade skill. PDR-042 and ADR-043 gain a link from this record rather than new text: neither said the wrong thing, and this decision states what both left unstated.

### Rejected: route every upgrade through the full loop

The full loop exists so a human accepts meaning the change introduces. An upgrade introduces no meaning the adopter can accept or decline artifact by artifact: the release's contract was accepted upstream, and the adopter's decision — whether to be on that release at all — is the one `clue-upgrade` already puts to the human before the first repository write. A CH identity, digest, and acceptance brief would ask the adopter to re-accept authorship it does not hold, at a cost measured against a preview, an apply, and a binary swap. Observed twice in the `model2diagram` adopter, this reading turned a mechanical version move into a full change with a workspace and a plan question it could not honestly answer.

### Rejected: decide the route from what the migration rewrites

Reading the route off the migration's file list makes the biggest release the fullest change and a corpus-shape migration full by construction, which is the diff-size heuristic PDR-042 rejected wearing a different hat. It also gets the real case backwards: a release that rewrites one line of a constraint an adopter must now satisfy deserves more of a decision than one that rewrites every skill byte.

### Rejected: leave the route to the shared tier text

That is the status quo, and it produces full every time. The tier text is generic by design; the skill knows something generic text cannot — that the contract change already happened somewhere else — and that knowledge belongs where the workflow is stated.
