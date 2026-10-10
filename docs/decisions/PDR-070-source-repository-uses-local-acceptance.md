---
id: PDR-070
type: decision
status: inferred
links: [G-016, CAP-012, CAP-004, PDR-069, PDR-015, C-012]
title: Source repositories follow acceptance policy and Cliewen uses local acceptance
author: agent
accepted-by: []
binds: adopter
---

# Source repositories use the same acceptance policy

The [recorded human approval policy](PDR-072-recorded-human-approval-allows-delegated-local-execution.md) permits delegated execution of an exact human decision and safe internal link snapshots. Its human decision boundary supersedes the earlier requirement that the human personally run the local command.

## Context

The owner selected local acceptance for Cliewen itself. A source-only PR restriction in the policy reader and local command prevented that choice even with explicit local configuration. Repository role identifies which methodology carriers apply; it need not dictate the acceptance mechanism.

## Decision

Source and adopter repositories use the same acceptance-policy reader, init selection and local preflight. Remove the source-only PR refusal. Preserve role-marker validation, strict policy parsing, agreement between effective base and candidate policies, exact candidate history, human confirmation and identity coordination. The shipped carrier is `internal/skills/source/shared/review-boundary.md.tmpl`; the policy reader and local command enforce that shared behavior.

Cliewen keeps `role: source` and selects `mode: local`, `branch: main`. Tracked work uses a human-run `clue accept` after verification and review. Direct work, including the administrative release cut, uses a human-controlled local merge after applicable checks. A release remains outside the tracked route; accepted main triggers the existing version, tag and publication gates. Release PR creation is optional. An agent never accepts its own work.

This supersedes the source-only PR exceptions in PDR-067 through PDR-069 and the mandatory release-PR handoff in PDR-015. A policy transition must be accepted under the previous accepted policy: the current PR-to-local candidate cannot accept itself using its new configuration. Complete that transition through the existing human-controlled PR boundary, then use local acceptance for subsequent work.

## Risk and authorization

The riskiest assumption is that local acceptance can replace hosted admission checks without hiding what was enforced. Local verification and review are required, but the local command does not enforce hosted CI or authenticate human presence. State those procedural limits in the handoff. The owner separately approved updating GitHub ruleset `protect-main-history` after acceptance of the transition: remove mandatory PR and pre-push status-check admission, retain deletion and non-fast-forward protection with no bypass actors, and keep CI workflows and local verification. Local acceptance does not provide hosted enforcement before integration; CI still judges published main. The owner explicitly selected local acceptance here, and the implementation follows the previously authorized direct route with an override recorded in Git history.
