---
id: AN-010
type: analysis
status: active
links: [P-007, M-029, PDR-002, PDR-011, P-013, PDR-026]
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

## Re-derived at head (P-013/M-066, 2026-08-10)

[PDR-026](../decisions/PDR-026-campaigns-close-on-re-derived-evidence.md) requires this cost surface to be measured again rather than read off a milestone table. The boundary extends from this finding's pin `86a4bd58514bcdc4d36f1dd374900e6eae3b29f3` (exclusive) through the adopter's current `origin/main` tip `2005b1d3a77f564a4714a7b0e64ed5545a1a396f` (2026-08-10), the merge commit for PR [#230](https://github.com/robocode-dev/tank-royale/pull/230): nine more first-parent accepted units, none beyond it as of this pin. Measurements repeat the original method — `git log --numstat` over each unit's first-parent/second-parent range, generated skills reported separately — on Microsoft Windows NT 10.0.26200.0 with PowerShell 5.1.26100.8875 and Git 2.55.0.windows.3.

| Accepted unit | Hosted tier claim or observable shape | Branch commits | Transient `changes/` churn | Durable corpus churn | Other churn |
|---|---|---:|---:|---:|---:|
| CH-006 / PR [#223](https://github.com/robocode-dev/tank-royale/pull/223) | PR body states **"full"** explicitly | 2 | +71 / -71 | +71 / -4 | skills +110 / -46; repository +104 / -1 |
| CH-007 / PR [#224](https://github.com/robocode-dev/tank-royale/pull/224) | PR body states **"light"** explicitly | 2 | 0 | 0 | skills +48 / -23; repository +33 / -67 |
| CH-008 / PR [#225](https://github.com/robocode-dev/tank-royale/pull/225) | PR body states **"light"** explicitly | 2 | 0 | +9 / -9 | skills +14 / -11; repository +8 / -7 |
| CH-009 / PR [#226](https://github.com/robocode-dev/tank-royale/pull/226) | No literal tier word; full workspace opened and removed within the unit | 6 | +90 / -90 | +145 / -44 | repository +427 / -26 |
| Direct commit `aa2018f1` | No pull request; single commit | 1 | 0 | 0 | repository (`AGENTS.md`) +4 / -0 |
| CH-010 / PR [#227](https://github.com/robocode-dev/tank-royale/pull/227) | No literal tier word; full workspace opened and removed within the unit | 7 | +92 / -92 | +93 / -5 | repository +175 / -115 |
| CH-011 / PR [#228](https://github.com/robocode-dev/tank-royale/pull/228) | No literal tier word; full workspace opened and removed within the unit; no product code touched | 5 | +80 / -80 | +94 / -3 | 0 |
| CH-012 / PR [#229](https://github.com/robocode-dev/tank-royale/pull/229) | No literal tier word; PR body self-describes as a paused draft | 3 | +75 / -0 (workspace opened, not closed) | +1 / -1 | repository (`.agents/principles-catalog/active.md`) +92 / -0 |
| CH-012 / PR [#230](https://github.com/robocode-dev/tank-royale/pull/230) | No literal tier word; PR body self-describes as a draft continuation | 2 | +5 / -4 (workspace still open, net +1) | +267 / -34 | 0 |

`.agents/principles-catalog/active.md` is reported on its own line rather than folded into "other," for the same reason generated skills are: one 92-line single-block addition would otherwise dominate and mask the unit's actual repository churn. It appeared once, in one unit, so it is a callout here rather than a fourth standing column.

**What changed:** the sample nearly doubled (5 → 14 accepted units) and, for the first time, contains a full-tier workspace that opens in one accepted unit and is still open two units later — CH-012 spans PR #229 and #230, both self-described as drafts, and `changes/CH-012-create-rumble-client/` remains in the tree at the pin. No unit among the nine modifies an existing acceptance criterion's meaning; every criterion-touching unit adds brand-new draft criteria for a new capability, the same pattern the original five showed. **What did not change:** the light/full split by literal tier word remains legible where a unit states one (2 full, 2 light, 5 unstated), the no-PR direct-commit pattern recurs (`aa2018f1`, no associated pull request, same shape as the original `5173b3f9`), and the CH-004-shaped tension recurs twice more — PR #224 and PR #225 both claim "light" while touching `.agents/skills/*`, the surface PDR-002's light-tier qualification excludes. **What the campaign declined:** re-deriving a transient-per-durable-line ratio across the extended sample the way the original Inferences section did for five units — CH-012's still-open workspace has no closing churn to divide by yet, and folding an in-progress unit into a ratio computed for closed ones would misstate what either number means; whether CH-012's open state is normal churn caught mid-flight or a gap in this adopter's close-out convention is left as an observation, not a finding this closing milestone adjudicates.
