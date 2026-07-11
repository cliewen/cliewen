---
id: G-001
type: goal
status: accepted
links: []
title: A verifiable thread from goal to test
---

# G-001 — A verifiable thread from goal to test

## Who wants it

Engineers and organizations doing agent-driven development who must
answer — to themselves, to reviewers, to auditors — *why does the
system look like this?* "The agent decided" is not an acceptable
answer.

## Why

Agent-written code outpaces human ability to keep documentation,
decisions and tests honest by hand. Existing SDD frameworks document
the *change* and let the system's durable truth rot. Cliewen makes the
documentation corpus the system-of-record and mechanically enforces the
chain goal → capability → acceptance criterion → test, so the thread
from intention to verified merge never breaks and never falls behind.

## Success looks like

- `git log docs/` in any Cliewen repo is a complete provenance archive.
- A build fails when an acceptance criterion lacks a test or a test
  lacks a criterion.
- A new user reaches a green `clue validate` in under 30 minutes
  ([CAP-001](../capabilities/CAP-001-onboarding/README.md)).
