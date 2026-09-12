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
  Scenario: The routing hub gives a fresh agent conditional orientation instructions
    Test-type: Unit
    Given clue init materializes the cross-agent repository routing hub
    When a fresh agent reads that hub
    Then the hub tells it to run clue next --all once after the mandatory clue latest --quiet check and before substantive work
    And the hub tells it to report repository position when the request leaves direction open or existing work materially affects the request
    And the hub tells it to read the selected candidate with clue context before recommending it
    But the hub rejects a routine status preamble for a concrete request with no relevant existing work
    And it says proposed work is never authorization to begin
```
