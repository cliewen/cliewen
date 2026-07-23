---
id: AN-005
type: analysis
status: active
provenance: inferred
links: [P-002, M-007, PDR-005, AN-004]
title: Foreign-soil trial — es-toolkit
---

# AN-005 — Foreign-soil trial: es-toolkit

## Classification

This is the second qualifying foreign-soil trial for P-002/M-007. [`toss/es-toolkit`](https://github.com/toss/es-toolkit) is an external public open-source repository with no shared maintainer, was not designed for Cliewen, and was explicitly selected by the human. The trial was read-and-apply analysis, not adoption: no tracked file or hosted state in the target repository was changed.

The spike is complete, so the analysis status is `verified`; its interpretation remains `provenance: inferred` because no es-toolkit maintainer verified the reconstructed intent.

## Risk named

Cliewen may mistake detailed agent instructions, active compatibility work, or a repository-declared verification command for an accepted goal, plan, or uniformly green evidence boundary when applying the methodology to a large repository with its own conventions.

## Evidence boundary

The trial used a detached disposable clone of `toss/es-toolkit` at [`fd4877295443a92655530735b0058c1d9ba1db4c`](https://github.com/toss/es-toolkit/commit/fd4877295443a92655530735b0058c1d9ba1db4c), the tip of `main` when CH-026 was proposed. Evidence came from committed files and history plus public repository metadata, issues, pull requests, and commits observed on 2026-07-18.

Verification ran on Windows 11 build 26200 in PowerShell. The repository pins Node.js 24.13.0 in `.nvmrc` and Yarn 4.12.0 in its package-manager declaration. An initial run used the installed Node.js 24.0.0; the qualifying rerun used a disposable official Node.js 24.13.0 Windows x64 distribution whose archive matched the SHA-256 published by nodejs.org, with Corepack 0.34.5 and Yarn 4.12.0.

Dependency installation and builds populated ignored outputs and local caches; the detached target checkout remained free of tracked changes. The trial created no branch, commit, pull request, issue, comment, or other hosted state in the target.

## What was tried

- Reconstructed a candidate goal, capability model, acceptance boundary, quality boundary, and possible future direction from the README, package metadata, root `AGENTS.md`, contributor guide, source, tests, documentation, workflows, history, and public GitHub activity.
- Applied `clue-analysis`, `clue-plan`, `clue-delta`, and `clue-verify` as read-only reasoning workflows without materializing a Cliewen corpus in the target.
- Compared the target's existing agent instructions and contribution workflow with the rules a Cliewen adoption would introduce.
- Ran immutable dependency installation, the full Vitest suite, package construction, TypeScript checking, full-tree linting, type tests, distribution checking, package-type validation, and bundle-size tests.
- Repeated the build and test boundary under the exact declared Node.js runtime after the installed runtime produced a different result.
- Rejected `clue init`, a target branch, and a target pull request because PDR-005 defines this as findings work rather than adoption.

## Evidence-qualified interpretation

| Claim | Confidence | Basis |
|---|---|---|
| es-toolkit is a modern, high-performance JavaScript utility library that emphasizes strong types, small bundles, and zero runtime dependencies | observed | `README.md`, package metadata, root `AGENTS.md` |
| The main entry point is intentionally strict and optimized for common cases, while `es-toolkit/compat` targets Lodash compatibility and accepts parity fixes rather than new functions | observed | Root `AGENTS.md` |
| Core utilities, compatibility utilities, types, documentation, packaging, and benchmarks are useful capability boundaries | inferred | Source layout, workspaces, focused tests, and workflows |
| Current public activity is concentrated on compatibility fixes, documentation, tests, and package-entry-point maintenance | observed | Recent commits, pull requests, and issues |
| That activity is an accepted future plan | unverified | No committed roadmap or open GitHub milestone was found; activity alone does not establish intent |
| Cliewen initialization preserves the target's existing `AGENTS.md` and therefore would not install Cliewen's routing rules there | observed | Cliewen's onboarding contract and the target's existing root `AGENTS.md` |

The inferred goal is: provide a fast, type-safe, compact set of modern JavaScript utilities, with an explicitly separate compatibility surface for users migrating from Lodash.

Likely capabilities are core utility behavior, Lodash-compatible behavior, type contracts, module and package entry points, multilingual reference documentation, bundle-size control, and performance comparison. These boundaries are useful for navigation but are not presented as maintainer-approved taxonomy.

No accepted plan could be reconstructed. Recent compatibility work is strong evidence of current activity and the root instructions define enduring contribution policy, but neither is a maintainer-approved sequence of milestones. `clue-plan` therefore stops at candidate direction rather than converting issues and pull requests into a plan.

## Native verification

| Invocation | Runtime | Result |
|---|---|---|
| `yarn install --immutable` | Node.js 24.0.0 | passed with peer-dependency warnings |
| `yarn prepack` | Node.js 24.0.0 | failed in a Node customization hook with `ERR_INVALID_RETURN_PROPERTY_VALUE` |
| `yarn vitest run` | Node.js 24.0.0 | failed: 640 files passed and one failed; 4,430 tests passed and four skipped; the distribution suite reached the same prepack failure |
| `yarn install --immutable` | exact Node.js 24.13.0 | passed with peer-dependency warnings |
| `yarn prepack` | exact Node.js 24.13.0 | passed; the build verified 297 compatibility entry points |
| `yarn vitest run` | exact Node.js 24.13.0 | failed: 640 files passed and one failed; 4,431 tests passed, one failed, and two skipped |
| `yarn tsc --noEmit` | exact Node.js 24.13.0 | passed |
| `yarn lint` | exact Node.js 24.13.0 | failed with eight errors and 2,469 warnings |
| `yarn lint --quiet` | exact Node.js 24.13.0 | failed with eight tracked missing-curly-brace errors in `docs/.vitepress/components/FpLazySimulation.vue` |
| `yarn workspace type-tests test` | exact Node.js 24.13.0 | passed: one file and one test |
| `node --test .scripts/check-dist.mjs` | exact Node.js 24.13.0 | passed |
| `yarn attw <temporary package archive>` | exact Node.js 24.13.0 | passed for all listed entry points |
| `yarn workspace benchmarks vitest run bundle-size` | exact Node.js 24.13.0 | passed: 100 files and 173 tests |

The remaining exact-runtime Vitest failure is localized to `tests/check-dist.spec.ts`. The test uses `path.join('es-toolkit', entrypoint)` to construct a package specifier embedded in generated JavaScript. On Windows, the resulting backslash is not a portable package-specifier separator and is interpreted inside the generated source. The direct distribution checker, package build, package-type validation, type checking, type tests, and bundle-size tests passed. The repository's main suite runs on Linux in CI, so that path does not expose this Windows-specific defect.

The full-tree lint result and the hosted lint gate also have different scopes. The contributor-facing `yarn lint` command finds eight existing errors in a tracked documentation component, while CI lints and formats only changed files. This trial records the declared command, clean-tree baseline, and hosted gate separately instead of reducing them to a single green or red repository state.

## Findings

### F1 — Exact declared runtimes are part of clean evidence

Node.js 24.0.0 and 24.13.0 are the same major release but did not produce the same build result. The exact pinned runtime converted the prepack failure into a pass. Recording only “Node 24” would make the first result difficult to interpret and the second difficult to reproduce. This independently validates the evidence-boundary rule added after AN-004.

### F2 — Declared commands, clean-tree baselines, and CI gates are distinct

The repository tells contributors to run full linting, but its CI lint job scopes itself to changed files. A detached clean checkout can therefore fail the documented full-tree command on existing tracked code while the hosted gate remains green. Cliewen verification must preserve those scopes rather than treating a command name as proof that the same surface is gated everywhere.

### F3 — A platform-specific test defect need not invalidate adjacent evidence

The Windows-only generated-specifier failure is reproducible and real, but it does not erase the passing build, direct distribution checker, package-type validation, type checks, type tests, and bundle-size suite. Evidence should identify the smallest demonstrated failure boundary instead of labeling the whole package either working or broken.

### F4 — Existing agent instructions are an adoption boundary

The target already uses a detailed root `AGENTS.md` to define product boundaries, contribution policies, verification, and changelog formatting. Cliewen's initializer deliberately preserves existing files, so running it would not install the routing hub that makes the generated skills discoverable. A real adoption would require an explicit integration choice.

The target's instruction to wrap changelog prose near 80 columns also conflicts with Cliewen's no-hard-wrap Markdown constraint. The trial surfaces that conflict without silently choosing one repository's rule over the other.

### F5 — Detailed policy is not a plan

The root instructions make the product boundary unusually recoverable: the strict main entry point and compatibility entry point have different acceptance policies. Recent public work also clusters around compatibility. Neither fact creates an accepted milestone sequence. The activity-versus-intent distinction from CH-025 prevented a plausible roadmap from being invented.

### F6 — Deterministic contracts and quality evidence remain different

Utility behavior, types, exports, compatibility examples, and bundle snapshots can be deterministic acceptance surfaces. Claims about outperforming alternatives, percentage bundle reductions, and benchmark leadership depend on comparison baselines and execution conditions even when current snapshot limits are deterministic. This confirms the acceptance/quality split from AN-003 and AN-004 in another ecosystem without requiring a new rule in this change.

## Cross-trial conclusions

The hyperfine and es-toolkit trials used different languages, repository structures, contribution models, and product surfaces, yet both produced useful candidate goals and capability boundaries without modifying their targets. In both, public activity required an explicit confidence boundary before it could be discussed as direction.

AN-004 showed that shell identity changed an integration-test result. AN-005 shows that an exact minor runtime changed a package build and that platform path semantics changed one distribution test. Together they validate the CH-025 rule that a reproducible evidence boundary includes source revision and the conditions that materially affect results.

Both repositories also separate deterministic behavior from environment-sensitive quality. Hyperfine measures statistical command performance; es-toolkit advertises utility performance and bundle efficiency. Cliewen's existing acceptance and quality taxonomy can represent both without turning measurements into deterministic promises.

The two trials expose adoption work that findings alone must not conceal: foreign repositories already have governance, contributor commands, CI scopes, documentation conventions, and agent instructions. PDR-005's choice to trial the skills without scaffolding was therefore useful isolation, not a substitute for a future adoption exercise.

## Methodology-adjustment assessment

No additional methodology adjustment is warranted by AN-005. The evidence-boundary rule introduced by CH-025 directly governs the reproduced runtime, platform, and activity-versus-intent risks and was successfully followed without dropping any of its obligations. The remaining findings concern adoption integration or distinctions already carried by Cliewen's initializer, plan, verification, constraint, and acceptance/quality contracts.

This is a positive cross-trial result rather than an attempt to manufacture a second change. M-007 requires at least one traceable adjustment; CH-025 supplies it, and AN-005 independently validates it.

## Rejected interpretations

- **The installed Node.js 24 runtime is close enough to the declared one.** The prepack result changed under the exact pinned 24.13.0 runtime.
- **One failing Vitest file means the package cannot be built or distributed.** Focused package construction, distribution, type, and bundle checks passed; the demonstrated failure is Windows-specific generated source in one test.
- **Green changed-file linting means the documented full lint command is green on the clean tree.** The two commands have different scopes, and the full-tree command finds existing tracked errors.
- **Recent compatibility activity is the plan.** The activity is observed; no accepted milestone sequence was found.
- **A pre-existing `AGENTS.md` means Cliewen can be initialized without integration work.** Initialization preserves that file and therefore does not install Cliewen's routing hub.
- **A useful trial must adopt Cliewen into the target.** PDR-005 deliberately separates methodology analysis from adoption mechanics.

## Trial result

Cliewen recovered a useful goal, differentiated product policies, candidate capability model, deterministic contract surface, quality boundary, and adoption conflicts without modifying the target. The trial independently validated the evidence-boundary adjustment produced by AN-004 and completed the second human-selected external trial required by M-007. Across both qualifying repositories, Cliewen's read-and-apply workflows proved useful while preserving uncertainty and foreign governance.
