---
id: CH-002
type: change
status: open
links: [P-001, M-002, G-001]
title: Minimal clue validate + CI wall
---

# CH-002 — Minimal `clue validate` + CI wall

## What

The first Go code: a `clue` binary whose only command is `validate` —
the deterministic judge for the corpus rules CH-001 established by
hand. Plus the CI workflow that runs the same binary on every PR and
push to `main`.

## Why

Serves P-001/M-002 (the thread is machine-checkable). CH-001 verified
the corpus with a throwaway script; that judgment must live in the
repo, deterministic and identical locally and in CI, or the wall is a
suggestion.

## Scope

`clue validate` checks: frontmatter presence and core fields, unique
IDs, link resolution (including M-xxx inside plan files), status
vocabularies per type, README.md in every `/docs` folder, index-block
integrity, and `--forbid-changes` (the digest-before-merge gate, used
by CI). New capability CAP-002 with Gherkin ACs; the yaml.v3 dependency
decision recorded as ADR-003 (`author: agent`, `inferred`).

**Not in scope:** the AC↔test traceability rule (the rest of M-002 —
needs the test-tag convention, next change), `init`/`scaffold`/`status`
commands, completed-plan immutability rule (needs git diff context).
M-002 moves to `in-progress`, not `done`.

## Operational note

Branch protection is deferred while the repo is private on the free
plan (decision: Flemming, 2026-07-12). The CI check is advisory until
then; merges happen only when it is green.
