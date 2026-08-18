---
id: ADR-025
type: decision
status: verified
links: [P-005, ADR-010, ADR-002, C-008, PDR-003]
title: One default status lifecycle — draft → active → retired — plus named exceptions
author: agent
accepted-by: Flemming N. Larsen (2026-08-02, conversation)
---

# ADR-025 — One default status lifecycle

## Context and problem statement

The validator carried many near-duplicate status vocabularies, making the few genuinely different lifecycles hard to see and forcing adopters to learn unnecessary tables.

## Decision outcome

**Every artifact type uses the default lifecycle `draft → active → retired` unless a semantic exception is named.** The exceptions are `goal: proposed → accepted`, `plan: draft → active → completed`, `decision: inferred → verified`, `change` and `tasks: open`, and `open-questions: open → resolved`.

ADR-034 corrects `retired` as a resting status: ordinary retired artifacts are deleted and their successor carries `supersedes:`; only criteria tombstones and completed plans remain as files. Architecture and analysis use the default rather than `verified` as a pseudo-provenance state, and constraints gain `draft`. Adopter-defined types use the default under ADR-026.

**Carrier:** `defaultLifecycle` and `statusVocabExceptions` in `internal/corpus/rules.go`, mirrored by `docs/README.md` and the init template. Keeping the per-type table is rejected because it hides the few exceptions; moving decisions to the default is rejected because `inferred → verified` carries human acceptance.
