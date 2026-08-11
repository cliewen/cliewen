---
id: CH-146-proposal
type: change
status: open
links: [P-015, M-072, C-008, C-022, ADR-057]
title: Repair measured corpus read-cost findings
---

# CH-146 — Repair measured corpus read-cost findings

Implement P-015/M-072 by repairing the live artifacts reported by the multi-document measurement, splitting only where each resulting file has a distinct primary consumer and recording a reason for every count that remains accepted. Completed plans remain outside the repair population because their immutable state makes alteration unavailable.

The change will also bring every reported default bounded context slice within the eight-artifact budget M-071 established, so the reported over-budget population reaches zero, using the measurement as the sole repair input. Historical analysis sections will move verbatim if a split is required.
