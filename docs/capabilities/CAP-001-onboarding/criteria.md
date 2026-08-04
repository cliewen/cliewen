---
id: CAP-001-criteria
type: criteria
status: active
links: [CAP-001]
title: Acceptance criteria for CAP-001
---

```gherkin
Feature: Onboarding — install to first green validate

  @AC-001 @retired
  Scenario: A new user reaches green validate in under 30 minutes
    # Retired 2026-07-17 (CH-020): the 30-minute clock spans a human
    # journey (reading, installing) no focused test pair can verify.
    # The mechanical path is covered by AC-002/AC-024/AC-025; the
    # 30-minute end-to-end promise is C-015, owned by the quickstart.

  @AC-002
  Scenario: init produces a corpus that validate accepts unchanged
    Given an empty git repository
    When the user runs "clue init"
    Then "clue validate" exits with code 0 without any manual edits

  @AC-003
  Scenario: validate fails loudly on a broken corpus
    Given a corpus scaffolded by "clue init"
    And an artifact whose links reference a non-existent ID
    When the user runs "clue validate"
    Then validate exits with a non-zero code
    And the output names the offending file and the missing ID

  @AC-024
  Scenario: init is idempotent and keeps hand-written prose
    Given a repository scaffolded by "clue init"
    And an artifact added to a taxonomy folder afterwards
    When the user runs "clue init" again
    Then the folder README's index block references the new artifact
    And prose outside the clue:index markers is unchanged
    And a pre-existing taxonomy README without markers gains an appended index block, its prose intact
    And a re-run with nothing new to index changes no file

  @AC-025
  Scenario: init never replaces an existing file
    Given a repository that already contains one of the files init emits
    When the user runs "clue init"
    Then the existing file is not replaced and its prose outside clue:index markers is unchanged
    And the report names it as skipped
    And every file the existing one did not shadow is still created

  @AC-036
  Scenario: The public guide gives an operator one supported next step
    Test-type: Unit
    Given a reader has completed the disposable Cliewen trial
    When they open the public guide's operations page
    Then it distinguishes shipped and verified support from methodology intent
    And it gives safe recovery paths for routine operating problems
    And every public-guide page ends with exactly one primary next action

  @AC-038
  Scenario: init does not write through a symlinked skills folder
    Given a repository whose ".claude/skills" is a symlink to a skills tree shared across checkouts
    When the user runs "clue init"
    Then no file is created inside the link's target
    And the report names ".claude/skills" as skipped because it is a symlink
    And the canonical ".agents/skills" skills and the rest of the convention are still created
    And a ".claude/skills" that is an ordinary directory is mirrored as before, with no such report line

  @AC-039
  Scenario: The published plugin bootstraps and ships none of the managed skills
    Given the plugin marketplace manifest published at the repository root
    When its listed plugin is resolved to the committed tree it names
    Then that tree ships exactly one skill, and it is the bootstrap
    And none of the managed Cliewen lifecycle skills appear among its components
    And nothing in the bootstrap pins a "clue" version
    And the plugin tree lies outside every directory the skill generator owns

  @AC-062
  Scenario: init emits a thin caller for an upstream validation workflow
    Test-type: Unit
    Given an empty git repository
    When the user runs "clue init"
    Then ".github/workflows/clue.yml" calls Cliewen's reusable validation workflow at the emitting source's 40-hex immutable commit reference, or at its matching protected release tag when module build metadata is unavailable
    And the caller exposes only the runner, clue source, and install-directory choices while pinning the clue version to the generated pair
    And the caller contains no copied checkout, scope-detection, warning, acceptance-brief, or validation steps

  @AC-064
  Scenario: A corpus migration previews and applies one safe coordinated upgrade
    Test-type: Unit
    Given a repository with an older corpus shape, recognized older generated carriers, and a thin caller at an older immutable reference
    When the user runs "clue migrate"
    Then it reports the exact corpus, generated-skill, mirror, and caller transformations without changing a file
    And "clue migrate --apply --reversal-cost=low" applies the complete preflighted plan while preserving unrelated prose and caller-owned choices
    And a second run is a no-op
    But missing semantic choices, ambiguous syntax, interrupted source changes, and locally modified carriers fail without partial writes

  @AC-080
  Scenario: A recognized preceding managed set receives a newly introduced skill directory
    Test-type: Unit
    Given a repository whose complete remaining generated skill set exactly matches a supported preceding release
    And that release did not ship a newly introduced managed skill directory
    When the user previews "clue migrate" with the newer binary
    Then the plan adds the new canonical skill directory with the target release's bytes
    And its preview names both the target release it writes and the preceding release it recognized
    And applying the plan produces a no-op on a second run
    But a partial or locally modified remaining set leaves the new directory a finding and blocks every write

  @AC-071
  Scenario: init emits a Claude Code entry point that only points at the hub
    Test-type: Unit
    Given an empty git repository
    When the user runs "clue init"
    Then a root "CLAUDE.md" imports "AGENTS.md" and says that rules belong in the hub
    And it duplicates no sentence of the hub it points at
    And a repository whose "CLAUDE.md" already exists keeps it byte-for-byte and sees it reported as skipped

  @AC-072
  Scenario: migrate reports an entry point that never reaches the hub
    Test-type: Unit
    Given an adopted repository whose Claude Code entry point is missing or imports nothing
    When the user runs "clue migrate"
    Then the plan reports that Claude Code reads no routing and names the remedy for that shape — "clue init" for an absent file, the import line for one the adopter wrote
    And a "CLAUDE.md" that imports the hub produces no such report, whether by an import that resolves to it from that file's own directory or by being a symlink to the hub
    But no run writes or rewrites that file, with or without "--apply"

  @AC-083
  Scenario: init emits a hub that asks the agent whether the repository is behind
    Test-type: Unit
    Given an empty git repository
    When the user runs "clue init"
    Then the materialized "AGENTS.md" asks the agent to run the quiet release check when it starts
    And it routes a non-empty answer to the upgrade skill rather than to an installation command
    But no file is emitted for any assistant's configuration, and no emitted file declares a hook

  @AC-084
  Scenario: migrate reports a hub that never asks whether the repository is behind
    Test-type: Unit
    Given an adopted repository whose routing hub is missing or never names the release check
    When the user runs "clue migrate"
    Then the plan reports that no session learns a release is available and names the remedy for that shape — "clue init" for an absent hub, the line for one the adopter wrote
    And a hub that names the check produces no such report
    But no run writes or rewrites that hub, with or without "--apply"

  @AC-085
  Scenario: a real session learns the repository is behind
    Test-type: Human
    Given a behind repository whose hub carries the instruction and whose clue carries the notifier
    When a human starts a coding-agent session in it and the session runs any ordinary clue workflow command
    Then the session is told the repository is behind, in exactly one line, without the human saying anything about releases
    And a current repository produces no such line, so the absence of one is informative
    But the evidence is what the session did, not what the hub says: an agent that skipped the instruction and learned it anyway from the notifier is this criterion passing, and neither channel reporting it is this criterion failing
```
