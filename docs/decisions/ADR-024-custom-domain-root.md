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

The custom domain makes ADR-023's `/cliewen/` base generate broken root-domain asset and navigation URLs.

## Decision outcome

**The public guide is canonical at `https://cliewen.dev/`, and VitePress builds with `base: "/"`.** GitHub Pages redirects legacy `https://cliewen.github.io/cliewen/<path>` links to the custom-domain root while preserving the page suffix. ADR-023's source, build, diagram, deployment, and generated-output decisions remain in force; only its published-URL clause is superseded. Keeping the project base or publishing duplicate root and project-path trees is rejected because either breaks the root or creates competing canonical URLs.
