---
id: ADR-060
type: decision
status: verified
links: [ADR-039, ADR-052, CAP-001]
title: Migration materializes a missing thin CI caller and reports a competing validation wall
author: agent
accepted-by: Flemming N. Larsen (2026-08-15, conversation)
---

# ADR-060 — Migration materializes a missing thin CI caller and reports a competing validation wall

## Context and problem statement

The thin CI caller is not an optional carrier an adopter has weighed and declined. It first shipped in v0.10.0, so every repository adopted before that release has no `.github/workflows/clue.yml` for a reason that has nothing to do with its CI policy: the file did not exist when it adopted. [ADR-052](ADR-052-missing-optional-carriers-do-not-block-safe-migrations.md) classified that absence as a notice naming `clue init` as the materialization route, and rejected creating the caller during migration because doing so *merely to unblock another migration* would take authority `init` deliberately owns.

That reasoning is correct for the case it addressed and returns the wrong answer for an upgrade. When an adopter is deliberately moving a repository from a release that predates the caller to one that expects it, materializing the caller is not a side effect of unblocking something else — it is the upgrade the adopter asked for. The route ADR-052 named does not serve it well either: `clue init` never overwrites, so it is safe, but it materializes the whole scaffold set, which in an established repository means corpus stubs and a constraint the adopter may deliberately not carry. The adopter is left hand-writing a file whose exact bytes the binary already embeds and already knows how to update once it exists.

The authority objection is weaker than it first appears. The caller's adopter-owned runner, clue-source, and install-directory choices are precisely the fields migration already preserves when it updates a caller that exists. Migration has demonstrably held those choices safely; what it lacked was permission to write the file the first time.

A second gap sits behind the same upgrade. A repository whose pre-caller CI installed and ran `clue validate` in its own workflow keeps that job after the caller arrives. The two walls then judge the same pull request under different rules, and the older one fails work the caller was configured to treat leniently. Migration could see this and said nothing.

## Decision outcome

**Migration creates a missing thin CI caller from its embedded template at the template's default adopter choices.** The absence becomes a planned change rather than a notice, previewed and applied under the same atomic preflight as every other migration write. `clue init` keeps its role as the materializer for a repository that has no corpus yet.

**A competing `clue validate` wall is reported, never rewritten.** When a job in a repository-owned workflow other than the caller runs the installed binary's `clue validate`, migration raises a finding naming the file and job so a human resolves it. A source build such as `go run ./cmd/clue validate` is how a repository dogfoods its own working tree and is not a wall; neither is a reusable workflow definition, which becomes one only where a caller references it. Migration gains no authority over a workflow it did not write.

Everything else ADR-052 decides is unchanged: a missing optional carrier still does not block an independent safe migration, a present caller whose content cannot be safely recognized remains a blocking finding, and ambiguous corpus meaning and locally modified managed carriers still fail without partial writes.

### Rejected: rewrite or delete the competing validation job

The job is repository-owned prose with adopter comments, pinned versions, and local intent that Cliewen never authored. Rewriting it would cross the line migration holds everywhere else — a locally modified carrier is a finding, not a write target — and [ADR-013](ADR-013-ships-generic-vs-repo-local.md) places an adopter's own CI outside the shipped surface. Naming the conflict costs the adopter one decision and costs Cliewen no authority.

### Rejected: keep the notice and make it recommend `init` more loudly

Better wording does not close the gap. The adopter still runs a command that writes files unrelated to the upgrade, or writes the caller by hand. The problem was the missing authority, not the missing explanation.

**Carrier:** `planCaller` plans the creation and the competing-wall detection in `internal/migrate/migrate.go`; AC-124's successor proves materialization and the new competing-wall criterion proves the finding; CAP-001's design states both. The `clue-upgrade` workflow reference carries the reconciliation step for a repository upgrading from before v0.10.0.
