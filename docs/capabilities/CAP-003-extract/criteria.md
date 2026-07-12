---
id: CAP-003-criteria
type: criteria
status: active
links: [CAP-003]
title: Acceptance criteria for CAP-003
---

```gherkin
Feature: Brownfield extraction — namespaced ACs, JVM harvesting, provenance

  @AC-014
  Scenario: A criteria file declares ACs in its own namespace
    Given a criteria.md whose frontmatter declares "ac-prefix: MG"
    And a scenario tagged with the ID MG-101
    When the user runs "clue validate"
    Then MG-101 is enforced like any AC: with no test the run fails naming the criteria file
    And a Go test named TestMG101_… satisfies it
    And criteria without an ac-prefix keep the default namespace AC unchanged

  @AC-015
  Scenario: A tag outside the file's namespace fails
    Given a criteria.md whose frontmatter declares "ac-prefix: MG"
    And a scenario tagged with the foreign ID PG-001
    When the user runs "clue validate"
    Then it exits with a non-zero code
    And the output names the file, the foreign tag and the file's namespace

  @AC-016
  Scenario: A JVM test tag satisfies AC coverage
    Given a declared AC MG-101 in active criteria
    And a Kotlin or Java test file containing @Tag("MG_101")
    When the user runs "clue validate"
    Then MG-101 counts as tested
    And tags outside every declared namespace are ignored as runner metadata

  @AC-017
  Scenario: A JVM tag referencing an unknown or retired AC fails
    Given a Kotlin test file containing @Tag("MG_999") which no criteria.md declares
    Or a tag referencing an AC retired by tombstone
    When the user runs "clue validate"
    Then it exits with a non-zero code
    And the output names the test file and the offending tag

  @AC-018
  Scenario: Provenance is linted and inferred artifacts are counted
    Given an artifact carrying "provenance" with a value outside inferred|verified
    Or a decision carrying a provenance field at all
    When the user runs "clue validate"
    Then it exits with a non-zero code naming the file
    But a valid corpus with inferred artifacts passes
    And the OK line reports how many artifacts are born inferred
```
