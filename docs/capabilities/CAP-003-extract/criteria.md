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

  @AC-109
  Scenario: A clean parity run passes and is deterministic
    Test-type: Unit
    Given a source manifest whose entries all match the derived target manifest in proof class, direction, and evidence location
    When the user runs "clue parity"
    Then it exits zero and reports no unmatched or altered entries
    And a second run against the same inputs produces byte-identical report content

  @AC-110
  Scenario: A missing criterion fails parity
    Test-type: Unit
    Given a source manifest entry with a proof class and no matching disposition or exclusion
    And no target entry exists for that criterion ID
    When the user runs "clue parity"
    Then it exits with a non-zero code
    And the report names the ID as a missing criterion
    But an entry the source manifest declares "excluded" with a reason is not reported missing

  @AC-111
  Scenario: An orphaned tag fails parity
    Test-type: Unit
    Given classified target evidence for a criterion ID absent from the source manifest entirely
    When the user runs "clue parity"
    Then it exits with a non-zero code
    And the report names the ID as an orphaned tag

  @AC-112
  Scenario: A changed direction or evidence location fails parity
    Test-type: Unit
    Given a source manifest entry recording a proof class, direction, and evidence location
    And the derived target entry for that ID has a different direction set or evidence location
    When the user runs "clue parity"
    Then it exits with a non-zero code
    And the report names the ID and shows the source and target values that disagree
    But an unchanged direction and evidence location produces no finding for that ID

  @AC-113
  Scenario: A stale source fingerprint fails parity
    Test-type: Unit
    Given a source manifest whose "source-revision" disagrees with the revision recorded on that ID's ledger entry
    When the user runs "clue parity"
    Then it exits with a non-zero code
    And the report names the ID and both the manifest and recorded revisions

  @AC-114
  Scenario: An unjustified draft, Human, or retirement disposition fails parity
    Test-type: Unit
    Given a target criterion classified "@draft", "Human", or retired
    And its source manifest entry has no matching "disposition" and "justification"
    When the user runs "clue parity"
    Then it exits with a non-zero code
    And the report names the ID and the unjustified disposition
    But a source entry carrying the matching disposition and justification produces no finding for that ID

  @AC-115
  Scenario: Extraction preserves in-flight source work as an inspectable imported-change record
    Test-type: Unit
    Given a source repository's pending change with a proposal, a design rationale, a dependency, and a proof-linked task
    When extraction writes the corresponding imported-change record
    Then its origin, intent, design rationale, dependency link, and proof-links table all remain readable from the record alone
    But a record with no "## Proof links" table declares no inspectable proof

  @AC-116
  Scenario: A complete imported-change record requires every proof link to be proven
    Test-type: Unit
    Given an imported-change record whose status is "complete"
    When the user runs "clue validate"
    Then a proof-links row naming a criterion that does not exist, is "@draft", or is retired fails, naming the record and the criterion
    But a complete record whose every proof-linked criterion exists, is undrafted, and is not retired passes

  @AC-117
  Scenario: An in-progress imported-change record may name unproven proof links
    Test-type: Unit
    Given an imported-change record whose status is "in-progress"
    And a proof-links row names a criterion that does not yet exist or is still "@draft"
    When the user runs "clue validate"
    Then it is not rejected for that criterion, because "in-progress" declares work still pending

  @AC-118
  Scenario: A clean carrier reconciliation run passes and is deterministic
    Test-type: Unit
    Given a carrier inventory whose every mapped entry's fingerprint matches its target's current content
    When the user runs "clue carriers"
    Then it exits zero and reports no findings
    But a mapped target that drifts after a clean baseline fails on a later run

  @AC-119
  Scenario: A carrier inventory entry maps to a target or blocks mutation, never both
    Test-type: Unit
    Given an inventory entry naming neither a target-path and fingerprint nor a blocked marker with a reason
    Or an entry combining a target-path and fingerprint with a blocked marker
    When the inventory is loaded
    Then it is rejected
    But a "blocked" entry with a reason is accepted and is not reconciled against any target

  @AC-120
  Scenario: A lost fingerprint fails carrier reconciliation
    Test-type: Unit
    Given a mapped carrier entry whose target path exists
    And its current content fingerprint disagrees with the inventory's pinned fingerprint
    When the user runs "clue carriers"
    Then it exits with a non-zero code
    And the report names the ID as a lost fingerprint
    But a matching fingerprint produces no finding for that ID

  @AC-121
  Scenario: A missing target fails carrier reconciliation
    Test-type: Unit
    Given a mapped carrier entry whose target-path does not exist in the reconciled corpus
    When the user runs "clue carriers"
    Then it exits with a non-zero code
    And the report names the ID as a missing asset
    But a present target produces no finding for that ID

  @AC-122
  Scenario: A stale deleted-path reference fails carrier reconciliation
    Test-type: Unit
    Given an inventory naming a source-repository path in "deleted-paths"
    And a local Markdown link anywhere in the reconciled corpus still resolves to that path
    When the user runs "clue carriers"
    Then it exits with a non-zero code
    And the report names the referencing file as a stale deleted-path reference
    But a link that does not resolve to a deleted path produces no finding

  @AC-123
  Scenario: Disposable brownfield fixtures prove the composed migration contract
    Test-type: Unit
    Given disposable numeric-archive and opaque-identifier source fixtures with pinned revisions, classified source evidence, pending work, and operational carriers
    When the approved fixture mutation creates their target corpora
    Then each target validates and its derived parity and carrier reports are clean
    And an archived numeric identity and an opaque identity cannot be reused
    And every required parity and carrier failure path is rejected by its deterministic command
    But the fixture source's own test result is not presented as Cliewen acceptance evidence

  @AC-125
  Scenario: A deferred migration criterion has inspectable accountability
    Test-type: Unit
    Given a source manifest disposition for a draft, Human, or retired criterion
    When the user runs "clue parity" against its target corpus
    Then the disposition names a source location and an existing milestone plan door
    And the report states the derived count of deferred criteria on both clean and failing runs
    But a missing accountability field is rejected and an unknown plan door fails parity

  @AC-126
  Scenario: A durable extraction report's figures are derived from its manifest
    Test-type: Unit
    Given an extraction report whose derived region names a pinned source manifest
    When the user runs "clue validate"
    Then a region rendered from that manifest passes
    But a typed count or mapping row, a region left stale by a revised manifest, and a region naming a manifest that is not there each fail

  @AC-127
  Scenario: The derived region is regenerated rather than written
    Test-type: Unit
    Given an extraction report whose derived region names a pinned source manifest
    When the user runs "clue report" against that report
    Then the region is rendered from the manifest and the report's prose outside it is unchanged
    But a report with no derived region is refused rather than given one

  @AC-128
  Scenario: Parity and carrier commands retain their deterministic failure behavior at assessment scale
    Test-type: Unit
    Given a disposable target with hundreds of criteria and classified evidence references plus tens of retired identities
    When the public `clue` command validates it and runs parity and carrier reconciliation
    Then clean parity and carrier reports pass, numeric allocation advances past every prefix's recorded population, and identity reservation refuses an existing, retired, or already-reserved identity across a ledger round trip
    And every required parity and carrier failure class still exits non-zero against that assessment-scale shape
    But this command-scale fixture does not claim extraction ordering, source-work preservation, or pinned-release evidence; AC-129 holds those

  @AC-129
  Scenario: The ordered migration path holds under a pinned release at assessment scale
    Test-type: Unit
    Given an assessment-scale source and target, a `clue` binary stamped with a pinned release, and installed skills carrying that same stamp
    When the rehearsal pins are written before the target exists and the approved mutation is verified against those unmodified pins
    Then the stamped binary validates the target, the rehearsal's own manifest and inventory still reconcile clean afterwards, and every in-flight imported change keeps its pinned origin, rationale sections, and proof links to live criteria
    But a skill stamp disagreeing with the pinned release is reported as drift, a parity or carrier claim made before the target exists fails, and an imported change claiming complete over an unproven criterion is rejected at that size
```
