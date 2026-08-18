---
id: PDR-009
type: decision
status: verified
links: [G-003]
title: Going public is a campaign — readiness first, one release, then the flip
author: agent
accepted-by: Flemming N. Larsen (2026-07-18, planning conversation — visibility goal, guide form, and release sequencing chosen explicitly)
---

# PDR-009 — Going public is a campaign

## Context and problem statement

Making the repository public is effectively irreversible, and the first public state must not expose an insider-only corpus, a stale install story, or an absent contributor front door. Public release therefore needs a deliberate readiness boundary.

## Decision outcome

**Readiness changes land first, release v0.4.0 is cut from them, and then repository visibility flips while the guide goes live in the same act.** Readiness includes a stranger-safe corpus, contribution, conduct, and security paths, and a deployment-ready guide whose build and dead-link checks pass.

The guide is a handwritten VitePress site that explains the methodology and links to the living corpus; it is not a verbatim rendering of corpus artifacts. Normative corpus citations must resolve publicly, while historical evidence naming private adopters remains identified as historical evidence. The campaign milestones and `CONTRIBUTING.md` carry this repository-local sequence.
