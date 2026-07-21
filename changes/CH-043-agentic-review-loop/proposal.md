---
id: CH-043
type: change
status: open
links: [G-001, ARCH-002, PDR-007, C-012]
title: Add an agentic review loop before publication
---

# CH-043 — Add an agentic review loop before publication

Plan-less. This change closes one bounded methodology gap in the existing change-verification handoff rather than advancing a campaign milestone.

## What

- Make adversarial review an automatic part of `clue-verify` before a change is published, rather than an extra review the human must remember to request.
- Prefer a new context-isolated, read-only reviewer when the agent host supports delegation; otherwise require an explicitly identified in-context fallback and disclose the weaker review mode.
- Review the committed candidate against its declared intent, corpus, diff, tests, constraints, and quality scenarios; return only actionable correctness, regression, security, evidence, or unjustified-complexity findings.
- Send findings back to the implementing context, which fixes and locally verifies them, then starts a fresh review pass; a change is locally ready only after a clean pass, and every substantive fix invalidates the prior clean result.
- Preserve the existing ready-PR, exact hosted-head, CI, human-merge, and review-fix boundaries.

This is a full change because it changes generated methodology carriers, the skills architecture, and the agent-enforced readiness constraint. It changes no CLI or corpus-format API and makes no semantic plan mutation.
