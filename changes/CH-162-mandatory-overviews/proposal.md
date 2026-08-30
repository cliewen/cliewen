---
id: CH-162
type: change
status: open
links: [P-018, M-080, PDR-051, CAP-001, CAP-002, CAP-006]
title: Architecture and design overviews are mandatory and maintained
---

# CH-162 — Architecture and design overviews are mandatory and maintained

## Proposal

Serve P-018's M-080 by making every Cliewen corpus carry concise, maintained system architecture and cross-cutting design overviews at their canonical README paths.

`clue init` currently creates an architecture index but no design overview, while extraction keeps the system-level design obligation capability-local. Neither the validator nor lifecycle skills ensure that an agent discovers, writes, or keeps the system picture current. This change makes the paths structural validation requirements, creates explicit bootstrap state, and directs agents to draft the real overview from repository evidence before green validation. It adds a documentation-impact feedback loop to normal change work, retains capability `design.md` for local detail, and makes relevant ADRs and IDRs point to their affected overview.

Migration and upgrade remain conservative: they can safely add an unactivated canonical template, but they do not infer where existing overview prose belongs or move it. The lifecycle skill surveys candidate source documents, obtains one explicit grouped consent before relocating them and rewriting links, and creates a short canonical pointer if the move is declined. The corpus remains current rather than accumulating a change log; a document is added only when it answers a reader question that existing material cannot answer without duplication.

## Scope boundary

This change does not make a CLI model system intent, mandate a diagram where it does not clarify a relationship or flow, retroactively rewrite historical decisions, or turn every simple edit into a whole-corpus rewrite.
