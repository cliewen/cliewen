---
id: CH-108
type: change-proposal
links: [P-010, M-046, CAP-001, CAP-004, ADR-021, ADR-022, ADR-031, ADR-039]
title: Upgrade routing reaches repositories that have already onboarded
---

# CH-108 — Upgrade routing reaches repositories that have already onboarded

## What

Add a generated, repository-managed upgrade skill with a distinct entry point for already-onboarded adopters. It will check whether a release is available, direct the agent to read that release's migration notes, ask the human whether to upgrade now, and, only after an affirmative answer, carry out the reviewed repository upgrade without merging it.

Extend the managed-carrier migration so it can add that newly introduced whole skill directory when other remaining managed carriers identify the installed release. Its preview will state the target release whose bytes it writes and the recognized sibling release that makes that addition safe. A missing or unrecognized managed set remains a finding and no write occurs.

Record and distribute the resulting six-skill topology coherently: canonical source, generated agent and embedded-template trees, migration manifest, capability criteria and guidance, decisions that previously fixed the set at five, release notes, and regression evidence.

## Why

An adopter who installed Cliewen before this release has its lifecycle skills but no entry point for staying current. `clue latest` can now say that a newer release exists, yet it cannot route an agent through the human-authorized, reviewed migration that moves the coordinated repository state. Treating a new managed skill as merely missing also blocks the complete migration exactly when an older, recognized release proves it is safe to add.

This serves P-010 / M-046 without turning `clue migrate` into a self-updater or granting an agent authority to alter or merge a repository without the human's decision and merge.
