---
id: IDR-007
type: decision
status: inferred
links: [ADR-010, CAP-002, CAP-003]
title: Frontmatter uses YAML and rejects hidden or duplicate blocks
author: agent
accepted-by: []
---

# IDR-007 — Frontmatter uses YAML and rejects hidden or duplicate blocks

## Context

The deterministic judge must accept real YAML while refusing byte and document shapes that conceal which metadata block is authoritative.

## Decision

Frontmatter is parsed with `gopkg.in/yaml.v3`. Validation rejects a UTF-8 BOM anywhere in corpus Markdown and a complete second frontmatter block opening an artifact body; extraction replaces source frontmatter and leaves exactly one block.
