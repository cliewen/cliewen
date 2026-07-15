---
id: PDR-006
type: decision
status: verified
links: [PDR-003, ADR-010]
title: Decision records are typed — ADRs for architecture, PDRs for project/process, log rows for the cheap
author: agent
accepted-by: Flemming N. Larsen (2026-07-15, PR #13 approval)
---

# PDR-006 — Decision records are typed

## Context and problem statement

ADR is the industry's one well-known decision record, and its name promises *architectural* decisions. When a corpus files process rules — change tiers, acceptance mechanics, trial plans — under ADR, the well-known term starts lying to every reader who knows it. The industry has already run this experiment: the MADR template project generalized its name from "Markdown **Architectural** Decision Records" to "Markdown **Any** Decision Records" in version 3.0 and renamed it **back** in 4.0. Yet the expensive-to-reverse process decisions still need a full record — context, options, rejected alternatives — that a one-line log row cannot carry. What record types does the corpus need, and what routes a decision to each?

## Decision outcome

**Three record types, routed by two questions: is reversing the decision cheap and local? If yes, a log row. If no, is it about architecture or about how the project works? Architecture is an ADR; the rest is a PDR.**

- **ADR — Architectural Decision Record**, in Nygard's strict sense: decisions about the structure of the software and of the corpus format (language, ID schemes, frontmatter contracts, lint rules, what ships to adopters). The well-known term keeps its well-known meaning.
- **PDR — Project/Process Decision Record**: expensive-to-reverse decisions about how the project works — change tiers, decision acceptance, validation strategy. Same MADR template, same two-tier provenance ([ADR-010](ADR-010-provenance-field.md)), same folder, own ID series.
- **Decision log** — unchanged from [PDR-003](PDR-003-decision-log.md): one dated row for the cheap-and-local-to-reverse.
- **Well-established practices are cited, not re-derived.** A decision that adopts a practice the industry already settled names it (with a source) and records only what is local: why it is adopted here and where this corpus deviates. The record earns its length by local reasoning, not by restating textbooks.
- Reclassification is ordinary change content, applied retroactively: a record filed under the wrong type is renamed into the right series with its text intact; git history keeps the provenance of the rename.

**Carrier:** the decisions folder README prose (ships as `clue init` template prose); the routing wording in the `clue-delta` and `clue-verify` skills (agent); the register entry [C-011](../constraints/C-011-decision-records-typed.md) that holds the routing rule.

### Rejected: one generic decision-record series for everything

The MADR project's 3.0 → 4.0 round trip is the evidence: "any decision record" erases the signal the ADR name carries, and the community pulled it back. It would also renumber every existing record for no gain in meaning.

### Rejected: process decisions as log rows only

A log row cannot carry context and rejected alternatives, and rejected alternatives are half of "why does the system look like this". Expensive-to-reverse process rules would lose exactly the record that keeps them from being re-argued.

### Rejected: keeping the stretched ADR definition ("architecture, methodology mechanics, …")

The definition PDR-003 shipped. It makes the folder honest only to insiders who have read the local redefinition — everyone else reads ADR and expects architecture. A methodology built for adoption cannot afford private meanings for public terms.
