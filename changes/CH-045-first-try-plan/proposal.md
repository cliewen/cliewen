---
id: CH-045
type: change
status: open
links: [G-001, P-004]
title: P-004 "Cliewen earns the first try" begins
---

# CH-045 — P-004 "Cliewen earns the first try" begins

## What

Create P-004 as the active campaign serving G-001 through three sequential milestones: make Cliewen's value, release-binary installation, and a safe validation failure concrete; explain how to arm and enforce the CI wall and what minimum viable practice requires; then document compatibility boundaries, maintenance, troubleshooting, evidence, and clear next actions.

Record the campaign's deliberately narrow scope in the decision log. The campaign changes the public guide and its durable onboarding contract, not `clue`, the validator, skills, scaffold templates, package-manager distribution, or QS-002's existing 30-minute promise.

## Why

The public guide explains the methodology in depth but asks a newcomer to absorb its vocabulary and initialize their own repository before showing the tool catch a concrete defect. Installation begins with the Go toolchain, the generated CI workflow hides its arming procedure in comments, and current support and maintenance boundaries are scattered across the corpus and release notes. G-001's verifiable thread is only useful when a newcomer can recognize the problem, try the judge safely, and understand the enforced path without mistaking intended conventions for implemented checks.

## Plan relationship

This is the plan-creating change and serves G-001 directly. After acceptance, M-013 through M-015 are implemented as separate Cliewen changes, one ready PR and human merge at a time.

## Decision boundary

This change decides only the campaign structure, milestone promises, and explicit deferrals. It changes no product behavior, acceptance criterion, quality scenario, methodology carrier, or completed plan. Any requirement to change the validator, add an installer or package-manager channel, claim unsupported compatibility, or promise a shorter timed journey becomes an open question and stops the change.
