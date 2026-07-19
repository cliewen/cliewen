---
id: CH-037
type: change
status: open
links: [P-003, M-011]
title: Let plain changes bypass Cliewen
---

# CH-037 — Let plain changes bypass Cliewen

Serves P-003/M-011. The v0.4.0 methodology must be cheap enough for adoption before the public visibility flip; the existing light tier still makes unrelated editorial work pay for CH identity, plan declaration, Cliewen verification, and full repository checks.

## What

- Add a plain-change route for work that changes no product behavior, intent, evidence, decision, plan, policy, or methodology carrier.
- Keep branch-from-main, pull-request review, and human merge universal while removing CH identity, corpus reads, proposal metadata, Cliewen skills, plan bookkeeping, changelog work, and the one-Cliewen-change-in-flight slot from plain changes.
- Run only checks relevant to a plain change's surface; in this repository, a guide-Markdown-only diff runs the guide build and whitespace check while mixed, protected, or unknown diffs fail closed to the full CI path.
- Update the corpus, routing hubs, canonical skill sources, generated skill trees, contributor surfaces, PR template, guide, and v0.4.0 release notes to carry the same boundary.

This is a full change because it changes methodology carriers and process decisions. It introduces no CLI or corpus-format API and changes no acceptance-criterion or capability meaning.
