---
id: OQ-082
type: open-questions
status: open
links: [CH-082, P-008]
title: Open questions for pricing adopter configuration
---

# CH-082 — open questions

None blocking.

## Non-blocking — recorded for the human, does not change the finding

**Was declining to vendor the `clue` binary a mandatory need or a preference?** CH-004 in the adopter rewrote the wall's install step to download the release asset instead of verifying a copy committed under `.github/tools/`. The commit message attributes the change to Cliewen becoming public, not to a constraint. It cannot be settled from the repository whether committing a ~10 MB binary into a product repository was unacceptable to the maintainer or merely undesirable once an alternative existed.

Recorded rather than inferred because the adopter's maintainer is the human who selected the adopter. It is non-blocking because AN-012's conclusion holds under either answer: a configuration key does not express *how the binary arrives* better than an edit does, and the candidate remedy — a wall that need not be forked — addresses the edit either way.

Expected source of blocking questions: the mandatory-versus-preference split. Tank Royale's maintainer is the human who selected it as the adopter, so a constraint that looks mandatory in the repository may be a preference, and vice versa. Any such case is recorded here for a human answer rather than inferred from repository activity ([`clue-analysis`](../../.agents/skills/clue-analysis/skill.md) step 2).
