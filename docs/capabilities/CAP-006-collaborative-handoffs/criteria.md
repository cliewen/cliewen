---
id: CAP-006-criteria
type: criteria
status: active
links: [CAP-006]
title: Acceptance criteria for CAP-006
---

# Acceptance criteria — CAP-006

```gherkin
Feature: Collaborative pull-request handoffs

  @AC-040
  Scenario: Review findings remain durable across agent sessions
    Given an agent reviews an open pull request at a named hosted head commit
    When it returns findings or a clean result
    Then each actionable finding is recorded as an unresolved hosted review conversation where the forge supports one
    And each conversation remains unresolved until a hosted commit contains the reviewed repair
    And a clean result names only the commit it covers
    And any later substantive edit invalidates that result
    And the resulting commit requires a new review
    But a reviewer unable to publish a resolvable finding reports the pull request as not merge-ready without claiming equivalent enforcement

  @AC-041
  Scenario: Any agent that edits completes the hosted handoff or fails safely
    Given an agent starts from the current hosted head of an open pull request
    When it repairs the change
    Then it commits and verifies the complete repair
    And a clean review covers the repaired commit
    And it pushes without force and confirms the pull request head equals that commit
    And it resolves satisfied findings only after that confirmation
    But if the hosted head changes or a normal push is rejected it does not overwrite remote work or claim the pull request is merge-ready
    And an open changed head is fetched and reconciled before verification and review repeat
    And if accepted main advances after publication, current main is merged into the pull request branch without rewriting hosted history before verification and review repeat
    And a merged or closed pull request stops with any local work reported as unpublished

  @AC-042
  Scenario: The human merge gate receives its remaining semantic decisions
    Given a full Cliewen change adds or changes acceptance criteria
    When the agent prepares its ready pull request
    Then the pull request starts with an acceptance brief naming the plan item, the criteria and scenarios, advisory scenario-resolution verdicts, and merge-binding decision effects
    And CI rejects a full pull request that leaves the brief's required placeholders unfilled
    But a scenario-resolution verdict is advisory and does not make `clue validate` fail

  @AC-048
  Scenario: The acceptance brief names newly declared Human-class criteria as their proof
    Test-type: Unit
    Given a full Cliewen change adds or materially revises a criterion declaring the Human test type
    When the agent prepares its ready pull request
    Then the acceptance brief's criteria line names that criterion and states that the brief is its proof
    But a change touching no Human-class criterion leaves the brief's existing criteria line unchanged

  @AC-063
  Scenario: An upstream workflow update leaves adopter choices in the caller
    Test-type: Unit
    Given two versions of a thin caller that differ only in the reusable workflow commit reference
    When the adopter updates that reference
    Then the runner, clue source, install directory, and clue version choices remain unchanged
    And the upstream workflow still owns scope detection, the armed warning, the acceptance-brief gate, and the "clue validate --forbid-changes" step
    And the branch-protection probe can identify the same stable "validate" check
```
