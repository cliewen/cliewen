---
id: G-008
type: goal
status: accepted
links: [G-003, PDR-042, CAP-006]
title: The acceptance-brief gate distinguishes a change's route from its CI scope
---

# G-008 — The acceptance-brief gate distinguishes a change's route from its CI scope

> Accepted 2026-08-13 with PDR-042 and CH-153: semantic routing is now simple or full, and CI check selection no longer supplies the full-loop signal.

**Who wants it:** contributors and agents integrating simple work or taking a full change to a ready pull request, in this repository and in any adopter running the shipped validation workflow (2026-08-11).

**Why:** the gate requiring a completed acceptance brief used to fire on the CI scope classifier's `full` output, which meant "this diff needs the full check suite" rather than "this work chose the full loop". Every non-full change touching `docs/` was therefore asked for a full change's artifact, and the adopter-facing workflow had the same shape. Route and check scope now travel independently: branch history and any complete user-override trailers select full-loop bookkeeping, while changed surfaces select relevant checks.

The accepted answer is that only a chosen full loop owes the acceptance brief. A simple integration carries no Cliewen form, including when the agent originally recommended full and the user chose simple; in that case the complete current-head trailers preserve the authorization and risk without turning relevant surface checks off.

**Success looks like:**

- The gate and the change loop agree about which pull requests owe an acceptance brief, and the corpus records which reading was chosen and why.
- A contributor who is asked for a brief is told a reason that matches the rule they are being held to.
- A simple change that legitimately owes no brief does not fail a required check, in this repository or in an adopter's.
