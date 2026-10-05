---
id: CH-198-questions
type: open-questions
status: open
links: [CH-198]
title: Scenario-trials plan open questions
---

# Open questions

Neither question blocks this change; both are answered by the maintainer in review.

1. **Where does the harness live, and in what language?** Recommendation: Go under `tools/`, because the repository is Go and the maintainer works on Windows, so one program avoids a bash and a PowerShell variant. It is not shipped to adopters. M-103 confirms or revises this when it starts.
2. **Does merging this change promote G-025 to `accepted` and P-024 to `active`?** This change writes them as `proposed` and `draft`, because promotion is a human act. If the maintainer wants them promoted in this change, they say so and the next commit does it before the pull request is marked ready.
