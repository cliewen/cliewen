---
id: CAP-005-criteria
type: criteria
status: active
links: [CAP-005]
title: Acceptance criteria for CAP-005
---

# Acceptance criteria — CAP-005

```gherkin
Feature: Index generation — clue scaffold

  @AC-026
  Scenario: scaffold regenerates index blocks and keeps prose
    Given a corpus with an artifact added after the last index regeneration
    When the user runs "clue scaffold"
    Then the folder README's index block references the new artifact
    And prose outside the clue:index markers is unchanged
    And "clue validate" accepts the result
    And a run with nothing new to index changes no file

  @AC-027
  Scenario: scaffold touches only index blocks and materializes nothing
    Given a corpus with an artifact to index
    When the user runs "clue scaffold"
    Then no file is created or deleted and no file outside the taxonomy READMEs is modified
    And on a path without a docs tree the command exits non-zero, names the problem, and creates nothing

  @AC-073 @retired
  Scenario: An appended index row states the record it links
    Test-type: Unit
    Given a taxonomy folder holding an artifact no index row references
    When the index block is regenerated
    Then the appended row states the artifact's id, its title, and its status
    And a title whose value contains a colon is spelled as the parsed value, not as its YAML quoting
    But an artifact missing any of a readable id, title, and status falls back to the plain link, rather than a row carrying an empty status
    And a row referencing a subfolder README carries no title or status, because it states a section rather than a record
    # Retired 2026-08-12 (CH-152): a constraint's register reader needs its
    # enforcement class in this position rather than its generic lifecycle
    # status; AC-138 carries the type-aware row and report contract.

  @AC-138
  Scenario: A constraint index badge states enforcement and a mismatch is visible
    Test-type: Unit
    Given a taxonomy folder holding an unindexed constraint with a readable enforcement class
    And an existing direct constraint index row whose badge differs from its target's enforcement class
    When the index block is regenerated and the user runs "clue validate --index-rows"
    Then the appended constraint row states its id, title, and enforcement class in the badge position
    And an appended row for an artifact other than a constraint still states its status in that position
    And the validation run exits zero, counts the mismatched constraint badge, and names the README and target a reader opens to
    But a constraint missing readable enforcement falls back to the plain link rather than carrying an empty or status badge
    And a matching direct constraint row, a subfolder row, and a curated row covering several targets are not counted as mismatches

  @AC-096
  Scenario: An appended index row says what the artifact is about
    Test-type: Unit
    Given a taxonomy folder holding an artifact no index row references
    When the index block is regenerated
    Then the appended row carries a description seeded from the artifact's own body after its status
    And a lede paragraph directly beneath the H1 is preferred over prose under a later heading
    And a body with no lede is read from the first paragraph under its first heading
    And structure is not read as that paragraph: a heading written either way, a table row, a bulleted or ordered list item, a blockquote, a fenced or indented code block, an HTML block, and a horizontal rule
    And every inline link in the seeded sentence is reduced to its label, so the row cannot cover a second target
    And a link written inside a code span is reduced too, because the block's link reading knows nothing about spans and a quoted placeholder target resolves to nothing
    And a sentence longer than the bound is cut at a word boundary, backing off to a rune boundary where the prose carries no spaces, so the row is never invalid UTF-8
    But a candidate that cannot be made safe is declined rather than repaired, and the shorter row is emitted instead: a link label carrying its own brackets defeats reduction, and an index marker or HTML comment would truncate the block it is written into

  @AC-097
  Scenario: A row is one shape or the other and never carries an empty description
    Test-type: Unit
    Given an artifact whose body holds no readable prose sentence
    When the index block is regenerated
    Then the appended row is exactly the row that states its record, with no trailing separator and no empty tail
    But an artifact whose frontmatter cannot be read degrades to the plain link, acquiring neither a status badge nor a description, even where its body holds a readable sentence

  @AC-098
  Scenario: A curated description outlives regeneration
    Test-type: Unit
    Given a taxonomy README whose row carries a description an author corrected by hand
    When the index block is regenerated
    Then the row is unchanged, because the seed is a first draft and never an assertion
    And no command writes a description into a row that already exists
    But a curated row whose target no longer resolves is dropped like any other, so a description cannot keep a dangling entry alive
```
