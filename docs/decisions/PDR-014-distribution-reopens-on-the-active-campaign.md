---
id: PDR-014
type: decision
status: verified
links: [P-007, CAP-001, C-015, ADR-030, PDR-008]
title: Installation distribution reopens on the active campaign, not a successor plan
author: agent
accepted-by: Flemming N. Larsen (2026-08-02, conversation)
---

# PDR-014 — Installation distribution reopens on the active campaign

## Context and problem statement

Installation distribution — installer scripts and package-manager channels alike — has been parked three times. The [decision log](PDR-046-decisions-route-by-subject.md) row of 2026-07-22 placed both outside the first-try campaign "until supported by their own evidence and scope", and P-004, P-005 and P-006 each restated them in their exclusion lists. The parking was conditional, and the condition was evidence — but nothing in the corpus says who decides the condition is met, or where the work lands when it is.

Two questions follow. Does the evidence now exist? And if it does, does distribution wait for a successor plan, or join the campaign already in flight?

## Decision outcome

**The condition is met, and the work joins P-007 as milestones M-030 and M-031 rather than waiting for a successor plan.**

The evidence the parking asked for is three things the corpus did not have in July 2026 and does have now: [C-015](../constraints/C-015-onboarding-under-30-minutes.md) states the onboarding bar as a checkable constraint rather than an aspiration; [CAP-004](../capabilities/CAP-004-ship/README.md) has shipped a released binary on six platforms, so there is something for a package manager to carry; and the install procedure is written down in `guide/getting-started.md`, so its length can be counted against the bar it must clear. Six manual steps — read a table, download, verify a checksum by hand, rename, `chmod`, edit `PATH` — stand at the front of a 30-minute journey. That is the demonstrated gap the parking wanted and the earlier campaigns lacked.

It joins P-007 for a reason specific to that plan: P-007's exclusion list, alone among the campaigns, never names installation. It excludes adopter hardening a real pilot must price first, the production feedback loop, and foreign documentation kinds. Distribution was excluded by P-004 through P-006 and then simply not carried forward, so appending it retracts nothing the active campaign promised. A successor plan would cost the six sequential red-line milestones standing in front of it, for a gap that is measured, bounded, and touches no core element.

The two milestones stay ordered rather than parallel: M-031's marketplace bootstrap consists of the install commands M-030 creates, so shipping it first would produce a bootstrap that hands the user the manual table.

**Carrier:** the milestone rows in [P-007](../plans/P-007-core-hardening.md) and its mutation rules; this record is the revision's backing under [PDR-008](PDR-008-plan-revisions-may-ride.md).

### Rejected: a successor plan P-008

Cleanest campaign hygiene, and the corpus's stated pattern for a deferred area is a goal proposed in its own change. But that pattern is written for scope whose shape is unknown — P-007 recommends it for foreign documentation kinds precisely because "the plan that pursues it should follow a real adopter's mandate rather than precede one". Distribution is the opposite case: the gap is already measured against an existing constraint, and the carriers that hold it — CAP-001 and CAP-004 — are both active. A plan exists to discover a campaign's shape; here the shape is known, and a plan would be ceremony that delays a fix to the first thing every newcomer meets.

### Rejected: a new goal alongside the plan

[G-001](../goals/G-001-verifiable-thread.md) already names the 30-minute journey in its success criteria, and CAP-001 already owns onboarding as a capability. A distribution goal would restate an accepted goal's own success condition as a second goal, splitting one intention across two records with nothing to keep them agreeing.

### Rejected: reopening by log row alone

The parking was a log row, so reversal is cheap and a row would be correctly typed for that half. But the decision being made is larger than un-parking: it revises an active plan, and [PDR-008](PDR-008-plan-revisions-may-ride.md) requires a plan revision to be backed by a typed record and called out for deliberate approval. A row carries neither the rejected alternatives above nor the reviewer's invitation to object to the revision specifically.
