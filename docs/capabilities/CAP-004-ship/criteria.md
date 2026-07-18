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

  @AC-020 @retired
  Scenario: A skill without a version stamp fails
    # Retired 2026-07-18 (CH-028): location alone no longer makes every
    # skill Cliewen-owned. AC-031 carries the version requirement after
    # AC-029 establishes ownership.

  @AC-021 @retired
  Scenario: The skills must agree on one version
    # Retired 2026-07-18 (CH-028): the version set is now explicitly the
    # marked Cliewen skills. AC-032 carries the scoped consistency rule.

  @AC-022 @retired
  Scenario: A released clue detects skill drift
    # Retired 2026-07-18 (CH-028): binary drift now applies only after a
    # skill declares Cliewen ownership. AC-033 carries that scoped rule.

  @AC-028
  Scenario: Versioned skills are generated as standalone artifacts
    Given canonical skill-specific templates and shared instruction fragments
    When a maintainer runs the repository skill generator
    Then it deterministically writes complete skills to both the agent and embedded-template trees
    And corresponding files in both trees are byte-identical
    But a missing, changed, or unexpected generated file fails the repository tests and names the drift

  @AC-029
  Scenario: The ownership marker scopes Cliewen skill validation
    Given .agents/skills contains a skill marked "cliewen-skill: true"
    And it contains an unmarked third-party skill
    When the user runs "clue validate"
    Then only the marked skill joins the Cliewen version set
    But a present "cliewen-skill" value other than boolean true fails and names the malformed skill

  @AC-030
  Scenario: Pre-marker Cliewen skill slots fail toward migration
    Given one of the five canonical Cliewen skill directories contains an unmarked skill.md
    When the user runs "clue validate"
    Then it exits with a non-zero code and tells the user to reinstall that legacy skill
    But an unmarked skill in any other directory is ignored

  @AC-031
  Scenario: A marked Cliewen skill requires a version stamp
    Given a skill marked "cliewen-skill: true" with no string version in its frontmatter
    When the user runs "clue validate"
    Then it exits with a non-zero code and names the missing or invalid version
    But a marked skill with a string version passes the stamp-presence check

  @AC-032
  Scenario: Marked Cliewen skills must agree on one version
    Given two marked Cliewen skills whose frontmatter declares different versions
    When the user runs "clue validate"
    Then it exits with a non-zero code and names the disagreeing skill
    But marked skills that agree and any unmarked skills pass the set-consistency check

  @AC-033
  Scenario: A released clue detects drift in marked Cliewen skills
    Given a clue binary stamped with a release version
    And a marked Cliewen skill whose version differs from that release
    When the user runs "clue validate"
    Then it exits with a non-zero code and reports the drift
    But an unmarked skill does not participate, a "dev" build skips the comparison, and a release matching the marked skills passes
```
