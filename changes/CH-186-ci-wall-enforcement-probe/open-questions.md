---
id: CH-186-open-questions
type: open-questions
status: open
links: [CH-186]
title: Open questions — the CI wall asks the host, and its setup instructions ship with it
---

# CH-186 — Open questions

None blocking.

One design call worth stating rather than asking, because PDR-059 already settled the shape it follows: on a fork-originated pull request, the `GITHUB_TOKEN` GitHub issues is scoped to the fork and cannot read the base repository's branch-protection or ruleset settings regardless of the permissions the workflow requests. The probe reports that case as unknown — the same answer PDR-059 gives a host it has no way to ask — rather than treating a same-repo pull request and a fork pull request differently in what they claim. This will be confirmed against a real fork PR while implementing AC-198, not decided in the abstract.
