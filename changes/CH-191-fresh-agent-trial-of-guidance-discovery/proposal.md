---
id: CH-191
type: change
status: open
links: [P-023, M-098, G-015]
title: A fresh agent finds the CONTRIBUTING.md discovery, uses it, and corrects it if conditions changed
---

# CH-191 — Fresh-agent trial of guidance discovery

[P-023](../../docs/plans/P-023-challenge-plans-and-retain-what-work-teaches.md)/M-098 requires proof that CH-190's "Writing Multi-Line PR and Commit Text on Windows" entry actually works as guidance, not merely that it exists: an agent with no prior session must locate it from the normal entry point, use it to complete a real task, and correct it if conditions had changed since it was written. `Human`-class evidence must name who ran it, under what conditions, and what happened.

**The trial.** A fresh subagent — no memory of this repository's recent sessions, not told the entry exists or that this is a test of it — is given one real, already-scoped task: [M-101](../../docs/plans/P-023-challenge-plans-and-retain-what-work-teaches.md)'s correction, with [AN-026](../../docs/analysis/AN-026-trialling-the-challenge-rule.md) already settling the routing (no decision record needed) and the specific edit (stop drawing a persistent Change → Capability edge in `guide/methodology.md`; match `guide/intent.md`'s Goal-only linkage). That task requires this repository's own PR workflow on Windows — precisely the condition CH-190's entry addresses — so the agent must reach the entry through the same route CH-190 documented (`AGENTS.md`'s source-repository-conventions table → `CONTRIBUTING.md`'s Verify Locally block, opened whole) rather than being pointed at it.

**Why M-101 as the vehicle.** It is real, already correctly scoped by AN-026 (no invented exercise), small and reversible, and ends in an actual PR with a multi-line body — the exact condition under which the CH-190 entry matters. Using it also delivers M-101 itself, so the trial is not throwaway process theater.

**What counts as passing.** The guidance existing is not enough, per M-098's own text. The trial records, as `Human`-class evidence: whether the fresh agent found the entry unprompted, whether it applied the procedure when building its commit/PR text, whether the text came out uncorrupted, and whether anything about the entry was stale or wrong under the conditions actually encountered (shell, `gh` version, task shape) — with a correction filed if so. Flemming observes the run and the resulting PR and states what happened; that observation is the proof, since `Human`-class evidence needs no code.

**Scope discipline.** This change's own tasks are the trial's setup and write-up, not M-101's content — M-101's diagram fix is the fresh agent's PR, reviewed and merged on its own, separately from this change.
