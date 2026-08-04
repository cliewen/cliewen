---
id: C-009
type: constraint
status: active
links: []
title: Type-specific frontmatter fields are present on their types
source: clue-verify checklist, docs/README.md
enforcement: machine
---

# C-009 — Type-specific frontmatter fields are present

Beyond the linted core (`id`, `type`, `status`, `links`, `title`), some types carry required extensions: decisions carry `author` and `accepted-by`, capabilities carry `goal`. Constraints' own `source` and `enforcement` were the first fields the register linted; the rest joined them when `clue validate` gained a per-type required-field map.

**Checked by:** `clue validate` ([AC-094](../capabilities/CAP-002-validate/criteria.md)) — the per-type required-field map in `internal/corpus/conventions.go` covers `author` and `accepted-by` on decisions and `goal` on capabilities. An empty `accepted-by` list is a present field: an unsigned decision declares itself unsigned, and that declaration is what [ADR-029](../decisions/ADR-029-accepted-by-is-cliewen-approval-only.md) wants recorded. A type the validator does not recognize is checked against the core fields alone, so an adopter's own types stay theirs.
