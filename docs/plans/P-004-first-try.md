---
id: P-004
type: plan
status: completed
links: [G-001]
title: Cliewen earns the first try — value, safe proof, enforced adoption
---

# P-004 — Cliewen earns the first try

> **Completed 2026-07-23** — all milestones done; frozen immutable. The public front door, safe demo, release-binary installation, enforced-CI guidance, adoption-fit guidance, and operating boundaries are all in place. No successor is designated.

P-001 established the verifiable thread, P-002 made it distributable, and P-003 made the repository and newcomer guide public. This campaign makes that public path earn a newcomer's first try: state the value before the vocabulary, demonstrate a real validation failure without risking an existing repository, make the CI safeguard actionable, and state the product's operating boundaries honestly. It serves [G-001](../goals/G-001-verifiable-thread.md) without changing the validator or [QS-002](../quality/QS-002-onboarding-under-30-minutes.md).

The campaign is deliberately sequential. Each milestone is a separate Cliewen change rooted at accepted `main`, published as a ready pull request, and accepted by human merge before the next begins. The scope decision and deferrals are recorded in the [decision log](../decisions/log.md) row dated 2026-07-22.

## Milestones

| ID | Milestone (exit criterion) | Status | Evidence |
|---|---|---|---|
| M-013 | **Value, installation, and safe proof are concrete**: the public front door names the problem, audience, outcome, and Cliewen/`clue`/corpus distinction before methodology terminology; release binaries are the primary installation route with Windows, macOS, and Linux asset, checksum, `PATH`, and version-result instructions while source installation is secondary; Getting Started uses a disposable Git repository, shows generated output and ownership boundaries, and turns a draft criterion active to produce the real missing-test diagnostic before cleanup; the change-loop narrative links the real AC↔test implementation PR; the strict guide build passes and the demo is rehearsed from a clean temporary repository | `done` | CH-046: front doors and concise direct-download instructions rewritten; v0.5.1 Windows asset checksum/version and fresh init verified; draft→active missing-test failure and recovery rehearsed in a disposable repository; `cliewen.dev` root-base repair and strict VitePress build passed; PR #2 linked as the real AC↔test change |
| M-014 | **The wall and minimum practice are actionable**: a dedicated CI guide gives exact pinned-binary vendoring, checksum, armed/unarmed, and `--forbid-changes` steps; GitHub ruleset instructions require the stable validation check on current `main`, protect the branch from bypass, force-push, and deletion, and include a disposable failing-PR probe that demonstrates merge blocking; other forges get the equivalent enforcement contract without unsupported UI instructions; the adoption guide states the minimum goal→capability→criterion→positive/negative-evidence path, when the wider taxonomy becomes useful, and when Cliewen is a poor fit | `done` | CH-047: v0.5.1 init/workflow pin and Linux CI asset checksum verified; dedicated wall guide covers vendoring, armed state, rulesets, forge-neutral enforcement, and cleanup; ordinary validation passed while `--forbid-changes` rejected the rehearsed disposable workspace; strict guide build passed; adoption guidance now states the minimum thread, wider-taxonomy triggers, and poor-fit cases |
| M-015 | **Boundaries and continued operation are explicit**: a support/operations page distinguishes shipped operating-system artifacts, verified test harvesters, installed agent-skill layouts, the GitHub workflow, and repository-local validation from broader methodology intent; unsupported “any framework” wording is corrected; coordinated binary/skill/CI upgrades, version drift, skipped init files, extraction rollback, unexpected validation failures, uninstallation, and adoption rollback have safe procedures; existing foreign-soil trials are linked only as trials; every guide page ends in one primary next action; guide, corpus, and repository verification pass | `done` | CH-049: operations guide names the release, Go/JVM harvester, generated-layout, GitHub CI, and repository-local boundaries; gives coordinated upgrade and recovery procedures; links trials only as trials; and adds AC-036 positive/negative guide checks alongside the strict guide build and repository verification. |

## Explicitly out of this campaign

Changes to `clue`, validation semantics, positive/negative-pair machine enforcement, skills, scaffold templates, or QS-002; installer scripts; Homebrew, Scoop, Winget, or another package-manager channel; a separate example repository; terminal recordings or GIFs; claims of external adoption; testimonials; and productivity or review-time measurements. Any of these needs its own evidence and accepted scope.

## Mutation rules (lintable)

Status and evidence fields in the milestone table may mutate in an implementing change's merge digest. Everything else in this file changes only through a declared plan revision backed by a decision record routed by reversal cost ([C-011](../constraints/C-011-decision-records-typed.md)). Plan adjustments are decisions.
