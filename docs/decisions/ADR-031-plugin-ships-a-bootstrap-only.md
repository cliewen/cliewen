---
id: ADR-031
type: decision
status: verified
links: [CAP-001, CAP-004, C-015, ADR-011, ADR-022, ADR-030]
title: The Claude Code plugin ships a bootstrap, and the managed skills never ride in it
author: agent
accepted-by: Flemming N. Larsen (2026-08-02, conversation)
---

# ADR-031 — The plugin ships a bootstrap, not the skills

## Context and problem statement

Claude Code's marketplace can put Cliewen in an existing agent session, but plugin components live in a per-user cache rather than the repository's stamped `.agents/skills/` tree. Bundling the six managed skills would therefore create an unversioned second instruction set that `clue validate` could not detect or reconcile with the repository's pinned release.

## Decision outcome

**The marketplace ships one plugin with one bootstrap skill whose only job is to install a real `clue` and then stop.** The bootstrap detects the host, uses ADR-030's installation channels, confirms a release binary rather than `dev`, and asks before running `clue init`; installing the plugin is not consent to scaffold a repository.

The bootstrap pins no `clue` version, and the plugin manifest omits `version`, so neither can drift from the release or require a second hand-maintained stamp. The six generated skills remain repository-scoped, stamped, and written only by `clue init` or reviewed migration; the hand-authored plugin tree stays under `plugins/`, outside the generator-owned directories.

**Carrier:** `.claude-plugin/marketplace.json`, `plugins/cliewen/`, `guide/plugin.md`, and AC-039 on [CAP-001](../capabilities/CAP-001-onboarding/criteria.md), whose tests reject a bundled managed skill or pinned plugin version.

## Rejected alternatives

Bundling the managed skills, installing them at project scope, or running `clue init` automatically would either bypass the repository-specific drift check or scaffold a repository without consent. A wider plugin catalog remains additive and deferred until the bootstrap has been used by someone outside its authorship.
