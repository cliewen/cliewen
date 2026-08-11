---
id: CH-150-proposal
type: change
status: open
links: [G-007]
title: State the changelog-scope rule in terms of the shipped surface
---

# CH-150 — The changelog-scope rule names the surface that reaches an adopter

This change is plan-less. It serves G-007 by making the release-note scope rule say what its list of terms already means — that each names something reaching an adopter — and by carrying a test an agent can apply without re-deriving the boundary from [ADR-013](../../docs/decisions/ADR-013-ships-generic-vs-repo-local.md).

The rule is stated in three live carriers: [C-002](../../docs/constraints/C-002-changelog-per-user-visible-change.md), the release-notes convention in `AGENTS.md`, and the digest sentence in `CONTRIBUTING.md`. It is not a methodology contract that reaches adopters: the generated skills and their scaffolded copies deliberately say only that a repository's own conventions may require a user-facing changelog entry, and name no scope test, so they state nothing this change contradicts and are not edited. The inventory is therefore those three carriers and the decision recording the tightening.

The tightening adds normative content rather than only rewording: the test that distinguishes a repo-local change from an adopter-visible one — whether the change touches an artifact `clue init` or `clue scaffold` materializes into an adopter repository, or the behaviour of `clue` or a generated skill — is not stated anywhere today. That is why a decision record accompanies it rather than a decision-log row, which could carry the change but not the test and its rationale.
