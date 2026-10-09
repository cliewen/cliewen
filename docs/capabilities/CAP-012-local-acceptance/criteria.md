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

@AC-221
Scenario: Local acceptance is explicit and independent of identity allocation
  Test-type: Integration
  Given a repository whose acceptance mechanism has not been changed on its accepted base
  When local acceptance is requested
  Then it is refused unless both base and candidate declare the same local integration branch
  And source repositories remain PR-only
  And an opted-in adopter accepts without changing its ID ledger or coordination configuration
  And the confirmation states that verification declarations and human control are procedural claims
```
