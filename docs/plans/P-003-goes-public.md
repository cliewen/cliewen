---
id: P-003
type: plan
status: active
links: [G-003]
title: Cliewen goes public — readiness, guide, release, visibility flip
---

# P-003 — Cliewen goes public

The baseline ([P-001](P-001-elaboration-baseline.md)) proved the methodology at home; [P-002](P-002-leaves-home.md) made it installable and trialed it on foreign soil. This campaign makes the repository, its releases, and a human-readable guide publicly reachable, in the order [PDR-009](../decisions/PDR-009-going-public.md) fixes: readiness first, the v0.4.0 release, then the flip — which takes the guide live in the same act. Serves [G-003](../goals/G-003-cliewen-is-public.md). Milestone numbering is corpus-global and continues from P-002 (M-004…M-007).

## Milestones

| ID | Milestone (exit criterion) | Status | Evidence |
|---|---|---|---|
| M-008 | **The corpus is stranger-safe**: no unresolvable reference to a private repository's artifacts remains as a normative reference (PDR-009's citation policy applied); the analysis index notes that some cited adopter repos are private; the root README's install text and status section describe the current state, not the private era | `done` | CH-029: private-target evidence boundary documented; normative references resolve within the corpus; README install and status text describe the active private-readiness state |
| M-009 | **The community front door exists**: CONTRIBUTING, CODE_OF_CONDUCT, SECURITY, and issue/PR templates exist and route contributors into the change loop ([C-005](../constraints/C-005-proposal-declares-plan-item.md), [C-012](../constraints/C-012-agents-never-merge-own-changes.md), [PDR-007](../decisions/PDR-007-review-boundary.md)) rather than around it | `done` | CH-030: Contributor Covenant 3.0, coordinated vulnerability disclosure, structured bug and proposed-goal intake, contribution guidance, and a change-loop PR template; repository sanity coverage guards the files and form structure |
| M-010 | **The guide exists and deploys**: a handwritten VitePress guide builds green locally and in CI with dead-link checking; diagrams render as mermaid ([C-007](../constraints/C-007-diagrams-inline-mermaid.md) parity); a GitHub Pages deploy workflow exists, gated so it cannot fail while the repository is private; the site's architecture is recorded as an ADR | `done` | CH-032: newcomer guide covers the methodology, corpus, change loop, and skills; locked VitePress/Mermaid build is green locally and in CI with strict links; Pages deploy is visibility-gated and triggered by the public event; ADR-023 records the site boundary |
| M-011 | **v0.4.0 is cut as the goes-public release** ([PDR-009](../decisions/PDR-009-going-public.md)): the changelog's unreleased section becomes the 0.4.0 section with an install story free of private-repo caveats; the skills carry the 0.4.0 stamp; the release workflow publishes the cross-platform binaries and checksums with the changelog section as the release body ([ADR-012](../decisions/ADR-012-release-notes-from-changelog.md)) | `todo` | CH-035: reviewed 0.4.0 changelog section and generated skill stamps prepared; tag publication and release verification remain before the exit criterion closes |
| M-012 | **The repository is public and reachable**: in one act per [PDR-009](../decisions/PDR-009-going-public.md), visibility flips, the deploy gate opens, and the guide goes live on Pages; `go install` of v0.4.0 and anonymous release-asset download succeed without credentials; the required CI check still gates merges | `todo` | |

## Explicitly out of this campaign

A custom domain for the guide (changes the Pages base path; a door for later); rendering the corpus into the site; announcements and marketing; a GitHub Discussions strategy; and everything P-002 already excluded (multi-agent orchestration, semantic consistency checking, `clue locate`, production feedback loop, non-OpenSpec extraction mappings).

## Mutation rules (lintable)

Status fields in the milestone table may mutate in any merge (bookkeeping). Everything else in this file changes only via a change that declares itself a plan revision, backed by a decision record routed by reversal cost ([C-011](../constraints/C-011-decision-records-typed.md)): a PDR for direction and process, an ADR if architectural, a decision-log row where reversing is cheap and local. Plan adjustments ARE decisions.
