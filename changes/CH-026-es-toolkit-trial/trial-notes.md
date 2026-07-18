---
id: CH-026-trial-notes
type: analysis
status: verified
links: [CH-026]
title: Evidence notes for the es-toolkit foreign-soil trial
---

# Evidence notes

These transient notes preserve the trial evidence before it is synthesized into AN-005. They are not an adoption record for the target repository.

## Boundary

- Target: [`toss/es-toolkit`](https://github.com/toss/es-toolkit) at commit [`fd4877295443a92655530735b0058c1d9ba1db4c`](https://github.com/toss/es-toolkit/commit/fd4877295443a92655530735b0058c1d9ba1db4c), the tip of `main` when CH-026 was proposed.
- Checkout: detached disposable clone outside both repositories; the target remained free of tracked changes.
- Host: Windows 11 build 26200, PowerShell.
- Declared runtime: Node.js 24.13.0 from `.nvmrc`; package manager Yarn 4.12.0 from the committed package-manager declaration.
- Initial runtime: system Node.js 24.0.0 with Corepack 0.32.0 and Yarn 4.12.0.
- Exact runtime: portable official Node.js 24.13.0 with Corepack 0.34.5 and Yarn 4.12.0. Its Windows x64 ZIP matched the SHA-256 published by nodejs.org: `ca2742695be8de44027d71b3f53a4bdb36009b95575fe1ae6f7f0b5ce091cb88`.
- Side effects: dependency installation populated ignored build/dependency outputs and the package-manager cache. No target branch, commit, pull request, issue, comment, or other hosted state was created.

## Public evidence

| Evidence | Observation | Confidence |
|---|---|---|
| `README.md`, root `AGENTS.md`, `.github/CONTRIBUTING.md`, package metadata | The repository describes a modern, high-performance JavaScript utility library with strong types, small bundles, and zero runtime dependencies | observed |
| Root `AGENTS.md` | The main entry point is intentionally strict and optimized for common cases; the `compat` entry point targets Lodash compatibility and accepts parity fixes rather than new functions | observed |
| Source layout and focused tests | Core utilities, compatibility utilities, types, documentation, packaging, and benchmarks are useful working capability boundaries | inferred |
| Open issues, recent pull requests, and recent commits | Current public activity is concentrated on compatibility fixes, documentation, tests, and package-entry-point maintenance | observed activity |
| No GitHub milestones and no committed roadmap artifact found | The activity represents an accepted future plan | unverified |
| Root `AGENTS.md` | Changelog prose should wrap near 80 columns | observed |
| Cliewen `clue init` behavior | A pre-existing `AGENTS.md` is preserved, so Cliewen's routing hub would not be installed automatically | observed from Cliewen's tested onboarding contract |

## Native verification

| Invocation | Runtime | Result |
|---|---|---|
| `yarn install --immutable` | Node.js 24.0.0 | passed with peer-dependency warnings |
| `yarn prepack` | Node.js 24.0.0 | failed in a Node customization hook with `ERR_INVALID_RETURN_PROPERTY_VALUE` |
| `yarn vitest run` | Node.js 24.0.0 | failed: 640 files passed and one failed; 4,430 tests passed and four skipped; the failing distribution suite reached the same prepack failure |
| `yarn install --immutable` | exact Node.js 24.13.0 | passed with peer-dependency warnings |
| `yarn prepack` | exact Node.js 24.13.0 | passed; the build verified 297 compatibility entry points |
| `yarn vitest run` | exact Node.js 24.13.0 | failed: 640 files passed and one failed; 4,431 tests passed, one failed, and two skipped |
| `yarn tsc --noEmit` | exact Node.js 24.13.0 | passed |
| `yarn lint` | exact Node.js 24.13.0 | failed with eight errors and 2,469 warnings |
| `yarn lint --quiet` | exact Node.js 24.13.0 | failed with eight tracked errors in `docs/.vitepress/components/FpLazySimulation.vue`, all requiring curly braces |
| `yarn workspace type-tests test` | exact Node.js 24.13.0 | passed: one file and one test |
| `node --test .scripts/check-dist.mjs` | exact Node.js 24.13.0 | passed |
| `yarn attw <temporary package archive>` | exact Node.js 24.13.0 | passed for all listed entry points |
| `yarn workspace benchmarks vitest run bundle-size` | exact Node.js 24.13.0 | passed: 100 files and 173 tests |

The remaining exact-runtime Vitest failure is localized to `tests/check-dist.spec.ts`. The test constructs package specifiers with `path.join('es-toolkit', entrypoint)` and embeds the result in generated JavaScript. On Windows the backslash becomes part of the generated source rather than a portable package-specifier separator, so a specifier such as `es-toolkit\array` is parsed incorrectly. The direct distribution checker, package build, package-type validation, type checking, type tests, and bundle-size tests passed. The target's main test workflow runs on Linux, so its current CI does not expose this Windows path.

The full-tree lint command and the CI lint gate have different scopes. The contributor-facing `yarn lint` command finds eight existing errors in a tracked documentation component, while CI invokes linting and formatting only for changed files. The trial therefore records the repository-declared command, the clean-tree baseline, and the hosted gate separately rather than calling the repository either green or broken.

## Cliewen workflow application

- `clue-analysis` kept the exact source revision, runtime, host, shell, and invocation deviations attached to the results. Its observed/inferred/unverified distinctions prevented recent compatibility activity from being promoted to maintainer intent.
- `clue-plan` permits reconstructing candidate capabilities and an unaccepted draft, but it prevents creating milestones from issue or pull-request volume. No explicit target plan was found.
- `clue-delta` would add plan linkage, typed decisions, transient change workspaces, a human merge boundary, and a ready-for-review pull-request rule to the target's existing contribution workflow.
- `clue-verify` can consume the native build and test evidence, but it must not collapse full-tree commands, changed-file CI gates, platform-specific failures, and successful focused checks into one green/red label.
- Adoption was not attempted. A real adoption would need a deliberate integration of the existing `AGENTS.md`, not a claim that `clue init` installed Cliewen's routing rules. It would also need to resolve the target's approximately 80-column changelog rule against Cliewen's no-hard-wrap Markdown constraint.

## Candidate conclusions

1. The evidence-boundary rule added by CH-025 is both usable and necessary. Exact Node.js 24.13.0 changed a build failure into a pass; recording only “Node 24” would have hidden the cause.
2. Repository-declared contributor commands, full-tree baselines, and hosted CI gates are distinct evidence. A clean checkout can legitimately expose pre-existing full-lint errors while changed-file CI stays green.
3. A broad cross-platform product surface can have a localized Windows defect in test-generated source even when package construction and focused distribution checks pass.
4. Existing agent instructions are an adoption boundary. Cliewen preserves a pre-existing `AGENTS.md`, so initialization alone cannot establish the methodology's routing contract.
5. Goal and enduring product policies were explicit, but public activity did not become an accepted plan. Cliewen's distinction between activity and intent held.
6. Deterministic behavior, types, exports, and compatibility examples fit acceptance criteria. Performance comparisons, bundle-reduction claims, and benchmark leadership need environment- and baseline-qualified quality evidence, although deterministic snapshot gates can enforce a current threshold.
7. No further methodology adjustment is warranted by this trial. CH-025's evidence-boundary adjustment addresses the reproduced risk, and the remaining observations are adoption/integration findings already governed by existing Cliewen contracts.
