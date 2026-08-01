---
id: CH-095-open-questions
type: open-questions
status: open
links: [CH-095]
title: Open questions for CH-095
---

# Open questions

## OQ-001 — Does the resolver run in CI, and can it ever block a merge?

The judge is the wall: deterministic, offline, and required by branch protection. The resolver is the opposite by design — it depends on other systems being up, so its answer can differ between two runs over the same revision.

Running it in CI as an advisory job is useful: a reference that rotted since the last release gets noticed without anyone remembering to look. Making it a required check is not, because an outage in someone else's infrastructure would then block an unrelated merge, which is the failure the two-layer split exists to prevent.

The recommendation is an advisory scheduled job rather than a pull-request gate: it reports on a cadence, opens its findings where a human reads them, and never sits between a change and its merge. The alternative — no CI role at all, resolver run by hand — is safe but relies on someone choosing to run it, which is how the 48 unqualified references accumulated in the first place.

This is not blocking implementation: the command works the same either way, and the CI role can be decided before the change is ready.
