---
id: CH-186-open-questions
type: open-questions
status: open
links: [CH-186]
title: Open questions — the CI wall asks the host, and its setup instructions ship with it
---

# CH-186 — Open questions

None blocking.

The proposal first assumed a fork pull request could not see the base repository's settings. That was wrong: a fork pull request's token is read-only but belongs to the base repository, and the rules endpoint needs only read access. What the workflow token genuinely cannot see is classic branch protection and a ruleset's bypass list, which need administration access. PDR-060 names those two as unseen in every report instead of treating forks as a special case.
