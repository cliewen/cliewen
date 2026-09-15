---
id: G-018
type: goal
status: accepted
links: [G-013, VIS-001]
title: Enabling identity coordination accounts for the identities already claimed on open branches
---

# G-018 — Enabling identity coordination accounts for the identities already claimed on open branches

**Who wants it:** any maintainer turning on Git-coordinated allocation in a repository that already has work in flight, and every contributor who branches afterwards (2026-09-12, found while enabling coordination in this repository).

**Why:** `clue id coordinate` seeds the allocator branch from the local ledger, and the local ledger is the one on the branch it is run from. Identities claimed on unmerged branches are invisible to it, because a branch's ledger additions reach `main` only at merge. So the allocator starts life believing every number above the last merged claim is free, and hands out identities that open work already holds.

This is not a corner case. It happened three times in one session in this repository: `CH-181`, then `PDR-058`, then `AC-190` through `AC-192`, each already held by a single open pull request, each detected only because someone went looking at that branch's ledger. Nothing warned, and `clue validate` cannot catch it either — the colliding artifact lives on a branch it never sees.

The consequence is worst where the feature is most valuable. A repository enables coordination precisely when several people are working at once, which is exactly when unmerged claims exist. The safety the command offers from that moment forward is real, and the gap it leaves behind is invisible until two artifacts carrying one identity meet at a merge.

The manual recovery is to notice, then allocate again until the number is above every in-flight claim — which requires knowing that the in-flight claims exist, and is precisely the discovery work coordination was meant to remove.

**Success looks like:**

- Turning on coordination accounts for identities claimed on branches that have not merged, or says clearly that it cannot see them and what the maintainer should do about it.
- A contributor who branches after coordination is enabled receives an identity no open branch already holds.
- A collision that does slip through is reported by something, rather than found by a human reading another branch's ledger.
- The remedy never involves hand-editing the append-only ledger or renumbering an artifact that already carries an identity.
- The command stays honest about what it checked: an unreachable branch or remote is reported as unknown, never assumed clear.
