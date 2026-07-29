---
id: CH-080
type: change
status: active
links: [P-008, M-033, AN-003, AN-004, AN-005]
title: Analysis evidence boundaries distinguish prepared environments and qualified populations
---

# CH-080 — Analysis evidence boundaries distinguish prepared environments and qualified populations

## Proposal

M-033 makes the analysis workflow preserve two distinctions demonstrated by the calibration and foreign-soil trials: a result from a clean disposable environment is not the same as one from a prepared environment with local prerequisites, and a percentage or statistical claim needs a declared population and uncertainty boundary. Update the canonical `clue-analysis` guidance, regenerate its managed copies, record the lasting workflow decision at the proportional decision-log tier, and add focused guards for the resulting rule.

## Why now

AN-003 found a passing prepared build that could not support an onboarding-reproducibility claim and showed why an informal bot population could bias a percentage result. AN-004 and AN-005 independently confirmed that environment conditions and quality claims materially affect how evidence is interpreted. P-008/M-033 consumes those findings before later extraction and adoption work builds on analysis guidance.

## Scope

This change changes analysis guidance and its generated distributions only. It does not change `clue validate`, add an adoption interface, or convert population-level quality claims into deterministic acceptance criteria.
