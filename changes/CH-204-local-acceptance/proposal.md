---
id: CH-204
type: change
status: open
links: [G-016, G-001, ARCH-003]
title: Human-controlled local acceptance without a pull request
---

# CH-204 — Local acceptance

This change is plan-less: it implements the plan accepted in conversation on 2026-10-08, including the final instruction to keep the existing identity ledger and coordination behavior unchanged. It serves G-016 under the draft, inferred VIS-001; this change does not confirm the vision.

## Contract

Adopters may explicitly select local acceptance. Agents prepare and verify a tracked change branch; a human runs `clue accept <candidate-sha> --base <base-sha> --brief <file>` to review and confirm an exact candidate. `--check` is read-only. The acceptance merge preserves the base and candidate as parents, the candidate tree, and the complete brief in its message. PR acceptance remains the default; this source repository still requires a PR and human merge.

The brief binds change ID, candidate/base identities, criteria and evidence, verification and review declarations. The command validates the isolated candidate and available proposal history, refuses stale or dirty integration state, and never resolves conflicts or pushes. Human control is procedural, not authenticated by the prompt or Git identity. No unattended acceptance option is provided. Identity allocation, test classes, and the deterministic judge's repository-state boundary remain unchanged.

## Challenge the commitment

The riskiest assumption is that local acceptance reduces friction while supporting an informed decision on a workstation whose agent can also run Git. A credible alternative is to retain PR-only acceptance or require separate credentials. The cheapest useful test is a disposable repository exercising successful acceptance, rejection of a green but unwanted candidate, cancellation, stale-state refusal, and recovery of the proposal and brief. Automated tests can demonstrate the mechanics; an unfamiliar maintainer's usability trial remains explicitly human evidence. Stop or revise if the flow cannot bind the reviewed tree without races, loses provenance, or claims to enforce human identity. A fully green implementation could still fail its user if its brief hides the important behavior or replaces PR ceremony with equally burdensome manual bookkeeping.

## Carrier inventory

Implementation: CLI dispatch/help, acceptance package, focused Git integration tests and evidence export. Methodology: canonical lifecycle review boundary, change loop, verification, planning/extraction references as affected, onboarding and agent hub templates, PR acceptance wording, local brief template. Durable truth: core/architecture/design overviews, CAP-012, PDR-067, affected existing acceptance criteria and decisions/constraints. Public guidance: README and guide passages asserting PR-only adoption. Generated skills and scaffold carriers are regenerated; local contributor/CI policy stays PR-based. CHANGELOG records the shipped addition.

## Acceptance

AC-219: exact-candidate local acceptance and durable provenance. AC-220: refusal and non-mutating preflight/cancellation. AC-221: declared opt-in and procedural human boundary with PR and identity compatibility. The final PR must report the human trial honestly rather than inventing a maintainer observation.
