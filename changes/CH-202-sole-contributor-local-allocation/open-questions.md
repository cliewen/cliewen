---
id: CH-202-questions
type: open-questions
status: open
links: [CH-202]
title: Sole-contributor local allocation open questions
---

# Open questions

No blocking question. The maintainer chose to let the agent proceed on local allocation and state it in the brief, over keeping the stop and having `clue init` explain it.

Settled in review: an unreadable host is not a case. The signals are what plain Git can read from the remote, so a repository hosted where `gh` cannot reach is treated like any other; the only stop beyond another contributor's branch is a remote the agent cannot read at all, where it could not push the proposal either.
