---
id: IDR-006
type: decision
status: verified
links: [ADR-018, ADR-028, CAP-001]
title: Init never writes through linked directories
author: agent
accepted-by: Flemming N. Larsen (2026-09-02, conversation)
---

# IDR-006 — Init never writes through linked directories

## Context

A linked directory below an initialization target can be shared by multiple checkouts, so writing through it changes content outside the repository being initialized.

## Decision

`clue init` skips any emitted subtree with a symlinked ancestor below the target root, reports it as `linked`, and offers no override. The target root itself is not inspected, and validation's deliberate link-following remains unchanged.
