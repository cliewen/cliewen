---
id: CH-124
type: change
status: open
links: [P-012, M-059, CAP-003, ADR-049]
title: Extraction reports are derived from parity manifests
---

# CH-124 — Extraction reports are derived from parity manifests

## Proposal

Serve P-012's M-059 by replacing hand-maintained extraction-report criterion counts and mapping tables with readable output derived from the same pinned parity manifest that `clue parity` consumes.

Today an extraction report can describe a different source or mapping from the manifest committed beside it: parity checks the manifest while report prose retains independently typed figures and tables. This change will establish one derived-report form, make its source and target state inspectable, and reject the legacy typed form so the two cannot silently diverge.

The change will record the durable report-contract decision, add focused acceptance evidence for derivation and disagreement failure paths, update migration guidance, and mark M-059 done only when its complete exit criterion is demonstrated. The derived report remains a human-readable document, not a second registry.

## Scope boundary

This change concerns durable extraction reports and the parity manifests they represent. It does not decide whether a per-criterion artifact must be committed to a pull request (M-060), nor does it undertake the order-of-magnitude proof (M-061).
