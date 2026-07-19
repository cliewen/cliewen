---
id: AN-006
type: analysis
status: verified
links: [P-003, M-011, PDR-002, PDR-007]
title: The light tier still overcharges changes outside Cliewen
---

# AN-006 — The light tier still overcharges changes outside Cliewen

## Risk

Cliewen may make an unrelated editorial change expensive enough that developers bypass or reject the methodology.

## Evidence boundary

The observed change is PR #35 at merge commit `b674407ea305ae5e437155ed66deb4afa403d67f`: one prose file under `guide/`, outside the `/docs` corpus, changed without code, configuration, test, decision, plan, or methodology-carrier impact. The reproduction environment was Windows amd64 with PowerShell 7.6.3, Go 1.26.5, Node.js 24.0.0, and npm 11.18.0. Command timings came from the agent tool output; token cost was visibly incurred but was not measured, so this finding makes no numerical token claim.

The maintainer explicitly identified the workflow as a critical adoption issue on 2026-07-20. That statement is evidence of maintainer intent, not an inference from repository activity.

## Observations

- PDR-002 removes the transient workspace for a light change but retains global CH numbering, a `ch-xxx-*` branch, proposal content in the PR, a plan item or plan-less declaration, and the full pre-PR checklist.
- AGENTS.md requires the corpus to be read before any task, so an agent pays the Cliewen context cost before it can conclude that Cliewen has nothing to contribute.
- The PR template requires traceability, decision, constraint, quality, digest, changelog, generation, full build, coverage, corpus validation, and review-boundary declarations for every PR.
- The repository CI runs the Go build and coverage suite, guide build, and corpus validation on every PR. For PR #35, the focused guide build completed in 7.3 seconds and the initial Go test pass completed in 7.2 seconds; the corpus validator could only confirm that the untouched corpus was still valid.
- `clue validate` scans `/docs` and `/changes`. It cannot produce change-specific evidence for a guide-only diff.
- Branching from accepted `main`, opening a PR, and leaving merge to a human match ordinary protected-repository practice and were not identified as the expensive part.

## Options assessed

Keeping the light tier unchanged preserves uniform provenance but repeats work that produces no evidence. Allowing direct pushes for small changes removes the review boundary and makes “small” a self-judged escape hatch. A path-only exemption is cheap to automate but lets a meaningful policy or contract change bypass Cliewen merely because it sits in an allowed folder.

A meaning-based plain route keeps the human merge boundary while excluding work that changes no behavior, intent, evidence, decision, plan, policy, or methodology carrier. Protected paths fail closed into Cliewen, and repositories may still choose relevant checks from the changed surface.

## Finding and consumer

Cliewen needs a route outside its light/full tiers. A plain change should use an ordinary branch, relevant checks, a PR, and human merge, with no CH identity or Cliewen artifacts. PDR-011 records the decision; CH-037 carries it into the routing hub, generated skills, contributor surfaces, and CI.
