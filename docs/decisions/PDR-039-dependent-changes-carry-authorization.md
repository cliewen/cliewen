---
id: PDR-039
type: decision
status: inferred
links: [P-013, AN-013, PDR-007, PDR-017, C-012, C-013, CAP-006]
title: An authorized dependent change records its base and what its merge would bind
author: agent
accepted-by: []
---

# PDR-039 — An authorized dependent change records its base and what its merge would bind

## Context and problem statement

PDR-007 permits work that genuinely depends on an unmerged change only after an explicit human decision. AN-013 shows why permission in conversation is insufficient: a dependent branch can absorb unaccepted meaning, and neither its corpus nor a green deterministic validation result says which base was relied on, who authorized that reliance, or what the dependent merge would bind.

The record must remain local and reviewable without mistaking a forge branch or merge status for truth. It must also reach the human who decides whether the proposed merge is acceptable.

## Decision outcome

**An authorized dependent change keeps its answered blocking question in its committed change workspace until digest.** The answer names the unmerged base change and branch or commit, the human authorization, and the specific unaccepted meaning that merging the dependent change would bind. It is not replaced by a private summary or a forge pointer.

**The ready pull request repeats that dependency in its acceptance brief's binding line.** The brief names the same base and says what accepting the dependent change would bind before that base is independently accepted. This is disclosure for the human merge judgment; it neither makes the base accepted nor authorizes an agent to merge either change.

The workspace record is deliberately transient. Once the base has been accepted and the dependent branch incorporates accepted `main`, the dependency no longer exists. The workspace is then digested as every full change is, while its committed history remains the repository record of the authorization and its scope.

No validator infers a dependency from Git graph or forge state. `clue validate` remains deterministic and state-based, and branch protection remains enforcement of admission rather than acceptance evidence. The rule is human-enforced because whether work semantically depends on unaccepted meaning is not derivable from files alone.

## Rejected: a permanent dependency registry

A registry would retain a relation after the base is accepted and the dependent branch has incorporated it, creating a second, stale representation of an event Git history already preserves. The need is visibility during the exceptional interval, not a new corpus artifact type.

## Rejected: infer the dependency from Git or a pull request

Git ancestry and a forge base can show a graph relation but cannot state the human authorization or the meaning at risk. Querying either would also make a deterministic local judgment depend on transient external state.

## Carrier

PDR-007's stacking clause is amended by this record. The shared review-boundary source generates the requirement into every lifecycle skill; C-012, the public change-loop guide, the repository and scaffolded pull-request templates, and CAP-006's handoff contract carry the same disclosure boundary.
