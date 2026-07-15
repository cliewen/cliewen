---
id: ADR-019
type: decision
status: inferred
links: [P-002]
title: Validation requires foreign soil — trials on external repos, as findings not adoptions
author: agent
accepted-by: —
---

# ADR-019 — Validation requires foreign soil

## Context and problem statement

Cliewen has been proven only on repositories that share a maintainer and a mindset — this repo and one sibling, both built with the methodology in mind. Self-referential validation cannot reveal the failures that matter for adoption: a repo with no changelog culture, docs that fight the corpus taxonomy, a language without test-tag conventions, an AGENTS.md that contradicts the skills. How does the campaign get honest evidence before real adopters arrive?

## Decision outcome

**The plan carries a foreign-soil milestone: the skills are trialed on at least two external open-source repositories, and the trials produce findings, not adoptions.**

- The trial repos are selected by a human and are deliberately foreign — no shared maintainer, not designed for the methodology. Selection is evidence, not plan structure: the milestone lands with repos unchosen.
- Each trial ends in an `AN-xxx` findings doc (the analysis skill's existing contract: no findings doc, no trial). The milestone closes only when at least one methodology adjustment traces back to trial findings — a trial that changes nothing was not looked at hard enough, or the methodology is more finished than assumed; either is a finding worth recording.
- **Trials are read-and-apply experiments, not adoptions**: no PRs against the foreign repos, no new extraction mappings (P-002 keeps non-OpenSpec mappings out of scope). A trial that demands a new mapping or tool surfaces that demand as a finding; acting on it is a later plan decision.

**Carrier:** the M-007 row in P-002 (the exit criterion is the rule); the analysis skill's findings-doc contract already covers the trial output. Repo-local — nothing ships to adopters.

### Rejected: waiting for organic adopters to supply the evidence

Real adopters arrive after the methodology already failed them silently — there is no feedback loop unless the first foreign contact is deliberate, instrumented, and cheap to repeat.

### Rejected: a full adoption of a foreign repo as the first test

Conflates two questions: whether onboarding tooling works (M-005's job) and whether the methodology's assumptions survive foreign ground. A failed full adoption cannot say which one broke; cheap trials isolate the second question first.
