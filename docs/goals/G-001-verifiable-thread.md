---
id: G-001
type: goal
status: accepted
links: [PDR-021]
title: A verifiable thread from goal to acceptance evidence
---

# G-001 — A verifiable thread from goal to acceptance evidence

## Who wants it

Engineers and organizations doing agent-driven development who must answer — to themselves, to reviewers, to auditors — *why does the system look like this?* "The agent decided" is not an acceptable answer.

## Why

Agent-written code outpaces human ability to keep documentation, decisions and acceptance evidence honest by hand. Existing SDD frameworks document the *change* and let the system's durable truth rot. Cliewen makes the documentation corpus the system-of-record and mechanically enforces the chain goal → capability → acceptance criterion → acceptance evidence, so the thread from intention to verified merge never breaks and never falls behind. Machine-proven criteria end in supported, classified test references; genuine Human-class criteria end in the pull request acceptance brief that the human merge judges.

## Success looks like

- Reachable Git history in any Cliewen repo is a complete provenance archive for a full change, with `git log docs/` exposing the durable corpus history.
- A build fails when an active machine-proven acceptance criterion lacks its required test reference or a test references no live criterion; Human-class proof remains explicit in the acceptance brief rather than being disguised as code.
- A new user reaches a green `clue validate` in under 30 minutes ([CAP-001](../capabilities/CAP-001-onboarding/README.md)).
