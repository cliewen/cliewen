---
id: AN-010
type: analysis
status: active
links: [P-007, M-029, PDR-002, PDR-011]
title: The first adopter history does not support widening the light tier
---

# AN-010 — The first adopter history does not support widening the light tier

## Risk

Cliewen may require a full proposal workspace for ordinary product changes without evidence that the extra work protects meaning often enough to justify its cost.

## Evidence boundary

The spike inspected the public Robocode Tank Royale repository, the first product repository to adopt Cliewen, from its adoption merge `b7fb320ccec1f4742ef923cb315e7dd84f7e824f` (PR [#218](https://github.com/robocode-dev/tank-royale/pull/218)) through fetched `origin/main` commit `86a4bd58514bcdc4d36f1dd374900e6eae3b29f3`, pinned on 2026-07-28. That boundary contains five accepted first-parent units between 2026-07-19 and 2026-07-23: four merge commits for PRs [#219](https://github.com/robocode-dev/tank-royale/pull/219)…[#222](https://github.com/robocode-dev/tank-royale/pull/222) and one direct commit. Measurements were reproduced on Microsoft Windows NT 10.0.26200.0 with PowerShell 7.6.4, Git 2.55.0.windows.3, and GitHub CLI against the hosted pull-request metadata.

An accepted unit is one first-parent mutation of `main`. For a merge, commit count and churn cover commits reachable from the pull request's second parent but not its first parent; for the direct commit, they cover that commit. `git show --numstat` supplies line churn. `changes/` is transient process material, `docs/` is durable corpus material, generated Cliewen skills are reported separately, and every other path is product or repository material. Binary rows have no line count. The numbers establish versioned activity, not author effort or maintainer intent.

## Observed facts

| Accepted unit | Hosted tier claim or observable shape | Branch commits | Transient `changes/` churn | Durable corpus churn | Other churn |
|---|---|---:|---:|---:|---:|
| CH-002 / PR [#219](https://github.com/robocode-dev/tank-royale/pull/219) | Full workspace; criteria and plan meaning reconciled to shipped behavior | 3 | +102 / -102 | +49 / -15 | 0 |
| CH-003 / PR [#220](https://github.com/robocode-dev/tank-royale/pull/220) | PR explicitly claims light; four provenance flags promoted | 1 | 0 | +4 / -4 | 0 |
| CH-004 / PR [#221](https://github.com/robocode-dev/tank-royale/pull/221) | PR explicitly claims light; generated skills, CI installation, and plan bookkeeping changed | 1 | 0 | +2 / -2 | skills +187 / -63; repository +9 / -14 |
| Dependency update `5173b3f9` | Direct commit with no associated pull request | 1 | 0 | 0 | repository +7 / -7 |
| CH-005 / PR [#222](https://github.com/robocode-dev/tank-royale/pull/222) | Full workspace; a new plan and decision-log row added | 3 | +42 / -42 | +27 / -1 | 0 |

The two full semantic changes authored 144 transient lines and deleted the same workspace during digest, against 76 durable corpus additions: 1.89 transient lines per durable line added. Each used three branch commits, while both hosted light changes used one. Across all five accepted units the median transient churn is zero because three units carried no workspace.

No accepted unit in the boundary changes product behavior under an existing acceptance criterion. CH-002 changes criterion meaning to describe behavior shipped before adoption; CH-003 changes provenance only; CH-004 updates method distribution and CI; the dependency update changes build inputs; CH-005 creates a plan. The history therefore measures the workspace cost on semantic corpus work, but contains no observed instance of the proposed wider category.

Two conformance observations are visible but do not establish intent. CH-004 calls itself light while changing Cliewen skill carriers, a surface PDR-002 excludes from the light tier. The dependency update has no associated pull request even though the branch-and-PR invariant applies to dependency bumps. Activity can show that routing was not carried consistently; it cannot say whether the cause was misunderstanding, deliberate bypass, or tooling outside the agent workflow.

## Inferences

The full workspace has material, countable cost when it applies: in this boundary it nearly doubles authored durable-corpus lines and adds two commits. It is not the cost paid by the modal accepted unit in this small sample, because light routing already removed it from three of five units.

The sample does not support widening light changes to product behavior under existing criteria. Such a change still modifies executable evidence of what the product promises, and the only relevant observed semantic correction, CH-002, demonstrates why that alignment deserves a durable proposal: the corpus was structurally green while describing the wrong shipped credential and dry-run behavior. Removing the workspace from an unobserved category would be a policy choice made in advance of evidence, the opposite of M-029's evidence requirement.

The more defensible response is to keep behavior changes full, state why in the public guide, and reduce a different cost the evidence and AN-008 both expose: finding the relevant durable context. A deterministic `clue context <id>` slice reduces reading without weakening the proposal, traceability, or human merge boundary.

## Rejected measurements

- **Pull-request wall-clock duration:** PRs in the boundary were open from under a minute to more than three hours, but that interval mixes human availability, CI, review, and idle time. It is elapsed hosting time, not change-loop effort.
- **Raw total lines as effort:** CH-004 is dominated by generated skill synchronization and one binary deletion. Generated and binary churn says little about cognitive cost, so the table separates it instead of folding it into an overhead ratio.
- **Pre-adoption product history:** it can describe the repository's older change mix but cannot measure Cliewen overhead because those changes did not run the Cliewen loop.
- **Cliewen's own history:** methodology-on-methodology changes select for full semantic work and are explicitly excluded by M-029 as a cost baseline.

## Consumer

PDR-018 consumes this finding by retaining the full tier for behavior changes and making the evidence-based defence explicit in the public guide. CH-076 consumes the reading-cost finding through `clue context <id>` and focused agent routing. A later adopter history containing behavior changes under existing criteria may reopen the tier boundary with direct evidence.
