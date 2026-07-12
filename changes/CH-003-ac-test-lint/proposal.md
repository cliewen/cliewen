---
id: CH-003
type: change
status: open
links: [P-001, M-002, G-001]
title: AC↔test traceability lint — completing M-002
---

# CH-003 — AC↔test traceability lint

## What

The missing half of the machine-checkable thread: `clue validate` learns to fail when an acceptance criterion in an **active** `criteria.md` has no test, or when a test references an AC that does not exist anywhere. This is the last link of the red thread (G → CAP → AC → test) becoming lintable.

## Why

Serves P-001/M-002, whose exit criterion reads: "linter fails the build when an AC lacks a test or a test lacks an AC." CH-002 built the judge; this change gives it the thread's final edge. On merge, M-002 is done.

## Design decisions inside

- **Test-reference convention** (ADR-005, `author: agent`, `inferred`): a Go test references AC-xxx by carrying `AC` + the digits in its function name (`TestAC004_ValidCorpus…`). Names are the one place Go tooling, humans and grep all already look.
- **Draft exemption**: ACs in `criteria.md` with `status: draft` require no tests yet — that is what draft means (CAP-001 today). Promotion of criteria to `active` arms the contract.
- **Deferred (door, not gap)**: enforcement of the positive+negative *pair* needs a labeling convention that hasn't earned its complexity yet; today the rule requires at least one test per AC.

## Scope

New rule in `internal/corpus`, scanning `*_test.go` outside `docs/`/`changes/`; new ACs (AC-009, AC-010) on CAP-002 with their tests; ADR-005; M-002 bookkeeping in the digest.
