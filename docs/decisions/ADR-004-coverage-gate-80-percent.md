---
id: ADR-004
type: decision
status: verified
links: [QS-001]
title: Default test-coverage gate at 80% total statement coverage
author: human
accepted-by: Flemming N. Larsen (2026-07-12)
---

# ADR-004 — Default test-coverage gate: 80% total statement coverage

## Context and problem statement

The AC↔test contract is the binding behavioral gate, but code coverage is a separate signal external quality tooling (e.g. SonarCloud) typically demands — commonly at 80%. Should Cliewen repos carry a coverage gate by default, and at what level?

## Decision drivers

- Repos ignoring coverage get hit later by external CI tooling defaults; retrofitting coverage is far more expensive than holding it.
- A gate must be a tripwire, not a target: high enough to catch untested code landing, low enough not to invite meaning-free tests written to touch lines (the Goodhart failure mode of 100%).
- ACs remain first priority: coverage never substitutes for the AC↔test contract, it backstops the code paths ACs don't reach.

## Considered options

1. **No coverage gate** — ACs only.
2. **80% total statement coverage** (the common external-tooling default).
3. **100% / per-package thresholds.**

## Decision outcome

**80% total statement coverage, enforced in CI, as the default.** Total rather than per-package: thin entry points (process-exit switches in `main`) legitimately resist in-process testing and would force contortions under per-package gates. Code smells and craft quality are not gated by `clue` or CI thresholds — they are the agents' job in each PR and the human's in review (consistent with dot-principles staying orthogonal, Foundation §6).

### Carrier (added 2026-07-12, per the carrier rule)

The CI workflow template `clue init` scaffolds: the coverage-gate step ships with the default workflow, adapted per language profile. This repo's own `.github/workflows/ci.yml` is the template source.

### Rejected: no gate

Fails the external-tooling reality and lets untested plumbing accumulate silently between AC-covered paths.

### Rejected: 100% / per-package

A Goodhart magnet: agents satisfy the number with tests that touch lines without verifying meaning — precisely the cheating the methodology exists to prevent.
