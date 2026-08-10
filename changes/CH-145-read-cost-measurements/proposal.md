---
id: CH-145
type: change
status: open
links: [P-015, M-071, CAP-007, PDR-029]
title: Measure corpus read cost
---

# CH-145 — Measure corpus read cost

Implement P-015/M-071 so `clue validate` reports the corpus populations that make read cost visible without turning either measurement into a validation failure. The command will count and, on request, name artifacts with more than one document and identities whose default bounded context slice exceeds the budget established by this change.

The change also records the corpus-wide one-primary-consumer-per-file principle, its enforcement class and residual, and the rejected measurement alternatives. This lets a later repair change act only on measured findings instead of inferring them from directory layout or subjective file size.
