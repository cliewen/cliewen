---
id: CAP-001-criteria
type: criteria
status: draft
links: [CAP-001]
title: Acceptance criteria for CAP-001
---

```gherkin
Feature: Onboarding — install to first green validate

  @AC-001
  Scenario: A new user reaches green validate in under 30 minutes
    Given a machine with no Cliewen tooling installed
    And an empty git repository
    When the user installs the clue binary
    And runs "clue init"
    And runs "clue validate"
    Then validate exits with code 0
    And the elapsed time from install to green is under 30 minutes

  @AC-002
  Scenario: init produces a corpus that validate accepts unchanged
    Given an empty git repository
    When the user runs "clue init"
    Then "clue validate" exits with code 0 without any manual edits

  @AC-003
  Scenario: validate fails loudly on a broken corpus
    Given a corpus scaffolded by "clue init"
    And an artifact whose links reference a non-existent ID
    When the user runs "clue validate"
    Then validate exits with a non-zero code
    And the output names the offending file and the missing ID
```
