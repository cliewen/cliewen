---
id: CAP-005-criteria
type: criteria
status: active
links: [CAP-005]
title: Acceptance criteria for CAP-005
---

# Acceptance criteria — CAP-005

```gherkin
Feature: Index generation — clue scaffold

  @AC-026
  Scenario: scaffold regenerates index blocks and keeps prose
    Given a corpus with an artifact added after the last index regeneration
    When the user runs "clue scaffold"
    Then the folder README's index block references the new artifact
    And prose outside the clue:index markers is unchanged
    And "clue validate" accepts the result
    And a run with nothing new to index changes no file

  @AC-027
  Scenario: scaffold touches only index blocks and materializes nothing
    Given a corpus with an artifact to index
    When the user runs "clue scaffold"
    Then no file is created or deleted and no file outside the taxonomy READMEs is modified
    And on a path without a docs tree the command exits non-zero, names the problem, and creates nothing
```
