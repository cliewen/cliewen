---
id: PDR-005
type: decision
status: verified
links: [P-002]
title: Validation requires foreign soil — trials on external repos, as findings not adoptions
author: agent
accepted-by: Flemming N. Larsen (2026-07-15, PR #12 approval)
---

# PDR-005 — Validation requires foreign soil

## Context and problem statement

Validation on repositories prepared by the maintainer cannot expose adoption failures caused by different documentation, language, test, or instruction conventions. Honest evidence requires ground the methodology did not prepare.

## Decision outcome

**The plan must trial the skills on at least two deliberately foreign open-source repositories, and each trial produces findings rather than an adoption.** A human selects repositories with no shared maintainer and no methodology preparation. Each trial ends in an `AN-xxx` findings document and must produce at least one methodology adjustment or an explicit finding that the expected adjustment was not needed.

Trials are read-and-apply experiments: they do not create pull requests against the foreign repositories or new extraction mappings. A demand for either becomes a finding for a later decision. The M-007 plan row and the analysis skill's findings-document contract carry this repository-local rule.
