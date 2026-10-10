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

@AC-220 @retired
Scenario: Preflight and refusal never integrate a candidate
  # Superseded by AC-238: recorded human approval permits delegated execution; refusal and fresh-state checks remain.

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

@AC-228 @retired
Scenario: Effective local policy is independent of identity allocation
  # Retired by PDR-070; superseded by AC-230.

@AC-229 @retired
Scenario: Shared strict policy parsing uses the local default
  # Retired by PDR-070; superseded by AC-231.

@AC-230
Scenario: Effective local policy works independently of repository role and identity allocation
  Test-type: Integration
  Given base and candidate with equivalent effective local policies in an adopter or source repository
  When a human requests local acceptance
  Then absent policy defaults to local on main and explicit local/main is equivalent
  And differing policies or an accepted PR policy are refused
  And acceptance leaves the ID ledger and coordination configuration unchanged
  And confirmation states the procedural limits of human presence and verification claims

@AC-231
Scenario: Shared strict policy parsing is independent of repository role
  Test-type: Unit
  Given acceptance configuration consumed by init, migration or local acceptance
  When its policy is read in an adopter or source repository
  Then absence means local acceptance on main
  And an explicit file requires local or pr with a valid explicit branch
  And malformed data, unknown fields and multiple YAML documents are refused
  And init supports the local default and explicit local or PR in either role
@AC-236
Scenario: Committed internal links are materialized safely and completely
  Test-type: Integration
  Given a committed revision containing internal file, directory or chained links
  When local acceptance constructs its isolated snapshot
  Then validation sees the committed targets as ordinary snapshot content without using checkout files or OS links
  And the original candidate link modes remain unchanged
  And external, missing, metadata, cyclic, excessive or submodule targets are refused

@AC-237
Scenario: An exact recorded human decision can be executed without a terminal prompt
  Test-type: Integration
  Given a reviewed candidate, complete brief and human approval record bound to candidate, base and brief hash
  When delegated local execution consumes the record
  Then it creates the exact-candidate-tree merge with ordered parents and retains the full brief and approval provenance
  And preflight alone does not integrate
  And missing, malformed, rejected or mismatched records cannot integrate
  And recorded provenance does not claim to authenticate human presence

@AC-238
Scenario: Preflight, refusal and fresh-state checks protect both approval paths
  Test-type: Integration
  Given interactive confirmation or an exact recorded human decision
  When preflight runs, confirmation is cancelled or inputs become stale
  Then preflight and refusal do not integrate
  And dirty checkout, advanced base, changed brief, invalid corpus, undigested workspace or missing proposal history is refused
  And noninteractive execution without an approval record is refused
  And changed approval bytes cannot be executed as the previously checked record
```
