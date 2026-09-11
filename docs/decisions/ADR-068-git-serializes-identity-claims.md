---
id: ADR-068
type: decision
status: verified
links: [G-013, CAP-010, ADR-048, ADR-044, C-013]
title: A Git remote ref serializes identity claims and the ledger is append-only
binds: adopter
author: agent
accepted-by: Flemming N. Larsen (2026-09-11)
---

# ADR-068 — Git serializes identity claims

## Context and problem statement

ADR-048 made the checked-in ledger remember issued identities, but its load-increment-save transaction is isolated to one checkout. Separate clones and worktrees starting from the same commit can therefore reserve the same number, and identical YAML additions can merge without a textual conflict.

The repair must retain sequential identifiers, repository-owned state, host independence, and [the deterministic judge's repository-state boundary](../architecture/README.md). A hosted allocator service would add vendor state outside the repository; random identifiers would discard the readable numeric contract rather than coordinate it.

## Decision outcome

**A repository may explicitly coordinate numeric allocation through the permanent `refs/heads/clue/id-allocator` branch on a configured Git remote.** Allocation commits form a fast-forward-only compare-and-swap sequence. A losing concurrent push refetches and retries; an unreachable, unauthorized, missing, or malformed remote state fails without local fallback. The ref stores permanent numeric claims only and is protected against force-push and deletion.

**The checked-in ledger becomes a versioned append-only event log merged by Git's built-in union driver.** It remains the lifecycle source read by validation, folds repeated events idempotently, advances state monotonically, and rejects conflicting identity metadata. The remote journal never imports feature-branch `live` or `retired` events.

Existing repositories migrate without renumbering and remain in local mode until a team enables coordination. Batch claims and read-only synchronization are the supported path for contributors who cannot push the allocator ref. The cross-cutting transaction and recovery flow is recorded in [the design overview](../design/README.md); the CLI and remote boundary are recorded in [the architecture overview](../architecture/README.md).

**Carrier:** `internal/idalloc`, the version-two ledger and migration, `clue id coordinate|next|sync|live`, the union `.gitattributes` rule, `internal/skills/source/skills/clue-delta.md.tmpl`, `internal/scaffold/templates/.gitattributes`, and the team-allocation guide.
