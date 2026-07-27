---
id: CH-068
type: change
status: open
links: [P-007, M-031]
title: Cliewen installs from Claude Code as a bootstrap plugin that deliberately ships none of the managed skills
---

# CH-068 — Cliewen installs from Claude Code

## What

A plugin marketplace manifest at the repository root makes `/plugin marketplace add cliewen/cliewen` the entry point for a coding agent. The marketplace lists exactly one plugin, and that plugin ships exactly one skill: a bootstrap that detects the host, installs `clue` through the channels [M-030](../../docs/plans/P-007-core-hardening.md) published, verifies the binary reports a release version rather than `dev`, and then **asks** before running `clue init`.

It deliberately does **not** ship `clue-analysis`, `clue-delta`, `clue-extract`, `clue-plan`, or `clue-verify`. That omission is the substance of the change, not an economy: a plugin's components are copied into the user's plugin cache, while Cliewen's five lifecycle skills are committed repository files version-locked to the binary that wrote them ([ADR-011](../../docs/decisions/ADR-011-version-stamping.md), [ADR-022](../../docs/decisions/ADR-022-skill-ownership-marker.md)). A cached copy sits outside `.agents/skills/`, so `checkSkillVersions` cannot see it, and it would carry one version for every repository the user opens — wrong for every repository pinned to a different release.

The plugin tree is hand-authored and lives under `plugins/`, outside both directories the skill generator owns (`.agents/skills/` and `internal/scaffold/templates/skills/`), so `go generate ./internal/skills` neither writes it nor reports it as drift.

## Why

[C-015](../../docs/constraints/C-015-onboarding-under-30-minutes.md) gives the whole journey — install, `clue init`, green `clue validate` — thirty minutes, and [CAP-001](../../docs/capabilities/CAP-001-onboarding/README.md) owns it. M-030 removed the six manual steps for a human at a terminal. It did nothing for the reader who is already sitting in a coding agent, which is the population Cliewen is built for: they still have to leave the session, find the guide, and run a shell command. `/plugin marketplace add cliewen/cliewen` closes that gap without any new hosting, because the manifest is a file in a repository they can already read.

The second reason is the one that needs writing down. A plugin is the obvious place to put five ready-made skills, and doing so would break the drift guarantee that makes `clue validate` trustworthy. The guide page added here exists to say that out loud — a page whose reason to exist is stating what the plugin does *not* install — so that the next person who proposes bundling the lifecycle skills finds the argument instead of rediscovering the failure.

## Decision boundary

One decision, expensive to reverse once published: a marketplace name is registered on other people's machines and a plugin name is a public identifier, so both are an [ADR](../../docs/decisions/README.md) under [C-011](../../docs/constraints/C-011-decision-records-typed.md).

- **ADR-031** — the marketplace ships a bootstrap, never the managed skills; the bootstrap pins no `clue` version; the plugin manifest omits `version` so no second stamp can be forgotten.

This changes neither what `clue validate` checks, nor what a merge accepts, nor the corpus graph, so it is **not** a red-line change under [C-013](../../docs/constraints/C-013-core-changes-need-decision.md). No skill source, lint rule, or corpus convention moves. It is a full change because it makes a decision and adds an acceptance criterion.

## Evidence

The rule this change exists to hold is a negative one, so it gets a criterion rather than a loose guard: **AC-039** on [CAP-001](../../docs/capabilities/CAP-001-onboarding/criteria.md) — the published plugin bootstraps and pins nothing, and the managed lifecycle skills are not among what it installs. Positive and negative tests read the committed plugin tree, so a future edit that quietly adds `clue-delta/` to the plugin, or writes a pinned `CLUE_VERSION` into the bootstrap, fails here rather than in a stranger's session.
