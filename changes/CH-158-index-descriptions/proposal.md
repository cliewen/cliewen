---
id: CH-158
type: change
status: open
links: [P-015, M-074, CAP-005]
title: Make every generated index row explain its artifact
---

# CH-158 — Make every generated index row explain its artifact

The generated taxonomy indexes still contain rows that identify an artifact by title and status but do not say what the artifact is about. A reader choosing an entry point must therefore open those artifacts merely to discover whether they are relevant, which leaves P-015's final read-cost milestone unfinished.

Curate every row named by `clue validate --index-rows` from the durable meaning of its target. Preserve existing relationship wording where a decision row states a live supersession or amendment, keep every row to one Markdown line, and do not alter the underlying artifacts or the index-generation contract.

Re-run the derived backlog until no undescribed row remains, then record M-074's evidence and close P-015. This is a corpus completion change: it changes no shipped `clue` behavior, generated skill, or adopter materialized by `clue init` or `clue scaffold`.
