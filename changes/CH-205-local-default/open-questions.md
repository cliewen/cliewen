---
id: CH-205-questions
type: open-questions
status: resolved
links: [CH-205]
title: Local-default scope choices
---

# Resolved choices

The human chose local default for new repos, explicit CLI selection for PR, preservation of existing workflows at upgrade, and continued explicit PR policy in this source repository. Integration branch defaults to the existing `main` convention and remains configurable in policy. Explicit options may create a missing policy but never overwrite an existing one; policy changes still require acceptance under the previous workflow.
