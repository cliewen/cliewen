---
id: CH-032
type: change
status: open
links: [P-003, G-003, CAP-001]
title: A deployment-ready public guide for Cliewen
---

# CH-032 — A deployment-ready public guide

## What

Build the handwritten VitePress guide required by P-003/M-010. The site will explain Cliewen's methodology, corpus taxonomy, change loop, and agent skills in newcomer-facing prose; render its diagrams with Mermaid; fail its local and CI build on dead links; and carry a GitHub Pages workflow that remains skipped while the repository is private.

## Why

P-003/M-010 is the next unfinished public-readiness milestone. The corpus and root README serve practitioners already inside the method, but G-003 requires a stranger to learn Cliewen without reading the corpus first. PDR-009 fixes the guide's form and the publication sequence: the site must be deployment-ready before v0.4.0, while actual publication waits for the final visibility flip.

## Decision boundary

The site's structural choices are expensive to reverse once URLs are public, so this change records ADR-023 for the guide root, repository-relative base path, corpus boundary, build ownership, and visibility gate. The VitePress choice and handwritten-not-rendered boundary are already fixed by PDR-009. The change does not alter CLI or methodology capability semantics and therefore adds no acceptance criterion; the M-010 exit criterion is verified by the production guide build in local development and CI.

## Constraints and quality scenarios

The guide and corpus prose follow C-001, user-visible impact is recorded under `[Unreleased]` per C-002, tasks tick immediately per C-003, and no existing check is weakened per C-004. The proposal serves P-003/M-010 per C-005; ADR-023 is timeless and correctly typed per C-006, C-009, and C-011; diagrams are inline Mermaid per C-007; completed plans remain untouched per C-008; P-003 retains its `todo`/`done` vocabulary per C-010; and the branch/PR remains inside C-012's review boundary. QS-001 remains enforced by the unchanged Go coverage gate. QS-002 remains intact because the root README quickstart stays the canonical shortest path and the guide adds no prerequisite before a first green validate.
