---
id: CH-155
type: change
status: active
links: [CAP-004, ADR-043, PDR-042]
title: An upgrade routes as simple work unless the release makes the adopter decide
---

# CH-155 — An upgrade routes as simple work unless the release makes the adopter decide

## Proposal

This change is plan-less: it corrects a routing gap observed in an adopter session rather than serving a milestone of the active plan P-015, whose remaining milestone M-074 concerns index-row descriptions.

`clue-upgrade` sends the agent to the shared change-routing reference before recommending a route — `internal/skills/generate.go:102` states the condition as "If the human chooses to upgrade now and before recommending its route" — and then never answers the question it raised. The agent answers it from the tier text alone, which routes anything changing "capability, policy, plan promise, decision, or methodology meaning" to full and makes uncertainty resolve the same way. An upgrade rewrites the six managed skills, the thin caller, and sometimes corpus shape, so the honest reading of the current carriers is full. `docs/architecture/skills.md:23` says the same thing outright, handing the affirmative choice to `clue-delta` — the full-loop skill by name.

That reading was observed twice in the `model2diagram` adopter: the 0.6.0 → 0.16.0 upgrade landed as a full change with a CH identity, and a 0.17.0 upgrade proposed the same shape. The cost is a CH identity, workspace, plan declaration, digest, acceptance brief, and mandatory agentic review for what is mechanically a preview, an apply, and a binary swap.

An upgrade does not change the adopting repository's accepted contract. The contract changes it carries were made and accepted upstream in Cliewen's own full-tier work and published in a release; the adopter is moving vendored carriers onto a version whose meaning someone else already accepted. That is maintenance in the routing rule's own terms. The exception is genuine and semantic rather than mechanical: when applying a release requires this repository to decide something of its own — an obligation it must choose how to meet, or a competing-wall reconciliation that changes what its own criteria promise — that decision is full work on its own terms.

The change will record the routing decision as a PDR extending PDR-042, state the route and its escalation condition in the canonical `clue-upgrade` source so every generated and scaffolded copy carries it, retire AC-081 and mint AC-142 for the revised skill contract with focused positive and negative Unit evidence, and move the remaining live carriers that state the upgrade's shape.

## Scope boundary

This change decides how an upgrade routes. It does not change what `clue migrate` does, what `clue latest` reports, the human authorization step, or the merge boundary — an upgrade remains a reviewed change no agent merges. It does not add a release route for adopters, which PDR-042 already rejected.
