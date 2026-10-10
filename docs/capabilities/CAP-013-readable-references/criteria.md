---
id: CAP-013-criteria
type: criteria
status: active
links: [CAP-013]
title: Readable-reference presentation and instruction criteria
---

# Criteria

```gherkin
@AC-232
Scenario: Standalone skills carry one reachable readable-reference instruction
  Test-type: Unit
  Given any independently distributed Cliewen lifecycle skill
  When its entry point and required references are generated
  Then authoring human-facing prose routes to the canonical readable-reference instruction
  And the instruction explains inline meaning, descriptive links, secondary IDs and literal identity exceptions
  And missing or drifted distributed guidance is rejected

@AC-233
Scenario: Human reports name local references without changing underlying data
  Test-type: Unit
  Given artifacts, criteria and milestones with canonical identities and readable names
  When human CLI reports render them
  Then local references show readable names and secondary identities with unchanged states and ordering
  And metadata, raw corpus content, machine-facing identities and validation outcomes are unchanged
  And absent, duplicate or unnamed references are described without guessing
  And terminal controls and embedded line breaks cannot create extra label rows

@AC-234
Scenario: Human diagnostics name only explicitly local subjects
  Test-type: Unit
  Given a diagnostic whose producing check records a local reference subject
  When its human presentation is rendered
  Then the subject is described by its local name or explicit uncertainty
  And raw finding paths and messages remain unchanged
  And same-spelled Cliewen source rules, foreign references and tokens inside inserted names are not relabeled as local subjects

@AC-235
Scenario: A reader understands a handoff without opening identity references
  Test-type: Human
  Given the self-contained acceptance summary and representative repaired text
  When the owner reads it without following the links or looking up identities
  Then the owner can explain the purpose, relevant reference meaning and required action
```
