---
id: G-011
type: goal
status: proposed
links: [G-003, G-007]
title: A repository states whether it is Cliewen's source or an adopter
---

# G-011 — A repository states whether it is Cliewen's source or an adopter

**Who wants it:** agents and contributors working in Cliewen's own repository and in the repositories that have adopted it (2026-09-01), found while reviewing whether the methodology's documentation output is proportionate. This repository's corpus had grown to a state where process history outweighed the durable system picture by a wide margin, and the largest single contributor was material about the methodology itself — work no adopter would ever produce.

**Why:** nothing in a repository says which kind it is. Cliewen's own repository and every adopter both carry `docs/`, `.clue/`, and `.agents/skills/`; the only structural difference is the incidental presence of `cmd/clue`. [ADR-013](../decisions/ADR-013-ships-generic-vs-repo-local.md) draws the boundary between what ships and what stays local, and [G-007](G-007-changelog-scope-names-the-shipped-surface.md) already had to repair one rule that read naturally as either side, but the boundary itself is enforced by discipline alone: a rule can be recorded in this repository's corpus and never reach the shipped skills and templates, or reach them and never be recorded, and no check fails either way.

An agent inherits that ambiguity. It reads a routing hub whose shared rules and repository-local rules sit in one document, and it cannot tell which of them bind where it is working. The cost compounds in both directions — this repository accumulates bookkeeping inside the same corpus shape it asks adopters to keep, and adopters inherit obligations that were only ever about Cliewen's own development.

**Success looks like:**

- A repository declares its role as an observable fact, so tooling and skills branch on it rather than inferring it from a directory listing.
- A repository that predates the declaration keeps working and is never blocked for lacking it.
- A rule that binds adopter behaviour is checkably present on the surface adopters actually receive.
- This repository's own local rules are separated from the shared routing text, so a reader can see which is which without reasoning from a decision record.
