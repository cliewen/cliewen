---
id: CAP-007-criteria
type: criteria
status: active
links: [CAP-007]
title: Acceptance criteria for CAP-007
---

```gherkin
Feature: Focused corpus context

  @AC-053
  Scenario: One identity yields its deterministic transitive dependency slice
    Test-type: Unit
    Given an artifact that links another artifact which links a third, plus an unrelated artifact
    When the user runs "clue context <id>"
    Then it emits the starting artifact first and each transitively linked artifact exactly once with its ID, repository-relative path, and complete markdown content
    And artifacts at the same traversal depth are ordered by repository-relative path
    And an acceptance-criterion or milestone ID starts from the artifact that declares it
    But an unknown or ambiguously declared ID exits non-zero, names the ID, and emits no artifact content
```
