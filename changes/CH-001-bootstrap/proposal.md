---
id: CH-001
type: change
status: open
links: [P-001, M-001, G-001]
title: Bootstrap the Cliewen corpus
---

# CH-001 — Bootstrap the Cliewen corpus

## What

Hand-scaffold the permanent `/docs` corpus and the method machinery
(`AGENTS.md`, `.agents/skills/`), executed as a Cliewen change loop:
this branch is the proposal, the merge is the acceptance, and this
folder is deleted in the digest commit before merge.

## Why

The Foundation Document (v0.4) mandates dogfooding from commit one: the
Cliewen repo runs Cliewen conventions before the `clue` CLI exists. The
hand-scaffolded taxonomy doubles as the future template source for
`clue init` — nothing here is throwaway.

## Serves

P-001 (elaboration baseline), milestone M-001: the change loop closes
once, end-to-end. This change is itself the first evidence.

## Scope

Docs only. No Go code, no CI workflow — those start with CH-002 so the
first delta stays small and the docs contract is reviewed on its own.
