---
id: CH-200-questions
type: open-questions
status: open
links: [CH-200]
title: Agent adapters open questions
---

# Open questions

No blocking question. The maintainer chose the tracked route and for the M-107 revision to ride with this change on 2026-10-06.

Where the Codex login comes from is decided here and is not blocking: a directory outside the repository, `~/.cliewen-trial/codex`, mounted into the container's `.codex`. It was created by `codex login --device-auth` run inside the container, and holds the maintainer's own login.
