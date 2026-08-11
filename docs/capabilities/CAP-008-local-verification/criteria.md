---
id: CAP-008-criteria
type: criteria
status: active
links: [CAP-008]
title: Acceptance criteria for CAP-008
---

# Acceptance criteria — CAP-008

```gherkin
Feature: Local verification commands

  @AC-136
  Scenario: The documented coverage report uses the portable invocation
    Test-type: Unit
    Given a contributor follows the local verification block in CONTRIBUTING.md
    When they inspect and run its coverage-report command on a supported contributor environment
    Then the command is "go tool cover -func coverage.out"
    And it renders a coverage profile
    And the deprecated equals-form invocation is absent from that block
```
