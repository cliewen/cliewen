# Constraints

C-xxx: rules that bind every change — laws, licenses, security policies, organizational mandates, and the **convention register**: every repo rule that would otherwise live only in prose.

External constraints come from outside the project: you do not decide them, you comply with them (the decision *how* to comply is an ADR or PDR in `decisions/`). The register holds the rest: a convention that binds every change but would otherwise live only in prose (an AGENTS.md rule, a README convention) registers here as a constraint artifact, so the rules have an inventory a validator can count instead of prose that drifts silently. Rules the versioned skills carry need no registration — the register starts with exactly the conventions AGENTS.md declares that no skill carries.

Each constraint names its `source` (the document, law, or catalog that states the rule) and an `enforcement` class: `machine` (a lint or CI check holds it), `agent` (an agent must hold it until a machine check exists — each such constraint states its promotion trigger, and the count of agent-enforced constraints on `clue validate`'s OK line is the visible promotion backlog), or `human` (only review can hold it). Every change is assessed against the active constraints before its PR. This index is the register table:

<!-- clue:index:start -->
- [C-001 — Markdown prose is never hard-wrapped](C-001-no-hard-wrapped-markdown.md) · `agent`
<!-- clue:index:end -->
