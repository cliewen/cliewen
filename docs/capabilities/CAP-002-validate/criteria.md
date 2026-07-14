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

  @AC-009
  Scenario: An acceptance criterion without a test fails
    Given a criteria.md with status active containing an @AC tag
    And no test function whose name references that AC
    When the user runs "clue validate"
    Then it exits with a non-zero code
    And the output names the criteria file and the untested AC
    But ACs in a criteria.md with status draft are exempt

  @AC-010
  Scenario: A test referencing an unknown AC fails
    Given a test function whose name references an AC
    And no criteria.md anywhere declares that AC
    When the user runs "clue validate"
    Then it exits with a non-zero code
    And the output names the test file and the unknown AC

  @AC-011
  Scenario: A test without a declared purpose fails
    Given a test function whose name matches no purpose class
    When the user runs "clue validate"
    Then it exits with a non-zero code
    And the output names the test file, the function and the taxonomy
    But tests declaring Unit, Sanity or Arch pass without referencing any AC

  @AC-012
  Scenario: A retired AC needs no test, and its surviving tests fail
    Given a criteria.md scenario tagged with an AC and "@retired" on its tag line
    When the user runs "clue validate"
    Then the retired AC requires no test
    But a test still referencing the retired AC exits with a non-zero code

  @AC-013
  Scenario: Duplicate AC declarations fail
    Given the same AC ID declared more than once across the corpus
    When the user runs "clue validate"
    Then it exits with a non-zero code
    And the output names both declaring files

  @AC-023
  Scenario: Constraint artifacts carry their register fields
    Given a constraint artifact missing a non-empty source or enforcement field
    Or carrying an enforcement value outside machine, agent, human
    When the user runs "clue validate"
    Then it exits with a non-zero code
    And the output names the file and the violated field
    And a valid corpus reports its count of agent-enforced constraints on the OK line
```
