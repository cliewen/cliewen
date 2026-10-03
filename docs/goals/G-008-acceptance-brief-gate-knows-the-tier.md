---
id: G-008
type: goal
status: accepted
links: [G-003, PDR-042, CAP-006]
title: The acceptance-brief gate distinguishes a change's route from its CI scope
---

# G-008 — The acceptance-brief gate distinguishes a change's route from its CI scope

> Accepted 2026-08-13 with PDR-042 and CH-153: semantic routing is now direct or tracked, and CI check selection no longer supplies the tracked-route signal.

**Who wants it:** contributors and agents integrating direct work or taking a tracked change to a ready pull request, in this repository and in any adopter running the shipped validation workflow (2026-08-11).

**Why:** the gate requiring a completed acceptance brief used to fire on the CI scope classifier's `full` output, which meant "this diff needs the full check suite" rather than "this work chose the tracked route". Every non-tracked change touching `docs/` was therefore asked for a tracked change's artifact, and the adopter-facing workflow had the same shape. Route and check scope now travel independently: branch history and any complete user-override trailers select tracked-route bookkeeping, while changed surfaces select relevant checks.

The accepted answer is that only a chosen tracked route owes the acceptance brief. A direct integration carries no Cliewen form, including when the agent originally recommended tracked and the user chose direct; in that case the complete current-head trailers preserve the authorization and risk without turning relevant surface checks off.

**Success looks like:**

- The gate and the change loop agree about which pull requests owe an acceptance brief, and the corpus records which reading was chosen and why.
- A contributor who is asked for a brief is told a reason that matches the rule they are being held to.
- A direct change that legitimately owes no brief does not fail a required check, in this repository or in an adopter's.
