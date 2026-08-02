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

  @AC-073
  Scenario: An appended index row states the record it links
    Test-type: Unit
    Given a taxonomy folder holding an artifact no index row references
    When the index block is regenerated
    Then the appended row states the artifact's id, its title, and its status
    And a title whose value contains a colon is spelled as the parsed value, not as its YAML quoting
    But an artifact missing any of a readable id, title, and status falls back to the plain link, rather than a row carrying an empty status
    And a row referencing a subfolder README carries no title or status, because it states a section rather than a record
```
