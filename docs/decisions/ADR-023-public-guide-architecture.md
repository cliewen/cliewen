---
id: ADR-023
type: decision
status: inferred
links: [G-003, P-003, CAP-001, PDR-009]
title: The public guide is an isolated VitePress site with a visibility-gated Pages deployment
author: agent
accepted-by: []
---

# ADR-023 — Public guide architecture

## Context and problem statement

[PDR-009](PDR-009-going-public.md) fixes the public guide as a handwritten VitePress site rather than a rendering of the corpus. The guide still needs a stable source boundary, URL base, build contract, diagram renderer, and deployment gate before its URLs become public. How does the guide remain useful to newcomers without coupling the permanent corpus to a site generator or attempting an unsupported private-repository deployment?

## Decision outcome

**The public guide is a self-contained VitePress source tree under `/guide`, built with repository-local Node dependencies, served from the `/cliewen/` GitHub Pages base path, and deployed from `main` only when GitHub reports that the repository is public.**

- **Source boundary:** `/guide` contains only newcomer-facing narrative pages and its `.vitepress` configuration. `/docs` remains the system-of-record and never becomes VitePress input. Guide pages link to stable repository paths when a reader needs the living corpus, while corpus validation remains independent of the site build.
- **Build ownership:** the root `package.json` and lockfile pin the documentation toolchain and expose `guide:dev`, `guide:build`, and `guide:preview`. The normal CI workflow runs a clean install and the production build on every pull request and push to `main`; VitePress dead-link failures remain enabled.
- **Published URL (superseded by [ADR-024](ADR-024-custom-domain-root.md)):** the site uses `/cliewen/` as its base because project Pages serves it below the repository name. A custom domain is outside P-003, so changing this base requires a later architectural decision with redirects for already-public URLs.
- **Diagrams:** guide diagrams are inline `mermaid` code fences rendered through `vitepress-plugin-mermaid`. The theme registers the renderer lazily so pages without diagrams do not preload the Mermaid runtime. Diagram source stays beside the prose, matching C-007's corpus convention without making the guide part of the corpus.
- **Deployment gate:** a dedicated Pages workflow listens for pushes to `main`, manual dispatch, and GitHub's `public` event. Its build job runs only when the event repository is public and the selected ref is `main`, so a manual run cannot publish an unaccepted branch. The deploy job alone receives Pages write and OIDC permissions and depends on that build. The ordinary CI build remains visibility-independent, so the guide is continuously verified while deployment is safely skipped and automatically becomes eligible when visibility flips.
- **Generated output:** VitePress cache and distribution directories are ignored build artifacts. GitHub Pages receives the built artifact; no generated site is committed to `main`.

### Rejected: use `/docs` as the VitePress source

That would couple corpus frontmatter, lifecycle conventions, and internal navigation to the public site's build, recreating the rendered-corpus option PDR-009 rejects.

### Rejected: commit generated HTML

Committed output duplicates source, obscures meaningful review with generated diffs, and gives the repository two candidates for the site's truth.

### Rejected: deploy from a long-lived publication branch

A publication branch introduces a second mutable history and delays the exact reviewed `main` state. GitHub Pages artifacts already provide the deployment boundary without another branch.
