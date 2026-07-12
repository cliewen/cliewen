---
id: CAP-002-criteria
type: criteria
status: active
links: [CAP-002]
title: Acceptance criteria for CAP-002
---

```gherkin
Feature: clue validate — deterministic corpus judgment

  @AC-004
  Scenario: A valid corpus passes
    Given a corpus whose artifacts all carry id, type, status, links and title
    And every link resolves and every index matches its folder
    When the user runs "clue validate"
    Then it exits with code 0
    And reports the number of artifacts discovered

  @AC-005
  Scenario: A missing core field fails loudly
    Given an artifact missing one of the core frontmatter fields
    When the user runs "clue validate"
    Then it exits with a non-zero code
    And the output names the file and the missing field

  @AC-006
  Scenario: A dangling link fails loudly
    Given an artifact whose links reference an ID no artifact carries
    When the user runs "clue validate"
    Then it exits with a non-zero code
    And the output names the file and the unresolved ID

  @AC-007
  Scenario: Index drift fails loudly
    Given a taxonomy README whose index block references a deleted file
    Or a sibling artifact its folder's index block does not reference
    When the user runs "clue validate"
    Then it exits with a non-zero code
    And the output names the README and the drifted entry

  @AC-008
  Scenario: The digest-before-merge gate
    Given a repository containing files under /changes
    When CI runs "clue validate --forbid-changes"
    Then it exits with a non-zero code
    And without the flag the same corpus passes
```
