---
id: CH-151
type: change
status: open
links: [P-015, M-073, CAP-004]
title: Make generated skills load their rules progressively
---

# CH-151 — Make generated skills load their rules progressively

The generated Cliewen skills currently render every shared rule and every workflow branch into `skill.md`, so invoking one lifecycle entry point loads instructions that are needed only later or only on a branch the work never takes. M-073 requires the carrier to cost less to load without dropping any rule it states.

Change the generated artifact from one complete file to one complete skill directory. Keep `skill.md` as the short routing entry point and generate relative `references/` files for cross-cutting and branch-specific instruction. Each route states exactly when its reference must be read, so deferred text is not silently omitted. Copying one generated skill directory remains sufficient; neither another skill nor the canonical source tree is a runtime dependency.

Record the architectural change in ADR-059, revise CAP-004's contract through AC-137, and keep both generated output trees byte-identical and fully drift-checked. Existing methodology rules keep their meaning; only their loading boundary changes.

