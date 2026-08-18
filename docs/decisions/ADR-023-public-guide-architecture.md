---
id: ADR-023
type: decision
status: verified
links: [G-003, P-003, CAP-001, PDR-009]
title: The public guide is an isolated VitePress site with a visibility-gated Pages deployment
author: agent
accepted-by: Flemming N. Larsen (2026-08-02, conversation)
---

# ADR-023 — Public guide architecture

## Context and problem statement

The public guide needs a stable source boundary, build contract, URL base, diagrams, and a safe deployment gate without coupling the permanent corpus to a site generator.

## Decision outcome

**Use a self-contained VitePress source tree under `/guide`, built with repository-local Node dependencies and deployed from `main` only when GitHub reports the repository public.** `/docs` remains the system of record. Root scripts expose development, build, and preview commands; CI builds with dead-link failures enabled. The project Pages base is `/cliewen/` until ADR-024's custom-domain decision applies.

Guide diagrams are inline Mermaid rendered by the VitePress plugin. The Pages workflow handles pushes to `main`, manual dispatch, and the public event; only a public `main` ref can deploy, with write and OIDC permissions limited to the deploy job. Generated VitePress output is never committed. Using `/docs` as input, committing HTML, or deploying from a long-lived publication branch is rejected because each creates a second mutable or generated source of truth.
