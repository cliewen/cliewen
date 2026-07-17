---
id: CAP-001-criteria
type: criteria
status: active
links: [CAP-001]
title: Acceptance criteria for CAP-001
---

```gherkin
Feature: Onboarding — install to first green validate

  @AC-001 @retired
  Scenario: A new user reaches green validate in under 30 minutes
    # Retired 2026-07-17 (CH-020): the 30-minute clock spans a human
    # journey (reading, installing) no focused test pair can verify.
    # The mechanical path is covered by AC-002/AC-024/AC-025; the
    # 30-minute end-to-end promise is QS-002, owned by the quickstart.

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

  @AC-024
  Scenario: init is idempotent and keeps hand-written prose
    Given a repository scaffolded by "clue init"
    And an artifact added to a taxonomy folder afterwards
    When the user runs "clue init" again
    Then the folder README's index block references the new artifact
    And prose outside the clue:index markers is unchanged
    And a re-run with nothing new to index changes no file

  @AC-025
  Scenario: init never overwrites an existing file
    Given a repository that already contains one of the files init emits
    When the user runs "clue init"
    Then the existing file's content is unchanged
    And the report names it as skipped
    And every file the existing one did not shadow is still created
```
