# Constraints

C-xxx: rules the system must not break — laws, licenses, security policies, organizational mandates.

Constraints come from outside the project: you do not decide them, you comply with them (the decision *how* to comply is an ADR or PDR in `decisions/`). Each constraint names its `source` (the document or catalog that states the rule) and an `enforcement` class: `machine` (a lint or CI check holds it), `agent` (an agent must hold it until a machine check exists — the visible promotion backlog), or `human` (only review can hold it). Every change is assessed against the active constraints before its PR.

<!-- clue:index:start -->
<!-- clue:index:end -->
