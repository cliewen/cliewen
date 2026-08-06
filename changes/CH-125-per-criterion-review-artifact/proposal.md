---
id: CH-125
type: change
status: open
links: [P-012, M-060, AN-017, ADR-049]
title: Make the per-criterion registry refusal an explicit migration limit
---

# CH-125 — Make the per-criterion registry refusal an explicit migration limit

## Proposal

Serve P-012's M-060 by recording the deliberate refusal to commit a second, per-criterion registry to an extraction pull request. The decision will name what a reviewer cannot read from such an artifact, the cost that imposes on an adopter, and the report-and-manifest route that replaces it.

AN-017 currently classifies the assessment request as declined, but the corpus has not yet made the resulting limit explicit and accepted. ADR-049 already keeps the parity report derived rather than editable; CH-124 makes that report internally trustworthy. This change will make the remaining review-boundary trade-off legible without claiming that the declined artifact was delivered.

## Scope boundary

This change decides and documents the migration review artifact boundary. It does not add a committed per-criterion registry, alter parity or report derivation, or attempt M-061's order-of-magnitude proof.
