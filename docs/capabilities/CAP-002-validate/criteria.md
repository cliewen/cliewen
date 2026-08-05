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

  @AC-023 @retired
  Scenario: Constraint artifacts carry their register fields
    Given a constraint artifact missing a non-empty source or enforcement field
    Or carrying an enforcement value outside machine, agent, human
    When the user runs "clue validate"
    Then it exits with a non-zero code
    And the output names the file and the violated field
    And a valid corpus reports its count of agent-enforced constraints on the OK line
    # Retired by ADR-045: partial joined the vocabulary and the register gained
    # required declarations, so this scenario's enforcement set is no longer the
    # rule. AC-089 states the widened register contract.

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

  @AC-057
  Scenario: JVM classified evidence belongs to one executable
    Test-type: Unit
    Given Java or Kotlin test files containing at least ten acceptance-criterion identities on distinct positive and negative executable methods
    When the user runs "clue validate"
    Then each criterion receives only the proof type and direction attached to its own executable
    And the complete mixed-file fixture passes
    But tags attached to unrelated methods cannot combine into classified evidence
    And an ambiguous executable carrying several identities, proof types, or directions receives no classified credit and produces a diagnostic

  @AC-059
  Scenario: Extended criterion identities remain canonical in validation and coverage
    Test-type: Unit
    Given active criteria files declare the IDs SNAP-SQS-001 and ADP-045b in their matching namespaces
    When the user runs "clue validate" and "clue validate --coverage"
    Then links, retirement tombstones, and derived coverage resolve those exact IDs without renumbering
    And duplicate IDs, wrong-namespace declarations, malformed near-matches, and colliding normalized prefixes fail with diagnostics
    But existing <PREFIX>-<digits> IDs remain valid without edits

  @AC-066
  Scenario: An external reference names its target or the corpus fails
    Test-type: Unit
    Given corpus prose containing a bare forge number, a full URL, a link whose label carries a forge number, a heading anchor, a colour literal in a code span, a fenced block, and a clue: identity
    When the user runs "clue validate"
    Then only the bare forge number fails, naming the file and the line a reader opens to
    And the failure says to write the full URL of what the reference points at
    But a full URL, a labelled link, an anchor, a fenced or code-span colour literal, fenced content, and a clue: identity produce no finding
    And a corpus artifact cited by its bare ID and a file cited by a relative path are untouched

  @AC-067
  Scenario: Foreign acceptance evidence is named but never counted
    Test-type: Unit
    Given corpus prose carrying the pointer clue:owner/repo@revision/ID and a second pointer that opens clue: without completing the form
    When the user runs "clue validate" and "clue validate --coverage"
    Then the well-formed pointer is accepted as a foreign identity naming its repository, revision, and identifier
    And the malformed pointer fails, naming the line and the form it should take
    But neither is ever counted as classified evidence, because the judge cannot see a run in another repository
    And --coverage lists it as named but locally unproven, apart from the per-capability states rather than inside them
    And a pointer written inside a fence or a code span is an example of the form, never a claim that something was proven

  @AC-074
  Scenario: Index rows that state only their own link are counted, never failed on
    Test-type: Unit
    Given a taxonomy README whose index row labels its link with the target's own filename
    When the user runs "clue validate"
    Then the run exits zero and the OK line counts the row as stating only its own link
    And "clue validate --index-rows" names the README and the target a reader opens to
    But a row stating its record's id and title, a row referencing a subfolder README, and a curated row covering several targets are never counted

  @AC-099
  Scenario: Index rows that state their record but not what it is about are counted, never failed on
    Test-type: Unit
    Given a taxonomy README whose index row states its record's id, title, and status and carries nothing after them
    When the user runs "clue validate"
    Then the run exits zero and the OK line counts the row as not saying what the artifact is about
    And "clue validate --index-rows" names the README and the target a reader opens to
    But a row carrying a description, a row referencing a subfolder README, and a curated row covering several targets are never counted
    And a row whose label restates only its own link is counted by the filler population instead, whether or not a hand edit has since given it a status badge, so the two populations stay disjoint
    And a row stating a record without the status badge the generator always writes is adopter prose in a shape no release produced, so it is left uncounted rather than graded

  @AC-089
  Scenario: The constraint register carries its fields and its declarations
    Test-type: Unit
    Given a constraint artifact missing a non-empty source or enforcement field
    Or carrying an enforcement value outside machine, partial, agent, human
    Or declaring partial or human without naming the machine that checks it and the residual judgment it leaves
    When the user runs "clue validate"
    Then it exits with a non-zero code
    And the output names the file and the violated field or missing declaration
    But a machine-enforced constraint needs neither declaration, and an agent-enforced one keeps its promotion trigger
    And a valid corpus reports its count of agent-enforced constraints on the OK line, saying nothing when that count is zero

  @AC-090
  Scenario: Hard-wrapped prose fails
    Test-type: Unit
    Given a corpus markdown file whose paragraph or list item is broken across two lines
    When the user runs "clue validate"
    Then it exits with a non-zero code
    And the output names the file and the line the continuation begins on
    But line breaks inside fenced code, indented code, frontmatter, tables written with or without outer pipes, an HTML block up to the blank line that ends it, and a comment up to its own closing marker are structure, not wrapping
    And a heading, a list item, a blockquote line, and a table row each begin their own block rather than continuing the one above
    And a fence documenting a shorter fence is closed only by a run of its own length, and an indented fence marker is code rather than a fence opening, so the rest of the file is still read

  @AC-091
  Scenario: A skipped task carries a reason
    Test-type: Unit
    Given a tasks.md under changes/ carrying a "[-]" task with nothing after its checkbox
    When the user runs "clue validate"
    Then it exits with a non-zero code
    And the output names the file and the line
    But a "[-]" task followed by prose passes, an unticked or ticked task is never asked for one, and one shown inside fenced or indented code is an example rather than a task

  @AC-092
  Scenario: A proposal declares the plan item it serves
    Test-type: Unit
    Given a proposal.md under changes/ whose links name no plan or milestone and whose body never says plan-less
    When the user runs "clue validate"
    Then it exits with a non-zero code
    And the output names the file and both ways to satisfy the rule
    But a proposal linking a P or M identity passes, and so does one declaring itself plan-less

  @AC-093 @retired
  Scenario: Diagrams are inline, not images
    Test-type: Unit
    Given a docs file carrying an image — an inline link, a reference link, or an img tag — or an image file stored under docs/
    When the user runs "clue validate"
    Then it exits with a non-zero code
    And the output names the file and says diagrams are inline Mermaid
    But an ordinary link to a markdown file passes, and an image inside a fence, an indented code block, a comment, or a code span is an example rather than a diagram
    And a table cell is not an exemption: an image there renders like any other
    # Retired by ADR-047: image links and assets are valid corpus content, and
    # choosing Mermaid, ASCII art, or SVG is human judgment rather than a
    # deterministic validation rule.

  @AC-100
  Scenario: Image links and assets remain valid corpus content
    Test-type: Unit (single-direction)
    Given a docs file carrying local and absolute image links in inline, reference, collapsed-reference, and HTML img forms
    And a local SVG asset stored under docs/
    When the user runs "clue validate" without network access
    Then it exits with code 0
    And it does not read, remove, or resolve the image targets
    # Single-direction because no image form is rejected any more: there is no
    # failing case to prove, and a second accepting test is not a negative one.

  @AC-094
  Scenario: Types carry the frontmatter fields their type requires
    Test-type: Unit
    Given a decision artifact missing author or accepted-by, or a capability missing goal
    When the user runs "clue validate"
    Then it exits with a non-zero code
    And the output names the file and the missing field
    But an empty accepted-by list is a present field, because an unsigned decision declares itself unsigned
    And a type carrying no extension requirement is checked against the core fields alone

  @AC-095
  Scenario: Milestone status cells follow one vocabulary
    Test-type: Unit
    Given a plan table whose header declares a Status column and whose row carries a value outside todo, doing, done, dropped
    When the user runs "clue validate"
    Then it exits with a non-zero code
    And the output names the plan, the milestone, and the value it read
    But a table that declares no Status column is not a milestone table, and a header row, a separator row, and an empty cell are never values
    And a table is read by its delimiter row, with or without outer pipes, with cells divided only by pipes that are unescaped and outside a code span

  @AC-101
  Scenario: clue id next allocates the next numeric ID through the ledger
    Test-type: Unit
    Given a repository whose identity ledger has been backfilled and whose counters map holds the last-issued numeric component for a prefix
    When the user runs "clue id next <prefix>"
    Then it prints the next sequential ID for that prefix as an increment of the stored counter, never a corpus scan
    And it persists the new entry as "reserved" and advances the counter
    But a canonical prefix absent from the ledger starts its counter at zero and issues the prefix's first ID; a malformed, lowercase, or improperly segmented prefix is rejected, and a repository with no ledger is told to run `clue migrate --apply` first

  @AC-108
  Scenario: clue id live promotes an allocated ID after its artifact is created
    Test-type: Unit
    Given the ledger marks an allocated ID "reserved"
    When the user runs "clue id live <id>" after creating the artifact
    Then it persists that ID as "live"
    But it rejects an ID that is missing or is already "live" or "retired"

  @AC-102
  Scenario: An opaque ID is preserved verbatim and never reused
    Test-type: Unit
    Given a ledger holding an opaque-kind entry with an exact canonical ID
    When a source mapping's documented generator proposes that same exact ID again
    Then the ledger rejects the ID as already used
    And it accepts a new opaque ID exactly as the generator produced it, without normalizing case, trimming, or reformatting

  @AC-103
  Scenario: checkLedger rejects a live artifact whose ID the ledger does not mark live
    Test-type: Unit
    Given a ".clue/id-ledger.yaml" that omits an ID or marks it "reserved" or "retired"
    And a live corpus artifact declaring that same ID
    When the user runs "clue validate"
    Then it exits with a non-zero code
    And the output names the file and says whether the ledger omits the ID or gives it a non-live state
    But a live artifact whose ID the ledger marks "live" passes, and a corpus with no ".clue/id-ledger.yaml" file is unaffected by this rule

  @AC-104
  Scenario: checkLedger rejects a malformed ledger entry shape
    Test-type: Unit
    Given a ".clue/id-ledger.yaml" holding an unknown kind, a numeric-kind entry whose ID, prefix, and component disagree, an opaque-kind entry carrying numeric fields, an invalid state, or duplicate canonical IDs
    When the user runs "clue validate"
    Then it exits with a non-zero code
    And the output names the entry and which shape or identity rule it breaks
    But a numeric entry whose canonical ID, prefix, and decimal component agree with a supported state, and an opaque entry with no numeric fields both pass

  @AC-105
  Scenario: An archived numeric ID above the live maximum is never reissued
    Test-type: Unit
    Given a ledger whose counter for a prefix sits below a "retired" entry's numeric component that once exceeded every live artifact's number
    When the user runs "clue id next <prefix>"
    Then the issued ID skips the retired number and every number the ledger already holds
    And "clue validate" rejects a live artifact later declaring that retired number
    But a live artifact declaring the newly issued number passes

  @AC-106
  Scenario: A UUID-like opaque ID cannot be reused after its source artifact is deleted
    Test-type: Unit
    Given a ledger holding a "retired" opaque entry with a UUID-like exact canonical ID, its source artifact having been deleted
    When a source mapping's generator proposes that same UUID-like ID again
    Then the ledger rejects it as already used
    And "clue validate" rejects a live artifact declaring that retired opaque ID
    But a freshly generated, distinct opaque ID passes both checks
```
