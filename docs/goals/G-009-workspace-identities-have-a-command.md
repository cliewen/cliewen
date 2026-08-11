---
id: G-009
type: goal
status: proposed
links: [G-001]
title: A change workspace's identities have a command, not a hand-edited ledger
---

# G-009 — A change workspace's identities have a command, not a hand-edited ledger

**Who wants it:** agents and contributors running the full change loop (2026-08-11), found while CH-150 opened and digested a workspace.

**Why:** `clue validate` requires every workspace artifact's identity to be present in `.clue/id-ledger.yaml`, and `clue id` cannot put it there. `clue id next` allocates a numeric identity and reserves it; `clue id live` promotes a reserved identity and refuses one it has never seen, which is every opaque workspace identity — the proposal, tasks, and open-questions records named after the change. There is no subcommand that retires an identity either, so the states a digested workspace leaves behind are reached the same way. The ledger is therefore hand-edited twice in the course of one ordinary change: once to register three opaque identities so the proposal commit validates, and once to retire them and the change identity so the digest matches what earlier changes left behind.

Hand-editing the ledger is the part of the loop with no machine and no guard. A wrong or reused identity is exactly what the ledger exists to prevent and what the judge treats as blocking, and the file's own integrity is being maintained by hand at the two moments a change is most likely to be rushed. The retirement half is worse than unhelpful: nothing checks it, so a digest that simply forgets leaves live entries pointing at deleted files, and the corpus looks the same either way.

**Success looks like:**

- Registering and retiring a workspace's identities is a command, so the ledger is not edited by hand in the ordinary loop.
- The states a digested change leaves in the ledger are either produced by that command or checked by the judge, so a forgotten retirement is visible.
