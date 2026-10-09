---
id: CAP-012-criteria
type: criteria
status: active
links: [CAP-012]
title: Local acceptance criteria
---

# Criteria

```gherkin
@AC-219
Scenario: A human accepts an exact candidate with recoverable provenance
  Test-type: Integration
  Given an explicitly opted-in adopter with a digested verified candidate containing its tracked proposal history
  When a human confirms local acceptance of that candidate and its current base
  Then the acceptance commit has the base and candidate as its ordered parents and exactly the candidate tree
  And the complete brief and proposal remain recoverable after deleting the candidate branch
  And an incomplete or revision-mismatched brief cannot accept a candidate

@AC-220
Scenario: Preflight and refusal never integrate a candidate
  Test-type: Integration
  Given a proposed local acceptance
  When preflight succeeds without confirmation or the human cancels
  Then no branch, index, or tracked file changes
  And stale base, dirty checkout, invalid corpus, undigested workspace, missing proposal history, or noninteractive acceptance is refused
  And a changed base or dirty checkout during confirmation cannot be accepted

@AC-221 @retired
Scenario: Local acceptance requires an explicit initial opt-in
  # Retired by CH-205: PDR-068 makes new adoption local by default while preserving declared policy and legacy PR.

@AC-225
Scenario: Local acceptance is explicit and independent of identity allocation
  Test-type: Integration
  Given a repository whose acceptance mechanism has not been changed on its accepted base
  When local acceptance is requested
  Then it is refused unless both base and candidate declare the same local integration branch
  And source repositories remain PR-only
  And an adopter with declared local policy accepts without changing its ID ledger or coordination configuration
  And the confirmation states that verification declarations and human control are procedural claims

@AC-223
Scenario: Fresh adoption defaults to local and PR is actively selected
  Test-type: Integration
  Given a fresh adoption with no prior Cliewen state
  When init runs without an acceptance option
  Then it materializes local acceptance on main
  And an explicit PR option instead materializes PR acceptance
  And repeating init preserves existing policy
  And a conflicting explicit option is refused before any writes

@AC-224
Scenario: Existing adoption retains acceptance during migration and init
  Test-type: Integration
  Given an existing Cliewen repository with implicit PR or an explicit acceptance policy
  When init or a reviewed migration resolves its policy
  Then implicit PR becomes explicit PR without changing the accepted workflow
  And existing valid local and PR policies remain byte-identical
  And preview writes nothing and malformed policy blocks apply
  And markerless Cliewen corpora are not mistaken for fresh adoption

@AC-226
Scenario: Policy parsing and source restrictions are shared and strict
  Test-type: Unit
  Given acceptance configuration consumed by init, migration or local acceptance
  When its policy is read
  Then only local or pr with a valid explicit branch is accepted
  And malformed data, unknown fields and multiple YAML documents are refused
  And source role permits PR only
  And missing configuration preserves the historical PR convention
```
