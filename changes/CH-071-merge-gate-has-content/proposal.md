---
id: CH-071-proposal
type: proposal
status: active
links: [P-007]
title: The merge gate has content (M-024)
---

# CH-071 — The merge gate has content

## What

Implements [M-024](../../docs/plans/P-007-core-hardening.md) of P-007. Four things, carried by one PDR:

1. The Cliewen digest now emits an **acceptance brief** at the top of the PR body, before the rest of the description — the human answers which plan item this serves and whether it's still wanted, reviews the added/changed criteria verbatim with their scenarios, and sees what becomes binding on merge (`inferred` decisions minted, records invalidated/superseded). The template bakes in the book's competence-heuristic warning and a one-screen-cap as prose pressure toward small deltas.
2. `clue-verify`'s review loop gains a **per-criterion scenario-resolution step**: for each added/changed criterion, the reviewer compares its scenario against the referenced tests' setup/action/assertions and records an advisory verdict (verifies / verifies-something-adjacent / undetermined). Non-blocking, feeds the brief, does not gate `clue validate` or the review loop's actionable-finding lifecycle.
3. An **opt-in spec-first pause** after Propose, for human review of the proposal before implementation begins. Default ready-PR loop is unchanged.
4. A new **shared skill fragment** states that an agent's private memory is never where work lives — anything needed to implement, continue, or hand off work lives in a corpus artifact, the change workspace, or the PR.

Carried by: a new PDR, `clue-delta`/`clue-verify` source templates, the new shared fragment, `.github/pull_request_template.md`, and a deterministic unfilled-brief check in this repository's CI and the scaffolded wall.

## Why

[AN-008](../../docs/analysis/AN-008-methodology-critiques.md) pattern A: the human half of the merge boundary is empty — ~5,600 words of agent-facing skill instruction and zero words telling the human what to verify at merge. "Machines enforce form; humans verify meaning" is theatre without a stated human verification surface. The book's "after" gate and its advisory test-matches-scenario check already answer this; Cliewen imports that design rather than reinventing it.

## Approach note

Given this is a C-013 red-line change (alters what a merge accepts), the decision content (PDR draft) is being surfaced for review before the mechanical implementation (skill templates, PR template, CI check, new criteria/tests) proceeds — using the same opt-in pre-implementation pause this milestone introduces.
