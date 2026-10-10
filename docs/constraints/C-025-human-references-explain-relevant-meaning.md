---
id: C-025
type: constraint
status: active
links: [G-026, PDR-071]
title: Human references explain relevant meaning without requiring an identity lookup
source: PDR-071
enforcement: agent
binds: adopter
---

# Explain references while writing

A reader must be able to understand the purpose, relevant effect and required action from the current text. Describe that meaning inline and use descriptive links for optional detail. Keep canonical identities secondary in human prose and unchanged in literal or machine-facing data. Changelogs omit internal rule identities.

The [readable-reference decision](../decisions/PDR-071-readable-references-explain-meaning-inline.md) establishes the policy. Its shipped instruction is `internal/skills/source/shared/readable-references.md.tmpl`; every standalone skill routes to it. Review judges comprehension. Structural formatter and carrier tests do not make this a machine-enforced readability score.
