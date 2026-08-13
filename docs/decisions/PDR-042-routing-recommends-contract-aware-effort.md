---
id: PDR-042
type: decision
status: verified
links: [G-001, G-008, PDR-002, PDR-007, PDR-011, PDR-012, PDR-018, PDR-030, PDR-040, PDR-041, C-002, C-012, C-013, CAP-006]
title: Routing recommends effort from accepted-contract change while the user retains integration authority
author: agent
accepted-by: Flemming N. Larsen (2026-08-13, planning conversation)
---

# PDR-042 — Routing recommends effort from accepted-contract change

## Context and problem statement

The plain/light/full hierarchy makes path and runtime behavior stand in for semantic cost. A findings document under `/docs` is automatically non-plain even when it records observations without deciding anything, while every bug fix is full merely because correcting a defect changes runtime behavior. The hierarchy also turns Cliewen's supported pull-request workflow into language that appears to govern what repository owners may do, even though Git permissions and local policy belong to them. How should an agent recommend enough process without making Cliewen the authority over its adopter?

## Decision outcome

**Cliewen has two semantic recommendations: simple work stays outside the change loop, and work that changes the accepted contract is recommended for the full loop. The user chooses the route and the repository chooses how work is integrated.**

Before editing, the agent inspects the smallest relevant durable context and states `Recommended route: simple` or `Recommended route: full`, why, and what discovery would change the recommendation. It reassesses when work exposes a semantic expansion and once more against the complete diff before integration. File paths, changed-line counts, and file counts are warning signals only; none proves what a change means.

Simple work leaves the accepted contract intact. It includes an observational analysis with a named consumer, a defect correction that restores an unchanged acceptance criterion, focused regression evidence for an unchanged criterion, a configuration adjustment inside the supported contract, a refactor, maintenance, and editorial correction. It carries no CH identity, workspace, plan declaration, digest, acceptance brief, or mandatory agentic review, and it runs the checks relevant to the surfaces it changes.

Full work changes the accepted contract: it adds, materially revises, or retires an acceptance criterion; introduces behavior outside the accepted criteria; changes capability, policy, plan promise, decision, or methodology meaning; or makes or rejects a consequential decision from an analysis. Uncertainty makes the full route the recommendation because the agent cannot honestly assert that the contract is unchanged.

Discovery changes the recommendation, not the user's authority. When simple work grows into full work, the agent pauses, explains the newly discovered meaning, and recommends entering the full loop before continuing. If the user explicitly declines, the agent proceeds as simple, keeps implementation, tests, and durable documentation truthful, and records the one-integration authorization in the final authored commit using all three trailers:

```text
Cliewen-Route: simple
Cliewen-Recommendation: full
Cliewen-Override: user chose simple; <concise semantic or evidence risk>
```

The trailers are Git history, not a corpus decision: the override chooses process for one integration and does not become permanent product meaning. An incomplete trailer set claims no override. A forge that rewrites commits may not preserve the record; that limitation is reported rather than repaired with a forge-specific registry.

Route and integration authority are separate. A route never authorizes a push. An agent pushes directly to an integration branch only with explicit user authorization and only when repository permissions allow it; otherwise it follows the repository's requested integration workflow. A human acting independently may integrate however the repository permits. A repository may add stricter local rules, including Cliewen's own rule that coding agents always use a pull request and human merge.

A release is not a Cliewen route. An adopter defines or omits its own release process, and release work is classified by whether it changes that adopter's accepted contract. The `cliewen/cliewen` version cut remains a repository-local specialization of simple work with its exact administrative surface and focused checks; it is not emitted as adopter methodology.

This replaces PDR-002's light tier, PDR-011's narrow plain boundary, and PDR-018's rule that every product behavior correction is full. It amends PDR-007 and PDR-040 only at their scope: their draft-PR durability and human merge rules govern a full loop the user chose, while simple integration follows explicit user authority and repository policy. The protected full-loop acceptance boundary remains intact; this decision changes when Cliewen recommends entering it and makes explicit that the method advises rather than owns the repository.

**Carrier inventory:** the core and system architecture; C-002, C-005, C-012, and C-013; CAP-005, CAP-006, and AC-139; PDR-002, PDR-007, PDR-011, PDR-012, PDR-013, PDR-016, PDR-018, PDR-040, and PDR-041; the repository and scaffolded routing hubs and corpus introductions; the canonical change-routing, review-boundary, local-conventions, `clue-delta`, and `clue-verify` skill sources and their generated directories; the repository and scaffolded pull-request templates; contributor guidance and the public change-loop, methodology, skills, adoption, operations, CI-wall, design, and introduction pages; repository and reusable CI scope detection; and the tests that pin generated routing, override recognition, release locality, and acceptance-brief gating. Historical analyses, completed plans, superseded decision text below its amendment notice, and changelog entries remain pinned history.

### Rejected: use diff size or protected paths as the route

A one-line condition can create a capability and a large refactor can preserve behavior. Paths still select relevant checks and warn the agent where meaning often lives, but they cannot decide intent.

### Rejected: make user override a corpus log row

Adding `/docs` bookkeeping to escape the `/docs` workflow recreates the cost the override declines and misstates one authorization as enduring project meaning. Git already carries the integration history.

### Rejected: define a release route for adopters

Release policy belongs to each repository. Promoting this repository's administrative cut into generated methodology would make a local convention appear universal.
