---
id: CH-047
type: change
status: open
links: [P-004, M-014, CAP-001, QS-002, C-012]
title: Make the CI wall and minimum adoption practice actionable
---

# CH-047 — Make the CI wall and minimum adoption practice actionable

## What

Complete P-004/M-014 with a dedicated CI guide that turns the generated wall from a warning into an enforced merge gate. Give exact release-binary vendoring and checksum steps, explain armed and unarmed behavior, show why `clue validate --forbid-changes` is the merge command, document the GitHub ruleset settings that protect `main`, and include a disposable failing-pull-request probe that demonstrates the required check blocking merge. State the equivalent enforcement contract for other forges without pretending their user interfaces are identical.

Revise the adoption guide around the minimum useful practice: goal → capability → acceptance criterion → focused positive and negative evidence. Explain when the wider corpus becomes useful and when Cliewen is a poor fit. Connect the new guide from the existing onboarding path, update CAP-001's accumulating guide requirements and the user-facing changelog, then close M-014 in the digest when the commands and probe have been rehearsed and the strict guide build passes.

## Why

The current guide says to arm the wall but leaves the exact vendoring commands inside a generated workflow comment and describes branch protection only in principle. A newcomer can therefore finish the safe demo without a reliable path from visible CI to an enforced merge boundary. The adoption page also explains greenfield and brownfield entry without stating the smallest practice worth adopting or the cases where the method's overhead is the wrong trade.

## Decision boundary

This change implements the already accepted M-014 scope. It does not alter `clue`, validation semantics, positive/negative-pair machine enforcement, skills, scaffold templates, generated workflows, QS-002, or hosted repository settings. GitHub instructions must match current supported ruleset behavior; other forges receive only the provider-neutral contract. Any need to change a product or methodology carrier, weaken a check, or make a new architectural or process decision becomes an open question and stops the change.
