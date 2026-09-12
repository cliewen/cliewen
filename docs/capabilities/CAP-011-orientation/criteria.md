---
id: CAP-011-criteria
type: criteria
status: active
links: [CAP-011]
title: Acceptance criteria for repository orientation
---

```gherkin
Feature: Repository orientation

  @AC-190
  Scenario: Next work is derived across campaign boundaries
    Test-type: Unit
    Given a corpus with an open change workspace, active doing and todo milestones, a draft milestone, and proposed goals
    When the user runs clue next
    Then the open change is the first resume choice, followed by active doing and todo milestones in stable order
    And draft milestones and proposed goals are labeled as proposed rather than actionable
    And clue next --all lists every recorded category without changing the corpus
    But when no recorded work offers a candidate, the command says to capture a proposed goal instead of inventing work

  @AC-191
  Scenario: A fresh agent context receives orientation only when useful
    Test-type: Unit
    Given an agent starts a fresh context from the repository routing hub
    When it has completed the mandatory clue latest --quiet check
    Then it runs clue next --all once before substantive work
    And it briefly reports repository position when the request leaves direction open or existing work materially affects the request
    And it reads the selected candidate with clue context before recommending it
    But a concrete request with no relevant existing work receives no routine status preamble
    And the agent never treats proposed work as authorization to begin
```
