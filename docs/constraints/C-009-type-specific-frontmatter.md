---
id: C-009
type: constraint
status: active
links: []
title: Type-specific frontmatter fields are present on their types
source: clue-verify checklist, docs/README.md
enforcement: agent
---

# C-009 — Type-specific frontmatter fields are present

Beyond the linted core (`id`, `type`, `status`, `links`, `title`), some types carry required extensions: decisions carry `author` and `accepted-by`, capabilities carry `goal`. Constraints' own `source` and `enforcement` were promoted to `machine` by the change that created this register; the rest are still agent-held.

**Promotion trigger:** `checkCoreFields` grows a per-type required-fields map covering `author`/`accepted-by`/`goal` — then `enforcement: machine`.
