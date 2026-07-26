---
id: CH-066-open-questions
type: open-questions
status: open
links: [CH-066]
title: Open questions for CH-066
---

# Open questions

None blocking. Two implementation unknowns are settled empirically by the goreleaser dry run rather than by a human decision, and are recorded here only so a reviewer can see they were not assumed:

1. Whether `formats: [binary]` carries the `.exe` extension into the uploaded Windows asset name. If it does not, the `bare` archive's `name_template` gains an explicit conditional suffix. Either way the published name is unchanged from today's.
2. Whether goreleaser's `SHA256SUMS` lines are GNU-coreutils-compatible, so an adopter's `sha256sum -c --ignore-missing SHA256SUMS` keeps passing.
