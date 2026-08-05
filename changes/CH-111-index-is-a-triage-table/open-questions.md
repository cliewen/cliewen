---
id: CH-111-open-questions
type: open-questions
status: resolved
links: [CH-111-proposal]
title: Open questions for CH-111
---

# CH-111 — open questions

Five questions were put to the human and answered before implementation began. All five are recorded in ADR-046 and the digest rather than left open here.

1. **Is the index shape predefined, or detected per adopter?** Answered: predefined. The purpose of the index is triage, and an agent cannot triage from a column set that varies by folder. The five differing schemas observed in one adopter argue for a fixed shape rather than against it.
2. **Where does the description come from?** Answered: extracted from the artifact body, preferring a lede beneath the H1. No new frontmatter field.
3. **Where does the date come from?** Answered: there is no date column. Sourcing it from the last commit was measured inert, the creation date carries the same defect through renames, and no corpus-wide date field exists. Recency belongs in a `reviewed:` field decided in its own change.
4. **Does the block become a table?** Answered: no, and this reversed the change's original shape. ADR-041 and C-016 settled the row's shape one release ago, and the information a table would carry is the information this change adds to the row instead. The question was reopened once ADR-041 was read, because the first proposal was written against v0.12.0 and did not know the row already carried id, title, and status.
5. **Is the extracted sentence authoritative, or a seed?** Answered: a seed. Measurement against this corpus showed extraction is good for capability READMEs and systematically wrong for decisions, whose first paragraph states the problem in the present tense. The author edits it, regeneration never rewrites it, and nothing is backfilled into rows that already exist.

None open.
