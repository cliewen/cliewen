---
id: CAP-010-criteria
type: criteria
status: active
links: [CAP-010]
title: Acceptance criteria for team-safe identity allocation
---

```gherkin
Feature: Team-safe identity allocation

  @AC-169
  Scenario: The checked-in ledger is an append-only lifecycle event log
    Test-type: Unit
    Given an identity moves from reserved to live to retired
    When the version-two ledger is saved and loaded
    Then each transition remains as one event and the effective state is retired
    But conflicting immutable metadata for the same identity is rejected

  @AC-170
  Scenario: Remote claims merge without downgrading local lifecycle state
    Test-type: Unit
    Given a local ledger marks an identity live and the remote journal still carries its permanent reservation
    When the remote claim is synchronized once or repeatedly
    Then the effective local state remains live and the claim is not duplicated
    But conflicting identity metadata is rejected

  @AC-171
  Scenario: Concurrent clones receive unique sequential identities
    Test-type: Integration (single-direction)
    Given ten clones start from a ledger whose last claimed change identity is CH-170
    When they allocate through the same Git allocator branch concurrently
    Then they receive CH-171 through CH-180 exactly once each

  @AC-172
  Scenario: Concurrent worktrees share the remote allocation boundary
    Test-type: Integration (single-direction)
    Given two worktrees share one local Git object store and one configured allocator remote
    When both allocate from the same accepted ledger state concurrently
    Then each receives a distinct sequential identity without changing the other's working tree

  @AC-173
  Scenario: A batch can be assigned and synchronized without allocator write access
    Test-type: Integration
    Given a maintainer atomically reserves several sequential IDs through the allocator branch
    When another checkout runs clue id sync through a readable remote
    Then every assigned ID appears locally as reserved without a push
    But sync rejects a repository that has not enabled Git coordination

  @AC-174
  Scenario: Local allocation remains explicit and warns about concurrency
    Test-type: Unit
    Given a version-two ledger remains in local mode
    When clue id next reserves one or a positive batch of IDs
    Then it prints every reserved ID and warns that clones and worktrees are not coordinated
    But a zero or negative batch is rejected

  @AC-175
  Scenario: Coordinated allocation fails closed when durable state is unavailable
    Test-type: Integration (single-direction)
    Given a repository has enabled Git coordination and its remote is unreachable, malformed, or missing its established allocator ref
    When clue id next attempts an allocation
    Then it exits nonzero and leaves the local ledger byte-identical
    And a missing established ref says to restore rather than recreate it

  @AC-176
  Scenario: Migration preserves version-one ledger meaning
    Test-type: Unit
    Given a version-one ledger contains numeric and opaque identities, lifecycle states, large components, and source provenance
    When clue migrate converts it to version-two events and adds the union merge attribute
    Then every identity retains its effective meaning and no number is reissued
    But migration leaves coordination local until the team explicitly enables Git coordination

  @AC-177
  Scenario: A remote claim remains singular and recoverable after local uncertainty
    Test-type: Integration (single-direction)
    Given the allocator remote accepts a new identity claim
    When the client loses the push response or cannot save the claim to its local ledger
    Then the command reports the remotely reserved identity
    And clue id sync recovers the reservation without allocating another number

  @AC-179
  Scenario: A ledger combined by Git's union merge recovers without losing an identity
    Test-type: Unit
    Given Git's union merge combined two branches that both changed the identity ledger, duplicating the lines that differ
    When a command reads the ledger
    Then every identity survives at its furthest-along state and the next allocation is past all of them
    And the command names what combined the file and how to rewrite it instead of reporting a parse failure
    And saving the ledger writes it back whole, so a command that allocates also repairs
    But a duplicated setting whose two values disagree, in a ledger still carrying its settings inline, is refused for a person to decide, naming both values

  @AC-181
  Scenario: Allocation settings live where a disagreement is an ordinary merge conflict
    Test-type: Integration
    Given two branches each enabled coordination against a different Git remote
    When the branches are merged
    Then Git raises a conflict on the coordination file for the person merging to resolve, and the union-merged ledger beside it still combines cleanly
    And an unresolved conflict left in that file is reported as one rather than as a parse failure
    But pointing an already-coordinated repository at a different remote is refused unless forced, because it abandons the claims recorded on the remote it leaves
```
