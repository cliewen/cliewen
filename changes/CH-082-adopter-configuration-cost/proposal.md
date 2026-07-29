---
id: CH-082
type: change
status: open
links: [P-008]
title: A real adopter prices configuration before Cliewen chooses an interface
---

# CH-082 — A real adopter prices configuration before Cliewen chooses an interface

## What

Measure, on the Robocode Tank Royale repository, what Cliewen's current adopter-facing assumptions actually cost the repository that lives with them. Three assumption families are in scope: the vendored CI wall's runner, action, download and binary-location choices; the placement of agent-facing directories (`.agents`, `.claude`) against what already occupies them; and the need for repository-local skills beside the managed `clue-*` set. The spike reproduces each assumption, records the concrete edits or failures the repository is forced into, separates what the repository genuinely needs from what its maintainer merely prefers, and weighs candidate configuration locations against [ADR-013](../../docs/decisions/ADR-013-ships-generic-vs-repo-local.md)'s existing boundary.

The change delivers one analysis record under `/docs/analysis`, explicit ADR/PDR candidates routed by reversal cost, a named successor-plan consumer, and the M-035 bookkeeping and any user-facing release note in its digest.

## Why

[ADR-013](../../docs/decisions/ADR-013-ships-generic-vs-repo-local.md) rejected a `cliewen.yaml` configuration file and left a precise door behind: "a machine config earns its place only when `clue` itself needs repo-local settings." P-008/M-035 exists to find out whether that condition is now met, from a real repository rather than from speculation. Choosing a configuration interface before pricing the need would violate the campaign's own working order — findings before interfaces — and would risk paying ADR-013's stated cost, a second source of truth for routing, to solve a problem no adopter has demonstrated.

Tank Royale is the human-selected adopter: it is the first product repository to adopt Cliewen, its adoption is merged and has lived history since, and it is public, so every observation here is independently reproducible.

## Scope

This is a full change serving P-008 milestone M-035. Per that milestone it implements **no** configuration file, mirror option, or extension mechanism, and it changes nothing inside Tank Royale — the working copy is read-only evidence, and any reproduction that would mutate it runs on a disposable clone. It does not open the distributed-work and cross-repository-evidence questions reserved for M-036, and it does not decide the candidates it ends with; those are routed for a successor plan to consume.

## Evidence discipline

M-033's `clue-analysis` obligations bind this spike. Every reproduced result is classified as a clean disposable environment or a prepared one, and any local prerequisite makes it prepared and disqualifies it as onboarding-reproducibility evidence. Repository activity is evidence of activity, not maintainer intent. Because the human who selected the adopter also maintains it, the mandatory-versus-preference split cannot be settled from the repository alone: ambiguous cases become open questions rather than inferred intent.
