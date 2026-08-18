---
id: ADR-017
type: decision
status: verified
links: [ADR-013, ADR-044, ADR-045]
title: Prose conventions register as constraint artifacts with enforcement classes
author: agent
accepted-by: Flemming N. Larsen (2026-07-14, PR #11 review)
---

# ADR-017 — Prose conventions register as constraint artifacts

## Context and problem statement

Rules living only in AGENTS.md, READMEs, or skills lack an inventory and a visible path from agent enforcement to machine enforcement.

## Decision outcome

**Every prose-only convention is a constraint artifact in `docs/constraints/`, and `enforcement: agent` is the promotion backlog.** Each constraint declares its `source:` and enforcement class (`machine|partial|agent|human` under ADR-045); `clue validate` checks them and reports the agent-enforced count. Partial and human constraints state their check and residual judgment, while an agent constraint states the trigger that would promote it to a machine check. The constraints README is the register, with no second inventory; repo-local prose carries pointers because shipped skills contain no corpus IDs.

**Carrier:** constraint-field lint and count in `clue validate`, the scaffolded register README, and the pre-PR checklist. A prose table is rejected because it is not a lintable artifact; promoting everything immediately is rejected because some rules need review or a machine with diff context.
