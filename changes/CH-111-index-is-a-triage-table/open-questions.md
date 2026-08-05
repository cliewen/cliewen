---
id: CH-111-open-questions
type: open-questions
status: resolved
links: [CH-111-proposal]
title: Open questions for CH-111
---

# CH-111 — open questions

Four questions were put to the human before the proposal was written and answered there. All four are recorded as decisions in ADR-046 and the digest rather than left open here.

1. **Is the index shape predefined, or detected per adopter?** Answered: predefined. The purpose of the index is triage, and an agent cannot triage from a column set that varies by folder. The five differing schemas observed in one adopter are the evidence for a fixed prefix rather than against it.
2. **Where does the description come from?** Answered: extracted from the artifact body, with a lede beneath the H1 preferred and a reported fallback where none exists. No new frontmatter field.
3. **Is the schema closed, or a required prefix an adopter may extend?** Answered: a required prefix. The predefined columns come first and are regenerated; anything after them is preserved per row, so an adopter keeps local columns without weakening the guarantee.
4. **Where does the date come from?** Answered: there is no date column. Sourcing it from the last commit was measured inert, the creation date carries the same defect through renames, and no corpus-wide date field exists. Recency belongs in a `reviewed:` field decided in its own change. The reasoning is recorded in the proposal and in ADR-046.

None open.
