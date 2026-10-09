---
id: PDR-069
type: decision
status: inferred
links: [G-016, CAP-001, CAP-012, PDR-067, PDR-068]
title: Local acceptance is the default and PR requires an explicit choice
author: agent
accepted-by: []
binds: adopter
---

# Local acceptance is the default

## Context

The owner requested one default for both new and existing adopters, with no compatibility layer for an implicit PR convention. Preserving that convention made an absent file mean different things depending on prior repository state.

## Decision

Without `.clue/acceptance.yaml`, acceptance is local on `main`. Init writes local policy regardless of adoption history; `--acceptance=pr` selects PR. Existing explicit policy is preserved, and conflicting init options fail before writes. The shared reader remains strict. Migration validates existing policy but does not create it. Remove MIG-021, prior-adoption detection and the historical PR fallback.

Base and candidate must resolve to the same local policy. An absent file and explicit local/main configuration are equivalent. An accepted PR policy cannot be overridden by a candidate to accept itself locally. Human confirmation, recorded evidence and identity allocation retain their existing meaning. Source repositories require explicit PR policy and remain PR-only.

This supersedes PDR-068's implicit PR preservation and PDR-067's initial opt-in requirement. The guide uses `/acceptance`, with no compatibility page at `/local-acceptance`. Shipped carriers are the shared policy reader, init, migration preconditions, acceptance preflight, `internal/skills/source/shared/review-boundary.md.tmpl` and `internal/scaffold/templates/AGENTS.md`.

## Risk and authorization

The riskiest assumption is that existing adopters have recorded a required PR workflow explicitly. A repository without policy now uses local acceptance. Such repositories must write `mode: pr` and their integration branch if PR acceptance is required; branch protection and repository permissions continue to apply. The owner explicitly chose this change without backward compatibility and authorized the direct route despite the tracked recommendation. No change identity or tracked workspace is created.
