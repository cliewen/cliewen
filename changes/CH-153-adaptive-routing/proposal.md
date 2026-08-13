---
id: CH-153
type: change
status: open
links: [G-008, CAP-006, PDR-007, PDR-011, PDR-018]
title: Route by accepted-contract change and leave integration authority with the user
---

# CH-153 — Route by accepted-contract change and leave integration authority with the user

This plan-less change replaces the plain/light/full methodology hierarchy with two recommendations: simple work stays outside the Cliewen loop, while work that changes the accepted contract is recommended for the full loop. The agent states its recommendation before editing, names what would change it, reassesses on semantic discovery and before integration, and treats paths or diff size only as warning signals.

Simple work includes observational analysis with a named consumer, fixes that restore unchanged acceptance criteria, regression evidence for unchanged criteria, in-contract configuration adjustments, refactors, maintenance, and editorial corrections. Full work changes acceptance criteria, capabilities, decisions, policy, plan meaning, methodology, or behavior not covered by the accepted contract. If the user declines a full recommendation, the agent follows the user's authority, records the override in vendor-neutral Git trailers, and keeps the repository truthful and appropriately tested.

The methodology does not prescribe how a repository owner integrates work. Agent pushes still require user authorization and repository permission, while a repository may impose stricter local rules. Cliewen's own release cut remains a repository-local specialization of simple work rather than an adopter-facing release route, and this repository continues to require agents to use pull requests and human merge.

Record the process decision in PDR-042, revise CAP-006 through AC-139, update every live carrier and generated output atomically, and separate CI check selection from full-loop bookkeeping so an analysis-only corpus change is checked without being charged for an acceptance brief.
