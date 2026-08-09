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

`cliewen/cliewen` release preparation had inherited the full Cliewen loop even though its only judgement is the chosen version and user-facing notes. It also treated a tag as a release, so a failed publication consumed a version that users could not install. An adopter's release is its own product operation and is outside Cliewen's method, skills, CLI, scaffold, and validation.

## Decision outcome

**Only `cliewen/cliewen` uses a short, ordinary PR for a `clue` release. Human merge authorizes that publication without a CH identity, workspace, digest, acceptance brief, or agentic review loop.**

- The version still comes from the canonical skill frontmatter stamp, and the merged `CHANGELOG.md` section remains the release body.
- A tag with no GitHub Release object is unpublished. The next merge retargets that tag to its repaired `main` commit, reruns the release gates, and publishes the same version.
- The existence of any GitHub Release object, including draft or prerelease, is the public boundary: its tag is never moved or republished.
- Tagging and publication are serialized, and a release runner checks that the tag still names its own commit before publishing. A stale runner cannot publish an older commit after a retry.

**Carrier:** this repository's `AGENTS.md`, `CONTRIBUTING.md`, and PR template; its CI scope classifier, release gates, `tag-on-merge.yml`, and `release.yml`; CAP-004's release-pipeline design; and their focused tests. No distributed skill, scaffold, guide, or `clue` command carries this process.

### Rejected: abandon a version after any failed run

A tag without a release is not a release. Consuming a patch number for an internal failure makes recovery slower and makes the public version sequence describe CI accidents rather than published software.

### Rejected: move a tag after a GitHub Release exists

Even an incomplete public release has escaped the private publication phase. Rewriting it changes public history and is less safe than issuing a new version.

## Consequences

Cutting this repository's `clue` release stays reviewed and human-merged, but it is deliberately short. A failed build before GitHub creates a release page is repaired on `main` and retried under the same version; a published version remains immutable. An adopter can release its own software in any way it chooses, with no Cliewen release rule or tooling involved.
