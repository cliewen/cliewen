---
id: CH-137
type: change
status: open
links: [G-001, G-002, P-012, P-013, ADR-012, ADR-028, ADR-030, ADR-039]
title: Release Cliewen 0.14.0
---

# CH-137 — Release Cliewen 0.14.0

## What

Prepare and publish the next minor Cliewen release as v0.14.0. The release turns the current `[Unreleased]` user-facing digest into the v0.14.0 release body, restores an empty `[Unreleased]` section, advances the generated-skill version stamp, and records the matching generated-carrier digests in the migration manifest.

## Why

The accepted work since v0.13.0 adds and changes user-facing commands, migration contracts, validation behavior, generated skills, and the hosted change workflow. Keeping that accumulated work unreleased prevents adopters from receiving the corresponding binary and coordinated migration path.

## How

- Promote the existing `[Unreleased]` changelog content to a dated 0.14.0 section without rewriting its user-facing history.
- Bump the canonical generated-skill frontmatter stamp and regenerate generated carriers.
- Add the generated carrier digests for 0.14.0 to `internal/migrate/migrate.go` so an adopter at this release can migrate forward without a false local-edit report.
- Run the release-specific local installation and migration preflight, then the repository's full verification gates. The release PR's merge triggers tagging and publication through the existing workflows.

## Plan

This change is explicitly plan-less. It publishes the completed work in P-012 and the completed milestones already merged under P-013, but it does not implement either active P-013 milestone; treating release packaging as one of those milestones would change their declared scope.
