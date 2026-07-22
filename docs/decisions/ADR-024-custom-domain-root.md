---
id: ADR-024
type: decision
status: verified
links: [ADR-023, G-003, P-004, CAP-001]
title: The public guide is canonical at the cliewen.dev root
author: human
accepted-by: Flemming N. Larsen (2026-07-22, Codex review conversation)
---

# ADR-024 — The public guide is canonical at the cliewen.dev root

## Context and problem statement

The public guide now has the custom domain `cliewen.dev`. The `/cliewen/` VitePress base chosen for project Pages in ADR-023 makes the custom-domain root generate asset and navigation URLs below a path that does not exist. The guide needs one canonical URL while preserving links published before the custom domain was attached.

## Decision outcome

**The public guide is canonical at `https://cliewen.dev/`, and VitePress builds it with `base: "/"`.**

- **Canonical root:** assets, navigation, search data, and page links are emitted below `/` on `cliewen.dev`.
- **Legacy links:** GitHub Pages owns the redirect from `https://cliewen.github.io/cliewen/<path>` to `https://cliewen.dev/<path>`. The redirect preserves the page suffix while removing the obsolete project path, so existing public guide links continue to resolve without a second site tree.
- **Deployment boundary:** ADR-023's source, build, diagram, deployment, and generated-output decisions remain in force. This decision supersedes only its published-URL clause.

### Rejected: retain the project base on the custom domain

Keeping `/cliewen/` would make the canonical root return HTML whose scripts and styles resolve to missing paths, leaving the page unrendered.

### Rejected: publish duplicate root and project-path builds

Two generated site trees would create competing canonical URLs and duplicate every page. GitHub Pages already redirects the former project URL to the custom-domain root.
