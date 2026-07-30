---
id: CH-086
type: change
status: open
links: [P-009, AN-013]
title: Restore P-009's deferred distributed-work boundaries
---

# CH-086 — Restore P-009's deferred distributed-work boundaries

## What

Revise [P-009](../../docs/plans/P-009-migration-readiness.md) without disturbing the six migration-readiness milestones accepted through CH-085. Add M-043 for durable dependent-change base and authorization data and M-044 for named but locally unproven foreign acceptance evidence, recovering the two unique load-bearing candidates preserved in the unmerged PR #89 history at `fefbacc3a5a96a5048723af1ecc0b738c398af2d`. The revision keeps M-037…M-042 first, removes only these two candidates from the plan's future-door exclusion, and records the changed scope and order as a reversible decision-log row.

## Why

[AN-013](../../docs/analysis/AN-013-distributed-work-and-evidence-boundaries.md) reproduced both gaps and named P-009 as their consumer. Concurrent plan-opening branches selected the same unreserved durable identities; CH-085 merged first with a migration-focused campaign that retained these candidates only as future doors. The second branch could not merge unchanged, but its unique meaning remains wanted and recoverable without duplicating CH-085's CI-wall or corpus-upgrade milestones.

## Decision boundary

This change decides only that P-009 continues from migration readiness into the two already-analyzed distributed-work boundaries, with M-043 before M-044 because the base-and-authorization rule is process-only while the foreign-evidence form changes corpus and validator contracts. It does not decide either milestone's PDR or ADR, weaken the human merge boundary, treat forge state as meaning, add network resolution to the judge, change acceptance criteria, or implement product behavior. Human merge accepts the semantic plan revision.
