---
id: CAP-003-criteria
type: criteria
status: active
links: [CAP-003]
title: Acceptance criteria for CAP-003
---

```gherkin
Feature: Brownfield analysis and extraction — evidence, namespaced ACs, JVM harvesting, provenance

  @AC-014
  Scenario: A criteria file declares ACs in its own namespace
    Given a criteria.md whose frontmatter declares "ac-prefix: MG"
    And a scenario tagged with the ID MG-101
    When the user runs "clue validate"
    Then MG-101 is enforced like any AC: with no test the run fails naming the criteria file
    And a Go test named TestMG101_… satisfies it
    And criteria without an ac-prefix keep the default namespace AC unchanged

  @AC-015
  Scenario: A tag outside the file's namespace fails
    Given a criteria.md whose frontmatter declares "ac-prefix: MG"
    And a scenario tagged with the foreign ID PG-001
    When the user runs "clue validate"
    Then it exits with a non-zero code
    And the output names the file, the foreign tag and the file's namespace

  @AC-016 @retired
  Scenario: A JVM test tag satisfies AC coverage
    Given a declared AC MG-101 in active criteria
    And a Kotlin or Java test file containing @Tag("MG_101")
    When the user runs "clue validate"
    Then MG-101 counts as tested
    And tags outside every declared namespace are ignored as runner metadata

  @AC-017 @retired
  Scenario: A JVM tag referencing an unknown or retired AC fails
    Given a Kotlin test file containing @Tag("MG_999") which no criteria.md declares
    Or a tag referencing an AC retired by tombstone
    When the user runs "clue validate"
    Then it exits with a non-zero code
    And the output names the test file and the offending tag

  @AC-018 @retired
  Scenario: Provenance is linted and inferred artifacts are counted
    Given an artifact carrying "provenance" with a value outside inferred|verified
    Or a decision carrying a provenance field at all
    When the user runs "clue validate"
    Then it exits with a non-zero code naming the file
    But a valid corpus with inferred artifacts passes
    And the OK line reports how many artifacts are born inferred

  @AC-054
  Scenario: Extraction phases evidence at criterion granularity
    Test-type: Unit
    Given a brownfield extraction contains proven and not-yet-proven criteria
    When the agent records their evidence state
    Then a genuinely not-yet-proven individual criterion may carry @draft while proven siblings and the capability remain active
    And a genuine Human-class criterion uses the pull request acceptance brief as proof without a code reference
    But whole-file draft phasing remains available when the capability's extracted criteria are not ready for active use

  @AC-055
  Scenario: Analysis qualifies verification environments and population claims
    Test-type: Unit
    Given analysis evidence from a clean disposable or prepared environment
    When the investigator records the evidence boundary
    Then it distinguishes the two environments
    And a clean result has no local prerequisites while any local prerequisite makes a result prepared
    And a prepared result names its prerequisites without claiming onboarding reproducibility
    And a statistical or percentage claim names its versioned corpus and population, eligibility rules, exclusions with reasons, sampling or repetition method, uncertainty, and deterministic-versus-quality boundary
    And adoption analysis names the governance or process changes it introduces
    But an environment-sensitive quality claim is not represented as a deterministic acceptance criterion
    And scaffolding is not described as neutral

  @AC-056
  Scenario: Extraction rehearses before it mutates
    Test-type: Unit
    Given an extraction full change has been proposed
    When the agent begins the extraction
    Then it first writes a report-only rehearsal under that change's workspace
    And the rehearsal inventories source formats and entry points, proposed mappings, preserved and minted IDs, confidence and reversal cost, test-purpose work, instruction conflicts, planned deletions, and plan doors
    And it changes no target source corpus, Cliewen corpus, tests, routing, or hosted state
    But an unresolved conflict is recorded as an open question and stops before mutation
    And mutation begins only with explicit human direction and digests the rehearsal into the durable extraction report

  @AC-058
  Scenario: JVM evidence uses a conservative executable carrier
    Test-type: Unit
    Given a Java or Kotlin executable carries one AC identity, one proof type, and one direction in its JUnit method annotations or stable test name
    When the user runs "clue validate"
    Then ordinary, parameterized, repeated, factory, template, and nested JUnit executables receive one statically attributable evidence credit
    And the named form test<PREFIX><digits><lowercase-suffix>_<Type><Direction>_<description> works without native framework tags, with segmented-prefix hyphens removed
    And an attached reference to an unknown or retired criterion fails
    But class-level AC tags, unsupported annotation syntax, and disagreeing tag-and-name carriers receive no credit and produce a diagnostic
    And structured proximity comments remain ignored

  @AC-060
  Scenario: Extended criterion identities work through supported evidence carriers
    Test-type: Unit
    Given declared IDs SNAP-SQS-001 and ADP-045b
    When supported Go, JVM, and Cucumber carriers reference those IDs using their documented normalized forms
    Then each carrier resolves the original canonical identity and classified evidence keeps its proof type and direction
    And lowercase letter suffixes remain part of the identity
    But uppercase suffixes, ambiguous normalized prefixes, and malformed normalized forms receive no credit and a diagnostic

  @AC-061
  Scenario: Extraction preserves and deterministically mints extended criterion identities
    Test-type: Unit
    Given a source corpus containing multi-segment and letter-suffixed criterion IDs plus requirements without IDs
    When the agent performs extraction
    Then source IDs are copied verbatim and only unlabelled requirements receive the namespace's next numeric IDs in source order
    And rerunning against the same source state produces the same preserved and minted mapping
    But extraction never renumbers a source ID to fit the default AC grammar
```
