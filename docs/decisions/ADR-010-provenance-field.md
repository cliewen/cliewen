---
id: ADR-010
type: decision
status: verified
links: [AN-002, ADR-008, CAP-003]
title: Extracted artifacts carry a provenance field, born inferred
author: agent
accepted-by: Flemming N. Larsen (2026-07-12)
---

# ADR-010 — The provenance field

## Context and problem statement

Decisions already carry two-tier provenance in their `status` (`inferred` → `verified`). Extraction ([ADR-008](ADR-008-extraction-is-a-skill.md)) mints dozens of *non-decision* artifacts in one PR — goals, plans, capabilities, criteria — all agent-reconstructed and none yet human truth. Their `status` vocabularies are lifecycle words (`draft`/`active`/…) and must stay that way: an extracted capability with passing tests is genuinely `active`, and marking it `draft` would exempt its criteria from the AC↔test contract (AC-009) — the opposite of what extraction needs.

## Decision outcome

**Provenance and lifecycle are separate axes.** Any artifact may carry an optional frontmatter field `provenance: inferred | verified`; extraction writes `provenance: inferred` on everything it mints. An absent field means human-authored — the existing corpus needs no touch. Decisions keep expressing provenance in `status` as today; a decision does not carry the separate field.

Promotion is the human review of the extraction PR: flipping `inferred` → `verified` file-by-file as meaning is confirmed (or in bulk when the reviewer accepts the whole corpus); the PR history is the audit record of who verified. Consumers: `clue validate` lints the vocabulary and reports how many artifacts remain `inferred` — the visible to-do list of unverified meaning.

**Carrier:** the vocabulary lint and inferred count in `clue` (machine); the born-inferred rule in the `clue-extract` skill (agent).

### Rejected: adding `inferred` to every status vocabulary

Conflates lifecycle with provenance: an extracted, tested capability would have to choose between "not yet active" (false, and it disables the test contract) and "active" (losing the unverified marker).

### Rejected: no marker at all — the PR review is enough

The review event is transient; the corpus is the system of record. Six months later nobody can tell which artifacts a human actually read. §7's rule holds — the field earns its place because `clue` reads it.
