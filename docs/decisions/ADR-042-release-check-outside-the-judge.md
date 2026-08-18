---
id: ADR-042
type: decision
status: verified
links: [P-010, CAP-002, CAP-004, ADR-011, ADR-030, ADR-039, ADR-040, C-004, C-013]
title: A release check reaches the network, reports, and writes nothing
author: agent
accepted-by: Flemming N. Larsen (2026-08-03, conversation)
---

# ADR-042 — Release check outside the judge

## Context and problem statement

Users need to know whether their binary or generated repository is behind a release, but network state and self-updating would make the deterministic judge fragile and turn a report into an unreviewed repository mutation.

## Decision outcome

**A release check is a separate read-only command that reaches the network, reports what it found, and never contributes to `clue validate` or a required status check.** It writes no repository file and never replaces its running binary; the upgrade remains the reviewed `clue migrate` path. It selects the one installation route for the current machine and prints the coordinated binary-and-repository sequence.

Offline, timeout, rate-limit, unrecognized, and unusable-cache cases exit successfully, write nothing, and mean “could not tell.” Only plain numeric release tags cross into the printed command or cache; running version stamps may include a prerelease suffix for comparison, and uncomparable or unpublished versions are reported as such. Quiet mode prints one line only when behind, and the bounded cache lives outside the repository. The drift message names how to upgrade and how to remain deliberately on the current release.

**Carrier:** the release-check command, its cache and parser, ADR-011/ADR-030/ADR-039 guidance, and the installation route output.
