---
id: CH-044
type: change
status: open
links: [G-001, ARCH-002, PDR-007, PDR-012]
title: Ground agentic findings in operative requirements
---

# CH-044 — Ground agentic findings in operative requirements

Plan-less. This change corrects a bounded review-quality gap discovered while using the automatic review loop introduced by CH-043.

## What

- Require an agentic-review finding to identify a violated operative requirement or a concrete failure of the declared change intent.
- Require reviewers to apply explicit lifecycle exceptions and authoritative decisions before treating nearby wording as contradictory, including the versioned changelog section created by a release cut and PDR-007's distinction between optional human review and mandatory human-controlled merge.
- Exclude alternative readings, historical wording, and checklist text that is truthful under the declared lifecycle from actionable findings when they do not change required behavior.
- Add regression evidence that holds these false-positive boundaries in the generated `clue-verify` skill.

This is a full change because it changes a generated methodology carrier. It changes no CLI behavior, acceptance criterion, capability meaning, or plan semantics.
