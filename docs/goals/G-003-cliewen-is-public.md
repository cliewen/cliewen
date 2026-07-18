---
id: G-003
type: goal
status: accepted
links: [G-001, G-002]
title: Cliewen is public
---

# G-003 — Cliewen is public

> Accepted 2026-07-18 with [P-003](../plans/P-003-goes-public.md), which carries it as milestones M-008…M-012.

**Who wants it:** the maintainer (2026-07-18), on completing P-002 — the methodology, its distribution story, and its foreign-soil validation are done, and every remaining question about adoption needs strangers.

**Why:** G-001 promises a methodology for engineers and organizations, but a repository no stranger can reach makes that promise unfalsifiable. The foreign-soil trials (P-002/M-007) exhausted what self-directed validation can show; real adopters, real issues, and real contributions require the repository, its releases, and a human-readable guide to be publicly reachable. Going public also completes the install story G-002 started: `go install github.com/cliewen/cliewen/cmd/clue@<version>` should work on a clean machine with no credentials.

**Success looks like:**

- The repository is public; `go install` and anonymous release-asset downloads work without authentication.
- A stranger can learn what Cliewen is and how to practice it from a published guide, without reading the corpus first.
- The front door exists: contribution, conduct, and security-reporting paths route newcomers into the change loop rather than around it.

**Door, not design:** what "readiness" includes — the guide's shape, the community files, the release cut, the flip mechanics — is the campaign's material ([P-003](../plans/P-003-goes-public.md)).
