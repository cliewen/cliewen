---
id: PDR-015
type: decision
status: verified
links: [CAP-004, ADR-011, ADR-012, C-012]
title: A short release PR cuts a recoverable release
author: agent
accepted-by: Flemming N. Larsen (2026-08-02, conversation); Flemming N. Larsen (2026-08-09, conversation)
---

# PDR-015 — A short release PR cuts a recoverable release

## Context and problem statement

Release preparation had inherited the full Cliewen loop even though its only judgement is the chosen version and user-facing notes. It also treated a tag as a release, so a failed publication consumed a version that users could not install.

## Decision outcome

**A repository release uses a short, ordinary PR containing only its declared release files. Human merge authorizes publication without a CH identity, workspace, digest, acceptance brief, or agentic review loop.**

- The version still comes from the canonical skill frontmatter stamp, and the merged `CHANGELOG.md` section remains the release body.
- A tag with no GitHub Release object is unpublished. The next merge retargets that tag to its repaired `main` commit, reruns the release gates, and publishes the same version.
- The existence of any GitHub Release object, including draft or prerelease, is the public boundary: its tag is never moved or republished.
- Tagging and publication are serialized, and a release runner checks that the tag still names its own commit before publishing. A stale runner cannot publish an older commit after a retry.

**Carrier:** `AGENTS.md`, `CONTRIBUTING.md`, the canonical change-tier source and generated skills, the PR template, CI scope classifier, release gates, `tag-on-merge.yml`, `release.yml`, CAP-004's release-pipeline design, and their focused tests.

### Rejected: abandon a version after any failed run

A tag without a release is not a release. Consuming a patch number for an internal failure makes recovery slower and makes the public version sequence describe CI accidents rather than published software.

### Rejected: move a tag after a GitHub Release exists

Even an incomplete public release has escaped the private publication phase. Rewriting it changes public history and is less safe than issuing a new version.

## Consequences

Cutting a release stays reviewed and human-merged, but it is deliberately short. A failed build before GitHub creates a release page is repaired on `main` and retried under the same version; a published version remains immutable.
