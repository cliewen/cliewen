---
id: G-007
type: goal
status: accepted
links: [G-003]
title: The changelog-scope rule names the surface that reaches an adopter
---

# G-007 — The changelog-scope rule names the surface that reaches an adopter

**Who wants it:** repository contributors and agents deciding whether a change owes a release note (2026-08-11), found while CH-148 removed an entry it had written for a fix to this repository's own contributor guidance.

**Why:** the rule asks for an entry from a change that affects "shipped behavior, a capability, a contract, a command, or a user workflow". Every term in that list reads naturally as this repository's own — this repository has capabilities, contracts, commands, and workflows that no adopter ever sees. Deciding correctly means reaching past the rule to [ADR-013](../decisions/ADR-013-ships-generic-vs-repo-local.md), which separates what ships from the repo-local layer, and re-deriving the boundary from it. An agent that does not make that leap writes an entry describing a verification block the reader does not have, and the release body carrying it is published verbatim. The two entries already in the changelog whose subject is contributor guidance are both genuinely adopter-visible for reasons the rule does not state, so the corpus looks inconsistent with itself to anyone checking precedent rather than reasoning from ADR-013.

**Success looks like:**

- The rule states that each term in its list names something reaching an adopter, so a repo-local change and an adopter-visible one are distinguishable from the rule alone.
- The rule carries a test an agent can apply without re-deriving the boundary from the decision behind it.
- Every live carrier stating the rule says the same thing.
