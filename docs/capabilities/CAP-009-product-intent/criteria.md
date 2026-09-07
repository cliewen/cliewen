---
id: CAP-009-criteria
type: criteria
status: active
links: [CAP-009, ADR-065, ADR-066, ADR-067, PDR-054, PDR-055]
title: Acceptance criteria for CAP-009
---

```gherkin
Feature: Product intent — one vision per corpus and optional use cases

  @AC-162
  Scenario: The judge holds the vision to one artifact at one address and claims nothing about its content
    Test-type: Unit
    Given a corpus carrying at most one artifact of type vision
    When the user runs "clue validate"
    Then a corpus with no vision at all passes
    And a single active vision at "docs/vision.md" passes without any rule about what it says
    But a second vision artifact fails and names both paths
    And a vision artifact anywhere other than "docs/vision.md" fails and names the required address
    And a "docs/vision.md" whose type is not vision fails
    And a vision still carrying the scaffold bootstrap marker fails and names the repair

  @AC-163
  Scenario: A use case is checked for form when it exists and is never required to exist
    Test-type: Unit
    Given a corpus that may contain any number of use-case artifacts
    When the user runs "clue validate"
    Then a corpus with no use cases passes, and no goal, capability, or criterion is reported for lacking one
    And a use case in "docs/use-cases" whose filename matches its UC identity, naming a goal and a capability and carrying the Actors, Trigger, Main flow, and Outcome sections, passes
    But a use-case artifact outside "docs/use-cases", or whose filename does not match its identity, fails
    And a use case whose links name no goal, or name no capability, fails and says which is missing
    And a use case missing any required structural section fails and names the section

  @AC-164
  Scenario: The intent report states what exists and computes no coverage figure
    Test-type: Unit
    Given a corpus whose vision and use cases are known
    When the user runs "clue validate --intent"
    Then it prints the vision's identity, status, and whether its meaning is inferred or human-authored
    And it prints each use case with the capabilities it crosses
    And a corpus with no vision reports the absence as a state rather than as an issue, and still exits 0
    But no percentage, ratio, or count of covered goals or capabilities is printed anywhere in the report

  @AC-165
  Scenario: Context names the use cases that reach an artifact without following them
    Test-type: Unit
    Given a use case whose links name a capability that the use case's goal does not name back
    When the user runs "clue context" on that capability
    Then the use case is named with its identity, title, and repository-relative path
    And no content of the use case is emitted and the slice's artifacts are unchanged
    And several such use cases are named in repository-relative path order
    But an artifact no use case names produces no such section, and the naming never follows an edge, so widening remains the reader's own act

  @AC-166
  Scenario: Init materializes the vision bootstrap and the optional use-case folder
    Test-type: Unit
    Given a repository with no Cliewen convention
    When the user runs "clue init"
    Then "docs/vision.md" is created carrying the bootstrap marker, and "docs/use-cases/README.md" is created with an empty index block
    And the corpus index references both
    But a repository that already carries either file has it reported and skipped, never overwritten

  @AC-167
  Scenario: Migration reports a missing vision and writes no vision content
    Test-type: Unit
    Given a repository that adopted Cliewen before the vision artifact existed
    When the user runs "clue migrate"
    Then the missing vision is reported as a notice that blocks nothing and leaves the plan appliable
    And the optional use-case folder and its index row are planned as structure
    And applying the plan writes no file under "docs/vision.md"
    But a repository that already carries a vision produces no such notice

  @AC-168
  Scenario: Intent discovery is proportionate and never presents inference as confirmed
    Test-type: Human
    Given a repository with no usable vision
    When an agent follows the shipped greenfield or brownfield discovery workflow
    Then a greenfield interview asks only questions that would change the vision, the initial goals, the boundary, or a candidate use case, adapts to the answers, and stops when another question would change nothing
    And a brownfield draft reads the repository first, cites the sources behind its material claims, separates observation from interpretation, and records contradictions rather than resolving them
    And drafted content is marked draft and inferred with its assumptions and open questions visible, and no statement of why the product exists is derived from implementation structure alone
```
