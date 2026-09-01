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
    # The mechanical path is covered by AC-150/AC-024/AC-025; the
    # 30-minute end-to-end promise is C-015, owned by the quickstart.

  @AC-002 @retired
  Scenario: init produces a corpus that validate accepts unchanged
    # Retired 2026-08-30 (CH-162): a CLI cannot honestly invent a repository's
    # system picture. AC-150 requires the explicit architecture and design
    # bootstrap that an agent activates from repository evidence.

  @AC-150
  Scenario: init and migration establish mandatory system-overview bootstraps
    Test-type: Unit
    Given an empty repository or an adopted Cliewen corpus missing a canonical system overview
    When the user runs "clue init" or previews and applies "clue migrate"
    Then the canonical architecture and design README paths are created with marked bootstraps without replacing existing overview prose
    And init tells the operator to replace the bootstraps before validation
    But migration neither writes through a symlink nor claims the bootstrap is repository-specific truth

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

  @AC-124 @retired
  Scenario: A missing thin CI caller does not block an unrelated safe migration
    # Retired 2026-08-15 (CH-154): ADR-060 reverses the clause that migration
    # neither creates nor rewrites the caller. The surviving half — a missing
    # optional carrier never blocks an independent safe migration — is carried
    # by AC-140, which also proves the materialization this criterion forbade.

  @AC-140
  Scenario: Migration materializes a missing thin CI caller
    Test-type: Unit
    Given an adopted repository without the optional ".github/workflows/clue.yml" thin caller
    And the migration plans a safe independent update
    When the user previews "clue migrate"
    Then the plan creates the caller from the embedded template at its default runner, clue-source, and install-directory choices
    And "clue migrate --apply" applies the creation and the independent update together
    And a second run is a no-op
    But a present caller keeps its adopter-owned choices and receives only its reference and version updates, and one whose content cannot be safely recognized remains a blocking finding

  @AC-141
  Scenario: Migration reports a competing validation wall beside the caller
    Test-type: Unit
    Given an adopted repository whose own workflow other than the thin caller runs the installed binary's "clue validate"
    When the user runs "clue migrate"
    Then it reports a finding naming that workflow file and its job
    And no file is changed
    But a repository whose only validation wall is the thin caller reports no such finding, and the finding never rewrites or deletes the repository-owned workflow

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

  @AC-107
  Scenario: migrate backfills the identity ledger from the current corpus scan without renumbering history
    Test-type: Unit
    Given a corpus with no ".clue/id-ledger.yaml" file and existing live artifacts across several native ID prefixes
    When the user runs "clue migrate --apply"
    Then it writes one "live" entry per currently-live ID, unchanged
    And it seeds each prefix's counter at that prefix's current maximum numeric component
    And a second run reports zero changes
    But a corpus that already carries a ".clue/id-ledger.yaml" file is left untouched by this step

  @AC-144
  Scenario: A legacy decision log requires reviewed classification before migration
    Test-type: Unit
    Given an adopted corpus whose legacy decision log contains several dated rows
    When the user previews or applies "clue migrate"
    Then the migration inventories every row and blocks all writes with guidance to classify future-shaping choices in a full change
    And a corpus whose reviewed conversion has removed the log receives no legacy-log finding or change
    But the migration never guesses a row's subject, durability, destination, or whether it is narrative that may be discarded

  @AC-145
  Scenario: A new corpus starts with only the subject-typed decision contract
    Test-type: Unit
    Given an empty git repository
    When the user runs "clue init" and then "clue scaffold"
    Then the decisions folder explains ADR, PDR, and IDR routing by subject and contains no decision log
    And the regenerated index remains valid without inventing a legacy register
    But existing decision files are not overwritten by either command

  @AC-153
  Scenario: A repository declares its Cliewen role, and an undeclared one is an adopter
    Test-type: Unit
    Given an empty repository, a repository already declaring a role, and a repository carrying no role marker at all
    When the user runs "clue init" or the tooling reads the repository's role
    Then init materializes an adopter marker in the empty repository and leaves an existing declaration untouched
    And a repository with no marker reads as an adopter without error and is blocked by nothing
    But a marker naming an unknown role, an empty role, or malformed YAML is rejected rather than defaulted

  @AC-156
  Scenario: Migration reports a spent analysis and never acts on it
    Test-type: Unit
    Given an analysis whose plans are all completed and whose declared carriers resolve, beside one whose plan is still active and one declaring no carrier
    When the user previews or applies "clue migrate"
    Then the spent analysis is reported as a non-blocking notice naming its completed plans and the artifacts carrying its findings
    And applying the migration writes nothing for it and removes no file
    But the analysis with an active plan and the analysis declaring no carrier are not reported
```
