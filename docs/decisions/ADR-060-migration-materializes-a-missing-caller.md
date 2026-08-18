---
id: ADR-060
type: decision
status: verified
links: [ADR-039, ADR-052, CAP-001]
title: Migration materializes a missing thin CI caller and reports a competing validation wall
author: agent
accepted-by: Flemming N. Larsen (2026-08-15, conversation)
---

# ADR-060 — Migration materializes a missing caller

## Context and problem statement

ADR-052 kept migration from creating a missing thin caller, but an adopter upgrading from before the caller existed cannot receive the reusable validation boundary through the upgrade path. Migration also needs to expose a repository-owned workflow that independently runs `clue validate` rather than overwriting it.

## Decision outcome

**Migration creates a missing thin CI caller from its embedded template using default adopter choices, previewed and applied under the normal atomic preflight; `clue init` remains the materializer for a repository with no corpus.** A repository-owned workflow job that runs the installed binary's `clue validate` is a competing validation wall and becomes a finding naming its file and job. Source builds and reusable workflow definitions are not walls, and migration never rewrites the competing job.

All other ADR-052 guarantees remain: a missing optional carrier does not block an independent safe migration, an unrecognized present caller blocks, and ambiguous corpus meaning or modified managed carriers prevent partial writes.

**Carrier:** `planCaller`, competing-wall detection, AC-124 and its successor criterion, CAP-001 design, and the `clue-upgrade` reconciliation guidance.
