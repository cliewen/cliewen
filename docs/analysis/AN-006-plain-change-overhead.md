---
id: AN-006
type: analysis
status: active
links: [P-003, M-011, PDR-002, PDR-007, P-013, PDR-026, PDR-011]
title: The light tier still overcharges changes outside Cliewen
---

# AN-006 — The light tier still overcharges changes outside Cliewen

## Risk

Cliewen may make an unrelated editorial change expensive enough that developers bypass or reject the methodology.

## Evidence boundary

The observed change is PR [#35](https://github.com/cliewen/cliewen/pull/35) at merge commit `b674407ea305ae5e437155ed66deb4afa403d67f`: one prose file under `guide/`, outside the `/docs` corpus, changed without code, configuration, test, decision, plan, or methodology-carrier impact. The reproduction environment was Windows amd64 with PowerShell 7.6.3, Go 1.26.5, Node.js 24.0.0, and npm 11.18.0. Command timings came from the agent tool output; token cost was visibly incurred but was not measured, so this finding makes no numerical token claim.

The maintainer explicitly identified the workflow as a critical adoption issue on 2026-07-20. That statement is evidence of maintainer intent, not an inference from repository activity.

## Observations

- PDR-002 removes the transient workspace for a light change but retains global CH numbering, a `ch-xxx-*` branch, proposal content in the PR, a plan item or plan-less declaration, and the full pre-PR checklist.
- AGENTS.md requires the corpus to be read before any task, so an agent pays the Cliewen context cost before it can conclude that Cliewen has nothing to contribute.
- The PR template requires traceability, decision, constraint, quality, digest, changelog, generation, full build, coverage, corpus validation, and review-boundary declarations for every PR.
- The repository CI runs the Go build and coverage suite, guide build, and corpus validation on every PR. For PR [#35](https://github.com/cliewen/cliewen/pull/35), the focused guide build completed in 7.3 seconds and the initial Go test pass completed in 7.2 seconds; the corpus validator could only confirm that the untouched corpus was still valid.
- `clue validate` scans `/docs` and `/changes`. It cannot produce change-specific evidence for a guide-only diff.
- Branching from accepted `main`, opening a PR, and leaving merge to a human match ordinary protected-repository practice and were not identified as the expensive part.

## Options assessed

Keeping the light tier unchanged preserves uniform provenance but repeats work that produces no evidence. Allowing direct pushes for small changes removes the review boundary and makes “small” a self-judged escape hatch. A path-only exemption is cheap to automate but lets a meaningful policy or contract change bypass Cliewen merely because it sits in an allowed folder.

A meaning-based plain route keeps the human merge boundary while excluding work that changes no behavior, intent, evidence, decision, plan, policy, or methodology carrier. Protected paths fail closed into Cliewen, and repositories may still choose relevant checks from the changed surface.

## Finding and consumer

Cliewen needs a route outside its light/full tiers. A plain change should use an ordinary branch, relevant checks, a PR, and human merge, with no CH identity or Cliewen artifacts. PDR-011 records the decision; CH-037 carries it into the routing hub, generated skills, contributor surfaces, and CI.

## Re-derived at head (P-013/M-066, 2026-08-10)

[PDR-026](../decisions/PDR-026-campaigns-close-on-re-derived-evidence.md) requires this cost surface to be measured again rather than read off a milestone table. The reproduction environment is Windows amd64 with PowerShell 7.6.4, Go 1.26.5, Git 2.55.0.windows.3, Node.js 24.0.0, and npm 11.18.0 — materially the same toolchain the original finding used. Every merged pull request on `github.com/cliewen/cliewen` was listed (`gh pr list --state merged`, 145 total) and filtered to the seven whose branch name carries no `ch-` prefix — the plain route's observable signature, since a plain change mints no CH identity or branch convention. Three sit after [PDR-011](../decisions/PDR-011-plain-changes-bypass-cliewen.md)'s acceptance (2026-07-20) and outside the later release tier: PR [#42](https://github.com/cliewen/cliewen/pull/42), PR [#66](https://github.com/cliewen/cliewen/pull/66), and PR [#75](https://github.com/cliewen/cliewen/pull/75).

PR #75, merged 2026-07-27, is the direct successor to this finding's original evidence: a one-line prose change to `guide/design.md`, outside the `/docs` corpus, with no code, configuration, test, decision, plan, or methodology-carrier impact — the same shape PR #35 had. It carries one commit, no CH identity, and opened and merged inside six minutes. Its hosted CI ran a single `validate` job, not the four separate checks (build, coverage-gated test, guide build, corpus validation) PR #35 paid: `Classify changed surface` and `Test CI scope classifier` ran in under a second, then `Require a completed acceptance brief`, `Build`, `Test (with coverage)`, `Coverage gate`, and `Validate corpus` all report `skipped`, leaving only `Build public guide` (19s, because the diff touched `guide/`) and `Check diff whitespace` to actually execute. The whole job completed in 27 seconds, start to finish.

**What changed:** the plain route now exists and is observably used — three plain PRs since acceptance, one of them the same shape as this finding's original evidence — and a change-scope classifier in CI skips the corpus validator, the coverage-gated Go test, and the acceptance-brief gate entirely for a guide-only diff, rather than running all of them and reporting only that the untouched corpus stayed valid. **What did not change:** a plain PR still pays for whatever check its actual changed surface triggers (here, the guide build, because the diff touched `guide/`) — the classifier narrows the check set to the changed surface, it does not remove checking; and this finding still makes no numerical token claim, for the same reason the original did not. **What the campaign declined:** re-deriving a churn or duration percentile across all seven plain PRs — three is too small a population for a distribution claim, and the single directly comparable pair (PR #35 and PR #75) is reported as a pair, not averaged into one.
