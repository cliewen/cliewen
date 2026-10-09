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

@AC-225 @retired
Scenario: Local acceptance is explicit and independent of identity allocation
  # Retired by PDR-069; superseded by AC-228.

@AC-223
Scenario: Fresh adoption defaults to local and PR is actively selected
  Test-type: Integration
  Given a fresh adoption with no prior Cliewen state
  When init runs without an acceptance option
  Then it materializes local acceptance on main
  And an explicit PR option instead materializes PR acceptance
  And repeating init preserves existing policy
  And a conflicting explicit option is refused before any writes

@AC-224 @retired
Scenario: Existing adoption retains acceptance during migration and init
  # Retired by PDR-069; superseded by AC-227.

@AC-226 @retired
Scenario: Policy parsing and source restrictions are shared and strict
  # Retired by PDR-069; superseded by AC-229.

@AC-227
Scenario: Missing policy defaults to local regardless of adoption history
  Test-type: Integration
  Given a repository with existing Cliewen state and no acceptance policy
  When init or migration resolves acceptance
  Then init materializes local acceptance on main and migration leaves policy absent
  And existing explicit local and PR policies remain byte-identical
  And preview writes nothing and malformed policy blocks apply

@AC-228
Scenario: Effective local policy is independent of identity allocation
  Test-type: Integration
  Given base and candidate with equivalent effective local policies
  When a human requests local acceptance
  Then absent policy defaults to local on main and explicit local/main is equivalent
  And differing policies or source repositories are refused
  And acceptance leaves the ID ledger and coordination configuration unchanged
  And confirmation states the procedural limits of human presence and verification claims

@AC-229
Scenario: Shared strict policy parsing uses the local default
  Test-type: Unit
  Given acceptance configuration consumed by init, migration or local acceptance
  When its policy is read
  Then absence means local acceptance on main
  And an explicit file requires local or pr with a valid explicit branch
  And malformed data, unknown fields and multiple YAML documents are refused
  And explicit source policy permits PR only
```
