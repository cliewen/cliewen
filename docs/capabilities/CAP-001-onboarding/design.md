---
id: CAP-001-design
type: design
status: active
links: [CAP-001, ADR-018]
title: Design for CAP-001 onboarding
---

# Design — CAP-001 Onboarding

## The command

`clue init [path]` materializes the whole convention in one call: the `/docs` taxonomy (folder READMEs that explain each record type, plus the decision log), the AGENTS.md routing hub, the five skills, and a CI workflow template. The scaffolding is embedded in the binary ([ADR-018](../../decisions/ADR-018-init-templates-embedded.md)): no network, no credentials, no checkout of this repo — the installed binary is the whole install.

- **Skills land twice**: `.agents/skills/` is the canonical location; `.claude/skills/` is the Claude Code mirror (`SKILL.md` spelling). The layout matches the first adopter repo, where it was proven by hand before being automated.
- **The CI template arms itself**: it carries the vendored-binary pattern from the adopter proof (checksum-verified release binary, `clue validate --forbid-changes`), with the vendoring commands in workflow comments and `CLUE_VERSION` pinned to the emitted skills' stamp. Until the binary is vendored the job passes with a visible warning naming the arming commands — a fresh `init` repo is never red before its first change, and an unarmed wall stays loudly unarmed on every run.
- **Idempotent**: a re-run regenerates README index blocks (between `<!-- clue:index:start/end -->` markers) from folder contents — existing entries whose targets survive keep their hand-written descriptions, missing entries are appended — and touches nothing else. A pre-existing taxonomy README without markers gains an appended index block (prose intact), so an existing repo still reaches green; a lone or reversed marker is ambiguous, and init errors naming the file rather than guess at the block's bounds. A pre-existing folder without the README `validate` requires is named in the report (`missing`) — init never invents it. Existing files are **never replaced**: they are reported as skipped, so an adopter with their own AGENTS.md keeps it (the repo-local layer, ADR-013).
- **The template-vs-validator contract is a test, not discipline**: AC-002's test runs `clue validate` over a fresh `init` output, so a validator rule the templates violate fails the build the moment it lands.

## The guide

The quickstart page lives in the repo README; the command is the guide's most important layer and must not require reading anything. Layers, kept strictly separate:

1. **Command (seconds):** `clue init` — the foundation in one call.
2. **Quickstart (5 minutes):** README — install, `clue init`, first change loop, watch `validate` go green. Owned by [QS-002](../../quality/QS-002-onboarding-under-30-minutes.md): under 30 minutes, reading nothing beyond the quickstart.
3. **Skills** — learned during use; the quickstart links each skill at the moment the reader's first change loop needs it (delta when they branch, verify before their first PR), never all four upfront. The system-level "why these skills" story lives in [architecture/skills.md](../../architecture/skills.md).
4. **Public guide** — the why and the full working model for newcomers who are not yet inside a Cliewen repository. It is a handwritten VitePress site under `/guide`, kept outside the corpus and deployed by the architecture in [ADR-023](../../decisions/ADR-023-public-guide-architecture.md), with [ADR-024](../../decisions/ADR-024-custom-domain-root.md) making `https://cliewen.dev/` its canonical root. The root README remains the shortest mechanical path; the guide adds depth without adding a prerequisite before first green validate.

The generated corpus explains itself: each emitted folder README says in plain language what its record type is and when a change updates it, and the emitted `docs/README.md` carries the map ("what lives where — and when a change updates it"). The guide points into the generated repo instead of duplicating it.

## Guide requirements (accumulating)

Lessons about what the guide must contain, appended the moment they are learned:

- **Prerequisites must be explicit** (learned 2026-07-12): git, the `clue` binary, and — for the PR-based change loop as practiced — an authenticated `gh` CLI. The loop works with plain git and any forge; `gh` is the convenient path, not a dependency of the method.
- **Skills are guide layer 3, learned during use** — see layer list above.
- **Value precedes vocabulary** (learned 2026-07-22): the public front door names the agent-written-change problem, intended audience, concrete outcome, and the distinction between Cliewen, `clue`, and the corpus before asking a newcomer to learn the taxonomy.
- **The first proof is disposable and observable** (learned 2026-07-22): Getting Started reaches the QS-002 green scaffold first in an empty Git repository, then activates a draft criterion without evidence so the released validator produces its real missing-test diagnostic; ownership, recovery, and full-directory cleanup are explicit.
- **Release binaries are the primary install path** (learned 2026-07-22): a compact platform table identifies the published architecture asset, and short manual steps verify it against `SHA256SUMS`, put it on the user `PATH`, and show the version result. The current unsigned macOS path names the checksum-first Gatekeeper exception explicitly. Source installation is secondary, so using Cliewen does not appear to require Go, and the first encounter does not present installer scripts.
- **The guide distinguishes machine and method guarantees** (learned 2026-07-22): `clue validate` currently requires at least one supported AC reference and does not run tests or count a positive/negative pair; the change loop and human review hold the stronger evidence rule.
- **CI enforcement must be reproducible** (learned 2026-07-22): the guide gives the exact pinned release vendoring and checksum path, distinguishes an unarmed warning from an armed `--forbid-changes` run, names the minimum protected-branch contract, and ends with a disposable failing-pull-request probe.
- **Adoption starts with the smallest useful thread** (learned 2026-07-22): begin with a goal, capability, acceptance criterion, and focused positive and negative evidence. Add plans, decisions, constraints, quality scenarios, architecture, and analysis only when their triggering problem exists, and state the repository and workflow conditions that make Cliewen a poor fit.
