---
id: CH-206-questions
type: questions
status: active
links: [CH-206]
title: Readable-reference decisions and observation boundary
---

# Questions

The owner chose the common rule and all human-facing surfaces, including CLI output and agent handoffs. No blocking implementation preference remains. Human comprehension is judged by a person and is not inferred from passing formatter or generator tests. The existing pilot's screenshot-method choice and Human observations remain pending outside this change; no observation is invented or accepted here.

## Source local preflight limitation

The accepted source tree contains managed Claude skill symlinks, while the current local command refuses all tracked symlinks. This affects native local preflight, independently of readable-reference behavior. A proposed follow-up goal records it in the digest; this change does not weaken materialization or claim a successful native local acceptance check. The owner must select an allowed integration path before acceptance if that limitation remains.
