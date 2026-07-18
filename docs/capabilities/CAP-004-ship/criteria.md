---
id: CAP-004-criteria
type: criteria
status: active
links: [CAP-004]
title: Acceptance criteria for CAP-004
---

```gherkin
Feature: clue ships — a versioned binary and versioned skills, drift made lintable

  @AC-019
  Scenario: clue reports its release version
    Given a clue binary stamped with a release version at build time
    When the user runs "clue version" or "clue --version"
    Then it prints that version
    But an unstamped source build reports "dev" rather than a release number

  @AC-020
  Scenario: A skill without a version stamp fails
    Given a repository whose .agents/skills holds a skill.md with no version in its frontmatter
    When the user runs "clue validate"
    Then it exits with a non-zero code
    And the output names the skill missing its version stamp

  @AC-021
  Scenario: The skills must agree on one version
    Given two skills whose frontmatter declares different versions
    When the user runs "clue validate"
    Then it exits with a non-zero code
    And the output names the disagreeing skill
    But skills that all declare the same version pass this check

  @AC-022
  Scenario: A released clue detects skill drift
    Given a clue binary stamped with a release version
    And a skill whose version differs from that release
    When the user runs "clue validate"
    Then it exits with a non-zero code and reports the drift
    But a "dev" build skips the comparison, and a release matching the skills passes

  @AC-028
  Scenario: Versioned skills are generated as standalone artifacts
    Given canonical skill-specific templates and shared instruction fragments
    When a maintainer runs the repository skill generator
    Then it deterministically writes complete skills to both the agent and embedded-template trees
    And corresponding files in both trees are byte-identical
    But a missing, changed, or unexpected generated file fails the repository tests and names the drift
```
