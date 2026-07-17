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
- **The CI template** carries the vendored-binary pattern from the adopter proof (checksum-verified release binary, `clue validate --forbid-changes`), with the vendoring commands in workflow comments and `CLUE_VERSION` pinned to the emitted skills' stamp.
- **Idempotent**: a re-run regenerates README index blocks (between `<!-- clue:index:start/end -->` markers) from folder contents — existing entries whose targets survive keep their hand-written descriptions, missing entries are appended — and touches nothing else. Existing files are **never overwritten**: they are reported as skipped, so an adopter with their own AGENTS.md keeps it (the repo-local layer, ADR-013).
- **The template-vs-validator contract is a test, not discipline**: AC-002's test runs `clue validate` over a fresh `init` output, so a validator rule the templates violate fails the build the moment it lands.

## The guide

The quickstart page lives in the repo README; the command is the guide's most important layer and must not require reading anything. Layers, kept strictly separate:

1. **Command (seconds):** `clue init` — the foundation in one call.
2. **Quickstart (5 minutes):** README — install, `clue init`, first change loop, watch `validate` go green. Owned by [QS-002](../../quality/QS-002-onboarding-under-30-minutes.md): under 30 minutes, reading nothing beyond the quickstart.
3. **Skills** — learned during use; the quickstart links each skill at the moment the reader's first change loop needs it (delta when they branch, verify before their first PR), never all four upfront. The system-level "why these skills" story lives in [architecture/skills.md](../../architecture/skills.md).
4. **Book** — the why; depth, secondary.

The generated corpus explains itself: each emitted folder README says in plain language what its record type is and when a change updates it, and the emitted `docs/README.md` carries the map ("what lives where — and when a change updates it"). The guide points into the generated repo instead of duplicating it.

## Guide requirements (accumulating)

Lessons about what the guide must contain, appended the moment they are learned:

- **Prerequisites must be explicit** (learned 2026-07-12): git, the `clue` binary, and — for the PR-based change loop as practiced — an authenticated `gh` CLI. The loop works with plain git and any forge; `gh` is the convenient path, not a dependency of the method.
- **Skills are guide layer 3, learned during use** — see layer list above.
