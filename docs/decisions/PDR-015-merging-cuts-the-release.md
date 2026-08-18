---
id: PDR-015
type: decision
status: verified
links: [CAP-004, ADR-011, ADR-012, C-012]
title: This repository's release PR cuts a recoverable clue release
author: agent
accepted-by: Flemming N. Larsen (2026-08-02, conversation); Flemming N. Larsen (2026-08-09, conversation)
---

# PDR-015 — This repository's release PR cuts a recoverable clue release

## Context and problem statement

This repository's release preparation needed only a chosen version and user-facing notes, while treating a tag as a release made failed publication consume a version. Adopter releases are their own product operations and are outside Cliewen's method and tooling.

## Decision outcome

**Only `cliewen/cliewen` uses a short ordinary release PR without CH identity, workspace, digest, acceptance brief, or agentic review; human merge authorizes publication.** The canonical skill stamp supplies the version and the merged `CHANGELOG.md` section supplies the release body.

A tag without a GitHub Release object may be retargeted to the repaired `main` commit and retried under the same version. Once any release object exists, including draft or prerelease, its tag is immutable. Tagging and publication are serialized, and the runner verifies that the tag names its own commit before publishing.

The repository's AGENTS.md, CONTRIBUTING.md, PR template, release workflows, CAP-004, and focused tests carry this local specialization; no adopter-facing skill, scaffold, guide, or `clue` command carries it.
