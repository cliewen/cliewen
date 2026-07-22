---
id: CH-049
type: change
status: open
links: [P-004, M-015, CAP-001, QS-002, C-012]
title: Make Cliewen's operating boundaries and recovery paths explicit
---

# CH-049 — Make Cliewen's operating boundaries and recovery paths explicit

## What

Complete P-004/M-015 with a support and operations guide for people who have already completed the first trial. It will distinguish the released operating-system binaries, supported harvested-test conventions, installed agent-skill layouts, generated GitHub workflow, and repository-local validation from the broader methodology. It will state the present framework and repository limits honestly; document safe binary, skill, workflow, and vendored-CI upgrades; and give recovery paths for drift, skipped init files, extraction rollback, unexpected validation failures, uninstallation, and adoption rollback.

The guide will link existing foreign-soil work only as trials, not external adoption evidence. Every public-guide page will end with one concrete primary next action. The durable onboarding design, user-facing changelog, and P-004 milestone bookkeeping will record the completed guide boundary.

## Why

The current guide explains a first safe trial and CI enforcement, but it does not yet tell an adopter what Cliewen actually ships, what it verifies, or how to safely recover when the setup changes or the method does not fit. That gap can turn a successful trial into avoidable operational risk or make the tool appear to promise support it does not provide.

## Decision boundary

This change implements the already accepted M-015 scope. It does not change `clue`, validation semantics, test-harvester support, skills, scaffolding, generated workflows, CI policy, QS-002, or any external repository. Instructions may describe only behavior verified from this repository and released artifacts. Any need to broaden supported frameworks, change a shipped contract, claim external adoption, or make a new architecture or process decision is an open question and stops the change.
