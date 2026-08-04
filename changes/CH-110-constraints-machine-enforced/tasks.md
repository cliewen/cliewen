---
id: CH-110-tasks
type: tasks
status: open
links: [CH-110-proposal]
title: Tasks for CH-110
---

# CH-110 — tasks

## Decisions first

- [x] Write the ADR stating that `clue validate` judges a repository state and never a transition, naming where transition rules live instead
- [x] Write the ADR amending ADR-017: the register names the machine holding each rule and the judgment that remains; add the `partial` class and widen `human`
- [x] Record the milestone status vocabulary as a decision-log row, checking every existing plan against it first

## Criteria before implementation

- [x] Retire AC-023, whose enforcement set the widened vocabulary makes untrue, and add AC-089…AC-095 to CAP-002's criteria with `Test-type: Unit`
- [-] Mint an acceptance criterion for the completed-plan CI guard — no capability promises it: it is this repository's own bookkeeping rather than shipped behaviour, so it carries `Unit`- and `Sanity`-purpose tests like the other workflow guards

## Validator checks

- [x] Accept `partial` in the enforcement vocabulary and require `Checked by`/`Residual` on `partial` and `human` constraints (AC-089)
- [x] Prose-layout lint for hard-wrapped paragraphs and list items (AC-090)
- [x] `[-]` task lines carry prose after the checkbox (AC-091)
- [x] Proposal declares a plan item or `plan-less` (AC-092)
- [x] No image links or image files under `docs/**` (AC-093)
  - [x] Scan collects non-markdown files so the image-file half can see them
- [x] Per-type required frontmatter fields (AC-094)
- [x] Milestone status cells match the declared vocabulary (AC-095)
- [x] Run each new check against this repository and repair what it legitimately finds — the corpus was already clean; the findings were in test fixtures carrying capabilities without a goal, a decision without its signature fields, and a proposal declaring no plan scope

## CI guard

- [x] Workflow step failing a pull request that modifies a plan whose `main`-side status is `completed`, with tests proving both directions and that the digest closing a campaign still passes

## The register

- [x] Reclassify C-001…C-013, each stating its machine or its residual and cost, replacing the promotion triggers
- [x] Give C-014, C-015, and C-016 the same declarations, so the register reads one way throughout
- [x] Update the register README: the vocabulary, what the OK-line count means, and the regenerated table
- [x] Resolve the register table's badge inconsistency — the rows state the enforcement class, and the README now says so, so nobody "repairs" it back to status

## Carriers

- [x] Amend ADR-017's vocabulary sentence and the architecture door note
- [x] CAP-002 README, criteria, and design state the new checks and the state-not-transition boundary
- [x] CLI comment naming the count's criterion and what the count now means
- [x] Scaffold templates: the constraints README, the shipped C-001, and the `AGENTS.md` note
- [x] Guide sample output, which printed an agent-enforced count a fresh scaffold no longer has
- [x] MIG-007 reports a scaffolded constraint still awaiting a check this release ships, and repairs nothing; CAP-001's migration inventory states it

## Close

- [x] `go test ./...`, `clue validate --forbid-changes`, coverage at or above the C-014 floor, guide build
- [x] `[Unreleased]` changelog entry written for users
