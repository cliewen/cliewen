---
id: ADR-005
type: decision
status: inferred
links: [CAP-002]
title: Tests reference ACs through their function names
author: agent
---

# ADR-005 — Tests reference ACs through their function names

## Context and problem statement

The AC↔test lint (M-002) needs a deterministic way to know which acceptance criteria a test verifies. Where does the reference live?

## Considered options

1. **Test function names carry the AC digits** — `TestAC004_ValidCorpusHasNoIssues` references AC-004.
2. **Structured comments** — `// verifies: AC-004` above the test.
3. **External mapping file** — a manifest linking test names to AC-IDs.

## Decision outcome

**Function names.** The name is the one identifier Go tooling (`go test -run`), CI logs, failure output and grep all already display — a reference there has maximal visibility and cannot silently detach from its test the way a comment or manifest row can. The convention: a test function matching `Test…AC<digits>…` references `AC-<digits>`; one test may reference one AC, and an AC may be referenced by many tests. Comments would drift; a mapping file is exactly the hand-maintained index the methodology forbids (first artifact type to rot).

### Deferred, not rejected

Enforcing the positive+negative **pair** per AC needs a labeling convention on top (e.g. name suffixes). It is left as a door until real use shows which labeling survives contact with table-driven tests.
