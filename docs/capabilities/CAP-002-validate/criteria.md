---
id: CAP-002-criteria
type: criteria
status: active
links: [CAP-002]
title: Acceptance criteria for CAP-002
---

```gherkin
Feature: clue validate — deterministic corpus judgment

  @AC-004
  Scenario: A valid corpus passes
    Given a corpus whose artifacts all carry id, type, status, links and title
    And every link resolves and every index matches its folder
    When the user runs "clue validate"
    Then it exits with code 0
    And reports the number of artifacts discovered

  @AC-005
  Scenario: A missing core field fails loudly
    Given an artifact missing one of the core frontmatter fields
    When the user runs "clue validate"
    Then it exits with a non-zero code
    And the output names the file and the missing field

  @AC-006
  Scenario: A dangling link fails loudly
    Given an artifact whose links reference an ID no artifact carries
    When the user runs "clue validate"
    Then it exits with a non-zero code
    And the output names the file and the unresolved ID

  @AC-007
  Scenario: Index drift fails loudly
    Given a taxonomy README whose index block references a deleted file
    Or a sibling artifact its folder's index block does not reference
    When the user runs "clue validate"
    Then it exits with a non-zero code
    And the output names the README and the drifted entry

  @AC-008
  Scenario: The digest-before-merge gate
    Given a repository containing files under /changes
    When CI runs "clue validate --forbid-changes"
    Then it exits with a non-zero code
    And without the flag the same corpus passes

  @AC-009
  Scenario: An acceptance criterion without a test fails
    Test-type: Unit
    Given a criteria.md with status active containing an @AC tag
    And no test function whose name references that AC
    When the user runs "clue validate"
    Then it exits with a non-zero code
    And the output names the criteria file and the untested AC
    But ACs in a criteria.md with status draft are exempt

  @AC-010
  Scenario: A test referencing an unknown AC fails
    Given a test function whose name references an AC
    And no criteria.md anywhere declares that AC
    When the user runs "clue validate"
    Then it exits with a non-zero code
    And the output names the test file and the unknown AC

  @AC-011
  Scenario: A test without a declared purpose fails
    Given a test function whose name matches no purpose class
    When the user runs "clue validate"
    Then it exits with a non-zero code
    And the output names the test file, the function and the taxonomy
    But tests declaring Unit, Sanity or Arch pass without referencing any AC

  @AC-012
  Scenario: A retired AC needs no test, and its surviving tests fail
    Given a criteria.md scenario tagged with an AC and "@retired" on its tag line
    When the user runs "clue validate"
    Then the retired AC requires no test
    But a test still referencing the retired AC exits with a non-zero code

  @AC-013
  Scenario: Duplicate AC declarations fail
    Given the same AC ID declared more than once across the corpus
    When the user runs "clue validate"
    Then it exits with a non-zero code
    And the output names both declaring files

  @AC-023
  Scenario: Constraint artifacts carry their register fields
    Given a constraint artifact missing a non-empty source or enforcement field
    Or carrying an enforcement value outside machine, agent, human
    When the user runs "clue validate"
    Then it exits with a non-zero code
    And the output names the file and the violated field
    And a valid corpus reports its count of agent-enforced constraints on the OK line

  @AC-034
  Scenario: A byte-order mark in a corpus file fails
    Given a corpus markdown file containing a UTF-8 byte-order mark, at the start or embedded
    When the user runs "clue validate"
    Then it exits with a non-zero code
    And the output names the file and says to strip the byte-order mark
    But a BOM-free corpus passes

  @AC-035
  Scenario: A second frontmatter block in an artifact fails
    Given an artifact whose body begins with a second frontmatter block after the authoritative closing fence
    When the user runs "clue validate"
    Then it exits with a non-zero code
    And the output names the file and calls the block a leftover second frontmatter
    But a thematic break without a matching closing fence does not trigger the check
    And thematic breaks enclosing ordinary markdown or nothing at all do not trigger the check

  @AC-043
  Scenario: Declared test types require classified positive and negative evidence
    Test-type: Unit
    Given an active acceptance criterion declaring the Unit test type
    And a test reference classified as Unit and positive
    When the user runs "clue validate"
    Then it exits with a non-zero code naming the missing negative evidence
    But a Unit negative reference makes the criterion covered

  @AC-044
  Scenario: Cucumber tags provide classified acceptance-criterion evidence
    Test-type: Unit
    Given an active acceptance criterion declaring the E2E test type
    And a Cucumber .feature scenario tagged with its AC ID, e2e, and positive
    When the user runs "clue validate"
    Then it exits with a non-zero code naming the missing negative evidence
    But a Cucumber scenario tagged e2e and negative makes the criterion covered

  @AC-045
  Scenario: A Human-class criterion validates without a code test
    Test-type: Unit
    Given an active acceptance criterion declaring the Human test type
    And no Go, JVM, or Cucumber evidence references it
    When the user runs "clue validate"
    Then it reports no issue for that criterion
    But a Human-class criterion additionally declaring "(single-direction)" fails as malformed

  @AC-046
  Scenario: A @draft criterion is exempt from the active-file test requirement
    Test-type: Unit
    Given an active acceptance criterion whose tag line also carries @draft
    And no test references it
    When the user runs "clue validate"
    Then it reports no issue for that criterion
    But removing @draft from an otherwise untested criterion in an active file restores the "has no test" issue

  @AC-047
  Scenario: Coverage reports derive per-capability state without a committed registry
    Test-type: Unit
    Given a capability whose non-retired criteria are all satisfied, Human-classed, or @draft-exempt
    When the user runs "clue validate --coverage"
    Then it reports that capability as covered
    But a capability with only some criteria satisfied reports partial, and one with none satisfied reports gap

  @AC-049
  Scenario: A supersedes entry naming a still-live artifact fails loudly
    Test-type: Unit
    Given an artifact whose supersedes field names an ID
    And an artifact carrying that same ID still exists in the corpus
    When the user runs "clue validate"
    Then it exits with a non-zero code naming the file and the not-actually-retired ID
    But once the still-live artifact is deleted and nothing else references it, the corpus passes
    And any artifact carrying status: retired is rejected outright, since retirement is deletion, not a status
    And an artifact naming its own ID in supersedes is rejected, since an artifact cannot retire itself
    And two artifacts both claiming to supersede the same ID are rejected, naming both claimants, rather than one being silently picked

  @AC-050
  Scenario: A dangling link to a superseded ID names its successor
    Test-type: Unit
    Given an artifact whose supersedes field names a retired ID
    And a separate artifact whose links field still references that retired ID
    When the user runs "clue validate"
    Then it exits with a non-zero code naming the retired ID and the successor that superseded it

  @AC-051
  Scenario: Reversal cost bounds inferred provenance at capability activation
    Test-type: Unit
    Given an active capability joined by a links edge to a non-decision artifact carrying provenance inferred and reversal-cost high
    When the user runs "clue validate"
    Then it exits with a non-zero code naming the inferred artifact and the active capability it blocks
    But reversal-cost low remains valid while inferred
    And the output separately counts high-cost inferred non-decisions and inferred decisions awaiting verification

  @AC-052
  Scenario: Incident analyses derive the capabilities where green met wrong
    Test-type: Unit
    Given an analysis carrying reality contradicted and a links edge to a capability or one of its acceptance criteria
    When the user runs "clue validate --reality-gaps"
    Then it lists the affected capability and the analysis that contradicted it
    But reality contradicted on a non-analysis or without a capability or live criterion link fails validation
```
