---
id: CH-110-tasks
type: tasks
status: active
links: [CH-110-proposal]
title: Tasks for CH-110
---

# CH-110 — tasks

## Decisions first

- [ ] Write the ADR stating that `clue validate` judges a repository state and never a transition, naming where transition rules live instead
- [ ] Write the ADR amending ADR-017: the register names the machine holding each rule and the judgment that remains; add the `partial` class and widen `human`
- [ ] Record the milestone status vocabulary as a decision-log row, checking every existing plan against it first

## Criteria before implementation

- [ ] Add AC-089…AC-095 to CAP-002's criteria with `Test-type: Unit`, and AC-096 for the completed-plan CI guard with `Test-type: Integration`

## Validator checks

- [ ] Accept `partial` in the enforcement vocabulary and require `Checked by`/`Residual` on `partial` and `human` constraints (AC-095)
- [ ] Prose-layout lint for hard-wrapped paragraphs and list items (AC-089)
- [ ] `[-]` task lines carry prose after the checkbox (AC-090)
- [ ] Proposal declares a plan item or `plan-less` (AC-091)
- [ ] No image links or image files under `docs/**` (AC-092)
  - [ ] Scan collects non-markdown files so the image-file half can see them
- [ ] Per-type required frontmatter fields (AC-093)
- [ ] Milestone status cells match the declared vocabulary (AC-094)
- [ ] Run each new check against this repository and repair what it legitimately finds

## CI guard

- [ ] Workflow step failing a pull request that modifies a plan whose `main`-side status is `completed` (AC-096)

## The register

- [ ] Reclassify C-001…C-013, each stating its machine or its residual and cost, replacing the promotion triggers
- [ ] Update the register README: the vocabulary, what the OK-line count means, and the regenerated table
- [ ] Resolve the register table's badge inconsistency — generated rows state `status`, hand-written rows state `enforcement`

## Carriers

- [ ] Amend ADR-017's vocabulary sentence and the architecture door note
- [ ] CAP-002 criteria and design state the new checks and the state-not-transition boundary
- [ ] CLI OK-line text and `clue help`, if either states the backlog framing
- [ ] Scaffold templates: the constraints README, the shipped C-001, and the `AGENTS.md` note
- [ ] Guide pages printing the OK line or describing enforcement classes
- [ ] Decide whether the scaffolded C-001 reclassification needs a migration, and add it if so

## Close

- [ ] `go test ./...`, `clue validate --forbid-changes`, coverage at or above the C-014 floor, guide build
- [ ] `[Unreleased]` changelog entry written for users
