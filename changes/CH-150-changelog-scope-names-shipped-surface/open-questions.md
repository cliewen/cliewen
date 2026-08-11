---
id: CH-150-open-questions
type: open-questions
status: open
links: [CH-150-proposal]
title: Open questions for CH-150
---

# Open questions

No blocking questions.

Not blocking, and out of this change's scope: the acceptance-brief gate in `.github/workflows/ci.yml` fires on the CI scope classifier's `full` output, which means "this diff needs the full check suite", while its failure message calls the pull request a "full Cliewen PR" — the tier sense of the word. The check has no tier signal, so a light change touching `docs/` is asked for a full change's artifact, and the adopter-facing `.github/workflows/clue-validation.yml` has the same shape. Whether every Cliewen pull request owes an acceptance brief or only a full-tier one does is a question about what the methodology requires, and answering it here would widen this change past the release-note scope rule. Recorded so it survives this workspace; it needs its own change and human scope.
