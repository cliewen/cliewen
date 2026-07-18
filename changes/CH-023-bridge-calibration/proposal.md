---
id: CH-023
type: change
status: open
links: [P-002, M-007]
title: Calibrate foreign-soil trials on robocode-api-bridge
---

# CH-023 — Calibrate foreign-soil trials on `robocode-api-bridge`

## What

Run a read-only calibration trial against the brownfield `robocode-api-bridge` repository, preserve the blind interpretation before comparing it with maintainer intent, and record the result as AN-003. The repository shares a maintainer with Cliewen, so the analysis is a control case and does not count as either of M-007's required foreign-soil trials.

Use trial feedback to make two corrections. First, express third-party software inspection as a jurisdiction-neutral evidence and permission process rather than embedding a legal conclusion from one jurisdiction. Second, carry the existing timeless-decision rule from C-006 into every skill that creates or verifies decision records: a decision states what is decided and only the enduring context and rationale needed to understand it; incident timelines and change history stay in the change workspace, PR, and Git history.

## Why

Serves P-002/M-007 by testing the trial method before it reaches genuinely foreign repositories. The calibration showed that a plausible brownfield interpretation needs explicit confidence layers, clean-environment evidence, population eligibility rules, and maintainer comparison. Review then exposed two methodology-carrier gaps: the skills did not make C-006's timelessness rule visible at the point where decisions are written, and the analysis mistook one jurisdiction's software exception for a durable global boundary.

No new decision is introduced. The skill edits implement the already active C-006 constraint, and the inspection correction removes an overbroad legal conclusion rather than adopting a new legal policy.
