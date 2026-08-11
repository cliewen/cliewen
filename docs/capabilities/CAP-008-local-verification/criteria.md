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
  Scenario: The documented coverage commands use the portable invocation
    Test-type: Unit
    Given a contributor follows the local verification block in CONTRIBUTING.md
    When they inspect and run its coverage commands on a supported contributor environment
    Then the block documents "go test ./... -coverprofile coverage.out" and "go tool cover -func coverage.out"
    And the report renders a coverage profile
    And no command in that block passes a single-dash "-flag=value" whose value contains a dot
```
