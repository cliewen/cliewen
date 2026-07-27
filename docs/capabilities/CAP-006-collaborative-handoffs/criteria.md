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
    And a merged or closed pull request stops with any local work reported as unpublished
```
