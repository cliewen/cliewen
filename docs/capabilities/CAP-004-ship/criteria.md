---
id: CAP-004-criteria
type: criteria
status: active
links: [CAP-004]
title: Acceptance criteria for CAP-004
---

```gherkin
Feature: clue ships — a versioned binary and versioned skills, drift made lintable

  @AC-019
  Scenario: clue reports its release version
    Given a clue binary stamped with a release version at build time
    When the user runs "clue version" or "clue --version"
    Then it prints that version
    But an unstamped source build reports "dev" rather than a release number

  @AC-020 @retired
  Scenario: A skill without a version stamp fails
    # Retired 2026-07-18 (CH-028): location alone no longer makes every
    # skill Cliewen-owned. AC-031 carries the version requirement after
    # AC-029 establishes ownership.

  @AC-021 @retired
  Scenario: The skills must agree on one version
    # Retired 2026-07-18 (CH-028): the version set is now explicitly the
    # marked Cliewen skills. AC-032 carries the scoped consistency rule.

  @AC-022 @retired
  Scenario: A released clue detects skill drift
    # Retired 2026-07-18 (CH-028): binary drift now applies only after a
    # skill declares Cliewen ownership. AC-033 carries that scoped rule.

  @AC-028
  Scenario: Versioned skills are generated as standalone artifacts
    Given canonical skill-specific templates and shared instruction fragments
    When a maintainer runs the repository skill generator
    Then it deterministically writes complete skills to both the agent and embedded-template trees
    And corresponding files in both trees are byte-identical
    But a missing, changed, or unexpected generated file fails the repository tests and names the drift

  @AC-029
  Scenario: The ownership marker scopes Cliewen skill validation
    Given .agents/skills contains a skill marked "cliewen-skill: true"
    And it contains an unmarked third-party skill
    When the user runs "clue validate"
    Then only the marked skill joins the Cliewen version set
    But a present "cliewen-skill" value other than boolean true fails and names the malformed skill

  @AC-030 @retired
  Scenario: Pre-marker Cliewen skill slots fail toward migration
    # Retired 2026-08-03 (CH-108): the managed set grew from five slots to
    # six. AC-082 carries the immutable six-slot meaning.

  @AC-082
  Scenario: Reserved Cliewen skill slots fail toward reinstall
    Test-type: Unit
    Given one of the six canonical Cliewen skill directories contains an unmarked skill.md
    When the user runs "clue validate"
    Then it exits with a non-zero code and tells the user to reinstall that legacy skill
    But an unmarked skill in any other directory is ignored

  @AC-031
  Scenario: A marked Cliewen skill requires a version stamp
    Given a skill marked "cliewen-skill: true" with no string version in its frontmatter
    When the user runs "clue validate"
    Then it exits with a non-zero code and names the missing or invalid version
    But a marked skill with a string version passes the stamp-presence check

  @AC-032
  Scenario: Marked Cliewen skills must agree on one version
    Given two marked Cliewen skills whose frontmatter declares different versions
    When the user runs "clue validate"
    Then it exits with a non-zero code and names the disagreeing skill
    But marked skills that agree and any unmarked skills pass the set-consistency check

  @AC-033
  Scenario: A released clue detects drift in marked Cliewen skills
    Given a clue binary stamped with a release version
    And a marked Cliewen skill whose version differs from that release
    When the user runs "clue validate"
    Then it exits with a non-zero code and reports the drift
    But an unmarked skill does not participate, a "dev" build skips the comparison, and a release matching the marked skills passes

  @AC-037
  Scenario: A skill's manifest is resolved the same way on every filesystem
    Given a skill directory whose manifest is named "SKILL.md" rather than "skill.md"
    When the user runs "clue validate"
    Then the skill joins the managed set exactly as a lowercase "skill.md" would, so the verdict does not depend on the host filesystem
    And a single directory holding two case-variants of the manifest name is reported as a named ambiguity rather than silently resolving to one
    But that ambiguity is reported for an unmarked third-party skill too, because ownership cannot be read until the manifest is resolved

  @AC-065
  Scenario: A release that changes corpus obligations names migration guidance
    Test-type: Unit
    Given a release tag whose first-parent diff changes the versioned migration registry
    When the release workflow extracts the changelog section
    Then it requires a non-empty "### Migration" section before publishing
    But a release that does not change the migration registry keeps the ordinary release-notes requirement without inventing migration guidance

  @AC-068
  Scenario: Resolving external addresses reports without condemning
    Test-type: Unit
    Given corpus addresses that answer 200, a bare 403, a redirect with a location, 404 whose owner root still answers, 404 whose owner root does not, a 403 carrying a rate-limit header, and a timeout
    When the user runs "clue refs"
    Then they are classified reachable, restricted, redirected, restricted, gone, unreachable, and unreachable in that order
    And only the gone address makes the command exit non-zero
    But restricted and unreachable never fail the run, because neither says the corpus is wrong
    And a clue: identity is never resolved

  @AC-069
  Scenario: A rewrite is offered, never taken behind the human
    Test-type: Unit
    Given a redirected address in an ordinary artifact and a redirected address in a completed plan
    When the user runs "clue refs" and then "clue refs --apply"
    Then the preview writes nothing and names both, marking the completed plan as pinned history
    And --apply rewrites only the ordinary artifact, to the location its host gave
    But the completed plan is left exactly as it was, because its address is part of what was observed
    And a gone address is never rewritten, since no correct replacement exists to offer

  @AC-075
  Scenario: A newer release is reported with the route for the machine it is running on
    Test-type: Unit
    Given a released clue whose version is behind the newest published release
    When the user runs "clue latest" on Windows, on macOS or Linux, and on a platform with no prebuilt asset
    Then each run names both versions and prints exactly one installation route — the PowerShell script, the shell script, and "go install" pinned to that release
    And each run also prints the coordinated migrate preview-and-apply sequence, because moving the binary alone produces drift rather than resolving it
    But no run writes a file in the repository, with or without flags, and no run replaces the binary

  @AC-081
  Scenario: A generated skill routes an existing adopter through a human-authorized upgrade
    Test-type: Unit
    Given a repository carrying the generated Cliewen skill set
    When an agent invokes the managed upgrade skill
    Then it runs the release check and reads the selected release's migration guidance
    And it asks the human whether to upgrade now or later before changing a repository
    And an affirmative answer makes the repository green, branches, moves the coordinated set, resolves every notice, verifies, and opens a pull request without merging it
    But the skill names no platform installation command and a later answer changes no repository file or hosted state

  @AC-076
  Scenario: The quiet mode is silent unless there is something to say
    Test-type: Unit
    Given a released clue behind the newest release and the same clue already on it
    When the user runs "clue latest --quiet" in each
    Then the run that is behind prints one line naming the newer release
    But the run that is current prints nothing at all, and both exit 0

  @AC-077
  Scenario: Not being able to tell is never a failure
    Test-type: Unit
    Given a release host that is offline, one that never answers within the timeout, one that answers with a rate limit, and one that answers with a body the command does not recognize
    When the user runs "clue latest" and "clue latest --quiet" against each
    Then every run exits 0 and writes nothing to standard error
    And the quiet runs print nothing at all
    But the ordinary runs say the release list could not be reached rather than claiming the repository is current
    And a running stamp the command cannot read as a release is reported as uncomparable rather than as current, while a pre-release stamp is compared as its release numbers and ordered before the release it names, build metadata orders with the release it decorates, and a stamp newer than anything published is reported as unpublished rather than as the newest release

  @AC-078
  Scenario: The answer is cached outside the repository, and a broken cache is absence
    Test-type: Unit
    Given a cached answer inside its lifetime, one older than its lifetime, and a cache file whose bytes cannot be read as an answer
    When the user runs "clue latest" against each
    Then the fresh cache is used and no request is made, and the stale and unreadable caches ask the host and store what it said
    And the cache is written under the user's cache directory, never inside the repository
    But a cache that cannot be read or written never fails the command, because a cache that can fail a command is worse than no cache

  @AC-079
  Scenario: The drift report names the way out and the way to stay
    Test-type: Unit
    Given a released clue whose version differs from its marked Cliewen skills
    When the user runs "clue validate"
    Then the drift issue names the command that reports the upgrade and the command that moves the repository
    And it names both halves of staying on the release the repository carries — the matching binary and the pinned CI caller
    But it still fails the run and still names both versions, because what the rule decides has not changed
```
